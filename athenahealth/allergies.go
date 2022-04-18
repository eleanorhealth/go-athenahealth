package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type AllergyReaction struct {
	ReactionName       string `json:"reactionname"`
	Severity           string `json:"severity"`
	SeveritySnomedCode int    `json:"severitysnomedcode"`
	SnomedCode         int    `json:"snomedcode"`
}

type Allergy struct {
	AllergenID           int               `json:"allergenid"`
	AllergenName         string            `json:"allergenname"`
	Categories           []string          `json:"categories"`
	Criticality          string            `json:"criticality"`
	DeactivateDate       string            `json:"deactivatedate"`
	DeactivatedUser      string            `json:"deactivateduser"`
	LastModifiedBy       string            `json:"lastmodifiedby"`
	LastModifiedDatetime string            `json:"lastmodifieddatetime"`
	Note                 string            `json:"note"`
	OnsetDate            string            `json:"onsetdate"`
	Reactions            []AllergyReaction `json:"reactions"`
}

type ListAllergiesResponse struct {
	Allergies                []*Allergy `json:"allergies"`
	LastModifiedBy           string     `json:"lastmodifiedby"`
	LastModifiedDatetime     string     `json:"lastmodifieddatetime"`
	LastUpdated              string     `json:"lastupdated"`
	NKDA                     bool       `json:"nkda"` // no known drug allergies
	NoteLastModifiedBy       string     `json:"notelastmodifiedby"`
	NoteLastModifiedDatetime string     `json:"notelastmodifieddatetime"`
	SectionNote              string     `json:"sectionnote"`
}
type ListAllergiesResult ListAllergiesResponse

// ListAllergies - Get patient's allergies
// GET /v1/{practiceid}/chart/{patientid}/allergies
// https://docs.athenahealth.com/api/api-ref/allergy#Get-patient's-allergies
func (h *HTTPClient) ListAllergies(ctx context.Context, patientID, departmentID string) (*ListAllergiesResult, error) {
	out := &ListAllergiesResponse{}

	q := url.Values{
		"departmentid": []string{departmentID},
	}

	_, err := h.Get(ctx, fmt.Sprintf("chart/%s/allergies", patientID), q, out)
	if err != nil {
		return nil, err
	}

	result := (*ListAllergiesResult)(out)

	return result, nil
}
