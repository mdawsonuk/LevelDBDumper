package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
)

// VERSION of LevelDB Dumper
const VERSION string = "3.0.0-beta.1"

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
		os.Exit(2)
	}

	testFile, err := os.Open(rootPath)
	if err != nil {
		color.Yellow.Println(fmt.Sprintf("Unable to open %s - make sure you haven't escaped the path with \\\"", rootPath))
		testFile.Close()
		os.Exit(2)
	}
	testFile.Close()

	if !isAdmin() {
		color.Red.Println("You should run LevelDB Dumper with root/Administrator privileges")
	} else {
		color.FgLightBlue.Println("Running LevelDB Dumper with root/Administrator privileges")
	}
	fmt.Println()

	start := time.Now()

	// See LevelDBSearch.go
	searchForDBs()
	// See LevelDBParse.go
	readDBs()
	// See FileOutputs.go
	writeDBInfo()

	elapsed := time.Now().Sub(start)
	color.FgLightBlue.Println(fmt.Sprintf("Completed search in %v", elapsed))

	if !offline {
		needsUpdate, latestVersion := checkUpdate(VERSION)

		if !needsUpdate {
			color.Magenta.Println("You are using the latest version of LevelDB Dumper")
			if checkForUpdate {
				os.Exit(0)
			}
		} else if checkForUpdate {
			color.Cyan.Println(fmt.Sprintf("Version %s is now available for LevelDB Dumper - please update!", latestVersion))
			os.Exit(0)
		}
	}

	os.Exit(0)
}
