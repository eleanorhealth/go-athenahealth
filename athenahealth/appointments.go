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
	AppointmentID              string            `json:"appointmentid"`
	AppointmentStatus          AppointmentStatus `json:"appointmentstatus"`
	AppointmentType            string            `json:"appointmenttype"`
	AppointmentTypeID          string            `json:"appointmenttypeid"`
	ChargeEntryNotRequired     bool              `json:"chargeentrynotrequired"`
	Date                       string            `json:"date"`
	DepartmentID               string            `json:"departmentid"`
	Duration                   int               `json:"duration"`
	EncounterID                string            `json:"encounterid"`
	PatientAppointmentTypeName string            `json:"patientappointmenttypename"`
	PatientID                  string            `json:"patientid"`
	ProviderID                 string            `json:"providerid"`
	RenderingProviderID        string            `json:"renderingproviderid"`
	StartTime                  string            `json:"starttime"`
}

// GetAppointment - Get single appointment
//
// GET /v1/{practiceid}/appointments/{appointmentid}
//
// https://docs.athenahealth.com/api/api-ref/appointment#Get-appointment-details
func (h *HTTPClient) GetAppointment(ctx context.Context, id string) (*Appointment, error) {
	var out []*Appointment

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

// ListAppointmentCustomFields - List of appointment custom fields (practice specific)
//
// GET /v1/{practiceid}/appointments/customfields
//
// https://docs.athenahealth.com/api/api-ref/appointment-custom-fields#Get-the-list-of-appointment-custom-fields
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
	AppointmentStatus          AppointmentStatus `json:"appointmentstatus"`
	AppointmentType            string            `json:"appointmenttype"`
	AppointmentTypeID          string            `json:"appointmenttypeid"`
	CancelReasonID             string            `json:"cancelreasonid"`
	CancelReasonName           string            `json:"cancelreasonname"`
	CancelReasonNoShow         bool              `json:"cancelreasonnoshow"`
	CancelReasonSlotAvailable  bool              `json:"cancelreasonslotavailable"`
	CancelledBy                string            `json:"cancelledby"`
	CancelledDatetime          string            `json:"cancelleddatetime"`
	ChargeEntryNotRequired     bool              `json:"chargeentrynotrequired"`
	CoordinatorEnterprise      bool              `json:"coordinatorenterprise"`
	Copay                      int               `json:"copay"`
	Date                       string            `json:"date"`
	DepartmentID               string            `json:"departmentid"`
	Duration                   int               `json:"duration"`
	EncounterID                string            `json:"encounterid"`
	HL7ProviderID              int               `json:"hl7providerid"`
	LastModified               string            `json:"lastmodified"`
	LastModifiedBy             string            `json:"lastmodifiedby"`
	PatientAppointmentTypeName string            `json:"patientappointmenttypename"`
	PatientID                  string            `json:"patientid"`
	ProviderID                 string            `json:"providerid"`
	RenderingProviderID        string            `json:"renderingproviderid"`
	ScheduledBy                string            `json:"scheduledby"`
	ScheduledDatetime          string            `json:"scheduleddatetime"`
	StartTime                  string            `json:"starttime"`
	TemplateAppointmentID      string            `json:"templateappointmentid"`
	TemplateAppointmentTypeID  string            `json:"templateappointmenttypeid"`
}

