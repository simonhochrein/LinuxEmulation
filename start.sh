#!/bin/sh
export GOPATH=$GOPATH:$(pwd)
go run src/main.go
# gomobile bind -target ios -o ./build/VM.framework vm