package cmd

import (
	"clonr/config"
	"clonr/utils"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
	"testing"
)

func Test_setup(t *testing.T) {
	// This is not actually a test, it just sets up the logging for these test cases
	config.ConfigureLogger()
}

func Test_GivenOneArg_and_GivenNoNameFlag_DetermineOutputDir(t *testing.T) {
	cmdArgs := cloneCmdArguments{
		args:     []string{"testing-resources"},
		nameFlag: config.GlobalConfig().DefaultProjectName,
	}

	result, err := determineProjectName(&cmdArgs)

	if result != cmdArgs.nameFlag {
		t.Fatalf("Result is not equal to providedNameFlag: %s", cmdArgs.nameFlag)
	}
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func Test_GivenOneArg_and_GivenOneNameFlag_DetermineOutputDir(t *testing.T) {
	cmdArgs := cloneCmdArguments{
		args:     []string{"testing-resources"},
		nameFlag: "custom-name-flag",
	}
	result, err := determineProjectName(&cmdArgs)
	if result != cmdArgs.nameFlag {
		t.Fatalf("Expected result was \"%s\", but got %s ", cmdArgs.nameFlag, result)
	}
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func Test_GivenTwoArgs_and_GivenNoNameFlag_DetermineOutputDir(t *testing.T) {
	var expectedResult = "should-be-this-name"

	cmdArgs := cloneCmdArguments{
		args:     []string{"testing-resources", expectedResult},
		nameFlag: config.GlobalConfig().DefaultProjectName,
	}

	result, err := determineProjectName(&cmdArgs)

	if result != expectedResult {
		t.Fatalf("Expected result was \"should-be-this-name\", but got %s", result)
	}

	if err != nil {
		t.Fatalf("Error: %e", err)
	}
}

func Test_GivenTwoArgs_and_GivenOneNameFlag_DetermineOutputDir(t *testing.T) {
	cmdArgs := cloneCmdArguments{
		args:     []string{"testing-resources", "should-not-be-this-name"},
		nameFlag: "something-is-wrong",
	}

	result, err := determineProjectName(&cmdArgs)
	if result != "" {
		t.Fatalf("Expected result was \"\", but got %s", result)
	}

	if err == nil {
		t.Fatal("Expected an error but none was thrown")
	}
}

func Test_GivenValidUrl_ShouldNotThrowError_ValidateAndExtractUrl(t *testing.T) {
	validUrl := []string{"https://google.com"}
	_, err := validateAndExtractUrl(validUrl)
	if err != nil {
		t.Fatalf("Unexpected Error: %e", err)
	}
}

func Test_GivenInvalidUrl_ShouldThrowError_ValidateAndExtractUrl(t *testing.T) {
	invalidUrl := []string{"https//google.com"}
	_, err := validateAndExtractUrl(invalidUrl)
	if err == nil {
		t.Fatal("Expected an error, but none was thrown")
	}
}
func Test_GivenValidUrl_ShouldReturnUrl_ValidateAndExtractUrl(t *testing.T) {
	validUrl := []string{"https://google.com"}
	url, err := validateAndExtractUrl(validUrl)
	if err != nil {
		t.Fatalf("Unexpected Error: %e", err)
	}
	if url != validUrl[0] {
		t.Fatalf("Expected %s, but got %s.", validUrl, url)
	}
}

func Test_givenTemplateFile_processFiles(t *testing.T) {
	config.ConfigureLogger()
	sourceDir := config.TestConfig().SourceDir
	outputDir := config.TestConfig().OutputDir
	answerKeyDir := config.TestConfig().AnswerKeyDir

	var cmdArguments = cloneCmdArguments{
		args:        []string{sourceDir},
		nameFlag:    outputDir,
		isLocalPath: true,
		inputMethod: func(input string) string {
			return input
		},
	}

	cloneProject(&cmdArguments)

	// Check the hash of the directories
	actualHash, actErr := dirhash.HashDir(outputDir, "test", dirhash.DefaultHash)
	utils.CheckForError(actErr)
	expectedHash, expErr := dirhash.HashDir(answerKeyDir, "test", dirhash.DefaultHash)
	utils.CheckForError(expErr)

	if actualHash != expectedHash {
		t.Fatal("output was not correct")
	}

	// Cleanup the test directory
	err := os.RemoveAll(outputDir)
	utils.CheckForError(err)
}
