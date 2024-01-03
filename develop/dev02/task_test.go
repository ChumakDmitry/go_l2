package main

import (
	"fmt"
	"testing"
)

func TestParseStr(t *testing.T) {
	passStr := "a4bc2d5e"
	expected := "aaaabccddddde"
	errStr := []string{
		"45",
		"",
	}

	//pass
	result, err := ParseStr(passStr)
	if err != nil {
		t.Error(err)
	}
	if expected != result {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, expected)
	}

	//fail
	for _, val := range errStr {
		result, err := ParseStr(val)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(result)
	}
}
