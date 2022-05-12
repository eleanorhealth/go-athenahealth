package athenahealth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// https://docs.athenahealth.com/api/workflows/health-history-forms-checkin
type HealthHistoryFormSection[FieldDataT, QuestionsT any] struct {
	FieldData           FieldDataT   `json:"fielddata"`
	PopulatedFrom       string       `json:"populatedfrom"`
	Name                string       `json:"name"`
	Populated           string       `json:"populated"`
	Questions           []QuestionsT `json:"questions"`
	PortalFormSectionID int          `json:"portalformsectionid"`
	AgeUnit             string       `json:"ageunit"`
	Empty               bool         `json:"empty"`
	Type                string       `json:"type"`
	Complete            bool         `json:"complete"`
}

type AllergySectionFieldData struct {
	Onset bool `json:"onset"`

	DeactivatedDate struct {
		InputType string `json:"inputtype"`
	} `json:"deactivateddate"`

	AllergyReactionID struct {
		InputType      string     `json:"inputtype"`
		DropdownValues [][]string `json:"dropdownvalues"`
		List           int        `json:"list"`
	} `json:"allergyreactionid"`

	Reaction bool `json:"reaction"`

	AllergyID struct {
		Required  int    `json:"required"`
		InputType string `json:"inputtype"`
	} `json:"allergyid"`

	OnsetDate struct {
		InputType string `json:"inputtype"`
	} `json:"onsetdate"`

	AllergyConceptType struct {
		InputType      string     `json:"inputtype"`
		DropdownValues [][]string `json:"dropdownvalues"`
	} `json:"allergyconcepttype"`

	SeverityID struct {
		InputType      string `json:"inputtype"`
		DropdownValues []struct {
			ParentCode  string `json:"parentcode"`
			Ordering    string `json:"ordering"`
			ID          string `json:"id"`
			SNOMEDCode  string `json:"snomedcode"`
			Description string `json:"description"`
		} `json:"dropdownvalues"`
		List int `json:"list"`
	} `json:"severityid"`

	AllergyCode struct {
		InputType string `json:"inputtype"`
	} `json:"allergycode"`
}

type AllergySectionQuestion struct {
	AllergyCode        string   `json:"allergycode"`
	AllergyConceptType string   `json:"allergyconcepttype"`
	AllergyID          string   `json:"allergyid"`
	AllergyReactionID  []string `json:"allergyreactionid"`
	Name               string   `json:"name"`
	SeverityID         []string `json:"severityid"`
}

type MedicationSectionFieldData struct {
	MedicationID struct {
		InputType string `json:"inputtype"`
	} `json:"medicationid"`

	StartDate struct {
		InputType string `json:"inputtype"`
	} `json:"startdate"`

	Frequency struct {
		InputType      string     `json:"inputtype"`
		DropdownValues [][]string `json:"dropdownvalues"`
	} `json:"frequency"`

	DosageQuantity struct {
		InputType string `json:"inputtype"`
	} `json:"dosagequantity"`

	StopDate struct {
		InputType string `json:"inputtype"`
	} `json:"stopdate"`

	RXNormName struct {
		InputType string `json:"inputtype"`
	} `json:"rxnormname"`

	RXNormID struct {
		InputType string `json:"inputtype"`
	} `json:"rxnormid"`

	MedDeactivateReasonID struct {
		InputType      string     `json:"inputtype"`
		DropdownValues [][]string `json:"dropdownvalues"`
	} `json:"meddeactivatereasonid"`

	DisplayDosageUnits struct {
		InputType string `json:"inputtype"`
	} `json:"displaydosageunits"`

	RXNormType struct {
		InputType      string     `json:"inputtype"`
		DropdownValues [][]string `json:"dropdownvalues"`
	} `json:"rxnormtype"`
}

type MedicationSectionQuestion struct {
	MedicationID       string `json:"medicationid"`
	Name               string `json:"name"`
	Frequency          string `json:"frequency"`
	DosageQuantity     string `json:"dosagequantity"`
	RXNormName         string `json:"rxnormname"`
	RXNormID           string `json:"rxnormid"`
	DisplayDosageUnits string `json:"displaydosageunits"`
	RXNormType         string `json:"rxnormtype"`
}

type SocialSectionFieldData struct {
	Default struct {
		InputType string `json:"inputtype"`
	} `json:"default"`
}

