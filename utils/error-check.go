package utils

import (
	"fmt"
	"os"
)

func CheckForError(err error, exitCode ...int) {
	if &exitCode == nil {
		exitCode[0] = 1
	}

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode[0])
	}
}