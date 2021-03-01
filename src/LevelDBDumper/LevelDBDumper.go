package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/hashicorp/go-version"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const VERSION string = "3.0.0-alpha"

var (
	// Info message colour
	Info = Teal
	// Warn message colour
	Warn = Yellow
	// Fatal message colour
	Fatal = Red
)

var (
	// Black message colour
	Black = Colour("\033[1;30m%s\033[0m")
	// Red message colour
	Red = Colour("\033[1;31m%s\033[0m")
	// Green message colour
	Green = Colour("\033[1;32m%s\033[0m")
	// Yellow message colour
	Yellow = Colour("\033[1;33m%s\033[0m")
	// Purple message colour
	Purple = Colour("\033[1;34m%s\033[0m")
	// Magenta message colour
	Magenta = Colour("\033[1;35m%s\033[0m")
	// Teal message colour
	Teal = Colour("\033[1;36m%s\033[0m")
	// White message colour
	White = Colour("\033[1;37m%s\033[0m")
)

// Colour the string based on the string given
func Colour(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var (
	searchResult []string

	help       bool
	rootPath   string
	quiet      bool
	outputType string = "csv"
	outputDir  string
	outputFile string
	batch      bool
	noColour   bool
)

func getArgs(args []string) {
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			help = true
			break
		}
		if (args[i] == "-d" || args[i] == "--dir") && i+1 < len(args) {
			path, err := filepath.Abs(args[i+1])
			if err != nil {
				printLine(fmt.Sprintf("Unable to get absolute path of %s", path), Fatal)
			} else {
				rootPath = path
			}
		}
		if args[i] == "-q" || args[i] == "--quiet" {
			quiet = true
		}
		if (args[i] == "-t" || args[i] == "--outputType") && i+1 < len(args) {
			outputType = args[i+1]
		}
		if (args[i] == "-o" || args[i] == "--outputDir") && i+1 < len(args) {
			outputDir = args[i+1]
		}
		if (args[i] == "-f" || args[i] == "--outputFile") && i+1 < len(args) {
			outputFile = args[i+1]
		}
		if args[i] == "-b" || args[i] == "--batch" {
			batch = true
		}
		if args[i] == "--no-colour" || args[i] == "--no-color" {
			noColour = true
		}
	}
}

func printUsage() {
	fmt.Println("      h/help              Display this help message.")
	fmt.Println("      d/dir               Directory to recursively process. This is required.")
	fmt.Println("      q/quiet             Don't output all key/value pairs to console. Default will output all key/value pairs")
	fmt.Println("      t/outputType        Output type. Can be \"csv\" or \"json\"")
	fmt.Println("      o/outputDir         Directory to save all output results to. Required for any file output")
	fmt.Println("      f/outputFile        Filename to use when saving output. This will be appended with path and date")
	fmt.Println("      b/batch             Combine all output files into one file. Supported by \"csv\" and \"json\" file types")
	fmt.Println("      no-colour/no-color  Don't colourise output")
	fmt.Println()
	fmt.Println("Short options (single letter) are prefixed with a single dash. Long commands are prefixed with two dashes")
	fmt.Println()
	fmt.Println("Examples: LevelDBParser.exe -d \"C:\\Temp\\leveldb\"")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -o \"C:\\Temp\" -q")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" --no-colour --quiet")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" --no-colour -b --outputType json -outputFile Evidence.json")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -t csv -f LevelDB.csv -o Evidence -b --no-colour --quiet")
	fmt.Println()
}

func main() {
	dumpDBs(os.Args)
}

func dumpDBs(args []string) {

	fmt.Println()
	fmt.Println(fmt.Sprintf("LevelDB Dumper %s", VERSION))
	fmt.Println()
	fmt.Println("Author: Matt Dawson")
	fmt.Println()

	getArgs(args)

	if help {
		printUsage()
		os.Exit(0)
	}

	needsUpdate, latestVersion := checkUpdate()

	if !needsUpdate {
		printLine("You are using the latest version of LevelDB Dumper", Purple)
		fmt.Println()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		printLine(fmt.Sprintf("Found %d results so far", len(searchResult)), Info)
		printLine("Ctrl+C detected, quitting...", Fatal)
		os.Exit(0)
	}()

	fmt.Println("Command Line:", strings.Join(args[1:], " "))
	fmt.Println()

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
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

	var data = [][]string{}
	for iter.Next() {
		key := iter.Key()
		keyName := string(key[:])

		byteValue, err := db.Get([]byte(key), nil)
		if err != nil {
			printLine(fmt.Sprintf("Error reading Key: %s", keyName), Fatal)
			printLine(err.Error(), Fatal)
			return
		}
		value := string(byteValue)

		data = append(data, []string{keyName, value})
	}

	if !quiet {
		if len(data) > 0 {
			if !quiet {
				printLine(fmt.Sprintf("%-56vValue:", "Key:"), Info)
			}
			for _, keyValue := range data {
				escapedKey := removeControlChars(keyValue[0])   //fmt.Sprintf("%q", keyName)
				escapedValue := removeControlChars(keyValue[1]) //fmt.Sprintf("%q", value)
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

	if outputDir != "" {
		// When batching, timestamp column should use time.Now().Unix()
		if len(data) > 0 {
			switch outputType {
			case "csv":
				createCsvOutput(dbPath, data)
			case "json":
				createJSONOutput(dbPath, data)
				break
			}
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

func createCsvOutput(dbPath string, data [][]string) {
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	escapedPath := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(dbPath, "/", "_"), "\\", "_"), ":", "")
	csvFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.csv", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), escapedPath)
	file, err := os.Create(filepath.Join(outputDir, csvFileName))
	checkError(err)
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"Key", "Value"})

	for _, value := range data {
		err := csvWriter.Write(value)
		checkError(err)
		csvWriter.Flush()
	}
}

func createJSONOutput(dbPath string, data [][]string) {
	var jsonData = map[string]string{}

	for _, keyValue := range data {
		jsonData[keyValue[0]] = keyValue[1]
	}

	json, _ := json.MarshalIndent(jsonData, "", " ")
	fmt.Println(string(json))

	timeNow := time.Now()
	year, month, day := timeNow.Date()
	escapedPath := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(dbPath, "/", "_"), "\\", "_"), ":", "")
	jsonFileName := fmt.Sprintf("%v%v%v%v%v%v_%v_LevelDBDumper.json", year, int(month), day, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), escapedPath)
	file, err := os.Create(filepath.Join(outputDir, jsonFileName))
	checkError(err)
	defer file.Close()
	file.Write(json)
}

func checkUpdate() (bool, string) {
	resp, err := http.Get("https://api.github.com/repos/mdawsonuk/LevelDBDumper/releases/latest")
	checkError(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var results map[string]interface{}

	json.Unmarshal(body, &results)

	tag := fmt.Sprintf("%s", results["tag_name"])[1:]

	currentVersion, _ := version.NewVersion(VERSION)
	latestVersion, _ := version.NewVersion(tag)

	return currentVersion.LessThan(latestVersion), tag
}

func removeControlChars(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(Fatal(err))
	}
}

type format func(...interface{}) string

func printLine(contents string, fn format) {
	if noColour {
		fmt.Println(contents)
	} else {
		fmt.Println(fn(contents))
	}
}
