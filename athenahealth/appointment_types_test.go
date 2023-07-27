package athenahealth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestHTTPClient_CreateAppointmentType(t *testing.T) {
	assert := assert.New(t)

	opts := &CreateAppointmentTypeOptions{
		Duration:         "99",
		Generic:          ptrBool(true),
		Name:             "1001",
		Patient:          true,
		ShortName:        "11",
		TemplateTypeOnly: ptrBool(true),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("duration"), opts.Duration)
		assert.Equal(r.Form.Get("generic"), strconv.FormatBool(*opts.Generic))
		assert.Equal(r.Form.Get("name"), opts.Name)
		assert.Equal(r.Form.Get("patient"), strconv.FormatBool(opts.Patient))
		assert.Equal(r.Form.Get("shortname"), opts.ShortName)
		assert.Equal(r.Form.Get("templatetypeonly"), strconv.FormatBool(*opts.TemplateTypeOnly))
		assert.Equal(r.URL.Path, "/appointmenttypes")
		b, _ := os.ReadFile("./resources/CreateAppointmentType.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	createAppointmentTypeResult, err := athenaClient.CreateAppointmentType(context.Background(), opts)

	assert.NotNil(createAppointmentTypeResult)
	assert.NoError(err)
	assert.Equal(5, createAppointmentTypeResult.AppointmentTypeID)
}
