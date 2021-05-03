package handler

import (
	"reflect"
	"testing"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/test"
	"github.com/edorosh/go-test-recipe-counter/pkg/delivery"
)

func TestDeliverySearchWorker(t *testing.T) {
	recipes := []app.Recipe{
		test.RecipeCreamyChicken,
		test.RecipeSpeedyFajitas,
		test.RecipeCreamyChicken2,
	}

	dFrom, err := delivery.NewTime("10 AM")
	if err != nil {
		t.Fatalf("NewTime: Unexpected error %v", err)
	}

	dTo, err := delivery.NewTime("03 PM")
	if err != nil {
		t.Fatalf("NewTime: Unexpected error %v", err)
	}

	expectedResult := app.NewResult()
	expectedResult.CountPerPostcodeAndTimeList = app.CountPerPostcodeAndTime{
		test.RecipeCreamyChicken.Postcode,
		dFrom,
		dTo,
		1,
	}

	query, err := PostcodeAndTimeQueryParse("10224 10AM-3PM")
	if err != nil {
		t.Fatalf("NewDeliverySearch: Unexpected error %v", err)
	}

	worker := NewDeliverySearch(query)

	for _, r := range recipes {
		worker.Handle(r)
	}

	result := app.NewResult()
	worker.UpdateResult(&result)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Count() got = %v, want %v", result, expectedResult)
	}
}
