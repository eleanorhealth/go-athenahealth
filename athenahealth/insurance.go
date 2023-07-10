package athenahealth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

// CreatePatientInsurancePackage - Create patient's insurance package
//
// POST /v1/{practiceid}/patients/{patientid}/insurances
//
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

// ReactivatePatientInsurancePackage - Reactivate a patient's insurance package
// POST /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}/reactivate
// https://docs.athenahealth.com/api/api-ref/patient-insurance#Reactivate-patient's-specific-insurance-package
func (h *HTTPClient) ReactivatePatientInsurancePackage(ctx context.Context, patientID, insuranceID string, expirationDate *time.Time) error {
	out := &ErrorMessageResponse{}

	form := url.Values{}

	if expirationDate != nil {
		form.Add("expirationdate", expirationDate.Format("01/02/2006"))
	}

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/insurances/%s/reactivate", patientID, insuranceID), form, &out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
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

// UpdatePatientInsurancePackage - Update a patient's specific insurance package
//
// PUT /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}
//
// https://docs.athenahealth.com/api/api-ref/patient-insurance#Update-patient's-specific-insurance-package
func (h *HTTPClient) UpdatePatientInsurancePackage(ctx context.Context, opts *UpdatePatientInsurancePackageOptions) error {
	if opts == nil {
		panic("opts is nil")
	}

	out := &MessageResponse{}

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

	_, err := h.PutForm(ctx, fmt.Sprintf("/patients/%s/insurances/%s", opts.PatientID, opts.InsuranceID), form, out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
}

// DeletePatientInsurancePackage - Delete a patient's specific insurance package
// DELETE /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}
// https://docs.athenahealth.com/api/api-ref/patient-insurance#Delete-patient's-specific-insurance-package
func (h *HTTPClient) DeletePatientInsurancePackage(ctx context.Context, patientID, insuranceID, cancellationNote string) error {
	out := &MessageResponse{}

	form := url.Values{}

	if len(cancellationNote) > 0 {
		form.Add("cancellationnote", cancellationNote)
	}

	_, err := h.DeleteForm(ctx, fmt.Sprintf("/patients/%s/insurances/%s", patientID, insuranceID), form, out)
	if err != nil {
		return err
	}

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
//
// GET /v1/{practiceid}/patients/{patientid}/insurances
//
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

type UploadPatientInsuranceCardImageOptions struct {
	DepartmentID string
	Image        []byte
}

type uploadPatientInsuranceCardImageResponse struct {
	Success bool `json:"success"`
}

type UploadPatientInsuranceCardImageResult struct {
	Success bool
}

// UploadPatientInsuranceCardImage - Uploads the patient's insurance card image
//
// POST /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}/image
//
// https://docs.athenahealth.com/api/api-ref/insurance-card-image#Upload-patient's-insurance-card-image
func (h *HTTPClient) UploadPatientInsuranceCardImage(ctx context.Context, patientID, insuranceID string, opts *UploadPatientInsuranceCardImageOptions) (*UploadPatientInsuranceCardImageResult, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := &uploadPatientInsuranceCardImageResponse{}

	form := url.Values{}

	if len(opts.DepartmentID) > 0 {
		form.Add("departmentid", opts.DepartmentID)
	}

	form.Add("image", base64.StdEncoding.EncodeToString(opts.Image))

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/insurances/%s/image", patientID, insuranceID), form, &out)
	if err != nil {
		return nil, err
	}

	return &UploadPatientInsuranceCardImageResult{
		Success: out.Success,
	}, nil
}

type UploadPatientInsuranceCardImageReaderOptions struct {
	DepartmentID string
	Image        io.Reader
}

// UploadPatientInsuranceCardImageReader - performs the same operation as UploadPatientInsuranceCardImage except is more memory efficient
// by streaming the image into the request, assuming you haven't already read the
// entire image into memory
// POST /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}/image
// https://docs.athenahealth.com/api/api-ref/insurance-card-image#Upload-patient's-insurance-card-image
func (h *HTTPClient) UploadPatientInsuranceCardImageReader(ctx context.Context, patientID, insuranceID string, opts *UploadPatientInsuranceCardImageReaderOptions) (*UploadPatientInsuranceCardImageResult, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := &uploadPatientInsuranceCardImageResponse{}

	form := NewFormURLEncoder()

	if len(opts.DepartmentID) > 0 {
		form.AddString("departmentid", opts.DepartmentID)
	}

	form.AddReader("image", newBase64Reader(opts.Image))

	_, err := h.PostFormReader(ctx, fmt.Sprintf("/patients/%s/insurances/%s/image", patientID, insuranceID), form, &out)
	if err != nil {
		return nil, err
	}

	return &UploadPatientInsuranceCardImageResult{
		Success: out.Success,
	}, nil
}

type getPatientInsuranceCardImageResponse struct {
	Image string `json:"image"` // base64 encoded image
}

type GetPatientInsuranceCardImageResult struct {
	Image string
}

// GetPatientInsuranceCardImage - Gets the patient's insurance card image
//
// GET /v1/{practiceid}/patients/{patientid}/insurances/{insuranceid}/image
//
// ERR: athenahealth.ErrNotFound if there's no image
//
// https://docs.athenahealth.com/api/api-ref/insurance-card-image#Get-patient's-insurance-card-image
func (h *HTTPClient) GetPatientInsuranceCardImage(ctx context.Context, patientID, insuranceID string) (*GetPatientInsuranceCardImageResult, error) {
	out := &getPatientInsuranceCardImageResponse{}

	_, err := h.Get(ctx, fmt.Sprintf("/patients/%s/insurances/%s/image", patientID, insuranceID), nil, &out)
	if err != nil {
		return nil, err
	}

	return &GetPatientInsuranceCardImageResult{
		Image: out.Image,
	}, nil
}
