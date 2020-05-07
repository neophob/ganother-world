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
PACKAGES_TO_TEST := ./anotherworld
DISTDIR := ./dist

## build: build all the things
build: build-native build-wasm

## release: build release build, SDL version could be compressed with UPX
release: build-native-release build-wasm-release

## build-native: build go SDL binary
build-native:
	@echo "  >  BUILD SDL version"
	@go build -o "$(DISTDIR)/main" $(SDLDIR)
	@echo "  DONE! run main in the dist directory"

build-native-release:
	@echo "  >  BUILD SDL release version"
	@env go build -o "$(DISTDIR)/main.release" $(RELEASE) $(SDLDIR)
	@echo "  DONE! run main.release in the dist directory"

wasm-common:
	@cp wasm/index.html $(DISTDIR)
	@cp wasm/main.js $(DISTDIR)
	@go build -o "$(DISTDIR)/devserver" cmd/devserver/main.go
	@cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(DISTDIR)
	@cp -r ./assets $(DISTDIR)
	@cp -r ./logo.png $(DISTDIR)

## build-wasm: builds the wasm app
build-wasm: wasm-common
	@echo "  >  BUILD WASM version"
	@env GOARCH=wasm GOOS=js go build -o "$(DISTDIR)/lib.wasm" $(WASMDIR)
	@echo "  DONE! run devserver in the dist directory"

## build-wasm: builds the wasm app
build-wasm-release: wasm-common
	@echo "  >  BUILD WASM release version"
	@env GOARCH=wasm GOOS=js go build -o "$(DISTDIR)/lib.wasm" $(RELEASE) $(WASMDIR)
	@echo "  DONE! run devserver in the dist directory"

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
	@go clean
	@rm -fr ./dist/*

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
