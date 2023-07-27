package athenahealth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetAppointment(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/GetAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	appointment, err := athenaClient.GetAppointment(context.Background(), "1")

	assert.NotNil(appointment)
	assert.NoError(err)
}

func TestHTTPClient_ListAppointmentCustomFields(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/ListAppointmentCustomFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	customFields, err := athenaClient.ListAppointmentCustomFields(context.Background())

	assert.Len(customFields, 2)
	assert.NoError(err)
}

func TestHTTPClient_ListBookedAppointments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("providerid"))
		assert.Equal("06/01/2020", r.URL.Query().Get("startdate"))
		assert.Equal("06/03/2020", r.URL.Query().Get("enddate"))
		assert.Equal(AppointmentStatusCancelled.String(), r.URL.Query().Get("appointmentstatus"))

		b, _ := os.ReadFile("./resources/ListBookedAppointments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListBookedAppointmentsOptions{
		ProviderID:        PtrStr("1"),
		StartDate:         time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:           time.Date(2020, 6, 3, 0, 0, 0, 0, time.UTC),
		AppointmentStatus: func() *AppointmentStatus { a := AppointmentStatusCancelled; return &a }(),
	}

	res, err := athenaClient.ListBookedAppointments(context.Background(), opts)

	assert.Len(res.BookedAppointments, 2)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 2)
	assert.NoError(err)
}

func TestHTTPClient_ListChangedAppointments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("providerid"))
		assert.Equal("06/01/2020 15:30:45", r.URL.Query().Get("showprocessedstartdatetime"))
		assert.Equal("06/02/2020 12:30:45", r.URL.Query().Get("showprocessedenddatetime"))

		b, _ := os.ReadFile("./resources/ListChangedAppointments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListChangedAppointmentsOptions{
		ProviderID:                 "1",
		ShowProcessedStartDatetime: time.Date(2020, 6, 1, 15, 30, 45, 0, time.UTC),
		ShowProcessedEndDatetime:   time.Date(2020, 6, 2, 12, 30, 45, 0, time.UTC),
	}

	appointments, err := athenaClient.ListChangedAppointments(context.Background(), opts)

	assert.Len(appointments, 2)
	assert.NoError(err)
}

func TestHTTPClient_CreateAppointmentNote(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "notetext=test+note")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &CreateAppointmentNoteOptions{
		AppointmentID: "1",
		NoteText:      "test note",
	}

	err := athenaClient.CreateAppointmentNote(context.Background(), "1", opts)

	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_ListAppointmentNotes(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("appointmentid"))

		b, _ := os.ReadFile("./resources/ListAppointmentNotes.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListAppointmentNotesOptions{
		AppointmentID: "1",
	}

	appointments, err := athenaClient.ListAppointmentNotes(context.Background(), "1", opts)

	assert.Len(appointments, 2)
	assert.NoError(err)
}

