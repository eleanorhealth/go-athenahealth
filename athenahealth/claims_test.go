package athenahealth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_ClaimCreateNote(t *testing.T) {
	assert := assert.New(t)

	claimID := "2"
	opts := &ClaimCreateNoteOptions{ClaimNote: "Claim Note Here"}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		assert.Equal(r.Form.Get("claimnote"), opts.ClaimNote)

		b, _ := os.ReadFile("./resources/ClaimCreateNote.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.ClaimCreateNote(context.Background(), claimID, opts)
	assert.NoError(err)
}

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

func TestHTTPClient_ClaimUpdateFinancial(t *testing.T) {
	assert := assert.New(t)

	claimID := "40"

	opts := &ClaimUpdateFinancialOptions{
		ClaimCharges: []*ClaimCharge{
			{
				AllowableAmount:     func() *json.Number { a := json.Number("1"); return &a }(),
				AllowableMax:        func() *json.Number { a := json.Number("2"); return &a }(),
				AllowableMin:        func() *json.Number { a := json.Number("3"); return &a }(),
				AllowableScheduleID: func() *int { a := 4; return &a }(),
				ICD10Code1:          "ICD10Code1",
				ICD10Code2:          "ICD10Code2",
				ICD10Code3:          "ICD10Code3",
				ICD10Code4:          "ICD10Code4",
				ICD9Code1:           "ICD9Code1",
				ICD9Code2:           "ICD9Code2",
				ICD9Code3:           "ICD9Code3",
				ICD9Code4:           "ICD9Code4",
				LineNote:            "LineNote",
				ProcedureCode:       "ProcedureCode",
				UnitAmount:          func() *json.Number { a := json.Number("5"); return &a }(),
				Units:               6,
			},
		},
		CustomFields: []*CustomFieldValue{
			{
				CustomFieldID:    "CustomFieldID",
				CustomFieldValue: "CustomFieldValue",
				OptionID:         "OptionID",
			},
		},
		OrderingProviderID:  func() *int { a := 7; return &a }(),
		ReferralAuthID:      func() *int { a := 8; return &a }(),
		ReferringProviderID: func() *int { a := 9; return &a }(),
		ServiceTypeAddons:   []string{"1", "2", "3"},
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())

		claimChargesJSON, jsonErr := json.Marshal(opts.ClaimCharges)
		assert.Nil(jsonErr)

		customFieldsJSON, jsonErr := json.Marshal(opts.CustomFields)
		assert.Nil(jsonErr)

		assert.Equal(r.Form.Get("claimcharges"), string(claimChargesJSON))
		assert.Equal(r.Form.Get("customfields"), string(customFieldsJSON))
		assert.Equal(r.Form.Get("orderingproviderid"), strconv.Itoa(*opts.OrderingProviderID))
		assert.Equal(r.Form.Get("referralauthid"), strconv.Itoa(*opts.ReferralAuthID))
		assert.Equal(r.Form.Get("referringproviderid"), strconv.Itoa(*opts.ReferringProviderID))
		assert.Equal(r.Form.Get("servicetypeaddons"), strings.Join(opts.ServiceTypeAddons, ","))

		assert.Contains(r.URL.Path, fmt.Sprintf("/claims/%s", claimID))

		b, _ := os.ReadFile("./resources/ClaimUpdateFinancial.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	res, err := athenaClient.ClaimUpdateFinancial(context.Background(), claimID, opts)
	assert.Nil(err)

	assert.Equal(1, res.CustomFields)
	assert.Equal(true, res.Success)
	assert.Equal(2, res.Transactions)
}
