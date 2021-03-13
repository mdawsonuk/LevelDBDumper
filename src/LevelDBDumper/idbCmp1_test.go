package main

import "testing"

func TestName(t *testing.T) {
	var comparator = idbCmp1{}
	ret := comparator.Name()
	if ret != "idb_cmp1" {
		t.Errorf("idbCmp1::Name() failed, actual: %s, expected: idb_cmp1", ret)
	}
}

func TestSeparator(t *testing.T) {
	var comparator = idbCmp1{}
	ret := comparator.Separator(nil, nil, nil)
	if ret != nil {
		t.Errorf("idbCmp1::Separator() failed, actual: %s, expected: nil", ret)
	}
}

func TestSuccessor(t *testing.T) {
	var comparator = idbCmp1{}
	ret := comparator.Successor(nil, nil)
	if ret != nil {
		t.Errorf("idbCmp1::Successor() failed, actual: %s, expected: nil", ret)
	}
}
