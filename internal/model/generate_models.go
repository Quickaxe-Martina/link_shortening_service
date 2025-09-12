package model

import (
	"fmt"
	"net/url"
)

// JSONGenerateURLRequest model for request
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
type JSONGenerateURLResponse struct {
	Result string `json:"result"`
}
