package athenahealth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/eleanorhealth/go-athenahealth/athenahealth/ratelimiter"
	"github.com/stretchr/testify/assert"
)

const testPracticeID = "123456"
const testAPIKey = "api-key"
const testAPISecret = "api-secret"
const testToken = "token"

func testClient(h http.HandlerFunc) (*HTTPClient, *httptest.Server) {
	if h == nil {
		h = func(w http.ResponseWriter, r *http.Request) {
			b, _ := json.Marshal(nil)
			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		}
	}

	ts := httptest.NewServer(h)

	athenaClient := NewHTTPClient(ts.Client(), testPracticeID, testAPIKey, testAPISecret).
		WithTokenProvider(&testTokenProvider{}).
		WithTokenCacher(&testTokenCacher{})

	athenaClient.baseURL = ts.URL

	return athenaClient, ts
}

type testTokenProvider struct {
}

func (t *testTokenProvider) Provide(ctx context.Context) (string, time.Time, error) {
	return testToken, time.Now().Add(time.Minute * 1), nil
}

type testTokenCacher struct {
}

func (t *testTokenCacher) Get(ctx context.Context) (string, error) {
	return testToken, nil
}

func (t *testTokenCacher) Set(context.Context, string, time.Time) error {
	return nil
}

type testRateLimiter struct {
	AllowedFunc func(preview bool) (time.Duration, error)
}

func (t *testRateLimiter) Allowed(ctx context.Context, preview bool) (time.Duration, error) {
	if t.AllowedFunc != nil {
		return t.AllowedFunc(preview)
	}

	return 0, nil
}

type testStats struct {
	RequestFunc         func(method, path string) error
	ResponseSuccessFunc func() error
	ResponseErrorFunc   func() error
}

func (t *testStats) Request(method, path string) error {
	if t.RequestFunc != nil {
		return t.RequestFunc(method, path)
	}

	return nil
}

func (t *testStats) ResponseSuccess() error {
	if t.ResponseSuccessFunc != nil {
		return t.ResponseSuccessFunc()
	}

	return nil
}

func (t *testStats) ResponseError() error {
	if t.ResponseErrorFunc != nil {
		return t.ResponseErrorFunc()
	}

	return nil
}

func TestNewHTTPClient(t *testing.T) {
	assert := assert.New(t)

	practiceID := "123456"
	key := "key"
	secret := "secret"

	athenaClient := NewHTTPClient(&http.Client{}, practiceID, key, secret)

	assert.Equal(practiceID, athenaClient.practiceID)
	assert.Equal(secret, athenaClient.secret)
	assert.Equal(key, athenaClient.clientID)

	// Preview mode should default to true.
	assert.True(athenaClient.preview)

	assert.NotNil(athenaClient.tokenProvider)
	assert.NotNil(athenaClient.tokenCacher)

	assert.NotEmpty(athenaClient.baseURL)
}

func TestAPIError_Error(t *testing.T) {
	assert := assert.New(t)

	err := &APIError{
		Err:                   ErrNotFound,
		AthenaError:           "unknown error",
		AthenaDetailedMessage: "something went wrong",
	}

	assert.True(errors.Is(err.Err, ErrNotFound))
	assert.Contains(err.Error(), err.AthenaError)
}

func TestAPIError_setBaseURL(t *testing.T) {
	assert := assert.New(t)

	practiceID := "123456"

	athenaClient := NewHTTPClient(&http.Client{}, practiceID, "", "")

	var expectedBaseURL string

	// Preview base URL
	expectedBaseURL = fmt.Sprintf("%s%s", PreviewBaseURL, practiceID)
	assert.Equal(expectedBaseURL, athenaClient.baseURL)

	// Production base URL
	athenaClient.preview = false
	athenaClient.setBaseURL()
	expectedBaseURL = fmt.Sprintf("%s%s", ProdBaseURL, practiceID)
	assert.Equal(expectedBaseURL, athenaClient.baseURL)
}

func TestHTTPClient_request(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(fmt.Sprintf("Bearer %s", testToken), r.Header.Get("Authorization"))
		assert.Equal(userAgent, r.UserAgent())

		w.Write([]byte(`{"msg":"Hello World!"}`))
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	var out map[string]string
	res, err := athenaClient.request(context.Background(), "GET", "/", nil, nil, &out)

	assert.NotNil(res)
	assert.NoError(err)
	assert.Equal("Hello World!", out["msg"])
}

func TestHTTPClient_request_error(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.request(context.Background(), "GET", "/", nil, nil, nil)

	assert.NotNil(res)
	assert.NotNil(err)
	assert.IsType(&APIError{}, err)
}

func TestHTTPClient_request_default_context_timeout(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(time.Millisecond * 100):
			w.WriteHeader(http.StatusOK)
		}
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()
	athenaClient.WithRequestTimeout(time.Millisecond * 10)

	res, err := athenaClient.request(context.Background(), "GET", "/", nil, nil, nil)

	assert.Nil(res)
	assert.NotNil(err)
	assert.ErrorIs(err, context.DeadlineExceeded)
}

