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

// Package server provides functions to configure, start and shutdown gRPC services
package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/utils"
	"github.com/scanoss/ipfilter/v2"
	"github.com/scanoss/zap-logging-helper/pkg/grpc/interceptor"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SetupGrpcServer configures the port, filtering & logging interceptors for a gRPC Server
func SetupGrpcServer(port, tlsCertFile, tlsKeyFile string, allowedIPs, deniedIPs []string, startTLS, blockedByDefault,
	trustProxy bool) (net.Listener, *grpc.Server, error) {
	port = utils.SetupPort(port)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return nil, nil, err
	}
	var interceptors []grpc.UnaryServerInterceptor
	// Configure the list of allowed/denied IPs to connect
	if len(allowedIPs) > 0 || len(deniedIPs) > 0 {
		ipFilter := ipfilter.New(ipfilter.Options{AllowedIPs: allowedIPs, BlockedIPs: deniedIPs,
			BlockByDefault: blockedByDefault, TrustProxy: trustProxy,
		})
		interceptors = append(interceptors, ipFilter.IPFilterUnaryServerInterceptor())
	}
	interceptors = append(interceptors, grpczap.UnaryServerInterceptor(zlog.L))
	interceptors = append(interceptors, interceptor.ContextPropagationUnaryServerInterceptor()) // Needs to be called after UnaryServerInterceptor to make sure the logger is set
	var opts []grpc.ServerOption
	if startTLS {
		creds, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
		if err != nil {
			zlog.S.Errorf("Problem loading TLS file: %s - %v", tlsCertFile, err)
			return nil, nil, fmt.Errorf("failed to load TLS credentials from file")
		}
		opts = append(opts, grpc.Creds(creds))
	}
	opts = append(opts, grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	// register service
	server := grpc.NewServer(opts...)

	return listen, server, nil
}

// StartGrpcServer starts the given gRPC server on the specified listener
func StartGrpcServer(listen net.Listener, server *grpc.Server, startTLS bool) {
	withTLS := ""
	if startTLS {
		withTLS = "with TLS "
	}
	zlog.S.Infof("starting gRPC server %son %v ...", withTLS, listen.Addr())
	httpErr := server.Serve(listen)
	if httpErr != nil && fmt.Sprintf("%s", httpErr) != "http: Server closed" {
		zlog.S.Panicf("issue encountered when starting service: %v", httpErr)
	}
}

// WaitServerComplete waits for a signal to terminate the
func WaitServerComplete(srv *http.Server, server *grpc.Server) error {
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c
	if srv != nil {
		zlog.S.Info("shutting down REST server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Set a deadline for gracefully shutting down
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			zlog.S.Warnf("error shutting down server %s", err)
			return fmt.Errorf("issue encountered while shutting down service")
		} else {
			zlog.S.Info("REST server gracefully stopped")
		}
	}
	if server != nil {
		zlog.S.Info("shutting down gRPC server...")
		server.GracefulStop()
		zlog.S.Info("gRPC server gracefully stopped")
	}
	return nil
}
