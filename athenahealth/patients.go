package athenahealth

import (
	"errors"
	"fmt"
	"net/url"
)

// Patient represents a patient in athenahealth.
type Patient struct {
	PatientID string `json:"patientid"`

	FirstName  string `json:"firstname"`
	MiddleName string `json:"middlename"`
	LastName   string `json:"lastname"`

	MobilePhone string `json:"mobilephone"`
	HasMobile   bool   `json:"hasmobile"`

	Email string `json:"email"`

	ConsentToCall bool `json:"consenttocall"`
	ConsentToText bool `json:"consenttotext"`
}

// GetPatient - Full view/update of patient demographics.
// GET /v1/{practiceid}/patients/{patientid}
// https://developer.athenahealth.com/docs/read/patientinfo/Patient_Information#section-5
func (h *HTTPClient) GetPatient(id string) (*Patient, error) {
	out := []*Patient{}

	_, err := h.Get(fmt.Sprintf("/patients/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type ListPatientsOptions struct {
	FirstName string
	LastName  string
}

type listPatientsResponse struct {
	Patients []*Patient `json:"patients"`
}

// ListPatients - Gets a set of patients or creates a patient.
// GET /v1/{practiceid}/patients
// https://developer.athenahealth.com/docs/read/patientinfo/Patient_Information#section-1
func (h *HTTPClient) ListPatients(opts *ListPatientsOptions) ([]*Patient, error) {
	out := &listPatientsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.FirstName) > 0 {
			q.Add("firstname", opts.FirstName)
		}

		if len(opts.LastName) > 0 {
			q.Add("lastname", opts.LastName)
		}
	}

	_, err := h.Get("/patients", q, out)
	if err != nil {
		return nil, err
	}

	return out.Patients, nil
}
