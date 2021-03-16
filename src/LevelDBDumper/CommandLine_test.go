package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintUsage(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printUsage()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if !strings.Contains(out, "h/help") {
		t.Errorf("Output does not contain string \"h/help\", output was instead:\n%s", out)
	}
	if !strings.Contains(out, "Display this help message") {
		t.Errorf("Output does not contain string \"Display this help message\", output was instead:\n%s", out)
	}
}

func TestArgsRootPathShort(t *testing.T) {
	args := []string{"-d", "."}
	getArgs(args)
	path, _ := os.Getwd()
	if rootPath != path {
		t.Errorf("rootPath was incorrect, actual: %s, expected: %s", rootPath, path)
	}
	rootPath = ""
}

func TestArgsRootPathLong(t *testing.T) {
	args := []string{"--dir", "."}
	getArgs(args)
	path, _ := os.Getwd()
	if rootPath != path {
		t.Errorf("rootPath was incorrect, actual: %s, expected: %s", rootPath, path)
	}
	rootPath = ""
}

func TestArgsQuietShort(t *testing.T) {
	args := []string{"-q"}
	getArgs(args)
	if quiet != true {
		t.Error("quiet was incorrect, actual: false, expected: true")
	}
	quiet = false
}

func TestArgsQuietLong(t *testing.T) {
	args := []string{"--quiet"}
	getArgs(args)
	if quiet != true {
		t.Error("quiet was incorrect, actual: false, expected: true")
	}
	quiet = false
}

func TestArgsOutputTypeShort(t *testing.T) {
	args := []string{"-t", "json"}
	getArgs(args)
	if outputType != "json" {
		t.Errorf("quiet was incorrect, actual: %s, expected: json", outputType)
	}
	outputType = ""
}

func TestArgsOutputTypeLong(t *testing.T) {
	args := []string{"--outputType", "json"}
	getArgs(args)
	if outputType != "json" {
		t.Errorf("outputType was incorrect, actual: %s, expected: json", outputType)
	}
	outputType = ""
}

func TestArgsOutputDirShort(t *testing.T) {
	args := []string{"-o", "test"}
	getArgs(args)
	if outputDir != "test" {
		t.Errorf("outputDir was incorrect, actual: %s, expected: test", outputDir)
	}
	outputDir = ""
}

func TestArgsOutputDirLong(t *testing.T) {
	args := []string{"--outputDir", "test"}
	getArgs(args)
	if outputDir != "test" {
		t.Errorf("outputDir was incorrect, actual: %s, expected: test", outputDir)
	}
	outputDir = ""
}

func TestArgsOutputFileShort(t *testing.T) {
	args := []string{"-f", "test"}
	getArgs(args)
	if outputFile != "test" {
		t.Errorf("outputFile was incorrect, actual: %s, expected: test", outputFile)
	}
	outputFile = ""
}

func TestArgsOutputFileLong(t *testing.T) {
	args := []string{"--outputFile", "test"}
	getArgs(args)
	if outputFile != "test" {
		t.Errorf("outputFile was incorrect, actual: %s, expected: test", outputFile)
	}
	outputFile = ""
}

func TestArgsNoHeader(t *testing.T) {
	args := []string{"--no-header"}
	getArgs(args)
	if noHeader != true {
		t.Errorf("noHeader was incorrect, actual: %t, expected: true", noHeader)
	}
}

func TestArgsUpdateShort(t *testing.T) {
	args := []string{"-u"}
	getArgs(args)
	if checkForUpdate != true {
		t.Errorf("checkForUpdate was incorrect, actual: %t, expected: true", checkForUpdate)
	}
	checkForUpdate = false
}

func TestArgsUpdateLong(t *testing.T) {
	args := []string{"--check-update"}
	getArgs(args)
	if checkForUpdate != true {
		t.Errorf("checkForUpdate was incorrect, actual: %t, expected: true", checkForUpdate)
	}
	checkForUpdate = false
}

func TestArgsHelpShort(t *testing.T) {
	args := []string{"-h"}
	getArgs(args)
	if help != true {
		t.Errorf("help was incorrect, actual: %t, expected: true", help)
	}
	help = false
}

func TestArgsHelpLong(t *testing.T) {
	args := []string{"--help"}
	getArgs(args)
	if help != true {
		t.Errorf("help was incorrect, actual: %t, expected: true", help)
	}
	help = false
}

func TestArgsCleanShort(t *testing.T) {
	args := []string{"-c"}
	getArgs(args)
	if cleanOutput != true {
		t.Errorf("cleanOutput was incorrect, actual: %t, expected: true", cleanOutput)
	}
	cleanOutput = false
}

func TestArgsCleanLong(t *testing.T) {
	args := []string{"--clean-output"}
	getArgs(args)
	if cleanOutput != true {
		t.Errorf("cleanOutput was incorrect, actual: %t, expected: true", cleanOutput)
	}
	cleanOutput = false
}

func TestArgsTimezoneShort(t *testing.T) {
	args := []string{"-z", "America/New_York"}
	getArgs(args)
	if timezone != "America/New_York" {
		t.Errorf("timezone was incorrect, actual: %s, expected: America/New_York", timezone)
	}
	timezone = ""
}

func TestArgsTimezoneLong(t *testing.T) {
	args := []string{"--timezone", "Europe/Berlin"}
	getArgs(args)
	if timezone != "Europe/Berlin" {
		t.Errorf("timezone was incorrect, actual: %s, expected: Europe/Berlin", timezone)
	}
	timezone = ""
}

func TestArgsBatchShort(t *testing.T) {
	args := []string{"-b"}
	getArgs(args)
	if batch != true {
		t.Errorf("batch was incorrect, actual: %t, expected: true", batch)
	}
	batch = false
}

func TestArgsBatchLong(t *testing.T) {
	args := []string{"--batch"}
	getArgs(args)
	if batch != true {
		t.Errorf("batch was incorrect, actual: %t, expected: true", batch)
	}
	batch = false
}
