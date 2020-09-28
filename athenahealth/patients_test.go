package athenahealth

import (
	"encoding/base64"
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

	id := "1"
	opts := &GetPatientOptions{
		ShowCustomFields: true,
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
