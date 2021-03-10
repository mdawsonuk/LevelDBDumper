package main

import "time"

// ParsedDB holds data for a parsed LevelDB database
type ParsedDB struct {
	path         string
	modifiedTime time.Time
	keys         []string
	values       []string
}
