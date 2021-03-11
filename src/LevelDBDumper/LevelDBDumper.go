package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// VERSION of LevelDB Dumper
const VERSION string = "3.0.0-alpha.2"

var (
	searchResult    []string
	parsedDatabases []ParsedDB
)

func main() {
	dumpDBs(os.Args)
}

func dumpDBs(args []string) {

	getArgs(args)

	if !noHeader {
		fmt.Println()
		fmt.Println(fmt.Sprintf("LevelDB Dumper %s", VERSION))
		fmt.Println()
		fmt.Println("Author: Matt Dawson")
		fmt.Println()
	}

	if help {
		printUsage()
		os.Exit(0)
	}

	fmt.Println("Command Line:", strings.Join(args[1:], " "))
	fmt.Println()

	needsUpdate, latestVersion := checkUpdate()

	if !needsUpdate {
		printLine("You are using the latest version of LevelDB Dumper", Purple)
		fmt.Println()
		if checkForUpdate {
			os.Exit(0)
		}
	} else if checkForUpdate {
		printLine(fmt.Sprintf("Version %s is now available for LevelDB Dumper - please update!", latestVersion), Purple)
		fmt.Println()
		os.Exit(0)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		printLine(fmt.Sprintf("Found %d results so far", len(searchResult)), Info)
		printLine("Ctrl+C detected, quitting...", Fatal)
		os.Exit(0)
	}()

	if rootPath == "" {
		printUsage()
		printLine("Missing -d argument", Fatal)
		os.Exit(1)
	}

	switch strings.ToLower(outputType) {
	case
		"csv",
		"json":
		break
	default:
		printLine(fmt.Sprintf("%s is not a recognised output type. Defaulting to CSV", outputType), Warn)
		fmt.Println()
		outputType = "csv"
	}

	dbPresent, _ := fileExists(rootPath)

	if !dbPresent {
		printLine(fmt.Sprintf("The path %s doesn't exist", rootPath), Fatal)
		fmt.Println()
		os.Exit(2)
	}

	testFile, err := os.Open(rootPath)
	if err != nil {
		printLine(fmt.Sprintf("Unable to open %s - make sure you haven't escaped the path with \\\"", rootPath), Warn)
		fmt.Println()
		os.Exit(2)
	}
	defer testFile.Close()

	if !isAdmin() {
		printLine("You should run LevelDB Dumper with root/Administrator privileges", Fatal)
	} else {
		printLine("Running LevelDB Dumper with root/Administrator privileges", Info)
	}
	fmt.Println()

	start := time.Now()

	searchForDBs()
	readDBs()

	elapsed := time.Now().Sub(start)
	printLine(fmt.Sprintf("Completed search in %v", elapsed), Info)
	fmt.Println()

	if needsUpdate {
		printLine(fmt.Sprintf("Version %s is now available for LevelDB Dumper - please update!", latestVersion), Purple)
		fmt.Println()
	}

	os.Exit(0)
}

func searchForDBs() {
	searchResult = []string{}

	start := time.Now()
	err := filepath.Walk(rootPath, findFile)
	if err != nil {
		return
	}
	elapsed := time.Now().Sub(start)

	if len(searchResult) > 0 {
		fmt.Println()
	}

	printLine(fmt.Sprintf("Searching for LevelDB databases from %s took %v", rootPath, elapsed), Info)
	fmt.Println()

	if len(searchResult) > 0 {
		printLine(fmt.Sprintf("%d LevelDB databases found", len(searchResult)), Warn)
	} else {
		printLine("0 LevelDB databases found", Fatal)
	}
	fmt.Println()
}

func readDBs() {
	for _, v := range searchResult {
		openDb(v)
	}
	// See FileOutputs.go
	writeDBInfo()
}

func findFile(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		if !quiet {
			printLine(fmt.Sprintf("Access denied for %s", path), Warn)
		}
		return nil
	}

	absolute, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if fileInfo.IsDir() {
		files, err := filepath.Glob(filepath.Join(absolute, "CURRENT"))
		checkError(err)
		if len(files) > 0 {
			files, err := filepath.Glob(filepath.Join(absolute, "MANIFEST-*"))
			checkError(err)
			if len(files) > 0 {
				searchResult = append(searchResult, absolute)
				if !quiet {
					printLine(fmt.Sprintf("Found database at %s", absolute), Purple)
				}
			}
		}
		return nil
	}

	return nil
}

