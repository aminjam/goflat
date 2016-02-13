GOTOOLS = github.com/mitchellh/gox github.com/FiloSottile/gvt
PACKAGES=$(shell go list ./... | grep -v vendor | sort | uniq)
FILES=$(shell find . -name "*.go" -type f | grep -v vendor)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
BINARY_NAME=$(shell basename ${PWD})
MAIN_PACKAGE="."

all: format init test vet build-dist

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
	@$(MAKE) init
	@gvt update --all

vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "--> Running go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $(FILES) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for reviewal."; \
	fi

.PHONY: all build build-dist format init test update-deps


