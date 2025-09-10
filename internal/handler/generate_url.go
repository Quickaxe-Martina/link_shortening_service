package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/model"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"
	"go.uber.org/zap"
)

// GenerateURL handles HTTP requests to create a shortened URL.
func (h *Handler) GenerateURL(w http.ResponseWriter, r *http.Request) {
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
	h.storageData.URLData[URLCode] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.cfg.ServerAddr + URLCode))
}

// JSONGenerateURL handles HTTP JSON requests to create a shortened URL.
func (h *Handler) JSONGenerateURL(w http.ResponseWriter, r *http.Request) {
	var req model.JSONGenerateURLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	URLCode, err := service.GenerateRandomString(6)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("URL code", zap.String("URLCode", URLCode))
	h.storageData.URLData[URLCode] = req.URL

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := model.JSONGenerateURLResponse{
		Result: h.cfg.ServerAddr + URLCode,
	}
	// сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}
