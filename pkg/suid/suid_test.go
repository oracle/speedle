//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package suid

import (
	"testing"
)

func TestSUIDLength(t *testing.T) {
	suid := New()
	suidStr := suid.String()
	t.Logf("SUID is %s.", suidStr)
	if len(suidStr) != 20 {
		t.Fatalf("SUID length is wrong, expected 20 but returned %d.", len(suidStr))
	}
}

func TestNoDuplicate(t *testing.T) {
	idmap := make(map[string]bool)
	// Test 1M times, test if there is any duplicated items
	for i := 0; i < 1000*1000; i++ {
		suid := New()
		str := suid.String()
		_, found := idmap[str]
		if found {
			t.Fatalf("SUID \"%s\" is duplicated.", str)
		}
		idmap[str] = true
	}
}
