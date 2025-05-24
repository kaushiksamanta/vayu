# Go binary names
GO := go

# Project variables
BINARY_NAME := vayu
MAIN_DIR := example

# Test directories
TEST_DIR := tests
UNIT_TEST_DIR := $(TEST_DIR)/unit
INTEGRATION_TEST_DIR := $(TEST_DIR)/integration

# Tools
GOFMT := gofmt

.PHONY: build test test-unit test-integration test-all test-silent test-verbose-logs lint fmt vet staticcheck run clean

build:
	$(GO) build -o $(BINARY_NAME) .

# Test targets
test:
	$(GO) test ./...

test-unit:
	$(GO) test ./$(UNIT_TEST_DIR)/...

test-integration:
	$(GO) test ./$(INTEGRATION_TEST_DIR)/...

test-all: test-unit test-integration

# Test with verbose output
test-v:
	$(GO) test -v ./...

test-unit-v:
	$(GO) test -v ./$(UNIT_TEST_DIR)/...

test-integration-v:
	$(GO) test -v ./$(INTEGRATION_TEST_DIR)/...

# Test with code coverage
test-cover:
	$(GO) test -cover ./...

test-cover-html:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Generated coverage report: coverage.html"

# Run benchmarks
bench:
	$(GO) test -bench=. -benchmem ./...

bench-unit:
	$(GO) test -bench=. -benchmem ./$(UNIT_TEST_DIR)/...

# Run tests in race detection mode
test-race:
	$(GO) test -race ./...

# Run tests with SilentMode explicitly enabled (suppress panic logs)
test-silent:
	$(GO) test -ldflags="-X 'github.com/kaushiksamanta/vayu.SilentMode=true'" ./...

# Run tests with SilentMode explicitly disabled (show all panic logs)
test-verbose-logs:
	$(GO) test -ldflags="-X 'github.com/kaushiksamanta/vayu.SilentMode=false'" ./...

# Format code
fmt:
	$(GOFMT) -w -s .

# Run go vet
vet:
	$(GO) vet ./...

# Run staticcheck
staticcheck:
	$(GO) run honnef.co/go/tools/cmd/staticcheck@latest ./...

# Combined lint target
lint: fmt vet staticcheck

# Check with both tests and linting
check: test lint

run:
	$(GO) run $(MAIN_DIR)/main.go

clean:
	rm -f $(BINARY_NAME)
	rm -rf bin/
