package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/octopunkio/taskflow-worker/internal/config"
	"github.com/octopunkio/taskflow-worker/internal/handler"
	"github.com/octopunkio/taskflow-worker/internal/job"
)

// EmailPayload represents the data needed to send an email
type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	HTML    bool   `json:"html"`
}

// EmailHandler handles email sending jobs
type EmailHandler struct {
	config *config.Config
}

// NewEmailHandler creates a new EmailHandler
func NewEmailHandler(cfg *config.Config) handler.Handler {
	return &EmailHandler{config: cfg}
}

// Handle processes an email job
func (h *EmailHandler) Handle(ctx context.Context, j *job.Job) error {
	var payload EmailPayload
	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	// TODO: Implement actual email sending via SMTP or API
	fmt.Printf("Sending email to %s: %s\n", payload.To, payload.Subject)

	return nil
}
