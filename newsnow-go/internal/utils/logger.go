package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	
	log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Success(v interface{}) {
	log.Info().Msgf("%v", v)
}

func Info(v interface{}) {
	log.Info().Msgf("%v", v)
}

func Warn(v interface{}) {
	log.Warn().Msgf("%v", v)
}

func Error(v interface{}) {
	log.Error().Msgf("%v", v)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Fatal(v interface{}) {
	log.Fatal().Msgf("%v", v)
}
