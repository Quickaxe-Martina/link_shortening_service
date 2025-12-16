package repository

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
)

// AuditEvent represents a user action event for auditing purposes.
// generate:reset
type AuditEvent struct {
	TS     int64  `json:"ts"`
	Action string `json:"action"`
	UserID int    `json:"-"`
	URL    string `json:"url"`
}

// MarshalJSON customizes JSON serialization for AuditEvent.
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

// AuditObserver defines an interface for receiving audit events.
type AuditObserver interface {
	Notify(event AuditEvent)
}

type observerWorker struct {
	obs AuditObserver
	ch  chan AuditEvent
}

// AuditPublisher publishes audit events to registered observers.
// generate:reset
type AuditPublisher struct {
	observers []observerWorker
	ch        chan AuditEvent
	wg        sync.WaitGroup
}

// NewAuditPublisher creates a new AuditPublisher with a buffered channel.
func NewAuditPublisher(buffer int) *AuditPublisher {
	p := &AuditPublisher{
		ch: make(chan AuditEvent, buffer),
	}

	p.wg.Add(1)
	go p.worker()

	return p
}

// Register adds a new observer to receive audit events.
func (p *AuditPublisher) Register(obs AuditObserver) {
	w := observerWorker{
		obs: obs,
		ch:  make(chan AuditEvent, 10),
	}

	p.observers = append(p.observers, w)

	p.wg.Add(1)
	go func(w observerWorker) {
		defer p.wg.Done()
		for ev := range w.ch {
			w.obs.Notify(ev)
		}
	}(w)
}

// Publish sends an audit event to all registered observers asynchronously.
func (p *AuditPublisher) Publish(event AuditEvent) {
	p.ch <- event
}

// Stop signals the worker to stop and waits for all events to be processed.
func (p *AuditPublisher) Stop() {
	close(p.ch)
	p.wg.Wait()
	logger.Log.Info("All workers stopped")
}

func (p *AuditPublisher) worker() {
	defer p.wg.Done()

	for event := range p.ch {
		for _, o := range p.observers {
			o.ch <- event
		}
	}

	for _, o := range p.observers {
		close(o.ch)
	}
}
