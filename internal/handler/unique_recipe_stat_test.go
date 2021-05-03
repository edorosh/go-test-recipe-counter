package handler

import (
	"reflect"
	"testing"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/test"
)

func TestCountUniqueRecipesStat(t *testing.T) {
	recipes := []app.Recipe{
		test.RecipeSpeedyFajitas,
		test.RecipeCreamyChicken,
		test.RecipeCreamyChicken,
	}
	expectedResult := app.NewResult()
	expectedResult.UniqueRecipeCount = 2
	expectedResult.CountPerRecipeList = []app.CountPerRecipe{
		{
			"Creamy Dill Chicken",
			2,
		},
		{
			"Speedy Steak Fajitas",
			1,
		},
	}

	worker := NewUniqueRecipeStat()

	for _, r := range recipes {
		worker.Handle(r)
	}

	result := app.NewResult()
	worker.UpdateResult(&result)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("UniqueRecipeStatWorker got = %v, want %v", result, expectedResult)
	}
}
