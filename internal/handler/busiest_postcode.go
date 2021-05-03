package handler

import (
	"sort"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
)

// BusiestPostcode find the busiest Postcode in the recipe list
type BusiestPostcode struct {
	postcodes map[app.Postcode]int
}

var _ app.Handler = &BusiestPostcode{}

// NewBusiestPostcode creates new BusiestPostcode instance by reference.
func NewBusiestPostcode() *BusiestPostcode {
	w := &BusiestPostcode{}
	w.postcodes = make(map[app.Postcode]int)

	return w
}

// Handle a recipe
func (w *BusiestPostcode) Handle(r app.Recipe) {
	postcode := r.Postcode
	if _, ok := w.postcodes[postcode]; !ok {
		w.postcodes[postcode] = 1
	} else {
		w.postcodes[postcode]++
	}
}

// UpdateResult updates given Result with counters.
func (w *BusiestPostcode) UpdateResult(r *app.Result) {
	if len(w.postcodes) == 0 {
		r.BusiestPostcode = app.BusiestPostcode{}
		return
	}

	max := w.busiestPostcodeInMap()

	r.BusiestPostcode = app.BusiestPostcode{
		max,
		app.DeliveryCount(w.postcodes[max]),
	}
}

func (w *BusiestPostcode) busiestPostcodeInMap() app.Postcode {
	var (
		postcode app.Postcode
		maxDel   int
	)

	if len(w.postcodes) == 0 {
		return postcode
	}

	keys := make(app.PostcodesByCode, 0, len(w.postcodes))
	for k := range w.postcodes {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	for _, k := range keys {
		if maxDel < w.postcodes[app.Postcode(k)] {
			maxDel = w.postcodes[app.Postcode(k)]
			postcode = app.Postcode(k)
		}
	}

	return postcode
}
