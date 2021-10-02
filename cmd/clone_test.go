package cmd

import (
	"testing"
)

/* CONFIGURATION VARIABLES */

func Test_DetermineOutputDir_GivenOneArg_and_GivenNoNameFlag(t *testing.T) {
	mockArgs := []string{"testing"}
	defaultNameFlag := "default-name-flag"
	providedNameFlag := ""
	result := determineOutputDir( defaultNameFlag, providedNameFlag, mockArgs)
	if result != defaultNameFlag {
		t.Fatal("FAILED TEST")
	}
}
func Test_DetermineOutputDir_GivenOneArg_and_GivenOneNameFlag(t *testing.T) {
	mockArgs := []string{"testing"}
	defaultNameFlag := "default-name-flag"
	providedNameFlag := "provided-name-flag"
	result := determineOutputDir( defaultNameFlag, providedNameFlag, mockArgs)
	if result != providedNameFlag {
		t.Fatal("FAILED TEST")
	}
}

func Test_DetermineOutputDir_GivenTwoArgs_and_GivenNoNameFlag(t *testing.T) {
	mockArgs := []string{"testing", "should-be-this-name"}
	defaultNameFlag := "default-name-flag"
	nameFlag := ""
	result := determineOutputDir( defaultNameFlag,nameFlag, mockArgs)
	if result != mockArgs[1] {
		t.Fatal("FAILED TEST")
	}
}

//func Test_DetermineOutputDir_GivenTwoArgs_and_GivenOneNameFlag(t *testing.T) {
//	mockArgs := []string{"testing", "should-be-this-name"}
//	defaultNameFlag := "default-name-flag"
//	nameFlag := "something-is-wrong"
//	result := determineOutputDir( defaultNameFlag,nameFlag, mockArgs)
//	if result != mockArgs[1] {
//		t.Fatal("FAILED TEST")
//	}
//}

