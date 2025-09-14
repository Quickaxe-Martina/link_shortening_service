package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // driver

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file" //
)

// PostgresStorage is DB implementation of the Storage interface
type PostgresStorage struct {
	DB *sql.DB
}

// NewPostgresStorage creates new PostgresStorage
func NewPostgresStorage(cfg *config.Config) *PostgresStorage {
	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		panic(err)
	}
	if err := runMigrations(db, cfg.MigrationsPath); err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}
	store := &PostgresStorage{
		DB: db,
	}
	return store
}

// SaveURL save a URL by code in DB
func (store *PostgresStorage) SaveURL(ctx context.Context, u URL) error {
	_, err := store.DB.ExecContext(ctx, "INSERT INTO urls (code, url) VALUES ($1, $2)", u.Code, u.URL)
	if err != nil {
		return err
	}
	return nil
}

// GetURL get URL by code from DB
func (store *PostgresStorage) GetURL(ctx context.Context, code string) (URL, error) {
	row := store.DB.QueryRowContext(ctx, "SELECT url FROM urls WHERE code = $1", code)
	var url string
	if err := row.Scan(&url); err != nil {
		if err == sql.ErrNoRows {
			return URL{}, fmt.Errorf("url with code %s not found", code)
		}
		return URL{}, err
	}

	return URL{Code: code, URL: url}, nil
}

// Close releases resources
func (store *PostgresStorage) Close() error {
	store.DB.Close()
	return nil
}

// Ping DB
func (store *PostgresStorage) Ping(ctx context.Context) error {
	if err := store.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// AllURLs returns all URLs
func (store *PostgresStorage) AllURLs(ctx context.Context) ([]URL, error) {
	// TODO
	return []URL{}, nil
}

// SaveBatchURL save list of URL
func (store *PostgresStorage) SaveBatchURL(ctx context.Context, urls []URL) error {
	tx, err := store.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (code, url) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, url := range urls {
		_, err := stmt.ExecContext(ctx, url.Code, url.URL)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("could not create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"pgx",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
