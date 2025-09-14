package storage

import (
	"context"
	"errors"
	"maps"
	"slices"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
)

// MemoryStorage is an in-memory implementation of the Storage interface
type MemoryStorage struct {
	Urls         map[string]URL
	UseFile      bool
	DataFilePath string
}

// NewMemoryStorage creates new MemoryStorage
func NewMemoryStorage(cfg *config.Config, useFile bool) *MemoryStorage {
	store := &MemoryStorage{
		Urls:         make(map[string]URL),
		UseFile:      useFile,
		DataFilePath: cfg.DataFilePath,
	}
	if useFile {
		LoadData(cfg.DataFilePath, store)
	}
	return store
}

// SaveURL save a URL by  code in memory
func (m *MemoryStorage) SaveURL(ctx context.Context, u URL) error {
	m.Urls[u.Code] = u
	return nil
}

// GetURL get URL by code from memory
func (m *MemoryStorage) GetURL(ctx context.Context, code string) (URL, error) {
	u, ok := m.Urls[code]
	if !ok {
		return URL{}, errors.New("user not found")
	}
	return u, nil
}

// Close releases resources
func (m *MemoryStorage) Close() error {
	if m.UseFile {
		SaveData(m.DataFilePath, m)
	}
	return nil
}

// Ping do nothing
func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

// AllURLs returns all URLs
func (m *MemoryStorage) AllURLs(ctx context.Context) ([]URL, error) {
	return slices.Collect(maps.Values(m.Urls)), nil
}
