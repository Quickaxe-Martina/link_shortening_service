/*
Package repository for repository
*/
package repository

import (
	"encoding/json"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// RemoteAuditObserver sends audit events to a remote HTTP endpoint.
type RemoteAuditObserver struct {
	client *resty.Client
	url    string
}

// NewRemoteAuditObserver creates a new RemoteAuditObserver for the given URL.
func NewRemoteAuditObserver(url string) *RemoteAuditObserver {
	return &RemoteAuditObserver{
		client: resty.New(),
		url:    url,
	}
}

// Notify sends the audit event to the configured remote endpoint.
func (a *RemoteAuditObserver) Notify(event AuditEvent) {
	body, _ := json.Marshal(event)

	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(a.url)

	if err != nil {
		logger.Log.Error("audit POST failed", zap.Error(err))
		return
	}

	if resp.IsError() {
		logger.Log.Warn("audit server returned error",
			zap.Int("status", resp.StatusCode()),
			zap.String("body", resp.String()),
		)
	}
}
