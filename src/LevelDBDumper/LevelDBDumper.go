package main

// With thanks to https://github.com/harshvsingh8/leveldb-reader for the bulk of the LevelDB Go code

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func main() {

	fmt.Println()
	fmt.Println("LevelDB Dumper 2.0.0.0")
	fmt.Println()
	fmt.Println("Author: Matt Dawson")
	fmt.Println()

	getArgs := func() (string, bool, string) {
		dbPath := ""
		quiet := false
		csvPath := ""

		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-d" && i+1 < len(os.Args) {
				dbPath = os.Args[i+1]
			}
			if os.Args[i] == "-q" {
				quiet = true
			}
			if os.Args[i] == "--csv" && i+1 < len(os.Args) {
				csvPath = os.Args[i+1]
			}
		}
		return dbPath, quiet, csvPath
	}

	printUsage := func() {
		fmt.Println("        d               Directory to recursively process. This is required.")
		fmt.Println("        q               Don't output all key/value pairs to console. Default is false")
		fmt.Println("        csv             Directory to save CSV formatted results to. Be sure to include the full path in double quotes")
		fmt.Println()
		fmt.Println("Examples: LevelDBParser.exe -f \"C:\\Temp\\leveldb\\\"")
		fmt.Println("          LevelDBParser.exe -f \"C:\\Temp\\leveldb\\\" --csv \"C:\\Temp\\\"")
		fmt.Println()
		fmt.Println("          Short options (single letter) are prefixed with a single dash. Long commands are prefixed with two dashes")
		fmt.Println()
	}

	fileExists := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}

	rootPath, quiet, csvPath := getArgs()

	if rootPath == "" {
		printUsage()
		fmt.Println(Fata("Missing -d argument"))
		fmt.Println()
		return
	}

	fmt.Println("Command Line:", strings.Join(os.Args[1:], " "))
	fmt.Println()

	if !checkForAdmin() {
		fmt.Println(Warn("Administrative privileges not found!"))
	}

	dbPresent, _ := fileExists(rootPath)

	if !dbPresent {
		fmt.Println(Fata("The DB path", rootPath, "doesn't exist"))
		printUsage()
		return
	}

	openDb(rootPath, quiet, csvPath)
}

func openDb(dbPath string, quiet bool, csvPath string) {
	fmt.Println(Info("Opening DB at ", dbPath))
	fmt.Println()

	options := &opt.Options{
		ErrorIfMissing: true,
	}

	start := time.Now()

	db, err := leveldb.OpenFile(dbPath, options)

	if err != nil {
		fmt.Println(Fata("Could not open DB at ", dbPath))
		return
	}
	defer db.Close()

	if csvPath != "" {
		csvWriter := csv.NewWriter(os.Stdout)
		csvWriter.Write([]string{"Key", "Value"})
		csvWriter.Flush()
	}

	iter := db.NewIterator(nil, nil)

	if !quiet {
		fmt.Println(Info(fmt.Sprintf("%-56vValue:", "Key:")))
	}

	for iter.Next() {
		key := iter.Key()
		keyName := string(key[:])

		byteValue, err := db.Get([]byte(key), nil)
		if err != nil {
			fmt.Println("Error reading Key: " + keyName)
			return
		}
		value := string(byteValue)

		escapedKey := removeControlChars(keyName)
		escapedValue := removeControlChars(value)

		if !quiet {
			if len(escapedValue) > 80 {
				fmt.Printf("%-64v | "+escapedValue[:80]+"...\n", Warn(escapedKey))
			} else {
				fmt.Printf("%-64v | "+escapedValue+"\n", Warn(escapedKey))
			}
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Println(Fata(err))
	}
	if !quiet {
		fmt.Println()
	}

	elapsed := time.Now().Sub(start)
	fmt.Println(Info("Dumping LevelDB database at ", dbPath, " took ", elapsed))
	fmt.Println()
}

func checkForAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func removeControlChars(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
}
