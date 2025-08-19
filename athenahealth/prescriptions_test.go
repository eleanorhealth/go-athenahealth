package athenahealth

import (
	"context"
	"fmt"
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

func TestHTTPClient_UpdatePrescriptionActionNote(t *testing.T) {
	assert := assert.New(t)

	departmentID := 1
	patientID := 67890
	documentID := 12345
	actionNote := "test"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())
		assert.Equal(strconv.Itoa(departmentID), r.Form.Get("departmentid"))
		assert.Equal(actionNote, r.Form.Get("actionnote"))
		assert.Empty(r.Form.Get("settotop"))

		// Extract patientID and documentID from the URL path
		// Path: /patients/{patientid}/documents/prescriptions/{documentid}
		var gotPatientID, gotDocumentID int
		_, err := fmt.Sscanf(r.URL.Path, "/patients/%d/documents/prescriptions/%d", &gotPatientID, &gotDocumentID)
		assert.NoError(err)
		assert.Equal(patientID, gotPatientID)
		assert.Equal(documentID, gotDocumentID)

		// Write a minimal valid JSON response directly
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.UpdatePrescription(context.Background(), departmentID, patientID, documentID, &UpdatePrescriptionOptions{ActionNote: &actionNote})
	assert.NoError(err)
	assert.True(res.Success)
}

func TestHTTPClient_UpdatePrescriptionActionNote_Error(t *testing.T) {
	assert := assert.New(t)

	departmentID := 1
	patientID := 67890
	documentID := 12345
	actionNote := "test"

	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.UpdatePrescription(context.Background(), departmentID, patientID, documentID, &UpdatePrescriptionOptions{ActionNote: &actionNote})
	assert.Error(err)
	assert.False(res.Success)
	assert.NotNil(res.ErrorMessage)
	assert.Contains(*res.ErrorMessage, "Bad Request")
}

func TestHTTPClient_UpdatePrescription_OtherFields_Success(t *testing.T) {
	assert := assert.New(t)

	departmentID := 2
	patientID := 11111
	documentID := 22222

	assignedTo := "user123"
	internalNote := "internal-only"
	note := "visible note"
	pin := true

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())
		assert.Equal(strconv.Itoa(departmentID), r.Form.Get("departmentid"))
		assert.Equal(assignedTo, r.Form.Get("assignedto"))
		assert.Equal(internalNote, r.Form.Get("internalnote"))
		assert.Equal(note, r.Form.Get("note"))
		assert.Equal(strconv.FormatBool(pin), r.Form.Get("pintotop"))

		var gotPatientID, gotDocumentID int
		_, err := fmt.Sscanf(r.URL.Path, "/patients/%d/documents/prescriptions/%d", &gotPatientID, &gotDocumentID)
		assert.NoError(err)
		assert.Equal(patientID, gotPatientID)
		assert.Equal(documentID, gotDocumentID)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.UpdatePrescription(context.Background(), departmentID, patientID, documentID, &UpdatePrescriptionOptions{
		AssignedTo:   &assignedTo,
		InternalNote: &internalNote,
		Note:         &note,
		PinToTop:     &pin,
	})
	assert.NoError(err)
	assert.True(res.Success)
}

func TestHTTPClient_UpdatePrescription_PinToTopFalse_Success(t *testing.T) {
	assert := assert.New(t)

	departmentID := 3
	patientID := 33333
	documentID := 44444
	pin := false

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())
		assert.Equal(strconv.Itoa(departmentID), r.Form.Get("departmentid"))
		assert.Equal("false", r.Form.Get("pintotop"))

		var gotPatientID, gotDocumentID int
		_, err := fmt.Sscanf(r.URL.Path, "/patients/%d/documents/prescriptions/%d", &gotPatientID, &gotDocumentID)
		assert.NoError(err)
		assert.Equal(patientID, gotPatientID)
		assert.Equal(documentID, gotDocumentID)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.UpdatePrescription(context.Background(), departmentID, patientID, documentID, &UpdatePrescriptionOptions{
		PinToTop: &pin,
	})
	assert.NoError(err)
	assert.True(res.Success)
}
