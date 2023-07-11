package athenahealth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type ClaimCreateNoteOptions struct {
	ClaimNote string `json:"claimnote"`
}

// ClaimCreateNote
// POST /v1/{practiceid}/claims/{claimid}/note
// https://docs.athenahealth.com/api/api-ref/claim-note#Create-new-claim-notes
func (h *HTTPClient) ClaimCreateNote(ctx context.Context, claimID string, opts *ClaimCreateNoteOptions) error {
	if claimID == "" {
		return fmt.Errorf("cannot ClaimCreateNote with empty claimID [%s]", claimID)
	}

	if opts == nil {
		return fmt.Errorf("cannot ClaimCreateNote with nil  ClaimCreateNoteOptions [%s]", opts)
	}

	form := url.Values{}

	form.Add("claimnote", opts.ClaimNote)

	out := ErrorMessageResponse{}
	_, err := h.PostForm(ctx, fmt.Sprintf("/claims/%s/note", claimID), form, &out)
	if err != nil {
		return err
	}

	if !out.Success {
		return fmt.Errorf("unexpected response with message: %s", out.Message)
	}

	return nil
}

type ClaimCharge struct {
	AllowableAmount     *json.Number `json:"allowableamount,omitempty"`
	AllowableMax        *json.Number `json:"allowablemax,omitempty"`
	AllowableMin        *json.Number `json:"allowablemin,omitempty"`
	AllowableScheduleID *int         `json:"allowablescheduleid,omitempty"`
	ICD10Code1          string       `json:"icd10code1"`
	ICD10Code2          string       `json:"icd10code2"`
	ICD10Code3          string       `json:"icd10code3"`
	ICD10Code4          string       `json:"icd10code4"`
	ICD9Code1           string       `json:"icd9code1"`
	ICD9Code2           string       `json:"icd9code2"`
	ICD9Code3           string       `json:"icd9code3"`
	ICD9Code4           string       `json:"icd9code4"`
	LineNote            string       `json:"linenote"`
	ProcedureCode       string       `json:"procedurecode"`
	UnitAmount          *json.Number `json:"unitamount,omitempty"`
	Units               int          `json:"units"`
}

type CreateClaimOptions struct {
	ClaimCharges                []*ClaimCharge
	CustomFields                []*CustomFieldValue
	DepartmentID                string
	OrderingProviderID          *string
	PatientID                   string
	PrimaryPatientInsuranceID   *string
	ReferralAuthID              *string
	ReferringProviderID         *string
	RenderingProviderID         *string
	Reserved19                  *string
	SecondaryPatientInsuranceID *string
	ServiceDate                 time.Time
	SupervisingProviderID       string
}

type createClaimResponse struct {
	ClaimIDs     []string `json:"claimids"`
	ErrorMessage string   `json:"errormessage"`
	Success      bool     `json:"success"`
}

