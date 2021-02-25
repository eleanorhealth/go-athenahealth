package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListAdminDocuments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(r.URL.Path, "/patients/123/")
		assert.Equal("3", r.URL.Query().Get("departmentid"))

		b, _ := ioutil.ReadFile("./resources/ListAdminDocuments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListAdminDocumentsOptions{
		DepartmentID: "3",
	}

	res, err := athenaClient.ListAdminDocuments("123", opts)

	assert.Len(res.AdminDocuments, 1)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 1)
	assert.NoError(err)
}
