#!/bin/sh
export GOPATH=$GOPATH:$(pwd)
go run src/main.go $1
# gomobile bind -target ios -o ./build/VM.framework vm