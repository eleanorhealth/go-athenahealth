package athenahealth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListSocialHistoryTemplates(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/ListSocialHistoryTemplates.json")
		w.Write(b)
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	templates, err := athenaClient.ListSocialHistoryTemplates(context.Background())
	assert.NoError(err)

	assert.Len(templates, 2)
}

func TestHTTPClient_GetPatientSocialHistory(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("2", r.URL.Query().Get("departmentid"))
		assert.Equal("patient", r.URL.Query().Get("recipientcategory"))
		assert.Equal("true", r.URL.Query().Get("shownotperformedquestions"))
		assert.Equal("true", r.URL.Query().Get("showunansweredquestions"))

		b, _ := os.ReadFile("./resources/GetPatientSocialHistory.json")
		w.Write(b)
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	opts := &GetPatientSocialHistoryOptions{
		DepartmentID:              "2",
		RecipientCategory:         "patient",
		ShowNotPerformedQuestions: true,
		ShowUnansweredQuestions:   true,
	}

	socialHistory, err := athenaClient.GetPatientSocialHistory(context.Background(), "1", opts)
	assert.NoError(err)

	assert.Len(socialHistory.Questions, 2)
}

func TestHTTPClient_UpdatePatientSocialHistory(t *testing.T) {
	assert := assert.New(t)

	questions := []*UpdatePatientSocialHistoryQuestion{
		{
			Key:  "KEY.1",
			Note: "9/10/2020",
		},
	}
	questionBytes, err := json.Marshal(questions)
	assert.NoError(err)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "departmentid=2")
		assert.Contains(string(reqBody), fmt.Sprintf("questions=%s", url.QueryEscape(string(questionBytes))))
		assert.Contains(string(reqBody), "sectionnote=foo")

		called = true
	}

	athenaClient, ts := TestClient(h)
	defer ts.Close()

	opts := &UpdatePatientSocialHistoryOptions{
		DepartmentID: "2",
		Questions:    questions,
		SectionNote:  "foo",
	}

	err = athenaClient.UpdatePatientSocialHistory(context.Background(), "1", opts)
	assert.NoError(err)

	assert.True(called)
}
