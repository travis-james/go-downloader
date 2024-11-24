# Variables
PROJECT_NAME := go_downloader
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Default target
.PHONY: all
all: test coverage

# Run tests with coverage
.PHONY: test
test:
	go test -v -coverprofile=$(COVERAGE_FILE) ./...

# Show coverage in the terminal, need to run 'test' first.
.PHONY: coverage
coverage:
	go tool cover -func=$(COVERAGE_FILE)

# Generate an HTML report of the coverage
.PHONY: coverage-html
coverage-html:
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Clean up coverage files
.PHONY: clean
clean:
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

# Notes for myself:
# 'PHONY' avoids name conflicts. If there was a file the same
# name as the make command, that would cause issues. It tells
# make to run the command and NOT the file of the same name.