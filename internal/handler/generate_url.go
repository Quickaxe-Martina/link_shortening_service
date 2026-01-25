package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/model"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"go.uber.org/zap"
)

// GenerateURL handles HTTP requests to create a shortened URL.
func (h *Handler) GenerateURL(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	shortURL, err := h.shortener.Shorten(ctx, user.ID, string(body))
	if errors.Is(err, storage.ErrURLAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(shortURL))
		return
	}
	if err != nil {
		logger.Log.Error("error shortening URL", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// JSONGenerateURL handles HTTP JSON requests to create a shortened URL.
func (h *Handler) JSONGenerateURL(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r.Context())

	var req model.JSONGenerateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	shortURL, err := h.shortener.Shorten(ctx, user.ID, req.URL)
	if errors.Is(err, storage.ErrURLAlreadyExists) {
		resp := model.JSONGenerateURLResponse{Result: shortURL}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Log.Error("error encoding response", zap.Error(err))
		}
		return
	}
	if err != nil {
		logger.Log.Error("error shortening URL", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := model.JSONGenerateURLResponse{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
	}
}

// BatchGenerateURL handles HTTP JSON requests to create a shortened URL.
func (h *Handler) BatchGenerateURL(w http.ResponseWriter, r *http.Request) {
	var requests []model.BatchGenerateURLRequest
	var urls []storage.URL
	var responses []model.BatchGenerateURLResponse

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&requests); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, req := range requests {
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		URLCode, err := service.GenerateRandomString(6)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		url := storage.URL{Code: URLCode, URL: req.URL}
		urls = append(urls, url)
		resp := model.BatchGenerateURLResponse{CorrelationID: req.CorrelationID, ShortURL: h.cfg.ServerAddr + URLCode}
		responses = append(responses, resp)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := h.store.SaveBatchURL(ctx, urls); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(responses); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}
