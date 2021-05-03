package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
	"github.com/edorosh/go-test-recipe-counter/internal/handler"
	"github.com/edorosh/go-test-recipe-counter/internal/jsonstream"
	"github.com/spf13/cobra"
)

// RootCmdOptions are the command flags
type RootCmdOptions struct {
	postcodeandtime string
	recipeNames     []string
}

// NewRootCmd creates a new instance of the root command.
func NewRootCmd(version string) *cobra.Command {
	opts := RootCmdOptions{}

	cmd := &cobra.Command{
		Use:     "recipecounter [filepath]",
		Version: version,
		Short:   "Recipe Counter analyses recipe data.",
		Long:    `Recipe Counter is a CLI application which parses recipe data and calculates some stats.`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRootCmdFunc(cmd, &opts, args)
		},
	}

	// Mute standard Cobra behaviour and use stdout for the result only (or --help) and errout for errors and infos
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	flags := cmd.PersistentFlags()
	flags.StringVarP(&opts.postcodeandtime, "postcode-and-time", "f", "", "(optional) Custom postcode and time window for search. Format: \"{postcode} {h}AM - {h}PM\"")
	flags.StringSliceVarP(&opts.recipeNames, "name", "n", []string{}, "(optional) Custom recipe names for search.")

	return cmd
}

func runRootCmdFunc(cmd *cobra.Command, opts *RootCmdOptions, args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error: %v", r)
		}
	}()

	var query handler.PostcodeAndTimeQuery
	if query, err = handler.PostcodeAndTimeQueryParse(opts.postcodeandtime); err != nil {
		return err
	}

	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error: Can not read the file \"%v\"", filename)
	}

	handlers := []app.Handler{
		handler.NewUniqueRecipeStat(),
		handler.NewBusiestPostcode(),
		handler.NewRecipeNamesSearch(opts.recipeNames),
		handler.NewDeliverySearch(query),
	}

	resulter := handler.NewSyncResulter(jsonstream.NewReader(file), handlers...)
	result, err := resulter.Result(cmd.Context())
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	cmd.Println(string(jsonResult))

	return nil
}
