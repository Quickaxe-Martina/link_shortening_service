package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const hostname = "http://localhost:8080/" // TODO: вынести в env
var urlData = make(map[string]string)     // TODO: БД в следующем спринте видимо?
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("🚀 Старт обработки %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("🏁 Завершено за %v", time.Since(start))
	})
}

func generateUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	urlCode := generateRandomString(6)
	log.Printf("Url code: %s", urlCode)
	urlData[urlCode] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(hostname + urlCode))
}

func redirectUrl(w http.ResponseWriter, r *http.Request) {
	urlCode := strings.TrimPrefix(r.URL.Path, "/")
	if len(urlCode) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	originalURL, exists := urlData[urlCode]
	if !exists {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		generateUrl(w, r)
	case http.MethodGet:
		redirectUrl(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Println("Server started")
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	loggedMux := loggingMiddleware(mux)

	err := http.ListenAndServe(`:8080`, loggedMux)
	if err != nil {
		panic(err)
	}
}
