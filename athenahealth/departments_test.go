package athenahealth

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_DepartmentGetRequiredCheckInFields(t *testing.T) {
	assert := assert.New(t)

	departmentID := "45"

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.URL.Path, "/departments/45/checkinrequired")
		b, _ := os.ReadFile("./resources/DepartmentGetRequiredCheckInFields.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.DepartmentGetRequiredCheckInFields(context.Background(), departmentID)

	assert.NoError(err)
	assert.True(reflect.DeepEqual(res.FieldList, []string{"1", "2", "3", "4", "5"}))
}

func TestHTTPClient_GetDepartment(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/GetDepartment.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	department, err := athenaClient.GetDepartment(context.Background(), "1")

	assert.NotNil(department)
	assert.NoError(err)
}

func TestHTTPClient_ListDepartments(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("1", r.URL.Query().Get("hospitalonly"))
		assert.Equal("1", r.URL.Query().Get("providerlist"))
		assert.Equal("1", r.URL.Query().Get("showalldepartments"))

		b, _ := os.ReadFile("./resources/ListDepartments.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListDepartmentsOptions{
		HospitalOnly:       true,
		ProviderList:       true,
		ShowAllDepartments: true,
	}

	res, err := athenaClient.ListDepartments(context.Background(), opts)

	assert.Len(res.Departments, 1)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 1)
	assert.NoError(err)
}
