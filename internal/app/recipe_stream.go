package app

// RecipeStream represents a stream of Recipes with handling errors.
// todo: define an interface and get rid of the concrete type
type RecipeStream struct {
	Recipe chan Recipe
	Err    chan error
}
