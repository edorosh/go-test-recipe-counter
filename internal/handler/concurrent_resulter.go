package handler

import (
	"context"
	"sync"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/jsonstream"
)

// ConcurrentResulter processes streams of Recipe concurrently and merges results from RecipeHandlers.
type ConcurrentResulter struct {
	handlers []app.Handler
	stream   *jsonstream.Reader
}

var _ app.Resulter = &ConcurrentResulter{}

// NewConcurrentResulter creates ConcurrentResulter by reference with the given handlers.
func NewConcurrentResulter(stream *jsonstream.Reader, handlers ...app.Handler) *ConcurrentResulter {
	r := &ConcurrentResulter{}
	r.handlers = handlers
	r.stream = stream

	return r
}

// Result returns Result with all counters being applied from Handler list.
func (r *ConcurrentResulter) Result(ctx context.Context) (app.Result, error) {
	result := app.NewResult()

	if err := r.handleRecipes(ctx); err != nil {
		return result, err
	}

	// get Result sync
	for i := 0; i < len(r.handlers); i++ {
		r.handlers[i].UpdateResult(&result)
	}

	return result, nil
}

func (r *ConcurrentResulter) handleRecipes(ctx context.Context) error {
	hNum := len(r.handlers)
	ctx, cancel := context.WithCancel(ctx)

	// Create a read channel for reach handler
	rChans := make([]chan app.Recipe, hNum)
	for i := 0; i < hNum; i++ {
		rChans[i] = make(chan app.Recipe)
	}

	// Start handlers in goroutines and bind specific channel
	var wg sync.WaitGroup
	wg.Add(hNum)
	for i := 0; i < hNum; i++ {
		go func(h app.Handler, c chan app.Recipe) {
			defer wg.Done()
			for r := range c {
				h.Handle(r)
			}
		}(r.handlers[i], rChans[i])
	}

	// Process Recipe Stream
	var err error
	r.stream.Handler(func(recipe app.Recipe) {
		// Broadcast Recipes to handlers' channels.
		// On cancel request it stops on next read from the chan
		for i := 0; i < hNum; i++ {
			rChans[i] <- recipe
		}
	}).ErrHandler(func(serr error) {
		err = serr
		cancel()
	}).Read(ctx)

	// Close handlers' channels and stop goroutines
	for i := 0; i < hNum; i++ {
		close(rChans[i])
	}

	// Wait till handlers' goroutines stop
	wg.Wait()

	return err
}
