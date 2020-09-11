package athenahealth

import "time"

// Client describes a client for the athenahealth API.
type Client interface {
	GetDepartment(string) (*Department, error)
	ListDepartments(*ListDepartmentsOptions) (*ListDepartmentsResult, error)

	GetPatient(string, *GetPatientOptions) (*Patient, error)
	ListPatients(*ListPatientsOptions) ([]*Patient, error)

	ListSocialHistoryTemplates() ([]*SocialHistoryTemplate, error)
	GetPatientSocialHistory(string, *GetPatientSocialHistoryOptions) (*GetPatientSocialHistoryResponse, error)
	UpdatePatientSocialHistory(string, *UpdatePatientSocialHistoryOptions) error

	GetAppointment(string) (*Appointment, error)
	ListBookedAppointments(*ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error)
	ListChangedAppointments(*ListChangedAppointmentsOptions) ([]*ChangedAppointment, error)

	ListAppointmentCustomFields() ([]*AppointmentCustomField, error)

	CreateAppointmentNote(string, *CreateAppointmentNoteOptions) error
	DeleteAppointmentNote(string, string, *DeleteAppointmentNoteOptions) error
	ListAppointmentNotes(string, *ListAppointmentNotesOptions) ([]*AppointmentNote, error)
	UpdateAppointmentNote(string, string, *UpdateAppointmentNoteOptions) error

	GetProvider(string) (*Provider, error)

	GetSubscription(string) (*Subscription, error)
	ListSubscriptionEvents(string) ([]*SubscriptionEvent, error)
	Subscribe(string, *SubscribeOptions) error
	Unsubscribe(string, *UnsubscribeOptions) error
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
