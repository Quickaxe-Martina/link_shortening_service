package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/auth"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"go.uber.org/zap"
)

type ctxKey string

const userKey = ctxKey("user")

// GetOrCreateUserMiddleware get or create user middware
func (h *Handler) GetOrCreateUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenExp := time.Hour * time.Duration(h.cfg.TokenExp)
		user, err := auth.GetOrCreateUser(w, r, h.store, h.cfg.SecretKey, tokenExp)
		if err != nil {
			logger.Log.Error("error get or create user", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser give user from context
func GetUser(ctx context.Context) storage.User {
	return ctx.Value(userKey).(storage.User)
}
