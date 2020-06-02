package athenahealth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokencacher"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokenprovider"
)

const (
	// PreviewBaseURL is the base URL used to make API requests in the preview environment.
	PreviewBaseURL = "https://api.athenahealth.com/preview1/"

	// ProdBaseURL is the base URL used to make API requests in the production environment.
	ProdBaseURL = "https://api.athenahealth.com/v1/"
)

type HTTPClient struct {
	httpClient *http.Client

	practiceID string
	key        string
	secret     string

	preview bool

	baseURL string

	tokenProvider TokenProvider
	tokenCacher   TokenCacher

	tokenLock sync.Mutex
}

var _ Client = &HTTPClient{}

// APIError represents an error response from the athenahealth API.
type APIError struct {
	Err             string `json:"error"`
	DetailedMessage string `json:"detailedmessage"`

	HTTPResponse *http.Response
}

func (a *APIError) Error() string {
	return fmt.Sprintf("athenahealth API error: %s", a.Err)
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

	h.tokenLock.Lock()

	token, err = h.tokenCacher.Get()
	if err != nil {
		if !errors.Is(err, tokencacher.ErrTokenNotExist) && !errors.Is(err, tokencacher.ErrTokenExpired) {
			return nil, err
		}

		token, expiresAt, err = h.tokenProvider.Provide()
		if err != nil {
			return nil, err
		}

		err = h.tokenCacher.Set(token, expiresAt)
		if err != nil {
			return nil, err
		}
	}

	h.tokenLock.Unlock()

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

	res, err := h.httpClient.Do(req)
	if err != nil {
		return res, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	res.Body.Close()

	// 200 OK
	// 300 Multiple Choices
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		err := &APIError{}
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
