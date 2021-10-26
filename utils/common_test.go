package utils

import (
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
