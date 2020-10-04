package athenahealth

import (
	"errors"
	"fmt"
	"net/url"
)

type Provider struct {
	FirstName                   string `json:"firstname"`
	Specialty                   string `json:"specialty"`
	SpecialtyID                 int    `json:"specialtyid"`
	SchedulingName              string `json:"schedulingname"`
	ProviderTypeID              string `json:"providertypeid"`
	Billable                    bool   `json:"billable"`
	DisplayName                 string `json:"displayname"`
	ANSINameCode                string `json:"ansinamecode"`
	LastName                    string `json:"lastname"`
	ProviderID                  int    `json:"providerid"`
	SupervisingProviderUsername string `json:"supervisingproviderusername"`
	ProviderUsername            string `json:"providerusername"`
	SupervisingProviderID       int    `json:"supervisingproviderid"`
	ANSISpecialtyCode           string `json:"ansispecialtycode"`
	HideInPortal                bool   `json:"hideinportal"`
	Sex                         string `json:"sex"`
	EntityType                  string `json:"entitytype"`
	NPI                         int    `json:"npi"`
	ProviderType                string `json:"providertype"`
	CreateEncounterOnCheckIn    bool   `json:"createencounteroncheckin"`
}

// GetProvider - Get details about a single provider.
// GET /v1/{practiceid}/providers/{providerid}
// https://developer.athenahealth.com/docs/read/administrative/Providers#section-3
func (h *HTTPClient) GetProvider(id string) (*Provider, error) {
	out := []*Provider{}

	_, err := h.Get(fmt.Sprintf("/providers/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type ChangedProvider struct {
	ProviderID string `json:"providerid"`
}

type ListChangedProviderOptions struct {
}

type listChangedProvidersResponse struct {
	Providers []*ChangedProvider `json:"providers"`
}

// ListChangedProviders ...
func (h *HTTPClient) ListChangedProviders(opts *ListChangedProviderOptions) ([]*ChangedProvider, error) {
	out := &listChangedProvidersResponse{}

	q := url.Values{}

	if opts != nil {
	}

	_, err := h.Get("/providers/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.Providers, nil
}
