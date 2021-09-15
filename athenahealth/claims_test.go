package athenahealth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_CreateClaim(t *testing.T) {
	assert := assert.New(t)

	allowableAmount := json.Number("1")
	allowableMax := json.Number("2")
	allowableMin := json.Number("3")
	allowableScheduleID := 4
	primaryPatientInsuranceID := "4"
	secondaryPatientInsuranceID := "8"
	unitAmount := json.Number("12")
	orderingProviderID := "2"
	referralAuthID := "5"
	referringProviderID := "6"
	renderingProviderID := "7"
	reserved19 := "foo"

	opts := &CreateClaimOptions{
		ClaimCharges: []*ClaimCharge{
			{
				AllowableAmount:     &allowableAmount,
				AllowableMax:        &allowableMax,
				AllowableMin:        &allowableMin,
				AllowableScheduleID: &allowableScheduleID,
				ICD10Code1:          "100",
				ICD10Code2:          "200",
				ICD10Code3:          "300",
				ICD10Code4:          "400",
				ICD9Code1:           "500",
				ICD9Code2:           "600",
				ICD9Code3:           "700",
				ICD9Code4:           "800",
				LineNote:            "foo",
				ProcedureCode:       "bar",
				UnitAmount:          &unitAmount,
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
		DepartmentID:                "1",
		OrderingProviderID:          &orderingProviderID,
		PatientID:                   "3",
		PrimaryPatientInsuranceID:   &primaryPatientInsuranceID,
		ReferralAuthID:              &referralAuthID,
		ReferringProviderID:         &referringProviderID,
		RenderingProviderID:         &renderingProviderID,
		Reserved19:                  &reserved19,
		SecondaryPatientInsuranceID: &secondaryPatientInsuranceID,
		ServiceDate:                 time.Now(),
		SupervisingProviderID:       "9",
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		claimCharges := []*ClaimCharge{}
		customFields := []*CustomFieldValue{}

		err := json.Unmarshal([]byte(r.Form.Get("claimcharges")), &claimCharges)
		assert.NoError(err)

		err = json.Unmarshal([]byte(r.Form.Get("customfields")), &customFields)
		assert.NoError(err)

		assert.Equal(opts.ClaimCharges, claimCharges)
		assert.Equal(opts.CustomFields, customFields)

		assert.Equal(r.Form.Get("departmentid"), opts.DepartmentID)
		assert.Equal(r.Form.Get("orderingproviderid"), *opts.OrderingProviderID)
		assert.Equal(r.Form.Get("patientid"), opts.PatientID)
		assert.Equal(r.Form.Get("primarypatientinsuranceid"), *opts.PrimaryPatientInsuranceID)
		assert.Equal(r.Form.Get("referralauthid"), *opts.ReferralAuthID)
		assert.Equal(r.Form.Get("referringproviderid"), *opts.ReferringProviderID)
		assert.Equal(r.Form.Get("renderingproviderid"), *opts.RenderingProviderID)
		assert.Equal(r.Form.Get("reserved19"), *opts.Reserved19)
		assert.Equal(r.Form.Get("secondarypatientinsuranceid"), *opts.SecondaryPatientInsuranceID)
		assert.Equal(r.Form.Get("servicedate"), opts.ServiceDate.Format("01/02/2006"))
		assert.Equal(r.Form.Get("supervisingproviderid"), opts.SupervisingProviderID)

		b, _ := ioutil.ReadFile("./resources/CreateClaim.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	claimIDs, err := athenaClient.CreateFinancialClaim(context.Background(), opts)
	assert.NoError(err)
	assert.Len(claimIDs, 2)
}