func openDb(dbPath string) {

	if noColour {
		fmt.Println("Opening DB at", dbPath)
	} else {
		fmt.Println(Info("Opening DB at ", Warn(dbPath)))
	}

	options := &opt.Options{
		ErrorIfMissing: true,
		ReadOnly:       true,
	}

	comparator := getComparator(dbPath)

	switch comparator {
	case "idb_cmp1":
		printLine("IndexedDB idb_cmp1 comparator not yet implemented, results will not be output", Fatal)
		options.Comparer = idbCmp1{}
		fmt.Println()
		return
	default:
		// Just leave it, as default is leveldb.bitwisecomparator
		break
	}

	start := time.Now()

	db, err := leveldb.OpenFile(dbPath, options)

	if err != nil {
		printLine(fmt.Sprintf("Could not open DB: %s", err.Error()), Fatal)
		fmt.Println()
		return
	}
	fmt.Println()

	defer db.Close()

	iter := db.NewIterator(nil, nil)

	// TODO: If you get created time, either use directory or LOG file
	files, err := filepath.Glob(filepath.Join(dbPath, "MANIFEST-*"))
	checkError(err)
	manifestPath := files[0]
	info, err := os.Stat(manifestPath)
	checkError(err)
	// Display the dates in UTC
	loc, _ := time.LoadLocation("UTC")
	if timezone != "" {
		// Display the dates in UTC
		loc, err = time.LoadLocation(timezone)
		checkError(err)
		if err != nil {
			printLine("Defaulting to using UTC timezone", Warn)
			fmt.Println()
			loc, _ = time.LoadLocation("UTC")
		}
	}

	var database = ParsedDB{path: dbPath, modifiedTime: info.ModTime().In(loc), keys: []string{}, values: []string{}}

	for iter.Next() {
		key := iter.Key()
		keyName := string(key[:])
		if cleanOutput {
			keyName = removeControlChars(keyName)
		}

		byteValue, err := db.Get([]byte(key), nil)
		if err != nil {
			printLine(fmt.Sprintf("Error reading Key: %s", keyName), Fatal)
			printLine(err.Error(), Fatal)
			return
		}
		value := string(byteValue)
		if cleanOutput {
			value = removeControlChars(value)
		}

		database.keys = append(database.keys, keyName)
		database.values = append(database.values, value)
	}
	parsedDatabases = append(parsedDatabases, database)

	if !quiet {
		if len(database.keys) > 0 {
			if !quiet {
				printLine(fmt.Sprintf("%-56vValue:", "Key:"), Info)
			}
			for index := range database.keys {
				escapedKey := removeControlChars(database.keys[index])     //fmt.Sprintf("%q", keyName)
				escapedValue := removeControlChars(database.values[index]) //fmt.Sprintf("%q", value)
				if len(escapedValue) > 80 {
					if noColour {
						fmt.Printf("%-53v | "+escapedValue[:80]+"...\n", escapedKey)
					} else {
						fmt.Printf("%-64v | "+escapedValue[:80]+"...\n", Warn(escapedKey))
					}
				} else {
					if noColour {
						fmt.Printf("%-53v | "+escapedValue+"\n", escapedKey)
					} else {
						fmt.Printf("%-64v | "+escapedValue+"\n", Warn(escapedKey))
					}

				}
			}
		} else {
			printLine("Parsed database but no key/value pairs were found", Warn)
		}
	}

	iter.Release()
	err = iter.Error()
	checkError(err)
	if !quiet {
		fmt.Println()
	}

	elapsed := time.Now().Sub(start)
	printLine(fmt.Sprintf("Dumping LevelDB database took %s", elapsed), Info)
	fmt.Println()
}

func getComparator(dbPath string) string {
	files, err := filepath.Glob(filepath.Join(dbPath, "MANIFEST-*"))
	checkError(err)
	manifestPath := files[0]

	f, err := os.Open(manifestPath)
	contents := make([]byte, 32)
	// The string containing the comparator type is always 9 bytes in
	f.Seek(9, 0)
	f.Read(contents)
	f.Close()

	for i, b := range contents {
		// Read until we reach the 0x02 byte at the end of the comparator
		if b == 0x02 {
			return string(contents[:i])
		}
	}

	return "Unknown"
}
