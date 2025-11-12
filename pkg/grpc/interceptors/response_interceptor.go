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
	"reflect"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc"

	common "github.com/scanoss/papi/api/commonv2"
)

// setStatusField uses reflection to set the Status field.
func setStatusField(resp interface{}, status *common.StatusResponse) {
	// Handle nil response
	if resp == nil {
		return
	}

	// Get the value
	v := reflect.ValueOf(resp)

	// Handle invalid value
	if !v.IsValid() {
		return
	}

	// If it's a pointer, get the element
	if v.Kind() == reflect.Ptr {
		// Check if pointer is nil
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	// Get the Status field
	statusField := v.FieldByName("Status")

	// Set it if possible
	if statusField.IsValid() && statusField.CanSet() {
		statusField.Set(reflect.ValueOf(status))
	}
}

// ResponseInterceptor is a simple interceptor that logs request information.
// Use this to verify that custom interceptors are working correctly.
func ResponseInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		s := ctxzap.Extract(ctx).Sugar()
		resp, err := handler(ctx, req)
		if err != nil {
			status := handle(ctx, s, err)
			setStatusField(resp, status)
			return resp, nil // Return nil so gateway uses our custom response format
		}
		return resp, err
	}
}
