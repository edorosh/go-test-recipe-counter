package app

import (
	"fmt"
	"sort"
)

// RecipeName represents a name in a Recipe.
type RecipeName string

var _ fmt.Stringer = RecipeName("")

func (r RecipeName) String() string {
	return string(r)
}

// NamesByAlpha implements sort.Interface alphabetically ordered based on the RecipeName.
type NamesByAlpha []RecipeName

var _ sort.Interface = NamesByAlpha{}

func (a NamesByAlpha) Len() int           { return len(a) }
func (a NamesByAlpha) Less(i, j int) bool { return a[i] < a[j] }
func (a NamesByAlpha) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
