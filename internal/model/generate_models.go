/*
Package model for models
*/
package model

import (
	"fmt"
	"net/url"
)

// JSONGenerateURLRequest model for request
// generate:reset
type JSONGenerateURLRequest struct {
	URL string `json:"url"`
}

// Validate validation method
func (r *JSONGenerateURLRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("url is required")
	}

	u, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("url must include scheme and host (e.g. https://example.com)")
	}

	return nil
}

// JSONGenerateURLResponse model for response
// generate:reset
type JSONGenerateURLResponse struct {
	Result string `json:"result"`
}

// BatchGenerateURLRequest model for request
// generate:reset
type BatchGenerateURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

// Validate validation method
func (r *BatchGenerateURLRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("url is required")
	}

	u, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("url must include scheme and host (e.g. https://example.com)")
	}

	return nil
}

// BatchGenerateURLResponse model for response
// generate:reset
type BatchGenerateURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURLsResponse model for response
// generate:reset
type UserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
