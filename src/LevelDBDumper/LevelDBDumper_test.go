package main

import (
	"os"
	"testing"
)

func TestArgsRootPath(t *testing.T) {
	args := []string{"-d", ".", "-q", "--csv", "test"}
	rootPath, _, _, _ := getArgs(args)
	path, _ := os.Getwd()
	if rootPath != path {
		t.Errorf("rootPath was incorrect, actual: %s, expected: %s", rootPath, path)
	}
}

func TestArgsQuiet(t *testing.T) {
	args := []string{"-d", ".", "-q"}
	_, quiet, _, _ := getArgs(args)
	if quiet != true {
		t.Error("quiet was incorrect, actual: false, expected: true")
	}
}

func TestArgsCSV(t *testing.T) {
	args := []string{"-d", ".", "--csv", "test"}
	_, _, csv, _ := getArgs(args)
	if csv != "test" {
		t.Errorf("quiet was incorrect, actual: %s, expected: test", csv)
	}
}

func TestArgsNoColour(t *testing.T) {
	args := []string{"-d", ".", "--no-colour"}
	_, _, _, noColour := getArgs(args)
	if noColour != true {
		t.Errorf("quiet was incorrect, actual: %t, expected: true", noColour)
	}
}

func TestRemoveControlChars(t *testing.T) {
	input := "\x41\x00\x42\x05\x43\x1F"

	output := removeControlChars(input)

	if output != "ABC" {
		t.Errorf("removeControlChars was incorrect, actual: %s, expected: ABC", output)
	}
}
