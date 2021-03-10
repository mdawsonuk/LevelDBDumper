package main

import (
	"bytes"
)

const (
	// From https://source.chromium.org/chromium/chromium/src/+/master:content/browser/indexed_db/indexed_db_leveldb_coding.h;l=139
	GLOBAL_METADATA   int = 0
	DATABASE_METADATA int = 1
	OBJECT_STORE_DATA int = 2
	EXISTS_ENTRY      int = 3
	INDEX_DATA        int = 4
	INVALID_TYPE      int = 5
	BLOB_ENTRY        int = 6
)

// TODO: https://source.chromium.org/chromium/chromium/src/+/master:content/browser/indexed_db/indexed_db_leveldb_operations.cc?q=idb_cmp1
// https://stackoverflow.com/questions/35074659/how-to-access-google-chromes-indexeddb-leveldb-files

type idbCmp1 struct{}

func (idbCmp1) Compare(a, b []byte) int {
	// TODO: https://source.chromium.org/chromium/chromium/src/+/master:content/browser/indexed_db/indexed_db_leveldb_coding.cc;l=840
	return bytes.Compare(a, b)
}

func (idbCmp1) Name() string {
	return "idb_cmp1"
}

func (idbCmp1) Separator(dst, a, b []byte) []byte {
	return nil
}

func (idbCmp1) Successor(dst, b []byte) []byte {
	return nil
}
