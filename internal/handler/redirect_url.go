package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"

	"github.com/go-chi/chi/v5"
)

// RedirectURL redirect by URLCode
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "URLCode")
	if code == "" {
		http.Error(w, "URLCode is empty", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	url, err := h.shortener.RedirectURL(ctx, code)
	if err != nil {
		if errors.Is(err, storage.ErrURLDeleted) {
			http.Error(w, "Status Gone", http.StatusGone)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
