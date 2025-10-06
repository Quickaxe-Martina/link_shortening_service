package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/logger"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"go.uber.org/zap"
)

// ErrWorkerStopped woker stopped
var ErrWorkerStopped = errors.New("woker stopped")

type deleteURLsTask struct {
	UserID int
	Codes  []string
}

// DeleteURLsWorkers struct
type DeleteURLsWorkers struct {
	store      storage.Storage
	inputCh    chan deleteURLsTask
	workerCh   chan deleteURLsTask
	doneCh     chan struct{}
	flushDelay time.Duration
	batchSize  int
	wg         sync.WaitGroup
}

// NewDeleteURLsWorkers create DeleteURLsWorkers
func NewDeleteURLsWorkers(store storage.Storage, numWorkers int, flushDelay time.Duration, batchSize int) *DeleteURLsWorkers {
	wm := &DeleteURLsWorkers{
		store:      store,
		inputCh:    make(chan deleteURLsTask, 100),
		workerCh:   make(chan deleteURLsTask, numWorkers*3),
		doneCh:     make(chan struct{}),
		flushDelay: flushDelay,
		batchSize:  batchSize,
	}

	wm.wg.Add(1)
	go wm.aggregator()
	for i := 0; i < numWorkers; i++ {
		wm.wg.Add(1)
		go wm.worker(i)
	}

	return wm
}

func (wm *DeleteURLsWorkers) aggregator() {
	logger.Log.Info("aggregator started")
	defer wm.wg.Done()
	count := 0
	ticker := time.NewTicker(wm.flushDelay)
	defer ticker.Stop()

	batch := make(map[int][]string)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		for userID, codes := range batch {
			wm.workerCh <- deleteURLsTask{UserID: userID, Codes: codes}
		}
		batch = make(map[int][]string)
	}

	for {
		select {
		case <-wm.doneCh:
			flush()
			close(wm.workerCh)
			logger.Log.Info("aggregator stopped")
			return
		case <-ticker.C:
			flush()
		case task := <-wm.inputCh:
			batch[task.UserID] = append(batch[task.UserID], task.Codes...)
			count += len(task.Codes)
			if count >= wm.batchSize {
				flush()
				count = 0
			}
		}
	}
}

func (wm *DeleteURLsWorkers) worker(id int) {
	defer wm.wg.Done()
	logger.Log.Info(fmt.Sprintf("worker-%d started", id))
	for {
		select {
		case <-wm.doneCh:
			logger.Log.Info(fmt.Sprintf("worker-%d stopping", id))
			return
		case task := <-wm.workerCh:
			wm.handleDeleteTask(task)
		}
	}
}

func (wm *DeleteURLsWorkers) handleDeleteTask(task deleteURLsTask) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := wm.store.DeleteUserURLs(ctx, task.UserID, task.Codes)
	if err != nil {
		logger.Log.Error("delete urls error", zap.Error(err))
	}
}

// AddTask add task to delete urls
func (wm *DeleteURLsWorkers) AddTask(userID int, codes []string) error {
	select {
	case <-wm.doneCh:
		return ErrWorkerStopped
	default:
	}

	wm.inputCh <- deleteURLsTask{UserID: userID, Codes: codes}
	return nil
}

// Stop end workers work
func (wm *DeleteURLsWorkers) Stop() {
	close(wm.doneCh)
	wm.wg.Wait()
	logger.Log.Info("All workers stopped")
	close(wm.inputCh)
}
