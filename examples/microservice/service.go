package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/restayway/rescode"
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
			"field":  "id",
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
