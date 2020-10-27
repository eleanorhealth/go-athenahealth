package athenahealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eleanorhealth/go-athenahealth/athenahealth/ratelimiter"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokencacher"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokenprovider"
)

const (
	// PreviewBaseURL is the base URL used to make API requests in the preview environment.
	PreviewBaseURL = "https://api.athenahealth.com/preview1/"

	// ProdBaseURL is the base URL used to make API requests in the production environment.
	ProdBaseURL = "https://api.athenahealth.com/v1/"

	// userAgent is the user agent that will be sent with every HTTP request.
	userAgent = "go-athenahealth/1.0"
)

var _ Client = (*HTTPClient)(nil)

type HTTPClient struct {
	httpClient *http.Client

	practiceID string
	key        string
	secret     string

	preview bool

	baseURL string

	tokenProvider TokenProvider
	tokenCacher   TokenCacher
	rateLimiter   RateLimiter

	requestLock sync.Mutex
}

var _ Client = &HTTPClient{}

// APIError represents an error response from the athenahealth API.
type APIError struct {
	Err                   error  `json:"-"`
	AthenaError           string `json:"error"`
	AthenaDetailedMessage string `json:"detailedmessage"`

	HTTPResponse *http.Response
}

func (a *APIError) Error() string {
	details := "no detailed message"
	if len(a.AthenaDetailedMessage) > 0 {
		details = a.AthenaDetailedMessage
	}

	return fmt.Sprintf("athenahealth API error: %s (%s)", a.AthenaError, details)
}

func (a *APIError) Unwrap() error {
	return a.Err
}

type PaginationOptions struct {
	Limit  int
	Offset int
}

type PaginationResult struct {
	NextOffset     int
	PreviousOffset int
	TotalCount     int
}

type PaginationResponse struct {
	Previous   string `json:"previous"`
	Next       string `json:"next"`
	TotalCount int    `json:"totalcount"`
}

func makePaginationResult(nextURL, previousURL string, totalCount int) *PaginationResult {
	var nextOffset, previousOffset int

	next, err := url.Parse(nextURL)
	if err == nil {
		nextOffset, _ = strconv.Atoi(next.Query().Get("offset"))
	}

	previous, err := url.Parse(previousURL)
	if err == nil {
		previousOffset, _ = strconv.Atoi(previous.Query().Get("offset"))
	}

	return &PaginationResult{
		NextOffset:     nextOffset,
		PreviousOffset: previousOffset,
		TotalCount:     totalCount,
	}
}

func NewHTTPClient(httpClient *http.Client, practiceID, key, secret string) *HTTPClient {
	preview := true

	c := &HTTPClient{
		httpClient: httpClient,

		practiceID: practiceID,
		key:        key,
		secret:     secret,

		preview: preview,

		tokenProvider: tokenprovider.NewDefault(httpClient, key, secret, preview),
		tokenCacher:   tokencacher.NewDefault(),
		rateLimiter:   ratelimiter.NewDefault(),
	}

	c.setBaseURL()

	return c
}

func (h *HTTPClient) setBaseURL() {
	if h.preview {
		h.baseURL = fmt.Sprintf("%s%s", PreviewBaseURL, h.practiceID)
	} else {
		h.baseURL = fmt.Sprintf("%s%s", ProdBaseURL, h.practiceID)
	}
}

func (h *HTTPClient) request(method, path string, body io.Reader, headers http.Header, out interface{}) (*http.Response, error) {
	var token string
	var err error
	var expiresAt time.Time

	h.requestLock.Lock()

	retryAfter, err := h.rateLimiter.Allowed(h.preview)
	if err != nil {
		h.requestLock.Unlock()

		if errors.Is(err, ratelimiter.ErrRateExceeded) {
			time.Sleep(retryAfter)
			return h.request(method, path, body, headers, out)
		}

		return nil, err
	}

	token, err = h.tokenCacher.Get()
	if err != nil {
		if !errors.Is(err, tokencacher.ErrTokenNotExist) && !errors.Is(err, tokencacher.ErrTokenExpired) {
			h.requestLock.Unlock()
			return nil, err
		}

		token, expiresAt, err = h.tokenProvider.Provide()
		if err != nil {
			h.requestLock.Unlock()
			return nil, err
		}

		// Remove 1 minute from the expiration time to create a buffer to see
		// if it resolves intermittent 401s.
		err = h.tokenCacher.Set(token, expiresAt.Add(-1*time.Minute))
		if err != nil {
			h.requestLock.Unlock()
			return nil, err
		}
	}

	h.requestLock.Unlock()

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	reqURL := fmt.Sprintf("%s%s", h.baseURL, path)

	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		req.Header = headers
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("User-Agent", userAgent)

	res, err := h.httpClient.Do(req)
	if err != nil {
		return res, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	res.Body.Close()

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	// 200 OK
	// 300 Multiple Choices
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		err := &APIError{}
		if res.StatusCode == http.StatusNotFound {
			err.Err = ErrNotFound
		}

		json.Unmarshal(resBody, err)

		err.HTTPResponse = res

		return res, err
	}

	if out != nil {
		err = json.Unmarshal(resBody, out)
		if err != nil {
			return res, fmt.Errorf("Error unmarshaling response body: %s", err)
		}
	}

	return res, nil
}

func (h *HTTPClient) WithPreview(preview bool) *HTTPClient {
	h.preview = preview
	h.setBaseURL()

	if _, ok := h.tokenProvider.(*tokenprovider.Default); ok {
		h.tokenProvider = tokenprovider.NewDefault(h.httpClient, h.key, h.secret, preview)
	}

	return h
}

func (h *HTTPClient) WithTokenProvider(provider TokenProvider) *HTTPClient {
	h.tokenProvider = provider

	return h
}

func (h *HTTPClient) WithTokenCacher(cacher TokenCacher) *HTTPClient {
	h.tokenCacher = cacher

	return h
}

func (h *HTTPClient) WithRateLimiter(rateLimiter RateLimiter) *HTTPClient {
	h.rateLimiter = rateLimiter

	return h
}

func (h *HTTPClient) Get(path string, query url.Values, out interface{}) (*http.Response, error) {
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	return h.request("GET", path, nil, nil, out)
}

func (h *HTTPClient) Post(path string, body io.Reader, out interface{}) (*http.Response, error) {
	return h.request("POST", path, body, nil, out)
}

func (h *HTTPClient) PostForm(path string, v url.Values, out interface{}) (*http.Response, error) {
	var body io.Reader
	var headers = http.Header{}

	if v != nil {
		body = strings.NewReader(v.Encode())
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return h.request("POST", path, body, headers, out)
}

func (h *HTTPClient) Put(path string, body io.Reader, out interface{}) (*http.Response, error) {
	return h.request("PUT", path, body, nil, out)
}

func (h *HTTPClient) PutForm(path string, v url.Values, out interface{}) (*http.Response, error) {
	var body io.Reader
	var headers = http.Header{}

	if v != nil {
		body = strings.NewReader(v.Encode())
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return h.request("PUT", path, body, headers, out)
}

func (h *HTTPClient) Delete(path string, body io.Reader, out interface{}) (*http.Response, error) {
	return h.request("DELETE", path, body, nil, out)
}

func (h *HTTPClient) DeleteForm(path string, v url.Values, out interface{}) (*http.Response, error) {
	var body io.Reader
	var headers = http.Header{}

	if v != nil {
		body = strings.NewReader(v.Encode())
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return h.request("DELETE", path, body, headers, out)
}
