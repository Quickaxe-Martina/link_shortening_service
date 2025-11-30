package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/handler"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/repository"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func setupRouter(cfg *config.Config, store storage.Storage, deleteWorker *repository.DeleteURLsWorkers, audit *repository.AuditPublisher) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewHandler(cfg, store, deleteWorker, audit)

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
	// Запустим HTTP-сервер для обработки запросов на профилирование
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	cfg := config.NewConfig()
	store, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Log.Error("Storage error", zap.Error(err))
	}

	deleteWorker := repository.NewDeleteURLsWorkers(store, 3, time.Duration(cfg.DeleteTimeDuration), cfg.DeleteBachSize)

	audit := repository.NewAuditPublisher(100)

	if cfg.AuditFile != "" {
		audit.Register(repository.NewFileAuditObserver(cfg.AuditFile))
	}

	if cfg.AuditURL != "" {
		audit.Register(repository.NewRemoteAuditObserver(cfg.AuditURL))
	}

	if err := logger.Initialize("info"); err != nil {
		log.Panic(err)
	}
	r := setupRouter(cfg, store, deleteWorker, audit)

	// Обработчик завершения (Ctrl+C, SIGTERM и т.п.)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		store.Close()
		deleteWorker.Stop()
		audit.Stop()
		os.Exit(0)
	}()

	logger.Log.Fatal("", zap.Error(http.ListenAndServe(cfg.RunAddr, r)))
}
