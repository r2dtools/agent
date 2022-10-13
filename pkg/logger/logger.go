package logger

import (
	"os"
	"path"
	"time"

	"github.com/r2dtools/agent/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerInterface interface {
	Error(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

type Logger struct {
	zapLogger *zap.SugaredLogger
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.zapLogger.Errorf(message, args...)
}

func (l *Logger) Warning(message string, args ...interface{}) {
	l.zapLogger.Warnf(message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.zapLogger.Infof(message, args...)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.zapLogger.Debugf(message, args...)
}

func NewLogger(config *config.Config) (LoggerInterface, error) {
	logDir := path.Dir(config.GetLoggerFileAbsPath())

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	var loggerConfig zap.Config
	outputPaths := []string{config.GetLoggerFileAbsPath()}

	if config.IsProdMode {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
		outputPaths = append(outputPaths, "stderr")
	}

	loggerConfig.OutputPaths = outputPaths
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger: logger.Sugar()}, nil
}
