package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

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
				fmt.Println(Fatal("Unable to get absolute path of ", path))
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
	fmt.Println("      t/outputType        Output type. Can be \"csv\", \"text\" or \"json\". JSON and text coming soon")
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
	fmt.Println("LevelDB Dumper 3.0.0-alpha")
	fmt.Println()
	fmt.Println("Author: Matt Dawson")
	fmt.Println()

	getArgs(args)

	if help {
		printUsage()
		os.Exit(0)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if noColour {
			fmt.Println("Found", len(searchResult), "results so far")
			fmt.Println("Ctrl+C detected, quitting...")
		} else {
			fmt.Println(Info("Found ", len(searchResult), " results so far"))
			fmt.Println(Fatal("Ctrl+C detected, quitting..."))
		}
		os.Exit(0)
	}()

	fmt.Println("Command Line:", strings.Join(args[1:], " "))
	fmt.Println()

	if rootPath == "" {
		printUsage()
		if noColour {
			fmt.Println("Missing -d argument")
		} else {
			fmt.Println(Fatal("Missing -d argument"))
		}
		os.Exit(1)
	}

	switch strings.ToLower(outputType) {
	case
		"csv",
		"text",
		"json":
		break
	default:
		if noColour {
			fmt.Println(outputType, "is not a recognised output type. Defaulting to CSV")
		} else {
			fmt.Println(Warn(outputType, " is not a recognised output type. Defaulting to CSV"))
		}
		fmt.Println()
		outputType = "csv"
	}

	dbPresent, _ := fileExists(rootPath)

	if !dbPresent {
		if noColour {
			fmt.Println("The path", rootPath, "doesn't exist")
		} else {
			fmt.Println(Fatal("The path ", rootPath, " doesn't exist"))
		}
		fmt.Println()
		os.Exit(2)
	}

	testFile, err := os.Open(rootPath)
	if err != nil {
		if noColour {
			fmt.Println("Unable to open", rootPath, "- make sure you haven't escaped the path with \\\"")
		} else {
			fmt.Println(Warn("Unable to open ", rootPath, " - make sure you haven't escaped the path with \\\""))
		}
		fmt.Println()
		os.Exit(2)
	}
	defer testFile.Close()

	searchForDBs()
	readDBs()
}

func searchForDBs() {
	searchResult = []string{}

	start := time.Now()
	err := filepath.Walk(rootPath, findFile)
	if err != nil {
		return
	}
	elapsed := time.Now().Sub(start)

	if noColour {
		fmt.Println("Searching for LevelDB databases from", rootPath, "took", elapsed)
	} else {
		fmt.Println(Info("Searching for LevelDB databases from ", rootPath, " took ", elapsed))
	}
	fmt.Println()

	if len(searchResult) > 0 {
		if noColour {
			fmt.Println(len(searchResult), "LevelDB databases found")
		} else {
			fmt.Println(Warn(len(searchResult), " LevelDB databases found"))
		}
	} else {
		if noColour {
			fmt.Println("0 LevelDB databases found")
		} else {
			fmt.Println(Fatal("0 LevelDB databases found"))
		}
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
			if noColour {
				fmt.Println("Access denied for", path)
			} else {
				fmt.Println(Warn("Access denied for ", path))
			}
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
					if noColour {
						fmt.Println("Found database at", absolute)
					} else {
						fmt.Println(Info("Found database at ", absolute))
					}
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

	// TODO: Instead of checking path, open MANIFEST-XXXX file and read string value
	if strings.Contains(dbPath, "\\IndexedDB\\") || strings.Contains(dbPath, "/IndexedDB/") {
		if noColour {
			fmt.Println("IndexedDB idb_cmp1 comparator not yet implemented, results will not be valid")
		} else {
			fmt.Println(Warn("IndexedDB idb_cmp1 comparator not yet implemented, results will not be valid"))
		}
		options.Comparer = idbCmp1{}
	}

	start := time.Now()

	db, err := leveldb.OpenFile(dbPath, options)

	if err != nil {
		if noColour {
			fmt.Println("Could not open DB:", err.Error())
		} else {
			fmt.Println(Fatal("Could not open DB: ", err.Error()))
		}
		fmt.Println()
		return
	}
	fmt.Println()

	defer db.Close()

	iter := db.NewIterator(nil, nil)

	if !quiet {
		if noColour {
			fmt.Println(fmt.Sprintf("%-56vValue:", "Key:"))
		} else {
			fmt.Println(Info(fmt.Sprintf("%-56vValue:", "Key:")))
		}

	}

	var data = [][]string{}
	for iter.Next() {
		key := iter.Key()
		keyName := string(key[:])

		byteValue, err := db.Get([]byte(key), nil)
		if err != nil {
			fmt.Println("Error reading Key: " + keyName)
			fmt.Println(err.Error())
			return
		}
		value := string(byteValue)

		data = append(data, []string{keyName, value})
	}

	if !quiet {
		if len(data) > 0 {
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
			if noColour {
				fmt.Println("Parsed database but no key/value pairs were found")
			} else {
				fmt.Println(Warn("Parsed database but no key/value pairs were found"))
			}
		}
	}

	if outputDir != "" {
		if len(data) > 0 {
			createCsvOutput(dbPath, data)
		}
	}

	iter.Release()
	err = iter.Error()
	checkError(err)
	if !quiet {
		fmt.Println()
	}

	elapsed := time.Now().Sub(start)
	if noColour {
		fmt.Println("Dumping LevelDB database took", elapsed)
	} else {
		fmt.Println(Info("Dumping LevelDB database took ", elapsed))
	}
	fmt.Println()
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
