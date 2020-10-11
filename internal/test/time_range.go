package test

import "github.com/edorosh/go-test-recipe-counter/pkg/delivery"

// newTimeRangeAsStrings creates TimeRange or panics on failure
func newTimeRangeAsStrings(from, to string) delivery.TimeRange {
	var (
		timeRange delivery.TimeRange
		err       error
	)

	timeRange, err = delivery.NewTimeRangeAsStrings(from, to)
	if err != nil {
		panic(err)
	}

	return timeRange
}
