package logging

import (
	"sync"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
)

func Init() {
	once.Do(func() {
		zap.ReplaceGlobals(newLocalLogger())
		zap.S().Info("logger initialized in development mode")
	})
}

func newLocalLogger() *zap.Logger {
	logCfg := zap.NewDevelopmentEncoderConfig()
	logCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(logCfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
}

func NewSugaredLogger(name string) *zap.SugaredLogger {
	Init()
	return zap.S().Named(name)
}
