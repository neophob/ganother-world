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
PACKAGES := $(SDLDIR) $(WASMDIR) ./logger ./anotherworld
PACKAGES_TO_TEST := ./anotherworld
DISTDIR := ./dist

## build: build all the things (size unoptimized build)
build: build-native build-wasm

## release: release build (size optimized build)
release: build-native-release build-wasm

## build-native: build go SDL2 binary
build-native:
	@echo "  >  BUILD SDL version"
	@go build -o "$(DISTDIR)/main" $(SDLDIR)
	@echo "  DONE! run main in the dist directory"

build-native-release:
	@echo "  >  BUILD SDL release version"
	@env go build -o "$(DISTDIR)/main.release" $(RELEASE) $(SDLDIR)
	@echo "  DONE! run main.release in the dist directory"

wasm-common:
	@mkdir -p $(DISTDIR)
	@cp -f wasm/index.html $(DISTDIR)
	@cp -f wasm/index.js $(DISTDIR)
	@go build -o "$(DISTDIR)/devserver" $(RELEASE) cmd/devserver/main.go
	@cp -f "$(GOROOT)/misc/wasm/wasm_exec.js" $(DISTDIR)
	@cp -fr ./assets $(DISTDIR)
	@cp -fr ./logo.png $(DISTDIR)

## build-wasm: builds the wasm app
build-wasm: wasm-common
	@echo "  >  BUILD WASM version"
	@env GOARCH=wasm GOOS=js go build -o "$(DISTDIR)/lib.wasm" $(WASMDIR)
	@echo "  DONE! run devserver in the dist directory"

## format: format code using go fmt
format:
	@gofmt -w .

## test: run unit tests
test:
	@go test -cover -v $(PACKAGES_TO_TEST)

## lint: static analyze source
lint:
	@env GOARCH=wasm GOOS=js go vet ./wasm/...
	@go vet ./sdl/...

## doc: create project documentation
doc:
	@go doc -all $(SDLDIR)
	@go doc -all $(WASMDIR)
	@go doc -all ./logger
	@go doc -all ./anotherworld

## install-go-deps: install go dependencies
install-go-deps:
	@go get -v ./...

## update-go-deps: update go dependencies
update-go-deps:
	@go get -t -v -d -u ./...

## clean: removes build files
clean:
	@go clean
	@rm -fr ./dist/*

all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: build release build-native build-native-release wasm-common build-wasm format test lint doc update-go-deps clean help
