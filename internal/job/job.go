package job

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Priority represents job priority levels
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 5
	PriorityHigh   Priority = 10
)

// Status represents the current state of a job
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusRetrying  Status = "retrying"
)

// Job represents a unit of work to be processed
type Job struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Priority    Priority        `json:"priority"`
	Status      Status          `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Error       string          `json:"error,omitempty"`
}

// NewJob creates a new job with the given type and payload
func NewJob(jobType string, payload interface{}) (*Job, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:          uuid.New().String(),
		Type:        jobType,
		Payload:     data,
		Priority:    PriorityNormal,
		Status:      StatusPending,
		MaxAttempts: 3,
		CreatedAt:   time.Now(),
	}, nil
}

// WithPriority sets the job priority
func (j *Job) WithPriority(p Priority) *Job {
	j.Priority = p
	return j
}

// WithMaxAttempts sets the maximum retry attempts
func (j *Job) WithMaxAttempts(n int) *Job {
	j.MaxAttempts = n
	return j
}

// ShouldRetry returns true if the job can be retried
func (j *Job) ShouldRetry() bool {
	return j.Attempts < j.MaxAttempts
}

// MarkStarted updates the job status to running
func (j *Job) MarkStarted() {
	now := time.Now()
	j.Status = StatusRunning
	j.StartedAt = &now
	j.Attempts++
}

// MarkCompleted updates the job status to completed
func (j *Job) MarkCompleted() {
	now := time.Now()
	j.Status = StatusCompleted
	j.CompletedAt = &now
}

// MarkFailed updates the job status to failed
func (j *Job) MarkFailed(err error) {
	j.Status = StatusFailed
	j.Error = err.Error()
}