func TestHTTPClient_UpdateAppointmentNote(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "notetext=test+note")
		assert.Contains(string(reqBody), "noteid=2")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &UpdateAppointmentNoteOptions{
		AppointmentID: "1",
		NoteID:        "2",
		NoteText:      "test note",
	}

	err := athenaClient.UpdateAppointmentNote(context.Background(), "1", "2", opts)

	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_DeleteAppointmentNote(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "noteid=1")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &DeleteAppointmentNoteOptions{
		AppointmentID: "1",
		NoteID:        "1",
	}

	err := athenaClient.DeleteAppointmentNote(context.Background(), "1", "1", opts)

	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_ListOpenAppointmentSlots(t *testing.T) {
	assert := assert.New(t)

	endDate := time.Now()
	startDate := time.Now()

	opts := &ListOpenAppointmentSlotOptions{
		AppointmentTypeID:           PtrStr("1"),
		BypassScheduleTimeChecks:    PtrBool(true),
		DepartmentIDs:               []string{"1"},
		EndDate:                     &endDate,
		IgnoreSchedulablePermission: PtrBool(false),
		ProviderIDs:                 []string{"4", "5"},
		ReasonIDs:                   []string{"2", "3"},
		ShowFrozenSlots:             PtrBool(true),
		StartDate:                   &startDate,

		Limit:  6,
		Offset: 7,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(*opts.AppointmentTypeID, r.URL.Query().Get("appointmenttypeid"))
		assert.Equal(strconv.FormatBool(*opts.BypassScheduleTimeChecks), r.URL.Query().Get("bypassscheduletimechecks"))
		assert.Equal(strings.Join(opts.DepartmentIDs, ","), r.URL.Query().Get("departmentid"))
		assert.Equal((*opts.EndDate).Format("01/02/2006"), r.URL.Query().Get("enddate"))
		assert.Equal(strconv.FormatBool(*opts.IgnoreSchedulablePermission), r.URL.Query().Get("ignoreschedulablepermission"))
		assert.Equal(strings.Join(opts.ProviderIDs, ","), r.URL.Query().Get("providerid"))
		assert.Equal(strings.Join(opts.ReasonIDs, ","), r.URL.Query().Get("reasonid"))
		assert.Equal(strconv.FormatBool(*opts.ShowFrozenSlots), r.URL.Query().Get("showfrozenslots"))
		assert.Equal((*opts.StartDate).Format("01/02/2006"), r.URL.Query().Get("startdate"))

		assert.Equal(strconv.Itoa(opts.Limit), r.URL.Query().Get("limit"))
		assert.Equal(strconv.Itoa(opts.Offset), r.URL.Query().Get("offset"))

		b, _ := os.ReadFile("./resources/ListOpenAppointmentSlots.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	slotsRes, err := athenaClient.ListOpenAppointmentSlots(context.Background(), opts)

	assert.Len(slotsRes.Appointments, 237)
	assert.NoError(err)
}

func TestHTTPClient_BookAppointment(t *testing.T) {
	assert := assert.New(t)

	apptID := "2"

	opts := &BookAppointmentOptions{
		AppointmentTypeID:           PtrStr("3"),
		BookingNote:                 PtrStr("Hello World!"),
		DepartmentID:                PtrStr("4"),
		DoNotSendConfirmationEmail:  PtrBool(true),
		IgnoreSchedulablePermission: PtrBool(false),
		NoPatientCase:               PtrBool(true),
		PatientID:                   "1",
		ReasonID:                    PtrStr("5"),
		Urgent:                      PtrBool(false),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("appointmenttypeid"), *opts.AppointmentTypeID)
		assert.Equal(r.Form.Get("bookingnote"), *opts.BookingNote)
		assert.Equal(r.Form.Get("departmentid"), *opts.DepartmentID)
		assert.Equal(r.Form.Get("donotsendconfirmationemail"), strconv.FormatBool(*opts.DoNotSendConfirmationEmail))
		assert.Equal(r.Form.Get("ignoreschedulablepermission"), strconv.FormatBool(*opts.IgnoreSchedulablePermission))
		assert.Equal(r.Form.Get("nopatientcase"), strconv.FormatBool(*opts.NoPatientCase))
		assert.Equal(r.Form.Get("patientid"), "1")
		assert.Equal(r.Form.Get("reasonid"), "5")
		assert.Equal(r.Form.Get("urgent"), strconv.FormatBool(*opts.Urgent))

		b, _ := os.ReadFile("./resources/BookAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	_, err := athenaClient.BookAppointment(context.Background(), apptID, opts)

	assert.NoError(err)
}

func TestHTTPClient_UpdateBookedAppointment_IntResponse(t *testing.T) {
	assert := assert.New(t)

	apptID := "1230322"

	opts := &UpdateBookedAppointmentOptions{
		AppointmentTypeID:     PtrStr("*opts.AppointmentTypeID"),
		DepartmentID:          PtrStr("*opts.DepartmentID"),
		ProviderID:            PtrStr("*opts.ProviderID"),
		SupervisingProviderID: PtrStr("*opts.SupervisingProviderID"),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("appointmenttypeid"), *opts.AppointmentTypeID)
		assert.Equal(r.Form.Get("departmentid"), *opts.DepartmentID)
		assert.Equal(r.Form.Get("providerid"), *opts.ProviderID)
		assert.Equal(r.Form.Get("supervisingproviderid"), *opts.SupervisingProviderID)

		assert.Equal(r.URL.Path, fmt.Sprintf("/appointments/booked/%s", apptID))

		b, _ := os.ReadFile("./resources/UpdateBookedAppointment_IntResponse.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	updateErr := athenaClient.UpdateBookedAppointment(context.Background(), apptID, opts)
	assert.NoError(updateErr)
}

func TestHTTPClient_UpdateBookedAppointment_StringResponse(t *testing.T) {
	assert := assert.New(t)

	apptID := "1230322"

	opts := &UpdateBookedAppointmentOptions{
		AppointmentTypeID:     PtrStr("opts.AppointmentTypeID"),
		DepartmentID:          PtrStr("opts.DepartmentID"),
		ProviderID:            PtrStr("opts.ProviderID"),
		SupervisingProviderID: PtrStr("opts.SupervisingProviderID"),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("appointmenttypeid"), *opts.AppointmentTypeID)
		assert.Equal(r.Form.Get("departmentid"), *opts.DepartmentID)
		assert.Equal(r.Form.Get("providerid"), *opts.ProviderID)
		assert.Equal(r.Form.Get("supervisingproviderid"), *opts.SupervisingProviderID)

		assert.Equal(r.URL.Path, fmt.Sprintf("/appointments/booked/%s", apptID))

		b, _ := os.ReadFile("./resources/UpdateBookedAppointment_StringResponse.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	updateErr := athenaClient.UpdateBookedAppointment(context.Background(), apptID, opts)
	assert.Error(updateErr)
	assert.Equal(updateErr.Error(), "Invalid PROVIDERID input")
}

func TestHTTPClient_RescheduleAppointment(t *testing.T) {
	assert := assert.New(t)

	apptID := "998877"

	opts := &RescheduleAppointmentOptions{
		AppointmentCancelReasonID:   PtrStr("2"),
		IgnoreSchedulablePermission: PtrBool(true),
		NewAppointmentID:            "123",
		NoPatientCase:               PtrBool(false),
		PatientID:                   "456",
		ReasonID:                    PtrStr("3"),
		RescheduleReason:            PtrStr("other commitments"),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("appointmentcancelreasonid"), *opts.AppointmentCancelReasonID)
		assert.Equal(r.Form.Get("ignoreschedulablepermission"), strconv.FormatBool(*opts.IgnoreSchedulablePermission))
		assert.Equal(r.Form.Get("newappointmentid"), opts.NewAppointmentID)
		assert.Equal(r.Form.Get("nopatientcase"), strconv.FormatBool(*opts.NoPatientCase))
		assert.Equal(r.Form.Get("patientid"), opts.PatientID)
		assert.Equal(r.Form.Get("reasonid"), *opts.ReasonID)
		assert.Equal(r.Form.Get("reschedulereason"), *opts.RescheduleReason)
		assert.Equal(r.URL.Path, fmt.Sprintf("/appointments/%s/reschedule", apptID))

		b, _ := os.ReadFile("./resources/RescheduleAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	rescheduleAppointmentResult, err := athenaClient.RescheduleAppointment(context.Background(), apptID, opts)
	assert.NotNil(rescheduleAppointmentResult)
	assert.NoError(err)
	assert.Equal("25.00", rescheduleAppointmentResult.AppointmentCopay)
	assert.Equal("APT12345", rescheduleAppointmentResult.AppointmentID)
	assert.Equal(AppointmentStatusFuture, rescheduleAppointmentResult.AppointmentStatus)
	assert.Equal("Checkup", rescheduleAppointmentResult.AppointmentType)
	assert.Equal("ATT12345", rescheduleAppointmentResult.AppointmentTypeID)
	assert.Len(rescheduleAppointmentResult.Claims, 1)
	assert.Equal("25.00", rescheduleAppointmentResult.Copay)
	assert.Equal("06/20/2023", rescheduleAppointmentResult.Date)
	assert.Equal("DPT12345", rescheduleAppointmentResult.DepartmentID)
	assert.Equal(30, rescheduleAppointmentResult.Duration)
	assert.Equal("true", rescheduleAppointmentResult.FrozenYN)
	assert.Equal(764012, rescheduleAppointmentResult.HL7ProviderID)
	assert.Equal("John Doe", rescheduleAppointmentResult.Patient)
	assert.Equal("Checkup", rescheduleAppointmentResult.PatientAppointmentTypeName)
	assert.Equal("P12345", rescheduleAppointmentResult.PatientID)
	assert.Equal("PR12345", rescheduleAppointmentResult.ProviderID)
	assert.Equal("RFP12345", rescheduleAppointmentResult.ReferringProviderID)
	assert.Equal("RNP12345", rescheduleAppointmentResult.RenderingProviderID)
	assert.Equal("APT98765", rescheduleAppointmentResult.RescheduledAppointmentID)
	assert.Equal("15:00", rescheduleAppointmentResult.StartCheckIn)
	assert.Equal("15:30", rescheduleAppointmentResult.StartTime)
	assert.Equal("15:10", rescheduleAppointmentResult.StopCheckIn)
	assert.Equal("SPP12345", rescheduleAppointmentResult.SupervisingProviderID)
	assert.Equal("true", rescheduleAppointmentResult.UrgentYN)
	assert.Len(rescheduleAppointmentResult.UseExpectedProcedureCodes, 1)
	assert.Equal("V12345", rescheduleAppointmentResult.VisitID)
}
