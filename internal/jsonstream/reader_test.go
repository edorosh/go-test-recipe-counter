package jsonstream

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/test"
)

const emptyJSON = ``
const emptyJSONList = `[]`
const invalidJSON = `[{`
const invalidJSONRecipeWithTrailingComma = `[
  {
    "postcode": "10224",
    "recipe": "Creamy Dill Chicken",
    "delivery": "Wednesday 9AM - 2PM"
  },
]`
const JSONRecipeMissingRequiredField = `[
  {
    "postcode": "10224",
    "rcp": "Creamy Dill Chicken",
    "delivery": "Wednesday 9AM - 2PM"
  }
]`

const JSONRecipePostcodeInt = `[
  {
    "postcode": 10224,
    "rcp": "Creamy Dill Chicken",
    "delivery": "Wednesday 9AM - 2PM"
  }
]`

const singleJSONRecipe = `[
  {
    "postcode": "10224",
    "recipe": "Creamy Dill Chicken",
    "delivery": "Wednesday 9AM - 2PM"
  }
]`

const coupleJSONRecipes = `[
  {
    "postcode": "10224",
    "recipe": "Creamy Dill Chicken",
    "delivery": "Wednesday 9AM - 2PM"
  },
  {
    "postcode": "10208",
    "recipe": "Speedy Steak Fajitas",
    "delivery": "Thursday 7AM - 5PM"
  }
]`

func TestJsonStreamDecoder(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		input io.Reader
	}
	tests := []struct {
		name        string
		args        args
		wantRecipes []app.Recipe
		wantErr     error
	}{
		{
			"Pass single Recipe",
			args{
				ctx,
				strings.NewReader(singleJSONRecipe),
			},
			[]app.Recipe{test.RecipeCreamyChicken},
			nil,
		},
		{
			"Pass a couple of Recipes",
			args{
				ctx,
				strings.NewReader(coupleJSONRecipes),
			},
			[]app.Recipe{
				test.RecipeCreamyChicken,
				test.RecipeSpeedyFajitas,
			},
			nil,
		},
		{
			"Fail single Recipe with trailing coma",
			args{
				ctx,
				strings.NewReader(invalidJSONRecipeWithTrailingComma),
			},
			[]app.Recipe{test.RecipeCreamyChicken},
			ErrJSONParse,
		},
		{
			"Pass empty JSON",
			args{
				ctx,
				strings.NewReader(emptyJSON),
			},
			nil,
			nil,
		},
		{
			"Pass JSON empty Recipe List",
			args{
				ctx,
				strings.NewReader(emptyJSONList),
			},
			nil,
			nil,
		},
		{
			"Fail Invalid JSON",
			args{
				ctx,
				strings.NewReader(invalidJSON),
			},
			nil,
			ErrJSONParse,
		},
		{
			"Fail JSON Recipe with missing required field",
			args{
				ctx,
				strings.NewReader(JSONRecipeMissingRequiredField),
			},
			nil,
			ErrJSONParse,
		},
		{
			"Fail JSON Recipe with Postcode int",
			args{
				ctx,
				strings.NewReader(JSONRecipePostcodeInt),
			},
			nil,
			ErrJSONParse,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var recipes []app.Recipe
			var err error

			decoder := NewReader(tt.args.input)
			decoder.Handler(func(r app.Recipe) {
				recipes = append(recipes, r)
			}).ErrHandler(func(serr error) {
				err = serr
			}).Read(ctx)

			if tt.wantErr == nil && err != nil {
				t.Errorf("Read() Expected no error, got %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Read() got error = %v, want %v", err, tt.wantErr)
			}

			if len(recipes) != len(tt.wantRecipes) {
				t.Errorf("Read() got # of recipes = %v, want %v", len(recipes), len(tt.wantRecipes))
			}

			if !reflect.DeepEqual(recipes, tt.wantRecipes) {
				t.Errorf("Read() got = %v, want %v", recipes, tt.wantRecipes)
			}
		})
	}
}
