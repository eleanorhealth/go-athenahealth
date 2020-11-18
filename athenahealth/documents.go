package athenahealth

import (
	"fmt"
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

// ListAdminDocuments - Get admin documents.
// GET /v1/{practiceid}/patients/{patientid}/documents/admin
// https://developer.athenahealth.com/docs/read/forms_and_documents/Document_Lists_By_Class#section-19
func (h *HTTPClient) ListAdminDocuments(patientID string, opts *ListAdminDocumentsOptions) (*ListAdminDocumentsResult, error) {
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

	_, err := h.Get(fmt.Sprintf("/patients/%s/documents/admin", patientID), q, out)
	if err != nil {
		return nil, err
	}

	return &ListAdminDocumentsResult{
		AdminDocuments: out.AdminDocuments,
		Pagination:     makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}
