package athenahealth

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
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

		b, _ := ioutil.ReadFile("./resources/CreatePatientInsurancePackage.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	_, err := athenaClient.CreatePatientInsurancePackage(context.Background(), opts)
	assert.NoError(err)
}

func TestHTTPClient_ListPatientInsurancePackages(t *testing.T) {
	assert := assert.New(t)

	opts := &ListPatientInsurancePackagesOptions{
		PatientID:  "1",
		Pagination: &PaginationOptions{},
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/ListPatientInsurancePackages.json")
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
