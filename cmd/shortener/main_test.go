package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var createdCodes []string

func TestGenerateUrl(t *testing.T) {
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
			wantPrefix:  hostname,
		},
		{
			name:        "неверный Content-Type",
			contentType: "application/json",
			body:        "https://example.com",
			wantStatus:  http.StatusBadRequest,
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
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			generateUrl(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantPrefix != "" && res.StatusCode == http.StatusCreated {
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.True(t, strings.HasPrefix(string(data), tt.wantPrefix))

				// Сохраняем код для дальнейшего теста
				code := strings.TrimPrefix(string(data), hostname)
				createdCodes = append(createdCodes, code)
			}
		})
	}
}

func TestRedirectUrl(t *testing.T) {
	for _, code := range createdCodes {
		urlData[code] = "https://example.com"
	}

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
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "неизвестный код",
			path:       "/unknown",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			redirectUrl(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantLoc != "" {
				assert.Equal(t, tt.wantLoc, res.Header.Get("Location"))
			}
		})
	}
}
