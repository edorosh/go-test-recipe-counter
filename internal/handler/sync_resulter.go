package handler

import (
	"context"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/jsonstream"
)

// SyncResulter processes streams of Recipe and merges results from RecipeHandlers.
type SyncResulter struct {
	handlers []app.Handler
	stream   *jsonstream.Reader
}

var _ app.Resulter = &SyncResulter{}

// NewSyncResulter creates SyncResulter by reference with the given handlers.
func NewSyncResulter(stream *jsonstream.Reader, handlers ...app.Handler) *SyncResulter {
	r := &SyncResulter{}
	r.handlers = handlers
	r.stream = stream

	return r
}

// Result returns Result with all counters being applied from Handler list.
func (r *SyncResulter) Result(ctx context.Context) (app.Result, error) {
	var resultErr error
	result := app.NewResult()
	ctx, cancel := context.WithCancel(ctx)

	r.stream.Handler(func(recipe app.Recipe) {
		for i := 0; i < len(r.handlers); i++ {
			r.handlers[i].Handle(recipe)
		}
	}).ErrHandler(func(err error) {
		resultErr = err
		cancel()
	}).Read(ctx)

	for i := 0; i < len(r.handlers); i++ {
		r.handlers[i].UpdateResult(&result)
	}

	return result, resultErr
}
