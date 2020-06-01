package athenahealth

import "time"

// Client describes a client for the athenahealth API.
type Client interface {
	GetPatient(id string) (*Patient, error)
	ListPatients(opts *ListPatientsOptions) ([]*Patient, error)
}

type TokenProvider interface {
	Provide() (string, time.Time, error)
}

type TokenCacher interface {
	Get() (string, error)
	Set(string, time.Time) error
}
