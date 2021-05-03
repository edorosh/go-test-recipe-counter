package app

import (
	"github.com/edorosh/go-test-recipe-counter/pkg/delivery"
)

// DeliveryCount represents the number of unique recipe names.
type DeliveryCount int

// CountPerRecipe represents an immutable counter per a Recipe RecipeName.
type CountPerRecipe struct {
	RecipeName `json:"Recipe"`
	Count      int
}

// Inc increases the counter and returns new CountPerRecipe.
func (c CountPerRecipe) Inc() CountPerRecipe {
	return CountPerRecipe{
		c.RecipeName,
		c.Count + 1,
	}
}

// CountPerPostcodeAndTime represents a Recipe counter found by a postcode and delivery range.
type CountPerPostcodeAndTime struct {
	Postcode      `json:"postcode,string"`
	From          delivery.Time `json:"from"`
	To            delivery.Time `json:"to"`
	DeliveryCount `json:"delivery_count"`
}

// BusiestPostcode represents the postcode with most delivered recipes.
type BusiestPostcode struct {
	Postcode      `json:"postcode,string"`
	DeliveryCount `json:"delivery_count"`
}

// Result represents final stat analyses.
type Result struct {
	UniqueRecipeCount           int                     `json:"unique_recipe_count"`
	CountPerRecipeList          []CountPerRecipe        `json:"count_per_recipe"`
	BusiestPostcode             BusiestPostcode         `json:"busiest_postcode"`
	CountPerPostcodeAndTimeList CountPerPostcodeAndTime `json:"count_per_postcode_and_time"`
	MatchByName                 []RecipeName            `json:"match_by_name"`
}

// NewResult creates new Statistics Result by a value.
func NewResult() Result {
	result := Result{}

	result.CountPerRecipeList = make([]CountPerRecipe, 0)
	result.MatchByName = make([]RecipeName, 0)

	return result
}
