package logger

import (
	"github.com/r2dtools/agent/config"
	"log"
	"os"
	"path"
)

type logType string
type logLevel int

const (
	errorType   logType = "error"
	warningType logType = "warning"
	infoType    logType = "info"
	debugType   logType = "debug"
)
const (
	errorLevel logLevel = iota
	warningLevel
	infoLevel
	debugLevel
)

var logTypeLevelMap = map[logType]logLevel{
	errorType:   errorLevel,
	warningType: warningLevel,
	infoType:    infoLevel,
	debugType:   debugLevel,
}

func init() {
	config := config.GetConfig()
	logDir := path.Dir(config.LogFile)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}
}

func write(message string, lType logType) {
	iConfig := config.GetConfig()
	configLogLevel := iConfig.LogLevel
	logLevel := logTypeLevelMap[lType]

	if int(logLevel) > configLogLevel {
		return
	}

	logFile, err := os.OpenFile(iConfig.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}

	defer logFile.Close()
	logger := log.New(logFile, string(lType)+": ", log.LstdFlags)
	logger.Println(message)
}

// Error logs the message as "error"
func Error(message string) {
	write(message, errorType)
}

// Warning logs the message as "warning"
func Warning(message string) {
	write(message, warningType)
}

// Debug logs the message as "debug"
func Debug(message string) {
	write(message, debugType)
}

// Info logs the message as "info"
func Info(message string) {
	write(message, infoType)
}
