package cmd

import (
	"clonr/config"
	"testing"
)

func Test_setup(t *testing.T) {
	// This is not actually a test, it just sets up the logging for these test cases
	config.ConfigureLogger()
}

func Test_DetermineOutputDir_GivenOneArg_and_GivenNoNameFlag(t *testing.T) {
	
	mockArgs := []string{"testing"}
	providedNameFlag := config.DefaultConfig().DefaultProjectName
	result, err  := determineOutputDir(providedNameFlag, mockArgs)
	if result != providedNameFlag || err != nil {
		t.Fatal("FAILED TEST")
	}
}


func Test_DetermineOutputDir_GivenOneArg_and_GivenOneNameFlag(t *testing.T) {
	mockArgs := []string{"testing"}
	providedNameFlag := "custom-name-flag"
	result, err  := determineOutputDir(providedNameFlag, mockArgs)
	if result != providedNameFlag|| err != nil  {
		t.Fatal("FAILED TEST")
	}
}


func Test_DetermineOutputDir_GivenTwoArgs_and_GivenNoNameFlag(t *testing.T) {
	mockArgs := []string{"testing", "should-be-this-name"}
	providedNameFlag := config.DefaultConfig().DefaultProjectName
	result, err := determineOutputDir(providedNameFlag, mockArgs)
	if result != "should-be-this-name" || err != nil {
		t.Fatal("FAILED TEST")
	}
}

func Test_DetermineOutputDir_GivenTwoArgs_and_GivenOneNameFlag(t *testing.T) {
	mockArgs := []string{"testing", "should-not-be-this-name"}
	providedNameFlag := "something-is-wrong"
	result, err := determineOutputDir(providedNameFlag, mockArgs)
	if err == nil || !(result == "should-be-this-name")  {
		t.Fatal("FAILED TEST")
	}
}

