package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_SearchAllergies(t *testing.T) {
	assert := assert.New(t)

	searchVal := "latex"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b, _ := os.ReadFile("./resources/SearchAllergies.json")
		w.Write(b)
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	allergies, err := athenaClient.SearchAllergies(context.Background(), searchVal)

	assert.NotNil(allergies)
	assert.Len(allergies, 2)
	assert.NoError(err)
}

func TestHTTPClient_SearchAllergies_None(t *testing.T) {
	assert := assert.New(t)

	searchVal := "sdfsadss"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b := []byte("[]")
		w.Write(b)
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	allergies, err := athenaClient.SearchAllergies(context.Background(), searchVal)

	assert.NotNil(allergies)
	assert.Len(allergies, 0)
	assert.NoError(err)
}
