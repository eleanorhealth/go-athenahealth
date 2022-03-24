package athenahealth

import (
	"bytes"
	"context"
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
	"github.com/eleanorhealth/go-athenahealth/athenahealth/stats"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokencacher"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokenprovider"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tracer"
	"github.com/rs/zerolog"
)

const (
	// PreviewBaseURL is the base URL used to make API requests in the preview environment.
	PreviewBaseURL = "https://api.preview.platform.athenahealth.com/v1/"

	// ProdBaseURL is the base URL used to make API requests in the production environment.
	ProdBaseURL = "https://api.platform.athenahealth.com/v1/"

	// userAgent is the user agent that will be sent with every HTTP request.
	userAgent = "go-athenahealth/1.0"

	// defaultRequestTimeout defines the HTTP request's context deadline if one is not specified by the caller.
	defaultRequestTimeout = 15 * time.Second
)

var _ Client = (*HTTPClient)(nil)

type HTTPClient struct {
	httpClient *http.Client

	practiceID string
	clientID   string
	secret     string

	preview bool

	baseURL string

	tokenProvider TokenProvider
	tokenCacher   TokenCacher
	rateLimiter   RateLimiter
	stats         Stats
	tracer        Tracer
	logger        *zerolog.Logger

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

	var status string
	if a.HTTPResponse != nil {
		status = a.HTTPResponse.Status
	}

	return fmt.Sprintf("athenahealth API error (%s): %s (%s)", status, a.AthenaError, details)
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

func NewHTTPClient(httpClient *http.Client, practiceID, clientID, secret string) *HTTPClient {
	preview := true

	c := &HTTPClient{
		httpClient: httpClient,

		practiceID: practiceID,
		clientID:   clientID,
		secret:     secret,

		preview: preview,

		tokenProvider: tokenprovider.NewDefault(httpClient, clientID, secret, preview),
		tokenCacher:   tokencacher.NewDefault(),
		rateLimiter:   ratelimiter.NewDefault(),
		stats:         stats.NewDefault(),
		tracer:        tracer.NewDefault(),
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

func (h *HTTPClient) request(ctx context.Context, method, path string, body io.Reader, headers http.Header, out interface{}) (*http.Response, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultRequestTimeout)
		defer cancel()
	}

	var token string
	var err error
	var expiresAt time.Time

	h.requestLock.Lock()

	retryAfter, err := h.rateLimiter.Allowed(ctx, h.preview)
	if err != nil {
		h.requestLock.Unlock()

		if errors.Is(err, ratelimiter.ErrRateExceeded) {
			time.Sleep(retryAfter)
			return h.request(ctx, method, path, body, headers, out)
		}

		return nil, err
	}

	token, err = h.tokenCacher.Get(ctx)
	if err != nil {
		if !errors.Is(err, tokencacher.ErrTokenNotExist) && !errors.Is(err, tokencacher.ErrTokenExpired) {
			h.requestLock.Unlock()
			return nil, err
		}

		token, expiresAt, err = h.tokenProvider.Provide(ctx)
		if err != nil {
			h.requestLock.Unlock()
			return nil, err
		}

		// Remove 1 minute from the expiration time to create a buffer to see
		// if it resolves intermittent 401s.
		err = h.tokenCacher.Set(context.Background(), token, expiresAt.Add(-1*time.Minute))
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

	ctx = h.tracer.Before(ctx, method, path)
	defer h.tracer.After(ctx)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		req.Header = headers
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("User-Agent", userAgent)

	if h.logger != nil {
		h.logger.Info().
			Str("method", method).
			Str("url", reqURL).
			Msg("athenahealth API request")
	}

	res, err := h.httpClient.Do(req)
	if err != nil {
		return res, err
	}

	err = h.stats.Request(method, path)
	if err != nil {
		return res, err
	}

	responseError := res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices
	if responseError {
		err = h.stats.ResponseError()
		if err != nil {
			return res, err
		}
	} else {
		err = h.stats.ResponseSuccess()
		if err != nil {
			return res, err
		}
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	res.Body.Close()

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	if h.logger != nil {
		h.logger.Info().
			Str("method", method).
			Str("url", reqURL).
			Int("statusCode", res.StatusCode).
			Int("bodyLength", len(resBody)).
			Msg("athenahealth API response")
	}

	if responseError {
		err := &APIError{}
		if res.StatusCode == http.StatusNotFound {
			err.Err = ErrNotFound
		}

		//nolint
		json.Unmarshal(resBody, err)

		err.HTTPResponse = res

		if h.logger != nil {
			h.logger.Info().
				Str("athenaError", err.AthenaError).
				Str("athenaDetailedMessage", err.AthenaDetailedMessage).
				Msg("athenahealth API error")
		}
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

func (h *HTTPClient) WithLogger(logger *zerolog.Logger) *HTTPClient {
	h.logger = logger

	return h
}

func (h *HTTPClient) WithPreview(preview bool) *HTTPClient {
	h.preview = preview
	h.setBaseURL()

	if _, ok := h.tokenProvider.(*tokenprovider.Default); ok {
		h.tokenProvider = tokenprovider.NewDefault(h.httpClient, h.clientID, h.secret, preview)
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

func (h *HTTPClient) WithStats(stats Stats) *HTTPClient {
	h.stats = stats

	return h
}

func (h *HTTPClient) WithTracer(tracer Tracer) *HTTPClient {
	h.tracer = tracer

	return h
}

func (h *HTTPClient) WithDatadogTracer(opts ...tracer.Option) *HTTPClient {
	h.tracer = tracer.NewDatadog(opts...)

	return h
}

func (h *HTTPClient) Get(ctx context.Context, path string, query url.Values, out interface{}) (*http.Response, error) {
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	return h.request(ctx, http.MethodGet, path, nil, nil, out)
}

func (h *HTTPClient) Post(ctx context.Context, path string, body io.Reader, out interface{}) (*http.Response, error) {
	return h.request(ctx, http.MethodPost, path, body, nil, out)
}

func (h *HTTPClient) PostForm(ctx context.Context, path string, v url.Values, out interface{}) (*http.Response, error) {
	var body io.Reader
	var headers = http.Header{}

	if v != nil {
		body = strings.NewReader(v.Encode())
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return h.request(ctx, http.MethodPost, path, body, headers, out)
}

func (h *HTTPClient) Put(ctx context.Context, path string, body io.Reader, out interface{}) (*http.Response, error) {
	return h.request(ctx, http.MethodPut, path, body, nil, out)
}

func (h *HTTPClient) PutForm(ctx context.Context, path string, v url.Values, out interface{}) (*http.Response, error) {
	var body io.Reader
	var headers = http.Header{}

	if v != nil {
		body = strings.NewReader(v.Encode())
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return h.request(ctx, http.MethodPut, path, body, headers, out)
}

func (h *HTTPClient) Delete(ctx context.Context, path string, body io.Reader, out interface{}) (*http.Response, error) {
	return h.request(ctx, http.MethodDelete, path, body, nil, out)
}

func (h *HTTPClient) DeleteForm(ctx context.Context, path string, v url.Values, out interface{}) (*http.Response, error) {
	var body io.Reader
	var headers = http.Header{}

	if v != nil {
		body = strings.NewReader(v.Encode())
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return h.request(ctx, http.MethodDelete, path, body, headers, out)
}
