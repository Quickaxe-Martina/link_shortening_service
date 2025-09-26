package storage

import (
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
)

// NewStorage выбирает и возвращает реализацию интерфейса Storage
func NewStorage(cfg *config.Config) (Storage, error) {
	if cfg.DatabaseDsn != "" {
		return NewPostgresStorage(cfg), nil
	}
	if cfg.DataFilePath != "" {
		return NewMemoryStorage(cfg, true), nil
	}
	return NewMemoryStorage(cfg, false), nil
}
