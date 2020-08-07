package athenahealth

type CustomFieldValue struct {
	CustomFieldID    string `json:"customfieldid"`
	CustomFieldValue string `json:"customfieldvalue"`
	OptionID         int    `json:"optionid"`
}

type CustomField struct {
	CaseSensitive  bool   `json:"casesensitive"`
	CustomFieldID  string `json:"customfieldid"`
	DisallowUpdate bool   `json:"disallowupdate"`
	Name           string `json:"name"`
	Searchable     bool   `json:"searchable"`
	Select         bool   `json:"select"`
	Type           string `json:"type"`
}

// ListCustomFields - Gets a set of patients or creates a patient.
// GET /v1/{practiceid}/customfields
// https://developer.athenahealth.com/docs/read/administrative/Custom_Fields_List#section-0
func (h *HTTPClient) ListCustomFields() ([]*CustomField, error) {
	var out []*CustomField

	_, err := h.Get("/customfields", nil, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
