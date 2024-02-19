package athenahealth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strconv"
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

// ListAdminDocuments - Get list of patient's admin documents
//
// GET /v1/{practiceid}/patients/{patientid}/documents/admin
//
// https://docs.athenahealth.com/api/api-ref/document-type-admin-document#Get-list-of-patient's-admin-documents
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

// AddDocument - Add document to patient's chart
//
// POST /v1/{practiceid}/patients/{patientid}/documents
//
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
		form = NewFormURLEncoder()

		if opts.ActionNote != nil {
			form.AddString("actionnote", *opts.ActionNote)
		}

		if opts.AppointmentID != nil {
			apptID := strconv.Itoa(*opts.AppointmentID)
			form.AddString("appointmentid", apptID)
		}

		form.AddReader("attachmentcontents", opts.AttachmentContents)

		if opts.AutoClose != nil {
			form.AddString("autoclose", *opts.AutoClose)
		}

		if opts.DepartmentID != nil {
			form.AddInt("departmentid", *opts.DepartmentID)
		}

		form.AddString("documentsubclass", opts.DocumentSubclass)

		if opts.InternalNote != nil {
			form.AddString("internalnote", *opts.InternalNote)
		}

		if opts.ProviderID != nil {
			form.AddInt("providerid", *opts.ProviderID)
		}
	}

	res := &addDocumentResponse{}

	_, err := h.PostFormReader(ctx, fmt.Sprintf("/patients/%s/documents", patientID), form, res)
	if err != nil {
		return "", err
	}

	return res.DocumentID, nil
}

type AddClinicalDocumentOptions struct {
	// The file contents that will be attached to this document. File must be Base64 encoded.
	AttachmentContents []byte
	AttachmentType     *string
	AutoClose          *string
	// The ID of the external provider/lab/pharmacy associated the document.
	ClinicalProviderID *int
	// The athenaNet department ID associated with the uploaded document. Mandatory.
	DepartmentID int
	// Text data stored with document
	DocumentData *string
	// Subclasses for CLINICALDOCUMENT documents
	DocumentSubclass string
	// A specific document type identifier.
	DocumentTypeID *int
	// Identifier of entity creating the document. entitytype is required while passing entityid.
	EntityID *int
	// Type of entity creating the document. entityid is required while passing entitytype
	EntityType *string
	// An internal note for the provider or staff. Updating this will append to any previous notes.
	InternalNote *string
	// The date an observation was made (mm/dd/yyyy).
	ObservationDate *string
	// The time an observation was made (hh24:mi). 24 hour time.
	ObservationTime *string
	// The original file name of this document without the file extension. Filename should not exceed 200 characters.
	OriginalFileName *string
	// Priority of this result. 1 is high; 2 is normal.
	Priority *string
	// The ID of the ordering provider.
	ProviderID *int
}

type AddClinicalDocumentResponse struct {
	ClinicalDocumentID int    `json:"clinicaldocumentid"`
	ErrorMessage       string `json:"errormessage"`
	Success            bool   `json:"success"`
}

