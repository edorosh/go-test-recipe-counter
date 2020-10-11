package delivery

import (
	"fmt"
	"strconv"
)

// TimeSchema represents 12h time string as "7PM".
type TimeSchema string

// PadHours makes 12h value valid for Time parsing. It adds 0 to one digit 12h value and a whitespace if needed.
// Ex. 1PM is changed to 01 PM.
func (t TimeSchema) PadHours() (string, error) {
	schema := t[len(t)-2:]

	if schema != "AM" && schema != "PM" {
		return "", fmt.Errorf("[delivery] PadHours() unsupported schema: %v", schema)
	}

	tm := t[:len(t)-2]

	if len(tm) == 2 {
		return fmt.Sprintf("%v %v", tm, schema), nil
	}

	iTime, err := strconv.Atoi(string(tm))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%02d %v", iTime, schema), nil
}
