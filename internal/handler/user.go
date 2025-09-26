package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/model"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"

	"go.uber.org/zap"
)

// GetUserURLs handles HTTP requests to retrieve all URLs associated with the authenticated user.
func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	var responses []model.UserURLsResponse
	user, err := service.GetUserByCookie(r, h.cfg.SecretKey)
	if err != nil {
		if errors.Is(err, service.ErrNoJWTInCookie) || errors.Is(err, service.ErrInvalidJWTToken) {
			user, err = service.GetOrCreateUser(w, r, h.store, h.cfg.SecretKey, time.Hour*time.Duration(h.cfg.TokenExp))
			if err != nil {
				logger.Log.Error("error get or create user", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	urls, err := h.store.GetURLsByUserID(r.Context(), user.ID)
	if err != nil {
		logger.Log.Error("error get urls by user id", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, url := range urls {
		responses = append(responses, model.UserURLsResponse{
			ShortURL:    h.cfg.ServerAddr + url.Code,
			OriginalURL: url.URL,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(responses); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}
