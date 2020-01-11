-include .env

# source: https://github.com/azer/go-makefile-example/blob/master/Makefile

BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# -X add string value definition of the form importpath.name=value
RELEASE := -ldflags "-s -w -X gaw.hello=world"
SRC := src/main.go src/resource.go src/vm.go src/parts.go src/decrunch.go
SRCDIR := ./src

## build: build go binary in dev mode
build:
	@echo "  >  BUILD"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(SRC)

## format: format code using go fmt
format:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go fmt $(SRC)

## build-cross: cross compile project in release mode (without debug symbols)
build-cross:
	@env GOOS=js GOARCH=wasm go build -o gaw.js $(RELEASE) $(SRC)
	@env GOOS=linux GOARCH=arm GOARM=7 go build -o gaw.lnx $(RELEASE) $(SRC)
	@env GOOS=windows GOARCH=amd64 go build -o gaw.win $(RELEASE) $(SRC)
	@env GOOS=darwin GOARCH=amd64 go build -o gaw.osx $(RELEASE) $(SRC)
	#env GOOS=android GOARCH=arm64 go build -o gaw.and

## test: run unit tests
test:
	@go test -cover -v $(SRCDIR)

## doc: create project documentation
doc:
	@go doc -all $(SRCDIR)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
