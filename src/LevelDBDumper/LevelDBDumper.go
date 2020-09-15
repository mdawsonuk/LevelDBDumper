package main

// With thanks to https://github.com/harshvsingh8/leveldb-reader for the bulk of the LevelDB Go code

import (
		"fmt"
		"os"
		"strings"
		"unicode"
)

import "github.com/syndtr/goleveldb/leveldb"

func main() {
	
	printUsage := func() {
		fmt.Println("Usage: LevelDBDumper.exe path/to/leveldb")
	}

	fileExists := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil { return true, nil }
		if os.IsNotExist(err) { return false, nil }
		return true, err
	}

	if len(os.Args) == 1 {
		fmt.Println("LevelDB folder path is not supplied")
		printUsage()
		return
	}

	dbPath := os.Args[1]

	dbPresent, err := fileExists(dbPath)

	if !dbPresent {
		fmt.Println("The DB path: " + dbPath + " does not exist")
		printUsage()
		return
	}

	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("Could not open DB from:", dbPath)
		printUsage()
		return
	}

	iter := db.NewIterator(nil, nil)
	
	for iter.Next() {
		key := iter.Key()
		keyName := string(key[:])
		
		data, err := db.Get([]byte(key), nil)
		if err != nil {
			fmt.Println("Error reading Key: " + keyName)
			return
		}
		s := string(data)
		
		s = removeControlChars(s)
		
		fmt.Println(removeControlChars(keyName))
		fmt.Println(s)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Println(err)
	}
}

func removeControlChars(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
}