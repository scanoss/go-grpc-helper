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

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ResponseError struct {
	Message      string                 // Human-readable error message
	HTTPCode     int                    // HTTP status code to return to client
	InternalCode string                 // Internal error code for logging/monitoring
	Err          error                  // Wrapped original error for error chain
	Details      map[string]interface{} // Optional additional context
}

// error implements the error interface.
func (e *ResponseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// getHTTPCode returns the HTTP status code for this error with fallback to 500.
func (e *ResponseError) getHTTPCode() int {
	if e.HTTPCode == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTPCode
}

// isResponseError checks if an error is a ResponseError.
func isResponseError(err error) bool {
	var serviceErr *ResponseError
	return errors.As(err, &serviceErr)
}

// getResponseError extracts a ResponseError from an error chain.
func getResponseError(err error) (*ResponseError, bool) {
	var responseError *ResponseError
	if errors.As(err, &responseError) {
		return responseError, true
	}
	return nil, false
}

// Handle converts a ResponseError to a gRPC response with proper HTTP status.
func handle(ctx context.Context, s *zap.SugaredLogger, err error) *common.StatusResponse {
	var responseError *ResponseError
	if isResponseError(err) {
		responseError, _ = getResponseError(err)
		// Set HTTP trailer based on custom error
		trailerErr := grpc.SetTrailer(ctx, metadata.Pairs("x-http-code", fmt.Sprintf("%d", responseError.getHTTPCode())))
		if trailerErr != nil {
			s.Debugf("error setting x-http-code to trailer: %v", trailerErr)
		}

		// Log with structured data for monitoring
		s.Errorw("service error",
			"error", responseError.Error(),
			"http_code", responseError.getHTTPCode(),
			"internal_code", responseError.InternalCode,
			"details", responseError.Details,
		)

		return &common.StatusResponse{
			Status:  common.StatusCode_FAILED,
			Message: responseError.Message,
		}
	}

	// Default to 500 for unknown errors
	trailerErr := grpc.SetTrailer(ctx, metadata.Pairs("x-http-code", fmt.Sprintf("%d", http.StatusInternalServerError)))
	if trailerErr != nil {
		s.Debugf("error setting x-http-code to trailer: %v", trailerErr)
	}

	s.Errorw("unhandled error", "error", err.Error())

	return &common.StatusResponse{
		Status:  common.StatusCode_FAILED,
		Message: "internal server error",
	}
}
