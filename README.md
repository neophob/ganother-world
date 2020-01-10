# ganother-world

Reference: https://fabiensanglard.net/anotherWorld_code_review/ and http://www.anotherworld.fr/anotherworld_uk/another_world.htm

- GOAL: Learn GO and implement "something fun"
- "something fun" - interpret the old game another world and create a web version!

# ARCHITECTURE

- Split Backend (parsing, interpreting) and Frontend (rendering audio & video).
- Backend should build as lib so it can be reusable in a web app (WASM target) or SDL2 app
- Local assets are here for dev purpose only - they will go away. thread them as they will be downloaded

# EPICS

- Basic Project Bootstrap
- Design GO library
- Resource Management: read MEMLIST.BIN, read resources
- Setup Video and Audio
- VM parser
- implement VM op codes (Polygons drawing)
- Audio / Sound FX

# GO

- https://awesome-go.com/ - A curated list of awesome Go frameworks, libraries and software
- https://github.com/golang-standards/project-layout
- GO111MODULE and modules?
- https://vincent.bernat.ch/en/blog/2019-makefile-build-golang
- https://github.com/golang/go/wiki/WebAssembly
- https://golang.org/doc/effective_go.html

## lang elements to check
- iota
- Channels
- range clause
- defer
- interface
- Method (bound functions?)

# Getting started

- Install make (`brew install make`?)
- Install go 1.13 https://golang.org/dl/
- Check out repo
- Run `make` to build, if it's green you're good
- Use `make help` for more
