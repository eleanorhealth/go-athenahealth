package athenahealth

import (
	"context"
	"errors"
	"net/url"
)

type Medication struct {
	Medication   string `json:"medication"`
	MedicationID int    `json:"medicationid"`
}

// SearchMedications - Retrieves a list of medications for a given search parameters.
// GET /v1/{practiceid}/reference/medications
// https://docs.athenahealth.com/api/api-ref/medication#Search-for-available-medications
func (h *HTTPClient) SearchMedications(ctx context.Context, searchVal string) ([]*Medication, error) {
	out := []*Medication{}

	q := url.Values{}
	q.Add("searchvalue", searchVal)

	_, err := h.Get(ctx, "/reference/medications", q, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("unexpected length returned")
	}

	return out, nil
}
