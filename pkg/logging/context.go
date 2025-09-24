package logging

import (
	"context"

	"go.uber.org/zap"
)

const loggerContextKey = "logger_context_key"

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger)
	if !ok || logger == nil {
		return NewSugaredLogger("")
	}
	return logger
}

func AddLoggerToContext(ctx context.Context, logger *zap.SugaredLogger) (ctxWithLog context.Context) {
	return context.WithValue(ctx, loggerContextKey, logger)
}
