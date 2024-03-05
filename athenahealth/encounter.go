package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type EncounterSummaryOptions struct {
	SkipAmendments bool
	Mobile         bool
}

type EncounterSummaryResponse struct {
	Summary string `json:"summaryhtml"`
}

// https://docs.athenahealth.com/api/api-ref/encounter-chart#Get-encounter-specific-encounter-summary-content
func (h *HTTPClient) EncounterSummary(ctx context.Context, encounterID string, opts *EncounterSummaryOptions) (*EncounterSummaryResponse, error) {
	out := &EncounterSummaryResponse{}

	if encounterID == "" {
		return nil, fmt.Errorf("encounterID required")
	}

	q := url.Values{}
	if opts != nil {
		if opts.SkipAmendments {
			q.Add("skipamendments", "1")
		}

		if opts.Mobile {
			q.Add("mobile", "1")
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("/chart/encounters/%s/summary", encounterID), q, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
