package log

import (
	"os"

	"github.com/hex-techs/blade/pkg/util/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger

var logLevel = zap.NewAtomicLevel()

func InitLogger() *zap.SugaredLogger {
	SetLevel()
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugaredLogger = logger.Sugar()
	return sugaredLogger
}

func SetLevel() {
	switch config.Read().Log.Level {
	case "debug":
		logLevel.SetLevel(zapcore.Level(zapcore.DebugLevel))
	case "info":
		logLevel.SetLevel(zapcore.Level(zapcore.InfoLevel))
	case "warn":
		logLevel.SetLevel(zapcore.Level(zapcore.WarnLevel))
	case "error":
		logLevel.SetLevel(zapcore.Level(zapcore.ErrorLevel))
	default:
		logLevel.SetLevel(zapcore.Level(zapcore.DebugLevel))
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if config.Read().Log.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	if config.Read().Log.Output == "file" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   config.Read().Log.Filename,
			MaxSize:    config.Read().Log.MaxSize,
			MaxBackups: config.Read().Log.MaxBackups,
			MaxAge:     config.Read().Log.MaxAge,
			Compress:   config.Read().Log.Compress,
		}
		return zapcore.AddSync(lumberJackLogger)
	}
	return zapcore.AddSync(zapcore.Lock(os.Stdout))
}

func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

func Infow(msg string, args ...interface{}) {
	sugaredLogger.Infow(msg, args...)
}

func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

func Warnw(msg string, args ...interface{}) {
	sugaredLogger.Warnw(msg, args...)
}

func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

func Errorw(msg string, args ...interface{}) {
	sugaredLogger.Errorw(msg, args...)
}

func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}

func Fatalw(msg string, args ...interface{}) {
	sugaredLogger.Fatalw(msg, args...)
}