type SocialSectionQuestion struct {
	Fields struct {
		Default string `json:"default"`
	} `json:"fields"`
	PortalFormQuestionID   string     `json:"portalformquestionid"`
	Text                   string     `json:"text"`
	InputType              string     `json:"inputtype"`
	DropdownValues         [][]string `json:"dropdownvalues,omitempty"`
	ClinicalElementID      string     `json:"clinicalelementid,omitempty"`
	LocalClinicalElementID string     `json:"localclinicalelementid,omitempty"`
	Ordering               string     `json:"ordering,omitempty"`
}

type SurgicalSectionFieldData struct {
	Default struct {
		InputType string `json:"inputtype"`
	} `json:"default"`

	SurgeryDate struct {
		InputType string `json:"inputtype"`
		List      int    `json:"list"`
	} `json:"surgerydate"`
}

type SurgicalSectionQuestionFields struct {
	Default     string   `json:"default"`     // Y/N checkbox which says whether they've had this surgery before.
	SurgeryDate []string `json:"surgerydate"` // optional list of surgery dates as an approximate date: A date in either YYYY, MM/YYYY, or MM/DD/YYYY format. Even one surgery date should be submitted as a (singleton) list.
}

type SurgicalSectionQuestion struct {
	Fields                     SurgicalSectionQuestionFields `json:"fields"`
	PortalFormQuestionID       string                        `json:"portalformquestionid"`
	SurgicalHistoryProcedureID string                        `json:"surgicalhistoryprocedureid"`
	Text                       string                        `json:"text"`
}

type MedicalSectionFieldData struct {
	Default struct {
		InputType string `json:"inputtype"`
	} `json:"default"`
}

type MedicalSectionQuestion struct {
	PastMedicalHistoryQuestionID string `json:"pastmedicalhistoryquestionid"`
	Fields                       struct {
		Default string `json:"default"`
	} `json:"fields"`
	PortalFormQuestionID string `json:"portalformquestionid"`
	Text                 string `json:"text"`
}

type FamilySectionFieldData struct {
	Relation struct {
		Required       int        `json:"required"`
		InputType      string     `json:"inputtype"`
		DropdownValues [][]string `json:"dropdownvalues"`
	} `json:"relation"`
}

type FamilySectionQuestion struct {
	Fields []struct {
		Default       string `json:"default"`       // Y/N checkbox that says whether this person has this problem.
		Relation      string `json:"relation"`      // relation this person has with the patient (Mother, Sister, etc.).
		RelationKeyID string `json:"relationkeyid"` // ID of the particular relation. Each relation type has its own keyspace. This means you can have a Brother 1, Brother 2, Sister 1, etc.
	} `json:"fields"`
	PortalFormQuestionID string `json:"portalformquestionid"`
	Text                 string `json:"text"`
}

type HealthHistoryForm struct {
	Allergy    []*HealthHistoryFormSection[AllergySectionFieldData, AllergySectionQuestion]
	Medication []*HealthHistoryFormSection[MedicationSectionFieldData, MedicationSectionQuestion]
	Social     []*HealthHistoryFormSection[SocialSectionFieldData, SocialSectionQuestion]
	Surgical   []*HealthHistoryFormSection[SurgicalSectionFieldData, SurgicalSectionQuestion]
	Medical    []*HealthHistoryFormSection[MedicalSectionFieldData, MedicalSectionQuestion]
	Family     []*HealthHistoryFormSection[FamilySectionFieldData, FamilySectionQuestion]
}

