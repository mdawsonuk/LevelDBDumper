package main

import (
	"os"
	"testing"
)

func TestArgsRootPathShort(t *testing.T) {
	args := []string{"-d", "."}
	getArgs(args)
	path, _ := os.Getwd()
	if rootPath != path {
		t.Errorf("rootPath was incorrect, actual: %s, expected: %s", rootPath, path)
	}
}

func TestArgsRootPathLong(t *testing.T) {
	args := []string{"-dir", "."}
	getArgs(args)
	path, _ := os.Getwd()
	if rootPath != path {
		t.Errorf("rootPath was incorrect, actual: %s, expected: %s", rootPath, path)
	}
}

func TestArgsQuietShort(t *testing.T) {
	args := []string{"-q"}
	getArgs(args)
	if quiet != true {
		t.Error("quiet was incorrect, actual: false, expected: true")
	}
}

func TestArgsQuietLong(t *testing.T) {
	args := []string{"--quiet"}
	getArgs(args)
	if quiet != true {
		t.Error("quiet was incorrect, actual: false, expected: true")
	}
}

func TestArgsOutputTypeShort(t *testing.T) {
	args := []string{"-t", "json"}
	getArgs(args)
	if outputType != "json" {
		t.Errorf("quiet was incorrect, actual: %s, expected: json", outputType)
	}
}

func TestArgsOutputTypeLong(t *testing.T) {
	args := []string{"-outputType", "json"}
	getArgs(args)
	if outputType != "json" {
		t.Errorf("quiet was incorrect, actual: %s, expected: json", outputType)
	}
}

func TestArgsOutputDirShort(t *testing.T) {
	args := []string{"-o", "test"}
	getArgs(args)
	if outputDir != "test" {
		t.Errorf("outputDir was incorrect, actual: %s, expected: test", outputDir)
	}
}

func TestArgsOutputDirLong(t *testing.T) {
	args := []string{"--outputDir", "test"}
	getArgs(args)
	if outputDir != "test" {
		t.Errorf("outputDir was incorrect, actual: %s, expected: test", outputDir)
	}
}

func TestArgsOutputFileShort(t *testing.T) {
	args := []string{"-f", "test"}
	getArgs(args)
	if outputFile != "test" {
		t.Errorf("outputFile was incorrect, actual: %s, expected: test", outputFile)
	}
}

func TestArgsOutputFileLong(t *testing.T) {
	args := []string{"--outputFile", "test"}
	getArgs(args)
	if outputFile != "test" {
		t.Errorf("outputFile was incorrect, actual: %s, expected: test", outputFile)
	}
}

func TestArgsNoColour(t *testing.T) {
	args := []string{"--no-colour"}
	getArgs(args)
	if noColour != true {
		t.Errorf("noColour was incorrect, actual: %t, expected: true", noColour)
	}
}

func TestArgsNoColor(t *testing.T) {
	args := []string{"--no-color"}
	getArgs(args)
	if noColour != true {
		t.Errorf("noColour was incorrect, actual: %t, expected: true", noColour)
	}
}
