package rescode

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
)

// LegacyErrorRegistry simulates a runtime map-based error system for comparison
type LegacyErrorRegistry struct {
	errors map[string]LegacyErrorDef
}

type LegacyErrorDef struct {
	Code     uint64
	Message  string
	HttpCode int
	RpcCode  codes.Code
}

// NewLegacyRegistry creates a legacy error registry with runtime lookups
func NewLegacyRegistry() *LegacyErrorRegistry {
	return &LegacyErrorRegistry{
		errors: map[string]LegacyErrorDef{
			"PolicyNotFound": {
				Code:     20001,
				Message:  "Policy not found",
				HttpCode: 404,
				RpcCode:  codes.NotFound,
			},
			"InvalidKind": {
				Code:     20002,
				Message:  "Invalid policy kind",
				HttpCode: 400,
				RpcCode:  codes.InvalidArgument,
			},
			"InternalError": {
				Code:     20003,
				Message:  "Internal server error",
				HttpCode: 500,
				RpcCode:  codes.Internal,
			},
		},
	}
}

func (r *LegacyErrorRegistry) CreateError(key string, err ...error) (*RC, error) {
	def, exists := r.errors[key]
	if !exists {
		return nil, fmt.Errorf("error definition not found: %s", key)
	}
	
	rc := &RC{
		Code:     def.Code,
		Message:  def.Message,
		HttpCode: def.HttpCode,
		RpcCode:  def.RpcCode,
	}
	
	if len(err) > 0 {
		rc.err = err[0]
	}
	
	return rc, nil
}

// Benchmarks comparing generated code vs legacy runtime approach
func BenchmarkGenerated_PolicyNotFound(b *testing.B) {
	// This uses the rescode.New approach with compile-time constants
	creator := New(20001, 404, codes.NotFound, "Policy not found")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = creator()
	}
}

func BenchmarkLegacy_PolicyNotFound(b *testing.B) {
	registry := NewLegacyRegistry()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.CreateError("PolicyNotFound")
	}
}

func BenchmarkGenerated_PolicyNotFound_WithError(b *testing.B) {
	creator := New(20001, 404, codes.NotFound, "Policy not found")
	err := errors.New("wrapped error")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = creator(err)
	}
}

func BenchmarkLegacy_PolicyNotFound_WithError(b *testing.B) {
	registry := NewLegacyRegistry()
	err := errors.New("wrapped error")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.CreateError("PolicyNotFound", err)
	}
}

func BenchmarkGenerated_MultipleErrors(b *testing.B) {
	policyNotFound := New(20001, 404, codes.NotFound, "Policy not found")
	invalidKind := New(20002, 400, codes.InvalidArgument, "Invalid policy kind")
	internalError := New(20003, 500, codes.Internal, "Internal server error")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 3 {
		case 0:
			_ = policyNotFound()
		case 1:
			_ = invalidKind()
		case 2:
			_ = internalError()
		}
	}
}

func BenchmarkLegacy_MultipleErrors(b *testing.B) {
	registry := NewLegacyRegistry()
	keys := []string{"PolicyNotFound", "InvalidKind", "InternalError"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%3]
		_, _ = registry.CreateError(key)
	}
}

func BenchmarkGenerated_ErrorMessage(b *testing.B) {
	creator := New(20001, 404, codes.NotFound, "Policy not found")
	err := creator(errors.New("wrapped error"))
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkLegacy_ErrorMessage(b *testing.B) {
	registry := NewLegacyRegistry()
	err, _ := registry.CreateError("PolicyNotFound", errors.New("wrapped error"))
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkGenerated_JSON(b *testing.B) {
	creator := New(20001, 404, codes.NotFound, "Policy not found", map[string]string{"resource": "policy_123"})
	err := creator(errors.New("wrapped error"))
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.JSON()
	}
}

func BenchmarkLegacy_JSON(b *testing.B) {
	registry := NewLegacyRegistry()
	err, _ := registry.CreateError("PolicyNotFound", errors.New("wrapped error"))
	err.SetData(map[string]string{"resource": "policy_123"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.JSON()
	}
}

// Test to verify legacy vs generated produce same results
func TestLegacyVsGenerated_Compatibility(t *testing.T) {
	registry := NewLegacyRegistry()
	creator := New(20001, 404, codes.NotFound, "Policy not found")
	
	// Test basic creation
	legacyErr, _ := registry.CreateError("PolicyNotFound")
	generatedErr := creator()
	
	if legacyErr.Code != generatedErr.Code {
		t.Errorf("Code mismatch: legacy %d, generated %d", legacyErr.Code, generatedErr.Code)
	}
	if legacyErr.Message != generatedErr.Message {
		t.Errorf("Message mismatch: legacy %q, generated %q", legacyErr.Message, generatedErr.Message)
	}
	if legacyErr.HttpCode != generatedErr.HttpCode {
		t.Errorf("HttpCode mismatch: legacy %d, generated %d", legacyErr.HttpCode, generatedErr.HttpCode)
	}
	if legacyErr.RpcCode != generatedErr.RpcCode {
		t.Errorf("RpcCode mismatch: legacy %d, generated %d", legacyErr.RpcCode, generatedErr.RpcCode)
	}
	
	// Test with wrapped error
	originalErr := errors.New("wrapped error")
	legacyWithErr, _ := registry.CreateError("PolicyNotFound", originalErr)
	generatedWithErr := creator(originalErr)
	
	if legacyWithErr.Error() != generatedWithErr.Error() {
		t.Errorf("Error() mismatch: legacy %q, generated %q", legacyWithErr.Error(), generatedWithErr.Error())
	}
}