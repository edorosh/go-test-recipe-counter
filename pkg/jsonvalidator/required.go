package jsonvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

// RequiredFields checks required fields of a structure on JSON Unmarshall.
// If a value of a field is the zero value for its type it fails the validation.
// RequiredFields stops validating on the first error found.
// Examples of struct field tags and their meanings:
//	type JSONRecipe struct {
//		Postcode   int    `json:",required"`
//		RecipeName string `json:"recipe,required"`
//		Delivery   string
//	}
//	jr := JSONRecipe{}
//	if err := json.Unmarshal(data, &jr); err != nil {
//		return err
//	}
//	if err := json.RequiredFields(&jr); err != nil {
//		return err
//	}
func RequiredFields(f interface{}) (err error) {
	// IsZero may panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[json] Argument error: %v", r)
		}
	}()

	fields := reflect.ValueOf(f).Elem()

	for i := 0; i < fields.NumField(); i++ {
		jsonTags := fields.Type().Field(i).Tag.Get("json")
		if strings.Contains(jsonTags, ",required") && fields.Field(i).IsZero() {
			err = fmt.Errorf("[json] Required field \"%v\" is missing", fields.Type().Field(i).Name)
			return
		}
	}

	return
}
