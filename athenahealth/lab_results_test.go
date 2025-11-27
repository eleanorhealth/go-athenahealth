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

func TestHTTPClient_ListLabResults_required_params_empty(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	startDate := time.Date(2011, 9, 22, 0, 0, 0, 0, time.UTC)

	h := func(w http.ResponseWriter, r *http.Request) {}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	_, err := athenaClient.ListLabResults(ctx, "", "", &ListLabResultsOptions{
		StartDate: &startDate,
	})
	assert.Error(err)
}

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

func TestHTTPClient_AddLabResultDocument_required_params_empty(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	h := func(w http.ResponseWriter, r *http.Request) {}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	b := bytes.NewReader([]byte(`test bytes`))
	_, err := athenaClient.AddLabResultDocumentReader(ctx, "", "", &AddLabResultDocumentReaderOptions{
		AttachmentContents: b,
		AttachmentType:     LabResultAttachmentTypeJPG,
	})

	assert.Error(err)
}

func TestHTTPClient_AddLabResultDocument(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	patientID := "123"
	departmentID := "456"

	observedAt := time.Date(2023, 9, 6, 16, 41, 3, 0, time.UTC)

	h := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(err)
		assert.Equal(departmentID, r.Form.Get("departmentid"))
		assert.Equal(string(LabResultAttachmentTypeJPG), r.Form.Get("attachmenttype"))

		obsDate := r.Form.Get("observationdate")
		assert.Equal("09/06/2023", obsDate)
		obsTime := r.Form.Get("observationtime")
		assert.Equal("16:41", obsTime)

		b, _ := os.ReadFile("./resources/AddLabResultDocument.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	b := bytes.NewReader([]byte(`test bytes`))
	res, err := athenaClient.AddLabResultDocumentReader(ctx, patientID, departmentID, &AddLabResultDocumentReaderOptions{
		AttachmentContents:  b,
		AttachmentType:      LabResultAttachmentTypeJPG,
		ObservationDateTime: NewObservationDateTime(observedAt),
	})

	assert.NoError(err)
	assert.Equal(res, 1083563)
}

func TestHTTPClient_AddLabResultDocument_observation_without_time(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	patientID := "123"
	departmentID := "456"

	observedAt := time.Date(2023, 9, 6, 16, 41, 3, 0, time.UTC)

	h := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(err)
		assert.Equal(departmentID, r.Form.Get("departmentid"))
		assert.Equal(string(LabResultAttachmentTypeJPG), r.Form.Get("attachmenttype"))

		obsDate := r.Form.Get("observationdate")
		assert.Equal("09/06/2023", obsDate)
		obsTime := r.Form.Get("observationtime")
		assert.Equal("", obsTime)

		b, _ := os.ReadFile("./resources/AddLabResultDocument.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	b := bytes.NewReader([]byte(`test bytes`))
	res, err := athenaClient.AddLabResultDocumentReader(ctx, patientID, departmentID, &AddLabResultDocumentReaderOptions{
		AttachmentContents:  b,
		AttachmentType:      LabResultAttachmentTypeJPG,
		ObservationDateTime: NewObservationDate(observedAt),
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
