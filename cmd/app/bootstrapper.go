package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kxddry/avito-backend-internship-2025/internal/service"
	"github.com/kxddry/avito-backend-internship-2025/internal/startup"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage/txmanager"
)

func fire(cfg *startup.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dsn := cfg.DBConfig.DSN()

	txMgr, err := txmanager.New(ctx, dsn)
	if err != nil {
		return err
	}
	defer txMgr.Close()

	svc := service.New(service.Dependencies{
		TransactionManager: txMgr,
	})

	app, err := startup.NewApplication(cfg, svc)
	if err != nil {
		return err
	}

	return app.Run(ctx)
}
