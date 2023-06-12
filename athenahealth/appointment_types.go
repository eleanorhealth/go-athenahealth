package athenahealth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type CreateAppointmentTypeOptions struct {
	Duration         string `json:"duration"`
	Generic          bool   `json:"generic"`
	Name             string `json:"name"`
	Patient          bool   `json:"patient"`
	ShortName        string `json:"shortname"`
	TemplateTypeOnly bool   `json:"templatetypeonly"`
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
		if opts.Name != "" {
			q.Set("name", opts.Name)
		}
		if opts.ShortName != "" {
			q.Set("shortname", opts.ShortName)
		}
		if opts.Generic {
			q.Set("generic", strconv.FormatBool(opts.Generic))
		}
		if opts.Patient {
			q.Set("patient", strconv.FormatBool(opts.Patient))
		}
		if opts.TemplateTypeOnly {
			q.Set("templatetypeonly", fmt.Sprintf("%v", opts.TemplateTypeOnly))
		}
	}

	_, err := h.PostForm(ctx, "/appointmenttypes", q, &out)

	return &out, err
}
