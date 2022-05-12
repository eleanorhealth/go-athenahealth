package athenahealth

import (
	"context"
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
		assert.Equal(fmt.Sprintf("/appointments/%s/healthhistoryforms/%s", apptID, formID), r.URL.String())

		b, _ := ioutil.ReadFile("./resources/HealthHistoryForm.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	hhf, err := athenaClient.GetHealthHistoryFormForAppointment(context.Background(), apptID, formID)

	assert.NotNil(hhf)
	assert.NoError(err)

}
