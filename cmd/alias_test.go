package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"os"
	"reflect"
	"testing"
)

func Test_GivenRelativePath_ExpectAbsolutePath_makeAliasMap(t *testing.T) {
	setupTests(t)
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
	setupTests(t)
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
	setupTests(t)
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
		ActualAliasName:     "",
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

func Test_GivenNoArgsAndNoNameFlagAndDeleteFlag_ExpectPromptForName_setNameForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"
	addAlias(testAliasName, testAliasLocation, false)

	args := types.AliasCmdArgs{
		Args:                []string{},
		AddFlag:             false,
		UpdateFlag:          false,
		DeleteFlag:          true,
		IsLocalFlag:         true,
		AliasNameFlag:       "",
		ActualAliasName:     "",
		AliasLocationFlag:   "",
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return testAliasName

		},
	}

	setNameForAlias(&args)

	if args.ActualAliasName != testAliasName {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasName, args.ActualAliasName)
	}
}

func Test_GivenNoArgsAndNoNameFlagAndAddFlag_ExpectPromptForName_setNameForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:                []string{},
		AddFlag:             true,
		UpdateFlag:          false,
		DeleteFlag:          false,
		IsLocalFlag:         true,
		AliasNameFlag:       "",
		ActualAliasName:     "",
		AliasLocationFlag:   "",
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return testAliasName
		},
	}

	setNameForAlias(&args)

	if args.ActualAliasName != testAliasName {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasName, args.ActualAliasName)
	}
}
func Test_GivenNoArgsAndNoNameFlagAndUpdateFlag_ExpectPromptForName_setNameForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"
	addAlias(testAliasName, testAliasLocation, false)

	args := types.AliasCmdArgs{
		Args:                []string{},
		AddFlag:             false,
		UpdateFlag:          true,
		DeleteFlag:          false,
		IsLocalFlag:         true,
		AliasNameFlag:       "",
		ActualAliasName:     "",
		AliasLocationFlag:   "",
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return testAliasName
		},
	}

	setNameForAlias(&args)

	if args.ActualAliasName != testAliasName {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasName, args.ActualAliasName)
	}
}

func Test_GivenNoArgsAndNoNameFlagAndDeleteFlag_ExpectPromptForNameTwiceWhenNonExistentNameProvided_setNameForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"
	addAlias(testAliasName, testAliasLocation, false)
	args := types.AliasCmdArgs{
		Args:                []string{},
		AddFlag:             false,
		UpdateFlag:          false,
		DeleteFlag:          true,
		IsLocalFlag:         true,
		AliasNameFlag:       "",
		ActualAliasName:     "",
		AliasLocationFlag:   "",
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			if utils.TestCounter == 0 {
				utils.TestCounter++
				return "doesnt-exist"
			}
			return testAliasName
		},
	}

	setNameForAlias(&args)

	if args.ActualAliasName != testAliasName {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasName, args.ActualAliasName)
	}
}

func Test_GivenInvalidFlags_ReturnFalse_isValidFlags(t *testing.T) {
	threeTrue := types.AliasCmdArgs{
		AddFlag:    true,
		UpdateFlag: true,
		DeleteFlag: true,
	}

	isValid := isValidFlags(&threeTrue)
	if isValid {
		t.Fatalf("Expected flags to be invalid, but they were valid.")
	}

	twoTrue := types.AliasCmdArgs{
		AddFlag:    false,
		UpdateFlag: true,
		DeleteFlag: true,
	}

	isValid2 := isValidFlags(&twoTrue)
	if isValid2 {
		t.Fatalf("Expected flags to be invalid, but they were valid.")
	}

	noneTrue := types.AliasCmdArgs{
		AddFlag:    false,
		UpdateFlag: false,
		DeleteFlag: false,
	}

	isValid3 := isValidFlags(&noneTrue)
	if isValid3 {
		t.Fatalf("Expected flags to be invalid, but they were valid.")
	}
}

func Test_GivenValidFlags_ReturnTrue_isValidFlags(t *testing.T) {

	delTrue := types.AliasCmdArgs{
		AddFlag:    false,
		UpdateFlag: false,
		DeleteFlag: true,
	}

	isValid := isValidFlags(&delTrue)
	if !isValid {
		t.Fatalf("Expected flags to be invalid, but they were valid.")
	}

	twoTrue := types.AliasCmdArgs{
		AddFlag:    false,
		UpdateFlag: true,
		DeleteFlag: false,
	}

	isValid2 := isValidFlags(&twoTrue)
	if !isValid2 {
		t.Fatalf("Expected flags to be invalid, but they were valid.")
	}

	addTrue := types.AliasCmdArgs{
		AddFlag:    true,
		UpdateFlag: false,
		DeleteFlag: false,
	}

	isValid3 := isValidFlags(&addTrue)
	if !isValid3 {
		t.Fatalf("Expected flags to be invalid, but they were valid.")
	}
}

