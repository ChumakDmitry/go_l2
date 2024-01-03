package main

import (
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {
	passingCase := time.Now()

	currTime, err := GetCurrentTime()
	if err != nil {
		t.Error(err)
	}

	//pass
	if currTime.Before(passingCase) {
		t.Errorf("%v != %v", passingCase, currTime)
	}
}
