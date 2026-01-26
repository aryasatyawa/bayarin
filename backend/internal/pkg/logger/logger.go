package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes logger
func Init(env string) {
	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Pretty print for development
	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		})
	}

	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// Info logs info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Debug logs debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Error logs error message
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Fatal logs fatal error and exits
func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}

// With creates a new logger with fields
func With() zerolog.Logger {
	return log.With().Logger()
}
