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
	Users        map[int]User
	UseFile      bool
	DataFilePath string
}

// NewMemoryStorage creates new MemoryStorage
func NewMemoryStorage(cfg *config.Config, useFile bool) *MemoryStorage {
	store := &MemoryStorage{
		Urls:         make(map[string]URL),
		Users:        make(map[int]User),
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
		return URL{}, errors.New("url not found")
	}
	return u, nil
}

// GetByURL get URL by url from memory
func (m *MemoryStorage) GetByURL(ctx context.Context, url string) (URL, error) {
	for _, v := range m.Urls {
		if v.URL == url {
			return v, nil
		}
	}
	return URL{}, errors.New("url not found")
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

// SaveBatchURL save list of URL
func (m *MemoryStorage) SaveBatchURL(ctx context.Context, urls []URL) error {
	for _, url := range urls {
		m.SaveURL(ctx, url)
	}
	return nil
}

// CreateUser creates a new user and returns it
func (m *MemoryStorage) CreateUser(ctx context.Context) (User, error) {
	newID := len(m.Users) + 1
	newUser := User{ID: newID}
	m.Users[newID] = newUser
	return newUser, nil
}

// GetURLsByUserID returns all URLs associated with a specific user ID
func (m *MemoryStorage) GetURLsByUserID(ctx context.Context, userID int) ([]URL, error) {
	var result []URL
	for _, url := range m.Urls {
		if url.UserID == userID {
			result = append(result, url)
		}
	}
	return result, nil
}

// GetAllUsers returns all users
func (m *MemoryStorage) GetAllUsers(ctx context.Context) ([]User, error) {
	return slices.Collect(maps.Values(m.Users)), nil
}
