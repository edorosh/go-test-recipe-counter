package handler

import (
	"reflect"
	"testing"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/test"
)

func TestBusiestPostcodeWorker(t *testing.T) {
	recipes := []app.Recipe{
		test.RecipeCreamyChicken,
		test.RecipeSpeedyFajitas,
		test.RecipeCreamyChicken,
	}
	expectedResult := app.NewResult()
	expectedResult.BusiestPostcode = app.BusiestPostcode{
		test.RecipeCreamyChicken.Postcode,
		app.DeliveryCount(2),
	}

	worker := NewBusiestPostcode()

	for _, r := range recipes {
		worker.Handle(r)
	}

	result := app.NewResult()
	worker.UpdateResult(&result)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Count() got = %v, want %v", result, expectedResult)
	}
}
