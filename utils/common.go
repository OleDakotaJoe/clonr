package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func CheckForError(err error, exitCode ...int) {
	if &exitCode == nil {
		exitCode[0] = 1
	}

	if err != nil {
		log.Error(err)
		os.Exit(exitCode[0])
	}
}