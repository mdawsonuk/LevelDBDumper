package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
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
