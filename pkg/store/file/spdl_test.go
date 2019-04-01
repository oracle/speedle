//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	testFile := `abcd
 1234   
	3456	
util        
 as   df   
`
	expected := []lineCtx{
		{no: 1, origin: "abcd"},
		{no: 2, origin: " 1234   "},
		{no: 3, origin: "\t3456\t"},
		{no: 4, origin: "util        "},
		{no: 5, origin: " as   df   "},
	}
	sr := strings.NewReader(testFile)
	lc := lineCtx{}
	r := bufio.NewReader(sr)
	for i := 0; i < len(expected); i++ {
		if err := readLine(r, &lc); err != nil {
			t.Fatalf("Failed to read line due to error %v\n", err)
		}
		expt := expected[i]
		if expt.no != lc.no || expt.origin != lc.origin {
			t.Fatalf("Unexpcted result at line %d with origin content [%s]", lc.no, lc.origin)
		}
	}
	if err := readLine(r, &lc); err != nil {
		if err != io.EOF {
			t.Fatal("Unexpected error at the end of a file.")
		}
		return
	}
	t.Fatal("Error with EOF should be here.")
}

func runDetermineTypeTestCase(t *testing.T, line string, expected *lineCtx, eerr error) {
	lineCtx := lineCtx{origin: line}
	err := determineType(&lineCtx)
	if eerr != nil && err == nil {
		t.Fatalf("Error %v should be found here for line \"%s\"", eerr, line)
	}
	if eerr != nil && err != nil {
		// No need to check for expected error
		return
	}
	if eerr == nil && err != nil {
		t.Fatalf("Error %v should not be found here for line \"%s\"", err, line)
	}

	if lineCtx.trimed != expected.trimed {
		t.Fatalf("Unexpected trimed line \"%s\" for line \"%s\"", lineCtx.trimed, line)
	}
	if lineCtx.ltype != expected.ltype {
		t.Fatalf("Unexpected line type %s for line \"%s\"", lineCtx.ltype, line)
	}
	if lineCtx.section != expected.section {
		t.Fatalf("Unexpected section name %s for line \"%s\"", lineCtx.section, line)
	}
}

func TestDetermineType(t *testing.T) {
	runDetermineTypeTestCase(t, "     ", &lineCtx{trimed: "", ltype: lineEmpty}, nil)
	runDetermineTypeTestCase(t, " \t   ", &lineCtx{trimed: "", ltype: lineEmpty}, nil)
	runDetermineTypeTestCase(t, "", &lineCtx{trimed: "", ltype: lineEmpty}, nil)
	runDetermineTypeTestCase(t, "[abcd]", &lineCtx{trimed: "[abcd]", ltype: lineSection, section: "abcd"}, nil)
	runDetermineTypeTestCase(t, "   \t[test]    ", &lineCtx{trimed: "[test]", ltype: lineSection, section: "test"}, nil)
	runDetermineTypeTestCase(t, "[abcd]  # comment", &lineCtx{trimed: "[abcd]", ltype: lineSection, section: "abcd"}, nil)
	runDetermineTypeTestCase(t, "#[abcd]  comment", &lineCtx{trimed: "", ltype: lineEmpty}, nil)
	runDetermineTypeTestCase(t, "[", &lineCtx{trimed: ""}, fmt.Errorf("error"))
	runDetermineTypeTestCase(t, "]fdsa", &lineCtx{trimed: ""}, fmt.Errorf("error"))
	runDetermineTypeTestCase(t, "[abcd", &lineCtx{trimed: ""}, fmt.Errorf("error"))
	runDetermineTypeTestCase(t, "[abcd] eutkl", &lineCtx{trimed: ""}, fmt.Errorf("error"))
	runDetermineTypeTestCase(t, "[]", &lineCtx{trimed: ""}, fmt.Errorf("error"))
	runDetermineTypeTestCase(t, "# grant user fdsa", &lineCtx{trimed: "", ltype: lineEmpty}, nil)
	runDetermineTypeTestCase(t, "    grant user fdsa", &lineCtx{trimed: "grant user fdsa", ltype: linePolicyDef}, nil)
	runDetermineTypeTestCase(t, "    grant user fdsa   # fdsaqwer", &lineCtx{trimed: "grant user fdsa", ltype: linePolicyDef}, nil)
}

func TestReadSPDLWithoutLock(t *testing.T) {
	store := Store{
		FileLocation: "./spdl_test.spdl",
	}
	ps, err := store.readSPDLWithoutLock()
	if err != nil {
		t.Fatalf("Can't read PDL file due to error %v", err)
	}
	if len(ps.Services) != 2 {
		t.Fatalf("Wrong service count, expected 2, but there are %d services", len(ps.Services))
	}

	for _, service := range ps.Services {
		switch service.Name {
		case "service1":
		case "service2":
		default:
			t.Fatalf("Unknown service name")
		}
	}
}
