package athenahealth

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListChangedPrescriptions(t *testing.T) {
	assert := assert.New(t)

	leaveUnprocessed := true

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(strconv.FormatBool(leaveUnprocessed), r.URL.Query().Get("leaveunprocessed"))

		b, _ := os.ReadFile("./resources/ListChangedPrescriptions.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListChangedPrescriptionsOptions{
		LeaveUnprocessed: leaveUnprocessed,
	}

	res, err := athenaClient.ListChangedPrescriptions(context.Background(), opts)
	prescriptions := res.ChangedPrescriptions

	assert.NoError(err)
	assert.Len(prescriptions, 3)
	assert.Equal(3, res.Pagination.TotalCount)
	assert.Len(prescriptions, res.Pagination.TotalCount)
}

func TestHTTPClient_UpdatePrescription(t *testing.T) {
	assert := assert.New(t)

	prescriptionID := "12345"
	opts := &UpdatePrescriptionOptions{
		PatientID: 67890,
		DocumentID: 12345,
		ActionNote: "test",
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())
		assert.Equal(r.Form.Get("providerid"), opts.ProviderID)
		assert.Equal(r.Form.Get("status"), opts.Status)

		b, _ := os.ReadFile("./resources/UpdatePrescription.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.UpdatePrescription(context.Background(), prescriptionID, opts)
	assert.NoError(err)
}
