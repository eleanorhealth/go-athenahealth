package athenahealth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// Patient represents a patient in athenahealth.
type Patient struct {
	Address1 string `json:"address1"`
	Balances []struct {
		Balance         json.Number `json:"balance"`
		DepartmentList  string      `json:"departmentlist"`
		ProviderGroupID int         `json:"providergroupid"`
		CleanBalance    bool        `json:"cleanbalance"`
	} `json:"balances"`
	CareSummaryDeliveryPreference      string             `json:"caresummarydeliverypreference"`
	City                               string             `json:"city"`
	ConsentToCall                      bool               `json:"consenttocall"`
	ConsentToText                      bool               `json:"consenttotext"`
	ContactHomePhone                   string             `json:"contacthomephone"`
	ContactName                        string             `json:"contactname"`
	ContactPreference                  string             `json:"contactpreference"`
	ContactPreferenceLabSMS            bool               `json:"contactpreference_lab_sms"`
	ContactPreferenceAnnouncementEmail bool               `json:"contactpreference_announcement_email"`
	ContactPreferenceAnnouncementPhone bool               `json:"contactpreference_announcement_phone"`
	ContactpreferenceAnnouncementSMS   bool               `json:"contactpreference_announcement_sms"`
	ContactPreferenceAppointmentEmail  bool               `json:"contactpreference_appointment_email"`
	ContactPreferenceAppointmentPhone  bool               `json:"contactpreference_appointment_phone"`
	ContactPreferenceAppointmentSMS    bool               `json:"contactpreference_appointment_sms"`
	ContactPreferenceBillingEmail      bool               `json:"contactpreference_billing_email"`
	ContactPreferenceBillingPhone      bool               `json:"contactpreference_billing_phone"`
	ContactPreferenceBillingSMS        bool               `json:"contactpreference_billing_sms"`
	ContactPreferenceLabEmail          bool               `json:"contactpreference_lab_email"`
	ContactPreferenceLabPhone          bool               `json:"contactpreference_lab_phone"`
	ContactRelationship                string             `json:"contactrelationship"`
	CountryCode                        string             `json:"countrycode"`
	CountryCode3166                    string             `json:"countrycode3166"`
	CustomFields                       []CustomFieldValue `json:"customfields"`
	DefaultPharmacyNCPDPID             string             `json:"defaultpharmacyncpdpid"`
	DepartmentID                       string             `json:"departmentid"`
	DOB                                string             `json:"dob"`
	DoNotCall                          bool               `json:"donotcall"`
	DriversLicense                     bool               `json:"driverslicense"`
	Email                              string             `json:"email"`
	EmployerAddress                    string             `json:"employeraddress"`
	EmployerID                         string             `json:"employerid"`
	EmployerName                       string             `json:"employername"`
	EmployerPhone                      string             `json:"employerphone"`
	EmployerState                      string             `json:"employerstate"`
	EthnicityCode                      string             `json:"ethnicitycode"`
	FirstAppointment                   string             `json:"firstappointment"`
	FirstName                          string             `json:"firstname"`
	GuarantorAddress1                  string             `json:"guarantoraddress1"`
	GuarantorAddressSameAsPatient      bool               `json:"guarantoraddresssameaspatient"`
	GuarantorCity                      string             `json:"guarantorcity"`
	GuarantorCountryCode               string             `json:"guarantorcountrycode"`
	GuarantorCountryCode3166           string             `json:"guarantorcountrycode3166"`
	GuarantorDOB                       string             `json:"guarantordob"`
	GuarantorEmployerID                string             `json:"guarantoremployerid"`
	GuarantorFirstName                 string             `json:"guarantorfirstname"`
	GuarantorLastName                  string             `json:"guarantorlastname"`
	GuarantorMiddleName                string             `json:"guarantormiddlename"`
	GuarantorPhone                     string             `json:"guarantorphone"`
	GuarantorRelationshipToPatient     string             `json:"guarantorrelationshiptopatient"`
	GuarantorSSN                       string             `json:"guarantorssn"`
	GuarantorState                     string             `json:"guarantorstate"`
	Guarantorzip                       string             `json:"guarantorzip"`
	GuarantorZmail                     string             `json:"guarantoremail"`
	HasMobile                          bool               `json:"hasmobile"`
	Homebound                          bool               `json:"homebound"`
	HomePhone                          string             `json:"homephone"`
	Language6392Code                   string             `json:"language6392code"`
	LastAppointment                    string             `json:"lastappointment"`
	LastEmail                          string             `json:"lastemail"`
	LastName                           string             `json:"lastname"`
	MaritalStatus                      string             `json:"maritalstatus"`
	MaritalStatusName                  string             `json:"maritalstatusname"`
	MobilePhone                        string             `json:"mobilephone"`
	OnlineStatementOnly                bool               `json:"onlinestatementonly"`
	PatientID                          string             `json:"patientid"`
	PatientPhoto                       bool               `json:"patientphoto"`
	PatientPhotoURL                    string             `json:"patientphotourl"`
	PortalAccessGiven                  bool               `json:"portalaccessgiven"`
	PortalTermsOnFile                  bool               `json:"portaltermsonfile"`
	PovertyLevelFamilySizeDeclined     bool               `json:"povertylevelfamilysizedeclined"`
	PovertyLevelIncomeDeclined         bool               `json:"povertylevelincomedeclined"`
	PovertyLevelIncomeRangeDeclined    bool               `json:"povertylevelincomerangedeclined"`
	PrimaryDepartmentID                string             `json:"primarydepartmentid"`
	PrimaryProviderID                  string             `json:"primaryproviderid"`
	PrivacyInformationVerified         bool               `json:"privacyinformationverified"`
	Race                               []string           `json:"race"`
	RaceName                           string             `json:"racename"`
	RegistrationDate                   string             `json:"registrationdate"`
	Sex                                string             `json:"sex"`
	SSN                                string             `json:"ssn"`
	State                              string             `json:"state"`
	Status                             string             `json:"status"`
	Zip                                string             `json:"zip"`
}

type GetPatientOptions struct {
	ShowCustomFields bool
}

// GetPatient - Full view/update of patient demographics.
// GET /v1/{practiceid}/patients/{patientid}
// https://developer.athenahealth.com/docs/read/patientinfo/Patient_Information#section-5
func (h *HTTPClient) GetPatient(id string, opts *GetPatientOptions) (*Patient, error) {
	out := []*Patient{}

	q := url.Values{}

	if opts != nil {
		if opts.ShowCustomFields {
			q.Add("showcustomfields", "true")
		}
	}

	_, err := h.Get(fmt.Sprintf("/patients/%s", id), q, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type ListPatientsOptions struct {
	FirstName string
	LastName  string
}

type listPatientsResponse struct {
	Patients []*Patient `json:"patients"`
}

// ListPatients - Gets a set of patients or creates a patient.
// GET /v1/{practiceid}/patients
// https://developer.athenahealth.com/docs/read/patientinfo/Patient_Information#section-1
func (h *HTTPClient) ListPatients(opts *ListPatientsOptions) ([]*Patient, error) {
	out := &listPatientsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.FirstName) > 0 {
			q.Add("firstname", opts.FirstName)
		}

		if len(opts.LastName) > 0 {
			q.Add("lastname", opts.LastName)
		}
	}

	_, err := h.Get("/patients", q, out)
	if err != nil {
		return nil, err
	}

	return out.Patients, nil
}
