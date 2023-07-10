package athenahealth

import (
	"context"
	"fmt"
)

// AppointmentCancelCheckIn cancels the check in process for the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/cancelcheckin
// https://docs.athenahealth.com/api/api-ref/appointment-check-in#Cancel-appointment-check-in-process
func (h *HTTPClient) AppointmentCancelCheckIn(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot AppointmentCancelCheckIn with empty apptID [%s]", apptID)
	}

	out := MessageResponse{}
	_, err := h.Post(ctx, fmt.Sprintf("/appointments/%s/cancelcheckin", apptID), nil, &out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
}

// AppointmentCheckIn checks in the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/checkin
// https://docs.athenahealth.com/api/api-ref/appointment-check-in#Check-in-this-appointment.
func (h *HTTPClient) AppointmentCheckIn(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot AppointmentCheckIn with empty apptID [%s]", apptID)
	}

	out := MessageResponse{}
	_, err := h.Post(ctx, fmt.Sprintf("/appointments/%s/checkin", apptID), nil, &out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
}

// AppointmentCheckOut checks out the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/checkout
// https://docs.athenahealth.com/api/api-ref/check-out#Complete-appointment-check-out-process
func (h *HTTPClient) AppointmentCheckOut(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot AppointmentCheckOut with empty apptID [%s]", apptID)
	}

	out := ErrorMessageResponse{}
	_, err := h.Post(ctx, fmt.Sprintf("/appointments/%s/checkout", apptID), nil, &out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
}

// StartCheckInAppointment starts the check in process for the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/startcheckin
// https://docs.athenahealth.com/api/api-ref/appointment-check-in#Initiate-appointment-check-in-process
func (h *HTTPClient) AppointmentStartCheckIn(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot AppointmentStartCheckIn with empty apptID [%s]", apptID)
	}

	out := MessageResponse{}
	_, err := h.Post(ctx, fmt.Sprintf("/appointments/%s/startcheckin", apptID), nil, &out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
}
