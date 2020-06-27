package athenahealth

import (
	"errors"
	"fmt"
)

type Department struct {
	MedicationHistoryConsent bool   `json:"medicationhistoryconsent"`
	TimeZoneOffset           int    `json:"timezoneoffset"`
	IsHospitalDepartment     bool   `json:"ishospitaldepartment"`
	ProviderGroupID          string `json:"providergroupid"`
	State                    string `json:"state"`
	PortalURL                string `json:"portalurl"`
	City                     string `json:"city"`
	PlaceOfServiceFacility   bool   `json:"placeofservicefacility"`
	ServiceDepartment        bool   `json:"servicedepartment"`
	ProviderGroupName        string `json:"providergroupname"`
	DoesNotObserveDST        bool   `json:"doesnotobservedst"`
	DepartmentID             string `json:"departmentid"`
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
