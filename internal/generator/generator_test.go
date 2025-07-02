package generator

import (
	"strings"
	"testing"
)

func TestParseInput_YAML(t *testing.T) {
	yamlInput := `
- code: 20001
  key: PolicyNotFound
  message: Policy not found
  http: 404
  grpc: 5
  desc: Policy could not be located in the database

- code: 20002
  key: InvalidKind
  message: Invalid policy kind
  http: 400
  grpc: 3
  desc: Policy kind is not supported
`

	errors, err := ParseInput(strings.NewReader(yamlInput), "test.yaml")
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	// Check first error
	if errors[0].Code != 20001 {
		t.Errorf("Expected code 20001, got %d", errors[0].Code)
	}
	if errors[0].Key != "PolicyNotFound" {
		t.Errorf("Expected key PolicyNotFound, got %s", errors[0].Key)
	}
	if errors[0].Message != "Policy not found" {
		t.Errorf("Expected message 'Policy not found', got %s", errors[0].Message)
	}
	if errors[0].HTTP != 404 {
		t.Errorf("Expected HTTP 404, got %d", errors[0].HTTP)
	}
	if errors[0].GRPC != 5 {
		t.Errorf("Expected GRPC 5, got %d", errors[0].GRPC)
	}
	if errors[0].Desc != "Policy could not be located in the database" {
		t.Errorf("Expected desc 'Policy could not be located in the database', got %s", errors[0].Desc)
	}

	// Check second error
	if errors[1].Code != 20002 {
		t.Errorf("Expected code 20002, got %d", errors[1].Code)
	}
	if errors[1].Key != "InvalidKind" {
		t.Errorf("Expected key InvalidKind, got %s", errors[1].Key)
	}
}

func TestParseInput_JSON(t *testing.T) {
	jsonInput := `[
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
    "key": "InvalidKind",
    "message": "Invalid policy kind",
    "http": 400,
    "grpc": 3,
    "desc": "Policy kind is not supported"
  }
]`

	errors, err := ParseInput(strings.NewReader(jsonInput), "test.json")
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	// Check first error
	if errors[0].Code != 20001 {
		t.Errorf("Expected code 20001, got %d", errors[0].Code)
	}
	if errors[0].Key != "PolicyNotFound" {
		t.Errorf("Expected key PolicyNotFound, got %s", errors[0].Key)
	}
}

func TestParseInput_AutoDetect_JSON(t *testing.T) {
	jsonInput := `[{"code": 20001, "key": "Test", "message": "Test message", "http": 400, "grpc": 3}]`

	errors, err := ParseInput(strings.NewReader(jsonInput), "test.unknown")
	if err != nil {
		t.Fatalf("Failed to auto-detect JSON: %v", err)
	}

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
}

func TestParseInput_AutoDetect_YAML(t *testing.T) {
	yamlInput := `- code: 20001
  key: Test
  message: Test message
  http: 400
  grpc: 3`

	errors, err := ParseInput(strings.NewReader(yamlInput), "test.unknown")
	if err != nil {
		t.Fatalf("Failed to auto-detect YAML: %v", err)
	}

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
}

func TestParseInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "missing code",
			input:   `[{"key": "Test", "message": "Test message", "http": 400, "grpc": 3}]`,
			wantErr: "code cannot be 0",
		},
		{
			name:    "missing key",
			input:   `[{"code": 20001, "message": "Test message", "http": 400, "grpc": 3}]`,
			wantErr: "key cannot be empty",
		},
		{
			name:    "missing message",
			input:   `[{"code": 20001, "key": "Test", "http": 400, "grpc": 3}]`,
			wantErr: "message cannot be empty",
		},
		{
			name:    "missing http",
			input:   `[{"code": 20001, "key": "Test", "message": "Test message", "grpc": 3}]`,
			wantErr: "http code cannot be 0",
		},
		{
			name:    "invalid grpc code",
			input:   `[{"code": 20001, "key": "Test", "message": "Test message", "http": 400, "grpc": 17}]`,
			wantErr: "grpc code must be between 0 and 16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseInput(strings.NewReader(tt.input), "test.json")
			if err == nil {
				t.Errorf("Expected error containing %q, got nil", tt.wantErr)
			} else if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	config := Config{
		Package: "testpkg",
		Errors: []ErrorDefinition{
			{
				Code:    20001,
				Key:     "PolicyNotFound",
				Message: "Policy not found",
				HTTP:    404,
				GRPC:    5,
				Desc:    "Policy could not be located in the database",
			},
			{
				Code:    20002,
				Key:     "InvalidKind",
				Message: "Invalid policy kind",
				HTTP:    400,
				GRPC:    3,
				Desc:    "Policy kind is not supported",
			},
		},
	}

	code, err := Generate(config)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	codeStr := string(code)

	// Check package declaration
	if !strings.Contains(codeStr, "package testpkg") {
		t.Error("Generated code should contain package declaration")
	}

	// Check imports
	if !strings.Contains(codeStr, `"github.com/restayway/rescode"`) {
		t.Error("Generated code should import rescode package")
	}
	if !strings.Contains(codeStr, `"google.golang.org/grpc/codes"`) {
		t.Error("Generated code should import grpc codes package")
	}

	// Check constants
	expectedConstants := []string{
		"PolicyNotFoundCode uint64",
		"= 20001",
		"PolicyNotFoundHTTP int",
		"= 404",
		"PolicyNotFoundGRPC codes.Code",
		"= 5",
		"PolicyNotFoundMsg  string",
		`= "Policy not found"`,
		"PolicyNotFoundDesc string",
		`= "Policy could not be located in the database"`,
		"InvalidKindCode uint64",
		"= 20002",
		"InvalidKindHTTP int",
		"= 400",
		"InvalidKindGRPC codes.Code",
		"= 3",
		"InvalidKindMsg  string",
		`= "Invalid policy kind"`,
		"InvalidKindDesc string",
		`= "Policy kind is not supported"`,
	}

	for _, expected := range expectedConstants {
		if !strings.Contains(codeStr, expected) {
			t.Errorf("Generated code should contain constant: %s", expected)
		}
	}

	// Check factory functions
	expectedFunctions := []string{
		"func PolicyNotFound(err ...error) *rescode.RC {",
		"return rescode.New(PolicyNotFoundCode, PolicyNotFoundHTTP, PolicyNotFoundGRPC, PolicyNotFoundMsg)(err...)",
		"func InvalidKind(err ...error) *rescode.RC {",
		"return rescode.New(InvalidKindCode, InvalidKindHTTP, InvalidKindGRPC, InvalidKindMsg)(err...)",
	}

	for _, expected := range expectedFunctions {
		if !strings.Contains(codeStr, expected) {
			t.Errorf("Generated code should contain function: %s", expected)
		}
	}

	// Check comments
	if !strings.Contains(codeStr, "// PolicyNotFound creates a new PolicyNotFound error.") {
		t.Error("Generated code should contain function comment")
	}
	if !strings.Contains(codeStr, "// Policy could not be located in the database") {
		t.Error("Generated code should contain description comment")
	}
}

func TestGenerate_DefaultPackage(t *testing.T) {
	config := Config{
		Package: "", // Empty package should default to "main"
		Errors: []ErrorDefinition{
			{
				Code:    20001,
				Key:     "TestError",
				Message: "Test message",
				HTTP:    400,
				GRPC:    3,
			},
		},
	}

	code, err := Generate(config)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	codeStr := string(code)
	if !strings.Contains(codeStr, "package main") {
		t.Error("Generated code should default to package main")
	}
}

func TestGenerate_NoDescription(t *testing.T) {
	config := Config{
		Package: "testpkg",
		Errors: []ErrorDefinition{
			{
				Code:    20001,
				Key:     "TestError",
				Message: "Test message",
				HTTP:    400,
				GRPC:    3,
				// No Desc field
			},
		},
	}

	code, err := Generate(config)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	codeStr := string(code)

	// Should not generate Desc constant
	if strings.Contains(codeStr, "TestErrorDesc") {
		t.Error("Generated code should not contain Desc constant when not provided")
	}

	// Should still generate function
	if !strings.Contains(codeStr, "func TestError(err ...error) *rescode.RC {") {
		t.Error("Generated code should contain function even without description")
	}
}

// Benchmark tests
func BenchmarkParseInput_YAML(b *testing.B) {
	yamlInput := `
- code: 20001
  key: PolicyNotFound
  message: Policy not found
  http: 404
  grpc: 5
  desc: Policy could not be located in the database
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseInput(strings.NewReader(yamlInput), "test.yaml")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseInput_JSON(b *testing.B) {
	jsonInput := `[{"code": 20001, "key": "PolicyNotFound", "message": "Policy not found", "http": 404, "grpc": 5, "desc": "Policy could not be located in the database"}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseInput(strings.NewReader(jsonInput), "test.json")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	config := Config{
		Package: "testpkg",
		Errors: []ErrorDefinition{
			{
				Code:    20001,
				Key:     "PolicyNotFound",
				Message: "Policy not found",
				HTTP:    404,
				GRPC:    5,
				Desc:    "Policy could not be located in the database",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}
