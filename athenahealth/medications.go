package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type MedicationEventType string

const (
	EventTypeStart MedicationEventType = "START"
	EventTypeEnd   MedicationEventType = "END"
	EventTypeOrder MedicationEventType = "ORDER"
	EventTypeEnter MedicationEventType = "ENTER"
	EventTypeFill  MedicationEventType = "FILL"
	EventTypeHide  MedicationEventType = "HIDE"
)

type MedicationEvent struct {
	EventDate       string              `json:"eventdate"`
	Type            MedicationEventType `json:"type"`
	UserDisplayName string              `json:"userdisplayname,omitempty"`
}

type MedicationStructuredSig struct {
	DosageAdditionalInstructions string       `json:"dosageadditionalinstructions,omitempty"`
	DosageFrequencyValue         NumberString `json:"dosagefrequencyvalue,omitempty"`
	DosageRoute                  string       `json:"dosageroute"`
	DosageAction                 string       `json:"dosageaction"`
	DosageFrequencyUnit          string       `json:"dosagefrequencyunit"`
	DosageQuantityUnit           string       `json:"dosagequantityunit"`
	DosageQuantityValue          NumberString `json:"dosagequantityvalue"`
	DosageFrequencyDescription   string       `json:"dosagefrequencydescription"`
	DosageDurationUnit           string       `json:"dosagedurationunit"`
	DosageDurationValue          NumberString `json:"dosagedurationvalue"`
}

type Medication struct {
	ClinicalOrderTypeID NumberString            `json:"clinicalordertypeid"`
	CreatedBy           string                  `json:"createdby"`
	Events              []MedicationEvent       `json:"events"`
	FDBMedicationID     NumberString            `json:"fdbmedicationid"`
	IsDiscontinued      bool                    `json:"isdiscontinued"`
	IsSafeToRenew       bool                    `json:"issafetorenew"`
	IsStructuredSig     bool                    `json:"isstructuredsig"`
	LastModifiedBy      string                  `json:"lastmodifiedby,omitempty"`
	LastModifiedDate    string                  `json:"lastmodifieddate,omitempty"`
	Medication          string                  `json:"medication"`
	MedicationEntryID   string                  `json:"medicationentryid"`
	MedicationID        NumberString            `json:"medicationid"`
	OrganClass          string                  `json:"organclass"`
	Pharmacy            string                  `json:"pharmacy"`
	PrescribedBy        string                  `json:"prescribedby"`
	Source              string                  `json:"source"`
	StructuredSig       MedicationStructuredSig `json:"structuredsig,omitempty"`
	UnstructuredSig     string                  `json:"unstructuredsig,omitempty"`
}

type ListMedicationsResult struct {
	LastUpdated                 string         `json:"lastupdated"`
	Medications                 [][]Medication `json:"medications"`
	NoMedicationsReported       bool           `json:"nomedicationsreported"`
	PatientDownloadConsent      bool           `json:"patientdownloadconsent"`
	PatientNeedsDownloadConsent bool           `json:"patientneedsdownloadconsent"`
}

const (
	MedicationTypeActive     string = "ACTIVE"
	MedicationTypeHistorical string = "HISTORICAL"
	MedicationTypeDenied     string = "DENIED"
)

type ListMedicationsOptions struct {
	DepartmentID   string
	MedicationType string
}

type SearchMedicationsResult struct {
	Medication   string `json:"medication"`
	MedicationID int    `json:"medicationid"`
}

// ListMedications - Retrieves a list of medications for a given patient
// GET /v1/{practiceid}/chart/{patientid}/medications
// https://docs.athenahealth.com/api/api-ref/medication#Get-patient's-medication-list
func (h *HTTPClient) ListMedications(ctx context.Context, patientID string, opts *ListMedicationsOptions) (*ListMedicationsResult, error) {
	out := &ListMedicationsResult{}

	q := url.Values{}
	if len(opts.DepartmentID) > 0 {
		q.Add("departmentid", opts.DepartmentID)
	}
	if len(opts.MedicationType) > 0 {
		if opts.MedicationType != MedicationTypeActive && opts.MedicationType != MedicationTypeHistorical && opts.MedicationType != MedicationTypeDenied {
			return nil, errors.New("invalid medication type parameter")
		}
		q.Add("medicationtype", opts.MedicationType)
	}

	_, err := h.Get(ctx, fmt.Sprintf("/chart/%s/medications", patientID), q, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// SearchMedications - Retrieves a list of medications for a given search parameters.
// GET /v1/{practiceid}/reference/medications
// https://docs.athenahealth.com/api/api-ref/medication#Search-for-available-medications
func (h *HTTPClient) SearchMedications(ctx context.Context, searchVal string) ([]*SearchMedicationsResult, error) {
	out := []*SearchMedicationsResult{}

	q := url.Values{}
	q.Add("searchvalue", searchVal)

	_, err := h.Get(ctx, "/reference/medications", q, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
