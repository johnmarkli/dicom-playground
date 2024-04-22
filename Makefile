BIN = $(GOPATH)/bin
GOLINT_VERSION=1.57.2
GOLINT = $(BIN)/golangci-lint
PROJECT_BASE=$(shell basename "$(PWD)")

all: lint test

echo:
	$(info Base=$(PROJECT_BASE) )
	$(info pwd=$(PWD) )

lint:
	$(info Linting code )
ifneq '$(GOLINT_VERSION)' '$(shell golangci-lint version 2>/dev/null | cut -d " " -f 4 -)'
	@echo Getting golangci-lint v$(GOLINT_VERSION)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) v$(GOLINT_VERSION)
	$Q$(GOLINT) run
else
	@echo Using current version
	$Q$(GOLINT) run
endif

test:
	$(info Running gotest )
	go test -coverprofile=coverage.out ./...
