package url

import (
	"testing"
)

func TestGetByIdReturnsExistingUrl(t *testing.T) {
	repo := NewRepository()
	url := &Url{Id: 1, Original: "https://example.com", Shortened: "abc"}
	repo.urls[1] = url

	result, err := repo.GetById(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != url {
		t.Errorf("Expected %v, got %v", url, result)
	}
}

func TestGetByIdReturnsNilForNonExistentId(t *testing.T) {
	repo := NewRepository()

	result, err := repo.GetById(999)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestGetByValueReturnsUrlWithMatchingShortenedValue(t *testing.T) {
	repo := NewRepository()
	url := &Url{Id: 1, Original: "https://example.com", Shortened: "abc"}
	repo.urls[1] = url

	result, err := repo.GetByValue("abc")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != url {
		t.Errorf("Expected %v, got %v", url, result)
	}
}

func TestGetByValueReturnsErrorForNonExistentValue(t *testing.T) {
	repo := NewRepository()

	result, err := repo.GetByValue("nonexistent")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestInsertAssignsIdWhenIdIsMinusOne(t *testing.T) {
	repo := NewRepository()
	url := &Url{Id: -1, Original: "https://example.com", Shortened: "abc"}

	result, err := repo.Insert(url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Id != 0 {
		t.Errorf("Expected Id to be 0, got %d", result.Id)
	}
}

func TestInsertKeepsExistingIdWhenNotMinusOne(t *testing.T) {
	repo := NewRepository()
	url := &Url{Id: 5, Original: "https://example.com", Shortened: "abc"}

	result, err := repo.Insert(url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Id != 5 {
		t.Errorf("Expected Id to be 5, got %d", result.Id)
	}
}

func TestInsertStoresUrlInRepository(t *testing.T) {
	repo := NewRepository()
	url := &Url{Id: 1, Original: "https://example.com", Shortened: "abc"}

	insert, err := repo.Insert(url)
	if err != nil {
		t.Errorf("Failed to insert record")
	}

	if insert == nil {
		t.Errorf("Expected non-nil result, got nil")
	}

	stored := repo.urls[1]
	if stored != url {
		t.Errorf("Expected %v to be stored, got %v", url, stored)
	}
}

func TestUpdateModifiesExistingUrl(t *testing.T) {
	repo := NewRepository()
	original := &Url{Id: 1, Original: "https://example.com", Shortened: "abc"}
	repo.urls[1] = original

	err := repo.Update(original)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if repo.urls[1].Visits != 1 {
		t.Errorf("Expected %v, got %v", 1, repo.urls[1].Visits)
	}
}

func TestUpdateReturnsNilForNonExistentUrl(t *testing.T) {
	repo := NewRepository()
	url := &Url{Id: 999, Original: "https://example.com", Shortened: "abc"}

	err := repo.Update(url)
	if err == nil {
		t.Errorf("Expected error, got none")
	}
}

func TestNextReturnsCurrentMapSize(t *testing.T) {
	repo := NewRepository()
	repo.urls[1] = &Url{Id: 1, Original: "https://example.com", Shortened: "abc"}
	repo.urls[2] = &Url{Id: 2, Original: "https://test.com", Shortened: "def"}

	result, err := repo.Next()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}

func TestNextReturnsZeroForEmptyRepository(t *testing.T) {
	repo := NewRepository()

	result, err := repo.Next()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestNewRepositoryCreatesEmptyRepository(t *testing.T) {
	repo := NewRepository()

	if repo.urls == nil {
		t.Errorf("Expected urls map to be initialized")
	}
	if len(repo.urls) != 0 {
		t.Errorf("Expected empty urls map, got %d items", len(repo.urls))
	}
}
