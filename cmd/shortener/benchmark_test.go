package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/model"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/repository"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"github.com/go-resty/resty/v2"
)

func setupBenchServer(b *testing.B) (*resty.Client, *httptest.Server) {
	cfg := &config.Config{
		RunAddr:    ":8080",
		ServerAddr: "http://localhost:8080/",
	}
	storageData, err := storage.NewStorage(cfg)
	assert.NoError(b, err)

	storageData.SaveURL(context.TODO(), storage.URL{Code: "bench", URL: "https://example.com"})

	deleteWorker := repository.NewDeleteURLsWorkers(storageData, 3, 2*time.Second, 50)
	audit := repository.NewAuditPublisher(100)
	router := setupRouter(cfg, storageData, deleteWorker, audit)
	srv := httptest.NewServer(router)

	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	return client, srv
}

func BenchmarkGenerateURL(b *testing.B) {
	client, srv := setupBenchServer(b)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		randStr, err := service.GenerateRandomString(6)
		if err != nil {
			return
		}
		resp, err := client.R().
			SetHeader("Content-Type", "text/plain").
			SetBody("https://example-" + randStr + ".com").
			Post(srv.URL + "/")

		if err != nil || resp.StatusCode() != 201 {
			b.Fatalf("bad response: %v, code=%d", err, resp.StatusCode())
		}
	}
}

func BenchmarkRedirectURL(b *testing.B) {
	client, srv := setupBenchServer(b)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.R().
			Get(srv.URL + "/bench")

		if err != nil {
			// Resty может бросать ошибку из-за NoRedirectPolicy, игнорируем
			if _, ok := err.(*url.Error); !ok {
				b.Fatalf("err: %v", err)
			}
		}

		if resp.StatusCode() != 307 {
			b.Fatalf("unexpected status %d", resp.StatusCode())
		}
	}
}

func BenchmarkJSONGenerateURL(b *testing.B) {
	client, srv := setupBenchServer(b)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		randStr, err := service.GenerateRandomString(6)
		if err != nil {
			return
		}
		req := model.JSONGenerateURLRequest{URL: "https://example-" + randStr + ".com"}
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(req).
			Post(srv.URL + "/api/shorten")

		if err != nil || resp.StatusCode() != 201 {
			b.Fatalf("bad response: %v, code=%d", err, resp.StatusCode())
		}

		var result model.JSONGenerateURLResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			b.Fatalf("json decode: %v", err)
		}
	}
}
