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
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/utils"
	"github.com/scanoss/ipfilter/v2"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// SetupGateway configures and returns an HTTP server that acts as a gateway to a gRPC service.
// The gateway is forced to connect to localhost regardless of the provided grpcPort hostname.
//
// Important note about localhost and certificates:
// The gateway always establishes its connection to the gRPC server through localhost
// (e.g., localhost:50051). Therefore, if the TLS certificate does not include "localhost"
// in its subject/SAN fields, the connection will fail with a certificate validation error.
// The commonName parameter allows you to override the hostname verification in such cases.
//
// For example:
//   - If your certificate is issued for "api.example.com" without "localhost" in SAN:
//     Set commonName="api.example.com" to match the certificate's subject.
//   - If your certificate includes "localhost" in SAN:
//     Set commonName="localhost" (or it can be left empty as "localhost" is the default).
func SetupGateway(grpcPort, httpPort, tlsCertFile, commonName string, allowedIPs, deniedIPs []string,
	blockByDefault, trustProxy, startTLS bool) (*http.Server, *runtime.ServeMux, string, []grpc.DialOption, error) {
	httpPort = utils.SetupPort(httpPort)
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitDefaultValues: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		runtime.WithForwardResponseOption(httpSuccessResponseModifier),
		runtime.WithErrorHandler(httpErrorResponseModifier),
	)
	srv := &http.Server{
		Addr:              httpPort,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           mux,
	}
	if len(allowedIPs) > 0 || len(deniedIPs) > 0 { // Configure the list of allowed/denied IPs to connect
		zlog.S.Debugf("Filtering requests by allowed: %v, denied: %v, block-by-default: %v, trust-proxy: %v",
			allowedIPs, deniedIPs, blockByDefault, trustProxy)
		handler := ipfilter.Wrap(mux, ipfilter.Options{AllowedIPs: allowedIPs, BlockedIPs: deniedIPs,
			BlockByDefault: blockByDefault, TrustProxy: trustProxy,
		})
		srv.Handler = handler // assign the filtered handler
	}
	var opts []grpc.DialOption
	if startTLS {
		creds, err := credentials.NewClientTLSFromFile(tlsCertFile, commonName)
		if err != nil {
			zlog.S.Errorf("Problem loading TLS file: %s - %v", tlsCertFile, err)
			return nil, nil, "", nil, fmt.Errorf("failed to load TLS credentials from file")
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
	// force the gateway to localhost
	var grpcGateway string
	if strings.Contains(grpcPort, ":") { // gRPC port has a hostname in it
		grpcGateway = "localhost:" + grpcPort[strings.LastIndex(grpcPort, ":")+1:]
	} else {
		grpcGateway = "localhost:" + grpcPort
	}
	return srv, mux, grpcGateway, opts, nil
}

func StartGateway(srv *http.Server, tlsCertFile, tlsKeyFile string, startTLS bool) {
	var httpErr error
	if startTLS {
		zlog.S.Infof("starting REST server with TLS on %v ...", srv.Addr)
		httpErr = srv.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	} else {
		zlog.S.Infof("starting REST server on %v ...", srv.Addr)
		httpErr = srv.ListenAndServe()
	}
	if httpErr != nil && fmt.Sprintf("%s", httpErr) != "http: Server closed" {
		zlog.S.Panicf("issue encountered when starting service: %v", httpErr)
	}
}

// httpSuccessResponseModifier is called for all successful gRPC responses (err == nil).
// It checks the x-http-code trailer and sets the appropriate HTTP status code.
// This allows the middleware to set custom HTTP codes even when returning err == nil.
func httpSuccessResponseModifier(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		fmt.Printf("httpSuccessResponseModifier: No server metadata found\n")
		return nil
	}
	// Check for custom HTTP status code in trailer
	if vals := md.TrailerMD.Get("x-http-code"); len(vals) > 0 {
		fmt.Printf("httpSuccessResponseModifier: Found x-http-code: %s\n", vals[0])
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			fmt.Printf("httpSuccessResponseModifier: Error parsing x-http-code: %v\n", err)
			return nil
		}
		fmt.Printf("httpSuccessResponseModifier: Setting HTTP status code: %d\n", code)
		w.WriteHeader(code)
	} else {
		fmt.Printf("httpSuccessResponseModifier: No x-http-code trailer found, using default\n")
	}
	return nil
}

// httpResponseModifier sets the HTTP status code based on the "x-http-code" trailer
// in the gRPC response metadata. This allows gRPC services to control the HTTP status
// when exposed via the gateway. This is called only for error responses.
func httpErrorResponseModifier(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		fmt.Printf("httpResponseModifier: No server metadata found\n")
	}
	// Check for custom HTTP status code
	if vals := md.TrailerMD.Get("x-http-code"); len(vals) > 0 {
		fmt.Printf("httpResponseModifier: Found x-http-code: %s\n", vals[0])
		code, parseErr := strconv.Atoi(vals[0])
		if parseErr != nil {
			fmt.Printf("httpResponseModifier: Error parsing x-http-code: %v\n", parseErr)
		}
		fmt.Printf("httpResponseModifier: Setting HTTP status code: %d\n", code)
		w.WriteHeader(code)
	} else {
		fmt.Printf("httpResponseModifier: No x-http-code header found\n")
	}
}
