package athenahealth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ListChangedPrescriptionsOptions struct {
	LeaveUnprocessed           bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time

	Pagination *PaginationOptions
}

type ChangedPrescription struct {
	AdministerYN                string `json:"administeryn"`
	ApprovedBy                  string `json:"approvedby"`
	ApprovedTimestamp           string `json:"approvedtimestamp"`
	AssignedUser                string `json:"assigneduser"`
	Class                       string `json:"class"`
	ClassDescription            string `json:"classdescription"`
	ClinicalOrderTypeID         int    `json:"clinicalordertypeid"`
	ClinicalProviderOrderTypeID int    `json:"clinicalproviderordertypeid"`
	DateOrdered                 string `json:"dateordered"`
	DeniedBy                    string `json:"deniedby"`
	DeniedTimestamp             string `json:"deniedtimestamp"`
	DepartmentID                int    `json:"departmentid"`
	Description                 string `json:"description"`
	DocumentationOnly           bool   `json:"documentationonly"`
	DocumentID                  int    `json:"documentid"`
	EncounterID                 int    `json:"encounterid"`
	ExternalNote                string `json:"externalnote"`
	OrderGenusName              string `json:"ordergenusname"`
	OrderingProvider            string `json:"orderingprovider"`
	OutOfNetworkReason          string `json:"outofnetworkreason"`
	PatientID                   int    `json:"patientid"`
	Status                      string `json:"status"`
}

type listChangedPrescriptionsResponse struct {
	ChangedPrescriptions []*ChangedPrescription `json:"prescriptions"`

	*PaginationResponse
}

type ListChangedPrescriptionsResult struct {
	ChangedPrescriptions []*ChangedPrescription `json:"changedprescriptions"`

	Pagination *PaginationResult
}

// ListChangedPrescriptions - List of changes in prescriptions based on subscribed events
//
// GET /v1/{practiceid}/prescriptions/changed
//
// https://docs.athenahealth.com/api/api-ref/document-type-prescription#Get-list-of-changes-in-prescriptions
func (h *HTTPClient) ListChangedPrescriptions(ctx context.Context, opts *ListChangedPrescriptionsOptions) (*ListChangedPrescriptionsResult, error) {
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
		if opts.Pagination != nil {
			if opts.Pagination.Limit > 0 {
				q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
			}

			if opts.Pagination.Offset > 0 {
				q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
			}
		}
	}
	out := &listChangedPrescriptionsResponse{}

	_, err := h.Get(ctx, "/prescriptions/changed", q, out)
	if err != nil {
		return nil, err
	}

	return &ListChangedPrescriptionsResult{
		ChangedPrescriptions: out.ChangedPrescriptions,
		Pagination:           makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type UpdatePrescriptionOptions struct {
	DocumentID int `json:"documentid"`
	// DepartmentID int    `json:"departmentid"`
	PatientID  int    `json:"patientid"`
	ActionNote string `json:"actionnote"`
}

type UpdatePrescriptionResult struct {
	ErrorMessage string `json:"errormessage,omitempty"`
	Success      bool   `json:"success"`
}

// ListChangedPrescriptions - List of changes in prescriptions based on subscribed events
//
// GET /v1/{practiceid}/prescriptions/changed
//
// https://docs.athenahealth.com/api/api-ref/document-type-prescription#Get-list-of-changes-in-prescriptions
func (h *HTTPClient) UpdatePrescription(ctx context.Context, patientID int, documentID int, actionNote string) (*UpdatePrescriptionResult, error) {
	q := url.Values{}
	q.Add("actionnote", actionNote)

	out := &listChangedPrescriptionsResponse{}

	if _, err := h.PutForm(ctx, fmt.Sprintf("/patients/%v/documents/prescriptions/%v", patientID, documentID), q, &out); err != nil {
		return &UpdatePrescriptionResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	return &UpdatePrescriptionResult{
		Success:      true,
		ErrorMessage: "",
	}, nil
}
