package athenahealth

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type ListAppointmentRemindersOptions struct {
	StartDate    time.Time
	EndDate      time.Time
	DepartmentID string

	AppointmentTypeID *int
	PatientID         *int
	ProviderID        *int
	ShowDeleted       *bool

	Pagination *PaginationOptions
}

type AppointmentReminder struct {
	AppointmentReminderID int    `json:"appointmentreminderid"`
	Note                  string `json:"note"`
	Status                string `json:"status"`
	Deleted               string `json:"deleted"`
	PatientID             int    `json:"patientid"`
	ProviderID            int    `json:"providerid"`
	EncounterID           int    `json:"encounterid"`
	DepartmentID          int    `json:"departmentid"`
	ApproximateDate       string `json:"approximatedate"`
	AppointmentTypeID     int    `json:"appointmenttypeid"`
	StaffInstructions     string `json:"staffinstructions"`
	PatientInstructions   string `json:"patientinstructions"`
}

type ListAppointmentRemindersResult struct {
	Reminders []AppointmentReminder `json:"reminders"`

	Pagination *PaginationResult
}

// ListAppointmentReminders - Retrieves a list of appointment reminders with the specified departmentid and
// approximatedate within the given date range
//
// GET /v1/{practiceid}/appointments/appointmentreminders
//
// https://docs.athenahealth.com/api/sandbox#/Appointments/getPracticeidAppointmentsAppointmentreminders
func (h *HTTPClient) ListAppointmentReminders(ctx context.Context, opts *ListAppointmentRemindersOptions) (*ListAppointmentRemindersResult, error) {
	out := &ListAppointmentRemindersResult{}

	q := url.Values{}
	if len(opts.DepartmentID) > 0 {
		q.Add("departmentid", opts.DepartmentID)
	}
	q.Add("startdate", opts.StartDate.Format("01/02/2006"))
	q.Add("enddate", opts.EndDate.Format("01/02/2006"))

	if opts.AppointmentTypeID != nil {
		q.Add("appointmenttypeid", strconv.Itoa(*opts.AppointmentTypeID))
	}
	if opts.PatientID != nil {
		q.Add("patientid", strconv.Itoa(*opts.PatientID))
	}
	if opts.ProviderID != nil {
		q.Add("providerid", strconv.Itoa(*opts.ProviderID))
	}
	if opts.ShowDeleted != nil {
		q.Add("showdeleted", strconv.FormatBool(*opts.ShowDeleted))
	}

	_, err := h.Get(ctx, "/appointments/appointmentreminders", q, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
