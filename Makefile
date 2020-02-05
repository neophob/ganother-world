-include .env

# source: https://github.com/azer/go-makefile-example/blob/master/Makefile

BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOFILES := $(wildcard *.go)
GOROOT := $(shell go env GOROOT)

# -X add string value definition of the form importpath.name=value
RELEASE := -ldflags "-s -w -X project.name=anotherworld"
SRC := main.go lib.go hal-sdl.go
SRCDIR := ./
PACKAGES := $(SRCDIR) anotherworld logger
DISTDIR := ./dist

## build: build all the things
build: build-native build-wasm

## build-native: build go binary in dev mode
build-native:
	@echo "  >  BUILD"
	@go build -o "$(DISTDIR)/main" $(SRC)

## build-wasm: builds the wasm app
build-wasm:
	@echo "  >  BUILD-WASM"
	@env GOARCH=wasm GOOS=js go build -o "$(DISTDIR)/lib.wasm" wasm/main.go
	@cp wasm/index.html $(DISTDIR)
	@go build -o "$(DISTDIR)/devserver" cmd/devserver/main.go
	@cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(DISTDIR)

## format: format code using go fmt
format:
	@go fmt $(SRC)

## build-release: build release build, could be compressed with UPX
build-release:
	#@env GOOS=js GOARCH=wasm go build -o gaw.js $(RELEASE) $(SRC)
	@env go build -o "$(DISTDIR)/main.release" $(RELEASE) $(SRC)

## test: run unit tests
test:
  # TODO test anotherworld lib too
	@go test -cover -v $(SRCDIR)

## doc: create project documentation
doc:
	@go doc -all $(SRCDIR)

## clean: removes build files
clean:
	@rm -r ./dist/*

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
