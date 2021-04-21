package tokenprovider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	assert := assert.New(t)

	key := "api-key"
	secret := "api-secret"
	preview := true

	p := NewDefault(&http.Client{}, key, secret, preview)

	assert.NotNil(p.httpClient)
	assert.Equal(key, p.clientID)
	assert.Equal(secret, p.secret)
	assert.Equal(PreviewAuthURL, p.authURL)

	preview = false
	p = NewDefault(&http.Client{}, "", "", preview)
	assert.Equal(ProdAuthURL, p.authURL)
}

func TestDefault_Provide(t *testing.T) {
	assert := assert.New(t)

	authRes := &authResponse{
		AccessToken: "foo",
		ExpiresIn:   "60",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal(r.FormValue("grant_type"), "client_credentials")

		b, _ := json.Marshal(authRes)
		w.Write(b)
	}))
	defer ts.Close()

	p := NewDefault(ts.Client(), "", "", false)
	p.authURL = ts.URL

	token, expiresAt, err := p.Provide()

	assert.Equal(authRes.AccessToken, token)
	assert.True(expiresAt.After(time.Now()))
	assert.NoError(err)
}
