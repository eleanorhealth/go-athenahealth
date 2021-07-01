package athenahealth

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ListAdminDocuments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(r.URL.Path, "/patients/123/")
		assert.Equal("3", r.URL.Query().Get("departmentid"))

		b, _ := ioutil.ReadFile("./resources/ListAdminDocuments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListAdminDocumentsOptions{
		DepartmentID: "3",
	}

	res, err := athenaClient.ListAdminDocuments("123", opts)

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
	attachmentContents := "test attachment contents"
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

		b, _ := ioutil.ReadFile("./resources/AddDocument.json")
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

	documentID, err := athenaClient.AddDocument("123", opts)

	assert.Equal("100", documentID)
	assert.NoError(err)
}
