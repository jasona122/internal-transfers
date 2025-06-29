package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "info"
	}
	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		level = zerolog.DebugLevel
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(consoleWriter).
		Level(level).
		With().
		Timestamp().
		Logger()

	log.Info().Str("level", level.String()).Msg("Logger initialized")
}
