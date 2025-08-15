package main

import (
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"log"

	"net/http"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/{URLCode}", handler.RedirectURL)
		r.Post("/", handler.GenerateURL)
	})
	return r
}

func main() {
	log.Println("Server started")
	config.ParseFlags()
	r := setupRouter()

	log.Fatal(http.ListenAndServe(config.FlagRunAddr, r))
}
