package athenahealth

import (
	"context"
	"time"
)

// Client describes a client for the athenahealth API.
type Client interface {
	GetDepartment(departmentID string) (*Department, error)
	ListDepartments(*ListDepartmentsOptions) (*ListDepartmentsResult, error)

	GetPatient(patientID string, opts *GetPatientOptions) (*Patient, error)
	ListPatients(*ListPatientsOptions) (*ListPatientsResult, error)

	ListSocialHistoryTemplates() ([]*SocialHistoryTemplate, error)
	GetPatientSocialHistory(patientID string, opts *GetPatientSocialHistoryOptions) (*GetPatientSocialHistoryResponse, error)
	UpdatePatientSocialHistory(patientID string, opts *UpdatePatientSocialHistoryOptions) error

	GetAppointment(appointmentID string) (*Appointment, error)
	ListBookedAppointments(*ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error)
	ListChangedAppointments(*ListChangedAppointmentsOptions) ([]*BookedAppointment, error)

	ListAppointmentCustomFields() ([]*AppointmentCustomField, error)

	CreateAppointmentNote(appointmentID string, opts *CreateAppointmentNoteOptions) error
	DeleteAppointmentNote(appointmentID string, noteID string, opts *DeleteAppointmentNoteOptions) error
	ListAppointmentNotes(appointmentID string, opts *ListAppointmentNotesOptions) ([]*AppointmentNote, error)
	UpdateAppointmentNote(appointmentID string, noteID string, opts *UpdateAppointmentNoteOptions) error

	ListProviders(*ListProvidersOptions) (*ListProvidersResult, error)
	GetProvider(providerID string) (*Provider, error)

	GetSubscription(feedType string) (*Subscription, error)
	ListSubscriptionEvents(feedType string) ([]*SubscriptionEvent, error)
	Subscribe(feedType string, opts *SubscribeOptions) error
	Unsubscribe(feedType string, opts *UnsubscribeOptions) error

	GetPatientPhoto(patientID string, opts *GetPatientPhotoOptions) (string, error)
	UpdatePatientPhoto(patientID string, data []byte) error

	ListChangedPatients(*ListChangedPatientOptions) ([]*Patient, error)
	ListChangedProviders(*ListChangedProviderOptions) ([]*Provider, error)
}

type TokenProvider interface {
	Provide() (string, time.Time, error)
}

type TokenCacher interface {
	Get() (string, error)
	Set(string, time.Time) error
}

type RateLimiter interface {
	Allowed(preview bool) (retryAfter time.Duration, err error)
}

type Stats interface {
	IncrRequests(context.Context) error
}
