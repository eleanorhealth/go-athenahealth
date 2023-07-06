package athenahealth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestHTTPClient_CancelCheckInAppointment(t *testing.T) {
	assert := assert.New(t)

	appointmentID := "54"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/appointments/54/cancelcheckin")
		b, _ := os.ReadFile("./resources/CancelCheckInAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.CancelCheckInAppointment(context.Background(), appointmentID)

	assert.NoError(err)
}

func TestHTTPClient_CheckInAppointment(t *testing.T) {
	assert := assert.New(t)

	appointmentID := "54"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/appointments/54/checkin")
		b, _ := os.ReadFile("./resources/CheckInAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.CheckInAppointment(context.Background(), appointmentID)

	assert.NoError(err)
}

func TestHTTPClient_GetRequiredCheckInFields(t *testing.T) {
	assert := assert.New(t)

	departmentID := "45"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/departments/45/checkinrequired")
		b, _ := os.ReadFile("./resources/GetRequiredCheckInFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.GetRequiredCheckInFields(context.Background(), departmentID)

	assert.NoError(err)
	assert.True(reflect.DeepEqual(res.FieldList, []string{"1", "2", "3", "4", "5"}))
}

func TestHTTPClient_StartCheckInAppointment(t *testing.T) {
	assert := assert.New(t)

	appointmentID := "54"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/appointments/54/startcheckin")
		b, _ := os.ReadFile("./resources/StartCheckInAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.StartCheckInAppointment(context.Background(), appointmentID)

	assert.NoError(err)
}
