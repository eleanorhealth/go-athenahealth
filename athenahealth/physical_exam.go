package athenahealth

import (
	"context"
	"net/url"
)

type GetPhysicalExamOpts struct {
	CCDOutput   bool     `,json:"ccdaoutputformat"`
	Structured  bool     `,json:"showstructured"`
	TemplateIDS []string `,json:"templateids"`
}

type PhysicalExam struct {
	LastModifiedBy                  string   `json:"lastmodifiedby"`
	LastModifiedDateTime            string   `json:"lastmodifieddatetime"`
	PhysicalExam                    []string `json:"physicalexam"`
	SectionNote                     string   `json:"sectionnote"`
	SectionNoteLastModifiedBy       string   `json:"sectionnotelastmodifiedby"`
	SectionNoteLastModifiedDateTime string   `json:"sectionnotelastmodifieddatetime"`
	SummaryText                     string   `json:"summarytext"`
	TemplateData                    []string `json:"templatedata"`
	Templates                       string   `json:"templates"`
}

// GetPhysicalExam - Get a physical exam
//
// GET /v1/{practiceid}/chart/encounter/{encounterid}/physicalexam
//
// https://docs.athenahealth.com/api/api-ref/physical-exam#Get-list-of-physical-exam-findings-and-notes-for-given-encounter
func (h *HTTPClient) GetPhysicalExam(ctx context.Context, opts *GetPhysicalExamOpts) (*PhysicalExam, error) {
	var out *PhysicalExam

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

	_, err := h.Get(ctx, "chart/encounter/{encounterid}/physicalexam", q, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
