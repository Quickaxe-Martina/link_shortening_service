package handler

import (
	"log"
	"net/http"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/go-chi/chi/v5"
)

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	URLCode := chi.URLParam(r, "URLCode")
	if len(URLCode) == 0 {
		log.Println("URLCode is empty")
		http.Error(w, "URLCode is empty", http.StatusBadRequest)
	}
	originalURL, exists := config.URLData[URLCode]
	if !exists {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
