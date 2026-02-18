package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/octopunkio/taskflow-worker/internal/config"
	"github.com/octopunkio/taskflow-worker/internal/handlers"
	"github.com/octopunkio/taskflow-worker/internal/worker"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.Load()

	w := worker.New(worker.Config{
		RedisAddr:   cfg.RedisAddr,
		Concurrency: cfg.Concurrency,
		Logger:      logger,
	})

	// Register job handlers
	w.RegisterHandler("email", handlers.NewEmailHandler(cfg))
	w.RegisterHandler("notification", handlers.NewNotificationHandler(cfg))
	w.RegisterHandler("webhook", handlers.NewWebhookHandler(cfg))
	w.RegisterHandler("cleanup", handlers.NewCleanupHandler(cfg))

	// Start worker
	ctx, cancel := context.WithCancel(context.Background())
	go w.Start(ctx)

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down worker...")
	cancel()
	w.Wait()
}
