package athenahealth

import (
	"context"
	"errors"
	"net/url"
)

type Allergy struct {
	AllergenID   int    `json:"allergenid"`
	AllergenName string `json:"allergenname"`
	Allergy      string `json:"allergy"`
	AllergyID    int    `json:"allergyid"`
}

// SearchAllergies - Retrieves a list of allergies matching the search criteria.
// GET /v1/{practiceid}/reference/allergies
// https://docs.athenahealth.com/api/api-ref/allergy#Search-for-available-allergies
func (h *HTTPClient) SearchAllergies(ctx context.Context, searchVal string) ([]*Allergy, error) {
	out := []*Allergy{}

	q := url.Values{}
	q.Add("searchvalue", searchVal)

	_, err := h.Get(ctx, "/reference/allergies", q, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("unexpected length returned")
	}

	return out, nil
}
