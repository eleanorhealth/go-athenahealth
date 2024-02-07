package athenahealth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type GetPhysicalExamOpts struct {
	CCDOutput   bool     `,json:"ccdaoutputformat"`
	Structured  bool     `,json:"showstructured"`
	TemplateIDS []string `,json:"templateids"`
}

type PhysicalExam struct {
	LastModifiedBy                  string                     `json:"lastmodifiedby"`
	LastModifiedDateTime            time.Time                  `json:"lastmodifieddatetime"`
	PhysicalExam                    []PhysicalExamParagraph    `json:"physicalexam"`
	SectionNote                     string                     `json:"sectionnote"`
	SectionNoteLastModifiedBy       string                     `json:"sectionnotelastmodifiedby"`
	SectionNoteLastModifiedDateTime string                     `json:"sectionnotelastmodifieddatetime"`
	SummaryText                     string                     `json:"summarytext"`
	TemplateData                    []PhysicalExamTemplateData `json:"templatedata"`
	Templates                       []string                   `json:"templates"`
}

type PhysicalExamTemplateData struct {
	TemplateID   int    `json:"templateid"`
	TemplateName string `json:"templatename"`
}

type PhysicalExamParagraph struct {
	ParagraphID   int                    `json:"paragraphid"`
	ParagraphName string                 `json:"paragraphname"`
	Sentences     []PhysicalExamSentence `json:"sentences"`
}

type PhysicalExamSentence struct {
	SentenceID   int                   `json:"sentenceid"`
	SentenceName string                `json:"sentencename"`
	Findings     []PhysicalExamFinding `json:"findings"`
}

type PhysicalExamFinding struct {
	MedcinID           int      `json:"medcinid"`
	FindingName        string   `json:"findingname"`
	ContradictionIDs   []string `json:"contradictionids"`
	OptionLists        []string `json:"optionlists"`
	FindingType        string   `json:"findingtype"`
	GenericFindingName string   `json:"genericfindingname"`
	SelectedOptions    []string `json:"selectedoptions"`
	FindingID          int      `json:"findingid"`
	Selected           bool     `json:"selected"`
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

		if len(opts.TemplateIDS) > 0 {
			s, err := json.Marshal(opts.TemplateIDS)
			if err != nil {
				return nil, err
			}

			q.Add("templateids", string(s))
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("chart/encounter/%s/physicalexam", encounterID), q, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
