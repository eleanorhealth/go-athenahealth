package athenahealth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Patient represents a patient in athenahealth.
type Patient struct {
	Address1 string `json:"address1"`
	Balances []struct {
		Balance         NumberString `json:"balance"`
		DepartmentList  string       `json:"departmentlist"`
		ProviderGroupID int          `json:"providergroupid"`
		CleanBalance    bool         `json:"cleanbalance"`
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
	Insurances                         []Insurance        `json:"insurances"`
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
	PortalStatus                       PortalStatus       `json:"portalstatus"`
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

type Insurance struct {
	EligibilityLastChecked              string `json:"eligibilitylastchecked"`
	EligibilityReason                   string `json:"eligibilityreason"`
	EligibilityStatus                   string `json:"eligibilitystatus"`
	ID                                  string `json:"id"`
	InsuranceID                         string `json:"insuranceid"`
	InsuranceIDNumber                   string `json:"insuranceidnumber"`
	InsurancePackageAddress1            string `json:"insurancepackageaddress1"`
	InsurancePackageCity                string `json:"insurancepackagecity"`
	InsurancePackageID                  int    `json:"insurancepackageid"`
	InsurancePackageState               string `json:"insurancepackagestate"`
	InsurancePackageZip                 string `json:"insurancepackagezip"`
	InsurancePhone                      string `json:"insurancephone"`
	InsurancePlanDisplayName            string `json:"insuranceplandisplayname"`
	InsurancePlanName                   string `json:"insuranceplanname"`
	InsurancePolicyHolder               string `json:"insurancepolicyholder"`
	InsurancePolicyHolderAddress1       string `json:"insurancepolicyholderaddress1"`
	InsurancePolicyHolderCity           string `json:"insurancepolicyholdercity"`
	InsurancePolicyHolderCountryCode    string `json:"insurancepolicyholdercountrycode"`
	InsurancePolicyHolderCountryISO3166 string `json:"insurancepolicyholdercountryiso3166"`
	InsurancePolicyHolderDOB            string `json:"insurancepolicyholderdob"`
	InsurancePolicyHolderFirstName      string `json:"insurancepolicyholderfirstname"`
	InsurancePolicyHolderLastName       string `json:"insurancepolicyholderlastname"`
	InsurancePolicyHolderSex            string `json:"insurancepolicyholdersex"`
	InsurancePolicyHolderState          string `json:"insurancepolicyholderstate"`
	InsurancePolicyHolderZip            string `json:"insurancepolicyholderzip"`
	InsuranceType                       string `json:"insurancetype"`
	InsuredAddress                      string `json:"insuredaddress"`
	InsuredCity                         string `json:"insuredcity"`
	InsuredCountryCode                  string `json:"insuredcountrycode"`
	InsuredCountryISO3166               string `json:"insuredcountryiso3166"`
	InsuredDOB                          string `json:"insureddob"`
	InsuredEntityTypeID                 int    `json:"insuredentitytypeid"`
	InsuredFirstName                    string `json:"insuredfirstname"`
	InsuredLastName                     string `json:"insuredlastname"`
	InsuredSex                          string `json:"insuredsex"`
	InsuredState                        string `json:"insuredstate"`
	InsuredZip                          string `json:"insuredzip"`
	IRCName                             string `json:"ircname"`
	RelationshipToInsured               string `json:"relationshiptoinsured"`
	RelationshipToInsuredID             int    `json:"relationshiptoinsuredid"`
	SequenceNumber                      int    `json:"sequencenumber"`
}

type PortalStatus struct {
	BlockedFailedLogins       bool   `json:"blockedfailedlogins"`
	EntityToDisplay           string `json:"entitytodisplay"`
	FamilyBlockedFailedLogins bool   `json:"familyblockedfailedlogins"`
	FamilyRegistered          bool   `json:"familyregistered"`
	NoPortal                  bool   `json:"noportal"`
	PortalRegistrationDate    string `json:"portalregistrationdate"`
	Registered                bool   `json:"registered"`
	Status                    string `json:"status"`
	TermsAccepted             bool   `json:"termsaccepted"`
}

type GetPatientOptions struct {
	ShowCustomFields bool
	ShowInsurance    bool
	ShowPortalStatus bool
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

		if opts.ShowInsurance {
			q.Add("showinsurance", "true")
		}

		if opts.ShowPortalStatus {
			q.Add("showportalstatus", "true")
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
	FirstName    string
	LastName     string
	DepartmentID int
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

		if opts.DepartmentID > 0 {
			q.Add("departmentid", strconv.Itoa(opts.DepartmentID))
		}
	}

	_, err := h.Get("/patients", q, out)
	if err != nil {
		return nil, err
	}

	return out.Patients, nil
}

type GetPatientPhotoOptions struct {
	JPEGOutput bool
}

type patientPhoto struct {
	Image string `json:"image"`
}

// GetPatientPhoto - Get a patient's photo.
// GET /v1/{practiceid}/patients/{patientid}/photo
// https://developer.athenahealth.com/docs/read/forms_and_documents/Patient_Photo#section-0
func (h *HTTPClient) GetPatientPhoto(patientID string, opts *GetPatientPhotoOptions) (string, error) {
	out := &patientPhoto{}

	q := url.Values{}

	if opts != nil {
		if opts.JPEGOutput {
			return "", errors.New("JPEGOutput is not supported")
		}
	}

	_, err := h.Get(fmt.Sprintf("/patients/%s/photo", patientID), q, &out)
	if err != nil {
		return "", err
	}

	return out.Image, nil
}

// UpdatePatientPhoto - Update a patient's photo.
// POST /v1/{practiceid}/patients/{patientid}/photo
// https://developer.athenahealth.com/docs/read/forms_and_documents/Patient_Photo#section-1
func (h *HTTPClient) UpdatePatientPhoto(patientID string, data []byte) error {
	form := url.Values{}
	form.Add("image", base64.StdEncoding.EncodeToString(data))

	_, err := h.PostForm(fmt.Sprintf("/patients/%s/photo", patientID), form, nil)
	return err
}

type ListChangedPatientOptions struct {
	DepartmentID               string
	IgnoreRestrictions         bool
	LeaveUnprocessed           bool
	PatientID                  string
	ReturnGlobalID             bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedPatientsResponse struct {
	ChangedPatients []*Patient `json:"patients"`
}

// ListChangedPatients - Gets changed patient records.
// GET /v1/{practiceid}/patients/changed
// https://developer.athenahealth.com/docs/read/patientinfo/Patients_Changed
func (h *HTTPClient) ListChangedPatients(opts *ListChangedPatientOptions) ([]*Patient, error) {
	out := &listChangedPatientsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if opts.IgnoreRestrictions {
			q.Add("ignorerestrictions", strconv.FormatBool(opts.IgnoreRestrictions))
		}

		if opts.LeaveUnprocessed {
			q.Add("leaveunprocessed", strconv.FormatBool(opts.LeaveUnprocessed))
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if opts.ReturnGlobalID {
			q.Add("returnglobalid", strconv.FormatBool(opts.ReturnGlobalID))
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
		}
	}

	_, err := h.Get("/patients/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedPatients, nil
}
