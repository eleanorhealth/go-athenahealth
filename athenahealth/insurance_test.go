package athenahealth

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_CreatePatientInsurancePackage(t *testing.T) {
	assert := assert.New(t)

	opts := &CreatePatientInsurancePackageOptions{
		PatientID:                      "1",
		InsurancePackageID:             2,
		InsuranceIDNumber:              "3",
		InsurancePolicyHolderFirstName: "John",
		InsurancePolicyHolderLastName:  "Smith",
		InsurancePolicyHolderDOB:       time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
		InsurancePolicyHolderSex:       "M",
		SequenceNumber:                 1,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("insurancepackageid"), strconv.Itoa(opts.InsurancePackageID))
		assert.Equal(r.Form.Get("insuranceidnumber"), opts.InsuranceIDNumber)
		assert.Equal(r.Form.Get("insurancepolicyholderfirstname"), opts.InsurancePolicyHolderFirstName)
		assert.Equal(r.Form.Get("insurancepolicyholderlastname"), opts.InsurancePolicyHolderLastName)
		assert.Equal(r.Form.Get("insurancepolicyholderdob"), opts.InsurancePolicyHolderDOB.Format("01/02/2006"))
		assert.Equal(r.Form.Get("insurancepolicyholdersex"), opts.InsurancePolicyHolderSex)
		assert.Equal(r.Form.Get("sequencenumber"), strconv.Itoa(opts.SequenceNumber))

		b, _ := os.ReadFile("./resources/CreatePatientInsurancePackage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	_, err := athenaClient.CreatePatientInsurancePackage(context.Background(), opts)
	assert.NoError(err)
}

func TestHTTPClient_UpdatePatientInsurancePackage(t *testing.T) {
	assert := assert.New(t)

	insuranceIDNumber := "3"
	insurancePolicyHolderFirstName := "John"
	insurancePolicyHolderLastName := "Smith"
	insurancePolicyHolderDOB := time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC)
	insurancePolicyHolderSex := "M"
	newSequenceNumber := 1

	opts := &UpdatePatientInsurancePackageOptions{
		PatientID:                      "1",
		InsuranceIDNumber:              &insuranceIDNumber,
		InsurancePolicyHolderFirstName: &insurancePolicyHolderFirstName,
		InsurancePolicyHolderLastName:  &insurancePolicyHolderLastName,
		InsurancePolicyHolderDOB:       &insurancePolicyHolderDOB,
		InsurancePolicyHolderSex:       &insurancePolicyHolderSex,
		NewSequenceNumber:              &newSequenceNumber,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("insuranceidnumber"), *opts.InsuranceIDNumber)
		assert.Equal(r.Form.Get("insurancepolicyholderfirstname"), *opts.InsurancePolicyHolderFirstName)
		assert.Equal(r.Form.Get("insurancepolicyholderlastname"), *opts.InsurancePolicyHolderLastName)
		assert.Equal(r.Form.Get("insurancepolicyholderdob"), opts.InsurancePolicyHolderDOB.Format("01/02/2006"))
		assert.Equal(r.Form.Get("insurancepolicyholdersex"), *opts.InsurancePolicyHolderSex)
		assert.Equal(r.Form.Get("newsequencenumber"), strconv.Itoa(*opts.NewSequenceNumber))

		b, _ := os.ReadFile("./resources/UpdatePatientInsurancePackage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.UpdatePatientInsurancePackage(context.Background(), opts)
	assert.NoError(err)
}

func TestHTTPClient_DeletePatientInsurancePackage(t *testing.T) {
	assert := assert.New(t)

	cancellationNote := "foo"

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "cancellationnote=foo")

		called = true

		b, _ := os.ReadFile("./resources/DeletePatientInsurancePackage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.DeletePatientInsurancePackage(context.Background(), "1", "2", cancellationNote)
	assert.NoError(err)

	assert.True(called)
}

func TestHTTPClient_ReactivatePatientInsurancePackage(t *testing.T) {
	assert := assert.New(t)

	patientID := "1"
	insuranceID := "2"
	expDate := time.Date(2022, time.January, 20, 0, 0, 0, 0, time.UTC)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/patients/1/insurances/2/reactivate", r.URL.Path)
		assert.Equal(http.MethodPost, r.Method)

		assert.NoError(r.ParseForm())
		assert.Equal(r.Form.Get("expirationdate"), expDate.Format("01/02/2006"))

		b, _ := os.ReadFile("./resources/ReactivatePatientInsurancePackage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.ReactivatePatientInsurancePackage(context.Background(), patientID, insuranceID, &expDate)
	assert.NoError(err)
}

func TestHTTPClient_ListPatientInsurancePackages(t *testing.T) {
	assert := assert.New(t)

	opts := &ListPatientInsurancePackagesOptions{
		PatientID:     "1",
		ShowCancelled: true,
		Pagination:    &PaginationOptions{},
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("showcancelled"), "true")

		b, _ := os.ReadFile("./resources/ListPatientInsurancePackages.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListPatientInsurancePackages(context.Background(), opts)
	assert.NoError(err)
	assert.Len(res.InsurancePackages, 1)
	assert.Equal(res.Pagination.TotalCount, 2)
	assert.NoError(err)
}

func TestHTTPClient_UploadPatientInsuranceCardImage(t *testing.T) {
	assert := assert.New(t)

	image := []byte("test attachment contents")
	deptID := "2"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Contains(r.URL.Path, "/patients/123/insurances/456/image")

		assert.Equal(base64.StdEncoding.EncodeToString([]byte(image)), r.FormValue("image"))
		assert.Equal(deptID, r.FormValue("departmentid"))

		b, _ := os.ReadFile("./resources/UploadPatientInsuranceCardImage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &UploadPatientInsuranceCardImageOptions{
		DepartmentID: deptID,
		Image:        image,
	}

	result, err := athenaClient.UploadPatientInsuranceCardImage(context.Background(), "123", "456", opts)
	assert.NoError(err)

	assert.True(result.Success)
}

func TestHTTPClient_GetPatientInsurancePackage(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patients/1/insurances/2/image", r.URL.String())

		b, _ := os.ReadFile("./resources/GetPatientInsuranceCardImage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.GetPatientInsuranceCardImage(context.Background(), "1", "2")
	assert.NoError(err)

	assert.True(called)

	d := base64.NewDecoder(base64.StdEncoding, strings.NewReader(res.Image))
	decoded := make([]byte, len(res.Image))
	n, err := d.Read(decoded)
	assert.NoError(err)

	assert.Equal("Hello World!", string(decoded[:n]))
}
