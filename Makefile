MAKEFLAGS += --silent
GOLANGCI_CONFIG ?= .golangci.yml
DOCKER=podman

all: clean lint test build

## help: Prints a list of available build targets.
help:
	echo "Usage: make <OPTIONS> ... <TARGETS>"
	echo ""
	echo "Available targets are:"
	echo ''
	sed -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /'
	echo
	echo "Targets run by default are: `sed -n 's/^all: //p' ./Makefile | sed -e 's/ /, /g' | sed -e 's/\(.*\), /\1, and /'`"


## build: Builds a custom 'k6' with the local extension. 
build: xk6-config
	xk6 build --with $(shell go list -m)=.

## linter-config: Checks if the linter config exists, if not, downloads it from the main k6 repository.
linter-config:
	test -s "${GOLANGCI_CONFIG}" || (echo "No linter config, downloading from main k6 repository..." && curl --silent --show-error --fail --no-location https://raw.githubusercontent.com/grafana/k6/master/.golangci.yml --output "${GOLANGCI_CONFIG}")

## check-linter-version: Checks if the linter version is the same as the one specified in the linter config.
check-linter-version:
	(golangci-lint version | grep "version $(shell head -n 1 .golangci.yml | tr -d '\# ')") || echo "Your installation of golangci-lint is different from the one that is specified in k6's linter config (there it's $(shell head -n 1 .golangci.yml | tr -d '\# ')). Results could be different in the CI."

## xk6-config: checks that xk6 is installed
xk6-config:
	command -v xk6 2>&1 > /dev/null || $ go install go.k6.io/xk6/cmd/xk6@latest

## test: Executes any tests.
test:
	go test -count=1 -timeout 60s ./...

## testcontainer: launch test dependencies
testcontainer:
	$(DOCKER) compose -f test/compose.yaml up --force-recreate -d

## fulltest: full test suite
fulltest: build testcontainer
	echo "Waiting for test containers to start..." && sleep 5
	make test

## lint: Runs the linters.
lint: linter-config check-linter-version
	echo "Running linters..."
	golangci-lint run ./...

## check: Runs the linters and tests.
check: lint test

## clean: Removes any previously created artifacts/downloads.
clean:
	echo "Cleaning up..."
	rm -f ./k6
	rm -f .golangci.yml	

.PHONY: test lint check build clean linter-config check-linter-version
