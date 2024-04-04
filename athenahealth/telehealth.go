package athenahealth

import (
	"context"
	"fmt"
)

type GetTelehealthInviteURLResult struct {
	AppointmentID string `json:"appointmentid"`
	JoinToken     string `json:"jointoken"`
	PatientURL    string `json:"patienturl"`
}

// GetTelehealthInviteURL - Retrieve athenaone telehealth invite url
//
// GET /v1/{practiceid}/appointments/{appointmentid}/nativeathenatelehealthroom
//
// https://docs.athenahealth.com/api/api-ref/appointment#Retrieve-athenaone-telehealth-invite-url
func (h *HTTPClient) GetTelehealthInviteURL(ctx context.Context, apptID string) (*GetTelehealthInviteURLResult, error) {
	if apptID == "" {
		return nil, fmt.Errorf("cannot GetTelehealthInviteURL with empty apptID [%s]", apptID)
	}

	out := &GetTelehealthInviteURLResult{}

	_, err := h.Get(ctx, fmt.Sprintf("appointments/%s/nativeathenatelehealthroom", apptID), nil, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
