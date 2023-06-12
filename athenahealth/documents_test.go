package athenahealth

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListAdminDocuments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(r.URL.Path, "/patients/123/")
		assert.Equal("3", r.URL.Query().Get("departmentid"))

		b, _ := os.ReadFile("./resources/ListAdminDocuments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListAdminDocumentsOptions{
		DepartmentID: "3",
	}

	res, err := athenaClient.ListAdminDocuments(context.Background(), "123", opts)

	assert.Len(res.AdminDocuments, 1)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 1)
	assert.NoError(err)
}

func TestHTTPClient_AddDocument(t *testing.T) {
	assert := assert.New(t)

	actionNote := "test action note"
	apptID := 1
	attachmentContents := []byte("test attachment contents")
	autoclose := "true"
	deptID := 2
	documentSubclass := "ADMIN_CONSENT"
	internalNote := "test internal note"
	providerID := 3

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Contains(r.URL.Path, "/patients/123/")

		assert.Equal(actionNote, r.FormValue("actionnote"))
		assert.Equal(strconv.Itoa(apptID), r.FormValue("appointmentid"))
		assert.Equal(base64.StdEncoding.EncodeToString([]byte(attachmentContents)), r.FormValue("attachmentcontents"))
		assert.Equal(autoclose, r.FormValue("autoclose"))
		assert.Equal(strconv.Itoa(deptID), r.FormValue("departmentid"))
		assert.Equal(documentSubclass, r.FormValue("documentsubclass"))
		assert.Equal(internalNote, r.FormValue("internalnote"))
		assert.Equal(strconv.Itoa(providerID), r.FormValue("providerid"))

		b, _ := os.ReadFile("./resources/AddDocument.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &AddDocumentOptions{
		ActionNote:         &actionNote,
		AppointmentID:      &apptID,
		AttachmentContents: attachmentContents,
		AutoClose:          &autoclose,
		DepartmentID:       &deptID,
		DocumentSubclass:   documentSubclass,
		InternalNote:       &internalNote,
		ProviderID:         &providerID,
	}

	documentID, err := athenaClient.AddDocument(context.Background(), "123", opts)

	assert.Equal("100", documentID)
	assert.NoError(err)
}

func TestHTTPClient_AddPatientCaseDocument(t *testing.T) {
	assert := assert.New(t)

	autoClose := true
	callbackName := "callback name"
	callbackNumber := "callback number"
	callbackNumberType := "callback number type"
	deptID := 5
	documentSource := "source"
	documentSubclass := "subclass"
	internalNote := "note"
	outboundOnly := true
	priority := "priority"
	providerID := 9
	subject := "subject"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal("1", r.FormValue("autoclose"))
		assert.Equal("callback name", r.FormValue("callbackname"))
		assert.Equal("callback number", r.FormValue("callbacknumber"))
		assert.Equal("callback number type", r.FormValue("callbacknumbertype"))
		assert.Equal("5", r.FormValue("departmentid"))
		assert.Equal("source", r.FormValue("documentsource"))
		assert.Equal("subclass", r.FormValue("documentsubclass"))
		assert.Equal("note", r.FormValue("internalnote"))
		assert.Equal("1", r.FormValue("outboundonly"))
		assert.Equal("priority", r.FormValue("priority"))
		assert.Equal("9", r.FormValue("providerid"))
		assert.Equal("subject", r.FormValue("subject"))

		b, _ := os.ReadFile("./resources/AddPatientCaseDocument.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &AddPatientCaseDocumentOptions{
		AutoClose:          &autoClose,
		CallbackName:       &callbackName,
		CallbackNumber:     &callbackNumber,
		CallbackNumberType: &callbackNumberType,
		DepartmentID:       deptID,
		DocumentSource:     documentSource,
		DocumentSubclass:   documentSubclass,
		InternalNote:       &internalNote,
		OutboundOnly:       &outboundOnly,
		Priority:           &priority,
		ProviderID:         &providerID,
		Subject:            &subject,
	}

	patientCaseID, err := athenaClient.AddPatientCaseDocument(context.Background(), "123", opts)

	assert.Equal(491696, patientCaseID)
	assert.NoError(err)
}
