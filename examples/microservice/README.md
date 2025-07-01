# Microservice Error Definitions Example

This example shows how to use rescode in a microservice context with JSON format and different error categories.

## errors.json

```json
[
  {
    "code": 10001,
    "key": "AuthenticationFailed",
    "message": "Authentication failed",
    "http": 401,
    "grpc": 16,
    "desc": "The provided credentials are invalid"
  },
  {
    "code": 10002,
    "key": "AuthorizationDenied",
    "message": "Authorization denied",
    "http": 403,
    "grpc": 7,
    "desc": "User does not have permission to perform this action"
  },
  {
    "code": 20001,
    "key": "PolicyNotFound",
    "message": "Policy not found",
    "http": 404,
    "grpc": 5,
    "desc": "Policy could not be located in the database"
  },
  {
    "code": 20002,
    "key": "InvalidPolicyKind",
    "message": "Invalid policy kind",
    "http": 400,
    "grpc": 3,
    "desc": "Policy kind is not supported"
  },
  {
    "code": 30001,
    "key": "RateLimitExceeded",
    "message": "Rate limit exceeded",
    "http": 429,
    "grpc": 8,
    "desc": "Request rate limit has been exceeded"
  },
  {
    "code": 99999,
    "key": "InternalServerError",
    "message": "Internal server error",
    "http": 500,
    "grpc": 13,
    "desc": "An unexpected internal error occurred"
  }
]
```

## service.go

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//go:generate go run github.com/restayway/rescode/cmd/rescodegen --input errors.json --output service_errors.go --package main

type PolicyService struct{}

type Policy struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Name string `json:"name"`
}

func (s *PolicyService) GetPolicy(id string) (*Policy, error) {
	// Simulate authentication check
	if id == "unauthorized" {
		return nil, AuthenticationFailed()
	}
	
	// Simulate authorization check
	if id == "forbidden" {
		return nil, AuthorizationDenied()
	}
	
	// Simulate policy not found
	if id == "notfound" {
		return nil, PolicyNotFound()
	}
	
	// Simulate invalid kind
	if id == "invalid" {
		originalErr := fmt.Errorf("unsupported kind: foobar")
		return nil, InvalidPolicyKind(originalErr)
	}
	
	// Return a valid policy
	return &Policy{
		ID:   id,
		Kind: "standard",
		Name: "Sample Policy",
	}, nil
}

func handleError(w http.ResponseWriter, err error) {
	if rcErr, ok := err.(*rescode.RC); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(rcErr.HttpCode)
		
		response := rcErr.JSON("code", "message", "data")
		json.NewEncoder(w).Encode(response)
	} else {
		// Fallback for unknown errors
		internalErr := InternalServerError(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(internalErr.HttpCode)
		
		response := internalErr.JSON("code", "message")
		json.NewEncoder(w).Encode(response)
	}
}

func policyHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		err := InvalidPolicyKind().SetData(map[string]string{
			"field": "id",
			"reason": "missing required parameter",
		})
		handleError(w, err)
		return
	}
	
	service := &PolicyService{}
	policy, err := service.GetPolicy(id)
	if err != nil {
		handleError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

func main() {
	http.HandleFunc("/policy", policyHandler)
	
	fmt.Println("Starting server on :8080")
	fmt.Println("Try these URLs:")
	fmt.Println("  http://localhost:8080/policy?id=123 (success)")
	fmt.Println("  http://localhost:8080/policy?id=unauthorized (401)")
	fmt.Println("  http://localhost:8080/policy?id=forbidden (403)")
	fmt.Println("  http://localhost:8080/policy?id=notfound (404)")
	fmt.Println("  http://localhost:8080/policy?id=invalid (400)")
	fmt.Println("  http://localhost:8080/policy (400)")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Generate and Run

```bash
# Generate the error code
go generate

# Run the service
go run .
```

## Example API Responses

### Success (200)
```bash
curl "http://localhost:8080/policy?id=123"
```
```json
{
  "id": "123",
  "kind": "standard", 
  "name": "Sample Policy"
}
```

### Authentication Failed (401)
```bash
curl "http://localhost:8080/policy?id=unauthorized"
```
```json
{
  "code": 10001,
  "message": "Authentication failed"
}
```

### Authorization Denied (403)
```bash
curl "http://localhost:8080/policy?id=forbidden"
```
```json
{
  "code": 10002,
  "message": "Authorization denied"
}
```

### Policy Not Found (404)
```bash
curl "http://localhost:8080/policy?id=notfound"
```
```json
{
  "code": 20001,
  "message": "Policy not found"
}
```

### Invalid Policy Kind (400)
```bash
curl "http://localhost:8080/policy?id=invalid"
```
```json
{
  "code": 20002,
  "message": "Invalid policy kind"
}
```

## Benefits Demonstrated

1. **Type Safety**: All error codes are compile-time constants
2. **Consistency**: Standardized error response format
3. **HTTP Integration**: Automatic HTTP status code mapping
4. **gRPC Ready**: gRPC status codes included
5. **Additional Data**: Structured error data support
6. **Error Wrapping**: Original error preservation
7. **Performance**: No runtime lookups or maps