package delivery

import (
	"fmt"
	"strings"
	"time"
)

// TimeRange represents a Time range with from and do time values.
type TimeRange struct {
	from time.Time
	to   time.Time
}

// DeliveryFrom returns Time from.
func (d TimeRange) DeliveryFrom() Time {
	return Time{&d.from}
}

// DeliveryTo returns Time to.
func (d TimeRange) DeliveryTo() Time {
	return Time{&d.to}
}

// InTimeSpan returns true in case TimeRange fits into the given TimeRange (inclusively)
func (d TimeRange) InTimeSpan(dd TimeRange) bool {
	return (d.from.Equal(dd.from) && d.to.Equal(dd.to)) ||
		(d.from.Equal(dd.from) && d.to.Before(dd.to)) ||
		(d.from.After(dd.from) && d.to.Before(dd.to)) ||
		(d.from.After(dd.from) && d.to.Equal(dd.to))
}

// NewTimeRangeAsStrings creates TimeRange from the given from and to as stings ("7AM", "9PM").
func NewTimeRangeAsStrings(from, to string) (TimeRange, error) {
	var (
		err      error
		dateFrom time.Time
		dateTo   time.Time
	)

	fromStr, err := TimeSchema(from).PadHours()
	if err != nil {
		return TimeRange{}, fmt.Errorf("NewTimeRangeAsStrings time error: %v", from)
	}

	toStr, err := TimeSchema(to).PadHours()
	if err != nil {
		return TimeRange{}, fmt.Errorf("NewTimeRangeAsStrings time error: %v", from)
	}

	dateFrom, err = time.Parse(TimeFormat, fromStr)
	if err != nil {
		return TimeRange{}, fmt.Errorf("NewTimeRangeAsStrings time error: %v", from)
	}

	dateTo, err = time.Parse(TimeFormat, toStr)
	if err != nil {
		return TimeRange{}, fmt.Errorf("NewTimeRangeAsStrings time error: %v", to)
	}

	return TimeRange{dateFrom, dateTo}, nil
}

// TimeRangeWithDateParse creates TimeRange from the Time and Date string ("Thursday 7AM - 9PM").
func TimeRangeWithDateParse(s string) (TimeRange, error) {
	tokens := strings.Split(s, " ")

	length := len(tokens)
	if length > 4 || length == 0 {
		return TimeRange{}, fmt.Errorf("TimeRangeWithDateParse error: %v", s)
	}

	return NewTimeRangeAsStrings(tokens[1], tokens[3])
}

// TimeRangeAsStringParse creates TimeRange from the given sting ("7AM-9PM").
func TimeRangeAsStringParse(s string) (TimeRange, error) {
	tokens := strings.Split(s, "-")
	length := len(tokens)
	if length > 2 || length == 0 {
		return TimeRange{}, fmt.Errorf("TimeRangeAsStringParse time error: %v", s)
	}

	deliveryRange, err := NewTimeRangeAsStrings(tokens[0], tokens[1])
	if err != nil {
		return TimeRange{}, fmt.Errorf("TimeRangeAsStringParse time error: %v", err)
	}

	return deliveryRange, nil
}
