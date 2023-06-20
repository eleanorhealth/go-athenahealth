package athenahealth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// AdminDocument represents an administrative document in athenahealth.
type AdminDocument struct {
	Priority             string `json:"priority"`
	AssignedTo           string `json:"assignedto"`
	DocumentClass        string `json:"documentclass"`
	CreatedDateTime      string `json:"createddatetime"`
	DepartmentID         string `json:"departmentid"`
	DocumentTypeID       int    `json:"documenttypeid"`
	InternalNote         string `json:"internalnote"`
	AdminID              int    `json:"adminid"`
	CreatedUser          string `json:"createduser"`
	Description          string `json:"description"`
	DocumentDate         string `json:"documentdate"`
	DocumentRoute        string `json:"documentroute"`
	DocumentSource       string `json:"documentsource"`
	CreatedDate          string `json:"createddate"`
	Status               string `json:"status"`
	ProviderID           int    `json:"providerid"`
	ProviderUsername     string `json:"providerusername"`
	LastModifiedDatetime string `json:"lastmodifieddatetime"`
	LastModifiedDate     string `json:"lastmodifieddate"`
}

type ListAdminDocumentsOptions struct {
	DepartmentID string

	Pagination *PaginationOptions
}

type ListAdminDocumentsResult struct {
	AdminDocuments []*AdminDocument

	Pagination *PaginationResult
}

type listAdminDocumentsResponse struct {
	AdminDocuments []*AdminDocument `json:"admins"`

	PaginationResponse
}

