package athenahealth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestHTTPClient_CancelCheckInAppointment(t *testing.T) {
	assert := assert.New(t)

	appointmentID := "54"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/appointments/54/cancelcheckin")
		b, _ := os.ReadFile("./resources/AppointmentCancelCheckIn.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.AppointmentCancelCheckIn(context.Background(), appointmentID)

	assert.NoError(err)
}

func TestHTTPClient_CheckInAppointment(t *testing.T) {
	assert := assert.New(t)

	appointmentID := "54"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/appointments/54/checkin")
		b, _ := os.ReadFile("./resources/AppointmentCheckIn.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.AppointmentCheckIn(context.Background(), appointmentID)

	assert.NoError(err)
}

func TestHTTPClient_StartCheckInAppointment(t *testing.T) {
	assert := assert.New(t)

	appointmentID := "54"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/appointments/54/startcheckin")
		b, _ := os.ReadFile("./resources/AppointmentStartCheckIn.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.AppointmentStartCheckIn(context.Background(), appointmentID)

	assert.NoError(err)
}