func (h *HealthHistoryForm) FromSections(sections []json.RawMessage) error {
	type section struct {
		Type string `json:"type"`
	}

	for _, secBytes := range sections {
		sec := &section{}

		err := json.Unmarshal(secBytes, sec)
		if err != nil {
			return fmt.Errorf("unmarshaling health history section: %w", err)
		}

		switch sec.Type {
		case "ALLERGY":
			allergySec := &HealthHistoryFormSection[AllergySectionFieldData, AllergySectionQuestion]{}
			err = json.Unmarshal(secBytes, allergySec)
			if err != nil {
				return fmt.Errorf("unmarshaling allergy section: %w", err)
			}

			h.Allergy = append(h.Allergy, allergySec)

		case "MEDICATION":
			medicationSec := &HealthHistoryFormSection[MedicationSectionFieldData, MedicationSectionQuestion]{}
			err = json.Unmarshal(secBytes, medicationSec)
			if err != nil {
				return fmt.Errorf("unmarshaling medication section: %w", err)
			}

			h.Medication = append(h.Medication, medicationSec)

		case "SOCIAL":
			socialSec := &HealthHistoryFormSection[SocialSectionFieldData, SocialSectionQuestion]{}
			err = json.Unmarshal(secBytes, socialSec)
			if err != nil {
				return fmt.Errorf("unmarshaling social section: %w", err)
			}

			h.Social = append(h.Social, socialSec)

		case "SURGICAL":
			surgicalSec := &HealthHistoryFormSection[SurgicalSectionFieldData, SurgicalSectionQuestion]{}
			err = json.Unmarshal(secBytes, surgicalSec)
			if err != nil {
				return fmt.Errorf("unmarshaling surgical section: %w", err)
			}

			h.Surgical = append(h.Surgical, surgicalSec)

		case "MEDICAL":
			medicalSec := &HealthHistoryFormSection[MedicalSectionFieldData, MedicalSectionQuestion]{}
			err = json.Unmarshal(secBytes, medicalSec)
			if err != nil {
				return fmt.Errorf("unmarshaling medical section: %w", err)
			}

			h.Medical = append(h.Medical, medicalSec)

		case "FAMILY":
			familySec := &HealthHistoryFormSection[FamilySectionFieldData, FamilySectionQuestion]{}
			err = json.Unmarshal(secBytes, familySec)
			if err != nil {
				return fmt.Errorf("unmarshaling family section: %w", err)
			}

			h.Family = append(h.Family, familySec)
		}
	}

	return nil
}

// Get specific health history form for given appointment
// GET /v1/{practiceid}/appointments/{appointmentid}/healthhistoryforms/{formid}
// https://docs.athenahealth.com/api/api-ref/appointment-health-history-form#Get-specific-health-history-forms-for-given-appointment
func (h *HTTPClient) GetHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string) (*HealthHistoryForm, error) {
	var out []json.RawMessage

	_, err := h.Get(ctx, fmt.Sprintf("/appointments/%s/healthhistoryforms/%s", url.QueryEscape(appointmentID), url.QueryEscape(formID)), nil, &out)
	if err != nil {
		return nil, fmt.Errorf("fetching health history form for appointment: %w", err)
	}

	hhf := &HealthHistoryForm{}
	err = hhf.FromSections(out)
	if err != nil {
		return nil, fmt.Errorf("health history form from sections: %w", err)
	}

	return hhf, nil
}

type updateHealthHistoryFormResponse struct {
	ErrorMessage string `json:"errormessage"`
	Success      bool   `json:"success"`
}

// Update specific health history form for given appointment
// PUT /v1/{practiceid}/appointments/{appointmentid}/healthhistoryforms/{formid}
// https://docs.athenahealth.com/api/api-ref/appointment-health-history-form#Update-specific-health-history-forms-for-given-appointment
func (h *HTTPClient) UpdateHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string, form *HealthHistoryForm) error {
	if form == nil {
		return errors.New("form is nil")
	}

	var secs []any

	for _, allergySec := range form.Allergy {
		secs = append(secs, allergySec)
	}

	for _, medicationSec := range form.Medication {
		secs = append(secs, medicationSec)
	}

	for _, socialSec := range form.Social {
		secs = append(secs, socialSec)
	}

	for _, surgicalSec := range form.Surgical {
		secs = append(secs, surgicalSec)
	}

	for _, medicalSec := range form.Medical {
		secs = append(secs, medicalSec)
	}

	for _, familySec := range form.Family {
		secs = append(secs, familySec)
	}

	b, err := json.Marshal(secs)
	if err != nil {
		return fmt.Errorf("marshaling health history form sections: %w", err)
	}

	payload := make(url.Values)
	payload.Set("healthhistoryform", string(b))

	out := &updateHealthHistoryFormResponse{}

	_, err = h.PutForm(ctx, fmt.Sprintf("/appointments/%s/healthhistoryforms/%s", url.QueryEscape(appointmentID), url.QueryEscape(formID)), payload, out)
	if err != nil {
		return fmt.Errorf("updating health history form for appointment: %w", err)
	}

	if !out.Success {
		return fmt.Errorf("error returned from athena: %s", out.ErrorMessage)
	}

	return nil
}
