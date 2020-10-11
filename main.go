package main

import (
	"context"
	"errors"
	"os"

	"github.com/edorosh/go-test-recipe-counter/cmd"
	"github.com/edorosh/go-test-recipe-counter/internal/signal"
)

var version = "0.0.0-dev"

func main() {
	ctx := signal.ContextWithSIGTERM(context.Background())
	// todo: make handlers injected and write a unit test
	rootCmd := cmd.NewRootCmd(version)

	// By default the command displays results in STDOUT as well as the help
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		// errors returned should be explicitly printed in STDERR
		// PrintErrln has bug in V1.0
		rootCmd.PrintErr(err)

		if !errors.Is(err, context.Canceled) {
			os.Exit(1)
		}
	}
}
