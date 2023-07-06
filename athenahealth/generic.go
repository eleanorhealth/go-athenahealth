package athenahealth

import (
	"encoding/json"
	"strconv"
)

type BoolString bool

func (bs *BoolString) UnmarshalJSON(data []byte) error {
	var stringVal string
	if err := json.Unmarshal(data, &stringVal); err != nil {
		return err
	}

	boolVal, parseErr := strconv.ParseBool(stringVal)
	if parseErr != nil {
		return parseErr
	}

	*bs = BoolString(boolVal)
	return nil
}

type MessageResponse struct {
	Message string     `json:"message"`
	Success BoolString `json:"success"`
}

type ErrorMessageResponse struct {
	Message string     `json:"errormessage"`
	Success BoolString `json:"success"`
}
