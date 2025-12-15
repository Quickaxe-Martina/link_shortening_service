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

// AuditPublisher publishes audit events to registered observers.
// generate:reset
type AuditPublisher struct {
	observers []AuditObserver
	ch        chan AuditEvent
	doneCh    chan struct{}
	wg        sync.WaitGroup
}

// NewAuditPublisher creates a new AuditPublisher with a buffered channel.
func NewAuditPublisher(buffer int) *AuditPublisher {
	p := &AuditPublisher{
		ch:     make(chan AuditEvent, buffer),
		doneCh: make(chan struct{}),
	}
	p.wg.Add(1)
	go p.worker()
	return p
}

// Register adds a new observer to receive audit events.
func (p *AuditPublisher) Register(obs AuditObserver) {
	p.observers = append(p.observers, obs)
}

// Publish sends an audit event to all registered observers asynchronously.
func (p *AuditPublisher) Publish(event AuditEvent) {
	p.ch <- event
}

// Stop signals the worker to stop and waits for all events to be processed.
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
