package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	r.Use(handler.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/{URLCode}", h.RedirectURL)
		r.Post("/", h.GenerateURL)
	})
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.JSONGenerateURL)
	})
	return r
}

func loadData(cfg *config.Config) {
	file, err := os.Open(cfg.DataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Log.Fatal("open file error", zap.Error(err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.URLData)
	if err != nil {
		logger.Log.Error("JSON decoding error", zap.Error(err))
	}
}

func saveData(cfg *config.Config) {
	file, err := os.Create(cfg.DataFilePath)
	if err != nil {
		logger.Log.Error("file creation error", zap.Error(err))
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg.URLData); err != nil {
		logger.Log.Error("JSON encoding error", zap.Error(err))
	}
}

func main() {
	cfg := config.NewConfig()
	loadData(cfg)

	if err := logger.Initialize("info"); err != nil {
		log.Panic(err)
	}
	r := setupRouter(cfg)

	// Обработчик завершения (Ctrl+C, SIGTERM и т.п.)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		saveData(cfg)
		os.Exit(0)
	}()

	logger.Log.Fatal("", zap.Error(http.ListenAndServe(cfg.RunAddr, r)))
}
