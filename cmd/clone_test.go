package cmd

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/utils"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
	"strings"
	"testing"
)

func Test_setup(t *testing.T) {
	// This is not actually a test, it just sets up the logging for these test cases
	config.ConfigureLogger()
}

func Test_GivenOneArg_and_GivenNoNameFlag_DetermineOutputDir(t *testing.T) {
	cmdArgs := CloneCmdArguments{
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
	cmdArgs := CloneCmdArguments{
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

	cmdArgs := CloneCmdArguments{
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
	cmdArgs := CloneCmdArguments{
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

	var cmdArguments = CloneCmdArguments{
		args:        []string{sourceDir},
		nameFlag:    outputDir,
		isLocalPath: true,
		inputMethod: func(input string) string {
			/*
			* This function is being used to simulate a user's response to questions being asked.
			* It's 'input' string is the question in the actual implementation, and gets its input from stdout.
			* in this mocked version of that 'answerQuestion' function, we are simply returning the input as the response,
			* and checking that that value is present in the output directory.
			 */

			if strings.Contains(input, "(da default)") {
				return "" // simulates user responding blank for default-test.txt
			}
			if strings.Contains(input, "(should-not-be-returned)") {
				return "file_sub_dir_multi_diff_2" // this is what would be returned if a default was not chosen
			}
			if strings.Contains(input, "(global-should-be-returned)") {
				return "" // simulates user responding blank for default-test.txt
			}
			if strings.Contains(input, "(global-should-not-be-returned)") {
				return "some-other-variable" // this is what would be returned if a default was not chosen
			}

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
