package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetDepartment(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/GetDepartment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	department, err := athenaClient.GetDepartment("1")

	assert.NotNil(department)
	assert.Nil(err)
}
