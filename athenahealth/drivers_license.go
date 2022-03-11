package athenahealth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
)

type AddPatientDriversLicenseDocumentOptions struct {
	DepartmentID string
	Image        []byte
}

type addPatientDriversLicenseDocumentResponse struct {
	Success bool `json:"success"`
}

type AddPatientDriversLicenseDocumentResult struct {
	Success bool
}

// AddPatientDriversLicenseDocument - Create a record of patient's driving license document
// POST /v1/{practiceid}/patients/{patientid}/driverslicense
// https://docs.athenahealth.com/api/api-ref/drivers-license#Add-patient's-driver's-license-document
func (h *HTTPClient) AddPatientDriversLicenseDocument(ctx context.Context, patientID string, opts *AddPatientDriversLicenseDocumentOptions) (*AddPatientDriversLicenseDocumentResult, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := &addPatientDriversLicenseDocumentResponse{}

	form := url.Values{}

	if len(opts.DepartmentID) > 0 {
		form.Add("departmentid", opts.DepartmentID)
	}

	form.Add("image", base64.StdEncoding.EncodeToString(opts.Image))

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/driverslicense", patientID), form, &out)
	if err != nil {
		return nil, err
	}

	return &AddPatientDriversLicenseDocumentResult{
		Success: out.Success,
	}, nil

}
