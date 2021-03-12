// +build windows

package main

import "testing"

func TestNotAdministrator(t *testing.T) {
	isAdmin := isAdmin()
	if isAdmin {
		t.Errorf("isAdmin, actual: %t, expected: false", isAdmin)
	}
}
