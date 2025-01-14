package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the smallest interface we use of go-uber/zap's sugared logger
// Defining this interface makes it easier for us to control what functions we use, and to pass the logger around to various packages
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})

	Infow(msg string, keysAndValues ...interface{})

	Warnw(msg string, keysAndValues ...interface{})

	Errorw(msg string, keysAndValues ...interface{})

	Fatal(args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Sync() error
}

type contextKeyType string

var (
	contextKey = contextKeyType("logger")
)

func FromContext(ctx context.Context) Logger {
	if v := ctx.Value(contextKey); v != nil {
		return v.(Logger)
	}

	log.Fatal("No logger found in context")
	return nil
}

func OnContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func New(logLevel zap.AtomicLevel, logDevelopment bool) Logger {
	zapConfig := zap.Config{
		Level:            logLevel,
		Development:      logDevelopment,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := zapConfig.Build()

	// Just until we've remove default log package usage in the project
	zap.RedirectStdLog(logger)

	return logger.Sugar()
}
