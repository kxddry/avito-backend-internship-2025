//go:build pprof

package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/rs/zerolog/log"
)

func init() {
	go func() {
		log.Info().Str("addr", "0.0.0.0:6060").Msg("Starting pprof server")
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
}
