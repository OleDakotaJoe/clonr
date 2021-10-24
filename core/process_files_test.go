package core

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
	"strings"
	"testing"
)

func Test_givenTemplateFile_checkForValidation_ProcessFiles(t *testing.T) {
	config.ConfigureLogger()
	/**
	* TEST SETUP
	 */
	sourceDir := config.TestConfig().ValidationTestSource
	outputDir := config.TestConfig().ValidationTestOutput
	answerKeyDir := config.TestConfig().ValidationTestAnswer
	// Since this is calling "ProcessFiles" we need to copy the directory for the test
	err := copy.Copy(sourceDir, outputDir)
	v, vErr := utils.ViperReadConfig(outputDir, config.Global().ConfigFileName, config.Global().ConfigFileType)
	utils.ExitIfError(vErr)

	utils.ExitIfError(err)
	retryCounter := 0
	var fileProcessorSettings = types.FileProcessorSettings{
		Viper:          *v,
		ConfigFilePath: outputDir,
		StringInputReader: func(input string) string {
			/**
			* This function is being used to simulate a user's response to questions being asked.
			* It's 'input' string is the question in the actual implementation, and gets its input from stdout.
			* in this mocked version of that 'answerQuestion' function, we are simply returning the input as the response,
			* and checking that that value is present in the output directory.
			 */

			if strings.Contains(input, "(pick_me)") {
				return "" // simulates user responding blank, indicating they've chosen the default answer
			}

			if strings.Contains(input, "!@#$%^") {
				if retryCounter < 5 {
					log.Info(retryCounter)
					retryCounter++
					return "!@#$%^"
				}
				return "is_not_valid"
			}

			return input
		},
	}

	ProcessFiles(&fileProcessorSettings)

	// Check the hash of the directories
	actualHash, actErr := dirhash.HashDir(outputDir, "test", dirhash.DefaultHash)
	utils.ExitIfError(actErr)
	expectedHash, expErr := dirhash.HashDir(answerKeyDir, "test", dirhash.DefaultHash)
	utils.ExitIfError(expErr)

	if actualHash != expectedHash {
		t.Fatal("output was not correct")
	}

	// Cleanup the test directory only if tests pass, so you can look at result if it fails
	osErr := os.RemoveAll(outputDir)
	utils.ExitIfError(osErr)
}
