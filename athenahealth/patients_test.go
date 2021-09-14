package athenahealth

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"
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
	patient, err := athenaClient.GetPatient(context.Background(), id, opts)

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

	res, err := athenaClient.ListPatients(context.Background(), opts)

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
	_, err := athenaClient.GetPatientPhoto(context.Background(), id, opts)

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
	err := athenaClient.UpdatePatientPhoto(context.Background(), id, data)
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

	patients, err := athenaClient.ListChangedPatients(context.Background(), opts)

	assert.Len(patients, 1)
	assert.NoError(err)
}

func TestHTTPClient_UpdatePatientInformationVerificationDetails(t *testing.T) {
	assert := assert.New(t)

	deptID := 1
	expirationDate := time.Now()
	insuredSignature := "true"
	patientSignature := "true"
	privacyNotice := "true"
	reasonPatientUnableToSign := "test reason"
	signatureDatetime := time.Now()
	signatureName := "John Smith"
	signerRelationshipToPatientID := "care provider"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Contains(r.URL.Path, "/patients/123/")

		assert.Equal(strconv.Itoa(deptID), r.FormValue("departmentid"))
		assert.Equal(expirationDate.Format("01/02/2006"), r.FormValue("expirationdate"))
		assert.Equal(insuredSignature, r.FormValue("insuredsignature"))
		assert.Equal(patientSignature, r.FormValue("patientsignature"))
		assert.Equal(privacyNotice, r.FormValue("privacynotice"))
		assert.Equal(reasonPatientUnableToSign, r.FormValue("reasonpatientunabletosign"))
		assert.Equal(signatureDatetime.Format("01/02/2006 15:04:05"), r.FormValue("signaturedatetime"))
		assert.Equal(signatureName, r.FormValue("signaturename"))
		assert.Equal(signerRelationshipToPatientID, r.FormValue("signerrelationshiptopatientid"))

		b, _ := ioutil.ReadFile("./resources/UpdatePatientInformationVerificationDetails.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &UpdatePatientInformationVerificationDetailsOptions{
		DepartmentID:                deptID,
		ExpirationDate:              &expirationDate,
		InsuredSignature:            &insuredSignature,
		PatientSignature:            &patientSignature,
		PrivacyNotice:               &privacyNotice,
		ReasonPatientUnableToSign:   &reasonPatientUnableToSign,
		SignatureDatetime:           signatureDatetime,
		SignatureName:               signatureName,
		SignerRelationshipToPatient: &signerRelationshipToPatientID,
	}

	err := athenaClient.UpdatePatientInformationVerificationDetails(context.Background(), "123", opts)
	assert.NoError(err)
}

func TestHTTPClient_GetPatientCustomFields(t *testing.T) {
	assert := assert.New(t)

	patientID := "1"
	departmentID := "2"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(r.URL.Path, "/patients/"+patientID+"/")

		assert.Equal(departmentID, r.URL.Query().Get("departmentid"))

		b, _ := ioutil.ReadFile("./resources/GetPatientCustomFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	customFields, err := athenaClient.GetPatientCustomFields(context.Background(), patientID, departmentID)
	assert.NoError(err)

	assert.Len(customFields, 2)

	assert.Contains(customFields, &CustomFieldValue{
		CustomFieldID:    "100",
		CustomFieldValue: "999999",
	})

	assert.Contains(customFields, &CustomFieldValue{
		CustomFieldID: "300",
		OptionID:      "3",
	})
}

func TestHTTPClient_UpdatePatientCustomFields(t *testing.T) {
	assert := assert.New(t)

	patientID := "1"
	departmentID := "2"
	customFields := []*CustomFieldValue{
		{
			CustomFieldID:    "3",
			CustomFieldValue: "foobar",
			OptionID:         "4",
		},
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Contains(r.URL.Path, "/patients/"+patientID+"/")

		assert.Equal(departmentID, r.FormValue("departmentid"))
		assert.Equal(`[{"customfieldid":"3","customfieldvalue":"foobar","optionid":"4"}]`, r.FormValue("customfields"))
		b, _ := ioutil.ReadFile("./resources/UpdatePatientCustomFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.UpdatePatientCustomFields(context.Background(), patientID, departmentID, customFields)
	assert.NoError(err)
}

func TestHTTPClient_ListPatientsMatchingCustomField(t *testing.T) {
	assert := assert.New(t)

	customFieldID := "1"
	customFieldValue := "foo"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(r.URL.Path, "/patients/customfields/"+customFieldID+"/"+customFieldValue)

		b, _ := ioutil.ReadFile("./resources/ListPatientsMatchingCustomField.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListPatientsMatchingCustomField(context.Background(), customFieldID, customFieldValue)

	assert.Len(res.Patients, 1)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 2)
	assert.NoError(err)
}
