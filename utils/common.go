package utils

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

func CheckForError(err error, exitCode ...int) {
	if err != nil {
		log.Error(err)
		os.Exit(determineExitStatus(exitCode[0])[0])
	}
}

func ThrowError(message string, exitCode ...int) ( []int, error) {
	err := errors.New(message)
	return determineExitStatus(exitCode[0]), err
}

func determineExitStatus(exitCode ...int) []int {
	var exitStatus int
	if len(exitCode) == 0 {
		exitStatus = 1
	} else {
		exitStatus = exitCode[0]
	}

	return []int {exitStatus}
}