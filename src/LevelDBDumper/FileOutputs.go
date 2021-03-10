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

func writeDBInfo() {
	if batch {

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

func createJSONOutput(db ParsedDB) {
	var jsonData = map[string]string{}

	for index := range db.values {
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
