package athenahealth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type SocialHistoryTemplate struct {
	Questions    []*SocialHistoryQuestion `json:"questions"`
	TemplateID   json.Number              `json:"templateid"`
	Templatename string                   `json:"templatename"`
}

type SocialHistoryQuestion struct {
	InputType  string        `json:"inputtype"`
	Key        string        `json:"key"`
	Options    []interface{} `json:"options"`
	Ordering   int           `json:"ordering"`
	Question   string        `json:"question"`
	QuestionID NumberString  `json:"questionid"`
}

// ListPatientSocialHistoryTemplates - List of social history questions and templates configured by this practice.
// GET /v1/{practiceid}/chart/configuration/socialhistory
// https://developer.athenahealth.com/docs/read/chart/Social_History#section-0
func (h *HTTPClient) ListSocialHistoryTemplates(ctx context.Context) ([]*SocialHistoryTemplate, error) {
	out := []*SocialHistoryTemplate{}

	_, err := h.Get(ctx, "/chart/configuration/socialhistory", nil, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type PatientSocialHistoryQuestion struct {
	Answer              string      `json:"answer"`
	Key                 string      `json:"key"`
	Lastupdated         string      `json:"lastupdated"`
	Note                string      `json:"note"`
	NoteLastUpdatedDate string      `json:"notelastupdateddate"`
	Ordering            int         `json:"ordering"`
	Question            string      `json:"question"`
	QuestionID          json.Number `json:"questionid"`
	TemplateID          json.Number `json:"templateid"`
}

type GetPatientSocialHistoryOptions struct {
	DepartmentID              string
	RecipientCategory         string
	ShowNotPerformedQuestions bool
	ShowUnansweredQuestions   bool
}

type GetPatientSocialHistoryResponse struct {
	Questions   []*PatientSocialHistoryQuestion `json:"questions"`
	SectionNote string                          `json:"sectionnote"`
	Templates   []*struct {
		TemplateID   json.Number `json:"templateid"`
		TemplateName string      `json:"templatename"`
	} `json:"templates"`
}

// GetPatientSocialHistory - List of social history data for this patient.
// GET /v1/{practiceid}/chart/{patientid}/socialhistory
// https://developer.athenahealth.com/docs/read/chart/Social_History#section-2
func (h *HTTPClient) GetPatientSocialHistory(ctx context.Context, patientID string, opts *GetPatientSocialHistoryOptions) (*GetPatientSocialHistoryResponse, error) {
	out := &GetPatientSocialHistoryResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.DepartmentID) != 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.RecipientCategory) != 0 {
			q.Add("recipientcategory", opts.RecipientCategory)
		}

		if opts.ShowNotPerformedQuestions {
			q.Add("shownotperformedquestions", "true")
		}

		if opts.ShowUnansweredQuestions {
			q.Add("showunansweredquestions", "true")
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("chart/%s/socialhistory", patientID), q, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type UpdatePatientSocialHistoryQuestion struct {
	Answer             string          `json:"answer"`
	Delete             bool            `json:"delete"`
	Key                string          `json:"key"`
	Note               string          `json:"note"`
	NotPerformedReason json.RawMessage `json:"notperformedreason"`
}

type UpdatePatientSocialHistoryOptions struct {
	DepartmentID string
	Questions    []*UpdatePatientSocialHistoryQuestion
	SectionNote  string
}

// UpdatePatientSocialHistory - Update the set of social history questions for this patient.
// PUT /v1/{practiceid}/chart/{patientid}/socialhistory
// https://developer.athenahealth.com/docs/read/chart/Social_History#section-2
func (h *HTTPClient) UpdatePatientSocialHistory(ctx context.Context, patientID string, opts *UpdatePatientSocialHistoryOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.DepartmentID) > 0 {
			form.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.Questions) > 0 {
			b, err := json.Marshal(opts.Questions)
			if err != nil {
				return err
			}

			form.Add("questions", string(b))
		}

		if len(opts.SectionNote) > 0 {
			form.Add("sectionnote", opts.SectionNote)
		}
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/chart/%s/socialhistory", patientID), form, nil)
	if err != nil {
		return err
	}

	return nil
}
