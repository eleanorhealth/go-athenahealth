package athenahealth

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

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
}

type ListProblemsOptions struct {
	DepartmentID string
	PatientID    string
}

type listProblemsResponse struct {
	Problems []*Problem `json:"problems"`
}

// ListProblems - Gets patient problems.
// GET /v1/{practiceid}/chart/{patientid}/problems
// https://developer.athenahealth.com/docs/read/chart/Problems#section-0
func (h *HTTPClient) ListProblems(patientID string, opts *ListProblemsOptions) ([]*Problem, error) {
	out := &listProblemsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}
	}

	_, err := h.Get(fmt.Sprintf("/chart/%s/problems", patientID), q, out)
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
	ChangedProblems []*Problem `json:"problems"`
}

// ListChangedProblems - Gets changed problems records
// GET /v1/{practiceid}/chart/healthhistory/problems/changed
// https://developer.athenahealth.com/docs/read/chart/Problems_Changed_Subscriptions#section-0
func (h *HTTPClient) ListChangedProblems(opts *ListChangedProblemsOptions) ([]*Problem, error) {
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

	_, err := h.Get("/chart/healthhistory/problems/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedProblems, nil
}
