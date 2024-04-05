package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_EncounterSummary(t *testing.T) {
	assert := assert.New(t)

	encounterID := "123"
	opts := &EncounterSummaryOptions{
		Mobile: true,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("mobile"))

		b, _ := os.ReadFile("./resources/EncounterSummary.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	summary, err := athenaClient.EncounterSummary(context.Background(), encounterID, opts)
	assert.NoError(err)
	assert.NotNil(summary)
	assert.Contains(summary.Summary, "Summary of the report")
}
