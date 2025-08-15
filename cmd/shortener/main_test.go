package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var createdCodes []string

func TestGenerateURL(t *testing.T) {
	router := setupRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	client := resty.New()

	tests := []struct {
		name        string
		contentType string
		body        string
		wantStatus  int
		wantPrefix  string
	}{
		{
			name:        "успешное создание",
			contentType: "text/plain",
			body:        "https://example.com",
			wantStatus:  http.StatusCreated,
			wantPrefix:  config.FlagServerAddr,
		},
		// {
		// 	name:        "неверный Content-Type",
		// 	contentType: "application/json",
		// 	body:        "https://example.com",
		// 	wantStatus:  http.StatusBadRequest,
		// },
		{
			name:        "пустое тело",
			contentType: "text/plain",
			body:        "",
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().
				SetBody(tt.body).
				Post(srv.URL + "/")
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode())
			if tt.wantStatus == http.StatusCreated {
				assert.Contains(t, resp.String(), config.FlagServerAddr)
				// Сохраняем код для TestRedirectURL
				code := strings.TrimPrefix(resp.String(), config.FlagServerAddr)
				createdCodes = append(createdCodes, code)
			}
		})
	}
}

func TestRedirectURL(t *testing.T) {
	router := setupRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	for _, code := range createdCodes {
		config.URLData[code] = "https://example.com"
	}
	log.Printf("createdCodes[0]: %s", string(createdCodes[0]))

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantLoc    string
	}{
		{
			name:       "успешный редирект",
			path:       "/" + createdCodes[0],
			wantStatus: http.StatusTemporaryRedirect,
			wantLoc:    "https://example.com",
		},
		{
			name:       "пустой код",
			path:       "/",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "неизвестный код",
			path:       "/unknown",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("url: %s", srv.URL+tt.path)
			resp, err := client.R().
				Get(srv.URL + tt.path)

			var urlErr *url.Error
			if err != nil && !(errors.As(err, &urlErr) && urlErr.Err.Error() == "auto redirect is disabled") {
				log.Printf("err.Error(): %s", err.Error())
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatus, resp.StatusCode())
			if tt.wantLoc != "" {
				assert.Equal(t, tt.wantLoc, resp.Header().Get("Location"))
			}
		})
	}
}
