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

// StaticCodeLegacy approach: uint64 constants with message map
const (
	StaticCodePolicyNotFound uint64 = 20001
	StaticCodeInvalidKind    uint64 = 20002
	StaticCodeInternalError  uint64 = 20003
)

var staticCodeMessages = map[uint64]string{
	StaticCodePolicyNotFound: "Policy not found",
	StaticCodeInvalidKind:    "Invalid policy kind",
	StaticCodeInternalError:  "Internal server error",
}

var staticCodeHttpCodes = map[uint64]int{
	StaticCodePolicyNotFound: 404,
	StaticCodeInvalidKind:    400,
	StaticCodeInternalError:  500,
}

var staticCodeRpcCodes = map[uint64]codes.Code{
	StaticCodePolicyNotFound: codes.NotFound,
	StaticCodeInvalidKind:    codes.InvalidArgument,
	StaticCodeInternalError:  codes.Internal,
}

func CreateStaticCodeError(code uint64, err ...error) (*RC, error) {
	message, exists := staticCodeMessages[code]
	if !exists {
		return nil, fmt.Errorf("error code not found: %d", code)
	}

	rc := &RC{
		Code:     code,
		Message:  message,
		HttpCode: staticCodeHttpCodes[code],
		RpcCode:  staticCodeRpcCodes[code],
	}

	if len(err) > 0 {
		rc.err = err[0]
	}

	return rc, nil
}

// VarDeclarationLegacy approach: var declarations with maps
var (
	VarCodePolicyNotFound uint64 = 20001
	VarCodeInvalidKind    uint64 = 20002
	VarCodeInternalError  uint64 = 20003
)

var varCodeMessages = map[uint64]string{
	VarCodePolicyNotFound: "Policy not found",
	VarCodeInvalidKind:    "Invalid policy kind",
	VarCodeInternalError:  "Internal server error",
}

var varCodeHttpCodes = map[uint64]int{
	VarCodePolicyNotFound: 404,
	VarCodeInvalidKind:    400,
	VarCodeInternalError:  500,
}

var varCodeRpcCodes = map[uint64]codes.Code{
	VarCodePolicyNotFound: codes.NotFound,
	VarCodeInvalidKind:    codes.InvalidArgument,
	VarCodeInternalError:  codes.Internal,
}

func CreateVarCodeError(code uint64, err ...error) (*RC, error) {
	message, exists := varCodeMessages[code]
	if !exists {
		return nil, fmt.Errorf("error code not found: %d", code)
	}

	rc := &RC{
		Code:     code,
		Message:  message,
		HttpCode: varCodeHttpCodes[code],
		RpcCode:  varCodeRpcCodes[code],
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
	err := creator()

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

// Benchmarks for Static Code approach
func BenchmarkStaticCode_PolicyNotFound(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CreateStaticCodeError(StaticCodePolicyNotFound)
	}
}

func BenchmarkStaticCode_PolicyNotFound_WithError(b *testing.B) {
	err := errors.New("wrapped error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CreateStaticCodeError(StaticCodePolicyNotFound, err)
	}
}

func BenchmarkStaticCode_MultipleErrors(b *testing.B) {
	codes := []uint64{StaticCodePolicyNotFound, StaticCodeInvalidKind, StaticCodeInternalError}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code := codes[i%3]
		_, _ = CreateStaticCodeError(code)
	}
}

func BenchmarkStaticCode_ErrorMessage(b *testing.B) {
	err, _ := CreateStaticCodeError(StaticCodePolicyNotFound, errors.New("wrapped error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkStaticCode_JSON(b *testing.B) {
	err, _ := CreateStaticCodeError(StaticCodePolicyNotFound, errors.New("wrapped error"))
	err.SetData(map[string]string{"resource": "policy_123"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.JSON()
	}
}

// Benchmarks for Var Declaration approach
func BenchmarkVarCode_PolicyNotFound(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CreateVarCodeError(VarCodePolicyNotFound)
	}
}

func BenchmarkVarCode_PolicyNotFound_WithError(b *testing.B) {
	err := errors.New("wrapped error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CreateVarCodeError(VarCodePolicyNotFound, err)
	}
}

func BenchmarkVarCode_MultipleErrors(b *testing.B) {
	codes := []uint64{VarCodePolicyNotFound, VarCodeInvalidKind, VarCodeInternalError}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code := codes[i%3]
		_, _ = CreateVarCodeError(code)
	}
}

func BenchmarkVarCode_ErrorMessage(b *testing.B) {
	err, _ := CreateVarCodeError(VarCodePolicyNotFound, errors.New("wrapped error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVarCode_JSON(b *testing.B) {
	err, _ := CreateVarCodeError(VarCodePolicyNotFound, errors.New("wrapped error"))
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

// Test static code approach compatibility
func TestStaticCodeVsGenerated_Compatibility(t *testing.T) {
	creator := New(20001, 404, codes.NotFound, "Policy not found")

	// Test basic creation
	staticErr, _ := CreateStaticCodeError(StaticCodePolicyNotFound)
	generatedErr := creator()

	if staticErr.Code != generatedErr.Code {
		t.Errorf("Code mismatch: static %d, generated %d", staticErr.Code, generatedErr.Code)
	}
	if staticErr.Message != generatedErr.Message {
		t.Errorf("Message mismatch: static %q, generated %q", staticErr.Message, generatedErr.Message)
	}
	if staticErr.HttpCode != generatedErr.HttpCode {
		t.Errorf("HttpCode mismatch: static %d, generated %d", staticErr.HttpCode, generatedErr.HttpCode)
	}
	if staticErr.RpcCode != generatedErr.RpcCode {
		t.Errorf("RpcCode mismatch: static %d, generated %d", staticErr.RpcCode, generatedErr.RpcCode)
	}
}

// Test var code approach compatibility
func TestVarCodeVsGenerated_Compatibility(t *testing.T) {
	creator := New(20001, 404, codes.NotFound, "Policy not found")

	// Test basic creation
	varErr, _ := CreateVarCodeError(VarCodePolicyNotFound)
	generatedErr := creator()

	if varErr.Code != generatedErr.Code {
		t.Errorf("Code mismatch: var %d, generated %d", varErr.Code, generatedErr.Code)
	}
	if varErr.Message != generatedErr.Message {
		t.Errorf("Message mismatch: var %q, generated %q", varErr.Message, generatedErr.Message)
	}
	if varErr.HttpCode != generatedErr.HttpCode {
		t.Errorf("HttpCode mismatch: var %d, generated %d", varErr.HttpCode, generatedErr.HttpCode)
	}
	if varErr.RpcCode != generatedErr.RpcCode {
		t.Errorf("RpcCode mismatch: var %d, generated %d", varErr.RpcCode, generatedErr.RpcCode)
	}
}
