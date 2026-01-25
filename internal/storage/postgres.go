package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // driver

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	_ "github.com/Quickaxe-Martina/link_shortening_service/internal/logger" // logger
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

// PostgresStorage is DB implementation of the Storage interface
// generate:reset
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
	userID := sql.NullInt64{}
	if u.UserID != 0 {
		userID.Int64 = int64(u.UserID)
		userID.Valid = true
	}
	_, err := store.DB.ExecContext(ctx, "INSERT INTO urls (code, url, user_id) VALUES ($1, $2, $3)", u.Code, u.URL, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case "idx_urls_code":
					return ErrCodeAlreadyExists
				case "idx_urls_url":
					return ErrURLAlreadyExists
				default:
					return err
				}
			}
		}
		return err
	}
	return nil
}

// GetURL get URL by code from DB
func (store *PostgresStorage) GetURL(ctx context.Context, code string) (URL, error) {
	row := store.DB.QueryRowContext(ctx, "SELECT url, is_deleted FROM urls WHERE code = $1", code)
	var url string
	var isDeleted bool
	if err := row.Scan(&url, &isDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return URL{}, fmt.Errorf("url with code %s not found", code)
		}
		return URL{}, err
	}
	if isDeleted {
		return URL{}, ErrURLDeleted
	}

	return URL{Code: code, URL: url, isDeleted: isDeleted}, nil
}

// GetByURL get URL by url from DB
func (store *PostgresStorage) GetByURL(ctx context.Context, url string) (URL, error) {
	row := store.DB.QueryRowContext(ctx, "SELECT code, is_deleted FROM urls WHERE url = $1", url)
	var code string
	var isDeleted bool
	if err := row.Scan(&code, &isDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return URL{}, fmt.Errorf("url with url %s not found", url)
		}
		return URL{}, err
	}
	if isDeleted {
		return URL{}, ErrURLDeleted
	}

	return URL{Code: code, URL: url, isDeleted: isDeleted}, nil
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
	return []URL{}, ErrNotImplemented
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

// CreateUser creates a new user and returns it
func (store *PostgresStorage) CreateUser(ctx context.Context) (User, error) {
	var id int
	err := store.DB.QueryRowContext(ctx, "INSERT INTO users DEFAULT VALUES RETURNING id").Scan(&id)
	if err != nil {
		return User{}, err
	}
	return User{ID: int(id)}, nil
}

// GetURLsByUserID returns all URLs associated with a specific user ID
func (store *PostgresStorage) GetURLsByUserID(ctx context.Context, userID int) ([]URL, error) {
	rows, err := store.DB.QueryContext(ctx, "SELECT code, url FROM urls WHERE user_id = $1 AND is_deleted=FALSE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.Code, &url.URL); err != nil {
			return nil, err
		}
		url.UserID = userID
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

// GetAllUsers returns all users
func (store *PostgresStorage) GetAllUsers(ctx context.Context) ([]User, error) {
	return nil, ErrNotImplemented
}

// DeleteUserURLs marks user's URLs as deleted.
// Security: ensures only URLs belonging to the given userID are affected.
func (store *PostgresStorage) DeleteUserURLs(ctx context.Context, userID int, codes []string) error {
	query := `
        UPDATE urls
        SET is_deleted = TRUE
        WHERE user_id = $1 AND code = ANY($2::text[]);
    `
	_, err := store.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}

// GetURLsCount returns URLs count
func (store *PostgresStorage) GetURLsCount(ctx context.Context) (int, error) {
	row := store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM urls WHERE is_deleted = false")

	var count int
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// GetUsersCount returns users count
func (store *PostgresStorage) GetUsersCount(ctx context.Context) (int, error) {
	row := store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")

	var count int
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}
