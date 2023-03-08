// SPDX-License-Identifier: MIT
/*
 * Copyright (c) 2023, SCANOSS
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
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	"golang.org/x/net/context"
)

func TestGatewaySetupNoTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs := []string{"127.0.0.1"}
	deniedIPs := []string{"192.168.0.1"}
	srv, mux, gateway, opts, err := SetupGateway("9443", "8443", "", allowedIPs, deniedIPs, true, false, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if srv == nil || mux == nil || len(gateway) == 0 || len(opts) == 0 {
		t.Errorf("Missing settings: %v, %v, %v, %v", srv, mux, gateway, opts)
	}
}

func TestGatewaySetupTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs := []string{"127.0.0.1"}
	deniedIPs := []string{"192.168.0.1"}
	srv, mux, gateway, opts, err := SetupGateway(":9443", "8443", "../../../tests/server.crt", allowedIPs, deniedIPs, true, false, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if srv == nil || mux == nil || len(gateway) == 0 || len(opts) == 0 {
		t.Errorf("Missing settings: %v, %v, %v, %v", srv, mux, gateway, opts)
	}
}

func TestStartGatewayNoTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	srv := &http.Server{
		Addr:              ":0",
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           runtime.NewServeMux(),
	}
	go func() {
		time.Sleep(3 * time.Second)
		err = srv.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}()
	StartGateway(srv, "", "", false)
}

func TestStartGatewayTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	srv := &http.Server{
		Addr:              ":0",
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           runtime.NewServeMux(),
	}
	go func() {
		time.Sleep(3 * time.Second)
		err = srv.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}()
	StartGateway(srv, "../../../tests/server.crt", "../../../tests/server.key", true)
}
