package athenahealth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
)

const (
	testPracticeID = "123456"
	testAPIKey     = "api-key"
	testAPISecret  = "api-secret"
)

func TestClient(h http.HandlerFunc) (*HTTPClient, *httptest.Server) {
	if h == nil {
		h = func(w http.ResponseWriter, r *http.Request) {
			b, _ := json.Marshal(nil)
			w.Header().Add("Content-Type", "application/json")
			_, _ = w.Write(b)
		}
	}

	ts := httptest.NewServer(h)

	athenaClient := NewHTTPClient(ts.Client(), testPracticeID, testAPIKey, testAPISecret).
		WithTokenProvider(&testTokenProvider{}).
		WithTokenCacher(&testTokenCacher{})

	athenaClient.baseURL = ts.URL

	return athenaClient, ts
}

const testToken = "token"

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

type testRateLimiter struct { //nolint:unused
	AllowedFunc func(preview bool) (time.Duration, error)
}

func (t *testRateLimiter) Allowed(ctx context.Context, preview bool) (time.Duration, error) { //nolint:unused
	if t.AllowedFunc != nil {
		return t.AllowedFunc(preview)
	}

	return 0, nil
}

type testStats struct { //nolint: unused
	RequestFunc         func(method, path string) error
	ResponseSuccessFunc func() error
	ResponseErrorFunc   func() error
}

func (t *testStats) Request(method, path string) error { // nolint: unused
	if t.RequestFunc != nil {
		return t.RequestFunc(method, path)
	}

	return nil
}

func (t *testStats) ResponseSuccess() error { // nolint: unused
	if t.ResponseSuccessFunc != nil {
		return t.ResponseSuccessFunc()
	}

	return nil
}

func (t *testStats) ResponseError() error { // nolint: unused
	if t.ResponseErrorFunc != nil {
		return t.ResponseErrorFunc()
	}

	return nil
}
