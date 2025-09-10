package handler

import (
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
)

// Handler data
type Handler struct {
	cfg         *config.Config
	storageData *storage.Storage
}

// NewHandler create Handler
func NewHandler(cfg *config.Config, storageData *storage.Storage) *Handler {
	return &Handler{
		cfg:         cfg,
		storageData: storageData,
	}
}
