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

type UpdatePatientInsurancePackageOptions struct {
	PatientID   string
	InsuranceID string

	ExpirationDate                 *time.Time
	InsuranceIDNumber              *string
	InsurancePolicyHolderDOB       *time.Time
	InsurancePolicyHolderFirstName *string
	InsurancePolicyHolderLastName  *string
	InsurancePolicyHolderSex       *string
	NewSequenceNumber              *int
}

type updatePatientInsurancePackageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdatePatientInsurancePackage - Update a patient's specific insurance package.
// PUT /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}
// https://docs.athenahealth.com/api/api-ref/patient-insurance#Update-patient's-specific-insurance-package
func (h *HTTPClient) UpdatePatientInsurancePackage(ctx context.Context, opts *UpdatePatientInsurancePackageOptions) error {
	if opts == nil {
		panic("opts is nil")
	}

	out := &updatePatientInsurancePackageResponse{}

	form := url.Values{}

	if opts.ExpirationDate != nil {
		form.Add("expirationdate", opts.ExpirationDate.Format("01/02/2006"))
	}

	if opts.InsuranceIDNumber != nil {
		form.Add("insuranceidnumber", *opts.InsuranceIDNumber)
	}

	if opts.InsurancePolicyHolderFirstName != nil {
		form.Add("insurancepolicyholderfirstname", *opts.InsurancePolicyHolderFirstName)
	}

	if opts.InsurancePolicyHolderLastName != nil {
		form.Add("insurancepolicyholderlastname", *opts.InsurancePolicyHolderLastName)
	}

	if opts.InsurancePolicyHolderDOB != nil {
		form.Add("insurancepolicyholderdob", opts.InsurancePolicyHolderDOB.Format("01/02/2006"))
	}

	if opts.InsurancePolicyHolderSex != nil {
		form.Add("insurancepolicyholdersex", *opts.InsurancePolicyHolderSex)
	}

	if opts.NewSequenceNumber != nil {
		form.Add("newsequencenumber", strconv.Itoa(*opts.NewSequenceNumber))
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/patients/%s/insurances/%s", opts.PatientID, opts.InsuranceID), form, &out)
	if err != nil {
		return err
	}

	// TODO: are we ok with including this message in the error?
	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
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
