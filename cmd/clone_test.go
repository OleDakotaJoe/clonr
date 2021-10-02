package cmd

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func Test_DetermineOutputDir_GivenOneArg(t *testing.T) {
	mockArgs := []string{"testing"}
	nameFlag := "clonr-app"
	result := determineOutputDir(nameFlag, mockArgs)
	if result != "testing" {
		t.Fatal("FAILED TEST")
	}
	log.Info(mockArgs[0])
}
