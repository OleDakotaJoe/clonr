package core

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"testing"
)

func Test_GivenScript_ExpectString_RunScriptAndReturnString(t *testing.T) {
	expectedResult := "some-string"
	testName := config.Global().GlobalsKeyName
	testVar := "globals-var"
	sampleGlobalsMap := types.ClonrVarMap{testVar: expectedResult}

	templateArgument := fmt.Sprintf("%s[%s]", testName, testVar)

	script := fmt.Sprintf("clonrResult = getClonrVar(\"%s\")", templateArgument)

	settings := types.FileProcessorSettings{
		GlobalsVarMap: sampleGlobalsMap,
	}

	actualResult, err := RunScriptAndReturnString(script, &settings)
	if err != nil {
		t.Fatalf("Something went wrong. Err: %s", err)
	}
	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func Test_GivenScript_ExpectBool_RunScriptAndReturnTrue(t *testing.T) {
	givenResult := "true"
	testName := config.Global().GlobalsKeyName
	testVar := "globals-var"
	sampleGlobalsMap := types.ClonrVarMap{testVar: givenResult}

	templateArgument := fmt.Sprintf("%s[%s]", testName, testVar)

	script := fmt.Sprintf("clonrResult = getClonrBool(\"%s\")", templateArgument)

	settings := types.FileProcessorSettings{
		GlobalsVarMap: sampleGlobalsMap,
	}

	actualResult, err := RunScriptAndReturnBool(script, &settings)
	if err != nil {
		t.Fatalf("Something went wrong. Err: %s", err)
	}
	if !actualResult {
		t.Fatalf("Expected true but got false")
	}
}

func Test_GivenScript_ExpectBool_RunScriptAndReturnFalse(t *testing.T) {
	givenResult := "false"
	testName := config.Global().GlobalsKeyName
	testVar := "globals-var"
	sampleGlobalsMap := types.ClonrVarMap{testVar: givenResult}

	templateArgument := fmt.Sprintf("%s[%s]", testName, testVar)

	script := fmt.Sprintf("clonrResult = getClonrBool(\"%s\")", templateArgument)

	settings := types.FileProcessorSettings{
		GlobalsVarMap: sampleGlobalsMap,
	}

	actualResult, err := RunScriptAndReturnBool(script, &settings)
	if err != nil {
		t.Fatalf("Something went wrong. Err: %s", err)
	}
	if actualResult {
		t.Fatalf("Expected false but got true")
	}
}

func Test_GivenValidInput_WhenGettingGlobalVariable_ReturnCorrectValue_getClonrVar(t *testing.T) {

	expectedResult := "globals-test"
	testName := config.Global().GlobalsKeyName
	testVar := "globals-var"
	sampleGlobalsMap := types.ClonrVarMap{testVar: expectedResult}

	templateArgument := fmt.Sprintf("%s[%s]", testName, testVar)
	templateDTO := types.ClonrVarDTO{
		Args:          []string{templateArgument},
		GlobalsVarMap: sampleGlobalsMap,
	}

	actualResult := resolveClonrVariable(&templateDTO)
	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func Test_GivenValidInput_WhenGettingTemplateVariable_ReturnCorrectValue_getClonrVar(t *testing.T) {
	testFilePath := "/test/file/path"
	testName := "test-file.txt"
	testVar := "template-var"
	expectedResult := "template-test"
	fullyQualifiedPathKey := fmt.Sprintf("%s/%s", testFilePath, testName)
	sampleMainTemplateMap := types.FileMap{fullyQualifiedPathKey: types.ClonrVarMap{testVar: expectedResult}}
	templateArgument := fmt.Sprintf("%s[%s]", testName, testVar)

	templateDTO := types.ClonrVarDTO{
		Args:            []string{templateArgument},
		MainTemplateMap: sampleMainTemplateMap,
		GlobalsVarMap:   nil,
		ConfigFilePath:  testFilePath,
	}

	actualResult := resolveClonrVariable(&templateDTO)
	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func Test_GivenValidScriptStringWithTags_ReturnScriptWithTags_ExtractScriptWithTags(t *testing.T) {
	expectedResult := config.Global().ConditionalExprPrefix + "some strings in the middle" + config.Global().ConditionalExprSuffix
	testScript := "some stuff before" + expectedResult + "some stuff after"

	actualResult := ExtractScriptWithTags(testScript)
	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func Test_GivenValidScriptStringWithTags_ReturnScriptWithoutTags_RemoveTagsFromScript(t *testing.T) {
	expectedResult := "test-script"
	testScript := config.Global().ConditionalExprPrefix + expectedResult + config.Global().ConditionalExprSuffix

	actualResult := RemoveTagsFromScript(testScript)
	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
