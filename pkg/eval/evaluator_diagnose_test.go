//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

func TestDiagnoseWithoutApp(t *testing.T) {
	const appStream = `
	{
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "bill",
			},
		},
	}
	result, err := evaluator.Diagnose(adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "dummy",
		Resource:    "dummy",
		Action:      "dummy",
	})
	if err == nil {
		t.Fatal("Error should be returned without an application.")
		return
	}
	t.Logf("Returned error [%v].", err)
	if result.Allowed {
		t.Fatal("Allowed should be false")
		return
	}
	if result.Reason != adsapi.SERVICE_NOT_FOUND {
		t.Fatal("reason should be SERVICE_NOT_FOUND, but is " + result.Reason.String())
		return
	}
	t.Logf("result [%v].", result)
}

func TestDiagnoseNotMatchApp(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp"
		}
		]
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "bill",
			},
		},
	}

	result, err := evaluator.Diagnose(adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "dummy",
		Resource:    "dummy",
		Action:      "dummy",
	})
	if err == nil {
		t.Fatal("Error should be returned if application is not matched.")
		return
	}
	t.Logf("Returned error [%v].", err)
	if result.Reason != adsapi.SERVICE_NOT_FOUND {
		t.Fatal("reason should be SERVICE_NOT_FOUND, but is " + result.Reason.String())
		return
	}
}

func TestDiagnoseNoPolicy(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp"
		}
		]
	}
	`

	preparePolicyDataInStore([]byte(appStream), t)
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "bill",
			},
		},
	}

	request := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "get",
	}
	result, err := evaluator.Diagnose(adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "dummy",
		Action:      "dummy",
	})
	if err != nil {
		t.Fatalf("Err: %v happens in diagnose", err)
	}
	if result.Allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if result.Reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Fatal("reason should be NO_APPLICABLE_POLICIES, but is " + result.Reason.String())
		return
	}
}

func TestDiagnoseNoPrincDef(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"policies": [
			{
				"name": "policy1",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["get"]
				}
				]
			}
			]
		}
		]
	}
	`

	preparePolicyDataInStore([]byte(appStream), t)

	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "bill",
			},
		},
	}

	request := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "get",
	}
	result, err := evaluator.Diagnose(adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "get",
	})
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !result.Allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
	if result.Reason != adsapi.GRANT_POLICY_FOUND {
		t.Fatal("reason should be GRANT_POLICY_FOUND, but is " + result.Reason.String())
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "post",
	}
	result, err = evaluator.Diagnose(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if result.Allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if result.Reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Fatal("reason should be NO_APPLICABLE_POLICIES, but is " + result.Reason.String())
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node2",
		Action:      "get",
	}
	result, err = evaluator.Diagnose(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if result.Allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if result.Reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Fatal("reason should be NO_APPLICABLE_POLICIES, but is " + result.Reason.String())
		return
	}

}
