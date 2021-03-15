package main

import (
	"os"
	"testing"
)

func TestGetComparator(t *testing.T) {
	f, _ := os.Create("MANIFEST-0001")
	f.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x54, 0x65, 0x73, 0x74, 0x02})
	comparator := getComparator(".")
	if comparator != "Test" {
		t.Errorf("getComparator was incorrect, actual: %s, expected: Test", comparator)
	}
	f.Close()
	os.Remove("MANIFEST-0001")
}

func TestGetMalformedComparator(t *testing.T) {
	f, _ := os.Create("MANIFEST-0001")
	f.Write([]byte{0x00})
	comparator := getComparator(".")
	if comparator != "Unknown" {
		t.Errorf("getComparator with malformed file was incorrect, actual: %s, expected: Unknown", comparator)
	}
	f.Close()
	os.Remove("MANIFEST-0001")
}
