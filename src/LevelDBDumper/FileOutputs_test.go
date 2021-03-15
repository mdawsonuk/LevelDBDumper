package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTestDBMakeOutputDirectory(t *testing.T) {
	outputDir, _ = filepath.Abs("./test")
	outputType = "invalid"

	writeDBInfo()

	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		t.Error("Directory", outputDir, "should exist")
	}

	os.Remove(outputDir)
}

func TestTestDBBatchCSV(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	outputType = "csv"
	batch = true
	parsedDatabases = []ParsedDB{}

	writeDBInfo()

	_, err := os.Stat("LevelDBDumper.csv")
	if os.IsNotExist(err) {
		t.Error("File LevelDBDumper.csv should exist")
	}

	os.Remove("LevelDBDumper.csv")
}

func TestTestDBBatchJSON(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	outputType = "json"
	batch = true
	parsedDatabases = []ParsedDB{}

	writeDBInfo()

	_, err := os.Stat("LevelDBDumper.json")
	if os.IsNotExist(err) {
		t.Error("File LevelDBDumper.json should exist")
	}

	os.Remove("LevelDBDumper.json")
}

func TestOutputCSV(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	batch = false
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	csvFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.csv", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), "Test DB Path")

	var db = ParsedDB{path: "Test DB Path", modifiedTime: time.Now(), keys: []string{"Test Key"}, values: []string{"Test Value"}}
	createCsvOutput(db)
	_, err := os.Stat(csvFileName)
	if os.IsNotExist(err) {
		t.Error("File", csvFileName, "should exist")
	}

	os.Remove(csvFileName)
}
