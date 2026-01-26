/*
Package storage for storages
*/
package storage

import (
	"context"
	"errors"
)

// ErrURLAlreadyExists url is already saved in DB
var ErrURLAlreadyExists = errors.New("url already exists")

// ErrCodeAlreadyExists code is already taken
var ErrCodeAlreadyExists = errors.New("url already taken")

// ErrURLDeleted code is already deleted
var ErrURLDeleted = errors.New("url has deleted")

// ErrNotImplemented not implemented
var ErrNotImplemented = errors.New("not implemented")

// URL code and original value
// generate:reset
type URL struct {
	Code      string
	URL       string
	UserID    int
	isDeleted bool
}

// User model
// generate:reset
type User struct {
	ID int
}

// URLStorage defines methods for saving and retrieving URLs
type URLStorage interface {
	SaveURL(ctx context.Context, u URL) error
	GetURL(ctx context.Context, code string) (URL, error)
	GetByURL(ctx context.Context, url string) (URL, error)
	GetURLsByUserID(ctx context.Context, userID int) ([]URL, error)
	AllURLs(ctx context.Context) ([]URL, error)
	SaveBatchURL(ctx context.Context, urls []URL) error
	DeleteUserURLs(ctx context.Context, userID int, codes []string) error
	GetURLsCount(ctx context.Context) (int, error)
}

// UserStorage defines methods for user management
type UserStorage interface {
	CreateUser(ctx context.Context) (User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUsersCount(ctx context.Context) (int, error)
}

// Storage defines methods
type Storage interface {
	URLStorage
	UserStorage
	Close() error
	Ping(ctx context.Context) error
}
