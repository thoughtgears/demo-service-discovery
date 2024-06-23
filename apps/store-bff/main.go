package main

import (
	"fmt"
	"time"

	"github.com/thoughtgears/demo-service-discovery/apps/store-bff/handlers"
	"github.com/thoughtgears/demo-service-discovery/apps/store-bff/pkg/cfg"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var config *cfg.Config

func init() {
	var err error
	config, err = cfg.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func main() {
	router := gin.New()
	router.Use(gin.Recovery(), logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	))

	router.GET("/items", handlers.GetItems(config))

	log.Info().Msgf("server running on port %s", config.Port)
	if err := router.Run(fmt.Sprintf(":%s", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
