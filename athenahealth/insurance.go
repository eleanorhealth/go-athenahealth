package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type CreatePatientInsurancePackageOptions struct {
	PatientID                      string
	InsurancePackageID             int
	InsuranceIDNumber              string
	InsurancePolicyHolderFirstName string
	InsurancePolicyHolderLastName  string
	InsurancePolicyHolderDOB       time.Time
	InsurancePolicyHolderSex       string
	SequenceNumber                 int
}

type InsurancePackage struct {
	InsurancePolicyHolderCountryCode    string `json:"insurancepolicyholdercountrycode"`
	SequenceNumber                      int    `json:"sequencenumber"`
	InsurancePolicyHolderLastName       string `json:"insurancepolicyholderlastname"`
	InsuredEntityTypeID                 int    `json:"insuredentitytypeid"`
	InsuranceIDNumber                   string `json:"insuranceidnumber"`
	InsurancePolicyHolderdDOB           string `json:"insurancepolicyholderdob"`
	RelationshipToInsured               string `json:"relationshiptoinsured"`
	EligibilityStatus                   string `json:"eligibilitystatus"`
	InsurancePackageAddress1            string `json:"insurancepackageaddress1"`
	InsurancePolicyHolderSex            string `json:"insurancepolicyholdersex"`
	InsurancePlanName                   string `json:"insuranceplanname"`
	InsuranceType                       string `json:"insurancetype"`
	InsurancePhone                      string `json:"insurancephone"`
	InsurancePackageState               string `json:"insurancepackagestate"`
	InsurancePackageCity                string `json:"insurancepackagecity"`
	RelationshipToInsuredID             int    `json:"relationshiptoinsuredid"`
	InsuranceID                         string `json:"insuranceid"`
	InsurancePolicyHolder               string `json:"insurancepolicyholder"`
	InsurancePolicyHolderFirstName      string `json:"insurancepolicyholderfirstname"`
	InsurancePackageID                  int    `json:"insurancepackageid"`
	InsurancePolicyHoldercountryiso3166 string `json:"insurancepolicyholdercountryiso3166"`
	InsurancePlanDisplayName            string `json:"insuranceplandisplayname"`
	InsurancePackageZip                 string `json:"insurancepackagezip"`
}

// CreatePatientInsurancePackage - Create patient's insurance package.
// POST /v1/{practiceid}/patients/{patientid}/insurances
// https://docs.athenahealth.com/api/api-ref/patient-insurance#Create-patient's-insurance-package
func (h *HTTPClient) CreatePatientInsurancePackage(ctx context.Context, opts *CreatePatientInsurancePackageOptions) (*InsurancePackage, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := []*InsurancePackage{}

	form := url.Values{}

	form.Add("insurancepackageid", strconv.Itoa(opts.InsurancePackageID))
	form.Add("insuranceidnumber", opts.InsuranceIDNumber)
	form.Add("insurancepolicyholderfirstname", opts.InsurancePolicyHolderFirstName)
	form.Add("insurancepolicyholderlastname", opts.InsurancePolicyHolderLastName)
	form.Add("insurancepolicyholderdob", opts.InsurancePolicyHolderDOB.Format("01/02/2006"))
	form.Add("insurancepolicyholdersex", opts.InsurancePolicyHolderSex)
	form.Add("sequencenumber", strconv.Itoa(opts.SequenceNumber))

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/insurances", opts.PatientID), form, &out)
	if err != nil {
		return nil, err
	}

	if len(out) != 1 {
		return nil, errors.New("unexpected response")
	}

	return out[0], nil
}

type ListPatientInsurancePackagesOptions struct {
	PatientID string

	Pagination *PaginationOptions
}

type listPatientInsurancePackagesResponse struct {
	Insurances []*InsurancePackage `json:"insurances"`
	PaginationResponse
}

type ListPatientInsurancePackagesResult struct {
	InsurancePackages []*InsurancePackage
	Pagination        *PaginationResult
}

// ListPatientInsurancePackages - Get patient's insurance packages
// GET /v1/{practiceid}/patients/{patientid}/insurances
// https://docs.athenahealth.com/api/api-ref/patient-insurance#Get-patient's-insurance-packages
func (h *HTTPClient) ListPatientInsurancePackages(ctx context.Context, opts *ListPatientInsurancePackagesOptions) (*ListPatientInsurancePackagesResult, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := &listPatientInsurancePackagesResponse{}

	q := url.Values{}

	if opts.Pagination != nil {
		if opts.Pagination.Limit > 0 {
			q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
		}

		if opts.Pagination.Offset > 0 {
			q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("/patients/%s/insurances", opts.PatientID), q, &out)
	if err != nil {
		return nil, err
	}

	return &ListPatientInsurancePackagesResult{
		InsurancePackages: out.Insurances,
		Pagination:        makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}
