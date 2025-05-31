#!/bin/sh

# Build the application
go build -o main ./cmd/api

# Run the application
./main 