func (h *HTTPClient) CreateFinancialClaim(ctx context.Context, opts *CreateClaimOptions) ([]string, error) {
	if opts == nil {
		panic("opts is nil")
	}

	form := url.Values{}

	claimChargesJSON, err := json.Marshal(opts.ClaimCharges)
	if err != nil {
		return []string{}, errors.Wrap(err, "marshaling claim charges")
	}

	customFieldsJSON, err := json.Marshal(opts.CustomFields)
	if err != nil {
		return []string{}, errors.Wrap(err, "marshaling custom fields")
	}

	form.Add("claimcharges", string(claimChargesJSON))
	form.Add("customfields", string(customFieldsJSON))
	form.Add("departmentid", opts.DepartmentID)

	if opts.OrderingProviderID != nil {
		form.Add("orderingproviderid", *opts.OrderingProviderID)
	}

	form.Add("patientid", opts.PatientID)

	if opts.PrimaryPatientInsuranceID != nil {
		form.Add("primarypatientinsuranceid", *opts.PrimaryPatientInsuranceID)
	}

	if opts.ReferralAuthID != nil {
		form.Add("referralauthid", *opts.ReferralAuthID)
	}

	if opts.ReferringProviderID != nil {
		form.Add("referringproviderid", *opts.ReferringProviderID)
	}

	if opts.RenderingProviderID != nil {
		form.Add("renderingproviderid", *opts.RenderingProviderID)
	}

	if opts.Reserved19 != nil {
		form.Add("reserved19", *opts.Reserved19)
	}

	if opts.SecondaryPatientInsuranceID != nil {
		form.Add("secondarypatientinsuranceid", *opts.SecondaryPatientInsuranceID)
	}

	if !opts.ServiceDate.IsZero() {
		form.Add("servicedate", opts.ServiceDate.Format("01/02/2006"))
	}

	form.Add("supervisingproviderid", opts.SupervisingProviderID)

	res := &createClaimResponse{}

	_, err = h.PostForm(ctx, "/claims", form, res)
	if err != nil {
		return []string{}, err
	}

	if !res.Success {
		return []string{}, errors.New(res.ErrorMessage)
	}

	return res.ClaimIDs, nil
}

type ClaimProcedure struct {
	ChargeAmount         string `json:"chargeamount"`
	ProcedureDescription string `json:"proceduredescription"`
	TransactionID        string `json:"transactionid"`
	ProcedureCode        string `json:"procedurecode"`
	ProcedureCategory    string `json:"procedurecategory"`
}

type ClaimDiagnosis struct {
	DiagnosisCategory    string `json:"diagnosiscategory"`
	DiagnosisID          string `json:"diagnosisid"`
	DiagnosisRawCode     string `json:"diagnosisrawcode"`
	DiagnosisCodeset     string `json:"diagnosiscodeset"`
	DiagnosisDescription string `json:"diagnosisdescription"`
	DeletedDiagnosis     string `json:"deleteddiagnosis"`
}

type Claim struct {
	Procedures        []ClaimProcedure   `json:"procedures"`
	ClaimCeatedDate   string             `json:"claimcreateddate"`
	BilledProviderID  int                `json:"billedproviderid"`
	ClaimID           string             `json:"claimid"`
	BilledServiceDate string             `json:"billedservicedate"`
	DepartmentID      int                `json:"departmentid"`
	Diagnoses         []ClaimDiagnosis   `json:"diagnoses"`
	PatientID         int                `json:"patientid"`
	CustomFields      []CustomFieldValue `json:"customfields"`
}

type ListClaimsOptions struct {
	PatientID        *string
	DepartmentID     *string
	ProviderID       *string
	ServiceStartDate *time.Time
	ServiceEndDate   *time.Time
	ShowCustomFields bool

	Pagination *PaginationOptions
}

type ListClaimsResult struct {
	Claims []*Claim

	Pagination *PaginationResult
}

type listClaimsResponse struct {
	Claims []*Claim `json:"claims"`

	PaginationResponse
}

