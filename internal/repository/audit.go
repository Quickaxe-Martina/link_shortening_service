package repository

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
)

// AuditEvent todo
type AuditEvent struct {
	TS     int64  `json:"ts"`
	Action string `json:"action"`
	UserID int    `json:"-"`
	URL    string `json:"url"`
}

// MarshalJSON todo
func (e AuditEvent) MarshalJSON() ([]byte, error) {
	type Alias AuditEvent
	return json.Marshal(&struct {
		UserID string `json:"user_id"`
		*Alias
	}{
		UserID: fmt.Sprintf("%d", e.UserID),
		Alias:  (*Alias)(&e),
	})
}

// AuditObserver todo
type AuditObserver interface {
	Notify(event AuditEvent)
}

// AuditPublisher todo
type AuditPublisher struct {
	observers []AuditObserver
	ch        chan AuditEvent
	doneCh    chan struct{}
	wg        sync.WaitGroup
}

// NewAuditPublisher todo
func NewAuditPublisher(buffer int) *AuditPublisher {
	p := &AuditPublisher{
		ch:     make(chan AuditEvent, buffer),
		doneCh: make(chan struct{}),
	}
	p.wg.Add(1)
	go p.worker()
	return p
}

// Register todo
func (p *AuditPublisher) Register(obs AuditObserver) {
	p.observers = append(p.observers, obs)
}

// Publish todo
func (p *AuditPublisher) Publish(event AuditEvent) {
	p.ch <- event
}

// Stop end workers work
func (p *AuditPublisher) Stop() {
	close(p.doneCh)
	p.wg.Wait()
	logger.Log.Info("All workers stopped")
	close(p.ch)
}

func (p *AuditPublisher) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.doneCh:
			return
		case event := <-p.ch:
			for _, o := range p.observers {
				go func(obs AuditObserver, ev AuditEvent) {
					obs.Notify(ev)
				}(o, event)
			}
		}
	}
}
