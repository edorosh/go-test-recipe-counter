Recipe Stats Calculator
====

Recipe Stats Calculator is a `CLI` application which processes a JSON file with recipe data and calculated some stats.
The JSON result is rendered to stdout. Errors, additional info or debug is printed it to stderr. It reads the files in
small chunks so it requires a small memory footprint. It may handle Gigabytes of data.

The application handles signal `SIGTERM` (Ctr+C) and shuts down the Application gracefully. It returns exit code `0` on
success or cancellation and code `1` on failures.

Features
------

1. Count the number of unique recipe names.
2. Count the number of occurrences for each unique recipe name (alphabetically ordered by recipe name).
3. Find the postcode with most delivered recipes.
4. Count the number of deliveries to postcode `10120` that lie within the delivery time between `10AM` and `3PM`,
   examples _(`12AM` denotes midnight)_:
    - `NO` - `9AM - 2PM`
    - `YES` - `10AM - 2PM`
5. List the recipe names (alphabetically ordered) that contain in their name one of the following words:
    - Potato
    - Veggie
    - Mushroom

Installation
---------------
The Application uses GO Modules. You may clone the repository at any location and build the binary. All dependencies are
downloaded automatically on `go build` or `go test` commands.

### Compile

In order to compile the binary you may want to use make scripts from the project folder:

```shell
make all
```

It runs unit test first and compiles the binary in `dist/recipecounter`.

Here are all make options:

| Parameter                           | Description                                                   |
|-------------------------------------|---------------------------------------------------------------|
|all|Cleans project folder and module list, runs unit-tests and compile the binary|
|clean|Cleans project folder from previous build|
|build|Compiles binary in `dist/recipecounter`|
|test-unit|Runs unit tests|
|test-integration|Runs integration tests with coverage|
|bench-rootcmd|Benchmarks the Application|
|lint|Runs golangci linter|

### Docker

Optionally you may want to use multi-stage Docker image to build and run the application. To use the Application with
Docker run the following commands from project path:

```shell
docker build -t recipecounter:latest -f build/Dockerfile .
docker run --rm -v $(PWD):/app -w /app recipecounter:latest /bin/recipecounter cmd/fixtures/small.json  
```

Example of usages
---------------

```shell
recipecounter [filepath] [flags]
```

All available options are listed in the help message: `recipecounter --help`. There are a couple of test files in the
project. You may want to check them:

```shell
recipecounter cmd/fixtures/small.json
recipecounter cmd/fixtures/small.json --name=Jack --name=Speedy
recipecounter cmd/fixtures/small.json --name=Jack --name=Speedy --postcode-and-time="10120 10AM-3PM"
recipecounter cmd/invalidJSONRecipe.json 2> err.txt
```

Example of the input
---------------

```json5
[
  {
    "postcode": "10208",
    "recipe": "Speedy Steak Fajitas",
    "delivery": "Thursday 7AM - 5PM"
  },
  {
    "postcode": "10120",
    "recipe": "Cherry Balsamic Pork Chops",
    "delivery": "Thursday 10AM - 2PM"
  },
  {
    "postcode": "10120",
    "recipe": "Cherry Balsamic Pork Chops",
    "delivery": "Thursday 9AM - 2PM"
  }
]
```

Example of the output
---------------

Generate a JSON file of the following format:

```json5
{
  "unique_recipe_count": 15,
  "count_per_recipe": [
    {
      "recipe": "Mediterranean Baked Veggies",
      "count": 1
    },
    {
      "recipe": "Speedy Steak Fajitas",
      "count": 1
    },
    {
      "recipe": "Tex-Mex Tilapia",
      "count": 3
    }
  ],
  "busiest_postcode": {
    "postcode": "10120",
    "delivery_count": 1000
  },
  "count_per_postcode_and_time": {
    "postcode": "10120",
    "from": "11AM",
    "to": "3PM",
    "delivery_count": 500
  },
  "match_by_name": [
    "Mediterranean Baked Veggies",
    "Speedy Steak Fajitas",
    "Tex-Mex Tilapia"
  ]
}
```

