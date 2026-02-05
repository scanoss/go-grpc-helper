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
	"fmt"
	"net/http"

	"github.com/scanoss/go-grpc-helper/pkg/grpc/interceptors"
)

// ResponseError represents a service-level error with HTTP status mapping and additional context.

// BadRequest
// Use for: missing required fields, malformed input, invalid parameters.
func BadRequest(message string, err error) *interceptors.ResponseError {
	return &interceptors.ResponseError{
		Message:      message,
		HTTPCode:     http.StatusBadRequest,
		InternalCode: "BAD_REQUEST",
		Err:          err,
	}
}

// NotFound
// Use for: ecosystem not found, dependencies not found, resource missing.
func NotFound(resource string) *interceptors.ResponseError {
	return &interceptors.ResponseError{
		Message:      fmt.Sprintf("%s not found", resource),
		HTTPCode:     http.StatusNotFound,
		InternalCode: "NOT_FOUND",
		Err:          nil,
	}
}

// InternalServerError
// Use for: unexpected errors, programming errors, unhandled exceptions.
func InternalServerError(message string, err error) *interceptors.ResponseError {
	return &interceptors.ResponseError{
		Message:      message,
		HTTPCode:     http.StatusInternalServerError,
		InternalCode: "INTERNAL_ERROR",
		Err:          err,
	}
}

// ServiceUnavailable
// Use for: database down, external service timeout, rate limits exceeded.
func ServiceUnavailable(message string, err error) *interceptors.ResponseError {
	return &interceptors.ResponseError{
		Message:      message,
		HTTPCode:     http.StatusServiceUnavailable,
		InternalCode: "SERVICE_UNAVAILABLE",
		Err:          err,
	}
}
