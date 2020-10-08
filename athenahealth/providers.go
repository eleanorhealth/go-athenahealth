package athenahealth

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
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
	ANSINameCode      string `json:"ansinamecode"`
	ANSISpecialtyCode string `json:"ansispecialtycode"`
	Billable          bool   `json:"billable"`
	DisplayName       string `json:"displayname"`
	EntityType        string `json:"entitytype"`
	FirstName         string `json:"firstname"`
	HideInPortal      bool   `json:"hideinportal"`
	LastName          string `json:"lastname"`
	ProviderID        int    `json:"providerid"`
	ProviderType      string `json:"providertype"`
	ProviderTypeID    string `json:"providertypeid"`
	SchedulingName    string `json:"schedulingname"`
	Specialty         string `json:"specialty"`
}

type ListChangedProviderOptions struct {
	LeaveUnprocessed           bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedProvidersResponse struct {
	ChangedProviders []*ChangedProvider `json:"providers"`
}

// ListChangedProviders - A list of all changes to providers.
// GET /v1/{practiceid}/providers/changed
// https://developer.athenahealth.com/docs/read/administrative/Providers#section-4
func (h *HTTPClient) ListChangedProviders(opts *ListChangedProviderOptions) ([]*ChangedProvider, error) {
	out := &listChangedProvidersResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.LeaveUnprocessed {
			q.Add("leaveunprocessed", strconv.FormatBool(opts.LeaveUnprocessed))
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
		}
	}

	_, err := h.Get("/providers/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedProviders, nil
}
