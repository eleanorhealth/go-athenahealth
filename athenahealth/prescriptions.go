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
	DepartmentID int     `json:"departmentid"`
	DocumentID   int     `json:"documentid"`
	PatientID    int     `json:"patientid"`
	ActionNote   *string `json:"actionnote"`
	AssignedTo   *string `json:"assignedto"`
	InternalNote *string `json:"internalnote"`
	Note         *string `json:"note"`
	PinToTop     *bool   `json:"pintotop"`
}

type UpdatePrescriptionResult struct {
	ErrorMessage *string `json:"errormessage,omitempty"`
	Success      bool    `json:"success"`
}

// UpdatePrescription - Update a prescription
//
//	/v1/{practiceid}/patients/{patientid}/documents/prescriptions/{prescriptionid}
//
// https://docs.athenahealth.com/api/api-ref/document-type-prescription#Update-specific-prescription-document-for-given-patient
func (h *HTTPClient) UpdatePrescription(ctx context.Context, departmentID int, patientID int, documentID int, opts *UpdatePrescriptionOptions) (*UpdatePrescriptionResult, error) {
	out := &UpdatePrescriptionResult{}

	form := url.Values{}

	form.Add("departmentid", strconv.Itoa(departmentID))

	if opts != nil {
		if opts.ActionNote != nil {
			form.Add("actionnote", *opts.ActionNote)
		}
		if opts.AssignedTo != nil {
			form.Add("assignedto", *opts.AssignedTo)
		}
		if opts.InternalNote != nil {
			form.Add("internalnote", *opts.InternalNote)
		}
		if opts.Note != nil {
			form.Add("note", *opts.Note)
		}
		if opts.PinToTop != nil {
			form.Add("pintotop", strconv.FormatBool(*opts.PinToTop))
		}
	}

	if _, err := h.PutForm(ctx, fmt.Sprintf("/patients/%d/documents/prescriptions/%d", patientID, documentID), form, out); err != nil {
		errMsg := err.Error()
		return &UpdatePrescriptionResult{
			Success:      false,
			ErrorMessage: &errMsg,
		}, err
	}

	if !out.Success {
		if out.ErrorMessage != nil && *out.ErrorMessage != "" {
			return out, fmt.Errorf(*out.ErrorMessage)
		}
		return out, fmt.Errorf("unexpected response from athena")
	}

	return out, nil
}
