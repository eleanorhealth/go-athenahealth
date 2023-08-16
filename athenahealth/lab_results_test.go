package athenahealth

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListLabResults(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	patientID := "123"
	departmentID := "456"
	startDate := time.Date(2011, 9, 22, 0, 0, 0, 0, time.UTC)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(startDate.Format("01/02/2006"), r.URL.Query().Get("startdate"))

		b, _ := os.ReadFile("./resources/ListLabResults.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListLabResults(ctx, patientID, departmentID, &ListLabResultsOptions{
		StartDate: &startDate,
	})
	assert.NoError(err)
	assert.Len(res.LabResults, 4)
	assert.Equal(res.Pagination.TotalCount, 4)
}

func TestHTTPClient_AddLabResultDocument(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	patientID := "123"
	departmentID := "456"

	h := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(err)
		assert.Equal(departmentID, r.Form.Get("departmentid"))
		assert.Equal(string(LabResultAttachmentTypeJPG), r.Form.Get("attachmenttype"))

		b, _ := os.ReadFile("./resources/AddLabResultDocument.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	b := bytes.NewReader([]byte(`test bytes`))
	res, err := athenaClient.AddLabResultDocumentReader(ctx, patientID, departmentID, &AddLabResultDocumentOptions{
		AttachmentContents: b,
		AttachmentType:     LabResultAttachmentTypeJPG,
	})

	assert.NoError(err)
	assert.Equal(res, 1083563)
}

func TestHTTPClient_ListChangedLabResults(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	showPortalOnly := false

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(strconv.FormatBool(showPortalOnly), r.URL.Query().Get("showportalonly"))

		b, _ := os.ReadFile("./resources/ListChangedLabResults.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListChangedLabResults(ctx, &ListChangedLabResultsOptions{
		ShowPortalOnly: &showPortalOnly,
	})
	labResults := res.ChangedLabResults

	assert.NoError(err)
	assert.Len(labResults, 3)
	assert.Equal(res.Pagination.TotalCount, 3)
}
