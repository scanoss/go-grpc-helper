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
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/utils"
	"github.com/scanoss/ipfilter/v2"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// SetupGateway configures.
func SetupGateway(grpcPort, httpPort, tlsCertFile, tlsKeyFile string, allowedIPs, deniedIPs []string,
	blockByDefault, trustProxy, startTLS bool, insecureSkipVerify bool) (*http.Server, *runtime.ServeMux, string, []grpc.DialOption, error) {
	httpPort = utils.SetupPort(httpPort)
	mux := runtime.NewServeMux()
	srv := &http.Server{
		Addr:              httpPort,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           mux,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
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
		var cred credentials.TransportCredentials
		var err error
		cred, err = credentials.NewClientTLSFromFile(tlsCertFile, "")
		if err != nil {
			zlog.S.Errorf("Problem loading TLS file: %s - %v", tlsCertFile, err)
			return nil, nil, "", nil, fmt.Errorf("failed to load TLS credentials from file")
		}

		if insecureSkipVerify == true {
			cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
			if err != nil {
				return nil, nil, "", nil, fmt.Errorf("failed to load TLS certificate and key: %v", err)
			}
			// Create custom TLS config that skips hostname validation
			config := &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			}
			cred = credentials.NewTLS(config)
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(cred)}
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
