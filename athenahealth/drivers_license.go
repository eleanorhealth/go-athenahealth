package athenahealth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"
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

type AddPatientDriversLicenseDocumentReaderOptions struct {
	DepartmentID string
	Image        io.Reader
}

type addPatientDriversLicenseDocumentReaderResponse struct {
	Success bool `json:"success"`
}

type AddPatientDriversLicenseDocumentReaderResult struct {
	Success bool
}

// AddPatientDriversLicenseDocumentReader - performs the same operation as AddPatientDriversLicenseDocument except is more memory efficient
// by streaming the image into the request, assuming you haven't already read the
// entire image into memory
// POST /v1/{practiceid}/patients/{patientid}/driverslicense
// https://docs.athenahealth.com/api/api-ref/drivers-license#Add-patient's-driver's-license-document
func (h *HTTPClient) AddPatientDriversLicenseDocumentReader(ctx context.Context, patientID string, opts *AddPatientDriversLicenseDocumentReaderOptions) (*AddPatientDriversLicenseDocumentReaderResult, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := &addPatientDriversLicenseDocumentReaderResponse{}

	form := newFormURLEncoder()

	if len(opts.DepartmentID) > 0 {
		form.Add("departmentid", strings.NewReader(opts.DepartmentID))
	}

	form.Add("image", newBase64Reader(opts.Image))

	_, err := h.PostFormReader(ctx, fmt.Sprintf("/patients/%s/driverslicense", patientID), form, &out)
	if err != nil {
		return nil, err
	}

	return &AddPatientDriversLicenseDocumentReaderResult{
		Success: out.Success,
	}, nil

}
