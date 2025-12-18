package main

import (
	"context"
	"fmt"
	"log"
	"net"
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
	"github.com/Quickaxe-Martina/link_shortening_service/internal/tools"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
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

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}

func setupSignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
}

func setupAudit(cfg *config.Config) *repository.AuditPublisher {
	audit := repository.NewAuditPublisher(100)

	if cfg.AuditFile != "" {
		audit.Register(repository.NewFileAuditObserver(cfg.AuditFile))
	}
	if cfg.AuditURL != "" {
		audit.Register(repository.NewRemoteAuditObserver(cfg.AuditURL))
	}

	return audit
}

func main() {
	printBuildInfo()
	mainCtx, stop := setupSignalContext()
	defer stop()

	cfg := config.NewConfig()

	if err := logger.Initialize("info"); err != nil {
		log.Panic(err)
	}

	store, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Log.Fatal("storage init error", zap.Error(err))
	}

	deleteWorker := repository.NewDeleteURLsWorkers(
		store,
		3,
		time.Duration(cfg.DeleteTimeDuration),
		cfg.DeleteBachSize,
	)

	audit := setupAudit(cfg)

	r := setupRouter(cfg, store, deleteWorker, audit)

	httpServer := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: r,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	pprofServer := &http.Server{
		Addr: "localhost:6060",
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	tools.RunServers(mainCtx, cfg, httpServer, pprofServer, store, deleteWorker, audit)
}