func TestHTTPClient_request_override_default_context_timeout(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(time.Millisecond * 100):
			w.WriteHeader(http.StatusOK)
		}
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()
	athenaClient.WithRequestTimeout(time.Millisecond * 10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*110)
	defer cancel()

	res, err := athenaClient.request(ctx, "GET", "/", nil, nil, nil)

	assert.NotNil(res)
	assert.Nil(err)
}

func TestHTTPClient_rate_limit(t *testing.T) {
	assert := assert.New(t)

	rateLimiter := &testRateLimiter{}

	rateLimited := false
	called := false
	rateLimiter.AllowedFunc = func(preview bool) (time.Duration, error) {
		if rateLimited {
			called = true
			return 0, nil
		}

		rateLimited = true

		return 100 * time.Millisecond, ratelimiter.ErrRateExceeded
	}

	athenaClient, ts := testClient(nil)
	athenaClient.WithRateLimiter(rateLimiter)

	defer ts.Close()

	var out map[string]string
	res, err := athenaClient.request(context.Background(), "GET", "/", nil, nil, &out)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_WithPreview(t *testing.T) {
	assert := assert.New(t)

	athenaClient := NewHTTPClient(&http.Client{}, "", "", "")

	athenaClient.WithPreview(false)

	assert.False(athenaClient.preview)
}

func TestHTTPClient_WithTokenProvider(t *testing.T) {
	assert := assert.New(t)

	athenaClient := NewHTTPClient(&http.Client{}, "", "", "")

	tokenProvider := &testTokenProvider{}
	athenaClient.WithTokenProvider(tokenProvider)

	assert.Equal(tokenProvider, athenaClient.tokenProvider)
}

func TestHTTPClient_WithTokenCacher(t *testing.T) {
	assert := assert.New(t)

	athenaClient := NewHTTPClient(&http.Client{}, "", "", "")

	tokenCacher := &testTokenCacher{}
	athenaClient.WithTokenCacher(tokenCacher)

	assert.Equal(tokenCacher, athenaClient.tokenCacher)
}

func TestHTTPClient_WithRateLimiter(t *testing.T) {
	assert := assert.New(t)

	athenaClient := NewHTTPClient(&http.Client{}, "", "", "")

	rateLimiter := &testRateLimiter{}
	athenaClient.WithRateLimiter(rateLimiter)

	assert.Equal(rateLimiter, athenaClient.rateLimiter)
}

func TestHTTPClient_WithStats(t *testing.T) {
	assert := assert.New(t)

	athenaClient := NewHTTPClient(&http.Client{}, "", "", "")

	stats := &testStats{}
	athenaClient.WithStats(stats)

	assert.Equal(stats, athenaClient.stats)
}

func TestHTTPClient_WithRequestTimeout(t *testing.T) {
	assert := assert.New(t)

	athenaClient := NewHTTPClient(&http.Client{}, "", "", "")
	assert.Equal(defaultRequestTimeout, athenaClient.requestTimeout)

	requestTimeout := time.Minute
	athenaClient.WithRequestTimeout(requestTimeout)

	assert.Equal(requestTimeout, athenaClient.requestTimeout)
}

func TestHTTPClient_Get(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("bar", r.URL.Query().Get("foo"))

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	var query = url.Values{}
	query.Add("foo", "bar")

	res, err := athenaClient.Get(context.Background(), "/", query, nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_Post(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.Post(context.Background(), "/", strings.NewReader("foo"), nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_PostForm(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		assert.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	var values = url.Values{}
	values.Add("foo", "bar")

	res, err := athenaClient.PostForm(context.Background(), "/", values, nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_PostFormReader(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		assert.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	values := NewFormURLEncoder()
	values.AddString("foo", "bar")

	res, err := athenaClient.PostFormReader(context.Background(), "/", values, nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_Put(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.Put(context.Background(), "/", strings.NewReader("foo"), nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_PutForm(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		assert.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	var values = url.Values{}
	values.Add("foo", "bar")

	res, err := athenaClient.PutForm(context.Background(), "/", values, nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_Delete(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	athenaClient.baseURL = ts.URL

	res, err := athenaClient.Delete(context.Background(), "/", strings.NewReader("foo"), nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_DeleteForm(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		assert.True(len(b) > 0)

		assert.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	var values = url.Values{}
	values.Add("foo", "bar")

	res, err := athenaClient.DeleteForm(context.Background(), "/", values, nil)

	assert.NotNil(res)
	assert.NoError(err)
	assert.True(called)
}

func Test_makePaginationResult(t *testing.T) {
	assert := assert.New(t)

	nextURL := "/booked?DEPARTMENTID=1&ENDDATE=12%2F31%2F2020&STARTDATE=01%2F01%2F2020&limit=10&offset=30"
	previousURL := "/booked?DEPARTMENTID=1&ENDDATE=12%2F31%2F2020&STARTDATE=01%2F01%2F2020&limit=10&offset=10"
	totalCount := 50

	p := makePaginationResult(nextURL, previousURL, totalCount)

	assert.Equal(30, p.NextOffset)
	assert.Equal(10, p.PreviousOffset)
	assert.Equal(totalCount, p.TotalCount)
}
