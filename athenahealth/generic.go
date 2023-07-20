package athenahealth

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ErrorMessageResponse struct {
	Message string `json:"errormessage"`
	Success bool   `json:"success"`
}

type StrInt struct {
	IntVal int
	StrVal string
}

func (is *StrInt) IsError() bool {
	if is.IntVal == 0 && is.StrVal == "" {
		return true
	}

	if is.StrVal != "" && is.IntVal == 0 {
		return true
	}

	return false
}

func (is *StrInt) Error() string {
	if is.IntVal == 0 && is.StrVal == "" {
		return fmt.Sprintf("StrInt is neither an int [%d] or a string [%s]", is.IntVal, is.StrVal)
	}

	if is.StrVal != "" && is.IntVal == 0 {
		return fmt.Sprintf("StrInt holds an Athena API Error message: [%s]", is.StrVal)
	}

	return ""
}

// StatusResponse
// https://docs.athenahealth.com/api/api-ref/appointment-booked#Appointment-Booked
// * This subroutine will return 1 on success, and will otherwise return an error message.
type StatusResponse struct {
	Status StrInt `json:"status"`
}

func (is *StrInt) UnmarshalJSON(data []byte) error {
	var tmp json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	strErr := json.Unmarshal(tmp, &is.StrVal)
	intErr := json.Unmarshal(tmp, &is.IntVal)

	if strErr != nil && intErr != nil {
		return fmt.Errorf("failed to unmarshal JSON string: [%+v] [%+v]", strErr, intErr)
	}

	if strErr == nil {
		if intValue, err := strconv.Atoi(is.StrVal); err == nil {
			is.IntVal = intValue
		}
	}

	if intErr == nil {
		is.StrVal = strconv.Itoa(is.IntVal)
	}

	return nil
}