Development
-----
The application has both unit and integration tests (with build tag `integration`). Integration test produces coverage
report. You may want to run test by using the following commands:

```shell
make test-unit
make test-integration
```

In order to benchmark some new features in the Command you may want to run the bench:

```shell
make bench-rootcmd
```

Implementation comments
-----------------------

* Project layout is based on [golang-standards](https://github.com/golang-standards/project-layout)
* Package `app` keeps domain objects: interfaces and core structs
* all dependencies in `inernal` point to `app`. Package `app` itself depends on `pkg` only
* `Postcode` and `RecipeName` have their dedicated types implementing Sort and Stringer interfaces accordingly. They are
  more likely to be optimized so working with domain interface simpler rather than pure types
* `NewRootCmd` handles all panic and prints errors in a user-friendly manner
* Error Handling is designed to return all the errors to the command (no panic, Fatalf)
* All JSON parsing errors are handled in JSON Decoder level. So all the handlers are guaranteed to have a valid data
* Handlers define sync interface. It makes them easy to be tested. Concurrency behaviour is done in separate structs.
* Concurrency communications is done by the package Context. Currently, JSON Decoder is cancellable only as the most
  time-consuming structure
* Any further design improvements should be first tested by Benchmarks if they do improve the performance
* In case there of >2 busiest postcodes the max postcode as a number is taken for the result

Performance Review
-----------------------
There were 3 possible solutions:

1. Fully Sync approach: read a Recipe from File -> (Sync) Each Recipe Handler counts a Recipe -> (Sync) Each Recipe
   Handler Writes into Result struct
2. Semi-Concurrent Count: read a Recipe from File -> (Concurrently) Each Recipe Handler counts a Recipe -> (Sync) Each
   Recipe Handler Writes into Result struct
3. Fully-Concurrent Count: read a Recipe from File -> (Concurrently) Each Recipe Handler counts a Recipe -> (
   Concurrently) Each Recipe Handler Writes into Result struct (Syncronization by Channels, immutable Result)

> All tests used the same command configuration.

Here are the Benchmarks results:

```shell
go test -tags=integration -run=XXX -bench=BenchmarkCmdParseFile ./cmd/
```

```shell
Fully Sync
BenchmarkCmdParseFile-4             3938            897256 ns/op
BenchmarkCmdParseFile-4             3067            928483 ns/op
BenchmarkCmdParseFile-4             4504           1023069 ns/op

Semi-Concurrent Count
BenchmarkCmdParseFile-4             3086            805381 ns/op
BenchmarkCmdParseFile-4             2606            713868 ns/op
BenchmarkCmdParseFile-4             2828            778577 ns/op

Fully-Concurrent Count
BenchmarkCmdParseFile-4             2098            642979 ns/op
BenchmarkCmdParseFile-4             2348            731585 ns/op
BenchmarkCmdParseFile-4             2762            796808 ns/op
```

Considering the Benchmarks above the fastest one is `Fully-Concurrent Count` but the difference is not that big.
Considering increasing design complexity there is no much sense to use it.

Here are results consuming test 1GB recipe data with `time`:

```shell
time dist/recipecounter 1GBRecipeList.json --name=Jack --name=Speedy --postcode-and-time="10120 10AM-3PM"
```

```shell
Fully Sync
real    0m51.525s
user    0m49.107s
sys     0m2.741s

Semi-Concurrent Count
real    1m44.879s
user    1m59.525s
sys     0m41.403s

Fully-Concurrent Count
real    1m37.532s
user    1m59.918s
sys     0m43.671s
```

On real data the simplest Fully Sync approach is faster. So it has been chosen for the final design. Memory profile is
required there.

TODO:
-----

* add travic-ci / github workflows
* add conv report
* add go quality report
* Make RecipeHandlers injectable into the Command and cover main.go by Unit Test
* analyze memory profile

