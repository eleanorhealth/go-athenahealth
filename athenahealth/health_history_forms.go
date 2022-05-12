package athenahealth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

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

type SurgicalSectionQuestion struct {
	Fields struct {
		Default string `json:"default"`
	} `json:"fields"`

	PortalFormQuestionID       string `json:"portalformquestionid"`
	Text                       string `json:"text"`
	SurgicalHistoryProcedureID string `json:"surgicalhistoryprocedureid"`
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
	}
}

type FamilySectionQuestion struct {
	Fields                       []interface{} `json:"fields"`
	PortalFormQuestionID         string        `json:"portalformquestionid"`
	Text                         string        `json:"text"`
	SNOMEDFamilyHistoryProblemID string        `json:"snomedfamilyhistoryproblemid"`
	SNOMEDCode                   string        `json:"snomedcode"`
	OrigText                     string        `json:"origtext"`
}

type HealthHistoryForm struct {
	Allergy    []*HealthHistoryFormSection[AllergySectionFieldData, AllergySectionQuestion]
	Medication []*HealthHistoryFormSection[MedicationSectionFieldData, MedicationSectionQuestion]
	Social     []*HealthHistoryFormSection[SocialSectionFieldData, SocialSectionQuestion]
	Surgical   []*HealthHistoryFormSection[SurgicalSectionFieldData, SurgicalSectionQuestion]
	Medical    []*HealthHistoryFormSection[MedicalSectionFieldData, MedicalSectionQuestion]
	Family     []*HealthHistoryFormSection[FamilySectionFieldData, FamilySectionQuestion]
}

// Get specific health history form for given appointment
// GET /v1/{practiceid}/appointments/{appointmentid}/healthhistoryforms/{formid}
// https://docs.athenahealth.com/api/api-ref/appointment-health-history-form#Get-specific-health-history-forms-for-given-appointment
func (h *HTTPClient) GetHealthHistoryFormForAppointment(ctx context.Context, appointmentID, formID string) (*HealthHistoryForm, error) {
	var out []json.RawMessage

	_, err := h.Get(ctx, fmt.Sprintf("/appointments/%s/healthhistoryforms/%s", url.QueryEscape(appointmentID), url.QueryEscape(formID)), nil, &out)
	if err != nil {
		return nil, fmt.Errorf("fetching healthy history form for appointment: %w", err)
	}

	type section struct {
		Type string `json:"type"`
	}

	hhf := &HealthHistoryForm{}

	for _, secBytes := range out {
		sec := &section{}

		err = json.Unmarshal(secBytes, sec)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling health history section: %w", err)
		}

		switch sec.Type {
		case "ALLERGY":
			allergySec := &HealthHistoryFormSection[AllergySectionFieldData, AllergySectionQuestion]{}
			err = json.Unmarshal(secBytes, allergySec)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling allergy section: %w", err)
			}

			hhf.Allergy = append(hhf.Allergy, allergySec)

		case "MEDICATION":
			medicationSec := &HealthHistoryFormSection[MedicationSectionFieldData, MedicationSectionQuestion]{}
			err = json.Unmarshal(secBytes, medicationSec)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling medication section: %w", err)
			}

			hhf.Medication = append(hhf.Medication, medicationSec)

		case "SOCIAL":
			socialSec := &HealthHistoryFormSection[SocialSectionFieldData, SocialSectionQuestion]{}
			err = json.Unmarshal(secBytes, socialSec)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling social section: %w", err)
			}

			hhf.Social = append(hhf.Social, socialSec)

		case "SURGICAL":
			surgicalSec := &HealthHistoryFormSection[SurgicalSectionFieldData, SurgicalSectionQuestion]{}
			err = json.Unmarshal(secBytes, surgicalSec)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling surgical section: %w", err)
			}

			hhf.Surgical = append(hhf.Surgical, surgicalSec)

		case "MEDICAL":
			medicalSec := &HealthHistoryFormSection[MedicalSectionFieldData, MedicalSectionQuestion]{}
			err = json.Unmarshal(secBytes, medicalSec)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling medical section: %w", err)
			}

			hhf.Medical = append(hhf.Medical, medicalSec)

		case "FAMILY":
			familySec := &HealthHistoryFormSection[FamilySectionFieldData, FamilySectionQuestion]{}
			err = json.Unmarshal(secBytes, familySec)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling family section: %w", err)
			}

			hhf.Family = append(hhf.Family, familySec)
		}
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
		return fmt.Errorf("updating healthy history form for appointment: %w", err)
	}

	if !out.Success {
		return fmt.Errorf("error returned from athena: %s", out.ErrorMessage)
	}

	return nil
}