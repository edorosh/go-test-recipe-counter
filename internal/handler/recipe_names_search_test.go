package handler

import (
	"reflect"
	"testing"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/test"
)

func TestRecipeNamesSearchWorker(t *testing.T) {
	recipes := []app.Recipe{
		test.RecipeCreamyChicken,
		test.RecipeSpeedyFajitas,
		test.RecipeCreamyChicken,
		test.RecipeMeltyBurgers,
	}
	expectedResult := app.NewResult()
	expectedResult.MatchByName = []app.RecipeName{
		"Melty Monterey Jack Burgers",
		"Speedy Steak Fajitas",
	}

	worker := NewRecipeNamesSearch([]string{"Jack", "Steak"})

	for _, r := range recipes {
		worker.Handle(r)
	}

	result := app.NewResult()
	worker.UpdateResult(&result)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Count() got = %v, want %v", result, expectedResult)
	}
}
