package worker

import (
	"context"
	"sync"
	"time"

	"github.com/octopunkio/taskflow-worker/internal/handler"
	"github.com/octopunkio/taskflow-worker/internal/job"
	"github.com/octopunkio/taskflow-worker/internal/queue"
	"go.uber.org/zap"
)

// Worker processes jobs from the queue
type Worker struct {
	queue       *queue.Queue
	handlers    map[string]handler.Handler
	concurrency int
	logger      *zap.Logger
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

// Config holds worker configuration
type Config struct {
	RedisAddr   string
	Concurrency int
	Logger      *zap.Logger
}

// New creates a new Worker instance
func New(cfg Config) *Worker {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 10
	}

	if cfg.Logger == nil {
		cfg.Logger, _ = zap.NewProduction()
	}

	return &Worker{
		queue: queue.New(queue.Config{
			RedisAddr: cfg.RedisAddr,
		}),
		handlers:    make(map[string]handler.Handler),
		concurrency: cfg.Concurrency,
		logger:      cfg.Logger,
	}
}

// RegisterHandler registers a handler for a job type
func (w *Worker) RegisterHandler(jobType string, h handler.Handler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[jobType] = h
	w.logger.Info("Registered handler", zap.String("type", jobType))
}

// Start begins processing jobs
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("Starting worker", zap.Int("concurrency", w.concurrency))

	// Start worker goroutines for each job type
	w.mu.RLock()
	jobTypes := make([]string, 0, len(w.handlers))
	for t := range w.handlers {
		jobTypes = append(jobTypes, t)
	}
	w.mu.RUnlock()

	for _, jobType := range jobTypes {
		for i := 0; i < w.concurrency; i++ {
			w.wg.Add(1)
			go w.processLoop(ctx, jobType)
		}
	}
}

// Wait blocks until all workers have finished
func (w *Worker) Wait() {
	w.wg.Wait()
	w.queue.Close()
}

func (w *Worker) processLoop(ctx context.Context, jobType string) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := w.processOne(ctx, jobType); err != nil {
				w.logger.Error("Process error", zap.Error(err))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (w *Worker) processOne(ctx context.Context, jobType string) error {
	j, err := w.queue.Dequeue(ctx, jobType)
	if err != nil {
		return err
	}
	if j == nil {
		return nil
	}

	w.mu.RLock()
	h, ok := w.handlers[jobType]
	w.mu.RUnlock()

	if !ok {
		w.logger.Warn("No handler for job type", zap.String("type", jobType))
		return nil
	}

	return w.executeJob(ctx, j, h)
}

func (w *Worker) executeJob(ctx context.Context, j *job.Job, h handler.Handler) error {
	w.logger.Info("Processing job",
		zap.String("id", j.ID),
		zap.String("type", j.Type),
		zap.Int("attempt", j.Attempts+1),
	)

	j.MarkStarted()
	w.queue.Update(ctx, j)

	if err := h.Handle(ctx, j); err != nil {
		j.MarkFailed(err)
		w.queue.Update(ctx, j)

		if j.ShouldRetry() {
			w.logger.Info("Retrying job",
				zap.String("id", j.ID),
				zap.Int("attempt", j.Attempts),
			)
			return w.queue.Requeue(ctx, j)
		}

		w.logger.Error("Job failed permanently",
			zap.String("id", j.ID),
			zap.Error(err),
		)
		return err
	}

	j.MarkCompleted()
	w.queue.Update(ctx, j)

	w.logger.Info("Job completed", zap.String("id", j.ID))
	return nil
}

// Enqueue adds a job to the queue
func (w *Worker) Enqueue(ctx context.Context, j *job.Job) error {
	return w.queue.Enqueue(ctx, j)
}
