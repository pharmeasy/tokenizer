NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# The import path is the unique absolute name of your repository.
# All subpackages should always be imported as relative to it.
# If you change this, run `make clean`.
PKG_SRC := bitbucket.org/pharmaeasyteam/tokenizer
VERSION := `cat VERSION`

.PHONY: all clean deps build

all: format clean deps test build docker

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/DATA-DOG/godog/cmd/godog

build: format
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "TOKENIZER_BUILD_ONLY_DEFAULT=$(TOKENIZER_BUILD_ONLY_DEFAULT) PKG_SRC=$(PKG_SRC) VERSION=$(VERSION) ./build.sh"

docker: all
	@docker build . --tag pe/tokenizer:$(VERSION)

run: build
	@build/tokenizer "start"

cleanrun: all
	@build/tokenizer "start"

test: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -cover ./...

test-integration: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -cover -tags=integration ./...

test-features:
	@/bin/sh -c "TOKENIZER_BUILD_ONLY_DEFAULT=1 PKG_SRC=$(PKG_SRC) ./build.sh"
	@/bin/sh -c "./features.sh"

format:
	@echo "$(OK_COLOR)==> checking code formating with 'gofmt' tool$(NO_COLOR)"
	@gofmt -l -s ./ | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

fixformat:
	@go fmt ./...

vet:
	@echo "$(OK_COLOR)==> checking code correctness with 'go vet' tool$(NO_COLOR)"
	@go vet ./...

lint: tools.golint
	@echo "$(OK_COLOR)==> checking code style with 'golint' tool$(NO_COLOR)"
	@go list ./... | xargs -n 1 golint -set_exit_status

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
	@rm -rf bin $GOPATH/bin

devwatchbuild:
	@/bin/sh -c "TOKENIZER_BUILD_ONLY_DEFAULT=1 PKG_SRC=$(PKG_SRC) ./watch_build.sh"

#---------------
#-- tools
#---------------

.PHONY: tools tools.golint
tools: tools.golint

tools.golint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		go get github.com/golang/lint/golint; \
	fi