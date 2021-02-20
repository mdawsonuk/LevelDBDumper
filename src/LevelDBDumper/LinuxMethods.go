// +build !windows

package main

import (
	"os"
)

func isAdmin() bool {
	euid := os.Geteuid()
	return euid == 0
}
