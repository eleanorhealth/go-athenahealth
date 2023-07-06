package athenahealth

import (
	"context"
	"fmt"
)

// CancelCheckInAppointment cancels the check in process for the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/cancelcheckin
// https://docs.athenahealth.com/api/api-ref/appointment-check-in#Cancel-appointment-check-in-process
func (h *HTTPClient) CancelCheckInAppointment(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot CancelCheckInAppointment with empty apptID [%s]", apptID)
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

// CheckInAppointment checks in the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/checkin
// https://docs.athenahealth.com/api/api-ref/appointment-check-in#Check-in-this-appointment.
func (h *HTTPClient) CheckInAppointment(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot CheckInAppointment with empty apptID [%s]", apptID)
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

// CheckOutAppointment checks out the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/checkout
// https://docs.athenahealth.com/api/api-ref/check-out#Complete-appointment-check-out-process
func (h *HTTPClient) CheckOutAppointment(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot CheckOutAppointment with empty apptID [%s]", apptID)
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

type GetRequiredCheckInFieldsResult struct {
	FieldList []string `json:"fieldlist"`
}

// GetRequiredCheckInFields gets the fields required in order to perform check in per practice
// GET /v1/{practiceid}/departments/{departmentid}/checkinrequired
// https://docs.athenahealth.com/api/api-ref/required-fields-check#Get-list-of-required-fields-for-patient-check-in
func (h *HTTPClient) GetRequiredCheckInFields(ctx context.Context, departmentID string) (*GetRequiredCheckInFieldsResult, error) {
	if departmentID == "" {
		return nil, fmt.Errorf("cannot GetRequiredCheckInFields with empty departmentID [%s]", departmentID)
	}

	out := GetRequiredCheckInFieldsResult{}
	_, err := h.Get(ctx, fmt.Sprintf("/departments/%s/checkinrequired", departmentID), nil, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// StartCheckInAppointment starts the check in process for the provider and patient who are previously booked for the appointment ID
// POST /v1/{practiceid}/appointments/{appointmentid}/startcheckin
// https://docs.athenahealth.com/api/api-ref/appointment-check-in#Initiate-appointment-check-in-process
func (h *HTTPClient) StartCheckInAppointment(ctx context.Context, apptID string) error {
	if apptID == "" {
		return fmt.Errorf("cannot StartCheckInAppointment with empty apptID [%s]", apptID)
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
