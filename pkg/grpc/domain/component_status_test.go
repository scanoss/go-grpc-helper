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

package domain

import (
	"fmt"
	pb "github.com/scanoss/papi/api/commonv2"
	"testing"
)

func TestStatusCodeToErrorCodeBuilder(t *testing.T) {

	tests := []struct {
		name     string
		input    StatusCode
		expected *pb.ErrorCode
	}{
		{
			name:     "Should_Component_NoFound",
			input:    ComponentNotFound,
			expected: pb.ErrorCode_COMPONENT_NOT_FOUND.Enum(),
		},
		{
			name:     "Should_ReturnInvalidPurl_WhenInputIsInvalidPurl",
			input:    InvalidPurl,
			expected: pb.ErrorCode_INVALID_PURL.Enum(),
		},
		{
			name:     "Should_ReturnNoInfo_WhenInputIsComponentWithoutInfo",
			input:    ComponentWithoutInfo,
			expected: pb.ErrorCode_NO_INFO.Enum(),
		},
		{
			name:     "Should_ReturnInvalidSemver_WhenInputIsInvalidSemver",
			input:    InvalidSemver,
			expected: pb.ErrorCode_INVALID_SEMVER.Enum(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StatusCodeToErrorCode(tt.input)
			fmt.Printf("Result %v", result)
			fmt.Printf("Expected %v", tt.expected)
			// Handle nil cases
			if *result != *tt.expected {
				t.Errorf("Expected %v, received %v", *tt.expected, *result)
			}
		})
	}

}
