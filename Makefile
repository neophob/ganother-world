-include .env

# source: https://github.com/azer/go-makefile-example/blob/master/Makefile

BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOFILES := $(wildcard *.go)

# -X add string value definition of the form importpath.name=value
RELEASE := -ldflags "-s -w -X project.name=anotherworld"
SRC := src/main.go src/resource.go src/vm.go src/parts.go src/decrunch.go \
	src/vm-ops.go src/assets.go src/video.go src/video-dummy.go src/video-sdl.go src/text.go \
	src/font.go src/log.go src/videoassets.go src/polygon.go
SRCDIR := ./src

## build: build go binary in dev mode
build:
	@echo "  >  BUILD"
	@go build $(SRC)

## format: format code using go fmt
format:
	@go fmt $(SRC)

## build-release: build release build, could be compressed with UPX
build-release:
	#@env GOOS=js GOARCH=wasm go build -o gaw.js $(RELEASE) $(SRC)
	@env go build -o main.release $(RELEASE) $(SRC)

## test: run unit tests
test:
	@go test -cover -v $(SRCDIR)

## doc: create project documentation
doc:
	@go doc -all $(SRCDIR)

## clean: removes build files
clean:
	@rm -f ./main
	@rm -f ./main.release

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
