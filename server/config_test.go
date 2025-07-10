package server

import (
	"github.com/rs/zerolog"
	"testing"
)

func TestLoadsConfigSuccessfully(t *testing.T) {
	config, err := LoadConfig("../.env")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if config.Port == "" {
		t.Errorf("Expected Port to be set, got empty string")
	}
	if config.ApiPrefix == "" {
		t.Errorf("Expected ApiPrefix to be set, got empty string")
	}
	if config.ApiVersion == 0 {
		t.Errorf("Expected ApiVersion to be set, got 0")
	}
	if config.LogLevel != zerolog.InfoLevel {
		t.Errorf("Expected LogLevel to be InfoLevel, got %v", config.LogLevel)
	}
}
