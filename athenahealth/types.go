package athenahealth

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ptrStr mimics aws.String which helps shortcut getting a pointer to a string
func ptrStr(input string) *string { // nolint
	return &input
}

// ptrBool mimics aws.Bool which helps shortcut getting a pointer to a bool
func ptrBool(input bool) *bool { // nolint
	return &input
}

type NumberString string

func (n *NumberString) UnmarshalJSON(data []byte) error {
	var aux interface{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	switch v := aux.(type) {
	case string:
		*n = NumberString(v)
	case int:
		*n = NumberString(strconv.Itoa(v))
	case float64:
		*n = NumberString(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return fmt.Errorf("unknown type: %T", v)
	}

	return nil
}
