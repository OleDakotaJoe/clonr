package utils

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

func CheckForError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func ThrowError(message string) error {
	err := errors.New(message)
	return err
}