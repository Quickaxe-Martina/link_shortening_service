package storage

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"go.uber.org/zap"
)

type savedURLItem struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int    `json:"user_id"`
}

// LoadData load data from file
func LoadData(filePath string, store Storage) {
	file, err := os.Open(filePath)
	var savedData []savedURLItem
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Log.Fatal("open file error", zap.Error(err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&savedData)
	if err != nil {
		logger.Log.Error("JSON decoding error", zap.Error(err))
	}
	for _, item := range savedData {
		store.SaveURL(context.TODO(), URL{Code: item.ShortURL, URL: item.OriginalURL, UserID: item.UserID})
	}
}

// SaveData save data in file
func SaveData(filePath string, store Storage) {
	var saveData []savedURLItem
	file, err := os.Create(filePath)
	if err != nil {
		logger.Log.Error("file creation error", zap.Error(err))
		return
	}
	defer file.Close()

	urls, err := store.AllURLs(context.TODO())
	if err != nil {
		logger.Log.Error("load urls error", zap.Error(err))
		return
	}
	i := 1
	for _, url := range urls {
		item := savedURLItem{
			UUID:        strconv.Itoa(i),
			ShortURL:    url.Code,
			OriginalURL: url.URL,
			UserID:      url.UserID,
		}
		saveData = append(saveData, item)
		i++
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(saveData); err != nil {
		logger.Log.Error("JSON encoding error", zap.Error(err))
	}
}
