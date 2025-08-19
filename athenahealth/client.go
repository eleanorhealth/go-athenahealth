package athenahealth

import (
	"context"
	"io"
	"time"
)

// Client describes a client for the athenahealth API.
type Client interface {
	// Department
	DepartmentGetRequiredCheckInFields(ctx context.Context, deptID string) (*GetRequiredCheckInFieldsResult, error)
	GetDepartment(ctx context.Context, departmentID string) (*Department, error)
	ListDepartments(context.Context, *ListDepartmentsOptions) (*ListDepartmentsResult, error)

	// Patient
	CreatePatient(ctx context.Context, opts *CreatePatientOptions) (string, error)
	GetPatient(ctx context.Context, patientID string, opts *GetPatientOptions) (*Patient, error)
	GetPatients(ctx context.Context, id string, opts *GetPatientOptions) ([]*Patient, error)
	ListPatients(context.Context, *ListPatientsOptions) (*ListPatientsResult, error)
	UpdatePatient(ctx context.Context, patientID string, opts *UpdatePatientOptions) (*UpdatePatientResult, error)
	UpdatePatientInformationVerificationDetails(ctx context.Context, patientID string, opts *UpdatePatientInformationVerificationDetailsOptions) error
	UpdatePatientMedicationHistoryConsent(ctx context.Context, patientID string, opts *UpdatePatientMedicationHistoryConsentOptions) error

	// Patient Photo
	GetPatientPhoto(ctx context.Context, patientID string, opts *GetPatientPhotoOptions) (string, error)
	UpdatePatientPhoto(ctx context.Context, patientID string, data []byte) error
	UpdatePatientPhotoReader(ctx context.Context, patientID string, r io.Reader) error

	// Patient Problems
	ListProblems(ctx context.Context, patientID string, opts *ListProblemsOptions) ([]*Problem, error)

	// Patient Documents
	ListAdminDocuments(ctx context.Context, patientID string, opts *ListAdminDocumentsOptions) (*ListAdminDocumentsResult, error)
	AddDocument(ctx context.Context, patientID string, opts *AddDocumentOptions) (string, error)
	AddDocumentReader(ctx context.Context, patientID string, opts *AddDocumentReaderOptions) (string, error)
	AddClinicalDocument(ctx context.Context, patientID string, opts *AddClinicalDocumentOptions) (*AddClinicalDocumentResponse, error)
	AddClinicalDocumentReader(ctx context.Context, patientID string, opts *AddClinicalDocumentReaderOptions) (*AddClinicalDocumentResponse, error)
	AddPatientCaseDocument(ctx context.Context, patientID string, opts *AddPatientCaseDocumentOptions) (int, error)
	DeleteClinicalDocument(ctx context.Context, patientID string, clinicalDocumentID string) (*DeleteClinicalDocumentResponse, error)

	// Patient Custom Fields
	ListPatientsMatchingCustomField(ctx context.Context, opts *ListPatientsMatchingCustomFieldOptions) (*ListPatientsMatchingCustomFieldResult, error)
	ListCustomFields(ctx context.Context) ([]*CustomField, error)
	GetPatientCustomFields(ctx context.Context, patientID, departmentID string) ([]*CustomFieldValue, error)
	UpdatePatientCustomFields(ctx context.Context, patientID, departmentID string, customFields []*CustomFieldValue) error

	// Patient Insurance
	CreatePatientInsurancePackage(ctx context.Context, opts *CreatePatientInsurancePackageOptions) (*InsurancePackage, error)
	DeletePatientInsurancePackage(ctx context.Context, patientID, insuranceID, cancellationNote string) error
	ListPatientInsurancePackages(ctx context.Context, opts *ListPatientInsurancePackagesOptions) (*ListPatientInsurancePackagesResult, error)
	UpdatePatientInsurancePackage(ctx context.Context, opts *UpdatePatientInsurancePackageOptions) error
	ReactivatePatientInsurancePackage(ctx context.Context, patientID, insuranceID string, expirationDate *time.Time) error
	UploadPatientInsuranceCardImage(ctx context.Context, patientID, insuranceID string, opts *UploadPatientInsuranceCardImageOptions) (*UploadPatientInsuranceCardImageResult, error)
	UploadPatientInsuranceCardImageReader(ctx context.Context, patientID, insuranceID string, opts *UploadPatientInsuranceCardImageReaderOptions) (*UploadPatientInsuranceCardImageResult, error)
	GetPatientInsuranceCardImage(ctx context.Context, patientID, insuranceID string) (*GetPatientInsuranceCardImageResult, error)

	// Patient Drivers License
	AddPatientDriversLicenseDocument(ctx context.Context, patientID string, opts *AddPatientDriversLicenseDocumentOptions) (*AddPatientDriversLicenseDocumentResult, error)
	AddPatientDriversLicenseDocumentReader(ctx context.Context, patientID string, opts *AddPatientDriversLicenseDocumentReaderOptions) (*AddPatientDriversLicenseDocumentResult, error)

	// Patient Lab Results
	AddLabResultDocumentReader(ctx context.Context, patientID string, departmentID string, opts *AddLabResultDocumentOptions) (int, error)
	ListLabResults(ctx context.Context, patientID string, departmentID string, opts *ListLabResultsOptions) (*ListLabResultsResult, error)
	ListChangedLabResults(ctx context.Context, opts *ListChangedLabResultsOptions) (*ListChangedLabResultsResult, error)

	// Health history
	ListSocialHistoryTemplates(context.Context) ([]*SocialHistoryTemplate, error)
	GetPatientSocialHistory(ctx context.Context, patientID string, opts *GetPatientSocialHistoryOptions) (*GetPatientSocialHistoryResponse, error)
	UpdatePatientSocialHistory(ctx context.Context, patientID string, opts *UpdatePatientSocialHistoryOptions) error
	GetHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string) (*HealthHistoryForm, error)
	UpdateHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string, form *HealthHistoryForm) error

	// Medications and Allergies
	SearchAllergies(ctx context.Context, searchVal string) ([]*Allergy, error)
	ListMedications(ctx context.Context, patientID string, opts *ListMedicationsOptions) (*ListMedicationsResult, error)
	SearchMedications(ctx context.Context, searchVal string) ([]*SearchMedicationsResult, error)

	// Appointment
	GetAppointment(ctx context.Context, appointmentID string) (*Appointment, error)
	ListBookedAppointments(context.Context, *ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error)
	ListChangedAppointments(context.Context, *ListChangedAppointmentsOptions) ([]*BookedAppointment, error)
	ListOpenAppointmentSlots(ctx context.Context, departmentID int, opts *ListOpenAppointmentSlotOptions) (*ListOpenAppointmentSlotsResult, error)
	BookAppointment(ctx context.Context, patientID, apptID string, opts *BookAppointmentOptions) (*BookedAppointment, error)
	UpdateBookedAppointment(ctx context.Context, apptID string, opts *UpdateBookedAppointmentOptions) error
	RescheduleAppointment(ctx context.Context, apptID int, opts *RescheduleAppointmentOptions) (*RescheduleAppointmentResult, error)
	ListAppointmentReminders(ctx context.Context, opts *ListAppointmentRemindersOptions) (*ListAppointmentRemindersResult, error)
	CreateAppointmentSlot(ctx context.Context, opts *CreateAppointmentSlotOptions) (*CreateAppointmentSlotResult, error)
	CreateAppointmentType(ctx context.Context, options *CreateAppointmentTypeOptions) (*CreateAppointmentTypeResult, error)
	ListAppointmentCustomFields(context.Context) ([]*AppointmentCustomField, error)
	FreezeAppointmentSlot(ctx context.Context, appointmentID string, opts *FreezeOrUnfreezeAppointmentSlotOptions) error
	UnfreezeAppointmentSlot(ctx context.Context, appointmentID string, opts *FreezeOrUnfreezeAppointmentSlotOptions) error

	// Appointment Check-In
	AppointmentCancelCheckIn(ctx context.Context, apptID string) error
	AppointmentCheckIn(ctx context.Context, apptID string) error
	AppointmentCheckOut(ctx context.Context, apptID string) error
	AppointmentStartCheckIn(ctx context.Context, apptID string) error

	// Appointment Note
	CreateAppointmentNote(ctx context.Context, appointmentID string, opts *CreateAppointmentNoteOptions) error
	DeleteAppointmentNote(ctx context.Context, appointmentID string, noteID string, opts *DeleteAppointmentNoteOptions) error
	ListAppointmentNotes(ctx context.Context, appointmentID string, opts *ListAppointmentNotesOptions) ([]*AppointmentNote, error)
	UpdateAppointmentNote(ctx context.Context, appointmentID string, noteID string, opts *UpdateAppointmentNoteOptions) error

	// Encounter
	GetPhysicalExam(ctx context.Context, encounterID string, opts *GetPhysicalExamOpts) (*PhysicalExam, error)
	ListEncounterDocuments(ctx context.Context, departmentID, patientID string, opts *ListEncounterDocumentsOptions) (*ListEncounterDocumentsResult, error)
	EncounterSummary(ctx context.Context, encounterID string, opts *EncounterSummaryOptions) (*EncounterSummaryResponse, error)

	// Provider
	ListProviders(context.Context, *ListProvidersOptions) (*ListProvidersResult, error)
	GetProvider(ctx context.Context, providerID string) (*Provider, error)

	// Subscription
	GetSubscription(ctx context.Context, feedType string) (*Subscription, error)
	ListSubscriptionEvents(ctx context.Context, feedType string) ([]*SubscriptionEvent, error)
	Subscribe(ctx context.Context, feedType string, opts *SubscribeOptions) error
	Unsubscribe(ctx context.Context, feedType string, opts *UnsubscribeOptions) error

	// List Changed
	ListChangedPatients(context.Context, *ListChangedPatientOptions) ([]*Patient, error)
	ListChangedProviders(context.Context, *ListChangedProviderOptions) ([]*Provider, error)
	ListChangedProblems(context.Context, *ListChangedProblemsOptions) ([]*ChangedProblem, error)
	ListChangedPrescriptions(ctx context.Context, options *ListChangedPrescriptionsOptions) (*ListChangedPrescriptionsResult, error)

	// Prescriptions
	UpdatePrescription(ctx context.Context, departmentID int, patientID int, documentID int, opts *UpdatePrescriptionOptions) (*UpdatePrescriptionResult, error)

	// Claims
	CreateFinancialClaim(ctx context.Context, opts *CreateClaimOptions) ([]string, error)
	ListClaims(ctx context.Context, opts *ListClaimsOptions) (*ListClaimsResult, error)

	// Telehealth
	GetTelehealthInviteURL(ctx context.Context, apptID string) (*GetTelehealthInviteURLResult, error)
}

type TokenProvider interface {
	Provide(context.Context) (string, time.Time, error)
}

type TokenCacher interface {
	Get(context.Context) (string, error)
	Set(context.Context, string, time.Time) error
}

type RateLimiter interface {
	Allowed(ctx context.Context, preview bool) (retryAfter time.Duration, err error)
}

type Stats interface {
	Request(method, path string) error
	ResponseSuccess() error
	ResponseError() error
}
