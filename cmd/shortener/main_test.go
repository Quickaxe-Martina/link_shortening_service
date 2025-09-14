package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/model"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateURL(t *testing.T) {
	cfg := &config.Config{
		RunAddr:    ":8080",
		ServerAddr: "http://localhost:8080/",
	}
	storageData, err := storage.NewStorage(cfg)
	assert.NoError(t, err)
	router := setupRouter(cfg, storageData)
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
			wantPrefix:  cfg.ServerAddr,
		},
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
				assert.Contains(t, resp.String(), cfg.ServerAddr)
			}
		})
	}
}

func TestRedirectURL(t *testing.T) {
	cfg := &config.Config{
		RunAddr:    ":8080",
		ServerAddr: "http://localhost:8080/",
	}
	store, err := storage.NewStorage(cfg)
	assert.NoError(t, err)
	store.SaveURL(context.TODO(), storage.URL{Code: "qwerty", URL: "https://example.com"})
	router := setupRouter(cfg, store)
	srv := httptest.NewServer(router)
	defer srv.Close()

	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantLoc    string
	}{
		{
			name:       "успешный редирект",
			path:       "/qwerty",
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

func TestJSONGenerateURL(t *testing.T) {
	cfg := &config.Config{
		RunAddr:    ":8080",
		ServerAddr: "http://localhost:8080/",
	}
	storageData, err := storage.NewStorage(cfg)
	assert.NoError(t, err)
	router := setupRouter(cfg, storageData)
	srv := httptest.NewServer(router)
	defer srv.Close()

	client := resty.New()

	tests := []struct {
		name       string
		req        interface{}
		wantStatus int
		wantPrefix string
	}{
		{
			name:       "успешное создание",
			req:        model.JSONGenerateURLRequest{URL: "https://example.com"},
			wantStatus: http.StatusCreated,
			wantPrefix: cfg.ServerAddr,
		},
		{
			name:       "пустой url",
			req:        model.JSONGenerateURLRequest{URL: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "некорректный url",
			req:        model.JSONGenerateURLRequest{URL: "not-a-url"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустое тело",
			req:        map[string]interface{}{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.req).
				Post(srv.URL + "/api/shorten")
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode())
			if tt.wantStatus == http.StatusCreated {
				assert.Contains(t, resp.String(), cfg.ServerAddr)
				var result model.JSONGenerateURLResponse
				err := json.Unmarshal(resp.Body(), &result)
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Result)
				assert.Contains(t, result.Result, tt.wantPrefix)
			}
		})
	}
}
