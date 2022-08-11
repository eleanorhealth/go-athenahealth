package athenahealth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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

		b, _ := os.ReadFile("./resources/CreateClaim.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	claimIDs, err := athenaClient.CreateFinancialClaim(context.Background(), opts)
	assert.NoError(err)
	assert.Len(claimIDs, 2)
}

func TestHTTPClient_ListClaims(t *testing.T) {
	assert := assert.New(t)

	patientID := "1"
	departmentID := "2"
	providerID := "3"
	startDate := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2021, 9, 5, 0, 0, 0, 0, time.UTC)

	opts := &ListClaimsOptions{
		PatientID:        &patientID,
		DepartmentID:     &departmentID,
		ProviderID:       &providerID,
		ServiceStartDate: &startDate,
		ServiceEndDate:   &endDate,
		ShowCustomFields: true,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(r.URL.Path, "/claims")

		assert.Equal(r.URL.Query().Get("patientid"), patientID)
		assert.Equal(r.URL.Query().Get("departmentid"), departmentID)
		assert.Equal(r.URL.Query().Get("providerid"), providerID)
		assert.Equal(r.URL.Query().Get("servicestartdate"), startDate.Format("01/02/2006"))
		assert.Equal(r.URL.Query().Get("serviceenddate"), endDate.Format("01/02/2006"))
		assert.Equal(r.URL.Query().Get("showcustomfields"), "true")

		b, _ := os.ReadFile("./resources/ListClaims.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ListClaims(context.Background(), opts)

	assert.Len(res.Claims, 1)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 1)
	assert.NoError(err)
}
