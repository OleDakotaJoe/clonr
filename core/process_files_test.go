package core

import (
	"clonr/config"
	"clonr/utils"
	"github.com/otiai10/copy"
	"golang.org/x/mod/sumdb/dirhash"
	"testing"
)

func Test_setup(t *testing.T) {
	// This is not really a test, it is just a setup function.
	config.ConfigureLogger()
}

func Test_givenTemplateFile_processFiles(t *testing.T) {
	config.ConfigureLogger()
	sourceTemplate := "../testing-resources/process_files_test/source_template"
	testOutputDirectory := "../testing-resources/process_files_test/test_output"
	testTemplateExpectedResult := "../testing-resources/process_files_test/answer_key"

	// This step just sets up the test. Copies the test project so that doesn't manipulate the original.
	copyErr := copy.Copy(sourceTemplate, testOutputDirectory)
	utils.CheckForError(copyErr)
	ProcessFiles(testOutputDirectory)

	// Create the sample map necessary to process the file.
	actualHash, actErr := dirhash.HashDir(testOutputDirectory, "test", dirhash.DefaultHash)
	utils.CheckForError(actErr)
	expectedHash, expErr := dirhash.HashDir(testTemplateExpectedResult, "test", dirhash.DefaultHash)
	utils.CheckForError(expErr)

	if actualHash != expectedHash {
		t.Fatal("output was not correct")
	}
}
