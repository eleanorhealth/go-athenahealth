package athenahealth

import (
	"encoding/json"
	"errors"
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

func (sr *StatusResponse) GetError() error {
	if !sr.IsValid {
		return fmt.Errorf("StatusResponse is neither an int [%d] or a string [%s]", sr.IntVal, sr.StrVal)
	}

	if sr.IsError {
		return errors.New(sr.StrVal)
	}

	return nil
}

// StatusResponse
// https://docs.athenahealth.com/api/api-ref/appointment-booked#Appointment-Booked
// * This subroutine will return 1 on success, and will otherwise return an error message.
type StatusResponse struct {
	StrVal  string
	IntVal  int
	IsValid bool
	IsError bool
}

type statusResponseStr struct {
	Status string `json:"status"`
}

func (sr *StatusResponse) UnmarshalJSON(data []byte) error {
	tempInt := 0
	tempSRStr := &statusResponseStr{}
	intErr := json.Unmarshal(data, &tempInt)
	strErr := json.Unmarshal(data, tempSRStr)

	if intErr == nil {
		sr.IntVal = tempInt
		sr.IsError = false
		sr.IsValid = true
		sr.StrVal = strconv.Itoa(sr.IntVal)
		return nil
	}

	if strErr == nil {
		sr.StrVal = tempSRStr.Status
		intVal, convErr := strconv.Atoi(sr.StrVal)
		sr.IntVal = intVal
		sr.IsValid = true
		if convErr != nil {
			sr.IsError = true
		} else {
			sr.IsError = false
		}
		return nil
	}

	return fmt.Errorf("StatusResponse is neither an int [%s] or a string [%s]", intErr.Error(), strErr.Error())
}
