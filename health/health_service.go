package health

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func New() *Service {
	return &Service{}
}

type Service struct{}

func (s Service) RegisterHandlers(mux *mux.Router) {
	mux.HandleFunc("/health", s.handleHealthCheck).Methods("GET")
}

func (s Service) handleHealthCheck(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(map[string]string{"status": "up"})
	if err != nil {
		return
	}
}
