package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/octopunkio/taskflow-worker/internal/config"
	"github.com/octopunkio/taskflow-worker/internal/handler"
	"github.com/octopunkio/taskflow-worker/internal/job"
)

// WebhookPayload represents webhook delivery data
type WebhookPayload struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    json.RawMessage   `json:"body"`
}

// WebhookHandler handles webhook delivery jobs
type WebhookHandler struct {
	config *config.Config
	client *http.Client
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(cfg *config.Config) handler.Handler {
	return &WebhookHandler{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Handle processes a webhook job
func (h *WebhookHandler) Handle(ctx context.Context, j *job.Job) error {
	var payload WebhookPayload
	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	method := payload.Method
	if method == "" {
		method = http.MethodPost
	}

	req, err := http.NewRequestWithContext(ctx, method, payload.URL, bytes.NewReader(payload.Body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range payload.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook failed with status %d", resp.StatusCode)
	}

	return nil
}
