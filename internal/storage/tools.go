package storage

import (
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
}

// LoadData load data from file
func LoadData(filePath string, storageData *Storage) {
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
		storageData.URLData[item.ShortURL] = item.OriginalURL
	}
}

// SaveData save data in file
func SaveData(filePath string, storageData *Storage) {
	var saveData []savedURLItem
	file, err := os.Create(filePath)
	if err != nil {
		logger.Log.Error("file creation error", zap.Error(err))
		return
	}
	defer file.Close()

	i := 1
	for short, original := range storageData.URLData {
		item := savedURLItem{
			UUID:        strconv.Itoa(i),
			ShortURL:    short,
			OriginalURL: original,
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
