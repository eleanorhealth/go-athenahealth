package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListProblems(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("d5", r.URL.Query().Get("departmentid"))
		assert.Equal("p1", r.URL.Query().Get("patientid"))
		assert.Equal("true", r.URL.Query().Get("showdiagnosisinfo"))

		b, _ := os.ReadFile("./resources/ListProblems.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListProblemsOptions{
		DepartmentID:      "d5",
		PatientID:         "p1",
		ShowDiagnosisInfo: true,
	}

	problems, err := athenaClient.ListProblems(context.Background(), opts.PatientID, opts)

	assert.Len(problems, 2)
	assert.NoError(err)
}

func TestHTTPClient_ListChangedProblems(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("leaveunprocessed"))
		assert.Equal("p1", r.URL.Query().Get("patientid"))
		assert.Equal("06/01/2020 15:30:45", r.URL.Query().Get("showprocessedstartdatetime"))
		assert.Equal("06/02/2020 12:30:45", r.URL.Query().Get("showprocessedenddatetime"))

		b, _ := os.ReadFile("./resources/ListChangedProblems.json")
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

	problems, err := athenaClient.ListChangedProblems(context.Background(), opts)

	assert.Len(problems, 1)
	assert.NoError(err)
}

func TestProblem_ICD10Code(t *testing.T) {
	tests := []struct {
		name    string
		problem *Problem
		want    string
	}{
		{
			name: "use code when codeset ICD10",
			problem: &Problem{
				Codeset:            "ICD10",
				Code:               "F49.1",
				BestMatchICD10Code: "F49",
			},
			want: "F49.1",
		},
		{
			name: "use bestmatchicd10code when codeset not ICD10",
			problem: &Problem{
				Codeset:            "SNOMED",
				Code:               "123456",
				BestMatchICD10Code: "F49",
			},
			want: "F49",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.problem.ICD10Code(); got != tt.want {
				t.Errorf("Problem.ICD10Code() = %v, want %v", got, tt.want)
			}
		})
	}
}
