//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
)

func TestGetResourcesWithoutApp(t *testing.T) {
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
	_, err = evaluator.GetAllGrantedPermissions(adsapi.RequestContext{Subject: &subject, ServiceName: "dummy"})
	if err == nil {
		t.Fatalf("Error should be returned without an application.")
		return
	}
	t.Logf("Returned error [%v].", err)
}

func TestGetResourcesNotMatchApp(t *testing.T) {
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
	_, err = evaluator.GetAllGrantedPermissions(adsapi.RequestContext{Subject: &subject, ServiceName: "dummy"})
	if err == nil {
		t.Fatalf("Error should be returned if application is not matched.")
		return
	}
	t.Logf("Returned error [%v].", err)
}

func TestGetResourcesNoPolicy(t *testing.T) {
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

	var resources []pms.Permission
	resources, err = evaluator.GetAllGrantedPermissions(adsapi.RequestContext{Subject: &subject, ServiceName: "erp"})
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(resources) > 0 {
		t.Fatalf("No resource should be returned, but returned %q.", resources)
		return
	}
}

func TestGetResourcesNoPrincDef(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"policies": [
			{
				"id": "policy1",
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

	var resActsList []pms.Permission
	resActsList, err = evaluator.GetAllGrantedPermissions(adsapi.RequestContext{Subject: &subject, ServiceName: "erp"})
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(resActsList) != 1 {
		t.Fatalf("One resource should be returned, but returned %q.", resActsList)
		return
	}
	t.Logf("resActions %v", resActsList)
	if "/node1" != resActsList[0].Resource || "get" != resActsList[0].Actions[0] {
		t.Fatalf("Resource name is incorrect, expected [%v], but returned [%v].", "/node1:get", resActsList[0].Resource)
		return
	}
}

func TestGetResourcesMatch(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1"],
				"principals": [
					"user:bill", "group:grp1"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2"],
				"principals": [
					"role:role1"
				]
			}
			],
			"policies": [
			{
				"id": "policy1",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["get"]
				}
				],
				"principals": [
					["group:grp1", "user:bill"]
				]
			},
			{
				"id": "policy2",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node2",
					"actions": ["get"]
				}
				],
				"principals": [
					["user:william", "role:role1"]
				]
			},
			{
				"id": "policy3",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node3",
					"actions": ["get"]
				}
				],
				"principals": [
					["group:grp2", "user:cynthia"]
				]
			},
			{
				"id": "policy4",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node4",
					"actions": ["get"]
				}
				],
				"principals": [
					["role:role2"]
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
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp1",
			},
		},
	}

	var resActsList []pms.Permission
	resActsList, err = evaluator.GetAllGrantedPermissions(adsapi.RequestContext{Subject: &subject, ServiceName: "erp"})
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(resActsList) != 2 {
		t.Fatalf("One resource actions should be returned, but returned %q.", resActsList)
		return
	}
	foundResActs1 := false
	foundResActs2 := false
	foundResActs3 := false
	foundResActs4 := false
	for _, resActs := range resActsList {
		switch resActs.Resource {
		case "/node1":
			foundResActs1 = true
			break
		case "/node2":
			foundResActs2 = true
			break
		case "/node3":
			foundResActs3 = true
			break
		case "/node4":
			foundResActs4 = true
			break
		}
	}
	if !foundResActs1 {
		t.Fatalf("Role [/node1] should be returned.")
	}
	if foundResActs2 {
		t.Fatalf("Role [/node2] shouldn't be returned.")
	}
	if foundResActs3 {
		t.Fatalf("Role [/node3] shouldn't be returned.")
	}
	if !foundResActs4 {
		t.Fatalf("Role [/node4] should be returned.")
	}
}
