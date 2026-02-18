package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/octopunkio/taskflow-worker/internal/config"
	"github.com/octopunkio/taskflow-worker/internal/handler"
	"github.com/octopunkio/taskflow-worker/internal/job"
)

// CleanupPayload represents cleanup job configuration
type CleanupPayload struct {
	Target    string `json:"target"`
	OlderThan string `json:"older_than"`
	DryRun    bool   `json:"dry_run"`
}

// CleanupHandler handles data cleanup jobs
type CleanupHandler struct {
	config *config.Config
}

// NewCleanupHandler creates a new CleanupHandler
func NewCleanupHandler(cfg *config.Config) handler.Handler {
	return &CleanupHandler{config: cfg}
}

// Handle processes a cleanup job
func (h *CleanupHandler) Handle(ctx context.Context, j *job.Job) error {
	var payload CleanupPayload
	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	duration, err := time.ParseDuration(payload.OlderThan)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	cutoff := time.Now().Add(-duration)

	if payload.DryRun {
		fmt.Printf("[DRY RUN] Would clean %s older than %v\n", payload.Target, cutoff)
		return nil
	}

	fmt.Printf("Cleaning %s older than %v\n", payload.Target, cutoff)
	// TODO: Implement actual cleanup logic based on target type

	return nil
}
