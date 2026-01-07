package athenahealth

import (
	"context"
	"net/url"
)

type ReferralSource struct {
	Name             string `json:"name"`
	ReferralSourceID string `json:"referralsourceid"`
}

type ListReferralSourcesResult struct {
	TotalCount      int               `json:"totalcount"`
	ReferralSources []*ReferralSource `json:"referralsources"`
}

// ListReferralSources - Get list of referral sources
//
// GET /v1/{practiceid}/referralsources
//
// https://docs.athenahealth.com/api/api-ref/referral-sources
func (h *HTTPClient) ListReferralSources(ctx context.Context) (*ListReferralSourcesResult, error) {
	out, q := &ListReferralSourcesResult{}, url.Values{}

	_, err := h.Get(ctx, "/referralsources", q, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
