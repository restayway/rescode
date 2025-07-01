# Rescode - Type-Safe Go Error Code Generator

[![Go Reference](https://pkg.go.dev/badge/github.com/restayway/rescode.svg)](https://pkg.go.dev/github.com/restayway/rescode)
[![Go Report Card](https://goreportcard.com/badge/github.com/restayway/rescode)](https://goreportcard.com/report/github.com/restayway/rescode)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Rescode is a robust, MIT-licensed, open-source Go code generator that produces type-safe error code constants and error creator functions from user-supplied JSON or YAML configuration. The generated code follows Go best practices, works seamlessly with `go:generate`, and enables high-performance, compile-time validated error handling without any runtime map/slice lookups.

## 🎯 Motivation

Traditional error handling in Go often involves:
- Runtime map lookups for error metadata
- String-based error codes prone to typos
- Inconsistent error response formats
- No compile-time validation
- Poor performance due to runtime lookups

Rescode solves these problems by generating type-safe, compile-time constants and factory functions that provide:
- **Zero runtime lookups** - All metadata is compile-time constants
- **Type safety** - No magic strings or runtime errors
- **High performance** - Up to 140x faster than map-based approaches
- **Consistent API** - Standardized error handling across services
- **HTTP/gRPC ready** - Built-in status code mapping

## 🚀 Features

- ✅ **Type-safe error constants** - Generate compile-time validated error codes
- ✅ **Factory functions** - Simple, consistent error creation API
- ✅ **YAML/JSON input** - Flexible configuration format with auto-detection
- ✅ **HTTP/gRPC integration** - Built-in status code mapping
- ✅ **Error wrapping** - Preserve original errors with full stack traces
- ✅ **Additional data** - Attach structured data to errors
- ✅ **go:generate compatible** - Seamless integration with Go toolchain
- ✅ **High performance** - No runtime maps or lookups
- ✅ **Comprehensive tests** - >95% test coverage with benchmarks
- ✅ **CLI tool** - Simple command-line interface

## 📦 Installation

```bash
go get github.com/restayway/rescode
```

## 🛠️ Quick Start

### 1. Define your errors (YAML)

Create an `errors.yaml` file:

```yaml
- code: 1001
  key: UserNotFound
  message: User not found
  http: 404
  grpc: 5
  desc: The specified user could not be found in the database

- code: 1002
  key: InvalidEmail
  message: Invalid email address
  http: 400
  grpc: 3
  desc: The provided email address is not valid
```

### 2. Generate Go code

```bash
go run github.com/restayway/rescode/cmd/rescodegen --input errors.yaml --output errors_gen.go --package main
```

Or use `go:generate`:

```go
//go:generate go run github.com/restayway/rescode/cmd/rescodegen --input errors.yaml --output errors_gen.go --package main
```

### 3. Use the generated errors

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Create simple error
    err := UserNotFound()
    fmt.Printf("Error: %v (HTTP: %d, gRPC: %d)\n", err, err.HttpCode, err.RpcCode)

    // Create error with wrapped error
    originalErr := fmt.Errorf("database connection failed")
    wrappedErr := UserNotFound(originalErr)
    fmt.Printf("Wrapped: %v\n", wrappedErr)
    fmt.Printf("Original: %v\n", wrappedErr.OriginalError())

    // Add additional data
    enrichedErr := InvalidEmail().SetData(map[string]string{
        "field": "email",
        "value": "invalid@",
    })
    fmt.Printf("JSON: %v\n", enrichedErr.JSON())

    // Use constants for logic
    if err.Code == UserNotFoundCode {
        log.Printf("Handling user not found error with code %d", UserNotFoundCode)
    }
}
```

## 📋 Error Definition Schema

### YAML Format

```yaml
- code: 1001              # Required: Unique numeric error code (uint64)
  key: UserNotFound       # Required: Go identifier for the error
  message: User not found # Required: Human-readable error message
  http: 404              # Required: HTTP status code
  grpc: 5                # Required: gRPC status code (0-16)
  desc: Description      # Optional: Detailed description for documentation
```

### JSON Format

```json
[
  {
    "code": 1001,
    "key": "UserNotFound",
    "message": "User not found",
    "http": 404,
    "grpc": 5,
    "desc": "The specified user could not be found in the database"
  }
]
```

### Field Validation

- **code**: Must be non-zero unique uint64
- **key**: Must be valid Go identifier (PascalCase recommended)
- **message**: Non-empty human-readable string
- **http**: Valid HTTP status code (typically 400-599)
- **grpc**: Valid gRPC status code (0-16)
- **desc**: Optional description for documentation

### gRPC Status Code Reference

| Code | gRPC Status | Description |
|------|-------------|-------------|
| 0 | OK | Success |
| 1 | CANCELLED | Operation cancelled |
| 2 | UNKNOWN | Unknown error |
| 3 | INVALID_ARGUMENT | Invalid argument |
| 4 | DEADLINE_EXCEEDED | Deadline exceeded |
| 5 | NOT_FOUND | Not found |
| 6 | ALREADY_EXISTS | Already exists |
| 7 | PERMISSION_DENIED | Permission denied |
| 8 | RESOURCE_EXHAUSTED | Resource exhausted |
| 9 | FAILED_PRECONDITION | Failed precondition |
| 10 | ABORTED | Aborted |
| 11 | OUT_OF_RANGE | Out of range |
| 12 | UNIMPLEMENTED | Unimplemented |
| 13 | INTERNAL | Internal error |
| 14 | UNAVAILABLE | Unavailable |
| 15 | DATA_LOSS | Data loss |
| 16 | UNAUTHENTICATED | Unauthenticated |

## 🎯 Generated Code

For the example above, rescode generates:

```go
// Code generated by rescodegen. DO NOT EDIT.

package main

import (
    "github.com/restayway/rescode"
    "google.golang.org/grpc/codes"
)

// Error code constants
const (
    UserNotFoundCode uint64     = 1001
    UserNotFoundHTTP int        = 404
    UserNotFoundGRPC codes.Code = 5
    UserNotFoundMsg  string     = "User not found"
    UserNotFoundDesc string     = "The specified user could not be found in the database"

    InvalidEmailCode uint64     = 1002
    InvalidEmailHTTP int        = 400
    InvalidEmailGRPC codes.Code = 3
    InvalidEmailMsg  string     = "Invalid email address"
    InvalidEmailDesc string     = "The provided email address is not valid"
)

// UserNotFound creates a new UserNotFound error.
// The specified user could not be found in the database
func UserNotFound(err ...error) *rescode.RC {
    return rescode.New(UserNotFoundCode, UserNotFoundHTTP, UserNotFoundGRPC, UserNotFoundMsg)(err...)
}

// InvalidEmail creates a new InvalidEmail error.
// The provided email address is not valid
func InvalidEmail(err ...error) *rescode.RC {
    return rescode.New(InvalidEmailCode, InvalidEmailHTTP, InvalidEmailGRPC, InvalidEmailMsg)(err...)
}
```

## 🏃‍♂️ CLI Usage

```bash
rescodegen [OPTIONS]

Options:
  --input     Path to YAML/JSON file containing error definitions (required)
  --output    Path to generated Go file (default: rescode_gen.go)
  --package   Go package name to use in generated code (default: directory name)
  --version   Show version information
  --help      Show help information

Examples:
  rescodegen --input errors.yaml --output errors_gen.go --package myservice
  go run github.com/restayway/rescode/cmd/rescodegen --input errors.json

For go:generate usage:
  //go:generate go run github.com/restayway/rescode/cmd/rescodegen --input errors.yaml --output errors_gen.go --package myservice
```

## 🔧 API Reference

### Core Types

```go
type RC struct {
    Code     uint64     // Unique error code
    Message  string     // Human-readable error message
    HttpCode int        // HTTP status code
    RpcCode  codes.Code // gRPC status code
    Data     any        // Optional additional data
}

type RcCreator func(...error) *RC
```

### Core Functions

```go
// New creates an RcCreator function with the specified parameters
func New(code uint64, hCode int, rCode codes.Code, message string, data ...any) RcCreator

// Error implements the error interface
func (r *RC) Error() string

// SetData sets additional data for the error and returns the RC for chaining
func (r *RC) SetData(data any) *RC

// JSON returns a map representation of the error, optionally filtering by keys
func (r *RC) JSON(keys ...string) map[string]interface{}

// OriginalError returns the wrapped original error, if any
func (r *RC) OriginalError() error

// String returns a string representation of the error
func (r *RC) String() string
```

## 📊 Performance Benchmarks

Rescode dramatically outperforms traditional runtime map-based error handling:

```
BenchmarkGenerated_PolicyNotFound-2         1000000000    0.31 ns/op    0 B/op    0 allocs/op
BenchmarkLegacy_PolicyNotFound-2                26802039   45.16 ns/op   80 B/op    1 allocs/op

BenchmarkGenerated_MultipleErrors-2         1000000000    0.94 ns/op    0 B/op    0 allocs/op
BenchmarkLegacy_MultipleErrors-2                25338333   45.99 ns/op   80 B/op    1 allocs/op
```

**Results Summary:**
- **145x faster** error creation
- **Zero allocations** vs 1 allocation per error
- **Zero memory usage** vs 80 bytes per error
- **Consistent performance** regardless of error count

The performance improvement comes from:
1. **Compile-time constants** - No runtime lookups
2. **Pre-allocated creators** - Functions generated at compile time
3. **Zero reflection** - Direct struct initialization
4. **No map access** - All metadata embedded in code

## 🏗️ Examples

### Basic Usage

See [examples/basic/](examples/basic/) for a complete basic example.

### Microservice Integration

See [examples/microservice/](examples/microservice/) for a full HTTP service example with error handling.

Key features demonstrated:
- HTTP status code mapping
- JSON error responses
- Error data attachment
- Error wrapping
- Type-safe error handling

## 🧪 Testing

Run tests with coverage:

```bash
go test -v -cover ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

The project maintains >95% test coverage with comprehensive tests for:
- Core error functionality
- Code generation
- CLI tool
- Generated code validation
- Performance benchmarks

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Clone the repository
2. Run tests: `go test ./...`
3. Run benchmarks: `go test -bench=.`
4. Generate examples: `cd examples/basic && go generate`

### Code Standards

- Follow Go best practices and idioms
- Maintain test coverage >95%
- Include benchmarks for performance-critical code
- Update documentation for API changes

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for type-safe, high-performance error handling in Go
- Built with Go's excellent tooling ecosystem
- Thanks to the Go community for feedback and suggestions

---

**Made with ❤️ for the Go community**