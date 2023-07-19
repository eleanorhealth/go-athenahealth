package enum

import "strconv"

type AppointmentStatus string

const (
	Cancelled     AppointmentStatus = "x"
	ChargeEntered AppointmentStatus = "4"
	CheckedIn     AppointmentStatus = "2"
	CheckedOut    AppointmentStatus = "3"
	Future        AppointmentStatus = "f"
	Open          AppointmentStatus = "o"
)

var appointmentStatusToString = map[AppointmentStatus]string{
	Cancelled:     "x",
	ChargeEntered: "4",
	CheckedIn:     "2",
	CheckedOut:    "3",
	Future:        "f",
	Open:          "o",
}

func (as AppointmentStatus) Int() int {
	intVal, convErr := strconv.Atoi(appointmentStatusToString[as])
	if convErr != nil {
		return -1
	}

	return intVal
}

func (as AppointmentStatus) String() string {
	return appointmentStatusToString[as]
}

func (as AppointmentStatus) Valid() bool {
	return appointmentStatusToString[as] != ""
}
