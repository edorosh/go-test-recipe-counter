package app

// Handler processes stream of Recipes.
type Handler interface {
	Handle(r Recipe)
	UpdateResult(r *Result)
}
