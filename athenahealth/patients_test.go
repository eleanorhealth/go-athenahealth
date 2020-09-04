package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetPatient(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("showcustomfields"))

		b, _ := ioutil.ReadFile("./resources/GetPatient.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	id := "1"
	opts := &GetPatientOptions{
		ShowCustomFields: true,
	}
	patient, err := athenaClient.GetPatient(id, opts)

	assert.NotNil(patient)
	assert.Nil(err)
}

func TestHTTPClient_ListPatients(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("John", r.URL.Query().Get("firstname"))
		assert.Equal("Smith", r.URL.Query().Get("lastname"))

		b, _ := ioutil.ReadFile("./resources/ListPatients.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListPatientsOptions{
		FirstName: "John",
		LastName:  "Smith",
	}

	patients, err := athenaClient.ListPatients(opts)

	assert.Len(patients, 2)
	assert.Nil(err)
}

func TestBalance_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *Balance
		args    args
		wantErr bool
	}{
		{
			name:    "string value",
			b:       new(Balance),
			args:    args{data: []byte(`"55.01"`)},
			wantErr: false,
		},
		{
			name:    "negative string value",
			b:       new(Balance),
			args:    args{data: []byte(`"-55.01"`)},
			wantErr: false,
		},
		{
			name:    "int value",
			b:       new(Balance),
			args:    args{data: []byte(`55`)},
			wantErr: false,
		},
		{
			name:    "float64 value",
			b:       new(Balance),
			args:    args{data: []byte(`55.01`)},
			wantErr: false,
		},
		{
			name:    "invalid bool value",
			b:       new(Balance),
			args:    args{data: []byte(`false`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Balance.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
