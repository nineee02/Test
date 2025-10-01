package logger

import (
	"context"
	"time"

	"github.com/nineee02/gotest/pkg/config"

	"github.com/nineee02/gotest/pkg/constant"
	"go.uber.org/zap"
	g_logger "gorm.io/gorm/logger"
)

const (
	slowThreshold = 200 * time.Millisecond
)

type Gorm interface {
	LogMode(level g_logger.LogLevel) g_logger.Interface
	Info(_ context.Context, msg string, data ...interface{})
	Warn(_ context.Context, msg string, data ...interface{})
	Error(_ context.Context, msg string, data ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
}

type GormLogger struct {
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

var _ Gorm = (*GormLogger)(nil)

func NewGormLogger(cfg *config.Configuration) *GormLogger {
	l := newZapLogger(cfg)
	return &GormLogger{
		sugarLogger: l.Sugar(),
		logger:      l,
	}
}

func (l *GormLogger) LogMode(level g_logger.LogLevel) g_logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	requestId := ctx.Value(constant.RequestIdKey).(string)
	l.sugarLogger.
		With(zap.String(constant.RequestIdKey, requestId)).
		Infof("[gorm] %s: %+v", msg, data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	requestId := ctx.Value(constant.RequestIdKey).(string)
	l.sugarLogger.
		With(zap.String(constant.RequestIdKey, requestId)).
		Warnf("[gorm] %s: %+v", msg, data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	requestId := ctx.Value(constant.RequestIdKey).(string)
	l.sugarLogger.
		With(zap.String(constant.RequestIdKey, requestId)).
		Errorf("[gorm] %s: %+v", msg, data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	var (
		sql  string
		rows int64
	)

	if fc != nil {
		sqlStr, rowCnt := fc()
		sql = sqlStr
		rows = rowCnt
	} else {
		sql = "N/A"
		rows = 0
	}

	requestId, _ := ctx.Value(constant.RequestIdKey).(string) // ป้องกัน panic หากไม่มี requestId
	lgr := l.sugarLogger.With(
		zap.String(constant.RequestIdKey, requestId),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	)

	switch {
	case err != nil:
		lgr.With(zap.Error(err)).Error("SQL execution error")

	case elapsed > slowThreshold:
		lgr.Warn("Slow SQL query")

	default:
		lgr.Debug("SQL executed")
	}
}
