package main

import (
	"os"
	"testing"
)

func TestSolving(t *testing.T) {
	file, err := os.Open("data/test.data")
	defer file.Close()

	result, err := Solve(file, true)
	if err != nil {
		t.Error(err)
	}
	if result != 240 {
		t.Error(`Wrong answer.`)
	}
}

func TestSolving2(t *testing.T) {
	file, err := os.Open("data/test.data")
	defer file.Close()

	result, err := Solve(file, false)
	if err != nil {
		t.Error(err)
	}
	if result != 4455 {
		t.Error(`Wrong answer.`)
	}
}
