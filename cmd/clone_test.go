package cmd

import (
	"clonr/config"
	"testing"
)

func Test_setup(t *testing.T) {
	// This is not actually a test, it just sets up the logging for these test cases
	config.ConfigureLogger()
}

func Test_GivenOneArg_and_GivenNoNameFlag_DetermineOutputDir(t *testing.T) {
	mockArgs := []string{"testing"}
	providedNameFlag := config.GlobalConfig().DefaultProjectName
	result, err  := determineProjectName(providedNameFlag, mockArgs)
	if result != providedNameFlag {
		t.Fatalf("Result is not equal to providedNameFlag: %s", providedNameFlag)
	}
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func Test_GivenOneArg_and_GivenOneNameFlag_DetermineOutputDir(t *testing.T) {
	mockArgs := []string{"testing"}
	providedNameFlag := "custom-name-flag"
	result, err  := determineProjectName(providedNameFlag, mockArgs)
	if result != providedNameFlag {
		t.Fatalf("Expected result was \"%s\", but got %s ", providedNameFlag, result)
	}
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func Test_GivenTwoArgs_and_GivenNoNameFlag_DetermineOutputDir(t *testing.T) {
	mockArgs := []string{"testing", "should-be-this-name"}
	providedNameFlag := config.GlobalConfig().DefaultProjectName
	result, err := determineProjectName(providedNameFlag, mockArgs)


	if result != "should-be-this-name" {
		t.Fatalf("Expected result was \"should-be-this-name\", but got %s", result)
	}

	if err != nil {
		t.Fatalf("Error: %e", err)
	}
}

func Test_GivenTwoArgs_and_GivenOneNameFlag_DetermineOutputDir(t *testing.T) {
	mockArgs := []string{"testing", "should-not-be-this-name"}
	providedNameFlag := "something-is-wrong"
	result, err := determineProjectName(providedNameFlag, mockArgs)
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
		t.Fatalf("Expected %s, but got %s.", validUrl, url )
	}
}

