package url

import (
	"testing"
)

var baseMapTest = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestReturnsFirstCharacterForZeroInput(t *testing.T) {
	result, err := baseConvert(0, baseMapTest)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "a" {
		t.Errorf("Expected 'a', got %s", result)
	}
}

func TestConvertsPositiveIntegerToBaseString(t *testing.T) {
	result, err := baseConvert(1, baseMapTest)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "b" {
		t.Errorf("Expected 'b', got %s", result)
	}
}

func TestConvertsLargeIntegerToBaseString(t *testing.T) {
	result, err := baseConvert(3844, baseMapTest)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result for large integer")
	}
	if result != "aab" {
		t.Errorf("Expected 'aab', got %s", result)
	}
}

func TestShortenURLCombinesRedirectUrlAndConvertedId(t *testing.T) {
	result, err := ShortenURL(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "b" {
		t.Errorf("Expected 'b', got %s", result)
	}
}

func TestShortenURLHandlesZeroId(t *testing.T) {
	result, err := ShortenURL(0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "a" {
		t.Errorf("Expected 'a', got %s", result)
	}
}

func TestShortenURLHandlesLargeId(t *testing.T) {
	result, err := ShortenURL(100)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == "https://example.com/" {
		t.Errorf("Expected shortened URL with encoded id, got %s", result)
	}
}
