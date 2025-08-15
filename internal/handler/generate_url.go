package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"
)

func GenerateURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	URLCode := service.GenerateRandomString(6)
	log.Printf("URL code: %s", URLCode)
	config.URLData[URLCode] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.FlagServerAddr + URLCode))
}