// AddClinicalDocument - Add clinical document to patient's chart
//
// POST /v1/{practiceid}/patients/{patientid}/documents/clinicaldocument
//
// https://docs.athenahealth.com/api/api-ref/document-type-clinical-document#Add-clinical-document-to-patient's-chart
func (h *HTTPClient) AddClinicalDocument(ctx context.Context, patientID string, opts *AddClinicalDocumentOptions) (*AddClinicalDocumentResponse, error) {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		form.Add("attachmentcontents", base64.StdEncoding.EncodeToString(opts.AttachmentContents))

		if opts.AttachmentType != nil {
			form.Add("attachmenttype", *opts.AttachmentType)
		}

		if opts.AutoClose != nil {
			form.Add("autoclose", *opts.AutoClose)
		}

		if opts.ClinicalProviderID != nil {
			form.Add("clinicalproviderid", strconv.Itoa(*opts.ClinicalProviderID))
		}

		form.Add("departmentid", strconv.Itoa(opts.DepartmentID))

		if opts.DocumentData != nil {
			form.Add("documentdata", *opts.DocumentData)
		}

		form.Add("documentsubclass", opts.DocumentSubclass)

		if opts.DocumentTypeID != nil {
			form.Add("documenttypeid", strconv.Itoa(*opts.DocumentTypeID))
		}

		if opts.EntityID != nil {
			form.Add("entityid", strconv.Itoa(*opts.EntityID))
		}

		if opts.EntityType != nil {
			form.Add("entitytype", *opts.EntityType)
		}

		if opts.InternalNote != nil {
			form.Add("internalnote", *opts.InternalNote)
		}

		if opts.ObservationDate != nil {
			form.Add("observationdate", *opts.ObservationDate)
		}

		if opts.ObservationTime != nil {
			form.Add("observationtime", *opts.ObservationTime)
		}

		if opts.OriginalFileName != nil {
			form.Add("originalfilename", *opts.OriginalFileName)
		}

		if opts.Priority != nil {
			form.Add("priority", *opts.Priority)
		}

		if opts.ProviderID != nil {
			form.Add("providerid", strconv.Itoa(*opts.ProviderID))
		}
	}

	res := &AddClinicalDocumentResponse{}

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/documents/clinicaldocument", patientID), form, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type AddClinicalDocumentReaderOptions struct {
	// The file contents that will be attached to this document. File must be Base64 encoded.
	AttachmentContents io.Reader
	AttachmentType     *string
	AutoClose          *string
	// The ID of the external provider/lab/pharmacy associated the document.
	ClinicalProviderID *int
	// The athenaNet department ID associated with the uploaded document. Mandatory.
	DepartmentID int
	// Text data stored with document
	DocumentData *string
	// Subclasses for CLINICALDOCUMENT documents
	DocumentSubclass string
	// A specific document type identifier.
	DocumentTypeID *int
	// Identifier of entity creating the document. entitytype is required while passing entityid.
	EntityID *int
	// Type of entity creating the document. entityid is required while passing entitytype
	EntityType *string
	// An internal note for the provider or staff. Updating this will append to any previous notes.
	InternalNote *string
	// The date an observation was made (mm/dd/yyyy).
	ObservationDate *string
	// The time an observation was made (hh24:mi). 24 hour time.
	ObservationTime *string
	// The original file name of this document without the file extension. Filename should not exceed 200 characters.
	OriginalFileName *string
	// Priority of this result. 1 is high; 2 is normal.
	Priority *string
	// The ID of the ordering provider.
	ProviderID *int
}

