package athenahealth

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

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
//
// GET /v1/{practiceid}/claims
//
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
