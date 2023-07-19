package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppointmentStatus_Int(t *testing.T) {
	tests := []struct {
		status      AppointmentStatus
		expectedInt int
	}{
		{Cancelled, -1},
		{ChargeEntered, 4},
		{CheckedIn, 2},
		{CheckedOut, 3},
		{Future, -1},
		{Open, -1},
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
		{Cancelled, "x"},
		{ChargeEntered, "4"},
		{CheckedIn, "2"},
		{CheckedOut, "3"},
		{Future, "f"},
		{Open, "o"},
		{"does not exist", ""},
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
		{Cancelled, true},
		{ChargeEntered, true},
		{CheckedIn, true},
		{CheckedOut, true},
		{Future, true},
		{Open, true},
		{"does not exist", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.status.Valid())
		})
	}
}
