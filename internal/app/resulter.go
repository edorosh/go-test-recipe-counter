package app

import (
	"context"
)

// Resulter processes stream of Recipes and returns Statistics Result.
type Resulter interface {
	Result(context.Context) (Result, error)
}
