package main

import (
	"log"
	"net/http"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/handler"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func setupRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewHandler(cfg)

	r.Use(logger.RequestLogger)
	r.Route("/", func(r chi.Router) {
		r.Get("/{URLCode}", h.RedirectURL)
		r.Post("/", h.GenerateURL)
	})
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.JSONGenerateURL)
	})
	return r
}

func main() {
	cfg := config.NewConfig()

	if err := logger.Initialize("info"); err != nil {
		log.Panic(err)
	}
	r := setupRouter(cfg)

	logger.Log.Fatal("", zap.Error(http.ListenAndServe(cfg.RunAddr, r)))
}
