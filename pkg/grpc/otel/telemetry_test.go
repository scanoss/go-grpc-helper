package otel

import (
	"testing"

	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestSetupTelemetry(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()

	oltpShutdown, err := InitTelemetryProviders("go-grpc-helper", "go-grpc", "0.0.1", "0.0.0.0:4317", GetTraceSampler("dev"), true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	oltpShutdown()
}

func TestSampleMode(t *testing.T) {
	GetTraceSampler("dev")
	GetTraceSampler("prod")
	GetTraceSampler("unknown")
}
