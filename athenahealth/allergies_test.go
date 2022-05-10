package athenahealth

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_SearchAllergies(t *testing.T) {
	assert := assert.New(t)

	searchVal := "latex"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b, _ := ioutil.ReadFile("./resources/SearchAllergies.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	allergies, err := athenaClient.SearchAllergies(context.Background(), searchVal)

	assert.NotNil(allergies)
	assert.NoError(err)
}
