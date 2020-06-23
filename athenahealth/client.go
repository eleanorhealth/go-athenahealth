package athenahealth

import "time"

// Client describes a client for the athenahealth API.
type Client interface {
	GetDepartment(string) (*Department, error)

	GetPatient(string) (*Patient, error)
	ListPatients(*ListPatientsOptions) ([]*Patient, error)

	GetAppointment(string) (*Appointment, error)
	ListBookedAppointments(*ListBookedAppointmentsOptions) ([]*BookedAppointment, error)
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
