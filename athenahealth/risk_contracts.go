package athenahealth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// RiskContract represents a risk contract associated with a patient
type RiskContract struct {
	ContractName  string `json:"contractname"`
	EffectiveDate string `json:"effectivedate"`
	ExpirationDate string `json:"expirationdate"`
	RiskContractID int    `json:"riskcontractid"`
}

// ListRiskContractsOptions represents options for listing risk contracts
type ListRiskContractsOptions struct {
	DepartmentID string
}

// CreateRiskContractOptions represents options for creating/updating a risk contract
type CreateRiskContractOptions struct {
	RiskContractID int
	EffectiveDate  string // Format: MM/DD/YYYY
	ExpirationDate string // Format: MM/DD/YYYY (optional)
}

// ListRiskContracts - Get a list of risk contracts associated with the patient
//
// GET /v1/{practiceid}/chart/{patientid}/riskcontract
//
// https://docs.athenahealth.com/api/api-ref/patient-risk-contract
func (h *HTTPClient) ListRiskContracts(ctx context.Context, patientID string, opts *ListRiskContractsOptions) ([]*RiskContract, error) {
	out := []*RiskContract{}

	q := url.Values{}

	if opts != nil && opts.DepartmentID != "" {
		q.Add("departmentid", opts.DepartmentID)
	}

	_, err := h.Get(ctx, fmt.Sprintf("/chart/%s/riskcontract", patientID), q, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// CreateRiskContract - Create a new risk contract for the patient
//
// PUT /v1/{practiceid}/chart/{patientid}/riskcontract
//
// https://docs.athenahealth.com/api/api-ref/patient-risk-contract
func (h *HTTPClient) CreateRiskContract(ctx context.Context, patientID string, opts *CreateRiskContractOptions) error {
	if opts == nil {
		panic("opts is nil")
	}

	out := &MessageResponse{}

	form := url.Values{}
	form.Add("riskcontractid", strconv.Itoa(opts.RiskContractID))
	form.Add("effectivedate", opts.EffectiveDate)

	if opts.ExpirationDate != "" {
		form.Add("expirationdate", opts.ExpirationDate)
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/chart/%s/riskcontract", patientID), form, &out)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRiskContract - Delete a risk contract for the patient
//
// DELETE /v1/{practiceid}/chart/{patientid}/riskcontract/{riskcontractid}
//
// https://docs.athenahealth.com/api/api-ref/patient-risk-contract
func (h *HTTPClient) DeleteRiskContract(ctx context.Context, patientID string, riskContractID int) error {
	out := &MessageResponse{}

	_, err := h.Delete(ctx, fmt.Sprintf("/chart/%s/riskcontract/%d", patientID, riskContractID), nil, &out)
	if err != nil {
		return err
	}

	return nil
}
