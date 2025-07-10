package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"thesilentcoder.com/m/server"
)

func main() {
	config, err := server.LoadConfig(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	zerolog.SetGlobalLevel(config.LogLevel)
	zerolog.DefaultContextLogger = &log.Logger

	if err := server.Start(context.Background(), *config); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
