package athenahealth

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_CreateAppointmentType(t *testing.T) {
	assert := assert.New(t)

	opts := &CreateAppointmentTypeOptions{
		Duration:         "99",
		Generic:          func() *bool { a := true; return &a }(),
		Name:             "1001",
		Patient:          true,
		ShortName:        "11",
		TemplateTypeOnly: func() *bool { a := false; return &a }(),
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("duration"), opts.Duration)
		assert.Equal(r.Form.Get("generic"), strconv.FormatBool(*opts.Generic))
		assert.Equal(r.Form.Get("name"), opts.Name)
		assert.Equal(r.Form.Get("patient"), strconv.FormatBool(opts.Patient))
		assert.Equal(r.Form.Get("shortname"), opts.ShortName)
		assert.Equal(r.Form.Get("templatetypeonly"), strconv.FormatBool(*opts.TemplateTypeOnly))
		assert.Equal(r.URL.Path, "/appointmenttypes")
		b, _ := os.ReadFile("./resources/CreateAppointmentType.json")
		_, _ = w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	createAppointmentTypeResult, err := athenaClient.CreateAppointmentType(context.Background(), opts)

	assert.NotNil(createAppointmentTypeResult)
	assert.NoError(err)
	assert.Equal(5, createAppointmentTypeResult.AppointmentTypeID)
}

func TestHTTPClient_GetAppointmentType(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/appointmenttypes/578", r.URL.Path)
		assert.Equal("true", r.URL.Query().Get("showappointmenttypeclasses"))
		b, _ := os.ReadFile("./resources/GetAppointmentType.json")
		_, _ = w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.GetAppointmentType(context.Background(), "578", &GetAppointmentTypeOptions{ShowAppointmentTypeClasses: true})

	assert.NoError(err)
	assert.Equal(578, res.AppointmentTypeID)
	assert.Equal(40, res.Duration)
	assert.Equal("CRP Intake 40 - Virtual", res.Name)
	assert.True(res.CreateEncounterOnCheckIn)
	assert.Equal(40, res.DepartmentID)
	assert.Equal(302, res.ProviderID)
	assert.Equal("Intake", res.AppointmentTypeClassName)
	assert.False(res.TemplateTypeOnly)
}

func TestHTTPClient_GetAppointmentType_Empty(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("[]"))
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	_, err := athenaClient.GetAppointmentType(context.Background(), "1", nil)

	assert.Error(t, err)
}

func TestHTTPClient_ListAppointmentTypes(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/appointmenttypes", r.URL.Path)
		query := r.URL.Query()
		assert.Equal("true", query.Get("hidegeneric"))
		assert.Equal("true", query.Get("hidenongeneric"))
		assert.Equal("true", query.Get("hidetemplatetypeonly"))
		assert.Equal("true", query.Get("showappointmenttypeclasses"))
		assert.Equal("false", query.Get("hidenonpatient"))
		assert.Equal("1,2", query.Get("providerids"))
		assert.Equal("40,41", query.Get("departmentids"))
		assert.Equal("2", r.URL.Query().Get("limit"))
		assert.Equal("4", r.URL.Query().Get("offset"))
		b, _ := os.ReadFile("./resources/ListAppointmentTypes.json")
		_, _ = w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListAppointmentTypes(context.Background(), &ListAppointmentTypesOptions{
		HideGeneric:                true,
		HideNonGeneric:             true,
		HideTemplateTypeOnly:       true,
		ShowAppointmentTypeClasses: true,
		HideNonPatient:             func() *bool { b := false; return &b }(),
		ProviderIDs:                []string{"1", "2"},
		DepartmentIDs:              []string{"40", "41"},
		Pagination:                 &PaginationOptions{Limit: 2, Offset: 4},
	})

	assert.NoError(err)
	assert.Len(res.AppointmentTypes, 2)
	assert.Equal(20, res.AppointmentTypes[0].Duration)
	assert.True(res.AppointmentTypes[0].TemplateTypeOnly)
	assert.Equal(60, res.AppointmentTypes[1].Duration)
	assert.Equal(3, res.Pagination.TotalCount)
	assert.Equal(2, res.Pagination.NextOffset)
}
