package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

func SetupLogger(debug bool) {
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	if term.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = zerolog.
			New(zerolog.ConsoleWriter{
				Out: os.Stdout,
			}).
			Level(level).
			With().Timestamp().
			Logger()
	} else {
		log.Logger = zerolog.
			New(os.Stdout).
			Level(level).
			With().Timestamp().
			Logger()
	}
}
