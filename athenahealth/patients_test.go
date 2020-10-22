package athenahealth

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetPatient(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("showcustomfields"))
		assert.Equal("true", r.URL.Query().Get("showinsurance"))
		assert.Equal("true", r.URL.Query().Get("showportalstatus"))
		assert.Equal("true", r.URL.Query().Get("showlocalpatientid"))

		b, _ := ioutil.ReadFile("./resources/GetPatient.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	id := "1"
	opts := &GetPatientOptions{
		ShowCustomFields:   true,
		ShowInsurance:      true,
		ShowPortalStatus:   true,
		ShowLocalPatientID: true,
	}
	patient, err := athenaClient.GetPatient(id, opts)

	assert.NotNil(patient)
	assert.NoError(err)
}

func TestHTTPClient_ListPatients(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("John", r.URL.Query().Get("firstname"))
		assert.Equal("Smith", r.URL.Query().Get("lastname"))
		assert.Equal("100", r.URL.Query().Get("departmentid"))
		assert.Equal("i", r.URL.Query().Get("status"))

		b, _ := ioutil.ReadFile("./resources/ListPatients.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListPatientsOptions{
		FirstName:    "John",
		LastName:     "Smith",
		DepartmentID: 100,
		Status:       "i",
	}

	res, err := athenaClient.ListPatients(opts)

	assert.Len(res.Patients, 2)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 2)
	assert.NoError(err)
}

func TestHTTPClient_GetPatientPhoto_JPEGOutputNotSupported(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	id := "1"
	opts := &GetPatientPhotoOptions{
		JPEGOutput: true,
	}
	_, err := athenaClient.GetPatientPhoto(id, opts)

	assert.Error(err)
}

func TestHTTPClient_UpdatePatientPhoto(t *testing.T) {
	assert := assert.New(t)

	data := []byte("Hello World!")

	h := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		assert.Equal(base64.StdEncoding.EncodeToString(data), r.Form.Get("image"))
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	id := "1"
	err := athenaClient.UpdatePatientPhoto(id, data)
	assert.NoError(err)
}

func TestHTTPClient_ListChangedPatients(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("d1", r.URL.Query().Get("departmentid"))
		assert.Equal("true", r.URL.Query().Get("ignorerestrictions"))
		assert.Equal("true", r.URL.Query().Get("leaveunprocessed"))
		assert.Equal("p1", r.URL.Query().Get("patientid"))
		assert.Equal("true", r.URL.Query().Get("returnglobalid"))
		assert.Equal("06/01/2020 15:30:45", r.URL.Query().Get("showprocessedstartdatetime"))
		assert.Equal("06/02/2020 12:30:45", r.URL.Query().Get("showprocessedenddatetime"))

		b, _ := ioutil.ReadFile("./resources/ListChangedPatients.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListChangedPatientOptions{
		DepartmentID:               "d1",
		IgnoreRestrictions:         true,
		LeaveUnprocessed:           true,
		PatientID:                  "p1",
		ReturnGlobalID:             true,
		ShowProcessedStartDatetime: time.Date(2020, 6, 1, 15, 30, 45, 0, time.UTC),
		ShowProcessedEndDatetime:   time.Date(2020, 6, 2, 12, 30, 45, 0, time.UTC),
	}

	patients, err := athenaClient.ListChangedPatients(opts)

	assert.Len(patients, 1)
	assert.NoError(err)
}
