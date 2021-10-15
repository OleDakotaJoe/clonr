package config

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"os"
)

type LoggerConfig struct {
	Writer    io.Writer
	InfoLevel log.Level
	Formatter log.Formatter
}

func ConfigureLogger() {
	config := LoggerConfig{
		Writer:    os.Stdout,
		InfoLevel: log.InfoLevel,
		Formatter: &easy.Formatter{
			LogFormat: "[%lvl%]: %msg%\n",
		},
	}
	log.SetOutput(config.Writer)
	log.SetLevel(config.InfoLevel)
	log.SetFormatter(config.Formatter)
}
