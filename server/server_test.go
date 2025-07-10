package server

import (
	"context"
	"testing"
	"time"
)

func TestStartsServerSuccessfully(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := Config{Port: ":0"}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := Start(ctx, config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestReturnsErrorWhenServerFailsToStart(t *testing.T) {
	ctx := context.Background()
	config := Config{Port: "invalid-port"}

	err := Start(ctx, config)
	if err == nil {
		t.Fatalf("Expected error when server fails to Start, got nil")
	}
}

func TestShutsDownGracefullyOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := Config{Port: ":0"}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := Start(ctx, config)
	if err != nil {
		t.Errorf("Expected graceful shutdown, got error: %v", err)
	}
}
