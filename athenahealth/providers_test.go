package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetProvider(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/GetProvider.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	provider, err := athenaClient.GetProvider("1")

	assert.NotNil(provider)
	assert.Nil(err)
}
