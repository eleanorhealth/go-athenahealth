package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Appointment struct {
	AppointmentID              string `json:"appointmentid"`
	AppointmentStatus          string `json:"appointmentstatus"`
	AppointmentType            string `json:"appointmenttype"`
	AppointmentTypeID          string `json:"appointmenttypeid"`
	ChargeEntryNotRequired     bool   `json:"chargeentrynotrequired"`
	Date                       string `json:"date"`
	DepartmentID               string `json:"departmentid"`
	Duration                   int    `json:"duration"`
	EncounterID                string `json:"encounterid"`
	PatientAppointmentTypeName string `json:"patientappointmenttypename"`
	ProviderID                 string `json:"providerid"`
	StartTime                  string `json:"starttime"`
}

// GetAppointment - Single appointment.
// GET /v1/{practiceid}/appointments/{appointmentid}
// https://developer.athenahealth.com/docs/read/appointments/Appointments#section-1
func (h *HTTPClient) GetAppointment(ctx context.Context, id string) (*Appointment, error) {
	out := []*Appointment{}

	_, err := h.Get(ctx, fmt.Sprintf("/appointments/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type AppointmentCustomField struct {
	CaseSensitive  bool   `json:"casesensitive"`
	CustomFieldID  int    `json:"customfieldid"`
	DisallowUpdate bool   `json:"disallowupdate"`
	Name           string `json:"name"`
	Searchable     bool   `json:"searchable,omitempty"`
	Select         bool   `json:"select"`
	Type           string `json:"type"`
	SelectList     []struct {
		OptionValue string `json:"optionvalue"`
		OptionID    int    `json:"optionid"`
	} `json:"selectlist,omitempty"`
}

type listAppointmentCustomFieldsResponse struct {
	AppointmentCustomFields []*AppointmentCustomField `json:"appointmentcustomfields"`
}

// ListAppointmentCustomFields - List of appointment custom fields (practice specific).
// GET /v1/{practiceid}/appointments/customfields
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Custom_Fields#section-0
func (h *HTTPClient) ListAppointmentCustomFields(ctx context.Context) ([]*AppointmentCustomField, error) {
	out := &listAppointmentCustomFieldsResponse{}

	_, err := h.Get(ctx, "/appointments/customfields", nil, &out)
	if err != nil {
		return nil, err
	}

	return out.AppointmentCustomFields, nil
}

type BookedAppointment struct {
	AppointmentID    string `json:"appointmentid"`
	AppointmentCopay struct {
		CollectedForOther       int `json:"collectedforother"`
		CollectedForAppointment int `json:"collectedforappointment"`
		InsuranceCopay          int `json:"insurancecopay"`
	} `json:"appointmentcopay"`
	AppointmentNotes []struct {
		DisplayOnSchedule bool   `json:"displayonschedule"`
		Text              string `json:"text"`
		ID                int    `json:"id"`
	} `json:"appointmentnotes"`
	AppointmentStatus          string `json:"appointmentstatus"`
	AppointmentType            string `json:"appointmenttype"`
	AppointmentTypeID          string `json:"appointmenttypeid"`
	CancelledBy                string `json:"cancelledby"`
	CancelledDatetime          string `json:"cancelleddatetime"`
	CancelReasonID             string `json:"cancelreasonid"`
	CancelReasonName           string `json:"cancelreasonname"`
	CancelReasonNoShow         bool   `json:"cancelreasonnoshow"`
	CancelReasonSlotAvailable  bool   `json:"cancelreasonslotavailable"`
	ChargeEntryNotRequired     bool   `json:"chargeentrynotrequired"`
	CoordinatorEnterprise      bool   `json:"coordinatorenterprise"`
	Copay                      int    `json:"copay"`
	Date                       string `json:"date"`
	DepartmentID               string `json:"departmentid"`
	Duration                   int    `json:"duration"`
	EncounterID                string `json:"encounterid"`
	HL7ProviderID              int    `json:"hl7providerid"`
	LastModified               string `json:"lastmodified"`
	LastModifiedBy             string `json:"lastmodifiedby"`
	PatientAppointmentTypeName string `json:"patientappointmenttypename"`
	PatientID                  string `json:"patientid"`
	ProviderID                 string `json:"providerid"`
	ScheduledBy                string `json:"scheduledby"`
	ScheduledDatetime          string `json:"scheduleddatetime"`
	StartTime                  string `json:"starttime"`
	TemplateAppointmentID      string `json:"templateappointmentid"`
	TemplateAppointmentTypeID  string `json:"templateappointmenttypeid"`
}

type ListBookedAppointmentsOptions struct {
	DepartmentID      string
	EndDate           time.Time
	PatientID         string
	ProviderID        string
	StartDate         time.Time
	AppointmentStatus string

	Pagination *PaginationOptions
}

type ListBookedAppointmentsResult struct {
	BookedAppointments []*BookedAppointment
	Pagination         *PaginationResult
}

type listBookedAppointmentsResponse struct {
	Appointments []*BookedAppointment `json:"appointments"`

	PaginationResponse
}

// ListBookedAppointments - Booked appointment slots.
// GET /v1/{practiceid}/appointments/booked
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-3
func (h *HTTPClient) ListBookedAppointments(ctx context.Context, opts *ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error) {
	out := &listBookedAppointmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.ProviderID) > 0 {
			q.Add("providerid", opts.ProviderID)
		}

		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if !opts.StartDate.IsZero() {
			q.Add("startdate", opts.StartDate.Format("01/02/2006"))
		}

		if !opts.EndDate.IsZero() {
			q.Add("enddate", opts.EndDate.Format("01/02/2006"))
		}

		if len(opts.AppointmentStatus) > 0 {
			q.Add("appointmentstatus", opts.AppointmentStatus)
		}

		if opts.Pagination != nil {
			if opts.Pagination.Limit > 0 {
				q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
			}

			if opts.Pagination.Offset > 0 {
				q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
			}
		}
	}

	_, err := h.Get(ctx, "/appointments/booked", q, out)
	if err != nil {
		return nil, err
	}

	return &ListBookedAppointmentsResult{
		BookedAppointments: out.Appointments,
		Pagination:         makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type ListChangedAppointmentsOptions struct {
	DepartmentID               string
	LeaveUnprocessed           bool
	PatientID                  string
	ProviderID                 string
	ShowPatientDetail          bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedAppointmentsResponse struct {
	ChangedAppointments []*BookedAppointment `json:"appointments"`
}

// ListChangedAppointments - Changed appointment slots.
// GET /v1/{practiceid}/appointments/changed
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-5
func (h *HTTPClient) ListChangedAppointments(ctx context.Context, opts *ListChangedAppointmentsOptions) ([]*BookedAppointment, error) {
	out := &listChangedAppointmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.ProviderID) > 0 {
			q.Add("providerid", opts.ProviderID)
		}

		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if opts.ShowPatientDetail {
			q.Add("showpatientdetail", strconv.FormatBool(opts.ShowPatientDetail))
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
		}

		if opts.LeaveUnprocessed {
			q.Add("leaveunprocessed", strconv.FormatBool(opts.LeaveUnprocessed))
		}
	}

	_, err := h.Get(ctx, "/appointments/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedAppointments, nil
}

type CreateAppointmentNoteOptions struct {
	AppointmentID     string
	DisplayOnSchedule bool
	NoteText          string
}

// CreateAppointmentNote - Notes for this appointment.
// POST /v1/{practiceid}/appointments/{appointmentid}/notes
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-0
func (h *HTTPClient) CreateAppointmentNote(ctx context.Context, appointmentID string, opts *CreateAppointmentNoteOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if opts.DisplayOnSchedule {
			form.Add("displayonschedule", strconv.FormatBool(opts.DisplayOnSchedule))
		}

		if len(opts.NoteText) > 0 {
			form.Add("notetext", opts.NoteText)
		}
	}

	_, err := h.PostForm(ctx, fmt.Sprintf("/appointments/%s/notes", appointmentID), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type AppointmentNote struct {
	Created           string `json:"created"`
	CreatedBy         string `json:"createdby"`
	DisplayOnSchedule bool   `json:"displayonschedule"`
	NoteID            string `json:"noteid"`
	NoteText          string `json:"notetext"`
}

type ListAppointmentNotesOptions struct {
	AppointmentID string
	ShowDeleted   bool
}

type listAppointmentNotesResponse struct {
	Notes []*AppointmentNote `json:"notes"`
}

// ListAppointmentNotes - Notes for this appointment.
// GET /v1/{practiceid}/appointments/{appointmentid}/notes
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-1
func (h *HTTPClient) ListAppointmentNotes(ctx context.Context, appointmentID string, opts *ListAppointmentNotesOptions) ([]*AppointmentNote, error) {
	out := &listAppointmentNotesResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.AppointmentID) > 0 {
			q.Add("appointmentid", opts.AppointmentID)
		}

		if opts.ShowDeleted {
			q.Add("showdeleted", strconv.FormatBool(opts.ShowDeleted))
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("/appointments/%s/notes", appointmentID), q, out)
	if err != nil {
		return nil, err
	}

	return out.Notes, nil
}

type UpdateAppointmentNoteOptions struct {
	AppointmentID     string
	DisplayOnSchedule bool
	NoteID            string
	NoteText          string
}

// UpdateAppointmentNote - Notes for this appointment.
// PUT /v1/{practiceid}/appointments/{appointmentid}/notes/{noteid}
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-3
func (h *HTTPClient) UpdateAppointmentNote(ctx context.Context, appointmentID, noteID string, opts *UpdateAppointmentNoteOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if opts.DisplayOnSchedule {
			form.Add("displayonschedule", strconv.FormatBool(opts.DisplayOnSchedule))
		}

		if len(opts.NoteID) > 0 {
			form.Add("noteid", opts.NoteID)
		}

		if len(opts.NoteText) > 0 {
			form.Add("notetext", opts.NoteText)
		}
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/appointments/%s/notes/%s", appointmentID, noteID), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type DeleteAppointmentNoteOptions struct {
	AppointmentID string
	NoteID        string
}

// DeleteAppointmentNote - Notes for this appointment.
// DELETE /v1/{practiceid}/appointments/{appointmentid}/notes/{noteid}
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-0
func (h *HTTPClient) DeleteAppointmentNote(ctx context.Context, appointmentID, noteID string, opts *DeleteAppointmentNoteOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if len(opts.NoteID) > 0 {
			form.Add("noteid", opts.NoteID)
		}
	}

	_, err := h.DeleteForm(ctx, fmt.Sprintf("/appointments/%s/notes/%s", appointmentID, noteID), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type ListOpenAppointmentSlotOptions struct {
	// Normally, an appointment reason ID should be used which will map to the correct underlying appointment type in athenaNet. This field will ignore the practice's existing setup for what should be scheduled. Please consult with athenahealth before using. Either an appointmenttypeid or a reasonid must be specified or no results will be returned.
	AppointmentTypeID int

	// The athenaNet patient appointment reason ID, from GET /patientappointmentreasons. While this is not technically required due to some unusual use cases, it is highly recommended for most calls. We do allow a special value of -1 for the reasonid. This reasonid will return open, web-schedulable slots regardless of reason. However, slots returned using a search of -1 may return slots that are not bookable by any reason ID (they may be bookable by specific appointment type IDs instead). This argument allows multiple valid reason IDs to be specified (e.g. reasonid=1,2,3), so if you are looking for slots that match "any" reason, it is recommended that you enumerate the set of reasons you are looking for. Either a reasonid or an appointmenttypeid must be specified or no results will be returned. If a reasonid other than -1 is specified then a providerid must also be specified.
	ReasonIDs []int

	// Bypass checks that usually require returned appointments to be some amount of hours in the future (as configured by the practice, defaulting to 24 hours), and also ignores the setting that only shows appointments for a certain number of days in the future (also configurable by the practice, defaulting to 90 days).
	BypassScheduleTimeChecks bool

	// End of the appointment search date range (mm/dd/yyyy). Inclusive. Defaults to seven days from startdate.
	EndDate time.Time

	// 	The athenaNet provider ID. Required if a reasonid other than -1 is specified.
	ProviderIDs []int

	// Start of the appointment search date range (mm/dd/yyyy). Inclusive. Defaults to today.
	StartDate time.Time

	// By default, we show only appointments that are are available to scheduled via the API. This flag allows you to bypass that restriction for viewing available appointments (but you still may not be able to schedule based on this permission!). This flag does not, however, show the full schedule (that is, appointments that are already booked).
	IgnoreSchedulablePermission bool

	// By default, we hide appointments that are frozen from being returned via the API. This flag allows you to show frozen slots in the set of results returned.
	ShowFrozenSlots bool

	// Number of entries to return (default 1000, max 10000)Please note that this endpoint has a different default and max than normal.
	Limit int

	// Starting point of entries; 0-indexed
	Offset int
}

type OpenAppointmentSlot struct {
	AppointmentID              int    `json:"appointmentid"`
	AppointmentType            string `json:"appointmenttype"`
	AppointmentTypeID          int    `json:"appointmenttypeid"`
	Date                       string `json:"date"`
	DepartmentID               int    `json:"departmentid"`
	Duration                   int    `json:"duration"`
	Frozen                     bool   `json:"frozen"`
	LocalProviderID            int    `json:"localproviderid"`
	PatientAppointmentTypeName string `json:"patientappointmenttypename"`
	ProviderID                 int    `json:"providerid"`
	StartTime                  string `json:"starttime"`
}

type listOpenAppointmentSlotsResponse struct {
	Appointments []*OpenAppointmentSlot `json:"appointments"`

	PaginationResponse
}

type ListOpenAppointmentSlotsResult struct {
	Appointments []*OpenAppointmentSlot

	Pagination *PaginationResult
}

// ListOpenAppointmentSlots - Get list of open appointment slots.
// GET /v1/{practiceid}/appointments/open
// https://docs.athenahealth.com/api/api-ref/appointment-slot#Get-list-of-open-appointment-slots
func (h *HTTPClient) ListOpenAppointmentSlots(ctx context.Context, departmentID int, opts *ListOpenAppointmentSlotOptions) (*ListOpenAppointmentSlotsResult, error) {
	out := &listOpenAppointmentSlotsResponse{}

	q := url.Values{}

	q.Add("departmentid", strconv.Itoa(departmentID))

	if opts != nil {
		if opts.AppointmentTypeID > 0 {
			q.Add("appointmenttypeid", strconv.Itoa(opts.AppointmentTypeID))
		}

		if len(opts.ReasonIDs) > 0 {
			var reasonIDs []string
			for _, reasonID := range opts.ReasonIDs {
				reasonIDs = append(reasonIDs, strconv.Itoa(reasonID))
			}

			q.Add("reasonid", strings.Join(reasonIDs, ","))
		}

		if opts.BypassScheduleTimeChecks {
			q.Add("bypassscheduletimechecks", "true")
		}

		if !opts.EndDate.IsZero() {
			q.Add("enddate", opts.EndDate.Format("01/02/2006"))
		}

		if len(opts.ProviderIDs) > 0 {
			var providerIDs []string
			for _, providerID := range opts.ProviderIDs {
				providerIDs = append(providerIDs, strconv.Itoa(providerID))
			}

			q.Add("providerid", strings.Join(providerIDs, ","))
		}

		if !opts.StartDate.IsZero() {
			q.Add("startdate", opts.StartDate.Format("01/02/2006"))
		}

		if opts.IgnoreSchedulablePermission {
			q.Add("ignoreschedulablepermission", "true")
		}

		if opts.ShowFrozenSlots {
			q.Add("showfrozenslots", "true")
		}

		if opts.Limit > 0 {
			q.Add("limit", strconv.Itoa(opts.Limit))
		}

		if opts.Offset > 0 {
			q.Add("offset", strconv.Itoa(opts.Offset))
		}
	}

	_, err := h.Get(ctx, "/appointments/open", q, out)
	if err != nil {
		return nil, err
	}

	return &ListOpenAppointmentSlotsResult{
		Appointments: out.Appointments,
		Pagination:   makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type BookAppointmentOptions struct {
	AppointmentTypeID           int
	BookingNote                 string
	DepartmentID                int
	DoNotSendConfirmationEmail  bool
	IgnoreSchedulablePermission bool
	NoPatientCase               bool
	ReasonID                    int
	Urgent                      bool
}

// BookAppointment - Create a single appointment for specific patient
// PUT /v1/{practiceid}/appointments/{appointmentid}
// https://docs.athenahealth.com/api/api-ref/appointment#Book-appointment
func (h *HTTPClient) BookAppointment(ctx context.Context, patientID, apptID int, opts *BookAppointmentOptions) (*BookedAppointment, error) {
	var out []*BookedAppointment

	form := url.Values{}

	form.Add("patientid", strconv.Itoa(patientID))

	if opts != nil {
		if opts.AppointmentTypeID > 0 {
			form.Add("appointmenttypeid", strconv.Itoa(opts.AppointmentTypeID))
		}

		if len(opts.BookingNote) > 0 {
			form.Add("bookingnote", opts.BookingNote)
		}

		if opts.DepartmentID > 0 {
			form.Add("departmentid", strconv.Itoa(opts.DepartmentID))
		}

		if opts.DoNotSendConfirmationEmail {
			form.Add("donotsendconfirmationemail", "true")
		}

		if opts.IgnoreSchedulablePermission {
			form.Add("ignoreschedulablepermission", "true")
		}

		if opts.NoPatientCase {
			form.Add("nopatientcase", "true")
		}

		if opts.ReasonID > 0 {
			form.Add("reasonid", strconv.Itoa(opts.ReasonID))
		}

		if opts.Urgent {
			form.Add("urgent", "true")
		}
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/appointments/%d", apptID), form, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("unexpected length returned")
	}

	return out[0], nil
}
