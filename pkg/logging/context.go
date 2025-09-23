package logging

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

var loggerContextKey = ctxKey{}

func LoggerFromContext(ctx context.Context) zerolog.Logger {
	logger, ok := ctx.Value(loggerContextKey).(zerolog.Logger)
	if !ok {
		return zerolog.Nop()
	}
	return logger
}

func AddLoggerToContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// NewLogger creates a new logger with pretty console output
func NewLogger() zerolog.Logger {
	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure console writer with colors
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			return zerolog.LevelFieldName + "=" + i.(string)
		},
		FormatMessage: func(i interface{}) string {
			return "msg=" + i.(string)
		},
		FormatFieldName: func(i interface{}) string {
			return i.(string) + "="
		},
		FormatFieldValue: func(i interface{}) string {
			return i.(string)
		},
	}

	return zerolog.New(output).With().Timestamp().Logger()
}
