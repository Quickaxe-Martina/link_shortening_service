package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"go.uber.org/zap"
)


// StatsResponse response
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// InternalStats returns the number of shortened URLs in the service and the number of users in the service
func (h *Handler) InternalStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	urls, err := h.store.GetURLsCount(ctx)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	users, err := h.store.GetUsersCount(ctx)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := StatsResponse{
		URLs:  urls,
		Users: users,
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}