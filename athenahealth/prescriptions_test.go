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

	ctx := context.Background()

	leaveUnprocessed := false

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(strconv.FormatBool(leaveUnprocessed), r.URL.Query().Get("leaveunprocessed"))

		b, _ := os.ReadFile("./resources/ListChangedPrescriptions.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListChangedPrescriptions(ctx, &ListChangedPrescriptionsOptions{
		LeaveUnprocessed: &leaveUnprocessed,
	})
	prescriptions := res.ChangedPrescriptions

	assert.NoError(err)
	assert.Len(prescriptions, 1)
	assert.Equal(res.Pagination.TotalCount, 1)
}
