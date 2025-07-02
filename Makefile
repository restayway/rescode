.PHONY: help test bench bench-all bench-compare clean build

# Default target
help:
	@echo "Available commands:"
	@echo "  test         - Run all tests"
	@echo "  bench        - Run all benchmarks"
	@echo "  bench-all    - Run all benchmarks with detailed output"
	@echo "  bench-compare- Run benchmarks with comparison between approaches"
	@echo "  bench-summary- Show clean performance comparison table"
	@echo "  bench-table  - Show all benchmarks in table format"
	@echo "  mod-tidy     - Run go mod tidy while preserving Go 1.20"
	@echo "  mod-verify   - Verify Go version is correctly set to 1.20"
	@echo "  mod-reset    - Reset go.mod to proper Go 1.20 state"
	@echo "  build        - Build the project"
	@echo "  clean        - Clean build artifacts"

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	@echo "Coverage report generated: coverage.out"

	@echo "To view the coverage report, run:"
	@echo "make view-coverage"

# View coverage report in browser
view-coverage:
	@echo "Opening coverage report in browser..."
	@go tool cover -html=coverage.out -o coverage.html
	@open coverage.html || xdg-open coverage.html || start coverage.html

# Run benchmarks
bench:
	go test -bench=. -benchmem

# Run benchmarks with more detailed output
bench-all:
	@echo "=== Running All Benchmarks ==="
	@echo ""
	@echo "Legend:"
	@echo "- Generated: Current optimized approach with compile-time constants"
	@echo "- Legacy: Runtime map-based error registry"
	@echo "- StaticCode: uint64 constants with message maps"
	@echo "- VarCode: var declarations with message maps"
	@echo ""
	go test -bench=. -benchmem -count=3

# Run benchmarks with comparison focus
bench-compare:
	@echo "=== Benchmark Comparison: Different Approaches ==="
	@echo ""
	@echo "Comparing 4 different error handling approaches:"
	@echo "1. Generated (current): Compile-time optimized"
	@echo "2. Legacy: Runtime map lookup"
	@echo "3. StaticCode: Static constants + maps"
	@echo "4. VarCode: Variable declarations + maps"
	@echo ""
	@echo "Performance Metrics:"
	@echo "- ns/op: nanoseconds per operation (lower is better)"
	@echo "- B/op: bytes allocated per operation (lower is better)"
	@echo "- allocs/op: allocations per operation (lower is better)"
	@echo ""
	go test -bench="PolicyNotFound$$|MultipleErrors$$|ErrorMessage$$|JSON$$" -benchmem -count=5

# Build the project
build:
	go build ./...

# Build command line tool
build-tool:
	go build -o bin/rescodegen ./cmd/rescodegen

# Clean build artifacts
clean:
	go clean ./...
	rm -rf bin/

# Run examples
run-example-basic:
	cd examples/basic && go run main.go

run-example-microservice:
	cd examples/microservice && go run service.go

# Generate code for examples
generate-basic:
	cd examples/basic && go run ../../cmd/rescodegen/main.go -input errors.yaml -output errors_gen.go

generate-microservice:
	cd examples/microservice && go run ../../cmd/rescodegen/main.go -input errors.json -output service_errors.go

# Development helpers
mod-tidy:
	@echo "Running go mod tidy while preserving Go 1.20..."
	@if grep -q "go 1\.20" go.mod; then \
		go mod tidy; \
		if ! grep -q "go 1\.20" go.mod; then \
			echo "WARNING: Go version was changed! Restoring Go 1.20..."; \
			sed -i '' 's/go [0-9]\+\.[0-9]\+.*/go 1.20/' go.mod; \
		fi; \
	else \
		echo "ERROR: go.mod doesn't contain 'go 1.20'. Please check manually."; \
		exit 1; \
	fi
	@echo "✓ Go version preserved at 1.20"

mod-verify:
	@echo "Verifying Go version in go.mod..."
	@if grep -q "go 1\.20" go.mod; then \
		echo "✓ Go version is correctly set to 1.20"; \
	else \
		echo "✗ Go version is not 1.20. Current version:"; \
		grep "^go " go.mod; \
		exit 1; \
	fi

# Reset go.mod to proper state if it gets corrupted
mod-reset:
	@echo "Resetting go.mod to Go 1.20 with compatible dependencies..."
	@rm -f go.mod go.sum
	@go mod init github.com/restayway/rescode
	@sed -i '' 's/go [0-9]\+\.[0-9]\+.*/go 1.20/' go.mod
	@go get google.golang.org/grpc@v1.56.3
	@go get gopkg.in/yaml.v3@v3.0.1
	@go mod edit -require=google.golang.org/grpc@v1.56.3 -require=gopkg.in/yaml.v3@v3.0.1
	@go mod tidy
	@make mod-verify
	@echo "✓ go.mod reset successfully"

# Full development check
check: fmt vet test bench

# Show benchmark results in a table format
bench-table:
	@echo "Running benchmarks and formatting results..."
	@go test -bench=. -benchmem | grep -E "(Benchmark|PASS)" | grep -v "PASS" | \
	awk 'BEGIN {printf "%-40s %12s %8s %12s\n", "Benchmark", "ns/op", "B/op", "allocs/op"; print "========================================================================="} \
	{printf "%-40s %12s %8s %12s\n", $$1, $$3, $$5, $$7}'

# Summary comparison of all approaches
bench-summary:
	@echo "=== Performance Comparison Summary ==="
	@echo ""
	@echo "Running benchmarks for PolicyNotFound scenario..."
	@go test -bench="PolicyNotFound$$" -benchmem -count=1 | grep "Benchmark" | \
	awk 'BEGIN {printf "%-35s %12s %10s %12s\n", "Approach", "ns/op", "B/op", "allocs/op"; print "==============================================================================="} \
	/BenchmarkGenerated_PolicyNotFound/ {printf "%-35s %12s %10s %12s\n", "Generated (Optimized)", $$3, $$5, $$7} \
	/BenchmarkLegacy_PolicyNotFound/ {printf "%-35s %12s %10s %12s\n", "Legacy (Map Lookup)", $$3, $$5, $$7} \
	/BenchmarkStaticCode_PolicyNotFound/ {printf "%-35s %12s %10s %12s\n", "Static Constants + Maps", $$3, $$5, $$7} \
	/BenchmarkVarCode_PolicyNotFound/ {printf "%-35s %12s %10s %12s\n", "Var Declarations + Maps", $$3, $$5, $$7}'
	@echo ""
	@echo "Key Insights:"
	@echo "- Generated approach: ~80x faster, zero allocations"
	@echo "- All map-based approaches: Similar performance, 1 allocation per call"
	@echo "- Lower ns/op = faster, Lower B/op = less memory, Lower allocs/op = better"
