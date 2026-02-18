package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/octopunkio/taskflow-worker/internal/job"
)

// Queue manages job storage and retrieval
type Queue struct {
	client *redis.Client
	prefix string
}

// Config holds queue configuration
type Config struct {
	RedisAddr string
	Prefix    string
}

// New creates a new Queue instance
func New(cfg Config) *Queue {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "taskflow:jobs"
	}

	return &Queue{
		client: client,
		prefix: prefix,
	}
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(ctx context.Context, j *job.Job) error {
	data, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	// Store job data
	key := q.jobKey(j.ID)
	if err := q.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("store job: %w", err)
	}

	// Add to priority queue
	score := float64(j.Priority)*1e12 - float64(j.CreatedAt.UnixNano())
	if err := q.client.ZAdd(ctx, q.queueKey(j.Type), &redis.Z{
		Score:  score,
		Member: j.ID,
	}).Err(); err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}

	return nil
}

// Dequeue retrieves the next job from the queue
func (q *Queue) Dequeue(ctx context.Context, jobType string) (*job.Job, error) {
	// Get highest priority job
	result, err := q.client.ZPopMax(ctx, q.queueKey(jobType), 1).Result()
	if err == redis.Nil || len(result) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dequeue job: %w", err)
	}

	jobID := result[0].Member.(string)

	// Fetch job data
	data, err := q.client.Get(ctx, q.jobKey(jobID)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("fetch job: %w", err)
	}

	var j job.Job
	if err := json.Unmarshal(data, &j); err != nil {
		return nil, fmt.Errorf("unmarshal job: %w", err)
	}

	return &j, nil
}

// Update persists job state changes
func (q *Queue) Update(ctx context.Context, j *job.Job) error {
	data, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	return q.client.Set(ctx, q.jobKey(j.ID), data, 24*time.Hour).Err()
}

// Requeue adds a failed job back to the queue for retry
func (q *Queue) Requeue(ctx context.Context, j *job.Job) error {
	j.Status = job.StatusRetrying
	return q.Enqueue(ctx, j)
}

func (q *Queue) jobKey(id string) string {
	return fmt.Sprintf("%s:data:%s", q.prefix, id)
}

func (q *Queue) queueKey(jobType string) string {
	return fmt.Sprintf("%s:queue:%s", q.prefix, jobType)
}

// Close closes the Redis connection
func (q *Queue) Close() error {
	return q.client.Close()
}
