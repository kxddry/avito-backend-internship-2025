package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/kxddry/avito-backend-internship-2025/internal/startup"
	"github.com/kxddry/avito-backend-internship-2025/pkg/logging"
)

func main() {
	cfg := new(startup.Config)
	if err := cfg.Read(); err != nil {
		log.Fatal().Err(err).Msg("failed to read configuration")
	}
	logging.SetupLogger(cfg.Debug)

	if err := fire(cfg); err != nil {
		log.Error().Err(err).Msg("failed to start fire")
		os.Exit(1)
	}
}
