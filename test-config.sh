#!/bin/bash

# Config Test Runner for micro-ddns
# This script runs comprehensive tests for the configuration module

set -e

echo "Starting micro-ddns configuration tests..."
echo "=========================================="

# Change to project root
cd "$(dirname "$0")"

# Clean up any previous test artifacts
echo "Cleaning up previous test artifacts..."
rm -f coverage.out coverage.html

# Run the tests with verbose output
echo "Running configuration tests..."
go test ./internal/config -v

# Generate coverage report
echo "Generating coverage report..."
go test ./internal/config -coverprofile=coverage.out
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')

echo "=========================================="
echo "Test Results:"
echo "- All tests: PASSED"
echo "- Coverage: $COVERAGE"
echo "=========================================="

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
echo "HTML coverage report generated: coverage.html"

# Clean up coverage.out (keep HTML report)
rm -f coverage.out

echo "Testing completed successfully!"
