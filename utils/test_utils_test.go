package utils

import "testing"

func Test_TestCounterIncrements_When_Called(t *testing.T) {
	if TestCounter != 0 {
		t.Fatalf("Something Unexpected Happend. TestCounter is expected to be 0 at the start of this test, but is %d", TestCounter)
	}
	TestCounter++
	if TestCounter != 1 {
		t.Fatalf("Something Unexpected Happend. TestCounter is expected to be 1 at the start of this test, but is %d", TestCounter)
	}
}
