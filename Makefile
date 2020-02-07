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
WASMDIR := ./wasm
SDLDIR := ./sdl
PACKAGES := $(SDLDIR) ./logger ./anotherworld
PACKAGES_TO_TEST := ./logger ./anotherworld
DISTDIR := ./dist

## build: build all the things
build: build-native build-wasm

## build-native: build go SDL binary
build-native:
	@echo "  >  BUILD"
	@go build -o "$(DISTDIR)/main" $(SDLDIR)

## build-wasm: builds the wasm app
build-wasm:
	@echo "  >  BUILD-WASM"
	@env GOARCH=wasm GOOS=js go build -o "$(DISTDIR)/lib.wasm" $(WASMDIR)/main.go
	@cp wasm/index.html $(DISTDIR)
	@go build -o "$(DISTDIR)/devserver" cmd/devserver/main.go
	@cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(DISTDIR)

## build-release: build release build, could be compressed with UPX
build-release:
	@env go build -o "$(DISTDIR)/main.release" $(RELEASE) $(SDLDIR)

## format: format code using go fmt
format:
	@go fmt $(PACKAGES)

## test: run unit tests
test:
	@go test -cover -v $(PACKAGES_TO_TEST)

## doc: create project documentation
doc:
	@go doc -all $(SDLDIR)
	@go doc -all $(WASMDIR)
	@go doc -all ./logger
	@go doc -all ./anotherworld

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
