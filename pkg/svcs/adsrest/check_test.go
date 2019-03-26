//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package adsrest

import (
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"net/http/httptest"

	"github.com/oracle/speedle/pkg/assertion"
	"github.com/oracle/speedle/pkg/svcs"
)

func TestSettingPrincipalHeader(t *testing.T) {
	assertserver := assertion.NewTestServer(t, nil)
	defer assertserver.Close()
	t.Logf("assert endpoint:%s", assertserver.URL)

	adsserver, err := newADSServerWithAsserter(assertserver.URL, t)
	if err != nil {
		t.Fatal("Failed to start ADS! Error:", err)
	}
	defer adsserver.Close()
	t.Logf("ads endpoint:%s", adsserver.URL)

	request := JsonContext{
		Subject: &JsonSubject{
			TokenType: "WERCKER",
			Token:     "testtoken",
		},
		ServiceName: "fakservice", //fakeservice is a predefined service in fakestore.json
	}
	isAllowedURL := adsserver.URL + "/authz-check/v1/is-allowed"
	buf, err := json.Marshal(request)
	if err != nil {
		t.Fatal("failed to marshal test request")
	}
	req, err := http.NewRequest("POST", isAllowedURL, bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal("failed to make test request")
	}

	var client *http.Client
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	t.Logf("Response status: %d", resp.StatusCode)
	principals := resp.Header.Get(svcs.PrincipalsHeader)
	if len(principals) == 0 {
		t.Fatal("No principal is returned!")
	} else {
		t.Logf("Principals in header:%s.", principals)
	}

	request = JsonContext{
		ServiceName: "fakservice",
	}
	buf, err = json.Marshal(request)
	if err != nil {
		t.Fatal("failed to marshal test request")
	}
	req, err = http.NewRequest("POST", isAllowedURL, bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal("failed to make test request")
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	t.Logf("Response status: %d", resp.StatusCode)
	principals = resp.Header.Get(svcs.PrincipalsHeader)
	if len(principals) != 0 {
		t.Logf("Principal returned is:%s.", principals)
		t.Fatal("Principal returned is not nil!")
	}

	request = JsonContext{
		Subject: &JsonSubject{
			TokenType: "FAKETYPE",
			Token:     "",
		},
		ServiceName: "fakservice",
	}
	buf, err = json.Marshal(request)
	if err != nil {
		t.Fatal("failed to marshal test request")
	}
	req, err = http.NewRequest("POST", isAllowedURL, bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal("failed to make test request")
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	t.Logf("Response status: %d", resp.StatusCode)
	principals = resp.Header.Get(svcs.PrincipalsHeader)
	if len(principals) != 0 {
		t.Logf("Principal returned is:%s.", principals)
		t.Fatal("Principal returned is not nil!")
	}
}

func newADSServerWithAsserter(assertserverendpoint string, t *testing.T) (*httptest.Server, error) {
	conf := GenerateServerConfig()
	asconfig := &assertion.AsserterConfig{
		Endpoint: assertserverendpoint + "/assert",
	}
	conf.AsserterWebhookConfig = asconfig
	return NewTestServerWithConfig(conf)
}
