package athenahealth

import (
	"context"
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

func TestHTTPClient_ListAppointmentReasons(t *testing.T) {
	assert := assert.New(t)

	deptID := "1"
	providerID := "2"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(deptID, r.URL.Query().Get("departmentid"))
		assert.Equal(providerID, r.URL.Query().Get("providerid"))

		b, _ := os.ReadFile("./resources/ListAppointmentReasons.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	reasons, err := athenaClient.ListAppointmentReasons(context.Background(), deptID, providerID)

	assert.Len(reasons, 2)
	assert.NoError(err)
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
		ProviderID:                  4,
		StartDate:                   startDate,
		IgnoreSchedulablePermission: true,
		ShowFrozenSlots:             true,
		Limit:                       5,
		Offset:                      6,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(strconv.Itoa(deptID), r.URL.Query().Get("departmentid"))

		b, _ := os.ReadFile("./resources/ListOpenAppointmentSlots.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	slots, err := athenaClient.ListOpenAppointmentSlots(context.Background(), deptID, opts)

	assert.Len(slots, 237)
	assert.NoError(err)
}

func TestHTTPClient_BookAppointment(t *testing.T) {
	assert := assert.New(t)

	patientID := 1
	apptID := 2

	h := func(w http.ResponseWriter, r *http.Request) {
		// assert.Equal(strconv.Itoa(deptID), r.URL.Query().Get("departmentid"))

		b, _ := os.ReadFile("./resources/BookAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

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

	_, err := athenaClient.BookAppointment(context.Background(), patientID, apptID, opts)

	assert.NoError(err)
}
