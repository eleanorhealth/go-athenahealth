package athenahealth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHTTPClient_CreateAppointmentSlot(t *testing.T) {
	assert := assert.New(t)

	opts := &CreateAppointmentSlotOptions{
		AppointmentDate:   "03/22/1945",
		AppointmentTime:   []string{"21:42", "01:59"},
		AppointmentTypeID: func() *int { a := 7; return &a }(),
		DepartmentID:      27,
		ProviderID:        99,
		ReasonID:          func() *int { a := 8; return &a }(),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("appointmentdate"), opts.AppointmentDate)
		assert.Equal(r.Form.Get("appointmenttime"), strings.Join(opts.AppointmentTime, ","))
		assert.Equal(r.Form.Get("appointmenttypeid"), strconv.Itoa(*opts.AppointmentTypeID))
		assert.Equal(r.Form.Get("departmentid"), strconv.Itoa(opts.DepartmentID))
		assert.Equal(r.Form.Get("providerid"), strconv.Itoa(opts.ProviderID))
		assert.Equal(r.Form.Get("reasonid"), strconv.Itoa(*opts.ReasonID))
		assert.Equal(r.URL.Path, "/appointments/open")
		b, _ := os.ReadFile("./resources/CreateAppointmentSlot.json")
		w.Write(b)
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	createAppointmentSlotResult, err := athenaClient.CreateAppointmentSlot(context.Background(), opts)

	assert.NotNil(createAppointmentSlotResult)
	assert.NoError(err)
	assert.Equal(1, len(createAppointmentSlotResult.AppointmentIDs))
}
