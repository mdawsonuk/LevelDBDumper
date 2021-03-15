package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilewalkLevelDBDatabaseTopLevelDatabase(t *testing.T) {
	rootPath = "."
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
		t.Errorf("Database exists in current directory, actual: %d, expected: 0", len(searchResult))
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
