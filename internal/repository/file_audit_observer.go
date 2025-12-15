package repository

import (
	"encoding/json"
	"os"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"go.uber.org/zap"
)

// FileAuditObserver writes audit events to a file.
type FileAuditObserver struct {
	filePath string
}

// NewFileAuditObserver creates a new FileAuditObserver for the given file path.
func NewFileAuditObserver(path string) *FileAuditObserver {
	return &FileAuditObserver{filePath: path}
}

// Notify appends an audit event to the file in JSON format.
func (a *FileAuditObserver) Notify(event AuditEvent) {
	file, err := os.OpenFile(a.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Error("cannot open audit file", zap.Error(err))
		return
	}
	defer file.Close()

	b, _ := json.Marshal(event)
	file.Write(append(b, '\n'))
}
