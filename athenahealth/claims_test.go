package athenahealth

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_CreateClaim(t *testing.T) {
	assert := assert.New(t)

	opts := &CreateClaimOptions{
		ClaimCharges: []*ClaimCharge{
			{
				AllowableAmount:     "1",
				AllowableMax:        "2",
				AllowableMin:        "1",
				AllowableScheduleID: 3,
				ICD10Code1:          4,
				ICD10Code2:          5,
				ICD10Code3:          6,
				ICD10Code4:          7,
				ICD9Code1:           8,
				ICD9Code2:           9,
				ICD9Code3:           10,
				ICD9Code4:           11,
				LineNote:            "foo",
				ProcedureCode:       "bar",
				UnitAmount:          "12",
				Units:               13,
			},
		},
		CustomFields: []*CustomFieldValue{
			{
				CustomFieldID:    "1",
				CustomFieldValue: "foo",
				OptionID:         "2",
			},
		},
		OrderingProviderID:          "1",
		PatientID:                   "2",
		PrimaryPatientInsuranceID:   "3",
		ReferralAuthID:              "4",
		ReferringProviderID:         "5",
		RenderingProviderID:         "6",
		Reserved19:                  "foo",
		SecondaryPatientInsuranceID: "7",
		ServiceDate:                 time.Now(),
		SupervisingProviderID:       "8",
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		// TODO: Assert HTTP request form.

		b, _ := ioutil.ReadFile("./resources/CreateClaim.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	claimIDs, err := athenaClient.CreateClaim(context.Background(), opts)
	assert.NoError(err)
	assert.Len(claimIDs, 2)
}