// ListAdminDocuments - Get admin documents.
// GET /v1/{practiceid}/patients/{patientid}/documents/admin
// https://developer.athenahealth.com/docs/read/forms_and_documents/Document_Lists_By_Class#section-19
func (h *HTTPClient) ListAdminDocuments(ctx context.Context, patientID string, opts *ListAdminDocumentsOptions) (*ListAdminDocumentsResult, error) {
	out := &listAdminDocumentsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
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

	_, err := h.Get(ctx, fmt.Sprintf("/patients/%s/documents/admin", patientID), q, out)
	if err != nil {
		return nil, err
	}

	return &ListAdminDocumentsResult{
		AdminDocuments: out.AdminDocuments,
		Pagination:     makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type AddDocumentOptions struct {
	ActionNote         *string
	AppointmentID      *int
	AttachmentContents []byte
	AutoClose          *string
	DepartmentID       *int
	DocumentSubclass   string
	InternalNote       *string
	ProviderID         *int
}

type addDocumentResponse struct {
	DocumentID string `json:"documentid"`
}

// AddDocument - Add document to patient's chart.
// POST /v1/{practiceid}/patients/{patientid}/documents
// https://docs.athenahealth.com/api/api-ref/document#Add-document-to-patient's-chart
// Document subclasses from https://docs.athenahealth.com/api/workflows/document-classification-guide:
// ADMIN_BILLING
// ADMIN_CONSENT
// ADMIN_HIPAA
// ADMIN_INSURANCEAPPROVAL
// ADMIN_INSURANCECARD
// ADMIN_INSURANCEDENIAL
// ADMIN_LEGAL
// ADMIN_MEDICALRECORDREQ
// ADMIN_REFERRAL
// ADMIN_SIGNEDFORMSLETTERS
// ADMIN_VACCINATIONRECORD
// CLINICALDOCUMENT_ADMISSIONDISCHARGE
// CLINICALDOCUMENT_CONSULTNOTE
// CLINICALDOCUMENT_MENTALHEALTH
// CLINICALDOCUMENT_OPERATIVENOTE
// CLINICALDOCUMENT_URGENTCARE
// ENCOUNTERDOCUMENT_IMAGEDOC
// ENCOUNTERDOCUMENT_PATIENTHISTORY
// ENCOUNTERDOCUMENT_PROCEDUREDOC
// ENCOUNTERDOCUMENT_PROGRESSNOTE
// MEDICALRECORD_CHARTTOABSTRACT
// MEDICALRECORD_COUMADIN
// MEDICALRECORD_GROWTHCHART
// MEDICALRECORD_HISTORICAL
// MEDICALRECORD_PATIENTDIARY
// MEDICALRECORD_VACCINATION
func (h *HTTPClient) AddDocument(ctx context.Context, patientID string, opts *AddDocumentOptions) (string, error) {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if opts.ActionNote != nil {
			form.Add("actionnote", *opts.ActionNote)
		}

		if opts.AppointmentID != nil {
			apptID := strconv.Itoa(*opts.AppointmentID)
			form.Add("appointmentid", apptID)
		}

		form.Add("attachmentcontents", base64.StdEncoding.EncodeToString(opts.AttachmentContents))

		if opts.AutoClose != nil {
			form.Add("autoclose", *opts.AutoClose)
		}

		if opts.DepartmentID != nil {
			deptID := strconv.Itoa(*opts.DepartmentID)
			form.Add("departmentid", deptID)
		}

		form.Add("documentsubclass", opts.DocumentSubclass)

		if opts.InternalNote != nil {
			form.Add("internalnote", *opts.InternalNote)
		}

		if opts.ProviderID != nil {
			providerID := strconv.Itoa(*opts.ProviderID)
			form.Add("providerid", providerID)
		}
	}

	res := &addDocumentResponse{}

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/documents", patientID), form, res)
	if err != nil {
		return "", err
	}

	return res.DocumentID, nil
}

type AddDocumentReaderOptions struct {
	ActionNote         *string
	AppointmentID      *int
	AttachmentContents io.Reader
	AutoClose          *string
	DepartmentID       *int
	DocumentSubclass   string
	InternalNote       *string
	ProviderID         *int
}

type addDocumentReaderResponse struct {
	DocumentID string `json:"documentid"`
}

// AddDocumentReader - performs the same operation as AddDocument except is more memory efficient
// by streaming the attachment contents into the request, assuming you haven't already read the
// entire attachment contents into memory
// POST /v1/{practiceid}/patients/{patientid}/documents
// https://docs.athenahealth.com/api/api-ref/document#Add-document-to-patient's-chart
// Document subclasses from https://docs.athenahealth.com/api/workflows/document-classification-guide:
// ADMIN_BILLING
// ADMIN_CONSENT
// ADMIN_HIPAA
// ADMIN_INSURANCEAPPROVAL
// ADMIN_INSURANCECARD
// ADMIN_INSURANCEDENIAL
// ADMIN_LEGAL
// ADMIN_MEDICALRECORDREQ
// ADMIN_REFERRAL
// ADMIN_SIGNEDFORMSLETTERS
// ADMIN_VACCINATIONRECORD
// CLINICALDOCUMENT_ADMISSIONDISCHARGE
// CLINICALDOCUMENT_CONSULTNOTE
// CLINICALDOCUMENT_MENTALHEALTH
// CLINICALDOCUMENT_OPERATIVENOTE
// CLINICALDOCUMENT_URGENTCARE
// ENCOUNTERDOCUMENT_IMAGEDOC
// ENCOUNTERDOCUMENT_PATIENTHISTORY
// ENCOUNTERDOCUMENT_PROCEDUREDOC
// ENCOUNTERDOCUMENT_PROGRESSNOTE
// MEDICALRECORD_CHARTTOABSTRACT
// MEDICALRECORD_COUMADIN
// MEDICALRECORD_GROWTHCHART
// MEDICALRECORD_HISTORICAL
// MEDICALRECORD_PATIENTDIARY
// MEDICALRECORD_VACCINATION
func (h *HTTPClient) AddDocumentReader(ctx context.Context, patientID string, opts *AddDocumentReaderOptions) (string, error) {
	var form *formURLEncoder

	if opts != nil {
		form = newFormURLEncoder()

		if opts.ActionNote != nil {
			form.Add("actionnote", strings.NewReader(*opts.ActionNote))
		}

		if opts.AppointmentID != nil {
			apptID := strconv.Itoa(*opts.AppointmentID)
			form.Add("appointmentid", strings.NewReader(apptID))
		}

		form.Add("attachmentcontents", newBase64Reader(opts.AttachmentContents))

		if opts.AutoClose != nil {
			form.Add("autoclose", strings.NewReader(*opts.AutoClose))
		}

		if opts.DepartmentID != nil {
			deptID := strconv.Itoa(*opts.DepartmentID)
			form.Add("departmentid", strings.NewReader(deptID))
		}

		form.Add("documentsubclass", strings.NewReader(opts.DocumentSubclass))

		if opts.InternalNote != nil {
			form.Add("internalnote", strings.NewReader(*opts.InternalNote))
		}

		if opts.ProviderID != nil {
			providerID := strconv.Itoa(*opts.ProviderID)
			form.Add("providerid", strings.NewReader(providerID))
		}
	}

	res := &addDocumentReaderResponse{}

	_, err := h.PostFormReader(ctx, fmt.Sprintf("/patients/%s/documents", patientID), form, res)
	if err != nil {
		return "", err
	}

	return res.DocumentID, nil
}
