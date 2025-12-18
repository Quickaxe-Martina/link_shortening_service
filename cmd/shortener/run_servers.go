package main

import (
	"context"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/repository"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func serve(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func shutdown(
	cfg *config.Config,
	httpServer *http.Server,
	pprofServer *http.Server,
	store storage.Storage,
	deleteWorker *repository.DeleteURLsWorkers,
	audit *repository.AuditPublisher,
) error {
	logger.Log.Info("shutdown signal received")
	timeout := time.Duration(cfg.ShutdownTimeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_ = httpServer.Shutdown(ctx)
	_ = pprofServer.Shutdown(ctx)

	deleteWorker.Stop()
	audit.Stop()
	store.Close()

	logger.Log.Info("graceful shutdown completed")
	return nil
}

func runServers(
	ctx context.Context,
	cfg *config.Config,
	httpServer *http.Server,
	pprofServer *http.Server,
	store storage.Storage,
	deleteWorker *repository.DeleteURLsWorkers,
	audit *repository.AuditPublisher,
) {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Log.Info("HTTP server started", zap.String("addr", httpServer.Addr))
		return serve(httpServer)
	})

	g.Go(func() error {
		logger.Log.Info("pprof server started", zap.String("addr", pprofServer.Addr))
		return serve(pprofServer)
	})

	g.Go(func() error {
		<-gCtx.Done()
		return shutdown(cfg, httpServer, pprofServer, store, deleteWorker, audit)
	})

	if err := g.Wait(); err != nil {
		logger.Log.Error("server exited with error", zap.Error(err))
	}
}
