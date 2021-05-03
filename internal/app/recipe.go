package app

import (
	"encoding/json"

	"github.com/edorosh/go-test-recipe-counter/pkg/delivery"
	"github.com/edorosh/go-test-recipe-counter/pkg/jsonvalidator"
)

// Recipe represents an entry in Recipe List.
type Recipe struct {
	Postcode
	RecipeName
	delivery.TimeRange
}

// JSONRecipe represents a Recipe in JSON format.
type JSONRecipe struct {
	Postcode   `json:",string,required"`
	RecipeName `json:"recipe,required"`
	Delivery   string `json:",required"`
}

// UnmarshalJSON recovers a Recipe from a JSON entry.
func (r *Recipe) UnmarshalJSON(data []byte) error {
	var (
		err       error
		timeRange delivery.TimeRange
	)

	jr := JSONRecipe{}

	if err = json.Unmarshal(data, &jr); err != nil {
		return err
	}
	if err = jsonvalidator.RequiredFields(&jr); err != nil {
		return err
	}

	timeRange, err = delivery.TimeRangeWithDateParse(jr.Delivery)
	if err != nil {
		return err
	}

	r.TimeRange = timeRange
	r.RecipeName = jr.RecipeName
	r.Postcode = jr.Postcode

	return nil
}
