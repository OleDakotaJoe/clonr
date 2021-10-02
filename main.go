package main

import (
	"clonr/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	cmd.Execute()
}

func init() {
	// Set up Logging
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}