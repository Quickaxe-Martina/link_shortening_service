package storage

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"go.uber.org/zap"
)

// generate:reset
type savedURLItem struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int    `json:"user_id"`
}

// generate:reset
type savedUserItem struct {
	UUID string `json:"uuid"`
	ID   int    `json:"id"`
}

const userFilePrefix = "user_"

func loadFromFile(filePath string, data any) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Log.Fatal("open file error", zap.Error(err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(data)
	if err != nil {
		logger.Log.Error("JSON decoding error", zap.Error(err))
	}
}

// LoadData load data from file
func LoadData(filePath string, store Storage) {
	var savedData []savedURLItem
	loadFromFile(filePath, &savedData)
	for _, item := range savedData {
		store.SaveURL(context.TODO(), URL{Code: item.ShortURL, URL: item.OriginalURL, UserID: item.UserID})
	}

	var savedUsers []savedUserItem
	loadFromFile(userFilePrefix+filePath, &savedUsers)
	for _, item := range savedUsers {
		store.(*MemoryStorage).Users[item.ID] = User{ID: item.ID}
	}
}

func saveDataToFile(filePath string, data any) {
	file, err := os.Create(filePath)
	if err != nil {
		logger.Log.Error("file creation error", zap.Error(err))
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		logger.Log.Error("JSON encoding error", zap.Error(err))
	}
}

// SaveData save data in file
func SaveData(filePath string, store Storage) {
	var saveURLData []savedURLItem

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
		saveURLData = append(saveURLData, item)
		i++
	}
	saveDataToFile(filePath, saveURLData)

	// Save users
	var saveUserData []savedUserItem
	users, err := store.GetAllUsers(context.TODO())
	if err != nil {
		logger.Log.Error("load users error", zap.Error(err))
		return
	}
	for _, user := range users {
		item := savedUserItem{
			UUID: strconv.Itoa(user.ID),
			ID:   user.ID,
		}
		saveUserData = append(saveUserData, item)
	}
	saveDataToFile(userFilePrefix+filePath, saveUserData)
}
