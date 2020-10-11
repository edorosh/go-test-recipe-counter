package handler

import (
	"fmt"
	"strings"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/pkg/delivery"
)

// PostcodeAndTimeQuery represents a query for searching recipes by a postcode and delivery range.
type PostcodeAndTimeQuery struct {
	delivery.TimeRange
	app.Postcode
}

// PostcodeAndTimeQueryParse creates new PostcodeAndTimeQuery from a string.
func PostcodeAndTimeQueryParse(s string) (PostcodeAndTimeQuery, error) {
	var query PostcodeAndTimeQuery
	var length int

	if len(s) == 0 {
		return query, nil
	}

	// 10120 7AM-9PM
	tokens := strings.Split(s, " ")

	length = len(tokens)
	if length > 2 || length == 0 {
		return PostcodeAndTimeQuery{}, fmt.Errorf("PostcodeAndTimeQueryParse postcode error: %v", s)
	}

	postCode, err := app.PostcodeFromString(tokens[0])
	if err != nil {
		return PostcodeAndTimeQuery{}, err
	}

	deliveryRange, err := delivery.TimeRangeAsStringParse(tokens[1])
	if err != nil {
		return PostcodeAndTimeQuery{}, err
	}

	query.Postcode = postCode
	query.TimeRange = deliveryRange

	return query, nil
}

// DeliverySearch counts deliveries by a postcode and delivery range.
type DeliverySearch struct {
	query   PostcodeAndTimeQuery
	counter int
}

var _ app.Handler = &DeliverySearch{}

// NewDeliverySearch creates DeliverySearch instance by reference.
func NewDeliverySearch(q PostcodeAndTimeQuery) *DeliverySearch {
	w := &DeliverySearch{}
	w.query = q

	return w
}

// Handle a recipe
func (w *DeliverySearch) Handle(r app.Recipe) {
	if r.Postcode != w.query.Postcode {
		return
	}

	if !r.TimeRange.InTimeSpan(w.query.TimeRange) {
		return
	}

	w.counter++
}

// UpdateResult updates given Result with counters.
func (w *DeliverySearch) UpdateResult(r *app.Result) {
	if w.counter == 0 {
		r.CountPerPostcodeAndTimeList = app.CountPerPostcodeAndTime{}
		return
	}

	r.CountPerPostcodeAndTimeList = app.CountPerPostcodeAndTime{
		w.query.Postcode,
		w.query.TimeRange.DeliveryFrom(),
		w.query.TimeRange.DeliveryTo(),
		app.DeliveryCount(w.counter),
	}
}
