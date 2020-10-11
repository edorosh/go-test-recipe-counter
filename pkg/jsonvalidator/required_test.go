package jsonvalidator

import (
	"encoding/json"
	"reflect"
	"testing"
)

const JSONCorrect = `{
	"postcode": 10224,
	"recipe": "Creamy Dill Chicken",
	"delivery": "Wednesday 1AM - 7PM"
}`
const JSONMissingProperty = `{
	"postcode": 10224,
	"rcp": "Creamy Dill Chicken",
	"delivery": "Wednesday 1AM - 7PM"
}`
const JSONPropertyIntDefaultValue = `{
	"postcode": 0,
	"recipe": "Creamy Dill Chicken",
	"delivery": "Wednesday 1AM - 7PM"
}`
const JSONPropertyStringDefaultValue = `{
	"postcode": 10224,
	"recipe": "",
	"delivery": "Wednesday 1AM - 7PM"
}`

type JSONRecipe struct {
	Postcode   int    `json:",required"`
	RecipeName string `json:"recipe,required"`
	Delivery   string
}

var recipe = JSONRecipe{
	10224,
	"Creamy Dill Chicken",
	"Wednesday 1AM - 7PM",
}

func TestValidateRequiredFields(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantRecipe *JSONRecipe
		wantErr    bool
	}{
		{
			"Pass on valid JSON",
			args{
				JSONCorrect,
			},
			&recipe,
			false,
		},
		{
			"Fail on missing property",
			args{
				JSONMissingProperty,
			},
			nil,
			true,
		},
		{
			"Fail on int with default value",
			args{
				JSONPropertyIntDefaultValue,
			},
			nil,
			true,
		},
		{
			"Fail on string with default value",
			args{
				JSONPropertyStringDefaultValue,
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r JSONRecipe
			if err := json.Unmarshal([]byte(tt.args.s), &r); err != nil {
				t.Errorf("Unexpected Unmarshal() error: %v", err)
			}

			if err := RequiredFields(&r); (err != nil) != tt.wantErr {
				t.Errorf("RequiredFields() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantRecipe != nil && !reflect.DeepEqual(r, *tt.wantRecipe) {
				t.Errorf("Read() got = %v, want %v", r, *tt.wantRecipe)
			}
		})
	}
}
