package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/model"

	"go.uber.org/zap"
)

// GetUserURLs handles HTTP requests to retrieve all URLs associated with the authenticated user.
func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	var responses []model.UserURLsResponse
	user := GetUser(r.Context())

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

// DeleteUserURLs delete users's urls
func (h *Handler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	var codes []string
	user := GetUser(r.Context())
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&codes); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.deleteWorker.AddTask(user.ID, codes)
	w.WriteHeader(http.StatusAccepted)
}