type ListBookedAppointmentsOptions struct {
	// Filter appointments by status.
	AppointmentStatus *AppointmentStatus
	// Filter by appointment type ID.
	AppointmentTypeID *string
	// The athenaNet department ID.
	DepartmentID *string
	// End of the appointment search date range (mm/dd/yyyy). Inclusive.
	EndDate time.Time
	// The athenaNet patient ID. If operating in a Provider Group Enterprise practice, this should be the enterprise patient ID.
	PatientID *string
	// The athenaNet provider ID. Multiple IDs (either as a comma delimited list or multiple POSTed values) are allowed.
	ProviderID *string
	// Start of the appointment search date range (mm/dd/yyyy). Inclusive.
	StartDate time.Time

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

// ListBookedAppointments - Booked appointment slots
//
// GET /v1/{practiceid}/appointments/booked
//
// https://docs.athenahealth.com/api/api-ref/appointment#Get-list-of-booked-appointments
func (h *HTTPClient) ListBookedAppointments(ctx context.Context, opts *ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error) {
	out := &listBookedAppointmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.AppointmentStatus != nil {
			if opts.AppointmentStatus.Valid() {
				q.Add("appointmentstatus", opts.AppointmentStatus.String())
			} else {
				return nil, fmt.Errorf("invalid AppointmentStatus [%s]", opts.AppointmentStatus.String())
			}
		}

		if opts.AppointmentTypeID != nil && *opts.AppointmentTypeID != "" {
			q.Add("appointmenttypeid", *opts.AppointmentTypeID)
		}

		if opts.DepartmentID != nil && *opts.DepartmentID != "" {
			q.Add("departmentid", *opts.DepartmentID)
		}

		if opts.EndDate.IsZero() {
			return nil, fmt.Errorf("cannot ListBookedAppointments when StartDate is zero [%+v]", opts.EndDate)
		} else {
			q.Add("enddate", opts.EndDate.Format("01/02/2006"))
		}

		if opts.PatientID != nil && *opts.PatientID != "" {
			q.Add("patientid", *opts.PatientID)
		}

		if opts.ProviderID != nil && *opts.ProviderID != "" {
			q.Add("providerid", *opts.ProviderID)
		}

		if opts.StartDate.IsZero() {
			return nil, fmt.Errorf("cannot ListBookedAppointments when StartDate is zero [%+v]", opts.StartDate)
		} else {
			q.Add("startdate", opts.StartDate.Format("01/02/2006"))
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
	LeaveUnprocessed           *bool
	PatientID                  string
	ProviderID                 string
	ShowPatientDetail          *bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedAppointmentsResponse struct {
	ChangedAppointments []*BookedAppointment `json:"appointments"`
}

// ListChangedAppointments - List of changes in appointment slots
//
// GET /v1/{practiceid}/appointments/changed
//
// https://docs.athenahealth.com/api/api-ref/appointment#Get-list-of-changes-in-appointment-slots-based-on-subscribed-events
func (h *HTTPClient) ListChangedAppointments(ctx context.Context, opts *ListChangedAppointmentsOptions) ([]*BookedAppointment, error) {
	out := &listChangedAppointmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if opts.LeaveUnprocessed != nil {
			q.Add("leaveunprocessed", strconv.FormatBool(*opts.LeaveUnprocessed))
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if len(opts.ProviderID) > 0 {
			q.Add("providerid", opts.ProviderID)
		}

		if opts.ShowPatientDetail != nil {
			q.Add("showpatientdetail", strconv.FormatBool(*opts.ShowPatientDetail))
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
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
	DisplayOnSchedule *bool
	NoteText          string
}

// CreateAppointmentNote - Create note for specific appointment
//
// POST /v1/{practiceid}/appointments/{appointmentid}/notes
//
// https://docs.athenahealth.com/api/api-ref/appointment-notes#Create-appointment-note
func (h *HTTPClient) CreateAppointmentNote(ctx context.Context, appointmentID string, opts *CreateAppointmentNoteOptions) error {
	form := url.Values{}

	if opts != nil {
		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if opts.DisplayOnSchedule != nil {
			form.Add("displayonschedule", strconv.FormatBool(*opts.DisplayOnSchedule))
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
	ShowDeleted   *bool
}

type listAppointmentNotesResponse struct {
	Notes []*AppointmentNote `json:"notes"`
}

// ListAppointmentNotes - List all notes for specific appointment
//
// GET /v1/{practiceid}/appointments/{appointmentid}/notes
//
// https://docs.athenahealth.com/api/api-ref/appointment-notes#Get-all-appointment-notes
func (h *HTTPClient) ListAppointmentNotes(ctx context.Context, appointmentID string, opts *ListAppointmentNotesOptions) ([]*AppointmentNote, error) {
	out := &listAppointmentNotesResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.AppointmentID) > 0 {
			q.Add("appointmentid", opts.AppointmentID)
		}

		if opts.ShowDeleted != nil {
			q.Add("showdeleted", strconv.FormatBool(*opts.ShowDeleted))
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
	DisplayOnSchedule *bool
	NoteID            string
	NoteText          string
}

// UpdateAppointmentNote - Update note for specific appointment
//
// PUT /v1/{practiceid}/appointments/{appointmentid}/notes/{noteid}
//
// https://docs.athenahealth.com/api/api-ref/appointment-notes#Update-appointment-note
func (h *HTTPClient) UpdateAppointmentNote(ctx context.Context, appointmentID, noteID string, opts *UpdateAppointmentNoteOptions) error {
	form := url.Values{}

	if opts != nil {
		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if opts.DisplayOnSchedule != nil {
			form.Add("displayonschedule", strconv.FormatBool(*opts.DisplayOnSchedule))
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

// DeleteAppointmentNote - Delete note for specific appointment
//
// DELETE /v1/{practiceid}/appointments/{appointmentid}/notes/{noteid}
//
// https://docs.athenahealth.com/api/api-ref/appointment-notes#Delete-appointment-note
func (h *HTTPClient) DeleteAppointmentNote(ctx context.Context, appointmentID, noteID string, opts *DeleteAppointmentNoteOptions) error {
	form := url.Values{}

	if opts != nil {
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
	AppointmentTypeID *string
	// Bypass checks that usually require returned appointments to be some amount of hours in the future (as configured by the practice, defaulting to 24 hours), and also ignores the setting that only shows appointments for a certain number of days in the future (also configurable by the practice, defaulting to 90 days).
	BypassScheduleTimeChecks *bool
	// The athenaNet department ID.
	DepartmentIDs []string
	// End of the appointment search date range (mm/dd/yyyy). Inclusive. Defaults to seven days from startdate.
	EndDate *time.Time
	// By default, we show only appointments that are available to scheduled via the API. This flag allows you to bypass that restriction for viewing available appointments (but you still may not be able to schedule based on this permission!). This flag does not, however, show the full schedule (that is, appointments that are already booked).
	IgnoreSchedulablePermission *bool
	// 	The athenaNet provider ID. Required if a reasonid other than -1 is specified.
	ProviderIDs []string
	// The athenaNet patient appointment reason ID, from GET /patientappointmentreasons. While this is not technically required due to some unusual use cases, it is highly recommended for most calls. We do allow a special value of -1 for the reasonid. This reasonid will return open, web-schedulable slots regardless of reason. However, slots returned using a search of -1 may return slots that are not bookable by any reason ID (they may be bookable by specific appointment type IDs instead). This argument allows multiple valid reason IDs to be specified (e.g. reasonid=1,2,3), so if you are looking for slots that match "any" reason, it is recommended that you enumerate the set of reasons you are looking for. Either a reasonid or an appointmenttypeid must be specified or no results will be returned. If a reasonid other than -1 is specified then a providerid must also be specified.
	ReasonIDs []string
	// Start of the appointment search date range (mm/dd/yyyy). Inclusive. Defaults to today.
	StartDate *time.Time
	// By default, we hide appointments that are frozen from being returned via the API. This flag allows you to show frozen slots in the set of results returned.
	ShowFrozenSlots *bool

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

// ListOpenAppointmentSlots - Get list of open appointment slots
//
// GET /v1/{practiceid}/appointments/open
//
// https://docs.athenahealth.com/api/api-ref/appointment-slot#Get-list-of-open-appointment-slots
func (h *HTTPClient) ListOpenAppointmentSlots(ctx context.Context, opts *ListOpenAppointmentSlotOptions) (*ListOpenAppointmentSlotsResult, error) {
	out := &listOpenAppointmentSlotsResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.AppointmentTypeID != nil && *opts.AppointmentTypeID != "" {
			q.Add("appointmenttypeid", *opts.AppointmentTypeID)
		}

		if opts.BypassScheduleTimeChecks != nil {
			q.Add("bypassscheduletimechecks", strconv.FormatBool(*opts.BypassScheduleTimeChecks))
		}

		if len(opts.DepartmentIDs) > 0 {
			q.Add("departmentid", strings.Join(opts.DepartmentIDs, ","))
		}

		if opts.EndDate != nil && !(*opts.EndDate).IsZero() {
			q.Add("enddate", opts.EndDate.Format("01/02/2006"))
		}

		if opts.IgnoreSchedulablePermission != nil {
			q.Add("ignoreschedulablepermission", strconv.FormatBool(*opts.IgnoreSchedulablePermission))
		}

		if len(opts.ProviderIDs) > 0 {
			q.Add("providerid", strings.Join(opts.ProviderIDs, ","))
		}

		if len(opts.ReasonIDs) > 0 {
			q.Add("reasonid", strings.Join(opts.ReasonIDs, ","))
		}

		if opts.StartDate != nil && !(*opts.StartDate).IsZero() {
			q.Add("startdate", opts.StartDate.Format("01/02/2006"))
		}

		if opts.ShowFrozenSlots != nil {
			q.Add("showfrozenslots", strconv.FormatBool(*opts.ShowFrozenSlots))
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
	// The appointment type to be booked. This field should never be used for booking appointments for web-based scheduling. The use of this field is reserved for digital check-in (aka "kiosk") or an application used by practice staff. One of this or reasonid is required.
	AppointmentTypeID *string
	// A note from the patient about why this appointment is being booked
	BookingNote *string
	// The athenaNet department ID.
	DepartmentID *string
	// For clients with athenaCommunicator, certain appointment types can be configured to have an appointment confirmation email sent to the patient at time of appointment booking. If this parameter is set to true, that email will not be sent. This should only be used if you plan on sending a confirmation email via another method.
	DoNotSendConfirmationEmail *bool
	// By default, we allow booking of appointments marked as schedulable via the web. This flag allows you to bypass that restriction for booking.
	IgnoreSchedulablePermission *bool
	// By default, we create a patient case upon booking an appointment for new patients. Setting this to true bypasses that patient case.
	NoPatientCase *bool
	// The athenaNet patient ID.
	PatientID string
	// The appointment reason ID to be booked. This field is required for booking appointments for web-based scheduling and is a reason that is retrieved from the /patientappointmentreasons call.
	ReasonID *string
	// Set this field in order to set the urgent flag in athena (if the practice settings allow for this).
	Urgent *bool
}

// BookAppointment - Create a single appointment for specific patient
// PUT /v1/{practiceid}/appointments/{appointmentid}
// https://docs.athenahealth.com/api/api-ref/appointment#Book-appointment
func (h *HTTPClient) BookAppointment(ctx context.Context, apptID string, opts *BookAppointmentOptions) (*BookedAppointment, error) {
	var out []*BookedAppointment

	form := url.Values{}

	if opts != nil {
		if opts.AppointmentTypeID != nil && *opts.AppointmentTypeID != "" {
			form.Add("appointmenttypeid", *opts.AppointmentTypeID)
		}

		if opts.BookingNote != nil && *opts.BookingNote != "" {
			form.Add("bookingnote", *opts.BookingNote)
		}

		if opts.DepartmentID != nil && *opts.DepartmentID != "" {
			form.Add("departmentid", *opts.DepartmentID)
		}

		if opts.DoNotSendConfirmationEmail != nil {
			form.Add("donotsendconfirmationemail", strconv.FormatBool(*opts.DoNotSendConfirmationEmail))
		}

		if opts.IgnoreSchedulablePermission != nil {
			form.Add("ignoreschedulablepermission", strconv.FormatBool(*opts.IgnoreSchedulablePermission))
		}

		if opts.NoPatientCase != nil {
			form.Add("nopatientcase", strconv.FormatBool(*opts.NoPatientCase))
		}

		form.Add("patientid", opts.PatientID)

		if opts.ReasonID != nil && *opts.ReasonID != "" {
			form.Add("reasonid", *opts.ReasonID)
		}

		if opts.Urgent != nil {
			form.Add("urgent", strconv.FormatBool(*opts.Urgent))
		}
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/appointments/%s", apptID), form, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("unexpected length returned")
	}

	return out[0], nil
}

type UpdateBookedAppointmentOptions struct {
	// New appointment type ID for this appointment.
	AppointmentTypeID *string `json:"appointmenttypeid"`
	// New department ID for this appointment.
	DepartmentID *string `json:"departmentid"`
	// New provider ID for this appointment.
	ProviderID *string `json:"providerid"`
	// New supervisingprovider ID for this appointment.
	SupervisingProviderID *string `json:"supervisingproviderid"`
}

// https://docs.athenahealth.com/api/api-ref/appointment-booked#Appointment-Booked
// status	string	This subroutine will return 1 on success, and will otherwise return an error message.
var updateBookedApptSuccess = "1"

// UpdateBookedAppointment
// PUT /v1/{practiceid}/appointments/booked/{appointmentid}
// https://docs.athenahealth.com/api/api-ref/appointment-booked#Appointment-Booked
func (h *HTTPClient) UpdateBookedAppointment(ctx context.Context, apptID string, opts *UpdateBookedAppointmentOptions) error {

	form := url.Values{}

	if opts != nil {
		if opts.AppointmentTypeID != nil {
			form.Add("appointmenttypeid", *opts.AppointmentTypeID)
		}

		if opts.DepartmentID != nil {
			form.Add("departmentid", *opts.DepartmentID)
		}

		if opts.ProviderID != nil {
			form.Add("providerid", *opts.ProviderID)
		}

		if opts.SupervisingProviderID != nil {
			form.Add("supervisingproviderid", *opts.SupervisingProviderID)
		}
	}

	var statusRes NumberString
	_, err := h.PutForm(ctx, fmt.Sprintf("/appointments/booked/%s", apptID), form, &statusRes)
	if err != nil {
		return err
	}

	if string(statusRes) != updateBookedApptSuccess {
		return errors.New(string(statusRes))
	}

	return nil
}

type UseExpectedProcedureCodes struct {
	// The ID of the code
	ProcedureCode string `json:"procedurecode"`
	// The description of the code
	ProcedureCodeDescription string `json:"procedurecodedescription"`
}

type RescheduleAppointmentResult struct {
	// Detailed information about the copay for this appointment. Gives more detail than the COPAY field. Note: this information is not yet available in all practices, we are rolling this out slowly.
	AppointmentCopay string `json:"appointmentcopay"`
	// Appointment ID of the booked appointment
	AppointmentID string `json:"appointmentid"`
	// The athenaNet appointment status. There are several possible statuses. x=cancelled f=future o=open. 2=checked in 3=checked out 4=charge entered
	AppointmentStatus AppointmentStatus `json:"appointmentstatus"`
	// The practice-friendly (not patient friendly) name for this appointment type. Note that this may not be the same as the booked appointment because of "generic" slots.
	AppointmentType string `json:"appointmenttype"`
	// This is the ID for the appointment type. Note that this may not be the same as the booked appointment because of "generic" slots.
	AppointmentTypeID string `json:"appointmenttypeid"`
	// As detailed in /claims, if requested.
	Claims []*Claim `json:"claims"`
	// Expected copay for this appointment. Based on the appointment type, the patient's primary insurance, and any copays collected. To see the amounts used in this calculated value, see the APPOINTMENTCOPAY fields.
	Copay string `json:"copay"`
	// The appointment date.
	Date string `json:"date"`
	// The athenaNet department id
	DepartmentID string `json:"departmentid"`
	// In minutes
	Duration int `json:"duration"`
	// If true, this appointment slot is frozen
	FrozenYN string `json:"frozenyn"`
	// This is the raw provider ID that should be used ONLY if using this appointment in conjunction with an HL7 message and with athenahealth's prior guidance. It is only available in some situations.
	HL7ProviderID int `json:"hl7providerid"`
	// As detailed in /patients, if requested.
	Patient string `json:"patient"`
	// The patient-friendly name for this appointment type. Note that this may not be the same as the booked appointment because of "generic" slots.
	PatientAppointmentTypeName string `json:"patientappointmenttypename"`
	// The athenaNet patient ID for this appointment
	PatientID string `json:"patientid"`
	// The athenaNet provider ID
	ProviderID string `json:"providerid"`
	// The referring provider ID.
	ReferringProviderID string `json:"referringproviderid"`
	// The rendering provider ID.
	RenderingProviderID string `json:"renderingproviderid"`
	// When an appointment is rescheduled, this is the ID of the replacement appointment.
	RescheduledAppointmentID string `json:"rescheduledappointmentid"`
	// The timestamp when the appointment started the check in process. If this is set while an appointment is still in status 'f', it means that the check-in process has begun but is not yet completed.
	StartCheckIn string `json:"startcheckin"`
	// As HH:MM (where HH is the 0-23 hour and MM is the minute). This time is local to the department.
	StartTime string `json:"starttime"`
	// The timestamp when the check-in process was finished for this appointment.
	StopCheckIn string `json:"stopcheckin"`
	// The supervising provider ID.
	SupervisingProviderID string `json:"supervisingproviderid"`
	// Urgent flag for the appointment.
	UrgentYN string `json:"urgentyn"`
	// An array of expected procedure codes attached to this appointment.
	UseExpectedProcedureCodes []UseExpectedProcedureCodes `json:"useexpectedprocedurecodes"`
	// Visit ID of the appointment. The VISITID property will only be visible if the following rollout toggle is ON : COLDEN_APPOINTMENT_WITH_VISITID_MDP_API
	VisitID string `json:"visitid"`
}

type RescheduleAppointmentOptions struct {
	// The appointment cancel reason id for cancellation of the original appointment. Use GET /appointmentcancelreasons to retrieve a list of cancel reasons.
	AppointmentCancelReasonID *string `json:"appointmentcancelreasonid"`
	// By default, we allow booking of appointments marked as schedulable via the web. This flag allows you to bypass that restriction for booking.
	IgnoreSchedulablePermission *bool `json:"ignoreschedulablepermission"`
	// The appointment ID of the new appointment. (The appointment ID in the URL is the ID of the currently scheduled appointment.)
	NewAppointmentID string `json:"newappointmentid"`
	// By default, we create a patient case upon booking an appointment for new patients. Setting this to true bypasses that patient case.
	NoPatientCase *bool `json:"nopatientcase"`
	// The athenaNet patient ID.
	PatientID string `json:"patientid"`
	// The appointment reason ID to be booked. If not provided, the same reason used in the original appointment will be used.
	ReasonID *string `json:"reasonid"`
	// A text explanation why the appointment is being rescheduled
	RescheduleReason *string `json:"reschedulereason"`
}

// RescheduleAppointment - Reschedule an existing appointment
// PUT /v1/{practiceid}/appointments/{appointmentid}/reschedule
// https://docs.athenahealth.com/api/api-ref/appointment#Reschedule-appointment
func (h *HTTPClient) RescheduleAppointment(ctx context.Context, apptID string, opts *RescheduleAppointmentOptions) (*RescheduleAppointmentResult, error) {
	var out []*RescheduleAppointmentResult

	form := url.Values{}

	if opts != nil {
		if opts.AppointmentCancelReasonID != nil && *opts.AppointmentCancelReasonID != "" {
			form.Set("appointmentcancelreasonid", *opts.AppointmentCancelReasonID)
		}

		if opts.IgnoreSchedulablePermission != nil {
			form.Set("ignoreschedulablepermission", strconv.FormatBool(*opts.IgnoreSchedulablePermission))
		}

		form.Set("newappointmentid", opts.NewAppointmentID)

		if opts.NoPatientCase != nil {
			form.Set("nopatientcase", strconv.FormatBool(*opts.NoPatientCase))
		}

		form.Set("patientid", opts.PatientID)

		if opts.ReasonID != nil && *opts.ReasonID != "" {
			form.Set("reasonid", *opts.ReasonID)
		}

		if opts.RescheduleReason != nil && *opts.RescheduleReason != "" {
			form.Set("reschedulereason", *opts.RescheduleReason)
		}
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/appointments/%s/reschedule", apptID), form, &out)

	if err != nil {
		return nil, err
	}

	return out[0], nil
}
