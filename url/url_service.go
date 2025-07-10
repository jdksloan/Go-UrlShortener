package url

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
	"thesilentcoder.com/m/repository"
)

type ShortLink struct {
	Url string `json:"url"`
}

type Url struct {
	Id        int
	Original  string
	Shortened string
	Url       string
	Visits    int
}

type ShortenedLink struct {
	Result string
}

func New(repository repository.Repository[Url], port string, redirectUrl string, apiPrefix string, apiVersion int) *Service {
	return &Service{repository, port, redirectUrl, apiPrefix, apiVersion}
}

type Service struct {
	repository  repository.Repository[Url]
	port        string
	redirectUrl string
	apiPrefix   string
	apiVersion  int
}

func (s Service) RegisterHandlers(router *mux.Router) {
	formattedUrl := fmt.Sprintf("/%s/v%d/", s.apiPrefix, s.apiVersion)
	router.HandleFunc(formattedUrl+"shorten", s.handleUrlShorten).Methods("POST")
	router.HandleFunc("/{shortened}/", s.handleUrlRedirect).Methods("GET")

	router.HandleFunc(formattedUrl+"stats/{id}", s.handleStats).Methods("GET")
}

func (s Service) handleUrlShorten(writer http.ResponseWriter, r *http.Request) {
	var short ShortLink
	err := json.NewDecoder(r.Body).Decode(&short)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if short.Url == "" {
		http.Error(writer, "URL cannot be empty", http.StatusBadRequest)
		return
	}
	parsedUrl, err := url.ParseRequestURI(short.Url)
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		http.Error(writer, "Invalid URL format", http.StatusBadRequest)
		return
	}

	next, err := s.repository.Next()
	if err != nil {
		http.Error(writer, "Failed to get next ID", http.StatusInternalServerError)
		return
	}

	shortened, err := ShortenURL(next)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	redirect := fmt.Sprintf("%s%s/%s", s.redirectUrl, s.port, shortened)
	u := Url{
		Id:        next,
		Original:  short.Url,
		Shortened: shortened,
		Visits:    0,
		Url:       redirect,
	}
	ret, err := s.repository.Insert(&u)
	if err != nil {
		http.Error(writer, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	shortLink := ShortenedLink{ret.Url}
	err = json.NewEncoder(writer).Encode(shortLink)
}

func (s Service) handleUrlRedirect(writer http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortened := params["shortened"]

	byValue, err := s.repository.GetByValue(shortened)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	err = s.repository.Update(byValue)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, r, byValue.Original, http.StatusFound)

}

func (s Service) handleStats(writer http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}
	res, err := s.repository.GetById(id)

	if res == nil {
		http.Error(writer, "URL not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(writer).Encode(map[string]int{"visits": res.Visits})

}
