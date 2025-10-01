package logger

import (
	"os"

	"github.com/nineee02/gotest/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func newZapLogger(cfg *config.Configuration) *zap.Logger {
	logLevel := zap.DebugLevel // Default level
	if level, exists := loggerLevelMap[cfg.Logger.Level]; exists {
		logLevel = level
	}

	var encoderCfg zapcore.EncoderConfig
	if cfg.Logger.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.LineEnding = zapcore.DefaultLineEnding
	encoderCfg.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.TimeKey = "ts"
	encoderCfg.NameKey = "name"
	encoderCfg.MessageKey = "msg"

	var encoder zapcore.Encoder
	if cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	cores := []zapcore.Core{}
	level := zap.NewAtomicLevelAt(logLevel)
	consoleWriter := zapcore.AddSync(zapcore.Lock(os.Stdout))
	cores = append(cores, zapcore.NewCore(encoder, consoleWriter, level))

	if cfg.Logger.Path != "" {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Logger.Path,
			MaxSize:    10, // Megabytes
			MaxBackups: 3,
			MaxAge:     28, // Days
			Compress:   true,
		})
		cores = append(cores, zapcore.NewCore(encoder, writer, level))
	}

	core := zapcore.NewTee(cores...)

	return zap.New(core, zap.AddCaller())
}
