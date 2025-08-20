package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"
)

// GenerateURL handles HTTP requests to create a shortened URL.
func (h *Handler) GenerateURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	URLCode, err := service.GenerateRandomString(6)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("URL code: %s", URLCode)
	h.cfg.URLData[URLCode] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.cfg.ServerAddr + URLCode))
}
