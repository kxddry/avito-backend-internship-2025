package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kxddry/avito-backend-internship-2025/internal/startup"
)

func fire(cfg *startup.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// placeholder
	_ = ctx
	return nil
}
