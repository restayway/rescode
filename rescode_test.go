package rescode

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestRC_Basic(t *testing.T) {
	creator := New(1001, 400, codes.InvalidArgument, "test error")
	rc := creator()

	if rc.Code != 1001 {
		t.Errorf("Expected Code 1001, got %d", rc.Code)
	}
	if rc.HttpCode != 400 {
		t.Errorf("Expected HttpCode 400, got %d", rc.HttpCode)
	}
	if rc.RpcCode != codes.InvalidArgument {
		t.Errorf("Expected RpcCode InvalidArgument, got %v", rc.RpcCode)
	}
	if rc.Message != "test error" {
		t.Errorf("Expected Message 'test error', got %s", rc.Message)
	}
	if rc.Data != nil {
		t.Errorf("Expected Data to be nil, got %v", rc.Data)
	}
	if rc.err != nil {
		t.Errorf("Expected err to be nil, got %v", rc.err)
	}
}

func TestRC_WithData(t *testing.T) {
	testData := map[string]string{"key": "value"}
	creator := New(1002, 404, codes.NotFound, "not found", testData)
	rc := creator()

	// Use type assertion and comparison since maps can't be compared directly
	if dataMap, ok := rc.Data.(map[string]string); !ok {
		t.Errorf("Expected Data to be map[string]string, got %T", rc.Data)
	} else if dataMap["key"] != "value" {
		t.Errorf("Expected Data['key'] to be 'value', got %v", dataMap["key"])
	}
}

func TestRC_WithWrappedError(t *testing.T) {
	originalErr := errors.New("original error")
	creator := New(1003, 500, codes.Internal, "internal error")
	rc := creator(originalErr)

	if rc.err != originalErr {
		t.Errorf("Expected wrapped error %v, got %v", originalErr, rc.err)
	}
	if rc.OriginalError() != originalErr {
		t.Errorf("Expected OriginalError() %v, got %v", originalErr, rc.OriginalError())
	}
}

func TestRC_Error(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		wrappedErr error
		expected   string
	}{
		{
			name:       "without wrapped error",
			message:    "simple error",
			wrappedErr: nil,
			expected:   "simple error",
		},
		{
			name:       "with wrapped error",
			message:    "parent error",
			wrappedErr: errors.New("child error"),
			expected:   "parent error: child error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator := New(1000, 400, codes.InvalidArgument, tt.message)
			var rc *RC
			if tt.wrappedErr != nil {
				rc = creator(tt.wrappedErr)
			} else {
				rc = creator()
			}

			if rc.Error() != tt.expected {
				t.Errorf("Expected Error() %q, got %q", tt.expected, rc.Error())
			}
		})
	}
}

func TestRC_SetData(t *testing.T) {
	creator := New(1004, 400, codes.InvalidArgument, "test error")
	rc := creator()

	testData := "new data"
	result := rc.SetData(testData)

	// Should return the same RC for chaining
	if result != rc {
		t.Error("SetData should return the same RC instance for chaining")
	}

	if rc.Data != testData {
		t.Errorf("Expected Data %v, got %v", testData, rc.Data)
	}
}

func TestRC_JSON(t *testing.T) {
	testData := map[string]interface{}{"test": "data"}
	originalErr := errors.New("wrapped error")
	creator := New(1005, 404, codes.NotFound, "test message", testData)
	rc := creator(originalErr)

	json := rc.JSON()

	expectedKeys := []string{"code", "message", "httpCode", "rpcCode", "data", "originalError"}
	for _, key := range expectedKeys {
		if _, exists := json[key]; !exists {
			t.Errorf("Expected JSON to contain key %s", key)
		}
	}

	if json["code"] != uint64(1005) {
		t.Errorf("Expected code 1005, got %v", json["code"])
	}
	if json["message"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", json["message"])
	}
	if json["httpCode"] != 404 {
		t.Errorf("Expected httpCode 404, got %v", json["httpCode"])
	}
	if json["rpcCode"] != int(codes.NotFound) {
		t.Errorf("Expected rpcCode %d, got %v", int(codes.NotFound), json["rpcCode"])
	}
	if dataMap, ok := json["data"].(map[string]interface{}); !ok {
		t.Errorf("Expected data to be map[string]interface{}, got %T", json["data"])
	} else if dataMap["test"] != "data" {
		t.Errorf("Expected data['test'] to be 'data', got %v", dataMap["test"])
	}
	if json["originalError"] != "wrapped error" {
		t.Errorf("Expected originalError 'wrapped error', got %v", json["originalError"])
	}
}

func TestRC_JSON_FilteredKeys(t *testing.T) {
	creator := New(1006, 400, codes.InvalidArgument, "test message")
	rc := creator()

	json := rc.JSON("code", "message")

	if len(json) != 2 {
		t.Errorf("Expected JSON to have 2 keys, got %d", len(json))
	}

	if json["code"] != uint64(1006) {
		t.Errorf("Expected code 1006, got %v", json["code"])
	}
	if json["message"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", json["message"])
	}

	// Should not contain other keys
	if _, exists := json["httpCode"]; exists {
		t.Error("JSON should not contain httpCode when filtered")
	}
}

func TestRC_String(t *testing.T) {
	testData := "test data"
	originalErr := errors.New("wrapped error")
	creator := New(1007, 400, codes.InvalidArgument, "test message", testData)
	rc := creator(originalErr)

	str := rc.String()

	// Should contain all the components
	expected := []string{
		"Code:1007",
		"HTTP:400",
		"gRPC:3", // InvalidArgument is code 3
		"Message:test message",
		"Data:test data",
		"OriginalError:wrapped error",
	}

	for _, exp := range expected {
		if !contains(str, exp) {
			t.Errorf("Expected String() to contain %q, got %q", exp, str)
		}
	}
}

func TestRC_String_Minimal(t *testing.T) {
	creator := New(1008, 200, codes.OK, "simple message")
	rc := creator()

	str := rc.String()

	// Should contain basic components
	expected := []string{
		"Code:1008",
		"HTTP:200",
		"gRPC:0", // OK is code 0
		"Message:simple message",
	}

	for _, exp := range expected {
		if !contains(str, exp) {
			t.Errorf("Expected String() to contain %q, got %q", exp, str)
		}
	}

	// Should not contain Data or OriginalError
	if contains(str, "Data:") {
		t.Error("String() should not contain Data when nil")
	}
	if contains(str, "OriginalError:") {
		t.Error("String() should not contain OriginalError when nil")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Benchmark tests for performance comparison
func BenchmarkRC_Creation(b *testing.B) {
	creator := New(1000, 400, codes.InvalidArgument, "benchmark error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = creator()
	}
}

func BenchmarkRC_CreationWithError(b *testing.B) {
	creator := New(1000, 400, codes.InvalidArgument, "benchmark error")
	err := errors.New("wrapped error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = creator(err)
	}
}

func BenchmarkRC_Error(b *testing.B) {
	creator := New(1000, 400, codes.InvalidArgument, "benchmark error")
	rc := creator(errors.New("wrapped error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rc.Error()
	}
}

func BenchmarkRC_JSON(b *testing.B) {
	creator := New(1000, 400, codes.InvalidArgument, "benchmark error", map[string]string{"key": "value"})
	rc := creator(errors.New("wrapped error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rc.JSON()
	}
}
