package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/octopunkio/taskflow-worker/internal/config"
	"github.com/octopunkio/taskflow-worker/internal/handler"
	"github.com/octopunkio/taskflow-worker/internal/job"
)

// NotificationType represents different notification channels
type NotificationType string

const (
	NotificationTypePush    NotificationType = "push"
	NotificationTypeInApp   NotificationType = "in_app"
	NotificationTypeSlack   NotificationType = "slack"
	NotificationTypeDiscord NotificationType = "discord"
)

// NotificationPayload represents notification data
type NotificationPayload struct {
	UserID  string           `json:"user_id"`
	Type    NotificationType `json:"type"`
	Title   string           `json:"title"`
	Message string           `json:"message"`
	Data    map[string]any   `json:"data,omitempty"`
}

// NotificationHandler handles notification jobs
type NotificationHandler struct {
	config *config.Config
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(cfg *config.Config) handler.Handler {
	return &NotificationHandler{config: cfg}
}

// Handle processes a notification job
func (h *NotificationHandler) Handle(ctx context.Context, j *job.Job) error {
	var payload NotificationPayload
	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	switch payload.Type {
	case NotificationTypePush:
		return h.sendPushNotification(ctx, &payload)
	case NotificationTypeInApp:
		return h.sendInAppNotification(ctx, &payload)
	case NotificationTypeSlack:
		return h.sendSlackNotification(ctx, &payload)
	case NotificationTypeDiscord:
		return h.sendDiscordNotification(ctx, &payload)
	default:
		return fmt.Errorf("unknown notification type: %s", payload.Type)
	}
}

func (h *NotificationHandler) sendPushNotification(ctx context.Context, p *NotificationPayload) error {
	fmt.Printf("Push notification to %s: %s\n", p.UserID, p.Title)
	return nil
}

func (h *NotificationHandler) sendInAppNotification(ctx context.Context, p *NotificationPayload) error {
	fmt.Printf("In-app notification to %s: %s\n", p.UserID, p.Title)
	return nil
}

func (h *NotificationHandler) sendSlackNotification(ctx context.Context, p *NotificationPayload) error {
	fmt.Printf("Slack notification: %s\n", p.Message)
	return nil
}

func (h *NotificationHandler) sendDiscordNotification(ctx context.Context, p *NotificationPayload) error {
	fmt.Printf("Discord notification: %s\n", p.Message)
	return nil
}
