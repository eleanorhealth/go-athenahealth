package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListChangedProblems(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("leaveunprocessed"))
		assert.Equal("p1", r.URL.Query().Get("patientid"))
		assert.Equal("06/01/2020 15:30:45", r.URL.Query().Get("showprocessedstartdatetime"))
		assert.Equal("06/02/2020 12:30:45", r.URL.Query().Get("showprocessedenddatetime"))

		b, _ := ioutil.ReadFile("./resources/ListChangedProblems.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListChangedProblemsOptions{
		LeaveUnprocessed:           true,
		PatientID:                  "p1",
		ShowProcessedStartDatetime: time.Date(2020, 6, 1, 15, 30, 45, 0, time.UTC),
		ShowProcessedEndDatetime:   time.Date(2020, 6, 2, 12, 30, 45, 0, time.UTC),
	}

	problems, err := athenaClient.ListChangedProblems(opts)

	assert.Len(problems, 1)
	assert.NoError(err)
}