func Test_GivenNoArgsAndNoLocationFlag_ExpectPromptForLocation_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{},
		IsLocalFlag:       true,
		AliasNameFlag:     "",
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return testAliasLocation
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenONEArgAndNoLocationFlag_ExpectPromptForLocation_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{testAliasLocation},
		IsLocalFlag:       true,
		AliasNameFlag:     "",
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return testAliasLocation
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenTWOArgsAndNoLocationFlag_ExpectActualLocationToBeSecondArg_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{testAliasName, testAliasLocation},
		IsLocalFlag:       true,
		AliasNameFlag:     "",
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			// When returning the prompt, it will not be equal to testAliasLocation
			// This proves we are not prompting for the location
			return prompt
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenNoArgsAndWITHLocationFlag_ExpectActualLocationToBeNameFlag_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{},
		IsLocalFlag:       true,
		AliasNameFlag:     "",
		AliasLocationFlag: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			// When returning the prompt, it will not be equal to testAliasLocation
			// This proves we are not prompting for the location
			return prompt
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenNoArgsAndWITHDeleteFlag_ExpectActualLocationToBeBlank_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{},
		IsLocalFlag:       true,
		DeleteFlag:        true,
		AliasNameFlag:     "",
		AliasLocationFlag: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			// When returning the prompt, it will not be equal to testAliasLocation
			// This proves we are not prompting for the location
			return prompt
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != "" {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenONEArgAndNoLocationFlagAndWITHNameFlag_ExpectActualLocationToBeTheArg_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{testAliasLocation},
		IsLocalFlag:       true,
		AliasNameFlag:     testAliasName,
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			// When returning the prompt, it will not be equal to testAliasLocation
			// This proves we are not prompting for the location
			return prompt
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenZEROArgAndNoLocationFlagAndWITHNameFlag_ExpectActualLocationToBeTheArg_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{},
		IsLocalFlag:       true,
		AliasNameFlag:     testAliasName,
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			// When prompted, provide the testAliasLocation
			return testAliasLocation
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}
}

func Test_GivenTWOArgsAndNOTLocal_ExpectActualLocationToBeSecondArgAndToBePromptedAboutLocal_setTemplateLocationForAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{testAliasName, testAliasLocation},
		IsLocalFlag:       false,
		AliasNameFlag:     "",
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			// In this scenario we give TWO args and 'isLocal=false' to ensure that the user is ONLY asked if this is a local machine,
			// Since 'isLocal' is currently false, returning 'y' to the prompt will trigger isLocalFlag to be set to true.
			return "y"
		},
	}

	setTemplateLocationForAlias(&args)

	if args.ActualAliasLocation != testAliasLocation {
		t.Fatalf("Names did not match. Expected: %s, Actual: %s", testAliasLocation, args.ActualAliasLocation)
	}

	if !args.IsLocalFlag {
		t.Fatalf("Expected IsLocalFlag to be true but was not.")
	}
}

func Test_GivenAddFlag_And_GivenPreExistingAlias_ExpectToUpdate_aliasManager(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"
	addAlias(testAliasName, testAliasLocation, true)
	expectedPrompt := fmt.Sprintf("Are you sure you want to update the alias: %s?", testAliasName)
	var actualPrompt string
	args := types.AliasCmdArgs{
		AddFlag:             true,
		UpdateFlag:          false,
		DeleteFlag:          false,
		IsLocalFlag:         false,
		ActualAliasName:     testAliasName,
		ActualAliasLocation: testAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
			actualPrompt = s
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return "unused"
		},
	}

	aliasManager(&args)

	aliases := config.Global().Aliases

	testAlias := cast.ToStringMapString(aliases[testAliasName])

	if cast.ToBool(testAlias[config.Global().AliasesLocalIndicatorKey]) {
		t.Fatalf("Expected local indicator in stored alias to be true but was not.")
	}

	if actualPrompt != expectedPrompt {
		t.Fatalf("Actual prompt was not equal to expected prompt. Actual: %s | Expected: %s", actualPrompt, expectedPrompt)
	}

}

func Test_GivenAddFlag_ExpectToAdd_aliasManager(t *testing.T) {
	setupTests(t)
	// Test inputs
	expectedAliasName := "lets-make-sure-this-is-unique"
	expectedAliasLocation := "/hello"
	expectedPrompt := fmt.Sprintf("Are you sure you want to add the alias: %s?", expectedAliasName)

	var actualPrompt string
	args := types.AliasCmdArgs{
		AddFlag:             true,
		UpdateFlag:          false,
		DeleteFlag:          false,
		IsLocalFlag:         false,
		ActualAliasName:     expectedAliasName,
		ActualAliasLocation: expectedAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
			actualPrompt = s
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return "unused"
		},
	}

	aliasManager(&args)

	aliases := config.Global().Aliases

	actualAlias := cast.ToStringMapString(aliases[expectedAliasName])

	actualAliasLocation := cast.ToString(actualAlias[config.Global().AliasesUrlKey])

	if len(actualAlias) == 0 {
		t.Fatalf("Something went wrong. Perhaps the alias didn't get created?")
	}

	if cast.ToBool(actualAlias[config.Global().AliasesLocalIndicatorKey]) {
		t.Fatalf("Expected local indicator in stored alias to be true but was not.")
	}

	if actualAliasLocation != expectedAliasLocation {
		t.Fatalf("Expected url to be %s but was actually %s.", expectedAliasLocation, actualAliasLocation)
	}

	if actualPrompt != expectedPrompt {
		t.Fatalf("Actual prompt was not equal to expected prompt. Actual: %s | Expected: %s", actualPrompt, expectedPrompt)
	}
}

