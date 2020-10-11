package handler

import (
	"sort"
	"strings"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
)

// RecipeNamesSearch finds recipes by names.
type RecipeNamesSearch struct {
	searches []string
	matches  map[app.RecipeName]struct{}
}

var _ app.Handler = &RecipeNamesSearch{}

// NewRecipeNamesSearch creates new NewRecipeNamesSearch instance by reference.
func NewRecipeNamesSearch(s []string) *RecipeNamesSearch {
	w := &RecipeNamesSearch{}
	w.searches = s
	w.matches = make(map[app.RecipeName]struct{})

	return w
}

// Handle a recipe
func (w *RecipeNamesSearch) Handle(r app.Recipe) {
	name := r.RecipeName

	for _, search := range w.searches {
		if !strings.Contains(name.String(), search) {
			continue
		}

		if _, ok := w.matches[name]; !ok {
			w.matches[name] = struct{}{}
		}
	}
}

// UpdateResult updates given Result with counters.
func (w *RecipeNamesSearch) UpdateResult(r *app.Result) {
	if len(w.matches) == 0 {
		r.MatchByName = app.NamesByAlpha{}
		return
	}

	r.MatchByName = w.keysSortByName()
}

// keysSortByName sorts list of recipe names (strings) by name in alphabetical order.
// todo: rewrite this with generics when GO 2.0 gets released
func (w *RecipeNamesSearch) keysSortByName() []app.RecipeName {
	keys := make(app.NamesByAlpha, 0, len(w.matches))
	for k := range w.matches {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	return keys
}
