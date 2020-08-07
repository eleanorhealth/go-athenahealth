package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetPatient(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("showcustomfields"))

		b, _ := ioutil.ReadFile("./resources/GetPatient.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &GetPatientOptions{
		ID:               "1",
		ShowCustomFields: true,
	}
	patient, err := athenaClient.GetPatient(opts)

	assert.NotNil(patient)
	assert.Nil(err)
}

func TestHTTPClient_ListPatients(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("John", r.URL.Query().Get("firstname"))
		assert.Equal("Smith", r.URL.Query().Get("lastname"))

		b, _ := ioutil.ReadFile("./resources/ListPatients.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListPatientsOptions{
		FirstName: "John",
		LastName:  "Smith",
	}

	patients, err := athenaClient.ListPatients(opts)

	assert.Len(patients, 2)
	assert.Nil(err)
}
