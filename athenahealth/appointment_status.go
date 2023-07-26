package athenahealth

// AppointmentStatus is derived from https://docs.athenahealth.com/api/api-ref/appointment#Book-appointment
// The athenaNet appointment status. There are several possible statuses.
// x=cancelled
// f=future (It can include appointments where were never checked in, even if the appointment date is in the past. It is up to a practice to cancel appointments as a no show when appropriate to do so.)
// o=open
// 2=checked in
// 3=checked out
// 4=charge entered (i.e. a past appointment)
type AppointmentStatus string

const (
	AppointmentStatusCancelled     AppointmentStatus = "x"
	AppointmentStatusChargeEntered AppointmentStatus = "4"
	AppointmentStatusCheckedIn     AppointmentStatus = "2"
	AppointmentStatusCheckedOut    AppointmentStatus = "3"
	AppointmentStatusFuture        AppointmentStatus = "f"
	AppointmentStatusOpen          AppointmentStatus = "o"
)

var appointmentStatuses = map[AppointmentStatus]struct{}{
	AppointmentStatusCancelled:     {},
	AppointmentStatusChargeEntered: {},
	AppointmentStatusCheckedIn:     {},
	AppointmentStatusCheckedOut:    {},
	AppointmentStatusFuture:        {},
	AppointmentStatusOpen:          {},
}

func (as AppointmentStatus) String() string {
	return string(as)
}

func (as AppointmentStatus) Valid() bool {
	_, exists := appointmentStatuses[as]
	return exists
}
