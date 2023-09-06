package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

type LabResult struct {
	AttachmentExists         bool   `json:"attachmentexists"`
	Description              string `json:"description"`
	ExactDuplicateDocumentID int    `json:"exactduplicatedocumentid"`
	FacilityID               int    `json:"facilityid"`
	InternalNote             string `json:"internalnote"`
	IsReviewedByProvider     string `json:"isreviewedbyprovider"`
	LabResultDate            string `json:"labresultdate"`
	LabResultDateTime        string `json:"labresultdatetime"`
	LabResultID              int    `json:"labresultid"`
	LabResultLOINC           string `json:"labresultloinc"`
	LabResultNote            string `json:"labresultnote"`
	OrderID                  int    `json:"orderid"`
	PatientNote              string `json:"patientnote"`
	Priority                 string `json:"priority"`
	ProviderID               int    `json:"providerid"`
	ResultStatus             string `json:"resultstatus"`
}

type ListLabResultsOptions struct {
	StartDate           *time.Time
	LabResultStatus     *string
	ShowHidden          *bool
	ShowAbnormalDetails *bool
	EndDate             *time.Time
	HideDuplicate       *bool

	Pagination *PaginationOptions
}

type listLabResultsResponse struct {
	LabResults []*LabResult `json:"results"`

	PaginationResponse
}

type ListLabResultsResult struct {
	LabResults []*LabResult

	Pagination *PaginationResult
}

