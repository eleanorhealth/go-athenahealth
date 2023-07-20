package athenahealth

import (
	"encoding/json"
	"errors"
	"io"
)

type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ErrorMessageResponse struct {
	Message string `json:"errormessage"`
	Success bool   `json:"success"`
}

// BodyToErrorString will handle the obtuse case where Athena returns an error
// message as a raw string with zero JSON formatting.
func BodyToErrorString(body io.ReadCloser) error {
	outStr := ""
	bodyBytes, readErr := io.ReadAll(body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(bodyBytes, &outStr)
	if jsonErr != nil {
		return jsonErr
	}

	return errors.New(outStr)
}
