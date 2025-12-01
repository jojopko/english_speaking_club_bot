package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"speaking-club-bot/internal/app"
	"speaking-club-bot/internal/config"
	"syscall"
)

func main() {
	config, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config:", "err", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := app.New(&config)
	if err != nil {
		slog.Error("Failed to create App:", "err", err)
		os.Exit(1)
	}

	if err := app.Start(ctx); err != nil {
		slog.Error("App finished with error:", "err", err)
		os.Exit(1)
	}
}
