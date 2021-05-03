package app

import (
	"fmt"
	"sort"
	"strconv"
)

// Postcode represents a customer postcode in the recipe delivery list. One postcode is not longer than 10 chars.
type Postcode uint64

// PostcodeFromString parses Postcode from a string.
func PostcodeFromString(s string) (Postcode, error) {
	code, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[recipe] Postcode from string error: %v", err)
	}

	if code == 0 {
		return 0, fmt.Errorf("[recipe] Postcode is empty: %v", s)
	}

	return Postcode(code), nil
}

// PostcodesByCode implements asc sort.Interface based on the Postcode.
type PostcodesByCode []Postcode

var _ sort.Interface = PostcodesByCode{}

func (a PostcodesByCode) Len() int           { return len(a) }
func (a PostcodesByCode) Less(i, j int) bool { return a[i] < a[j] }
func (a PostcodesByCode) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
