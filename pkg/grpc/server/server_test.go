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

package server

import (
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/scanoss/go-grpc-helper/pkg/grpc/otel"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestSetupGrpcServerNoTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs := []string{"127.0.0.1"}
	deniedIPs := []string{"192.168.0.1"}

	listen, server, err := SetupGrpcServer(":0", "", "", allowedIPs, deniedIPs, false, true, false, false, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("Listening on %v\n", listen.Addr().String())

	go func() {
		time.Sleep(3 * time.Second)
		server.GracefulStop()
	}()
	StartGrpcServer(listen, server, false)
}

func TestSetupGrpcServerWithTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs := []string{"127.0.0.1"}
	deniedIPs := []string{"192.168.0.1"}
	listen, server, err := SetupGrpcServer(":0", "../../../tests/server.crt", "../../../tests/server.key", allowedIPs, deniedIPs, true, true, false, false, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("Listening on %v\n", listen.Addr().String())
	go func() {
		time.Sleep(3 * time.Second)
		server.GracefulStop()
	}()
	StartGrpcServer(listen, server, true)
}

func TestSetupGrpcServerNoTLSShutdown(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs := []string{"127.0.0.1"}
	deniedIPs := []string{"192.168.0.1"}
	listen, server, err := SetupGrpcServer(":0", "", "", allowedIPs, deniedIPs, false, true, false, false, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	srv := &http.Server{
		Addr:              ":0",
		Handler:           runtime.NewServeMux(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		httpErr := srv.ListenAndServe()
		if httpErr != nil && fmt.Sprintf("%s", httpErr) != "http: Server closed" {
			t.Errorf("Unexpected error: %v", httpErr)
		}
	}()
	go func() {
		fmt.Printf("Listening on %v\n", listen.Addr().String())
		StartGrpcServer(listen, server, false)
	}()
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Sending signal...")
		killErr := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if killErr != nil {
			fmt.Printf("Problem running kill: %v", killErr)
		}
	}()
	fmt.Println("Waiting for the app to finish...")
	err = WaitServerComplete(srv, server)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSetupGrpcServerNoTLSTelemetry(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs := []string{"127.0.0.1"}
	deniedIPs := []string{"192.168.0.1"}

	otelShutdown, err := otel.InitTelemetryProviders("go-grpc-helper", "go-grpc", "0.0.1", "0.0.0.0:4317", otel.GetTraceSampler("dev"), true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	listen, server, err := SetupGrpcServer(":0", "", "", allowedIPs, deniedIPs, false, true, false, true, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("Listening on %v\n", listen.Addr().String())

	go func() {
		time.Sleep(3 * time.Second)
		server.GracefulStop()
	}()
	StartGrpcServer(listen, server, false)
	otelShutdown()
}
