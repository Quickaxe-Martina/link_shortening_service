package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"

	"github.com/go-chi/chi/v5"
)

// RedirectURL redirect by URLCode
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	URLCode := chi.URLParam(r, "URLCode")
	if len(URLCode) == 0 {
		logger.Log.Info("URLCode is empty")
		http.Error(w, "URLCode is empty", http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	url, err := h.store.GetURL(ctx, URLCode)
	if err != nil {
		if errors.Is(err, storage.ErrURLDeleted) {
			http.Error(w, "Status Gone", http.StatusGone)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		return
	}
	w.Header().Set("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
