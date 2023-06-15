package athenahealth

import (
	"context"
	"time"
)

// Client describes a client for the athenahealth API.
type Client interface {
	GetDepartment(ctx context.Context, departmentID string) (*Department, error)
	ListDepartments(context.Context, *ListDepartmentsOptions) (*ListDepartmentsResult, error)

	GetPatient(ctx context.Context, patientID string, opts *GetPatientOptions) (*Patient, error)
	ListPatients(context.Context, *ListPatientsOptions) (*ListPatientsResult, error)
	UpdatePatient(ctx context.Context, patientID string, opts *UpdatePatientOptions) (*UpdatePatientResult, error)
	UpdatePatientInformationVerificationDetails(ctx context.Context, patientID string, opts *UpdatePatientInformationVerificationDetailsOptions) error

	ListSocialHistoryTemplates(context.Context) ([]*SocialHistoryTemplate, error)
	GetPatientSocialHistory(ctx context.Context, patientID string, opts *GetPatientSocialHistoryOptions) (*GetPatientSocialHistoryResponse, error)
	UpdatePatientSocialHistory(ctx context.Context, patientID string, opts *UpdatePatientSocialHistoryOptions) error

	GetAppointment(ctx context.Context, appointmentID string) (*Appointment, error)
	ListBookedAppointments(context.Context, *ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error)
	ListChangedAppointments(context.Context, *ListChangedAppointmentsOptions) ([]*BookedAppointment, error)
	ListOpenAppointmentSlots(ctx context.Context, departmentID int, opts *ListOpenAppointmentSlotOptions) (*ListOpenAppointmentSlotsResult, error)
	BookAppointment(ctx context.Context, patientID, apptID int, opts *BookAppointmentOptions) (*BookedAppointment, error)
	RescheduleAppointment(ctx context.Context, apptID int, opts *RescheduleAppointmentOptions) (*RescheduleAppointmentResult, error)

	CreateAppointmentSlot(ctx context.Context, opts *CreateAppointmentSlotOptions) (*CreateAppointmentSlotResult, error)
	CreateAppointmentType(ctx context.Context, options *CreateAppointmentTypeOptions) (*CreateAppointmentTypeResult, error)

	ListAppointmentCustomFields(context.Context) ([]*AppointmentCustomField, error)

	CreateAppointmentNote(ctx context.Context, appointmentID string, opts *CreateAppointmentNoteOptions) error
	DeleteAppointmentNote(ctx context.Context, appointmentID string, noteID string, opts *DeleteAppointmentNoteOptions) error
	ListAppointmentNotes(ctx context.Context, appointmentID string, opts *ListAppointmentNotesOptions) ([]*AppointmentNote, error)
	UpdateAppointmentNote(ctx context.Context, appointmentID string, noteID string, opts *UpdateAppointmentNoteOptions) error

	ListProviders(context.Context, *ListProvidersOptions) (*ListProvidersResult, error)
	GetProvider(ctx context.Context, providerID string) (*Provider, error)

	GetSubscription(ctx context.Context, feedType string) (*Subscription, error)
	ListSubscriptionEvents(ctx context.Context, feedType string) ([]*SubscriptionEvent, error)
	Subscribe(ctx context.Context, feedType string, opts *SubscribeOptions) error
	Unsubscribe(ctx context.Context, feedType string, opts *UnsubscribeOptions) error

	GetPatientPhoto(ctx context.Context, patientID string, opts *GetPatientPhotoOptions) (string, error)
	UpdatePatientPhoto(ctx context.Context, patientID string, data []byte) error

	ListChangedPatients(context.Context, *ListChangedPatientOptions) ([]*Patient, error)
	ListChangedProviders(context.Context, *ListChangedProviderOptions) ([]*Provider, error)
	ListChangedProblems(context.Context, *ListChangedProblemsOptions) ([]*ChangedProblem, error)

	ListProblems(ctx context.Context, patientID string, opts *ListProblemsOptions) ([]*Problem, error)
	ListAdminDocuments(ctx context.Context, patientID string, opts *ListAdminDocumentsOptions) (*ListAdminDocumentsResult, error)
	AddDocument(ctx context.Context, patientID string, opts *AddDocumentOptions) (string, error)
	AddPatientCaseDocument(ctx context.Context, patientID string, opts *AddPatientCaseDocumentOptions) (int, error)

	ListPatientsMatchingCustomField(ctx context.Context, opts *ListPatientsMatchingCustomFieldOptions) (*ListPatientsMatchingCustomFieldResult, error)

	GetPatientCustomFields(ctx context.Context, patientID, departmentID string) ([]*CustomFieldValue, error)
	UpdatePatientCustomFields(ctx context.Context, patientID, departmentID string, customFields []*CustomFieldValue) error

	CreateFinancialClaim(ctx context.Context, opts *CreateClaimOptions) ([]string, error)

	CreatePatient(ctx context.Context, opts *CreatePatientOptions) (string, error)

	ListClaims(ctx context.Context, opts *ListClaimsOptions) (*ListClaimsResult, error)

	CreatePatientInsurancePackage(ctx context.Context, opts *CreatePatientInsurancePackageOptions) (*InsurancePackage, error)
	DeletePatientInsurancePackage(ctx context.Context, patientID, insuranceID, cancellationNote string) error
	ListPatientInsurancePackages(ctx context.Context, opts *ListPatientInsurancePackagesOptions) (*ListPatientInsurancePackagesResult, error)
	UpdatePatientInsurancePackage(ctx context.Context, opts *UpdatePatientInsurancePackageOptions) error
	ReactivatePatientInsurancePackage(ctx context.Context, patientID, insuranceID string, expirationDate *time.Time) error

	UploadPatientInsuranceCardImage(ctx context.Context, patientID, insuranceID string, opts *UploadPatientInsuranceCardImageOptions) (*UploadPatientInsuranceCardImageResult, error)
	GetPatientInsuranceCardImage(ctx context.Context, patientID, insuranceID string) (*GetPatientInsuranceCardImageResult, error)

	AddPatientDriversLicenseDocument(ctx context.Context, patientID string, opts *AddPatientDriversLicenseDocumentOptions) (*AddPatientDriversLicenseDocumentResult, error)

	GetHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string) (*HealthHistoryForm, error)
	UpdateHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string, form *HealthHistoryForm) error

	SearchAllergies(ctx context.Context, searchVal string) ([]*Allergy, error)

	ListMedications(ctx context.Context, patientID string, opts *ListMedicationsOptions) (*ListMedicationsResult, error)
	SearchMedications(ctx context.Context, searchVal string) ([]*SearchMedicationsResult, error)
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
