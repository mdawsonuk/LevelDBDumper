package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/color"
)

func searchForDBs() {
	searchResult = []string{}

	start := time.Now()
	err := filepath.Walk(rootPath, findFile)
	if err != nil {
		return
	}
	elapsed := time.Now().Sub(start)

	if len(searchResult) > 0 && !quiet {
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
