package athenahealth

import (
	"context"
	"fmt"
	"net/url"
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

type TelehealthDeeplinkResult struct {
	ExpirationISOTimestamp string `json:"expiration"`
	Deeplink               string `json:"telehealthdeeplink"`
}

// TelehealthDeeplink - Create athenaone telehealth deep link join url
//
// GET /v1/{practiceid}/appointments/telehealth/deeplink
//
// https://docs.athenahealth.com/api/api-ref/appointment#Create-athenaone-telehealth-deep-link-join-url
func (h *HTTPClient) TelehealthDeeplink(ctx context.Context, apptID, patientID string) (*TelehealthDeeplinkResult, error) {
	if apptID == "" || patientID == "" {
		return nil, fmt.Errorf("apptID [%s] and patientID [%s] are required", apptID, patientID)
	}

	q := url.Values{}
	q.Set("patientid", patientID)
	q.Set("appointmentid", apptID)

	out := &TelehealthDeeplinkResult{}

	_, err := h.PostForm(ctx, "appointments/telehealth/deeplink", q, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
