package athenahealth

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

type CreateAppointmentSlotOptions struct {
	AppointmentDate   string   `json:"appointmentdate"`
	AppointmentTime   []string `json:"appointmenttime"`
	AppointmentTypeID *int     `json:"appointmenttypeid"`
	DepartmentID      int      `json:"departmentid"`
	ProviderID        int      `json:"providerid"`
	ReasonID          *int     `json:"reasonid"`
}

type CreateAppointmentSlotResult struct {
	AppointmentIDs map[string]string `json:"appointmentids"`
}

// CreateAppointmentSlot creates an Appointment Slot
// POST /v1/{practiceid}/appointments/open
// https://docs.athenahealth.com/api/api-ref/appointment-slot#Create-a-new-appointment-slot
func (h *HTTPClient) CreateAppointmentSlot(ctx context.Context, opts *CreateAppointmentSlotOptions) (*CreateAppointmentSlotResult, error) {
	out := CreateAppointmentSlotResult{}

	q := url.Values{}
	if opts != nil {
		if opts.AppointmentDate != "" {
			q.Set("appointmentdate", opts.AppointmentDate)
		}

		if len(opts.AppointmentTime) > 0 {
			q.Set("appointmenttime", strings.Join(opts.AppointmentTime, ","))
		}

		if opts.AppointmentTypeID != nil {
			q.Set("appointmenttypeid", strconv.Itoa(*opts.AppointmentTypeID))
		}

		q.Set("departmentid", strconv.Itoa(opts.DepartmentID))

		q.Set("providerid", strconv.Itoa(opts.ProviderID))

		if opts.ReasonID != nil {
			q.Set("reasonid", strconv.Itoa(*opts.ReasonID))
		}
	}

	_, err := h.PostForm(ctx, "/appointments/open", q, &out)

	return &out, err
}
