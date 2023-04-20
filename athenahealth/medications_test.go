package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListMedications(t *testing.T) {
	assert := assert.New(t)

	patientID := "123"
	departmentID := "789"
	opts := &ListMedicationsOptions{DepartmentID: departmentID, MedicationType: MedicationTypeActive}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(departmentID, r.URL.Query().Get("departmentid"))
		assert.Equal(MedicationTypeActive, r.URL.Query().Get("medicationtype"))

		b, _ := os.ReadFile("./resources/ListMedications.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	medicationsResult, err := athenaClient.ListMedications(context.Background(), patientID, opts)

	assert.NotNil(medicationsResult)
	assert.Len(medicationsResult.Medications, 5)
	assert.Equal("Staff Stafferson", medicationsResult.Medications[0][0].PrescribedBy)
	assert.NoError(err)
}

func TestHTTPClient_ListMedications_None(t *testing.T) {
	assert := assert.New(t)

	patientID := "123"
	departmentID := "789"
	opts := &ListMedicationsOptions{DepartmentID: departmentID}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(departmentID, r.URL.Query().Get("departmentid"))

		b := []byte(`{
			"patientneedsdownloadconsent": false,"medications": [],"patientdownloadconsent": true,"nomedicationsreported": false}`)
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	medicationsResult, err := athenaClient.ListMedications(context.Background(), patientID, opts)

	assert.NotNil(medicationsResult)
	assert.Len(medicationsResult.Medications, 0)
	assert.NoError(err)
}

func TestHTTPClient_SearchMedications(t *testing.T) {
	assert := assert.New(t)

	searchVal := "ibuprofen"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b, _ := os.ReadFile("./resources/SearchMedications.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	medications, err := athenaClient.SearchMedications(context.Background(), searchVal)

	assert.NotNil(medications)
	assert.Len(medications, 2)
	assert.NoError(err)
}

func TestHTTPClient_SearchMedications_None(t *testing.T) {
	assert := assert.New(t)

	searchVal := "sdfsadss"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(searchVal, r.URL.Query().Get("searchvalue"))

		b := []byte("[]")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	meds, err := athenaClient.SearchAllergies(context.Background(), searchVal)

	assert.NotNil(meds)
	assert.Len(meds, 0)
	assert.NoError(err)
}
