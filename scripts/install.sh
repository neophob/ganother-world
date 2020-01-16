#!/bin/bash

echo "Install SDL2 with brew:"
brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config

echo "Installing go packages:"
go get -v github.com/veandco/go-sdl2/sdl
go get -v github.com/veandco/go-sdl2/gfx
