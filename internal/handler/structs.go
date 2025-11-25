package handler

import (
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/repository"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
)

// Handler data
type Handler struct {
	cfg          *config.Config
	store        storage.Storage
	deleteWorker *repository.DeleteURLsWorkers
	audit        *repository.AuditPublisher
}

// NewHandler create Handler
func NewHandler(cfg *config.Config, store storage.Storage, deleteWorker *repository.DeleteURLsWorkers, audit *repository.AuditPublisher) *Handler {
	return &Handler{
		cfg:          cfg,
		store:        store,
		deleteWorker: deleteWorker,
		audit:        audit,
	}
}
