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

// Package domain defines core domain types and mappings for the gRPC helper.
package domain

import (
	pb "github.com/scanoss/papi/api/commonv2"
)

// ComponentStatus represents the result status of a component operation,
// containing a human-readable message and a machine-readable status code.
type ComponentStatus struct {
	Message    string
	StatusCode StatusCode
}

// StatusCode represents the possible outcome codes for a component operation.
type StatusCode string

const (
	// ComponentNotFound indicates the requested component was not found.
	ComponentNotFound StatusCode = "COMPONENT_NOT_FOUND"
	// InvalidPurl indicates the provided PURL is malformed or invalid.
	InvalidPurl StatusCode = "INVALID_PURL"
	// ComponentWithoutInfo indicates the component exists but has no available information.
	ComponentWithoutInfo StatusCode = "COMPONENT_WITHOUT_INFO"
	// Success indicates the operation completed successfully.
	Success StatusCode = "SUCCESS"
	// InvalidSemver indicates the provided semantic version is invalid.
	InvalidSemver StatusCode = "INVALID_SEMVER"
)

// StatusCodeToErrorCode maps a domain StatusCode to its corresponding protobuf ErrorCode.
// Returns nil for Success or any unrecognized status code.
func StatusCodeToErrorCode(code StatusCode) *pb.ErrorCode {
	switch code {
	case InvalidPurl:
		return pb.ErrorCode_INVALID_PURL.Enum()
	case ComponentNotFound:
		return pb.ErrorCode_COMPONENT_NOT_FOUND.Enum()
	case InvalidSemver:
		return pb.ErrorCode_INVALID_SEMVER.Enum()
	case ComponentWithoutInfo:
		return pb.ErrorCode_NO_INFO.Enum()
	case Success:
		return nil
	default:
		return nil
	}
}
