package main

import (
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"log"

	"net/http"
)

func setupRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewHandler(cfg)
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/{URLCode}", h.RedirectURL)
		r.Post("/", h.GenerateURL)
	})
	return r
}

func main() {
	log.Println("Server started")
	cfg := config.NewConfig()
	r := setupRouter(cfg)

	log.Fatal(http.ListenAndServe(cfg.RunAddr, r))
}
