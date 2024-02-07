package athenahealth

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type GetPhysicalExamOpts struct {
	CCDOutput   bool     `,json:"ccdaoutputformat"`
	Structured  bool     `,json:"showstructured"`
	TemplateIDS []string `,json:"templateids"`
}

type TemplateData struct {
	TemplateID   int    `json:"templateid"`
	TemplateName string `json:"templatename"`
}

type PhysicalExam struct {
	LastModifiedBy                  string         `json:"lastmodifiedby"`
	LastModifiedDateTime            time.Time      `json:"lastmodifieddatetime"`
	PhysicalExam                    []string       `json:"physicalexam"`
	SectionNote                     string         `json:"sectionnote"`
	SectionNoteLastModifiedBy       string         `json:"sectionnotelastmodifiedby"`
	SectionNoteLastModifiedDateTime string         `json:"sectionnotelastmodifieddatetime"`
	SummaryText                     string         `json:"summarytext"`
	TemplateData                    []TemplateData `json:"templatedata"`
	Templates                       []string       `json:"templates"`
}

// GetPhysicalExam - Get a physical exam
//
// GET /v1/{practiceid}/chart/encounter/{encounterid}/physicalexam
//
// https://docs.athenahealth.com/api/api-ref/physical-exam#Get-list-of-physical-exam-findings-and-notes-for-given-encounter
func (h *HTTPClient) GetPhysicalExam(ctx context.Context, encounterID string, opts *GetPhysicalExamOpts) (*PhysicalExam, error) {
	var out PhysicalExam

	if encounterID == "" {
		return nil, fmt.Errorf("encounterID empty")
	}

	q := url.Values{}

	if opts != nil {
		if opts.CCDOutput {
			q.Add("ccdaoutputformat", "true")
		}

		if opts.Structured {
			q.Add("showstructured", "true")
		}

		for _, id := range opts.TemplateIDS {
			q.Add("templateids", id)
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("chart/encounter/%s/physicalexam", encounterID), q, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
