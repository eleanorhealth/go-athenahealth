package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type CreateAppointmentTypeOptions struct {
	Duration         string `json:"duration"`
	Generic          bool   `json:"generic"`
	Name             string `json:"name"`
	Patient          bool   `json:"patient"`
	ShortName        string `json:"shortname"`
	TemplateTypeOnly bool   `json:"templatetypeonly"`
}

// CreateAppointmentType creates an Appointment Type
// POST /v1/{practiceid}/appointmenttypes
// https://docs.athenahealth.com/api/api-ref/appointment-types
func (h *HTTPClient) CreateAppointmentType(ctx context.Context, opts *CreateAppointmentTypeOptions) (int, error) {
	out := 0

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
			q.Set("generic", fmt.Sprintf("%v", opts.Generic))
		}
		if opts.Patient {
			q.Set("patient", fmt.Sprintf("%v", opts.Patient))
		}
		if opts.TemplateTypeOnly {
			q.Set("templatetypeonly", fmt.Sprintf("%v", opts.TemplateTypeOnly))
		}
	}

	_, err := h.PostForm(ctx, "/appointmenttypes", q, &out)

	return out, err
}
