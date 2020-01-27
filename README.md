# ganother-world

GOAL:
- Learn GO and implement "something fun"
- "something fun" - interpret the old game another world and create a web version!

Reference:
- https://github.com/cyxx/rawgl
- https://fabiensanglard.net/anotherWorld_code_review/
- https://fabiensanglard.net/another_world_polygons/index.html
- http://www.anotherworld.fr/anotherworld_uk/another_world.htm
- https://www.gdcvault.com/play/1014630/Classic-Game-Postmortem-OUT-OF

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
- Audio / Sound FX (Protracker module)

# GO

- https://awesome-go.com/ - A curated list of awesome Go frameworks, libraries and software
- https://github.com/golang-standards/project-layout
- GO111MODULE and modules?
- https://vincent.bernat.ch/en/blog/2019-makefile-build-golang
- https://github.com/golang/go/wiki/WebAssembly
- https://golang.org/doc/effective_go.html
- https://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/

## lang elements to check

- iota
- Channels
- range clause
- defer
- interface
- Method (bound functions?)

# Getting started

- Install make (autotools on linux, xcode on OSX)
- Install go 1.13 https://golang.org/dl/
- Check out repo
- Run `./scripts/osx-install.sh` to install/download dependencies or make sure SDL2 (sdl2, sdl2_gfx, sdl2_image, sdl2_mixer, sdl_net) and pkg-config are installed correctly
- Run `make` to build, if it's green you're good
- Use `make help` for more
