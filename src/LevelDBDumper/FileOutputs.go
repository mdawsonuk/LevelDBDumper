package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// See https://golang.org/src/time/format.go
const timeFormat = "2006-01-02T15:04:05 MST"

// JSONDB holds all of the JSON data for the array of databases
type JSONDB struct {
	ModifiedTimestamp string            `json:"modified_timestamp"`
	Path              string            `json:"path"`
	Data              map[string]string `json:"data"`
}

func writeDBInfo() {
	if outputDir == "" {
		return
	}
	err := os.MkdirAll(outputDir, os.ModePerm)
	checkError(err)
	if err != nil {
		return
	}
	if batch {
		switch outputType {
		case "csv":
			createBatchCsvOutput()
		case "json":
			createBatchJSONOutput()
			break
		}
	} else {
		for _, database := range parsedDatabases {
			if outputDir != "" {
				if len(database.keys) > 0 {
					switch outputType {
					case "csv":
						createCsvOutput(database)
					case "json":
						createJSONOutput(database)
						break
					}
				}
			}
		}
	}
}

func createCsvOutput(db ParsedDB) {
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	escapedPath := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(db.path, "/", "_"), "\\", "_"), ":", "")
	csvFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.csv", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), escapedPath)
	file, err := os.Create(filepath.Join(outputDir, csvFileName))
	checkError(err)
	if err != nil {
		return
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"Key", "Value"})

	for index := range db.keys {
		err := csvWriter.Write([]string{db.keys[index], db.values[index]})
		checkError(err)
		csvWriter.Flush()
	}
	file.Close()
}

func createBatchCsvOutput() {
	csvFileName := "LevelDBDumper.csv"
	file, err := os.Create(filepath.Join(outputDir, csvFileName))
	checkError(err)
	if err != nil {
		return
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"Timestamp (Last Modified)", "Key", "Value", "Path"})

	for _, db := range parsedDatabases {
		for index := range db.keys {
			err := csvWriter.Write([]string{db.modifiedTime.Format(timeFormat), db.keys[index], db.values[index], db.path})
			checkError(err)
			csvWriter.Flush()
		}
	}
	file.Close()
}

func createJSONOutput(db ParsedDB) {
	var jsonData = map[string]string{}

	for index := range db.keys {
		jsonData[db.keys[index]] = db.values[index]
	}

	json, _ := json.MarshalIndent(jsonData, "", " ")

	timeNow := time.Now()
	year, month, day := timeNow.Date()
	escapedPath := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(db.path, "/", "_"), "\\", "_"), ":", "")
	jsonFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.json", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), escapedPath)
	file, err := os.Create(filepath.Join(outputDir, jsonFileName))
	checkError(err)
	if err != nil {
		return
	}
	file.Write(json)
	file.Close()
}

func createBatchJSONOutput() {
	var databases = []JSONDB{}

	for _, db := range parsedDatabases {
		var data = map[string]string{}
		for index := range db.keys {
			data[db.keys[index]] = db.values[index]
		}
		databases = append(databases, JSONDB{ModifiedTimestamp: db.modifiedTime.Format(timeFormat), Path: db.path, Data: data})
	}

	json, _ := json.MarshalIndent(databases, "", " ")

	jsonFileName := "LevelDBDumper.json"
	file, err := os.Create(filepath.Join(outputDir, jsonFileName))
	checkError(err)
	if err != nil {
		return
	}
	file.Write(json)
	file.Close()
}
