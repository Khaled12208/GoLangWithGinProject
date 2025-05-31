.PHONY: test test-unit test-integration test-coverage test-report clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=golangwithgin

# Test parameters
COVERAGE_FILE=coverage.out
JUNIT_REPORT=test-report.xml

all: test

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	@cd internal/domain/mocks && \
	mockgen -destination=task_repository_mock.go -package=mocks golangwithgin/internal/domain TaskRepository && \
	mockgen -destination=task_processor_mock.go -package=mocks golangwithgin/internal/domain TaskProcessor && \
	mockgen -destination=task_service_mock.go -package=mocks golangwithgin/internal/domain TaskService

# Run unit tests
test-unit: generate-mocks
	@echo "Running unit tests..."
	$(GOTEST) -v ./internal/... -coverprofile=$(COVERAGE_FILE) -json | tee test-output.json
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@cat test-output.json | go-junit-report > $(JUNIT_REPORT)

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v ./tests/integration/... -tags=integration -coverprofile=integration-$(COVERAGE_FILE) -json | tee integration-test-output.json
	@go tool cover -html=integration-$(COVERAGE_FILE) -o integration-coverage.html
	@cat integration-test-output.json | go-junit-report > integration-$(JUNIT_REPORT)

# Run all tests
test: test-unit test-integration
	@echo "All tests completed"

# Generate test coverage report
test-coverage: test
	@echo "Generating test coverage report..."
	@go tool cover -func=$(COVERAGE_FILE)
	@go tool cover -func=integration-$(COVERAGE_FILE)

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	@rm -f $(COVERAGE_FILE) coverage.html test-output.json $(JUNIT_REPORT)
	@rm -f integration-$(COVERAGE_FILE) integration-coverage.html integration-test-output.json integration-$(JUNIT_REPORT)

# Install test dependencies
install-test-deps:
	@echo "Installing test dependencies..."
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/jstemmer/go-junit-report@latest 