package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RedirectURL redirect by URLCode
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	URLCode := chi.URLParam(r, "URLCode")
	if len(URLCode) == 0 {
		log.Println("URLCode is empty")
		http.Error(w, "URLCode is empty", http.StatusBadRequest)
	}
	originalURL, exists := h.cfg.URLData[URLCode]
	if !exists {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
