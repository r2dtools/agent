package logger

import (
	"os"
	"path"
	"time"

	"github.com/r2dtools/agent/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

type logger struct {
	zapLogger *zap.SugaredLogger
}

func (l *logger) Error(message string, args ...interface{}) {
	l.zapLogger.Errorf(message, args...)
}

func (l *logger) Warning(message string, args ...interface{}) {
	l.zapLogger.Warnf(message, args...)
}

func (l *logger) Info(message string, args ...interface{}) {
	l.zapLogger.Infof(message, args...)
}

func (l *logger) Debug(message string, args ...interface{}) {
	l.zapLogger.Debugf(message, args...)
}

func NewLogger(config *config.Config) (Logger, error) {
	logDir := path.Dir(config.GetLoggerFileAbsPath())

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)

		if err != nil {
			return nil, err
		}
	}

	var loggerConfig zap.Config
	outputPaths := []string{config.GetLoggerFileAbsPath()}

	if config.IsDevMode {
		loggerConfig = zap.NewDevelopmentConfig()
		outputPaths = append(outputPaths, "stderr")
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	loggerConfig.OutputPaths = outputPaths
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	zLogger, err := loggerConfig.Build()

	if err != nil {
		return nil, err
	}

	return &logger{zapLogger: zLogger.Sugar()}, nil
}

type NilLogger struct{}

func (l *NilLogger) Error(message string, args ...interface{}) {
}

func (l *NilLogger) Warning(message string, args ...interface{}) {
}

func (l *NilLogger) Info(message string, args ...interface{}) {
}

func (l *NilLogger) Debug(message string, args ...interface{}) {
}
