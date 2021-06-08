package athenahealth

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Department struct {
	MedicationHistoryConsent bool   `json:"medicationhistoryconsent"`
	TimeZoneOffset           int    `json:"timezoneoffset"`
	IsHospitalDepartment     bool   `json:"ishospitaldepartment"`
	ProviderGroupID          string `json:"providergroupid"`
	State                    string `json:"state"`
	PortalURL                string `json:"portalurl"`
	City                     string `json:"city"`
	ClinicalProviderFax      string `json:"clincalproviderfax"`
	PlaceOfServiceFacility   bool   `json:"placeofservicefacility"`
	ServiceDepartment        bool   `json:"servicedepartment"`
	ProviderGroupName        string `json:"providergroupname"`
	DoesNotObserveDST        bool   `json:"doesnotobservedst"`
	DepartmentID             string `json:"departmentid"`
	Fax                      string `json:"fax"`
	Address                  string `json:"address"`
	PlaceOfServiceTypeID     string `json:"placeofservicetypeid"`
	Clinicals                string `json:"clinicals"`
	TimeZone                 int    `json:"timezone"`
	PatientDepartmentName    string `json:"patientdepartmentname"`
	ChartSharingGroupID      string `json:"chartsharinggroupid"`
	Name                     string `json:"name"`
	PlaceOfServiceTypeName   string `json:"placeofservicetypename"`
	Phone                    string `json:"phone"`
	Address2                 string `json:"address2"`
	Zip                      string `json:"zip"`
	TimeZoneName             string `json:"timezonename"`
	CommunicatorBrandID      string `json:"communicatorbrandid"`
}

// GetDepartment - Details about a single department.
// GET /v1/{practiceid}/departments/{departmentid}
// https://developer.athenahealth.com/docs/read/administrative/Departments#section-1
func (h *HTTPClient) GetDepartment(id string) (*Department, error) {
	out := []*Department{}

	_, err := h.Get(fmt.Sprintf("/departments/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type ListDepartmentsOptions struct {
	HospitalOnly       bool
	ProviderList       bool
	ShowAllDepartments bool

	Pagination *PaginationOptions
}

type ListDepartmentsResult struct {
	Departments []*Department

	Pagination *PaginationResult
}

type listDepartmentsResponse struct {
	Departments []*Department `json:"departments"`

	PaginationResponse
}

// ListDepartments - List of all departments available for this practice
// GET /v1/{practiceid}/departments
// https://developer.athenahealth.com/docs/read/administrative/Departments#section-0
func (h *HTTPClient) ListDepartments(opts *ListDepartmentsOptions) (*ListDepartmentsResult, error) {
	out := &listDepartmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.HospitalOnly {
			q.Add("hospitalonly", "1")
		}

		if opts.ProviderList {
			q.Add("providerlist", "1")
		}

		if opts.ProviderList {
			q.Add("showalldepartments", "1")
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

	_, err := h.Get("/departments", q, out)
	if err != nil {
		return nil, err
	}

	return &ListDepartmentsResult{
		Departments: out.Departments,
		Pagination:  makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}
