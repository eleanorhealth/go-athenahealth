package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_SearchMedications(t *testing.T) {
	assert := assert.New(t)

	searchVal := "ibuprofen"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b, _ := os.ReadFile("./resources/SearchMedications.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	medications, err := athenaClient.SearchMedications(context.Background(), searchVal)

	assert.NotNil(medications)
	assert.Len(medications, 2)
	assert.NoError(err)
}

func TestHTTPClient_SearchMedications_None(t *testing.T) {
	assert := assert.New(t)

	searchVal := "sdfsadss"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b := []byte("[]")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	meds, err := athenaClient.SearchAllergies(context.Background(), searchVal)

	assert.NotNil(meds)
	assert.Len(meds, 0)
	assert.NoError(err)
}
