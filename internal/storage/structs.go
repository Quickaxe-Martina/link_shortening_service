package storage

import (
	"database/sql"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib" // driver
)

// Storage variables
type Storage struct {
	URLData map[string]string
	DB      *sql.DB
}

// NewStorage create Storage
func NewStorage(cfg *config.Config) *Storage {
	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		panic(err)
	}
	var storage = Storage{
		URLData: make(map[string]string),
		DB:      db,
	}
	return &storage
}

func (s *Storage) Close() {
	s.DB.Close()
}
