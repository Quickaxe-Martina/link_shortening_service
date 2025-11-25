package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/handler"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/repository"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func setupRouter(cfg *config.Config, store storage.Storage, deleteWorker *repository.DeleteURLsWorkers) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewHandler(cfg, store, deleteWorker)

	r.Use(logger.RequestLogger)
	r.Use(handler.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/{URLCode}", h.RedirectURL)
		r.With(h.GetOrCreateUserMiddleware).Post("/", h.GenerateURL)
	})
	r.Route("/api/shorten", func(r chi.Router) {
		r.With(h.GetOrCreateUserMiddleware).Post("/", h.JSONGenerateURL)
		r.Post("/batch", h.BatchGenerateURL)
	})
	r.Route("/api/user", func(r chi.Router) {
		r.With(h.GetOrCreateUserMiddleware).Get("/urls", h.GetUserURLs)
		r.With(h.GetOrCreateUserMiddleware).Delete("/urls", h.DeleteUserURLs)
	})
	r.Route("/ping", func(r chi.Router) {
		r.Get("/", h.Ping)
	})
	return r
}

func main() {
	cfg := config.NewConfig()
	store, err := storage.NewStorage(cfg)
	deleteWorker := repository.NewDeleteURLsWorkers(store, 3, time.Duration(cfg.DeleteTimeDuration), cfg.DeleteBachSize)
	if err != nil {
		logger.Log.Error("Storage error", zap.Error(err))
	}

	if err := logger.Initialize("info"); err != nil {
		log.Panic(err)
	}
	r := setupRouter(cfg, store, deleteWorker)

	// Обработчик завершения (Ctrl+C, SIGTERM и т.п.)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		store.Close()
		deleteWorker.Stop()
		os.Exit(0)
	}()

	logger.Log.Fatal("", zap.Error(http.ListenAndServe(cfg.RunAddr, r)))
}
