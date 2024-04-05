package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetTelehealthInviteURL(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/GetTelehealthInviteURL.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	result, err := athenaClient.GetTelehealthInviteURL(context.Background(), "12588")

	assert.NotNil(result)
	assert.NoError(err)
	assert.Equal("12588", result.AppointmentID)
}
