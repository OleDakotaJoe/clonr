package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"
)

func Test_GivenRelativePath_ExpectAbsolutePath_makeAliasMap(t *testing.T) {
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "./hello"

	args := types.AliasCmdArgs{
		IsLocalFlag:         true,
		ActualAliasName:     testAliasName,
		ActualAliasLocation: testAliasLocation,
	}

	// Set up expected result
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Something went wrong with the test, this is probably not a problem with the code.")
	}
	expectedLocation := fmt.Sprintf("%s/hello", pwd)
	expectedResult := map[string]interface{}{testAliasName: map[string]interface{}{
		config.Global().AliasesUrlKey:            expectedLocation,
		config.Global().AliasesLocalIndicatorKey: true,
	}}

	// Run the test
	actualResult := makeAliasMap(&args)

	equal := reflect.DeepEqual(expectedResult, actualResult)

	if !equal {
		t.Fatalf("The maps did not match! Actual: %s | Expected: %s", actualResult, expectedResult)
	}

}

func Test_GivenAbsolutePath_ExpectAbsolutePath_makeAliasMap(t *testing.T) {
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		IsLocalFlag:         true,
		ActualAliasName:     testAliasName,
		ActualAliasLocation: testAliasLocation,
	}

	// Set up expected result
	expectedResult := map[string]interface{}{testAliasName: map[string]interface{}{
		config.Global().AliasesUrlKey:            testAliasLocation,
		config.Global().AliasesLocalIndicatorKey: true,
	}}

	// Run the test
	actualResult := makeAliasMap(&args)

	equal := reflect.DeepEqual(expectedResult, actualResult)

	if !equal {
		t.Fatalf("The maps did not match! Actual: %s | Expected: %s", actualResult, expectedResult)
	}

}

func Test_GivenOneArgAndNoNameFlag_ExpectArgToBeSetAsName_setNameForAlias(t *testing.T) {
	config.ConfigureLogger()
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "./hello"

	args := types.AliasCmdArgs{
		Args:                []string{testAliasName},
		AddFlag:             true,
		UpdateFlag:          false,
		DeleteFlag:          false,
		IsLocalFlag:         true,
		AliasNameFlag:       "",
		ActualAliasName:     testAliasName,
		AliasLocationFlag:   "",
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return prompt
		},
	}

	setNameForAlias(&args)

	if args.ActualAliasName != testAliasName {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasName, args.ActualAliasName)
	}
}

func Test_GivenNoArgsNoNameFlagsAndDeleteFlag_ExpectPromptForName_setNameForAlias(t *testing.T) {
	config.ConfigureLogger()
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"
	addAlias(testAliasName, testAliasLocation, false)

	args := types.AliasCmdArgs{
		Args:                []string{testAliasName},
		AddFlag:             false,
		UpdateFlag:          false,
		DeleteFlag:          true,
		IsLocalFlag:         true,
		AliasNameFlag:       "",
		ActualAliasName:     testAliasName,
		AliasLocationFlag:   "",
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			if prompt == "Which alias do you want to delete?" {
				return testAliasName
			}
			return prompt
		},
	}

	setNameForAlias(&args)

	if args.ActualAliasName != testAliasName {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasName, args.ActualAliasName)
	}
}

func addAlias(aliasName string, aliasLocation string, isLocal bool) {
	args := types.AliasCmdArgs{
		Args:                nil,
		AddFlag:             true,
		UpdateFlag:          false,
		DeleteFlag:          false,
		IsLocalFlag:         isLocal,
		AliasNameFlag:       "",
		ActualAliasName:     aliasName,
		AliasLocationFlag:   "",
		ActualAliasLocation: aliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return prompt
		},
	}

	aliasManager(&args)
}
