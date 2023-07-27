package athenahealth

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// PtrStr mimics aws.String which helps shortcut getting a pointer to a string
func PtrStr(input string) *string {
	return &input
}

// PtrBool mimics aws.Bool which helps shortcut getting a pointer to a bool
func PtrBool(input bool) *bool {
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
