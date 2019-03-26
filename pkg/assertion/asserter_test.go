//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package assertion

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

func TestAssertion(t *testing.T) {

	server := NewTestServer(t, nil)
	defer server.Close()

	asserter, err := getAsserter(server.URL+"/assert", t)

	t.Logf("load asserter: %v, err: %v", asserter, err)
	if asserter == nil || err != nil {
		t.Error("asserter is nil, err ", err)

	}

	ar, errAssert := asserter.AssertToken("testtoken", "WERCKER", "", nil)
	t.Logf("auth result: %v, err: %v", ar, errAssert)
	if errAssert != nil {
		t.Error("assertion, err ", errAssert)
	}

	if len(ar.Principals) != 1 {
		t.Fatalf("One principals should be returned but returned %d principals.\n", len(ar.Principals))
		t.FailNow()
	}

	if ar.Principals[0].Name != "testUser" || ar.Principals[0].Type != adsapi.PRINCIPAL_TYPE_USER {
		t.Fatalf("returned user should be testUser, actually returned %v .\n", ar.Principals[0])
		t.FailNow()
	}

	// no token
	ar, errAssert = asserter.AssertToken("", "", "", nil)
	t.Logf("auth result: %v, err: %v", ar, errAssert)
	if errAssert == nil {
		t.Fatalf("should failed with error")
		t.FailNow()
	}

	// invalid token
	ar, errAssert = asserter.AssertToken("test-token", "", "", nil)
	t.Logf("auth result: %v, err: %v", ar, errAssert)
	if errAssert == nil {
		t.Fatalf("should failed with error")
		t.FailNow()
	}
	if !strings.Contains(errAssert.Error(), strconv.Itoa(http.StatusBadRequest)) {
		t.Fatalf("should failed with error code: %d, actually failed with error code: %s.\n", http.StatusBadRequest, errAssert.Error())
		t.FailNow()
	}

}

func getAsserter(endpoint string, t *testing.T) (TokenAsserter, error) {
	conf := &AsserterConfig{
		Endpoint: endpoint,
	}
	t.Logf("endpoint: %s", endpoint)
	return NewAsserter(conf, nil)

}
