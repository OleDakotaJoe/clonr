package core

import (
	"github.com/oledakotajoe/clonr/config"
	"testing"
)

func Test_GivenValidScriptStringWithTags_ReturnScriptWithoutTags_RemoveTagsFromScript(t *testing.T) {
	expectedResult := "test-script"
	testScript := config.Global().ConditionalExprPrefix + expectedResult + config.Global().ConditionalExprSuffix

	actualResult := RemoveTagsFromScript(testScript)
	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
