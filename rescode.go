// Package rescode provides type-safe error code functionality.
package rescode

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
)

// RC represents a structured error with multiple code formats and optional data.
type RC struct {
	Code     uint64     // Unique error code
	Message  string     // Human-readable error message
	HttpCode int        // HTTP status code
	RpcCode  codes.Code // gRPC status code
	Data     any        // Optional additional data
	err      error      // Wrapped original error
}

// RcCreator is a function type that creates an RC with optional wrapped errors.
type RcCreator func(...error) *RC

// New creates an RcCreator function with the specified parameters.
// This is designed to be used by generated code for efficient error creation.
func New(code uint64, hCode int, rCode codes.Code, message string, data ...any) RcCreator {
	var d any
	if len(data) > 0 {
		d = data[0]
	}
	
	return func(errs ...error) *RC {
		rc := &RC{
			Code:     code,
			Message:  message,
			HttpCode: hCode,
			RpcCode:  rCode,
			Data:     d,
		}
		
		if len(errs) > 0 {
			rc.err = errs[0]
		}
		
		return rc
	}
}

// Error implements the error interface.
func (r *RC) Error() string {
	if r.err != nil {
		return fmt.Sprintf("%s: %v", r.Message, r.err)
	}
	return r.Message
}

// SetData sets additional data for the error and returns the RC for chaining.
func (r *RC) SetData(data any) *RC {
	r.Data = data
	return r
}

// JSON returns a map representation of the error, optionally filtering by keys.
func (r *RC) JSON(keys ...string) map[string]interface{} {
	result := map[string]interface{}{
		"code":     r.Code,
		"message":  r.Message,
		"httpCode": r.HttpCode,
		"rpcCode":  int(r.RpcCode),
	}
	
	if r.Data != nil {
		result["data"] = r.Data
	}
	
	if r.err != nil {
		result["originalError"] = r.err.Error()
	}
	
	// If specific keys are requested, filter the result
	if len(keys) > 0 {
		filtered := make(map[string]interface{})
		for _, key := range keys {
			if val, exists := result[key]; exists {
				filtered[key] = val
			}
		}
		return filtered
	}
	
	return result
}

// OriginalError returns the wrapped original error, if any.
func (r *RC) OriginalError() error {
	return r.err
}

// String returns a string representation of the error.
func (r *RC) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Code:%d", r.Code))
	parts = append(parts, fmt.Sprintf("HTTP:%d", r.HttpCode))
	parts = append(parts, fmt.Sprintf("gRPC:%d", r.RpcCode))
	parts = append(parts, fmt.Sprintf("Message:%s", r.Message))
	
	if r.Data != nil {
		parts = append(parts, fmt.Sprintf("Data:%v", r.Data))
	}
	
	if r.err != nil {
		parts = append(parts, fmt.Sprintf("OriginalError:%v", r.err))
	}
	
	return fmt.Sprintf("RC{%s}", strings.Join(parts, ", "))
}