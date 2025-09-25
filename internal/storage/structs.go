package storage

import (
	"context"
)

// URL code and original value
type URL struct {
	Code string
	URL  string
}

// URLStorage defines methods for saving and retrieving URLs
type URLStorage interface {
	SaveURL(ctx context.Context, u URL) error
	GetURL(ctx context.Context, code string) (URL, error)
	AllURLs(ctx context.Context) ([]URL, error)
	SaveBatchURL(ctx context.Context, urls []URL) error
}

// Storage defines methods
type Storage interface {
	URLStorage
	Close() error
	Ping(ctx context.Context) error
}
