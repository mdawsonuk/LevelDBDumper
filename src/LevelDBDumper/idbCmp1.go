package main

import (
	"bytes"
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
