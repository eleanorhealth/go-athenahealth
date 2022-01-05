package athenahealth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ChangedProblem struct {
	Code      string                `json:"code"`
	CodeSet   string                `json:"codeset"`
	Events    []ChangedProblemEvent `json:"events"`
	Name      string                `json:"name"`
	PatientID int                   `json:"patientid"`
	ProblemID int                   `json:"problemid"`
}

type ChangedProblemEvent struct {
	Diagnoses  []ChangedProblemEventDiagnosis `json:"diagnoses"`
	EndDate    string                         `json:"enddate"`
	EventType  string                         `json:"eventtype"`
	Laterality string                         `json:"laterality"`
	Note       string                         `json:"note"`
	Source     string                         `json:"source"`
	StartDate  string                         `json:"startdate"`
	Status     string                         `json:"status"`
}

type ChangedProblemEventDiagnosis struct {
	Code    string `json:"code"`
	Codeset string `json:"codeset"`
	Name    string `json:"name"`
}

type ProblemEvent struct {
	EventType   string `json:"eventtype"`
	StartDate   string `json:"startdate"`
	CreatedDate string `json:"createddate"`
	OnsetDate   string `json:"onsetdate"`
	CreatedBy   string `json:"createdby"`
}

type Problem struct {
	LastModifiedDatetime string         `json:"lastmodifieddatetime"`
	LastModifiedBy       string         `json:"lastmodifiedby"`
	Name                 string         `json:"name"`
	PatientID            int            `json:"patientid"`
	ProblemID            int            `json:"problemid"`
	Events               []ProblemEvent `json:"events"`
	Codeset              string         `json:"codeset"`
	Code                 string         `json:"code"`
	BestMatchICD10Code   string         `json:"bestmatchicd10code"`
}

func (p *Problem) ICD10Code() string {
	if p.Codeset == "ICD10" {
		return p.Code
	} else {
		return p.BestMatchICD10Code
	}
}

type ListProblemsOptions struct {
	DepartmentID      string
	PatientID         string
	ShowDiagnosisInfo bool
}

type listProblemsResponse struct {
	Problems []*Problem `json:"problems"`
}

// ListProblems - Gets patient problems.
// GET /v1/{practiceid}/chart/{patientid}/problems
// https://docs.athenahealth.com/api/api-ref/problems#Get-patient's-problem-list
func (h *HTTPClient) ListProblems(ctx context.Context, patientID string, opts *ListProblemsOptions) ([]*Problem, error) {
	out := &listProblemsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if opts.ShowDiagnosisInfo {
			q.Add("showdiagnosisinfo", "true")
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("/chart/%s/problems", patientID), q, out)
	if err != nil {
		return nil, err
	}

	return out.Problems, nil
}

type ListChangedProblemsOptions struct {
	LeaveUnprocessed           bool
	PatientID                  string
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedProblemsResponse struct {
	ChangedProblems []*ChangedProblem `json:"problems"`
}

// ListChangedProblems - Gets changed problems records
// GET /v1/{practiceid}/chart/healthhistory/problems/changed
// https://docs.athenahealth.com/api/api-ref/problems#Get-list-of-changes-in-problems-based-on-subscribed-events
func (h *HTTPClient) ListChangedProblems(ctx context.Context, opts *ListChangedProblemsOptions) ([]*ChangedProblem, error) {
	out := &listChangedProblemsResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.LeaveUnprocessed {
			q.Add("leaveunprocessed", strconv.FormatBool(opts.LeaveUnprocessed))
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
		}
	}

	_, err := h.Get(ctx, "/chart/healthhistory/problems/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedProblems, nil
}
