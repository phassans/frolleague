package common

import (
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// InitLogger is to initialize a logger
func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "viraagh").
		Logger()
}

// GetLogger is to get Logger
func GetLogger() zerolog.Logger {
	return logger
}
