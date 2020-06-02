package athenahealth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// Patient represents a patient in athenahealth.
type Patient struct {
	RaceName                           string `json:"racename"`
	DoNotCall                          bool   `json:"donotcall"`
	Email                              string `json:"email"`
	DepartmentID                       string `json:"departmentid"`
	HomePhone                          string `json:"homephone"`
	GuarantorState                     string `json:"guarantorstate"`
	PortalAccessGiven                  bool   `json:"portalaccessgiven"`
	DriversLicense                     bool   `json:"driverslicense"`
	ContactPreferenceAppointmentEmail  bool   `json:"contactpreference_appointment_email"`
	Homebound                          bool   `json:"homebound"`
	ContactPreferenceAppointmentSMS    bool   `json:"contactpreference_appointment_sms"`
	ContactPreferenceBillingPhone      bool   `json:"contactpreference_billing_phone"`
	EthnicityCode                      string `json:"ethnicitycode"`
	ContactPreferenceAnnouncementPhone bool   `json:"contactpreference_announcement_phone"`
	LastEmail                          string `json:"lastemail"`
	ContactRelationship                string `json:"contactrelationship"`
	EmployerID                         string `json:"employerid"`
	ContactHomePhone                   string `json:"contacthomephone"`
	GuarantorSSN                       string `json:"guarantorssn"`
	ContactOreferenceLabSMS            bool   `json:"contactpreference_lab_sms"`
	GuarantorDOB                       string `json:"guarantordob"`
	Zip                                string `json:"zip"`
	GuarantorAddressSameAsPatient      bool   `json:"guarantoraddresssameaspatient"`
	EmployerPhone                      string `json:"employerphone"`
	PortalTermsOnFile                  bool   `json:"portaltermsonfile"`
	Status                             string `json:"status"`
	LastName                           string `json:"lastname"`
	GuarantorFirstName                 string `json:"guarantorfirstname"`
	City                               string `json:"city"`
	SSN                                string `json:"ssn"`
	LastAppointment                    string `json:"lastappointment"`
	PovertyLevelIncomeDeclined         bool   `json:"povertylevelincomedeclined"`
	GuarantorZmail                     string `json:"guarantoremail"`
	GuarantorCity                      string `json:"guarantorcity"`
	Guarantorzip                       string `json:"guarantorzip"`
	Sex                                string `json:"sex"`
	PrivacyInformationVerified         bool   `json:"privacyinformationverified"`
	PrimaryDepartmentID                string `json:"primarydepartmentid"`
	ContactPreferenceLabEmail          bool   `json:"contactpreference_lab_email"`
	Balances                           []struct {
		Balance         json.Number `json:"balance"`
		DepartmentList  string      `json:"departmentlist"`
		ProviderGroupID int         `json:"providergroupid"`
		CleanBalance    bool        `json:"cleanbalance"`
	} `json:"balances"`
	ContactpreferenceAnnouncementSMS   bool     `json:"contactpreference_announcement_sms"`
	Race                               []string `json:"race"`
	EmployerState                      string   `json:"employerstate"`
	FirstAppointment                   string   `json:"firstappointment"`
	Language6392Code                   string   `json:"language6392code"`
	PrimaryProviderID                  string   `json:"primaryproviderid"`
	PatientPhoto                       bool     `json:"patientphoto"`
	ConsentToCall                      bool     `json:"consenttocall"`
	DefaultPharmacyNCPDPID             string   `json:"defaultpharmacyncpdpid"`
	PovertyLevelIncomeRangeDeclined    bool     `json:"povertylevelincomerangedeclined"`
	ContactPreferenceBillingEmail      bool     `json:"contactpreference_billing_email"`
	PatientPhotoURL                    string   `json:"patientphotourl"`
	ContactName                        string   `json:"contactname"`
	MobilePhone                        string   `json:"mobilephone"`
	ContactPreferenceAnnouncementEmail bool     `json:"contactpreference_announcement_email"`
	HasMobile                          bool     `json:"hasmobile"`
	RegistrationDate                   string   `json:"registrationdate"`
	CareSummaryDeliveryPreference      string   `json:"caresummarydeliverypreference"`
	GuarantorLastName                  string   `json:"guarantorlastname"`
	FirstName                          string   `json:"firstname"`
	GuarantorCountryCode               string   `json:"guarantorcountrycode"`
	State                              string   `json:"state"`
	ContactPreferenceAppointmentPhone  bool     `json:"contactpreference_appointment_phone"`
	PatientID                          string   `json:"patientid"`
	DOB                                string   `json:"dob"`
	GuarantorRelationshipToPatient     string   `json:"guarantorrelationshiptopatient"`
	Address1                           string   `json:"address1"`
	GuarantorEmployerID                string   `json:"guarantoremployerid"`
	ContactPreferenceBillingSMS        bool     `json:"contactpreference_billing_sms"`
	GuarantorPhone                     string   `json:"guarantorphone"`
	PovertyLevelFamilySizeDeclined     bool     `json:"povertylevelfamilysizedeclined"`
	EmployerName                       string   `json:"employername"`
	EmployerAddress                    string   `json:"employeraddress"`
	MaritalStatus                      string   `json:"maritalstatus"`
	CountryCode                        string   `json:"countrycode"`
	GuarantorAddress1                  string   `json:"guarantoraddress1"`
	MaritalStatusName                  string   `json:"maritalstatusname"`
	GuarantorMiddleName                string   `json:"guarantormiddlename"`
	ConsentToText                      bool     `json:"consenttotext"`
	CountryCode3166                    string   `json:"countrycode3166"`
	OnlineStatementOnly                bool     `json:"onlinestatementonly"`
	ContactPreferenceLabPhone          bool     `json:"contactpreference_lab_phone"`
	GuarantorCountryCode3166           string   `json:"guarantorcountrycode3166"`
}

// GetPatient - Full view/update of patient demographics.
// GET /v1/{practiceid}/patients/{patientid}
// https://developer.athenahealth.com/docs/read/patientinfo/Patient_Information#section-5
func (h *HTTPClient) GetPatient(id string) (*Patient, error) {
	out := []*Patient{}

	_, err := h.Get(fmt.Sprintf("/patients/%s", id), nil, &out)
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
