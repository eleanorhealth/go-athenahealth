package athenahealth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type CreateAppointmentSlotOptions struct {
	// The appointment date for the new open appointment slot (mm/dd/yyyy).
	AppointmentDate string `json:"appointmentdate"`
	// The time (hh24:mi) for the new appointment slot. Multiple times (either as a comma delimited list or multiple POSTed values) are allowed. 24 hour time.
	AppointmentTime []string `json:"appointmenttime"`
	// The appointment type ID to be created. Either this or a reason must be provided.
	AppointmentTypeID *string `json:"appointmenttypeid"`
	// The athenaNet department ID.
	DepartmentID string `json:"departmentid"`
	// The athenaNet provider ID.
	ProviderID string `json:"providerid"`
	// The appointment reason (/patientappointmentreasons) to be created. Either this or a raw appointment type ID must be provided.
	ReasonID *string `json:"reasonid"`
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
		if opts.AppointmentDate == "" {
			return nil, fmt.Errorf("cannot CreateAppointmentSlot when AppointmentDate is empty [%s]", opts.AppointmentDate)
		} else {
			q.Set("appointmentdate", opts.AppointmentDate)
		}

		if len(opts.AppointmentTime) > 0 {
			q.Set("appointmenttime", strings.Join(opts.AppointmentTime, ","))
		} else {
			return nil, fmt.Errorf("cannot CreateAppointmentSlot without at least one AppointmentTime [%+v]", opts.AppointmentTime)
		}

		if opts.AppointmentTypeID != nil && *opts.AppointmentTypeID != "" {
			q.Set("appointmenttypeid", *opts.AppointmentTypeID)
		}

		q.Set("departmentid", opts.DepartmentID)

		q.Set("providerid", opts.ProviderID)

		if opts.ReasonID != nil {
			q.Set("reasonid", *opts.ReasonID)
		}
	}

	_, err := h.PostForm(ctx, "/appointments/open", q, &out)

	return &out, err
}
