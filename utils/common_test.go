package utils

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_GivenTwoMaps_MergeThem_ThenReturnOne_MergeStringMapString(t *testing.T) {
	mapOne := map[string]interface{}{"hello": "hello"}
	mapTwo := map[string]interface{}{"world": "world"}
	expectedResult := map[string]interface{}{"hello": "hello", "world": "world"}

	actualResult := MergeStringMaps(mapOne, mapTwo)

	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Fatal("Maps are not the same\n")
	}
}

func Test_GivenInput_ReturnInput_StringInputReader(t *testing.T) {
	expectedOutput := "some-input"
	prompt := "some-prompt"

	content := []byte(expectedOutput)
	tmpfile, err := ioutil.TempFile("", "example")
	ExitIfError(err)

	defer os.Remove(tmpfile.Name()) // clean up

	_, wErr := tmpfile.Write(content)
	ExitIfError(wErr)

	_, sErr := tmpfile.Seek(0, 0)
	ExitIfError(sErr)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	actualOutput := StringInputReader(prompt)

	cErr := tmpfile.Close()
	ExitIfError(cErr)

	if expectedOutput != actualOutput {
		t.Fatalf("Input was not returned, expected: %s, got: %s", expectedOutput, actualOutput)
	}
}
