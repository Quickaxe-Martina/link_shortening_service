package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const hostname = "http://localhost:8080/" // TODO: вынести в env
var URLData = make(map[string]string)     // TODO: БД в следующем спринте видимо?
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func generateURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	URLCode := generateRandomString(6)
	log.Printf("URL code: %s", URLCode)
	URLData[URLCode] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(hostname + URLCode))
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	URLCode := chi.URLParam(r, "URLCode")
	if len(URLCode) == 0 {
		log.Println("URLCode is empty")
		http.Error(w, "URLCode is empty", http.StatusBadRequest)
	}
	originalURL, exists := URLData[URLCode]
	if !exists {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/{URLCode}", redirectURL)
		r.Post("/", generateURL)
	})
	return r
}

func main() {
	log.Println("Server started")
	r := setupRouter()

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
