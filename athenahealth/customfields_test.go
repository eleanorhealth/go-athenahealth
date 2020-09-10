package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListCustomFields(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/ListCustomFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	customFields, err := athenaClient.ListCustomFields()

	assert.Len(customFields, 2)
	assert.NoError(err)
}
