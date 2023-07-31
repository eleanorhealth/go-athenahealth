package athenahealth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestHTTPClient_ListAppointmentReminders(t *testing.T) {
	assert := assert.New(t)

	appointmentTypeID := 123
	patientID := 456
	providerID := 789
	showDeleted := false

	opts := &ListAppointmentRemindersOptions{
		StartDate:         time.Date(2011, 9, 22, 0, 0, 0, 0, time.UTC),
		EndDate:           time.Date(2022, 2, 25, 0, 0, 0, 0, time.UTC),
		DepartmentID:      "1",
		AppointmentTypeID: &appointmentTypeID,
		PatientID:         &patientID,
		ProviderID:        &providerID,
		ShowDeleted:       &showDeleted,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal("09/22/2011", query.Get("startdate"))
		assert.Equal("02/25/2022", query.Get("enddate"))
		assert.Equal(opts.DepartmentID, query.Get("departmentid"))

		assert.Equal("123", query.Get("appointmenttypeid"))
		assert.Equal("456", query.Get("patientid"))
		assert.Equal("789", query.Get("providerid"))
		assert.Equal("false", query.Get("showdeleted"))

		b, _ := os.ReadFile("./resources/ListAppointmentReminders.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	apppointmentRemindersResult, err := athenaClient.ListAppointmentReminders(context.Background(), opts)
	assert.NoError(err)
	assert.Len(apppointmentRemindersResult.Reminders, 2)
}