// ListLabResults retrieves the patients laboratory results
//
// GET /v1/{practiceid}/chart/{patientid}/labresults
//
// https://docs.athenahealth.com/api/api-ref/lab-result#Get-patient's-lab-results
func (h *HTTPClient) ListLabResults(ctx context.Context, patientID string, departmentID string, opts *ListLabResultsOptions) (*ListLabResultsResult, error) {
	var requiredParamErrors []error
	if len(patientID) == 0 {
		requiredParamErrors = append(requiredParamErrors, errors.New("patientID is required"))
	}
	if len(departmentID) == 0 {
		requiredParamErrors = append(requiredParamErrors, errors.New("departmentID is required"))
	}
	if len(requiredParamErrors) > 0 {
		return nil, errors.Join(requiredParamErrors...)
	}

	q := url.Values{}
	q.Add("departmentid", departmentID)

	if opts != nil {
		if opts.StartDate != nil && !opts.StartDate.IsZero() {
			q.Add("startdate", opts.StartDate.Format("01/02/2006"))
		}
		if opts.LabResultStatus != nil {
			q.Add("labresultstatus", *opts.LabResultStatus)
		}
		if opts.ShowHidden != nil {
			q.Add("showhidden", strconv.FormatBool(*opts.ShowHidden))
		}
		if opts.ShowAbnormalDetails != nil {
			q.Add("showabnormaldetails", strconv.FormatBool(*opts.ShowAbnormalDetails))
		}
		if opts.EndDate != nil && !opts.EndDate.IsZero() {
			q.Add("enddate", opts.EndDate.Format("01/02/2006"))
		}
		if opts.HideDuplicate != nil {
			q.Add("hideduplicate", strconv.FormatBool(*opts.HideDuplicate))
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

	out := &listLabResultsResponse{}

	_, err := h.Get(ctx, fmt.Sprintf("chart/%s/labresults", patientID), q, out)
	if err != nil {
		return nil, err
	}

	return &ListLabResultsResult{
		LabResults: out.LabResults,
		Pagination: makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type LabResultAttachmentType string

const (
	LabResultAttachmentTypeBMP  LabResultAttachmentType = "BMP"
	LabResultAttachmentTypeGIF  LabResultAttachmentType = "GIF"
	LabResultAttachmentTypeJPG  LabResultAttachmentType = "JPG"
	LabResultAttachmentTypeJPEG LabResultAttachmentType = "JPEG"
	LabResultAttachmentTypePDF  LabResultAttachmentType = "PDF"
	LabResultAttachmentTypeTIF  LabResultAttachmentType = "TIF"
	LabResultAttachmentTypeTIFF LabResultAttachmentType = "TIFF"
)

type AddLabResultDocumentOptions struct {
	// AttachmentContents must be Base64 encoded
	AttachmentContents io.Reader
	AttachmentType     LabResultAttachmentType
	InternalNote       *string
	NoteToPatient      *string
	// Sets both observationdate and observationtime if not nil
	ObservedAt       *time.Time
	OriginalFilename *string
	// 1 = high, 2 = normal
	Priority    *string
	ResultNotes *string
	// Final, Partial, Pending, Preliminary, Corrected, Cancelled
	ResultStatus *string
	TieToOrderID *int
}

type addLabResultDocumentResponse struct {
	ErrorMessage string `json:"errormessage"`
	LabResultID  int    `json:"labresultid"`
	Success      bool   `json:"success"`
}

// AddLabResultDocument creates a lab result document record of a specific patient
//
// POST /v1/{practiceid}/patients/{patientid}/documents/labresult
//
// https://docs.athenahealth.com/api/api-ref/document-type-lab-result#Add-lab-result-document-to-patient's-chart
func (h *HTTPClient) AddLabResultDocumentReader(ctx context.Context, patientID string, departmentID string, opts *AddLabResultDocumentOptions) (int, error) {
	var requiredParamErrors []error
	if len(patientID) == 0 {
		requiredParamErrors = append(requiredParamErrors, errors.New("patientID is required"))
	}
	if len(departmentID) == 0 {
		requiredParamErrors = append(requiredParamErrors, errors.New("departmentID is required"))
	}
	if len(requiredParamErrors) > 0 {
		return 0, errors.Join(requiredParamErrors...)
	}

	form := NewFormURLEncoder()

	form.AddString("departmentid", departmentID)

	if opts != nil {
		if opts.AttachmentContents != nil {
			form.AddReader("attachmentcontents", opts.AttachmentContents)
		}
		if len(opts.AttachmentType) > 0 {
			form.AddString("attachmenttype", string(opts.AttachmentType))
		}
		if opts.InternalNote != nil {
			form.AddString("internalnote", string(*opts.InternalNote))
		}
		if opts.NoteToPatient != nil {
			form.AddString("notetopatient", string(*opts.NoteToPatient))
		}
		if opts.ObservedAt != nil {
			form.AddString("observationdate", opts.ObservedAt.Format("01/02/2006"))
			form.AddString("observationtime", opts.ObservedAt.Format("15:04"))
		}
		if opts.OriginalFilename != nil {
			form.AddString("originalfilename", string(*opts.OriginalFilename))
		}
		if opts.Priority != nil {
			form.AddString("priority", string(*opts.Priority))
		}
		if opts.ResultNotes != nil {
			form.AddString("resultnotes", string(*opts.ResultNotes))
		}
		if opts.ResultStatus != nil {
			form.AddString("resultstatus", string(*opts.ResultStatus))
		}
		if opts.TieToOrderID != nil {
			form.AddInt("tietoorderid", *opts.TieToOrderID)
		}
	}

	out := &addLabResultDocumentResponse{}

	_, err := h.PostFormReader(ctx, fmt.Sprintf("patients/%s/documents/labresult", patientID), form, out)
	if err != nil {
		return 0, err
	}

	if !out.Success {
		return 0, errors.New(out.ErrorMessage)
	}

	return out.LabResultID, nil
}

type ListChangedLabResultsOptions struct {
	ShowPortalOnly             *bool
	LeaveUnprocessed           *bool
	ShowProcessedEndDateTime   time.Time
	ShowProcessedStartDateTime time.Time

	Pagination *PaginationOptions
}

type ChangedLabResult struct {
	AppointmentID        int    `json:"appointmentid"`
	AssignedTo           string `json:"assignedto"`
	CreatedDate          string `json:"createddate"`
	CreatedDateTime      string `json:"createddatetime"`
	CreatedUser          string `json:"createduser"`
	DeletedDateTime      string `json:"deleteddatetime"`
	DepartmentID         string `json:"departmentid"`
	Description          string `json:"description"`
	DocumentClass        string `json:"documentclass"`
	DocumentRoute        string `json:"documentroute"`
	DocumentSource       string `json:"documentsource"`
	DocumentSubclass     string `json:"documentsubclass"`
	EncounterDate        string `json:"encounterdate"`
	EncounterID          string `json:"encounterid"`
	ExternalNoteOnly     string `json:"externalnoteonly"`
	FacilityID           int    `json:"facilityid"`
	FileExtension        string `json:"fileextension"`
	InternalNote         string `json:"internalnote"`
	Interpretation       string `json:"interpretation"`
	IsConfidential       string `json:"isconfidential"`
	IsReviewedByProvider string `json:"isreviewedbyprovider"`
	LabResultID          int    `json:"labresultid"`
	LabResultLOINC       string `json:"labresultloinc"`
	LastModifiedDate     string `json:"lastmodifieddate"`
	LastModifiedDateTime string `json:"lastmodifieddatetime"`
	NoteFromLab          string `json:"notefromlab"`
	ObservationDate      string `json:"observationdate"`
	ObservationDateTime  string `json:"observationdatetime"`
	OrderType            string `json:"ordertype"`
	OriginalFilename     string `json:"originalfilename"`
	PatientID            int    `json:"patientid"`
	PatientNote          string `json:"patientnote"`
	Priority             string `json:"priority"`
	ProviderID           int    `json:"providerid"`
	ReportStatus         string `json:"reportstatus"`
	ResultStatus         string `json:"resultstatus"`
	Status               string `json:"status"`
	Subject              string `json:"subject"`
	TieToOrderID         int    `json:"tietoorderid"`
}

type listChangedLabResultsResponse struct {
	LabResults []*ChangedLabResult `json:"labresults"`

	*PaginationResponse
}

type ListChangedLabResultsResult struct {
	ChangedLabResults []*ChangedLabResult `json:"changedlabresults"`

	Pagination *PaginationResult
}

// ListChangedLabResults retrieves list of records of modified lab results of patient
//
// GET /v1/{practiceid}/labresults/changed
//
// https://docs.athenahealth.com/api/api-ref/document-type-lab-result#Get-list-of-changes-in-lab-results-based-on-subscription
func (h *HTTPClient) ListChangedLabResults(ctx context.Context, opts *ListChangedLabResultsOptions) (*ListChangedLabResultsResult, error) {
	q := url.Values{}

	if opts != nil {
		if opts.ShowPortalOnly != nil {
			q.Add("showportalonly", strconv.FormatBool(*opts.ShowPortalOnly))
		}
		if opts.LeaveUnprocessed != nil {
			q.Add("leaveunprocessed", strconv.FormatBool(*opts.LeaveUnprocessed))
		}
		if !opts.ShowProcessedEndDateTime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDateTime.Format("01/02/2006 15:04:05"))
		}
		if !opts.ShowProcessedStartDateTime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDateTime.Format("01/02/2006 15:04:05"))
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

	out := &listChangedLabResultsResponse{}

	_, err := h.Get(ctx, "labresults/changed", q, out)

	return &ListChangedLabResultsResult{
		ChangedLabResults: out.LabResults,
		Pagination:        makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, err
}
