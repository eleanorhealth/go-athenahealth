package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListCustomFields(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/ListCustomFields.json")
		w.Write(b)
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	customFields, err := athenaClient.ListCustomFields(context.Background())

	assert.Len(customFields, 2)
	assert.NoError(err)
}