// AddClinicalDocumentReader - performs the same operation as AddClinicalDocument except is more memory efficient
// by streaming the attachment contents into the request, assuming you haven't already read the
// entire attachment contents into memory
//
// POST /v1/{practiceid}/patients/{patientid}/documents/clinicaldocument
//
// https://docs.athenahealth.com/api/api-ref/document-type-clinical-document#Add-clinical-document-to-patient's-chart
func (h *HTTPClient) AddClinicalDocumentReader(ctx context.Context, patientID string, opts *AddClinicalDocumentReaderOptions) (*AddClinicalDocumentResponse, error) {
	var form *formURLEncoder

	if opts != nil {
		form = NewFormURLEncoder()

		form.AddReader("attachmentcontents", opts.AttachmentContents)

		if opts.AttachmentType != nil {
			form.AddString("attachmenttype", *opts.AttachmentType)
		}

		if opts.AutoClose != nil {
			form.AddString("autoclose", *opts.AutoClose)
		}

		if opts.ClinicalProviderID != nil {
			form.AddInt("clinicalproviderid", *opts.ClinicalProviderID)
		}

		form.AddInt("departmentid", opts.DepartmentID)

		if opts.DocumentData != nil {
			form.AddString("documentdata", *opts.DocumentData)
		}

		form.AddString("documentsubclass", opts.DocumentSubclass)

		if opts.DocumentTypeID != nil {
			form.AddInt("documenttypeid", *opts.DocumentTypeID)
		}

		if opts.EntityID != nil {
			form.AddInt("entityid", *opts.EntityID)
		}

		if opts.EntityType != nil {
			form.AddString("entitytype", *opts.EntityType)
		}

		if opts.InternalNote != nil {
			form.AddString("internalnote", *opts.InternalNote)
		}

		if opts.ObservationDate != nil {
			form.AddString("observationdate", *opts.ObservationDate)
		}

		if opts.ObservationTime != nil {
			form.AddString("observationtime", *opts.ObservationTime)
		}

		if opts.OriginalFileName != nil {
			form.AddString("originalfilename", *opts.OriginalFileName)
		}

		if opts.Priority != nil {
			form.AddString("priority", *opts.Priority)
		}

		if opts.ProviderID != nil {
			form.AddInt("providerid", *opts.ProviderID)
		}
	}

	res := &AddClinicalDocumentResponse{}

	_, err := h.PostFormReader(ctx, fmt.Sprintf("/patients/%s/documents/clinicaldocument", patientID), form, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type AddPatientCaseDocumentOptions struct {
	AutoClose          *bool
	CallbackName       *string
	CallbackNumber     *string
	CallbackNumberType *string
	DepartmentID       int
	DocumentSource     string
	DocumentSubclass   string
	InternalNote       *string
	OutboundOnly       *bool
	Priority           *string
	ProviderID         *int
	Subject            *string
}

type addPatientCaseDocumentResponse struct {
	PatientCaseID int `json:"patientcaseid"`
}

// AddDocument - Add patient case document for a patient
// POST /v1/{practiceid}/patients/{patientid}/documents/patientcase
// https://docs.athenahealth.com/api/api-ref/document-type-patient-case#Add-patient-case-document-for-a-patient
func (h *HTTPClient) AddPatientCaseDocument(ctx context.Context, patientID string, opts *AddPatientCaseDocumentOptions) (int, error) {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if opts.AutoClose != nil && *opts.AutoClose {
			form.Add("autoclose", "true")
		}

		if opts.CallbackName != nil {
			form.Add("callbackname", *opts.CallbackName)
		}

		if opts.CallbackNumber != nil {
			form.Add("callbacknumber", *opts.CallbackNumber)
		}

		if opts.CallbackNumberType != nil {
			form.Add("callbacknumbertype", *opts.CallbackNumberType)
		}

		deptID := strconv.Itoa(opts.DepartmentID)
		form.Add("departmentid", deptID)
		form.Add("documentsource", opts.DocumentSource)
		form.Add("documentsubclass", opts.DocumentSubclass)

		if opts.InternalNote != nil {
			form.Add("internalnote", *opts.InternalNote)
		}

		if opts.OutboundOnly != nil && *opts.OutboundOnly {
			form.Add("outboundonly", "true")
		}

		if opts.Priority != nil {
			form.Add("priority", *opts.Priority)
		}

		if opts.ProviderID != nil {
			providerID := strconv.Itoa(*opts.ProviderID)
			form.Add("providerid", providerID)
		}

		if opts.Subject != nil {
			form.Add("subject", *opts.Subject)
		}
	}

	res := &addPatientCaseDocumentResponse{}

	_, err := h.PostForm(ctx, fmt.Sprintf("/patients/%s/documents/patientcase", patientID), form, res)
	if err != nil {
		return 0, err
	}

	return res.PatientCaseID, nil
}

type DeleteClinicalDocumentResponse struct {
	ClinicalDocumentID int    `json:"clinicaldocumentid"`
	ErrorMessage       string `json:"errormessage"`
	Success            bool   `json:"success"`
}

// DeleteClinicalDocument - Mark patient's clinical document as deleted
//
// DELETE /v1/{practiceid}/patients/{patientid}/documents/clinicaldocument/{clinicaldocumentid}
//
// https://docs.athenahealth.com/api/api-ref/document-type-clinical-document#Mark-patient's-clinical-document-as-deleted
func (h *HTTPClient) DeleteClinicalDocument(ctx context.Context, patientID string, clinicalDocumentID string) (*DeleteClinicalDocumentResponse, error) {

	res := &DeleteClinicalDocumentResponse{}

	_, err := h.Delete(ctx, fmt.Sprintf("/patients/%s/documents/clinicaldocument/%s", patientID, clinicalDocumentID), nil, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type EncounterDocument struct {
	AppointmentID      int    `json:"appointmentid"`
	AssignedTo         string `json:"assignedto"`
	ClinicalProviderID int    `json:"clinicalproviderid"`
	// ContraindicationReason
	CreatedDate        string `json:"createddate"`
	CreatedDateTime    string `json:"createddatetime"`
	DeclinedReasonText string `json:"declinedreasontext"`
	DeletedDateTime    string `json:"deleteddatetime"`
	DepartmentID       string `json:"departmentid"`
	// DeclinedReason
	Description          string `json:"description"`
	DocumentClass        string `json:"documentclass"`
	DocumentDate         string `json:"documentdate"`
	DocumentRoute        string `json:"documentroute"`
	DocumentSource       string `json:"documentsource"`
	DocumentSubClass     string `json:"documentsubclass"`
	DocumentType         string `json:"documenttype"`
	DocumentTypeID       int    `json:"documenttypeid"`
	Encounterdocumentid  int    `json:"encounterdocumentid"`
	Encounterid          string `json:"encounterid"`
	Externalaccessionid  string `json:"externalaccessionid"`
	InternalNote         string `json:"internalnote"`
	LastModifiedDate     string `json:"lastmodifieddate"`
	LastModifiedDatetime string `json:"lastmodifieddatetime"`
	Lastmodifieduser     string `json:"lastmodifieduser"`
	Observationdatetime  string `json:"observationdatetime"`
	Patientid            int    `json:"patientid"`
	Priority             int    `json:"priority"`
	ProviderID           int    `json:"providerid"`
	ProviderUsername     string `json:"providerusername"`
	Receivernote         string `json:"receivernote"`
	Status               string `json:"status"`
	Subject              string `json:"subject"`
	Tietoorderid         int    `json:"tietoorderid"`
}

type ListEncounterDocumentsOptions struct {
	DocumentSubclass string
	ShowDeleted      bool
	Encounterid      string

	Pagination *PaginationOptions
}

type ListEncounterDocumentsResult struct {
	EncounterDocuments []*EncounterDocument

	Pagination *PaginationResult
}

type listEncounterDocumentsResponse struct {
	EncounterDocuments []*EncounterDocument `json:"encounterdocuments"`

	PaginationResponse
}

func (h *HTTPClient) ListEncounterDocuments(ctx context.Context, departmentID, patientID string, opts *ListEncounterDocumentsOptions) (*ListEncounterDocumentsResult, error) {
	out := &listEncounterDocumentsResponse{}

	if departmentID == "" || patientID == "" {
		return nil, fmt.Errorf("DepartmentID and patientID are required")
	}

	q := url.Values{}
	if opts != nil {
		if len(departmentID) > 0 {
			q.Add("departmentid", departmentID)
		}

		if opts.Encounterid != "" {
			q.Add("encounterid", opts.Encounterid)
		}

		if opts.DocumentSubclass != "" {
			q.Add("documentsubclass", opts.DocumentSubclass)
		}

		if opts.ShowDeleted {
			q.Add("showdeleted", "true")
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

	_, err := h.Get(ctx, fmt.Sprintf("/patients/%s/documents/encounterdocument", patientID), q, out)
	if err != nil {
		return nil, err
	}

	return &ListEncounterDocumentsResult{
		EncounterDocuments: out.EncounterDocuments,
		Pagination:         makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}
