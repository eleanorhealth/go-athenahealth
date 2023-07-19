package athenahealth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppointmentStatus_Int(t *testing.T) {
	tests := []struct {
		status      AppointmentStatus
		expectedInt int
	}{
		{AppointmentStatusCancelled, -1},
		{AppointmentStatusChargeEntered, 4},
		{AppointmentStatusCheckedIn, 2},
		{AppointmentStatusCheckedOut, 3},
		{AppointmentStatusFuture, -1},
		{AppointmentStatusOpen, -1},
		{"unknown", -1},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expectedInt, tt.status.Int())
		})
	}
}

func TestAppointmentStatus_String(t *testing.T) {
	tests := []struct {
		status      AppointmentStatus
		expectedStr string
	}{
		{AppointmentStatusCancelled, "x"},
		{AppointmentStatusChargeEntered, "4"},
		{AppointmentStatusCheckedIn, "2"},
		{AppointmentStatusCheckedOut, "3"},
		{AppointmentStatusFuture, "f"},
		{AppointmentStatusOpen, "o"},
		{"does not exist", "does not exist"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expectedStr, tt.status.String())
		})
	}
}

func TestAppointmentStatus_Valid(t *testing.T) {
	tests := []struct {
		status AppointmentStatus
		valid  bool
	}{
		{AppointmentStatusCancelled, true},
		{AppointmentStatusChargeEntered, true},
		{AppointmentStatusCheckedIn, true},
		{AppointmentStatusCheckedOut, true},
		{AppointmentStatusFuture, true},
		{AppointmentStatusOpen, true},
		{"does not exist", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.status.Valid())
		})
	}
}