func Test_GivenUpdateFlag_ExpectToUpdate_aliasManager(t *testing.T) {
	setupTests(t)
	// Test inputs
	expectedAliasName := "lets-make-sure-this-is-unique"
	expectedAliasLocation := "/hello"
	expectedPrompt := fmt.Sprintf("Are you sure you want to update the alias: %s?", expectedAliasName)

	var actualPrompt string
	args := types.AliasCmdArgs{
		AddFlag:             false,
		UpdateFlag:          true,
		DeleteFlag:          false,
		IsLocalFlag:         false,
		ActualAliasName:     expectedAliasName,
		ActualAliasLocation: expectedAliasLocation,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
			actualPrompt = s
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return "unused"
		},
	}

	aliasManager(&args)

	aliases := config.Global().Aliases

	actualAlias := cast.ToStringMapString(aliases[expectedAliasName])

	actualAliasLocation := cast.ToString(actualAlias[config.Global().AliasesUrlKey])

	if len(actualAlias) == 0 {
		t.Fatalf("Something went wrong. Perhaps the alias didn't get created?")
	}

	if cast.ToBool(actualAlias[config.Global().AliasesLocalIndicatorKey]) {
		t.Fatalf("Expected local indicator in stored alias to be true but was not.")
	}

	if actualAliasLocation != expectedAliasLocation {
		t.Fatalf("Expected url to be %s but was actually %s.", expectedAliasLocation, actualAliasLocation)
	}

	if actualPrompt != expectedPrompt {
		t.Fatalf("Actual prompt was not equal to expected prompt. Actual: %s | Expected: %s", actualPrompt, expectedPrompt)
	}
}

func Test_GivenDeleteFlag_And_GivenPreExistingAlias_ExpectToDelete_aliasManager(t *testing.T) {
	setupTests(t)
	// Test inputs
	testAliasName := "test-alias"
	testAliasLocation := "/hello"
	addAlias(testAliasName, testAliasLocation, true)
	expectedPrompt := fmt.Sprintf("Are you sure you want to delete the alias: %s?", testAliasName)
	var actualPrompt string
	args := types.AliasCmdArgs{
		AddFlag:         false,
		UpdateFlag:      false,
		DeleteFlag:      true,
		IsLocalFlag:     false,
		ActualAliasName: testAliasName,
		ConfirmFunction: func(s string) {
			log.Infoln(s)
			actualPrompt = s
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return "unused"
		},
	}

	aliasManager(&args)

	aliases := config.Global().Aliases

	testAlias := cast.ToStringMapString(aliases[testAliasName])

	if len(testAlias) != 0 {
		t.Fatalf("Expected alias to be deleted but was: %s.", testAlias)
	}

	if actualPrompt != expectedPrompt {
		t.Fatalf("Actual prompt was not equal to expected prompt. Actual: %s | Expected: %s", actualPrompt, expectedPrompt)
	}

}

func Test_GivenAddFlag_ExpectToAdd_processAlias(t *testing.T) {
	setupTests(t)
	// Test inputs
	expectedAliasName := "test-alias"
	expectedAliasLocation := "/hello"

	args := types.AliasCmdArgs{
		Args:              []string{expectedAliasName, expectedAliasLocation},
		AddFlag:           true,
		UpdateFlag:        false,
		DeleteFlag:        false,
		IsLocalFlag:       true,
		AliasNameFlag:     "",
		AliasLocationFlag: "",
		ConfirmFunction: func(s string) {
			log.Infoln(s)
		},
		StringInputReader: func(prompt string) string {
			log.Infoln(prompt)
			return "unused"
		},
	}

	processAlias(&args)

	aliases := config.Global().Aliases

	actualAlias := cast.ToStringMapString(aliases[expectedAliasName])

	actualAliasLocation := cast.ToString(actualAlias[config.Global().AliasesUrlKey])

	if len(actualAlias) == 0 {
		t.Fatalf("Something went wrong. Perhaps the alias didn't get created?")
	}

	if !cast.ToBool(actualAlias[config.Global().AliasesLocalIndicatorKey]) {
		t.Fatalf("Expected local indicator in stored alias to be true but was not.")
	}

	if actualAliasLocation != expectedAliasLocation {
		t.Fatalf("Expected url to be %s but was actually %s.", expectedAliasLocation, actualAliasLocation)
	}

}

/**
* Utils
 */
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

func setupTests(t *testing.T) {
	config.ConfigureLogger()
	t.Cleanup(func() {
		config.SetPropertyAndSave("Aliases", "")
	})
}
