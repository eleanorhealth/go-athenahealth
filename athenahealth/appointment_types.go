package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type CreateAppointmentTypeOptions struct {
	Duration         string `json:"duration"`
	Generic          *bool  `json:"generic"`
	Name             string `json:"name"`
	Patient          bool   `json:"patient"`
	ShortName        string `json:"shortname"`
	TemplateTypeOnly *bool  `json:"templatetypeonly"`
}

type CreateAppointmentTypeResult struct {
	AppointmentTypeID int `json:"appointmenttypeid"`
}

// CreateAppointmentType creates an Appointment Type
// POST /v1/{practiceid}/appointmenttypes
// https://docs.athenahealth.com/api/api-ref/appointment-types
func (h *HTTPClient) CreateAppointmentType(ctx context.Context, opts *CreateAppointmentTypeOptions) (*CreateAppointmentTypeResult, error) {
	out := CreateAppointmentTypeResult{}

	q := url.Values{}
	if opts != nil {
		if opts.Duration != "" {
			q.Set("duration", opts.Duration)
		}

		if opts.Generic != nil {
			q.Set("generic", strconv.FormatBool(*opts.Generic))
		}

		if opts.Name != "" {
			q.Set("name", opts.Name)
		}

		q.Set("patient", strconv.FormatBool(opts.Patient))

		if opts.ShortName != "" {
			q.Set("shortname", opts.ShortName)
		}

		if opts.TemplateTypeOnly != nil {
			q.Set("templatetypeonly", strconv.FormatBool(*opts.TemplateTypeOnly))
		}
	}

	_, err := h.PostForm(ctx, "/appointmenttypes", q, &out)

	return &out, err
}

// AppointmentType is an appointment type of the practice. Duration is in minutes: a booking of the type
// takes that long, whatever the length of the slot it is booked into.
type AppointmentType struct {
	AppointmentTypeID        int    `json:"appointmenttypeid"`
	CreateEncounterOnCheckIn bool   `json:"createencounteroncheckin"`
	DepartmentID             int    `json:"departmentid"`
	Duration                 int    `json:"duration"`
	Generic                  bool   `json:"generic"`
	Group                    bool   `json:"group"`
	Name                     string `json:"name"`
	Patient                  bool   `json:"patient"`
	PatientDisplayName       string `json:"patientdisplayname"`
	ProviderID               int    `json:"providerid"`
	ShortName                string `json:"shortname"`
	TemplateTypeOnly         bool   `json:"templatetypeonly"`

	// Only returned with ShowAppointmentTypeClasses.
	AppointmentTypeClassID   string `json:"appointmenttypeclassid"`
	AppointmentTypeClassName string `json:"appointmenttypeclassname"`
}

type GetAppointmentTypeOptions struct {
	ShowAppointmentTypeClasses bool
}

// GetAppointmentType returns one appointment type.
// GET /v1/{practiceid}/appointmenttypes/{appointmenttypeid}
// https://docs.athenahealth.com/api/api-ref/appointment-types#Get-one-appointment-type.
func (h *HTTPClient) GetAppointmentType(ctx context.Context, id string, opts *GetAppointmentTypeOptions) (*AppointmentType, error) {
	out := []*AppointmentType{}

	q := url.Values{}
	if opts != nil && opts.ShowAppointmentTypeClasses {
		q.Set("showappointmenttypeclasses", "true")
	}

	_, err := h.Get(ctx, fmt.Sprintf("/appointmenttypes/%s", id), q, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("unexpected length returned")
	}

	return out[0], nil
}

type ListAppointmentTypesOptions struct {
	HideGeneric                bool
	HideNonGeneric             bool
	HideTemplateTypeOnly       bool
	ShowAppointmentTypeClasses bool

	// HideNonPatient defaults to true in athena when not set, so non-patient types are hidden.
	// Set it to false to show them.
	HideNonPatient *bool

	// ProviderIDs and DepartmentIDs filter to the types of those providers and departments.
	ProviderIDs   []string
	DepartmentIDs []string

	Pagination *PaginationOptions
}

type ListAppointmentTypesResult struct {
	AppointmentTypes []*AppointmentType

	Pagination *PaginationResult
}

type listAppointmentTypesResponse struct {
	AppointmentTypes []*AppointmentType `json:"appointmenttypes"`

	PaginationResponse
}

// ListAppointmentTypes lists the appointment types of the practice.
// GET /v1/{practiceid}/appointmenttypes
// https://docs.athenahealth.com/api/api-ref/appointment-types#Get-list-of-appointment-types
func (h *HTTPClient) ListAppointmentTypes(ctx context.Context, opts *ListAppointmentTypesOptions) (*ListAppointmentTypesResult, error) {
	out := &listAppointmentTypesResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.HideGeneric {
			q.Set("hidegeneric", "true")
		}

		if opts.HideNonGeneric {
			q.Set("hidenongeneric", "true")
		}

		if opts.HideTemplateTypeOnly {
			q.Set("hidetemplatetypeonly", "true")
		}

		if opts.ShowAppointmentTypeClasses {
			q.Set("showappointmenttypeclasses", "true")
		}

		if opts.HideNonPatient != nil {
			q.Set("hidenonpatient", strconv.FormatBool(*opts.HideNonPatient))
		}

		if len(opts.ProviderIDs) > 0 {
			q.Set("providerids", strings.Join(opts.ProviderIDs, ","))
		}

		if len(opts.DepartmentIDs) > 0 {
			q.Set("departmentids", strings.Join(opts.DepartmentIDs, ","))
		}

		if opts.Pagination != nil {
			if opts.Pagination.Limit > 0 {
				q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
			}

			if opts.Pagination.Offset > 0 {
				q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
			}
		}
	}

	_, err := h.Get(ctx, "/appointmenttypes", q, out)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentTypesResult{
		AppointmentTypes: out.AppointmentTypes,
		Pagination:       makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}