// ListClaims - Get list of claims
// GET /v1/{practiceid}/claims
// https://docs.athenahealth.com/api/api-ref/claim#Get-list-of-claim-details
func (h *HTTPClient) ListClaims(ctx context.Context, opts *ListClaimsOptions) (*ListClaimsResult, error) {
	if opts == nil {
		panic("opts is nil")
	}

	out := &listClaimsResponse{}
	q := url.Values{}

	if opts.PatientID != nil {
		q.Add("patientid", *opts.PatientID)
	}

	if opts.DepartmentID != nil {
		q.Add("departmentid", *opts.DepartmentID)
	}

	if opts.ProviderID != nil {
		q.Add("providerid", *opts.ProviderID)
	}

	if opts.ServiceStartDate != nil {
		q.Add("servicestartdate", opts.ServiceStartDate.Format("01/02/2006"))
	}

	if opts.ServiceEndDate != nil {
		q.Add("serviceenddate", opts.ServiceEndDate.Format("01/02/2006"))
	}

	if opts.ShowCustomFields {
		q.Add("showcustomfields", "true")
	}

	if opts.Pagination != nil {
		if opts.Pagination.Limit > 0 {
			q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
		}

		if opts.Pagination.Offset > 0 {
			q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
		}
	}

	_, err := h.Get(ctx, "/claims", q, out)
	if err != nil {
		return nil, err
	}

	return &ListClaimsResult{
		Claims:     out.Claims,
		Pagination: makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type ClaimUpdateFinancialOptions struct {
	// List of charges for this claim whose allowable values should be updated. This should be a JSON string representing an array of charge objects.
	ClaimCharges []*ClaimCharge `json:"claimcharges"`
	// A list of custom field JSON objects to populate on creation of a claim.
	CustomFields []*CustomFieldValue `json:"customfields"`
	// The ordering provider ID. 'Ordering Provider' service type add-on must be enabled. Default is no ordering provider ID. Any entry in this field will override any ordering provider ID in the service type add-ons field.
	OrderingProviderID *int `json:"orderingproviderid"`
	// The referral authorization ID to associate with this claim.
	ReferralAuthID *int `json:"referralauthid"`
	// The referring provider ID (not the same from /providers) associated with this claim.
	ReferringProviderID *int `json:"referringproviderid"`
	// Array of service type add-ons (STAOs) for the claim. Some claim level STAO fields do not support multiple values. These fields will save only the first value if more than one is passed in. The functionality behind this parameter is toggled by COLDEN_CLAIM_STAO_MDP_API. It is part of a feature that is scheduled to rollout in or before March 2023.
	ServiceTypeAddons []string `json:"servicetypeaddons"`
}

type ClaimUpdateFinancialResult struct {
	// Number of custom fields updated.
	CustomFields int `json:"customfields"`
	// Whether the operation was successful.
	Success bool `json:"success"`
	// Number of transactions with allowables information updated. Note that this does not include transactions that had Service type add-on information updated.
	Transactions int `json:"transactions"`
}

// ClaimUpdateFinancial
// PUT /v1/{practiceid}/claims/{claimid}
// https://docs.athenahealth.com/api/api-ref/claim#Update-individual-claim-details
func (h *HTTPClient) ClaimUpdateFinancial(ctx context.Context, claimID string, opts *ClaimUpdateFinancialOptions) (*ClaimUpdateFinancialResult, error) {
	if claimID == "" {
		return nil, fmt.Errorf("cannot ClaimUpdateFinancial with empty claimID [%s]", claimID)
	}

	form := url.Values{}

	if len(opts.ClaimCharges) > 0 {
		claimChargesJSON, jsonErr := json.Marshal(opts.ClaimCharges)
		if jsonErr != nil {
			return nil, errors.Wrap(jsonErr, "marshaling claim charges")
		}
		form.Add("claimcharges", string(claimChargesJSON))
	}

	if len(opts.CustomFields) > 0 {
		customFieldsJSON, jsonErr := json.Marshal(opts.CustomFields)
		if jsonErr != nil {
			return nil, errors.Wrap(jsonErr, "marshaling custom fields")
		}
		form.Add("customfields", string(customFieldsJSON))
	}

	if opts.OrderingProviderID != nil {
		form.Add("orderingproviderid", strconv.Itoa(*opts.OrderingProviderID))
	}

	if opts.ReferralAuthID != nil {
		form.Add("referralauthid", strconv.Itoa(*opts.ReferralAuthID))
	}

	if opts.ReferringProviderID != nil {
		form.Add("referringproviderid", strconv.Itoa(*opts.ReferringProviderID))
	}

	if len(opts.ServiceTypeAddons) > 0 {
		form.Add("servicetypeaddons", strings.Join(opts.ServiceTypeAddons, ","))
	}

	out := ClaimUpdateFinancialResult{}
	_, err := h.PostForm(ctx, fmt.Sprintf("/claims/%s", claimID), form, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
