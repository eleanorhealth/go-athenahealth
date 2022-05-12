package athenahealth

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_SearchMedications(t *testing.T) {
	assert := assert.New(t)

	searchVal := "ibuprofen"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b, _ := ioutil.ReadFile("./resources/SearchMedications.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	medications, err := athenaClient.SearchMedications(context.Background(), searchVal)

	assert.NotNil(medications)
	assert.NoError(err)
}
