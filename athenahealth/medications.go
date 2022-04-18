package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type MedicationEvent struct {
	Eventdate       string `json:"eventdate"`
	Type            string `json:"type"`
	Userdisplayname string `json:"userdisplayname"`
}

type MedicationStructuredSig struct {
	DosageAction                 string `json:"dosageaction"`
	DosageAdditionalInstructions string `json:"dosageadditionalinstructions"`
	DosageDurationUnit           string `json:"dosagedurationunit"`
	DosageDurationValue          int    `json:"dosagedurationvalue"`
	DosageFrequencyDescription   string `json:"dosagefrequencydescription"`
	DosageFrequencyUnit          string `json:"dosagefrequencyunit"`
	DosageFrequencyValue         int    `json:"dosagefrequencyvalue"`
	DosageQuantityUnit           string `json:"dosagequantityunit"`
	DosageQuantityValue          int    `json:"dosagequantityvalue"`
	DosageRoute                  string `json:"dosageroute"`
}

type Medication struct {
	ApprovedBy        string                  `json:"approvedby"`
	BillingNDC        string                  `json:"billingndc"`
	CreatedBy         string                  `json:"createdby"`
	EarliestFillDate  string                  `json:"earliestfilldate"`
	EncounterID       int                     `json:"encounterid"`
	Events            []MedicationEvent       `json:"events"`
	FutureSubmitDate  string                  `json:"futuresubmitdate"`
	IsSafeToRenew     bool                    `json:"issafetorenew"`
	IsDiscontinued    bool                    `json:"isdiscontinued"`
	IsStructuredSig   bool                    `json:"isstructuredsig"`
	Medication        string                  `json:"medication"`
	MedicationEntryID string                  `json:"medicationentryid"`
	MedicationID      int                     `json:"medicationid"`
	NDCOptions        string                  `json:"ndcoptions"`
	OrderingMode      string                  `json:"orderingmode"`
	OrganClass        string                  `json:"organclass"`
	PatientNote       string                  `json:"patientnote"`
	Pharmacy          string                  `json:"pharmacy"`
	PharmacyNCPDPID   string                  `json:"pharmacyncpdpid"`
	PrescribedBy      string                  `json:"prescribedBy"`
	ProviderNote      string                  `json:"providernote"`
	Quantity          int                     `json:"quantity"`
	QuantityUnit      string                  `json:"quantityunit"`
	RefillsAllowed    int                     `json:"refillsallowed"`
	Route             string                  `json:"route"`
	RXNorm            string                  `json:"rxnorm"`
	Source            string                  `json:"source"`
	Status            string                  `json:"status"`
	StopReason        string                  `json:"stopreason"`
	StructuredSig     MedicationStructuredSig `json:"structuredsig"`
	TherapeuticClass  string                  `json:"therapeuticclass"`
	UnstructuredSig   string                  `json:"unstructuredsig"`
}

type ListMedicationsResponse struct {
	LastDownloadDenialReason    string          `json:"lastdownloaddenialreason"`
	LastDownloadDenied          string          `json:"lastdownloaddenied"`
	LastDownloadedDate          string          `json:"lastdownloadeddate"`
	LastUpdated                 string          `json:"lastupdated"`
	Medications                 [][]*Medication `json:"medications"`
	NoMedicationsReported       bool            `json:"nomedicationsreported"`
	PatientDownloadConsent      bool            `json:"patientdownloadconsent"`
	PatientNeedsDownloadConsent bool            `json:"patientneedsdownloadconsent"`
	SectionNote                 string          `json:"sectionnote"`
}
type ListMedicationsResult ListMedicationsResponse

func (h *HTTPClient) ListMedications(ctx context.Context, patientID, departmentID string) (*ListMedicationsResult, error) {
	out := &ListMedicationsResponse{}

	q := url.Values{
		"departmentid": []string{departmentID},
	}

	_, err := h.Get(ctx, fmt.Sprintf("chart/%s/medications", patientID), q, out)
	if err != nil {
		return nil, err
	}

	result := (*ListMedicationsResult)(out)

	return result, nil
}
