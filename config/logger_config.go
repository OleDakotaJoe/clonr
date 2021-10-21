package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"os"
	"strings"
)

type LoggerConfig struct {
	Writer    io.Writer
	InfoLevel log.Level
	Formatter log.Formatter
}

func ConfigureLogger() {
	config := LoggerConfig{
		Writer:    os.Stdout,
		InfoLevel: getLogLevel(),
		Formatter: &easy.Formatter{
			LogFormat: "[%lvl%]: %msg%\n",
		},
	}
	log.SetOutput(config.Writer)
	log.SetLevel(getLogLevel())
	log.SetFormatter(config.Formatter)
}

func getLogLevel() log.Level {
	logLevel := strings.ToLower(Global().LogLevel)

	switch logLevel {
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		fmt.Printf("The log level you specified {%s} is not valid.", logLevel)
		return log.InfoLevel
	}
}
