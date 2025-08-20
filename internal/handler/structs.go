package handler

import (
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
)

// Handler data
type Handler struct {
	cfg *config.Config
}

// NewHandler create Handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}
