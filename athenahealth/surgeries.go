package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type Procedure struct {
	Description   string `json:"description"`
	Note          string `json:"note"`
	ProcedureCode string `json:"procedurecode"`
	ProcedureDate string `json:"proceduredate"`
	ProcedureID   int    `json:"procedureid"`
	ProviderID    int    `json:"providerid"`
	Source        string `json:"source"`
}

type ListSurgeriesResponse struct {
	Procedures  []*Procedure `json:"procedures"`
	SectionNote string       `json:"sectionnote"`
}
type ListSurgeriesResult ListSurgeriesResponse

// ListSurgeries - Get patient's surgical history data
// GET /v1/{practiceid}/chart/{patientid}/surgicalhistory
// https://docs.athenahealth.com/api/api-ref/surgical-history#Get-patient's-surgical-history-data
func (h *HTTPClient) ListSurgeries(ctx context.Context, patientID, departmentID string) (*ListSurgeriesResult, error) {
	out := &ListSurgeriesResponse{}

	q := url.Values{
		"departmentid": []string{departmentID},
	}

	_, err := h.Get(ctx, fmt.Sprintf("chart/%s/surgicalhistory", patientID), q, out)
	if err != nil {
		return nil, err
	}

	result := (*ListSurgeriesResult)(out)

	return result, nil
}
