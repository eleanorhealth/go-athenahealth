package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetAppointment(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/GetAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	appointment, err := athenaClient.GetAppointment("1")

	assert.NotNil(appointment)
	assert.Nil(err)
}

func TestHTTPClient_ListAppointmentCustomFields(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/ListAppointmentCustomFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	customFields, err := athenaClient.ListAppointmentCustomFields()

	assert.Len(customFields, 2)
	assert.Nil(err)
}

func TestHTTPClient_ListBookedAppointments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("providerid"))
		assert.Equal("06/01/2020", r.URL.Query().Get("startdate"))
		assert.Equal("06/03/2020", r.URL.Query().Get("enddate"))

		b, _ := ioutil.ReadFile("./resources/ListBookedAppointments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListBookedAppointmentsOptions{
		ProviderID: "1",
		StartDate:  time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2020, 6, 3, 0, 0, 0, 0, time.UTC),
	}

	appointments, err := athenaClient.ListBookedAppointments(opts)

	assert.Len(appointments, 2)
	assert.Nil(err)
}

func TestHTTPClient_ListChangedAppointments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("providerid"))
		assert.Equal("06/01/2020 15:30:45", r.URL.Query().Get("showprocessedstartdatetime"))
		assert.Equal("06/02/2020 12:30:45", r.URL.Query().Get("showprocessedenddatetime"))

		b, _ := ioutil.ReadFile("./resources/ListChangedAppointments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListChangedAppointmentsOptions{
		ProviderID:                 "1",
		ShowProcessedStartDatetime: time.Date(2020, 6, 1, 15, 30, 45, 0, time.UTC),
		ShowProcessedEndDatetime:   time.Date(2020, 6, 2, 12, 30, 45, 0, time.UTC),
	}

	appointments, err := athenaClient.ListChangedAppointments(opts)

	assert.Len(appointments, 2)
	assert.Nil(err)
}

func TestHTTPClient_CreateAppointmentNote(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
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

	err := athenaClient.CreateAppointmentNote("1", opts)

	assert.Nil(err)
	assert.True(called)
}

func TestHTTPClient_ListAppointmentNotes(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("appointmentid"))

		b, _ := ioutil.ReadFile("./resources/ListAppointmentNotes.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListAppointmentNotesOptions{
		AppointmentID: "1",
	}

	appointments, err := athenaClient.ListAppointmentNotes("1", opts)

	assert.Len(appointments, 2)
	assert.Nil(err)
}

func TestHTTPClient_UpdateAppointmentNote(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
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

	err := athenaClient.UpdateAppointmentNote("1", "2", opts)

	assert.Nil(err)
	assert.True(called)
}

func TestHTTPClient_DeleteAppointmentNote(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
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

	err := athenaClient.DeleteAppointmentNote("1", "1", opts)

	assert.Nil(err)
	assert.True(called)
}
