package athenahealth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testPracticeID = "123456"
const testAPIKey = "api-key"
const testAPISecret = "api-secret"
const testToken = "token"

func testClient(h http.HandlerFunc) (*HTTPClient, *httptest.Server) {
	if h == nil {
		h = func(w http.ResponseWriter, r *http.Request) {
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

func (t *testTokenProvider) Provide() (string, time.Time, error) {
	return testToken, time.Now().Add(time.Minute * 1), nil
}

type testTokenCacher struct {
}

func (t *testTokenCacher) Get() (string, error) {
	return testToken, nil
}

func (t *testTokenCacher) Set(string, time.Time) error {
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
	assert.Equal(key, athenaClient.key)

	// Preview mode should default to true.
	assert.True(athenaClient.preview)

	assert.NotNil(athenaClient.tokenProvider)
	assert.NotNil(athenaClient.tokenCacher)

	assert.NotEmpty(athenaClient.baseURL)
}

func TestAPIError_Error(t *testing.T) {
	assert := assert.New(t)

	err := &APIError{
		Err:             "unknown error",
		DetailedMessage: "something went wrong",
	}

	assert.Contains(err.Error(), err.Err)
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
	res, err := athenaClient.request("GET", "/", nil, &out)

	assert.NotNil(res)
	assert.Nil(err)
	assert.Equal("Hello World!", out["msg"])
}

func TestHTTPClient_request_error(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.request("GET", "/", nil, nil)

	assert.NotNil(res)
	assert.NotNil(err)
	assert.IsType(&APIError{}, err)
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

	res, err := athenaClient.Get("/", query, nil)

	assert.NotNil(res)
	assert.Nil(err)
	assert.True(called)
}

func TestHTTPClient_Post(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.Post("/", nil, nil)

	assert.NotNil(res)
	assert.Nil(err)
	assert.True(called)
}

func TestHTTPClient_Put(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.Put("/", nil, nil)

	assert.NotNil(res)
	assert.Nil(err)
	assert.True(called)
}

func TestHTTPClient_Delete(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	athenaClient.baseURL = ts.URL

	res, err := athenaClient.Delete("/", nil)

	assert.NotNil(res)
	assert.Nil(err)
	assert.True(called)
}
