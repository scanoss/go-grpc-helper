// SPDX-License-Identifier: MIT
/*
 * Copyright (c) 2026, SCANOSS
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

package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func TestRegisterHealthEndpoint(t *testing.T) {
	mux := runtime.NewServeMux()
	if err := RegisterHealthEndpoint(mux); err != nil {
		t.Fatalf("unexpected error registering health endpoint: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, HealthPath, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	if body := rec.Body.String(); body != `{"status":"ok"}` {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestRegisterHealthEndpointNilMux(t *testing.T) {
	if err := RegisterHealthEndpoint(nil); err == nil {
		t.Error("expected an error for a nil mux, got nil")
	}
}

// GET-only, so it can't be confused with the POST-only Echo RPCs.
func TestRegisterHealthEndpointRejectsPost(t *testing.T) {
	mux := runtime.NewServeMux()
	if err := RegisterHealthEndpoint(mux); err != nil {
		t.Fatalf("unexpected error registering health endpoint: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, HealthPath, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Errorf("expected POST /health to not return 200, got %d", rec.Code)
	}
}
