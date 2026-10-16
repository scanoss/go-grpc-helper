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
	"fmt"

	common "github.com/scanoss/papi/api/commonv2"

	"net/http"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestServiceError_getHTTPCode(t *testing.T) {
	tests := []struct {
		name     string
		err      *ResponseError
		expected int
	}{
		{
			name: "with explicit code",
			err: &ResponseError{
				HTTPCode: http.StatusBadRequest,
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "with zero code - defaults to 500",
			err: &ResponseError{
				HTTPCode: 0,
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.getHTTPCode()
			if result != tt.expected {
				t.Errorf("expected HTTP code %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestIsServiceError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "is ResponseError",
			err: &ResponseError{
				Message:      "test",
				HTTPCode:     http.StatusBadRequest,
				InternalCode: "FAILED",
				Err:          nil,
			},
			expected: true,
		},
		{
			name: "is wrapped ResponseError",
			err: fmt.Errorf("wrapped: %w", &ResponseError{
				Message:      "resource not found",
				HTTPCode:     http.StatusNotFound,
				InternalCode: "FAILED",
				Err:          nil,
			}),
			expected: true,
		},
		{
			name:     "is not ResponseError",
			err:      errors.New("standard error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isResponseError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetServiceError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		expectOk  bool
		expectNil bool
	}{
		{
			name: "direct ResponseError",
			err: &ResponseError{
				Message:      "test",
				HTTPCode:     http.StatusBadRequest,
				InternalCode: "FAILED",
				Err:          nil,
			},
			expectOk:  true,
			expectNil: false,
		},
		{
			name: "wrapped ResponseError",
			err: fmt.Errorf("wrapped: %w", &ResponseError{
				Message:      "test",
				HTTPCode:     http.StatusNotFound,
				InternalCode: "FAILED",
				Err:          nil,
			}),
			expectOk:  true,
			expectNil: false,
		},
		{
			name:      "not a ResponseError",
			err:       errors.New("standard error"),
			expectOk:  false,
			expectNil: true,
		},
		{
			name:      "nil error",
			err:       nil,
			expectOk:  false,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceErr, ok := getResponseError(tt.err)
			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got %v", tt.expectOk, ok)
			}
			if tt.expectNil && serviceErr != nil {
				t.Errorf("expected nil ResponseError, got %v", serviceErr)
			}
			if !tt.expectNil && serviceErr == nil {
				t.Error("expected non-nil ResponseError, got nil")
			}
		})
	}
}

func TestHandle(t *testing.T) {
	logger := zap.NewNop().Sugar()

	tests := []struct {
		name             string
		err              error
		expectedStatus   common.StatusCode
		expectedMessage  string
		expectedHTTPCode string
		checkMetadata    bool
	}{
		{
			name: "BadRequest error",
			err: &ResponseError{
				Message:      "invalid input",
				HTTPCode:     http.StatusBadRequest,
				InternalCode: "FAILED",
				Err:          nil,
			},
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "invalid input",
			expectedHTTPCode: "400",
			checkMetadata:    true,
		},
		{
			name: "NotFound error",
			err: &ResponseError{
				Message:      "Resource not found",
				HTTPCode:     http.StatusNotFound,
				InternalCode: "FAILED",
				Err:          nil,
			},
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "Resource not found",
			expectedHTTPCode: "404",
			checkMetadata:    true,
		},
		{
			name: "Internal error",
			err: &ResponseError{
				Message:      "unexpected error",
				HTTPCode:     http.StatusInternalServerError,
				InternalCode: "FAILED",
				Err:          nil,
			},
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "unexpected error",
			expectedHTTPCode: "500",
			checkMetadata:    true,
		},
		{
			name: "ServiceUnavailable error",
			err: &ResponseError{
				Message:      "service down",
				HTTPCode:     http.StatusServiceUnavailable,
				InternalCode: "FAILED",
				Err:          nil,
			},
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "service down",
			expectedHTTPCode: "503",
			checkMetadata:    true,
		},
		{
			name:             "Standard error",
			err:              errors.New("some error"),
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "internal server error",
			expectedHTTPCode: "500",
			checkMetadata:    true,
		},
		{
			name:             "gRPC Unimplemented error",
			err:              status.Error(codes.Unimplemented, "method Echo not implemented"),
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "method Echo not implemented",
			expectedHTTPCode: "501",
			checkMetadata:    true,
		},
		{
			name:             "gRPC Internal error keeps the generic message",
			err:              status.Error(codes.Internal, "db password is hunter2"),
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "internal server error",
			expectedHTTPCode: "500",
			checkMetadata:    true,
		},
		{
			name:             "ResponseError with zero HTTP code",
			err:              &ResponseError{Message: "test", HTTPCode: 0},
			expectedStatus:   common.StatusCode_FAILED,
			expectedMessage:  "test",
			expectedHTTPCode: "500",
			checkMetadata:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create context with metadata support
			ctx := metadata.NewOutgoingContext(context.Background(), metadata.MD{})

			response := handle(ctx, logger, tt.err)

			if response.Status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, response.Status)
			}
			if response.Message != tt.expectedMessage {
				t.Errorf("expected message %q, got %q", tt.expectedMessage, response.Message)
			}

			// Note: In a real gRPC server, we would be able to check the trailer metadata
			// For unit tests, we just verify that Handle doesn't panic when setting trailers
		})
	}
}

// trailerStream captures the trailers handle sets, so the HTTP code can be checked.
type trailerStream struct{ trailer metadata.MD }

func (s *trailerStream) Method() string                  { return "/test/Method" }
func (s *trailerStream) SetHeader(metadata.MD) error     { return nil }
func (s *trailerStream) SendHeader(metadata.MD) error    { return nil }
func (s *trailerStream) SetTrailer(md metadata.MD) error { s.trailer = metadata.Join(s.trailer, md); return nil }

// TestHandleHTTPCode pins the x-http-code each kind of error is answered with.
func TestHandleHTTPCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ResponseError", &ResponseError{Message: "bad", HTTPCode: http.StatusBadRequest}, "400"},
		{"Unimplemented", status.Error(codes.Unimplemented, "method Echo not implemented"), "501"},
		{"InvalidArgument", status.Error(codes.InvalidArgument, "bad purl"), "400"},
		{"NotFound", status.Error(codes.NotFound, "no such thing"), "404"},
		{"Internal", status.Error(codes.Internal, "boom"), "500"},
		{"plain error", errors.New("boom"), "500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := &trailerStream{}
			ctx := grpc.NewContextWithServerTransportStream(context.Background(), stream)
			handle(ctx, zap.NewNop().Sugar(), tt.err)
			if got := stream.trailer.Get("x-http-code"); len(got) != 1 || got[0] != tt.want {
				t.Errorf("x-http-code = %v, want %s", got, tt.want)
			}
		})
	}
}
