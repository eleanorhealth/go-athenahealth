package athenahealth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestHTTPClient_CreateAppointmentType(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/CreateAppointmentType.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	createAppointmentTypeResult, err := athenaClient.CreateAppointmentType(context.Background(), nil)

	assert.NotNil(createAppointmentTypeResult)
	assert.NoError(err)
	assert.Equal(5, createAppointmentTypeResult.AppointmentTypeID)
}
