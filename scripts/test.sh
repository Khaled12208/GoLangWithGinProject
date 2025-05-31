#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo -e "\n${GREEN}=== $1 ===${NC}\n"
}

# Function to run tests and check result
run_test() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $1 passed${NC}"
    else
        echo -e "${RED}✗ $1 failed${NC}"
        exit 1
    fi
}

# Clean test cache
print_header "Cleaning test cache"
go clean -testcache

# Run unit tests
print_header "Running unit tests"
go test ./... -v -short
run_test "Unit tests"

# Run integration tests
print_header "Running integration tests"
go test ./... -v -run Integration
run_test "Integration tests"

# Generate test coverage
print_header "Generating test coverage"
mkdir -p coverage
go test ./... -coverprofile=coverage/coverage.out
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
run_test "Coverage generation"

# Display coverage statistics
print_header "Coverage Statistics"
go tool cover -func=coverage/coverage.out

print_header "Testing completed successfully!" 