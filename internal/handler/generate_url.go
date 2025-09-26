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
	tokenExp := time.Hour * time.Duration(h.cfg.TokenExp)
	user, err := service.GetOrCreateUser(w, r, h.store, h.cfg.SecretKey, tokenExp)
	if err != nil {
		logger.Log.Error("error get or create user", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	URLCode, err := service.GenerateRandomString(6)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Log.Info("URL code", zap.String("URLCode", URLCode))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := h.store.SaveURL(ctx, storage.URL{Code: URLCode, URL: string(body), UserID: user.ID}); err != nil {
		if errors.Is(err, storage.ErrURLAlreadyExists) {
			url, err := h.store.GetByURL(ctx, string(body))
			if err != nil {
				logger.Log.Error("error get by url", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(h.cfg.ServerAddr + url.Code))
			return
		}
		logger.Log.Error("", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.cfg.ServerAddr + URLCode))
}

// JSONGenerateURL handles HTTP JSON requests to create a shortened URL.
func (h *Handler) JSONGenerateURL(w http.ResponseWriter, r *http.Request) {
	tokenExp := time.Hour * time.Duration(h.cfg.TokenExp)
	user, err := service.GetOrCreateUser(w, r, h.store, h.cfg.SecretKey, tokenExp)
	if err != nil {
		logger.Log.Error("error get or create user", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var req model.JSONGenerateURLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	URLCode, err := service.GenerateRandomString(6)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("URL code", zap.String("URLCode", URLCode))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := h.store.SaveURL(ctx, storage.URL{Code: URLCode, URL: req.URL, UserID: user.ID}); err != nil {
		if errors.Is(err, storage.ErrURLAlreadyExists) {
			logger.Log.Info("ErrURLAlreadyExists")
			url, err := h.store.GetByURL(ctx, req.URL)
			if err != nil {
				logger.Log.Error("error get by url", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)

			resp := model.JSONGenerateURLResponse{
				Result: h.cfg.ServerAddr + url.Code,
			}
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				logger.Log.Error("error encoding response", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
		logger.Log.Error("", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := model.JSONGenerateURLResponse{
		Result: h.cfg.ServerAddr + URLCode,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
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
