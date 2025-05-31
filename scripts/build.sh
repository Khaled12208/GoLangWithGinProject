#!/bin/sh

# Clean previous builds
rm -f main

# Build for the current platform
go build -o main ./cmd/api

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o main-linux-amd64 ./cmd/api
GOOS=darwin GOARCH=amd64 go build -o main-darwin-amd64 ./cmd/api
GOOS=windows GOARCH=amd64 go build -o main-windows-amd64.exe ./cmd/api 