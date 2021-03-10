package main

import (
	"os"
	"testing"
)

func TestRemoveControlChars(t *testing.T) {
	input := "\x41\x00\x42\x05\x43\x1F"

	output := removeControlChars(input)

	if output != "ABC" {
		t.Errorf("removeControlChars was incorrect, actual: %s, expected: ABC", output)
	}
}

func TestFileExists(t *testing.T) {
	file, _ := os.Create("test.txt")

	exists, _ := fileExists("test.txt")

	if !exists {
		t.Errorf("fileExists was incorrect, actual: %t, expected: true", exists)
	}
	file.Close()
	os.Remove("test.txt")
}

func TestFileExistsNotExist(t *testing.T) {
	exists, _ := fileExists("test.txt")

	if exists {
		t.Errorf("fileExists was incorrect, actual: %t, expected: false", exists)
	}
}
