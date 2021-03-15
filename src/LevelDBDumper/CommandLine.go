package main

import (
	"fmt"
	"path/filepath"

	"github.com/gookit/color"
)

var (
	help           bool
	rootPath       string
	quiet          bool
	outputType     string = "csv"
	outputDir      string
	outputFile     string
	batch          bool
	timezone       string
	noHeader       bool
	checkForUpdate bool
	cleanOutput    bool
)

func getArgs(args []string) {
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			help = true
		}
		if (args[i] == "-d" || args[i] == "--dir") && i+1 < len(args) {
			path, err := filepath.Abs(args[i+1])
			if err != nil {
				color.Red.Println(fmt.Sprintf("Unable to get absolute path of %s", path))
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
		if args[i] == "--no-header" {
			noHeader = true
		}
		if args[i] == "-u" || args[i] == "--check-update" {
			checkForUpdate = true
		}
		if args[i] == "-c" || args[i] == "--clean-output" {
			cleanOutput = true
		}
		if (args[i] == "-z" || args[i] == "--timezone") && i+1 < len(args) {
			timezone = args[i+1]
		}
	}
}

func printUsage() {
	fmt.Println("      h/help              Display this help message")
	fmt.Println("      d/dir               Directory to recursively process. This is required")
	fmt.Println("      q/quiet             Don't output all key/value pairs to console. This happens by default")
	fmt.Println("      t/outputType        Output type. Can be \"csv\" or \"json\"")
	fmt.Println("      o/outputDir         Directory to save all output results to. Required for any file output")
	fmt.Println("      f/outputFile        Filename to use when saving output. This will be appended with path and date")
	fmt.Println("      b/batch             Combine all output files into one file. Supported by \"csv\" and \"json\" file types")
	fmt.Println("      c/clean-output      Clean the file output of non-visual characters, such as \\u001")
	fmt.Println("      z/timezone          Specify the IANA timezone to use when using timestamps. Default is UTC")
	fmt.Println("      no-header           Don't display the header")
	fmt.Println("      u/check-update      Check for updates only")
	fmt.Println()
	fmt.Println("Short options (single letter) are prefixed with a single dash. Long commands are prefixed with two dashes")
	fmt.Println()
	fmt.Println("Examples: LevelDBParser.exe -d \"C:\\Temp\\leveldb\"")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -o \"C:\\Temp\" -q")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" --quiet --no-header --clean-output")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -b --outputType json -outputFile Evidence.json")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -t csv -f LevelDB.csv -o Evidence -b --quiet")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -t csv -o Evidence -b --timezone America/New_York")
	fmt.Println("          LevelDBParser.exe -d \"C:\\Temp\\leveldb\" -t json -o Evidence -b -z Local --quiet --clean-output")
	fmt.Println("          LevelDBParser.exe --check-update")
	fmt.Println("          LevelDBParser.exe --help")
	if !help {
		fmt.Println()
	}
}
