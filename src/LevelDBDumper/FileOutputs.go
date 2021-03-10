package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// JSONDB holds all of the JSON data for the array of databases
type JSONDB struct {
	ModifiedTimestamp int64             `json:"modified_timestamp`
	Path              string            `json:"path"`
	Data              map[string]string `json:"data"`
}

func writeDBInfo() {
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
				// When batching, timestamp column should use time.Now().Unix()
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
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"Key", "Value"})

	for index := range db.keys {
		err := csvWriter.Write([]string{db.keys[index], db.values[index]})
		checkError(err)
		csvWriter.Flush()
	}
}

func createBatchCsvOutput() {
	csvFileName := "LevelDBDumper.csv"
	file, err := os.Create(filepath.Join(outputDir, csvFileName))
	checkError(err)
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"Timestamp (Last Modified)", "Key", "Value", "Path"})

	for _, db := range parsedDatabases {
		for index := range db.keys {
			err := csvWriter.Write([]string{strconv.FormatInt(db.modifiedTime, 10), db.keys[index], db.values[index], db.path})
			checkError(err)
			csvWriter.Flush()
		}
	}
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
	defer file.Close()
	file.Write(json)
}

func createBatchJSONOutput() {
	var databases = []JSONDB{}

	for _, db := range parsedDatabases {
		var data = map[string]string{}
		for index := range db.keys {
			data[db.keys[index]] = db.values[index]
		}
		databases = append(databases, JSONDB{ModifiedTimestamp: db.modifiedTime, Path: db.path, Data: data})
	}

	json, _ := json.MarshalIndent(databases, "", " ")

	jsonFileName := "LevelDBDumper.json"
	file, err := os.Create(filepath.Join(outputDir, jsonFileName))
	checkError(err)
	defer file.Close()
	file.Write(json)
}
