.PHONY: run build test vet lint test-coverage clean

# Run the application
run:
	go run ./cmd/web

# Build the binary
build:
	go build ./cmd/web

# Run all tests
test:
	go test ./...

# Run go vet
vet:
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Run tests with coverage, excluding infrastructure/wiring packages
COVERAGE_EXCLUDE := cmd/web|internal/view|internal/config|internal/db|internal/http/routes|internal/http/handlers/api|internal/model

test-coverage:
	@go test ./... -coverprofile=coverage.out
	@grep -v -E "$(COVERAGE_EXCLUDE)" coverage.out > coverage_filtered.out
	@go tool cover -func=coverage_filtered.out
	@rm -f coverage_filtered.out

# Remove coverage artifacts
clean:
	@rm -f coverage.out coverage_filtered.out
