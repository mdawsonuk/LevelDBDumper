package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unicode"

	"github.com/gookit/color"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func readDBs() {
	for _, v := range searchResult {
		openDb(v)
	}
}

func openDb(dbPath string) {

	fmt.Println(fmt.Sprintf("%s %s", color.FgWhite.Render("Opening DB at"), color.FgYellow.Render(dbPath)))
	fmt.Println()

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
	case "leveldb.BytewiseComparator":
		color.FgLightBlue.Println("Using leveldb.BytewiseComparator")
		break
	default:
		// We don't know this comparator, break out
		fmt.Println(fmt.Sprintf("%s %s", color.FgWhite.Render("Using unrecognised comparator:"), color.FgYellow.Render(comparator)))
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
	if err != nil {
		color.Red.Println(fmt.Sprintf("Could not load MANIFEST file: %s", err.Error()))
	}
	checkError(err)
	manifestPath := files[0]
	info, err := os.Stat(manifestPath)
	checkError(err)
	// Display the dates in UTC
	loc, _ := time.LoadLocation("UTC")
	if timezone != "" {
		// Display the dates in UTC
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			color.Yellow.Println("Unable to load", timezone, "defaulting to using UTC timezone")
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
				escapedKey := removeControlChars(database.keys[index])
				escapedValue := removeControlChars(database.values[index])
				if len(escapedValue) > 80 {
					fmt.Printf("%-64v | "+escapedValue[:80]+"...\n", color.Yellow.Render(escapedKey))
				} else {
					fmt.Printf("%-64v | "+escapedValue+"\n", color.Yellow.Render(escapedKey))
				}
			}
		} else {
			color.Yellow.Println("Parsed database but no key/value pairs were found")
			color.Yellow.Println("It is possible that key/value pairs were present, but have been deleted")
			color.Yellow.Println("LevelDB Dumper does not currently support retrieving deleted keys")
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

	f, err := os.Open(files[0])
	if !quiet {
		fmt.Println("Using", files[0], "to parse comparator")
	}
	contents := make([]byte, 32)
	// The string containing the comparator type is always 9 bytes in
	f.Seek(9, 0)
	f.Read(contents)
	f.Close()

	for i, b := range contents {
		// Read until we reach the first non-graphic byte, identifying the end of the comparator string
		if !unicode.IsGraphic(rune(b)) {
			return string(contents[:i])
		}
	}

	color.Red.Println("Unable to parse comparator from", string(contents))

	return "Unknown"
}
