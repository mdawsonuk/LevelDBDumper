package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteDBNoErrorIfOutputEmpty(t *testing.T) {
	outputDir = ""
	writeDBInfo()
}

func TestDBMakeOutputDirectory(t *testing.T) {
	outputDir, _ = filepath.Abs("./test")
	outputType = "invalid"

	writeDBInfo()

	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		t.Error("Directory", outputDir, "should exist")
	}

	os.Remove(outputDir)
}

func TestWriteDBBatchCSV(t *testing.T) {
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

func TestWriteDBBatchJSON(t *testing.T) {
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

func TestWriteDBCSV(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	batch = false
	outputType = "csv"
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	csvFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.csv", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), "Test DB Path")

	var db = ParsedDB{path: "Test DB Path", modifiedTime: time.Now(), keys: []string{"Test Key"}, values: []string{"Test Value"}}
	parsedDatabases = append(parsedDatabases, db)
	writeDBInfo()
	_, err := os.Stat(csvFileName)
	if os.IsNotExist(err) {
		t.Error("File", csvFileName, "should exist")
	}

	os.Remove(csvFileName)
}

func TestWriteDBJSON(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	batch = false
	outputType = "json"
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	jsonFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.json", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), "Test DB Path")

	var db = ParsedDB{path: "Test DB Path", modifiedTime: time.Now(), keys: []string{"Test Key"}, values: []string{"Test Value"}}
	parsedDatabases = append(parsedDatabases, db)
	writeDBInfo()
	_, err := os.Stat(jsonFileName)
	if os.IsNotExist(err) {
		t.Error("File", jsonFileName, "should exist")
	}

	os.Remove(jsonFileName)
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

func TestOutputJSON(t *testing.T) {
	outputDir, _ = filepath.Abs(".")
	batch = false
	outputType = "json"
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	jsonFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.json", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), "Test DB Path")

	var db = ParsedDB{path: "Test DB Path", modifiedTime: time.Now(), keys: []string{"Test Key"}, values: []string{"Test Value"}}
	createJSONOutput(db)
	_, err := os.Stat(jsonFileName)
	if os.IsNotExist(err) {
		t.Error("File", jsonFileName, "should exist")
	}

	os.Remove(jsonFileName)
}
