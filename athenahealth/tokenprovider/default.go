package tokenprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// PreviewAuthURL is the URL used to authenticate in the preview environment.
	PreviewAuthURL = "https://api.athenahealth.com/oauthpreview/token"

	// ProdAuthURL is the URL used to authenticate in the production environment.
	ProdAuthURL = "https://api.athenahealth.com/oauth/token"
)

type Default struct {
	httpClient *http.Client

	key    string
	secret string

	authURL string
}

func NewDefault(httpClient *http.Client, key, secret string, preview bool) *Default {
	d := &Default{
		httpClient: httpClient,

		key:    key,
		secret: secret,
	}

	if preview {
		d.authURL = PreviewAuthURL
	} else {
		d.authURL = ProdAuthURL
	}

	return d
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (d *Default) Provide() (string, time.Time, error) {
	vals := url.Values{
		"grant_type": {"client_credentials"},
	}

	req, err := http.NewRequest("POST", d.authURL, bytes.NewBufferString(vals.Encode()))
	if err != nil {
		return "", time.Now(), err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.key, d.secret)

	res, err := d.httpClient.Do(req)
	if err != nil {
		return "", time.Now(), err
	}

	if res.StatusCode != http.StatusOK {
		return "", time.Now(), fmt.Errorf("%s", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", time.Now(), err
	}
	res.Body.Close()

	authRes := &authResponse{}
	err = json.Unmarshal(b, authRes)
	if err != nil {
		return "", time.Now(), err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(authRes.ExpiresIn))

	return authRes.AccessToken, expiresAt, nil
}
