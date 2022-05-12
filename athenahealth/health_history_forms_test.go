package athenahealth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetHealthHistoryFormForAppointment(t *testing.T) {
	assert := assert.New(t)

	apptID := "123"
	formID := "1"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal(fmt.Sprintf("/appointments/%s/healthhistoryforms/%s", apptID, formID), r.URL.String())

		b, _ := ioutil.ReadFile("./resources/GetHealthHistoryFormForAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	hhf, err := athenaClient.GetHealthHistoryFormForAppointment(context.Background(), apptID, formID)

	assert.NotNil(hhf)
	assert.NoError(err)
}

func TestHTTPClient_UpdateHealthHistoryFormForAppointment(t *testing.T) {
	assert := assert.New(t)

	apptID := "123"
	formID := "1"

	hhfBytes, err := ioutil.ReadFile("./resources/GetHealthHistoryFormForAppointment.json")
	assert.NoError(err)

	var sections []json.RawMessage
	err = json.Unmarshal(hhfBytes, &sections)
	assert.NoError(err)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal(fmt.Sprintf("/appointments/%s/healthhistoryforms/%s", apptID, formID), r.URL.String())
		assert.Equal(string(hhfBytes), r.FormValue("healthhistoryform"))

		b, _ := ioutil.ReadFile("./resources/UpdateHealthHistoryFormForAppointment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	hhf := &HealthHistoryForm{}
	err = hhf.fromSections(sections)
	assert.NoError(err)

	err = athenaClient.UpdateHealthHistoryFormForAppointment(context.Background(), apptID, formID, hhf)
	assert.NoError(err)
}
