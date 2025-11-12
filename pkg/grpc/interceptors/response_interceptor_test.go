// SPDX-License-Identifier: MIT
/*
 * Copyright (c) 2025, SCANOSS
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package interceptors

import (
	"context"
	"errors"
	"net/http"
	"testing"

	common "github.com/scanoss/papi/api/commonv2"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MockResponse is a test response with a Status field
type MockResponse struct {
	Status *common.StatusResponse
	Data   string
}

// MockResponseWithoutStatus is a test response without a Status field
type MockResponseWithoutStatus struct {
	Data string
}

// MockResponseWithPrivateStatus has an unexported Status field
type MockResponseWithPrivateStatus struct {
	status *common.StatusResponse
	Data   string
}

func TestSetStatusField(t *testing.T) {
	tests := []struct {
		name           string
		resp           interface{}
		status         *common.StatusResponse
		expectSet      bool
		expectedStatus common.StatusCode
		expectedMsg    string
	}{
		{
			name: "struct with Status field",
			resp: &MockResponse{
				Data: "test data",
			},
			status: &common.StatusResponse{
				Status:  common.StatusCode_FAILED,
				Message: "error occurred",
			},
			expectSet:      true,
			expectedStatus: common.StatusCode_FAILED,
			expectedMsg:    "error occurred",
		},
		{
			name: "struct without Status field",
			resp: &MockResponseWithoutStatus{
				Data: "test data",
			},
			status: &common.StatusResponse{
				Status:  common.StatusCode_FAILED,
				Message: "error occurred",
			},
			expectSet: false,
		},
		{
			name: "struct with unexported status field",
			resp: &MockResponseWithPrivateStatus{
				Data: "test data",
			},
			status: &common.StatusResponse{
				Status:  common.StatusCode_FAILED,
				Message: "error occurred",
			},
			expectSet: false,
		},
		{
			name: "non-pointer struct",
			resp: MockResponse{
				Data: "test data",
			},
			status: &common.StatusResponse{
				Status:  common.StatusCode_FAILED,
				Message: "error occurred",
			},
			expectSet: false,
		},
		{
			name: "nil response",
			resp: nil,
			status: &common.StatusResponse{
				Status:  common.StatusCode_FAILED,
				Message: "error occurred",
			},
			expectSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setStatusField(tt.resp, tt.status)

			if tt.expectSet {
				mockResp, ok := tt.resp.(*MockResponse)
				if !ok {
					t.Fatal("expected MockResponse type")
				}
				if mockResp.Status == nil {
					t.Fatal("expected Status to be set, but it was nil")
				}
				if mockResp.Status.Status != tt.expectedStatus {
					t.Errorf("expected status %v, got %v", tt.expectedStatus, mockResp.Status.Status)
				}
				if mockResp.Status.Message != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, mockResp.Status.Message)
				}
			} else {
				// Verify nothing panicked and function handled gracefully
				switch resp := tt.resp.(type) {
				case *MockResponse:
					// Should not reach here in these tests
					if resp.Status != nil {
						t.Error("Status should not be set for this test case")
					}
				case *MockResponseWithoutStatus:
					// No status field, nothing to check
				case *MockResponseWithPrivateStatus:
					// Private status field should not be set
					if resp.status != nil {
						t.Error("private status field should not be set")
					}
				case MockResponse:
					// Non-pointer case, nothing to check
				}
			}
		})
	}
}

func TestResponseInterceptor(t *testing.T) {
	tests := []struct {
		name            string
		handlerResp     interface{}
		handlerErr      error
		expectErr       bool
		expectStatusSet bool
		expectedStatus  common.StatusCode
		expectedMessage string
		validateResp    func(t *testing.T, resp interface{})
	}{
		{
			name: "successful request - no error",
			handlerResp: &MockResponse{
				Data: "success",
			},
			handlerErr:      nil,
			expectErr:       false,
			expectStatusSet: false,
		},
		{
			name: "ResponseError with BadRequest",
			handlerResp: &MockResponse{
				Data: "failed",
			},
			handlerErr: &ResponseError{
				Message:      "invalid input",
				HTTPCode:     http.StatusBadRequest,
				InternalCode: "FAILED",
			},
			expectErr:       false,
			expectStatusSet: true,
			expectedStatus:  common.StatusCode_FAILED,
			expectedMessage: "invalid input",
		},
		{
			name: "ResponseError with NotFound",
			handlerResp: &MockResponse{
				Data: "failed",
			},
			handlerErr: &ResponseError{
				Message:      "resource not found",
				HTTPCode:     http.StatusNotFound,
				InternalCode: "FAILED",
			},
			expectErr:       false,
			expectStatusSet: true,
			expectedStatus:  common.StatusCode_FAILED,
			expectedMessage: "resource not found",
		},
		{
			name: "ResponseError with InternalServerError",
			handlerResp: &MockResponse{
				Data: "failed",
			},
			handlerErr: &ResponseError{
				Message:      "unexpected error",
				HTTPCode:     http.StatusInternalServerError,
				InternalCode: "FAILED",
			},
			expectErr:       false,
			expectStatusSet: true,
			expectedStatus:  common.StatusCode_FAILED,
			expectedMessage: "unexpected error",
		},
		{
			name: "standard error",
			handlerResp: &MockResponse{
				Data: "failed",
			},
			handlerErr:      errors.New("some standard error"),
			expectErr:       false,
			expectStatusSet: true,
			expectedStatus:  common.StatusCode_FAILED,
			expectedMessage: "internal server error",
		},
		{
			name:        "nil response with error",
			handlerResp: nil,
			handlerErr: &ResponseError{
				Message:      "service unavailable",
				HTTPCode:     http.StatusServiceUnavailable,
				InternalCode: "FAILED",
			},
			expectErr:       false,
			expectStatusSet: false, // Can't set status on nil response
			validateResp: func(t *testing.T, resp interface{}) {
				if resp != nil {
					t.Errorf("expected nil response, got %v", resp)
				}
			},
		},
		{
			name: "response without Status field",
			handlerResp: &MockResponseWithoutStatus{
				Data: "failed",
			},
			handlerErr: &ResponseError{
				Message:      "error message",
				HTTPCode:     http.StatusBadRequest,
				InternalCode: "FAILED",
			},
			expectErr:       false,
			expectStatusSet: false, // Can't set on response without Status field
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a logger for the context
			logger := zap.NewNop()
			ctx := ctxzap.ToContext(context.Background(), logger)
			ctx = metadata.NewOutgoingContext(ctx, metadata.MD{})

			// Create a mock handler
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return tt.handlerResp, tt.handlerErr
			}

			// Create the interceptor
			interceptor := ResponseInterceptor()

			// Execute the interceptor
			resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
				FullMethod: "/test.Service/Method",
			}, handler)

			// Validate error
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// Custom validation function
			if tt.validateResp != nil {
				tt.validateResp(t, resp)
				return
			}

			// Validate Status field was set correctly
			if tt.expectStatusSet {
				mockResp, ok := resp.(*MockResponse)
				if !ok {
					t.Fatalf("expected MockResponse type, got %T", resp)
				}

				if mockResp.Status == nil {
					t.Fatal("expected Status to be set, but it was nil")
				}

				if mockResp.Status.Status != tt.expectedStatus {
					t.Errorf("expected status %v, got %v", tt.expectedStatus, mockResp.Status.Status)
				}

				if mockResp.Status.Message != tt.expectedMessage {
					t.Errorf("expected message %q, got %q", tt.expectedMessage, mockResp.Status.Message)
				}
			}

			// For successful cases, verify the response is returned unchanged
			if !tt.expectStatusSet && tt.handlerErr == nil {
				if resp != tt.handlerResp {
					t.Error("expected response to be unchanged")
				}
			}
		})
	}
}

func TestResponseInterceptor_Integration(t *testing.T) {
	// This test simulates a more realistic scenario with proper context setup
	logger := zap.NewNop()
	ctx := ctxzap.ToContext(context.Background(), logger)
	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{})

	// Test case: Handler returns error
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Simulate some processing that results in an error
		return &MockResponse{Data: "processing"}, &ResponseError{
			Message:      "validation failed",
			HTTPCode:     http.StatusUnprocessableEntity,
			InternalCode: "VALIDATION_ERROR",
			Details: map[string]interface{}{
				"field": "email",
				"issue": "invalid format",
			},
		}
	}

	interceptor := ResponseInterceptor()

	resp, err := interceptor(ctx, map[string]string{"test": "request"}, &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.UserService/CreateUser",
	}, handler)

	// Interceptor should return nil error (converts to custom response)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Verify Status was set
	mockResp, ok := resp.(*MockResponse)
	if !ok {
		t.Fatalf("expected MockResponse type, got %T", resp)
	}

	if mockResp.Status == nil {
		t.Fatal("expected Status to be set")
	}

	if mockResp.Status.Status != common.StatusCode_FAILED {
		t.Errorf("expected FAILED status, got %v", mockResp.Status.Status)
	}

	if mockResp.Status.Message != "validation failed" {
		t.Errorf("expected message 'validation failed', got %q", mockResp.Status.Message)
	}
}
