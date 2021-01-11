package athenahealth

import (
	"net/url"
	"strconv"
	"time"
)

type ProblemEvent struct {
	EventType string `json:"eventtype"`
	StartDate string `json:"startdate"`
}

type Problem struct {
	Name      string         `json:"name"`
	PatientID int            `json:"patientid"`
	ProblemID int            `json:"problemid"`
	Events    []ProblemEvent `json:"events"`
	Codeset   string         `json:"codeset"`
	Code      string         `json:"code"`
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
