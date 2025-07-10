package url

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepository struct {
	urls map[int]*Url
}

func (m *mockRepository) GetById(id int) (*Url, error) {
	url, exists := m.urls[id]
	if !exists {
		return nil, nil
	}
	return url, nil
}

func (m *mockRepository) GetByValue(value string) (*Url, error) {
	for _, url := range m.urls {
		if url.Shortened == value {
			return url, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockRepository) Insert(url *Url) (*Url, error) {
	if url.Id == -1 {
		url.Id = len(m.urls)
	}
	m.urls[url.Id] = url
	return url, nil
}

func (m *mockRepository) Update(url *Url) error {
	if _, exists := m.urls[url.Id]; exists {
		m.urls[url.Id].Visits += 1
	}
	return nil
}

func (m *mockRepository) Next() (int, error) {
	return len(m.urls), nil
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		urls: make(map[int]*Url),
	}
}

func TestHandleUrlShortenReturnsSuccessfulResponse(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	shortLink := ShortLink{Url: "https://example.com"}
	body, _ := json.Marshal(shortLink)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	service.handleUrlShorten(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result ShortenedLink
	json.NewDecoder(w.Body).Decode(&result)

	if result.Result != "http://localhost:8080/a" {
		t.Errorf("http://localhost:8080/a, got %s", result.Result)
	}
}

func TestHandleUrlShortenReturnsBadRequestForInvalidJson(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	service.handleUrlShorten(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleUrlShortenSetsCorrectContentType(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	shortLink := ShortLink{Url: "https://example.com"}
	body, _ := json.Marshal(shortLink)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	service.handleUrlShorten(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got %s", contentType)
	}
}

func TestHandleUrlShortenCreatesCorrectShortenedUrl(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	shortLink := ShortLink{Url: "https://example.com"}
	body, _ := json.Marshal(shortLink)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	service.handleUrlShorten(w, req)

	var result ShortenedLink
	json.NewDecoder(w.Body).Decode(&result)

	if result.Result != "http://localhost:8080/a" {
		t.Errorf("Expected shortened URL to be 'http://localhost:8080/a', got %s", result.Result)
	}
}

func TestRegisterHandlersCreatesCorrectRoute(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	router := mux.NewRouter()
	service.RegisterHandlers(router)

	shortLink := ShortLink{Url: "https://example.com"}
	body, _ := json.Marshal(shortLink)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNewServiceCreatesServiceWithCorrectFields(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	if service.repository != repo {
		t.Errorf("Expected repository to be set correctly")
	}

	if service.port != ":8080" {
		t.Errorf("Expected port to be ':8080', got %s", service.port)
	}

	if service.redirectUrl != "http://localhost" {
		t.Errorf("Expected redirectUrl to be 'http://localhost', got %s", service.redirectUrl)
	}

	if service.apiPrefix != "api" {
		t.Errorf("Expected apiPrefix to be 'api', got %s", service.apiPrefix)
	}

	if service.apiVersion != 1 {
		t.Errorf("Expected apiVersion to be 1, got %d", service.apiVersion)
	}
}

func TestHandleUrlRedirectNotFound(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	req := httptest.NewRequest(http.MethodGet, "/abc/", nil)
	w := httptest.NewRecorder()

	service.handleUrlRedirect(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandleUrlRedirectIncrementsVisits(t *testing.T) {
	repo := newMockRepository()
	repo.urls[1] = &Url{Id: 1, Original: "https://example.com", Shortened: "abc", Visits: 0}
	service := New(repo, ":8080", "http://localhost", "api", 1)

	req := httptest.NewRequest(http.MethodGet, "/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"shortened": "abc"})
	w := httptest.NewRecorder()

	service.handleUrlRedirect(w, req)

	if repo.urls[1].Visits != 1 {
		t.Errorf("Expected visits to be incremented to 1, got %d", repo.urls[1].Visits)
	}
	if w.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", w.Code)
	}
}

func TestHandleStatsReturnsVisits(t *testing.T) {
	repo := newMockRepository()
	repo.urls[1] = &Url{Id: 1, Original: "https://example.com", Shortened: "abc", Visits: 5}
	service := New(repo, ":8080", "http://localhost", "api", 1)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/1", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	service.handleStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]int
	json.NewDecoder(w.Body).Decode(&result)
	if result["visits"] != 5 {
		t.Errorf("Expected visits to be 5, got %d", result["visits"])
	}
}

func TestHandleStatsNotFound(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/999", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "999"}
	req = mux.SetURLVars(req, vars)

	service.handleStats(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandleStatsInvalidId(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/invalid", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "invalid"}
	req = mux.SetURLVars(req, vars)

	service.handleStats(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleUrlShortenReturnsBadRequestForEmptyUrl(t *testing.T) {
	repo := newMockRepository()
	service := New(repo, ":8080", "http://localhost", "api", 1)

	shortLink := ShortLink{Url: ""}
	body, _ := json.Marshal(shortLink)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	service.handleUrlShorten(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
