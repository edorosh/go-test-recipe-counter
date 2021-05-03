package test

import (
	"github.com/edorosh/go-test-recipe-counter/internal/app"
)

var (
	// RecipeCreamyChicken a Recipe fixture for Tests
	RecipeCreamyChicken = app.Recipe{
		10224,
		"Creamy Dill Chicken",
		newTimeRangeAsStrings("9AM", "2PM"),
	}

	// RecipeCreamyChicken2 a Recipe fixture for Tests
	RecipeCreamyChicken2 = app.Recipe{
		10224,
		"Creamy Dill Chicken",
		newTimeRangeAsStrings("10AM", "2PM"),
	}

	// RecipeSpeedyFajitas a Recipe fixture for Tests
	RecipeSpeedyFajitas = app.Recipe{
		10208,
		"Speedy Steak Fajitas",
		newTimeRangeAsStrings("7AM", "5PM"),
	}
	// RecipeMeltyBurgers a Recipe fixture for Tests
	RecipeMeltyBurgers = app.Recipe{
		10124,
		"Melty Monterey Jack Burgers",
		newTimeRangeAsStrings("9AM", "6PM"),
	}
)
