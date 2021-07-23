package tokenprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// PreviewAuthURL is the URL used to authenticate in the preview environment.
	PreviewAuthURL = "https://api.preview.platform.athenahealth.com/oauth2/v1/token"

	// ProdAuthURL is the URL used to authenticate in the production environment.
	ProdAuthURL = "https://api.platform.athenahealth.com/oauth2/v1/token"
)

type Default struct {
	httpClient *http.Client

	clientID string
	secret   string

	authURL string
}

func NewDefault(httpClient *http.Client, clientID, secret string, preview bool) *Default {
	d := &Default{
		httpClient: httpClient,

		clientID: clientID,
		secret:   secret,
	}

	if preview {
		d.authURL = PreviewAuthURL
	} else {
		d.authURL = ProdAuthURL
	}

	return d
}

type authResponse struct {
	AccessToken string      `json:"access_token"`
	ExpiresIn   json.Number `json:"expires_in"`
}

func (d *Default) Provide(ctx context.Context) (string, time.Time, error) {
	vals := url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"athena/service/Athenanet.MDP.*"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.authURL, bytes.NewBufferString(vals.Encode()))
	if err != nil {
		return "", time.Now(), err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.clientID, d.secret)

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

	expiresIn, err := authRes.ExpiresIn.Int64()
	if err != nil {
		return "", time.Now(), err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(expiresIn))

	return authRes.AccessToken, expiresAt, nil
}
