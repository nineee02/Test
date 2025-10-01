package logger

import (
	"context"

	"github.com/nineee02/gotest/pkg/app_context"
	"github.com/nineee02/gotest/pkg/config"
	"go.uber.org/zap"
)

type Logger interface {
	Named(s string) Logger
	With(fields ...zap.Field) Logger
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	WithFields(fields ...zap.Field) *zap.Logger
	Sync() error
	Close()
}

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(cfg *config.Configuration) Logger {
	log := newZapLogger(cfg)
	return &zapLogger{logger: log}
}

func NewNoOpLogger() Logger {
	log := zap.NewNop()
	return &zapLogger{logger: log}
}

func extractRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(app_context.RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

func (l *zapLogger) Named(s string) Logger {
	return &zapLogger{logger: l.logger.Named(s)}
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	requestId := extractRequestID(ctx)
	if requestId != "" {
		fields = append(fields, zap.String("request_id", requestId))
	}
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	requestId := extractRequestID(ctx)
	if requestId != "" {
		fields = append(fields, zap.String("request_id", requestId))
	}
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	requestId := extractRequestID(ctx)
	if requestId != "" {
		fields = append(fields, zap.String("request_id", requestId))
	}
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	requestId := extractRequestID(ctx)
	if requestId != "" {
		fields = append(fields, zap.String("request_id", requestId))
	}
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *zapLogger) Close() {
	_ = l.logger.Sync()
}

func (l *zapLogger) WithFields(fields ...zap.Field) *zap.Logger {
	return l.logger.With(fields...)
}
