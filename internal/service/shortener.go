package service

import (
	"context"
	"errors"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/repository"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
)

type ShortenerService struct {
	store storage.Storage
	cfg   *config.Config
	audit *repository.AuditPublisher
}

func NewShortenerService(
	store storage.Storage,
	cfg *config.Config,
	audit *repository.AuditPublisher,
) *ShortenerService {
	return &ShortenerService{store: store, cfg: cfg, audit: audit}
}

func (s *ShortenerService) Shorten(
	ctx context.Context,
	userID int,
	rawURL string,
) (string, error) {

	if rawURL == "" {
		return "", errors.New("empty url")
	}

	s.audit.Publish(repository.AuditEvent{
		TS:     time.Now().Unix(),
		Action: "shorten",
		UserID: userID,
		URL:    rawURL,
	})

	code, err := GenerateRandomString(6)
	if err != nil {
		return "", err
	}

	err = s.store.SaveURL(ctx, storage.URL{
		Code:   code,
		URL:    rawURL,
		UserID: userID,
	})

	if errors.Is(err, storage.ErrURLAlreadyExists) {
		url, err := s.store.GetByURL(ctx, rawURL)
		if err != nil {
			return "", err
		}
		return s.cfg.ServerAddr + url.Code, storage.ErrURLAlreadyExists
	}

	if err != nil {
		return "", err
	}

	return s.cfg.ServerAddr + code, nil
}
