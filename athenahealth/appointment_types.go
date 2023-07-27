package athenahealth

import (
	"context"
	"net/url"
	"strconv"
)

type CreateAppointmentTypeOptions struct {
	// The expected duration, in minutes, of the appointment type. Note, this value cannot be changed after creation, so please choose carefully.
	Duration string `json:"duration"`
	// If set to true, this type serves as a "generic" type, that will match any type when searching. Defaults to false.
	Generic *bool `json:"generic"`
	// The name of the appointment type. Maximum length of 30 characters.
	Name string `json:"name"`
	// If set to true, this type serves as a "patient" type, meaning that is is a type that can be used for booking patients. If set to false, then it this type will not be used for patient (e.g. "Lunch" or "Vacation"). Non-patient types are mostly used to reserving time for providers to not see patients.
	Patient bool `json:"patient"`
	// The short name code of the appointment type. Maximum length of 4 characters. Used for making schedule templates. Note, this value cannot be changed after creation, so please choose carefully.
	ShortName string `json:"shortname"`
	// If set to true, this type serves as a "template-only" type, meaning that it can be used for building schedule templates, but cannot be used for booking appointments (i.e. another type must be chosen). Defaults to false.
	TemplateTypeOnly *bool `json:"templatetypeonly"`
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
