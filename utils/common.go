package utils

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

func CheckForError(err error, exitCode ...int) {
	var exitStatus int
	if len(exitCode) == 0 {
		exitStatus = 1
	} else {
		exitStatus = exitCode[0]
	}
	if err != nil {
		log.Error(err)
		os.Exit(exitStatus)
	}
}

func ThrowError(message string, exitCode ...int) {
	err := errors.New(message)
	CheckForError(err, exitCode[0])
}