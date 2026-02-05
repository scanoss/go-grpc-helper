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
