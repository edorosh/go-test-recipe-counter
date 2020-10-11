NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

VERSION ?= "dev-$(shell git rev-parse --short HEAD)"
GO_LINKER_FLAGS=-ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: all clean lint test-unit test-integration bench-rootcmd build

all: clean test-unit build

# Cleans our project: deletes binaries
clean:
	@printf "$(OK_COLOR)==> Cleaning project$(NO_COLOR)\n"
	@if [ -d dist ] ; then rm -rf dist/* ; fi; go mod tidy

build:
	@echo "$(OK_COLOR)==> Building default binary... $(NO_COLOR)"
	@CGO_ENABLED=0 go build ${GO_LINKER_FLAGS} -o "dist/recipecounter"

test-unit:
	@echo "$(OK_COLOR)==> Running unit tests$(NO_COLOR)"
	@CGO_ENABLED=0 go test -short ./...

test-integration:
	@printf "$(OK_COLOR)==> Running integration tests$(NO_COLOR)\n"
	@CGO_ENABLED=1 go test -tags=integration -race -p=1 -cover ./... -coverprofile=coverage.txt -covermode=atomic

bench-rootcmd:
	@printf "$(OK_COLOR)==> Running Benchmarking on room cmd$(NO_COLOR)\n"
	@go test -tags=integration -run=XXX -bench=BenchmarkCmdParseFile ./cmd/

lint:
	@echo "$(OK_COLOR)==> Linting with golangci-lint running in docker container$(NO_COLOR)"
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.30.0 golangci-lint run -v
