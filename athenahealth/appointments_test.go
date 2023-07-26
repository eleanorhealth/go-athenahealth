package athenahealth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
		assert.Equal("x", r.URL.Query().Get("appointmentstatus"))

		b, _ := os.ReadFile("./resources/ListBookedAppointments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListBookedAppointmentsOptions{
		ProviderID:        "1",
		StartDate:         time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:           time.Date(2020, 6, 3, 0, 0, 0, 0, time.UTC),
		AppointmentStatus: "x",
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

	deptID := 1

	endDate := time.Now()
	startDate := time.Now()

	opts := &ListOpenAppointmentSlotOptions{
		AppointmentTypeID:           1,
		ReasonIDs:                   []int{2, 3},
		BypassScheduleTimeChecks:    true,
		EndDate:                     endDate,
		ProviderIDs:                 []int{4, 5},
		StartDate:                   startDate,
		IgnoreSchedulablePermission: true,
		ShowFrozenSlots:             true,
		Limit:                       6,
		Offset:                      7,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(strconv.Itoa(deptID), r.URL.Query().Get("departmentid"))
		assert.Equal("1", r.URL.Query().Get("appointmenttypeid"))
		assert.Equal("2,3", r.URL.Query().Get("reasonid"))
		assert.Equal("true", r.URL.Query().Get("bypassscheduletimechecks"))
		assert.Equal(endDate.Format("01/02/2006"), r.URL.Query().Get("enddate"))
		assert.Equal("4,5", r.URL.Query().Get("providerid"))
		assert.Equal(startDate.Format("01/02/2006"), r.URL.Query().Get("startdate"))
		assert.Equal("true", r.URL.Query().Get("ignoreschedulablepermission"))
		assert.Equal("true", r.URL.Query().Get("showfrozenslots"))
		assert.Equal("6", r.URL.Query().Get("limit"))
		assert.Equal("7", r.URL.Query().Get("offset"))

		b, _ := os.ReadFile("./resources/ListOpenAppointmentSlots.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	slotsRes, err := athenaClient.ListOpenAppointmentSlots(context.Background(), deptID, opts)

	assert.Len(slotsRes.Appointments, 237)
	assert.NoError(err)
}

func TestHTTPClient_BookAppointment(t *testing.T) {
	assert := assert.New(t)

	patientID := "1"
	apptID := "2"

	opts := &BookAppointmentOptions{
		AppointmentTypeID:           3,
		BookingNote:                 "Hello World!",
		DepartmentID:                4,
		DoNotSendConfirmationEmail:  true,
		IgnoreSchedulablePermission: true,
		NoPatientCase:               true,
		ReasonID:                    5,
		Urgent:                      true,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "appointmenttypeid=3")
		assert.Contains(string(reqBody), "bookingnote=Hello+World%21")
		assert.Contains(string(reqBody), "departmentid=4")
		assert.Contains(string(reqBody), "donotsendconfirmationemail=true")
		assert.Contains(string(reqBody), "ignoreschedulablepermission=true")
		assert.Contains(string(reqBody), "nopatientcase=true")
		assert.Contains(string(reqBody), "reasonid=5")
		assert.Contains(string(reqBody), "urgent=true")

		b, _ := os.ReadFile("./resources/BookAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	_, err := athenaClient.BookAppointment(context.Background(), patientID, apptID, opts)

	assert.NoError(err)
}

func TestHTTPClient_UpdateBookedAppointment_IntResponse(t *testing.T) {
	assert := assert.New(t)

	apptID := "1230322"

	opts := &UpdateBookedAppointmentOptions{
		AppointmentTypeID:     func() *string { a := "opts.AppointmentTypeID"; return &a }(),
		DepartmentID:          func() *string { a := "opts.DepartmentID"; return &a }(),
		ProviderID:            func() *string { a := "opts.ProviderID"; return &a }(),
		SupervisingProviderID: func() *string { a := "opts.SupervisingProviderID"; return &a }(),
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
		AppointmentTypeID:     func() *string { a := "opts.AppointmentTypeID"; return &a }(),
		DepartmentID:          func() *string { a := "opts.DepartmentID"; return &a }(),
		ProviderID:            func() *string { a := "opts.ProviderID"; return &a }(),
		SupervisingProviderID: func() *string { a := "opts.SupervisingProviderID"; return &a }(),
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

	opts := &RescheduleAppointmentOptions{
		AppointmentCancelReasonID:   func() *int { a := 2; return &a }(),
		IgnoreSchedulablePermission: func() *bool { a := true; return &a }(),
		NewAppointmentID:            123,
		NoPatientCase:               func() *bool { a := false; return &a }(),
		PatientID:                   456,
		ReasonID:                    func() *int { a := 3; return &a }(),
		RescheduleReason:            func() *string { a := "Other commitments"; return &a }(),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("appointmentcancelreasonid"), strconv.Itoa(*opts.AppointmentCancelReasonID))
		assert.Equal(r.Form.Get("ignoreschedulablepermission"), strconv.FormatBool(*opts.IgnoreSchedulablePermission))
		assert.Equal(r.Form.Get("newappointmentid"), strconv.Itoa(opts.NewAppointmentID))
		assert.Equal(r.Form.Get("nopatientcase"), strconv.FormatBool(*opts.NoPatientCase))
		assert.Equal(r.Form.Get("patientid"), strconv.Itoa(opts.PatientID))
		assert.Equal(r.Form.Get("reasonid"), strconv.Itoa(*opts.ReasonID))
		assert.Equal(r.Form.Get("reschedulereason"), *opts.RescheduleReason)
		assert.Equal(r.URL.Path, fmt.Sprintf("/appointments/%d/reschedule", 998877))

		b, _ := os.ReadFile("./resources/RescheduleAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	rescheduleAppointmentResult, err := athenaClient.RescheduleAppointment(context.Background(), 998877, opts)
	assert.NotNil(rescheduleAppointmentResult)
	assert.NoError(err)
	assert.Equal("25.00", rescheduleAppointmentResult.AppointmentCopay)
	assert.Equal("APT12345", rescheduleAppointmentResult.AppointmentID)
	assert.Equal("f", rescheduleAppointmentResult.AppointmentStatus)
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
