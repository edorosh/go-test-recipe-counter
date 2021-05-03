package handler

import (
	"sort"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
)

// UniqueRecipeStat counts all unique recipe names.
type UniqueRecipeStat struct {
	uniqueR map[app.RecipeName]app.CountPerRecipe
}

var _ app.Handler = &UniqueRecipeStat{}

// NewUniqueRecipeStat creates UniqueRecipeStat instance by reference
func NewUniqueRecipeStat() *UniqueRecipeStat {
	w := &UniqueRecipeStat{}
	w.uniqueR = make(map[app.RecipeName]app.CountPerRecipe)

	return w
}

// Handle a recipe
func (w *UniqueRecipeStat) Handle(r app.Recipe) {
	recipeName := r.RecipeName
	if cr, ok := w.uniqueR[recipeName]; !ok {
		w.uniqueR[recipeName] = app.CountPerRecipe{r.RecipeName, 1}
	} else {
		w.uniqueR[recipeName] = cr.Inc()
	}
}

// UpdateResult updates given Result with counters.
func (w *UniqueRecipeStat) UpdateResult(r *app.Result) {
	r.CountPerRecipeList = w.mapToSliceByKeys(w.keysSortByName())
	r.UniqueRecipeCount = len(r.CountPerRecipeList)
}

// todo: rewrite this with generics when GO 2.0 gets released
func (w *UniqueRecipeStat) keysSortByName() []app.RecipeName {
	keys := make(app.NamesByAlpha, 0, len(w.uniqueR))
	for k := range w.uniqueR {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	return keys
}

func (w *UniqueRecipeStat) mapToSliceByKeys(keys []app.RecipeName) []app.CountPerRecipe {
	v := make([]app.CountPerRecipe, 0, len(w.uniqueR))

	for _, r := range keys {
		v = append(v, w.uniqueR[r])
	}

	return v
}
