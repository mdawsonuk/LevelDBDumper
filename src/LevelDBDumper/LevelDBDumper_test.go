package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFilewalkLevelDBDatabaseTopLevelDatabase(t *testing.T) {
	searchResult = []string{}

	currentFile, _ := os.Create("CURRENT")
	currentFile.Close()
	manifestFile, _ := os.Create("MANIFEST-0000")
	manifestFile.Close()

	path, _ := os.Getwd()
	err := filepath.Walk(path, findFile)
	if err != nil {
		t.Errorf("filewalk produced error")
	}
	if len(searchResult) != 1 {
		t.Errorf("Database exists in current directory, acual: %d, expected: 0", len(searchResult))
	}

	os.Remove("CURRENT")
	os.Remove("MANIFEST-0000")
}

func TestFilewalkLevelDBDatabaseNoDatabases(t *testing.T) {
	searchForDBs()
	if len(searchResult) > 0 {
		t.Errorf("Database exists in current directory, actual: %d, expected: 0", len(searchResult))
	}
}

func TestSearchForDBsTopLevel(t *testing.T) {
	currentFile, _ := os.Create("CURRENT")
	currentFile.Close()
	manifestFile, _ := os.Create("MANIFEST-0000")
	manifestFile.Close()

	searchForDBs()
	if len(searchResult) != 1 {
		t.Errorf("Database should exist in current directory, actual: %d, expected: 1", len(searchResult))
	}

	os.Remove("CURRENT")
	os.Remove("MANIFEST-0000")
}

func TestSearchForDBsTopLevelAndSubDir(t *testing.T) {
	os.MkdirAll("Discord/Local Storage/leveldb", os.ModePerm)

	currentFile, _ := os.Create("CURRENT")
	currentFile.Close()
	manifestFile, _ := os.Create("MANIFEST-0000")
	manifestFile.Close()
	currentFile, _ = os.Create("Discord/Local Storage/leveldb/CURRENT")
	currentFile.Close()
	manifestFile, _ = os.Create("Discord/Local Storage/leveldb/MANIFEST-0000")
	manifestFile.Close()

	searchForDBs()
	if len(searchResult) != 2 {
		t.Errorf("Databases should exist in current directory, actual: %d, expected: 2", len(searchResult))
	}

	os.Remove("CURRENT")
	os.Remove("MANIFEST-0000")
	os.RemoveAll("Discord")
}

func TestRemoveControlChars(t *testing.T) {
	input := "\x41\x00\x42\x05\x43\x1F"

	output := removeControlChars(input)

	if output != "ABC" {
		t.Errorf("removeControlChars was incorrect, actual: %s, expected: ABC", output)
	}
}

func TestOutputCSV(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	csvFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.csv", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), "Test DB Path")
	createCsvOutput("Test DB Path", [][]string{
		{"Test Key", "Test Value"},
	})
	_, err := os.Stat(csvFileName)
	if os.IsNotExist(err) {
		t.Error("File", csvFileName, "should exist")
	}

	os.Remove(csvFileName)
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
