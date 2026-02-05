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

package responseerror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/scanoss/go-grpc-helper/pkg/grpc/interceptors"
)

func TestNewBadRequestError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
	}{
		{
			name:    "with wrapped error",
			message: "invalid input",
			err:     errors.New("field missing"),
		},
		{
			name:    "without wrapped error",
			message: "invalid input",
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceErr := BadRequest(tt.message, tt.err)

			if serviceErr.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, serviceErr.Message)
			}
			if serviceErr.HTTPCode != http.StatusBadRequest {
				t.Errorf("expected HTTP code %d, got %d", http.StatusBadRequest, serviceErr.HTTPCode)
			}
			if serviceErr.InternalCode != "BAD_REQUEST" {
				t.Errorf("expected internal code %q, got %q", "BAD_REQUEST", serviceErr.InternalCode)
			}
			if !errors.Is(serviceErr.Err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, serviceErr.Err)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected string
	}{
		{
			name:     "ecosystem not found",
			resource: "Ecosystem",
			expected: "Ecosystem not found",
		},
		{
			name:     "dependency not found",
			resource: "Dependency",
			expected: "Dependency not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceErr := NotFound(tt.resource)

			if serviceErr.Message != tt.expected {
				t.Errorf("expected message %q, got %q", tt.expected, serviceErr.Message)
			}
			if serviceErr.HTTPCode != http.StatusNotFound {
				t.Errorf("expected HTTP code %d, got %d", http.StatusNotFound, serviceErr.HTTPCode)
			}
			if serviceErr.InternalCode != "NOT_FOUND" {
				t.Errorf("expected internal code %q, got %q", "NOT_FOUND", serviceErr.InternalCode)
			}
			if serviceErr.Err != nil {
				t.Errorf("expected error to be nil, got %v", serviceErr.Err)
			}
		})
	}
}

func TestNewInternalError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
	}{
		{
			name:    "with wrapped error",
			message: "unexpected error occurred",
			err:     errors.New("database error"),
		},
		{
			name:    "without wrapped error",
			message: "unexpected error occurred",
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceErr := InternalServerError(tt.message, tt.err)

			if serviceErr.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, serviceErr.Message)
			}
			if serviceErr.HTTPCode != http.StatusInternalServerError {
				t.Errorf("expected HTTP code %d, got %d", http.StatusInternalServerError, serviceErr.HTTPCode)
			}
			if serviceErr.InternalCode != "INTERNAL_ERROR" {
				t.Errorf("expected internal code %q, got %q", "INTERNAL_ERROR", serviceErr.InternalCode)
			}
			if !errors.Is(serviceErr.Err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, serviceErr.Err)
			}
		})
	}
}

func TestNewServiceUnavailableError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
	}{
		{
			name:    "database down",
			message: "database unavailable",
			err:     errors.New("connection timeout"),
		},
		{
			name:    "external service timeout",
			message: "service timeout",
			err:     errors.New("timeout exceeded"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceErr := ServiceUnavailable(tt.message, tt.err)

			if serviceErr.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, serviceErr.Message)
			}
			if serviceErr.HTTPCode != http.StatusServiceUnavailable {
				t.Errorf("expected HTTP code %d, got %d", http.StatusServiceUnavailable, serviceErr.HTTPCode)
			}
			if serviceErr.InternalCode != "SERVICE_UNAVAILABLE" {
				t.Errorf("expected internal code %q, got %q", "SERVICE_UNAVAILABLE", serviceErr.InternalCode)
			}
			if !errors.Is(serviceErr.Err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, serviceErr.Err)
			}
		})
	}
}

func TestServiceError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *interceptors.ResponseError
		expected string
	}{
		{
			name: "with wrapped error",
			err: &interceptors.ResponseError{
				Message: "invalid input",
				Err:     errors.New("field missing"),
			},
			expected: "invalid input: field missing",
		},
		{
			name: "without wrapped error",
			err: &interceptors.ResponseError{
				Message: "not found",
				Err:     nil,
			},
			expected: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("expected error string %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestServiceError_WithDetails(t *testing.T) {
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid",
	}

	err := &interceptors.ResponseError{
		Details: details,
	}

	if err.Details["field"] != "email" {
		t.Errorf("expected field=email, got %v", err.Details["field"])
	}
	if err.Details["value"] != "invalid" {
		t.Errorf("expected value=invalid, got %v", err.Details["value"])
	}
}

func TestServiceError_ErrorChaining(t *testing.T) {
	originalErr := errors.New("original error")
	serviceErr := InternalServerError("service failed", originalErr)
	wrappedErr := fmt.Errorf("wrapper: %w", serviceErr)

	// Test that errors.As works with ResponseError
	var targetErr *interceptors.ResponseError
	if !errors.As(wrappedErr, &targetErr) {
		t.Error("expected errors.As to find ResponseError in chain")
	}

	if targetErr.Message != "service failed" {
		t.Errorf("expected message 'service failed', got %q", targetErr.Message)
	}

	// Test that we can access the original error via the ResponseError's Err field
	if !errors.Is(originalErr, targetErr.Err) {
		t.Errorf("expected wrapped error to be %v, got %v", originalErr, targetErr.Err)
	}

	// Verify the error message includes both the service message and the original error
	expectedErrorMsg := "service failed: original error"
	if serviceErr.Error() != expectedErrorMsg {
		t.Errorf("expected error message %q, got %q", expectedErrorMsg, serviceErr.Error())
	}
}
