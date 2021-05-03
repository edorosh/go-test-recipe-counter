package delivery

import (
	"encoding/json"
	"time"
)

// TimeFormat represents a Time Format being used in a Recipe
const TimeFormat = "03 PM"

// CompactTimeFormat represents a compact Recipe Time Format
const CompactTimeFormat = "3PM"

// Time wraps time.Time type keeping it as a reference inside. Typical use case is JSON serialization if empty value by
// default.
type Time struct {
	*time.Time
}

// NewTime creates Time structure and am/pm schema
func NewTime(t string) (Time, error) {
	time, err := time.Parse(TimeFormat, t)
	if err != nil {
		return Time{}, err
	}

	return Time{&time}, nil
}

func (t Time) String() string {
	if t.Time == nil {
		return ""
	}

	return t.Time.Format(CompactTimeFormat)
}

// MarshalJSON implements json.Marshaler interface
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
