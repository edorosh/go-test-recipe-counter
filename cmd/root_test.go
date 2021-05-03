// +build integration

package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"testing"
)

const VERSION = "1.0"

const FixtureFilepath = "fixtures/small.json"
const InvalidJsonFixtureFilepath = "fixtures/invalidJson.json"
const EmptyFixtureFilepath = "fixtures/empty.json"
const InvalidJsonRecipeFixtureFilepath = "fixtures/invalidJSONRecipe.json"

const ExpectedEmptyFileOutput = `{
  "unique_recipe_count": 0,
  "count_per_recipe": [],
  "busiest_postcode": {
    "postcode": "0",
    "delivery_count": 0
  },
  "count_per_postcode_and_time": {
    "postcode": "0",
    "from": "",
    "to": "",
    "delivery_count": 0
  },
  "match_by_name": []
}
`
const ExpectedFileOutput = `{
  "unique_recipe_count": 8,
  "count_per_recipe": [
    {
      "Recipe": "Cherry Balsamic Pork Chops",
      "Count": 5
    },
    {
      "Recipe": "Chicken Sausage Pizzas",
      "Count": 1
    },
    {
      "Recipe": "Hot Honey Barbecue Chicken Legs",
      "Count": 1
    },
    {
      "Recipe": "Mediterranean Baked Veggies",
      "Count": 1
    },
    {
      "Recipe": "Melty Monterey Jack Burgers",
      "Count": 2
    },
    {
      "Recipe": "Parmesan-Crusted Pork Tenderloin",
      "Count": 1
    },
    {
      "Recipe": "Speedy Steak Fajitas",
      "Count": 3
    },
    {
      "Recipe": "Steakhouse-Style New York Strip",
      "Count": 1
    }
  ],
  "busiest_postcode": {
    "postcode": "10120",
    "delivery_count": 2
  },
  "count_per_postcode_and_time": {
    "postcode": "0",
    "from": "",
    "to": "",
    "delivery_count": 0
  },
  "match_by_name": []
}
`
const ExpectedFileOutputWithSearchByName = `{
  "unique_recipe_count": 8,
  "count_per_recipe": [
    {
      "Recipe": "Cherry Balsamic Pork Chops",
      "Count": 5
    },
    {
      "Recipe": "Chicken Sausage Pizzas",
      "Count": 1
    },
    {
      "Recipe": "Hot Honey Barbecue Chicken Legs",
      "Count": 1
    },
    {
      "Recipe": "Mediterranean Baked Veggies",
      "Count": 1
    },
    {
      "Recipe": "Melty Monterey Jack Burgers",
      "Count": 2
    },
    {
      "Recipe": "Parmesan-Crusted Pork Tenderloin",
      "Count": 1
    },
    {
      "Recipe": "Speedy Steak Fajitas",
      "Count": 3
    },
    {
      "Recipe": "Steakhouse-Style New York Strip",
      "Count": 1
    }
  ],
  "busiest_postcode": {
    "postcode": "10120",
    "delivery_count": 2
  },
  "count_per_postcode_and_time": {
    "postcode": "0",
    "from": "",
    "to": "",
    "delivery_count": 0
  },
  "match_by_name": [
    "Melty Monterey Jack Burgers",
    "Speedy Steak Fajitas"
  ]
}
`
const ExpectedFileOutputWithPostcodeAndDate = `{
  "unique_recipe_count": 8,
  "count_per_recipe": [
    {
      "Recipe": "Cherry Balsamic Pork Chops",
      "Count": 5
    },
    {
      "Recipe": "Chicken Sausage Pizzas",
      "Count": 1
    },
    {
      "Recipe": "Hot Honey Barbecue Chicken Legs",
      "Count": 1
    },
    {
      "Recipe": "Mediterranean Baked Veggies",
      "Count": 1
    },
    {
      "Recipe": "Melty Monterey Jack Burgers",
      "Count": 2
    },
    {
      "Recipe": "Parmesan-Crusted Pork Tenderloin",
      "Count": 1
    },
    {
      "Recipe": "Speedy Steak Fajitas",
      "Count": 3
    },
    {
      "Recipe": "Steakhouse-Style New York Strip",
      "Count": 1
    }
  ],
  "busiest_postcode": {
    "postcode": "10120",
    "delivery_count": 2
  },
  "count_per_postcode_and_time": {
    "postcode": "10120",
    "from": "10AM",
    "to": "3PM",
    "delivery_count": 1
  },
  "match_by_name": []
}
`

func TestCmdFailures(t *testing.T) {
	var tests = []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			"Fail on file not found",
			[]string{"non_existing_file"},
			"Error: Can not read the file \"non_existing_file\"",
		},
		{
			"Fail on invalid JSON file",
			[]string{InvalidJsonFixtureFilepath},
			"Error: [jsonstream] JSON parse error: unexpected EOF",
		},
		{
			"Fail on invalid JSON Recipe in file",
			[]string{InvalidJsonRecipeFixtureFilepath},
			"Error: [jsonstream] JSON parse error: [json] Required field \"RecipeName\" is missing",
		},
		{
			"Fail on PostcodeAndTimeQueryParse flag",
			[]string{
				FixtureFilepath,
				"--postcode-and-time=blablabla",
			},
			"[recipe] Postcode from string error: strconv.ParseUint: parsing \"blablabla\": invalid syntax",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rootCmd := NewRootCmd(VERSION)
			output, err := executeCommand(rootCmd, tt.args...)

			if err == nil {
				t.Errorf("Expected the comamnd returning an error")
			}

			if tt.wantErr != err.Error() {
				t.Errorf("got: '%s' want: '%s'", err, tt.wantErr)
			}

			if output != "" {
				t.Errorf("Expected empty output, got: %s", output)
			}
		})
	}
}

func TestCmdPass(t *testing.T) {
	var tests = []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{
			"Pass on empty file",
			[]string{EmptyFixtureFilepath},
			ExpectedEmptyFileOutput,
		},
		{
			"Pass on real JSON file",
			[]string{FixtureFilepath},
			ExpectedFileOutput,
		},
		{
			"Pass on Search by RecipeName",
			[]string{FixtureFilepath, "--name=Jack", "--name=Speedy"},
			ExpectedFileOutputWithSearchByName,
		},
		{
			"Pass on Search by Postcode and Date",
			[]string{FixtureFilepath, "--postcode-and-time=10120 10AM-3PM"},
			ExpectedFileOutputWithPostcodeAndDate,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rootCmd := NewRootCmd(VERSION)
			output, err := executeCommand(rootCmd, tt.args...)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.wantOutput != output {
				t.Errorf("Expected output: %v, given %v", tt.wantOutput, output)
			}
		})
	}
}

func BenchmarkCmdParseFile(b *testing.B) {
	rootCmd := NewRootCmd(VERSION)
	for n := 0; n < b.N; n++ {
		executeCommand(rootCmd, FixtureFilepath, "--postcode-and-time=10120 10AM-3PM", "--name=Jack", "--name=Speedy")
	}
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
