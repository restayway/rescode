package main

import (
	"fmt"
	"log"
)

//go:generate go run github.com/restayway/rescode/cmd/rescodegen --input errors.yaml --output errors_gen.go --package main

func main() {
	// Create errors without wrapping
	userErr := UserNotFound()
	fmt.Printf("User error: %v\n", userErr)
	fmt.Printf("HTTP code: %d\n", userErr.HttpCode)
	fmt.Printf("gRPC code: %d\n", userErr.RpcCode)

	// Create error with wrapped error
	dbErr := fmt.Errorf("connection timeout")
	wrappedErr := DatabaseError(dbErr)
	fmt.Printf("Database error: %v\n", wrappedErr)
	fmt.Printf("Original error: %v\n", wrappedErr.OriginalError())

	// Add additional data
	emailErr := InvalidEmail().SetData(map[string]string{
		"email":    "invalid@",
		"field":    "email",
		"location": "signup form",
	})
	fmt.Printf("Email error JSON: %v\n", emailErr.JSON())

	// Use constants directly
	if userErr.Code == UserNotFoundCode {
		log.Printf("Handling user not found with code %d", UserNotFoundCode)
	}
}