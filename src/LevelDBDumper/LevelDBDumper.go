package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// VERSION of LevelDB Dumper
const VERSION string = "3.0.0-alpha.3"

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
		color.Magenta.Println("You are using the latest version of LevelDB Dumper")
		fmt.Println()
		if checkForUpdate {
			os.Exit(0)
		}
	} else if checkForUpdate {
		color.Cyan.Println(fmt.Sprintf("Version %s is now available for LevelDB Dumper - please update!", latestVersion))
		fmt.Println()
		os.Exit(0)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		color.FgLightBlue.Println(fmt.Sprintf("Found %d results so far", len(searchResult)))
		color.Red.Println("Ctrl+C detected, quitting...")
		os.Exit(0)
	}()

	if rootPath == "" {
		printUsage()
		color.Red.Println("Missing -d argument")
		os.Exit(1)
	}

	switch strings.ToLower(outputType) {
	case
		"csv",
		"json":
		break
	default:
		color.Yellow.Println(fmt.Sprintf("%s is not a recognised output type. Defaulting to CSV", outputType))
		fmt.Println()
		outputType = "csv"
	}

	dbPresent, _ := fileExists(rootPath)

	if !dbPresent {
		color.Red.Println(fmt.Sprintf("The path %s doesn't exist", rootPath))
		fmt.Println()
		os.Exit(2)
	}

	testFile, err := os.Open(rootPath)
	if err != nil {
		color.Yellow.Println(fmt.Sprintf("Unable to open %s - make sure you haven't escaped the path with \\\"", rootPath))
		fmt.Println()
		os.Exit(2)
	}
	defer testFile.Close()

	if !isAdmin() {
		color.Red.Println("You should run LevelDB Dumper with root/Administrator privileges")
	} else {
		color.FgLightBlue.Println("Running LevelDB Dumper with root/Administrator privileges")
	}
	fmt.Println()

	start := time.Now()

	searchForDBs()
	readDBs()

	elapsed := time.Now().Sub(start)
	color.FgLightBlue.Println(fmt.Sprintf("Completed search in %v", elapsed))
	fmt.Println()

	if needsUpdate {
		color.Magenta.Println(fmt.Sprintf("Version %s is now available for LevelDB Dumper - please update!", latestVersion))
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

	color.FgLightBlue.Println(fmt.Sprintf("Searching for LevelDB databases from %s took %v", rootPath, elapsed))
	fmt.Println()

	if len(searchResult) > 0 {
		if len(searchResult) == 1 {
			color.Yellow.Println("1 LevelDB database found")
		} else {
			color.Yellow.Println(fmt.Sprintf("%d LevelDB databases found", len(searchResult)))
		}

	} else {
		color.Red.Println("0 LevelDB databases found")
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
			color.Yellow.Println(fmt.Sprintf("Access denied for %s", path))
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
					color.Magenta.Println(fmt.Sprintf("Found database at %s", absolute))
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
		fmt.Println(fmt.Sprintf("%s %s", color.FgWhite.Render("Opening DB at"), color.FgYellow.Render(dbPath)))
	}

	options := &opt.Options{
		ErrorIfMissing: true,
		ReadOnly:       true,
	}

	comparator := getComparator(dbPath)

	switch comparator {
	case "idb_cmp1":
		color.Red.Println("IndexedDB idb_cmp1 comparator not yet implemented, results will not be output")
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
		color.Red.Println(fmt.Sprintf("Could not open DB: %s", err.Error()))
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
			color.Yellow.Println("Defaulting to using UTC timezone")
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
			color.Red.Println(fmt.Sprintf("Error reading Key: %s", keyName))
			color.Red.Println(err.Error())
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
				color.FgLightBlue.Println(fmt.Sprintf("%-58vValue:", "Key:"))
			}
			for index := range database.keys {
				escapedKey := removeControlChars(database.keys[index])     //fmt.Sprintf("%q", keyName)
				escapedValue := removeControlChars(database.values[index]) //fmt.Sprintf("%q", value)
				if len(escapedValue) > 80 {
					if noColour {
						fmt.Printf("%-53v | "+escapedValue[:80]+"...\n", escapedKey)
					} else {
						fmt.Printf("%-64v | "+escapedValue[:80]+"...\n", color.Yellow.Render(escapedKey))
					}
				} else {
					if noColour {
						fmt.Printf("%-53v | "+escapedValue+"\n", escapedKey)
					} else {
						fmt.Printf("%-64v | "+escapedValue+"\n", color.Yellow.Render(escapedKey))
					}

				}
			}
		} else {
			color.Yellow.Println("Parsed database but no key/value pairs were found")
		}
	}

	iter.Release()
	err = iter.Error()
	checkError(err)
	if !quiet {
		fmt.Println()
	}

	elapsed := time.Now().Sub(start)
	color.FgLightBlue.Println(fmt.Sprintf("Dumping LevelDB database took %s", elapsed))
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
