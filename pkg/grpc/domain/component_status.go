package domain

import (
	pb "github.com/scanoss/papi/api/commonv2"
)

type ComponentStatus struct {
	Message    string
	StatusCode StatusCode
}

type StatusCode string

const (
	ComponentNotFound    StatusCode = "COMPONENT_NOT_FOUND"
	InvalidPurl          StatusCode = "INVALID_PURL"
	ComponentWithoutInfo StatusCode = "COMPONENT_WITHOUT_INFO"
	Success              StatusCode = "SUCCESS"
	InvalidSemver        StatusCode = "INVALID_SEMVER"
)

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
