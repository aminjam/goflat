GOTOOLS = github.com/mitchellh/gox \
	  github.com/FiloSottile/gvt \

PACKAGES=$(shell go list ./... | grep -v vendor | sort | uniq)
BINARY_NAME=$(shell basename ${PWD})
MAIN_PACKAGE="."

all: format test build-dist

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

update-deps:
	@echo "--> Updating dependencies"
	@go get -v $(GOTOOLS)
	@gvt update --all

.PHONY: all build build-dist format test update-deps


