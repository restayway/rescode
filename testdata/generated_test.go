package testdata

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
)

// Test the generated error factories
func TestPolicyNotFound(t *testing.T) {
	err := PolicyNotFound()

	if err.Code != PolicyNotFoundCode {
		t.Errorf("Expected Code %d, got %d", PolicyNotFoundCode, err.Code)
	}
	if err.HttpCode != PolicyNotFoundHTTP {
		t.Errorf("Expected HttpCode %d, got %d", PolicyNotFoundHTTP, err.HttpCode)
	}
	if err.RpcCode != PolicyNotFoundGRPC {
		t.Errorf("Expected RpcCode %d, got %d", PolicyNotFoundGRPC, err.RpcCode)
	}
	if err.Message != PolicyNotFoundMsg {
		t.Errorf("Expected Message %q, got %q", PolicyNotFoundMsg, err.Message)
	}
	if err.Error() != PolicyNotFoundMsg {
		t.Errorf("Expected Error() %q, got %q", PolicyNotFoundMsg, err.Error())
	}
}

func TestPolicyNotFound_WithWrappedError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	err := PolicyNotFound(originalErr)

	if err.OriginalError() != originalErr {
		t.Errorf("Expected OriginalError() %v, got %v", originalErr, err.OriginalError())
	}

	expectedErrorMsg := "Policy not found: database connection failed"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected Error() %q, got %q", expectedErrorMsg, err.Error())
	}
}

func TestInvalidKind(t *testing.T) {
	err := InvalidKind()

	if err.Code != InvalidKindCode {
		t.Errorf("Expected Code %d, got %d", InvalidKindCode, err.Code)
	}
	if err.HttpCode != InvalidKindHTTP {
		t.Errorf("Expected HttpCode %d, got %d", InvalidKindHTTP, err.HttpCode)
	}
	if err.RpcCode != InvalidKindGRPC {
		t.Errorf("Expected RpcCode %d, got %d", InvalidKindGRPC, err.RpcCode)
	}
	if err.Message != InvalidKindMsg {
		t.Errorf("Expected Message %q, got %q", InvalidKindMsg, err.Message)
	}
}

func TestInternalError(t *testing.T) {
	err := InternalError()

	if err.Code != InternalErrorCode {
		t.Errorf("Expected Code %d, got %d", InternalErrorCode, err.Code)
	}
	if err.HttpCode != InternalErrorHTTP {
		t.Errorf("Expected HttpCode %d, got %d", InternalErrorHTTP, err.HttpCode)
	}
	if err.RpcCode != InternalErrorGRPC {
		t.Errorf("Expected RpcCode %d, got %d", InternalErrorGRPC, err.RpcCode)
	}
	if err.Message != InternalErrorMsg {
		t.Errorf("Expected Message %q, got %q", InternalErrorMsg, err.Message)
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that constants have expected values
	if PolicyNotFoundCode != 20001 {
		t.Errorf("Expected PolicyNotFoundCode 20001, got %d", PolicyNotFoundCode)
	}
	if PolicyNotFoundHTTP != 404 {
		t.Errorf("Expected PolicyNotFoundHTTP 404, got %d", PolicyNotFoundHTTP)
	}
	if PolicyNotFoundGRPC != codes.NotFound {
		t.Errorf("Expected PolicyNotFoundGRPC %d, got %d", codes.NotFound, PolicyNotFoundGRPC)
	}
	if PolicyNotFoundMsg != "Policy not found" {
		t.Errorf("Expected PolicyNotFoundMsg 'Policy not found', got %q", PolicyNotFoundMsg)
	}

	if InvalidKindCode != 20002 {
		t.Errorf("Expected InvalidKindCode 20002, got %d", InvalidKindCode)
	}
	if InvalidKindHTTP != 400 {
		t.Errorf("Expected InvalidKindHTTP 400, got %d", InvalidKindHTTP)
	}
	if InvalidKindGRPC != codes.InvalidArgument {
		t.Errorf("Expected InvalidKindGRPC %d, got %d", codes.InvalidArgument, InvalidKindGRPC)
	}
	if InvalidKindMsg != "Invalid policy kind" {
		t.Errorf("Expected InvalidKindMsg 'Invalid policy kind', got %q", InvalidKindMsg)
	}
}

func TestJSON_GeneratedErrors(t *testing.T) {
	err := PolicyNotFound().SetData(map[string]string{"resource": "policy_123"})

	json := err.JSON()

	if json["code"] != PolicyNotFoundCode {
		t.Errorf("Expected JSON code %d, got %v", PolicyNotFoundCode, json["code"])
	}
	if json["httpCode"] != PolicyNotFoundHTTP {
		t.Errorf("Expected JSON httpCode %d, got %v", PolicyNotFoundHTTP, json["httpCode"])
	}
	if json["rpcCode"] != int(PolicyNotFoundGRPC) {
		t.Errorf("Expected JSON rpcCode %d, got %v", int(PolicyNotFoundGRPC), json["rpcCode"])
	}
	if json["message"] != PolicyNotFoundMsg {
		t.Errorf("Expected JSON message %q, got %v", PolicyNotFoundMsg, json["message"])
	}

	if dataMap, ok := json["data"].(map[string]string); !ok {
		t.Errorf("Expected data to be map[string]string, got %T", json["data"])
	} else if dataMap["resource"] != "policy_123" {
		t.Errorf("Expected data[resource] 'policy_123', got %v", dataMap["resource"])
	}
}

// Benchmarks for generated code performance
func BenchmarkPolicyNotFound_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PolicyNotFound()
	}
}

func BenchmarkPolicyNotFound_CreationWithError(b *testing.B) {
	err := errors.New("wrapped error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PolicyNotFound(err)
	}
}

func BenchmarkPolicyNotFound_Error(b *testing.B) {
	err := PolicyNotFound(errors.New("wrapped error"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkPolicyNotFound_JSON(b *testing.B) {
	err := PolicyNotFound().SetData(map[string]string{"resource": "policy_123"})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.JSON()
	}
}
