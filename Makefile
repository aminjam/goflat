GOTOOLS = github.com/mitchellh/gox github.com/FiloSottile/gvt
PACKAGES=$(shell go list ./... | grep -v vendor | sort | uniq)
BINARY_NAME=$(shell basename ${PWD})
MAIN_PACKAGE="."

all: format init test build-dist

build:
	@$(CURDIR)/scripts/build.bash $(BINARY_NAME) $(MAIN_PACKAGE) dev

build-dist:
	@$(CURDIR)/scripts/build.bash $(BINARY_NAME) $(MAIN_PACKAGE)

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

test:
	@echo "--> Running Tests"
	@go test

init:
	@echo "--> Init build tools"
	@go get -v $(GOTOOLS)

update-deps:
	@echo "--> Updating dependencies"
	@(MAKE) tools
	@gvt update --all

.PHONY: all build build-dist format init test update-deps


