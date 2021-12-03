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
	Cancelled                           string `json:"cancelled"`
	EligibilityStatus                   string `json:"eligibilitystatus"`
	ExpirationDate                      string `json:"expirationdate"`
	InsuranceID                         string `json:"insuranceid"`
	InsuranceIDNumber                   string `json:"insuranceidnumber"`
	InsurancePackageAddress1            string `json:"insurancepackageaddress1"`
	InsurancePackageCity                string `json:"insurancepackagecity"`
	InsurancePackageID                  int    `json:"insurancepackageid"`
	InsurancePackageState               string `json:"insurancepackagestate"`
	InsurancePackageZip                 string `json:"insurancepackagezip"`
	InsurancePhone                      string `json:"insurancephone"`
	InsurancePlanDisplayName            string `json:"insuranceplandisplayname"`
	InsurancePlanName                   string `json:"insuranceplanname"`
	InsurancePolicyHolder               string `json:"insurancepolicyholder"`
	InsurancePolicyHolderCountryCode    string `json:"insurancepolicyholdercountrycode"`
	InsurancePolicyHoldercountryiso3166 string `json:"insurancepolicyholdercountryiso3166"`
	InsurancePolicyHolderdDOB           string `json:"insurancepolicyholderdob"`
	InsurancePolicyHolderFirstName      string `json:"insurancepolicyholderfirstname"`
	InsurancePolicyHolderLastName       string `json:"insurancepolicyholderlastname"`
	InsurancePolicyHolderSex            string `json:"insurancepolicyholdersex"`
	InsuranceType                       string `json:"insurancetype"`
	InsuredEntityTypeID                 int    `json:"insuredentitytypeid"`
	RelationshipToInsured               string `json:"relationshiptoinsured"`
	RelationshipToInsuredID             int    `json:"relationshiptoinsuredid"`
	SequenceNumber                      int    `json:"sequencenumber"`
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
	PatientID     string
	ShowCancelled bool

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
		if opts.ShowCancelled {
			q.Add("showcancelled", "true")
		}

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
