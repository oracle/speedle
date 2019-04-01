//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

func TestGetRolesWithoutApp(t *testing.T) {
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
	ctx := adsapi.RequestContext{Subject: &subject}
	_, err = evaluator.GetAllGrantedRoles(ctx)
	if err == nil {
		t.Fatalf("Error should be returned without an application.")
		return
	}
	t.Logf("Returned error [%v].", err)
}

func TestGetRolesNotMatchApp(t *testing.T) {
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
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "dummy"}
	_, err = evaluator.GetAllGrantedRoles(ctx)
	if err == nil {
		t.Fatalf("Error should be returned if application is not matched.")
		return
	}
	t.Logf("Returned error [%v].", err)
}

func TestGetRolesNoRole(t *testing.T) {
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) > 0 {
		t.Fatalf("No role should be returned, but returned %q.", roles)
		return
	}
}

func TestGetDirectRolesMatch(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"type": "Applications",
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
					"user:william", "group:grp2"
				]
			},
			{
                "id": "rp3",
				"effect": "grant",
				"roles": ["role3"],
				"principals": [
					"user:cynthia", "group:grp3"
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
				Name: "william",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp2",
			},
		},
	}

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 1 {
		t.Fatalf("One role should be returned, but returned %q.", roles)
		return
	}
	foundRole1 := false
	foundRole2 := false
	for _, role := range roles {
		switch role {
		case "role1":
			foundRole1 = true
			break
		case "role2":
			foundRole2 = true
			break
		}
	}
	if foundRole1 {
		t.Fatalf("Role [role1] shouldn't be returned.")
	}
	if !foundRole2 {
		t.Fatalf("Role [role2] should be returned.")
	}
}

func TestGetHierarchyRolesMatch(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1"],
				"principals": [
					"user:bill", "group:grp2"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2"],
				"principals": [
					"user:william", "role:role1"
				]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role3"],
				"principals": [
					"user:bill"
				]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role4"],
				"principals": [
					"user:cynthia", "group:william"
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
				Name: "grp2",
			},
		},
	}

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"role3", "role2", "role1"}
	returnedRolesMap := make(map[string]bool)
	for _, role := range roles {
		returnedRolesMap[role] = true
	}
	if len(expectedRoles) != len(roles) {
		t.Fatalf("Expected %d roles, but returned %d roles", len(expectedRoles), len(roles))
	}
	for _, expectedRole := range expectedRoles {
		if _, ok := returnedRolesMap[expectedRole]; !ok {
			t.Fatalf("Expected role %s, but not returned", expectedRole)
		}
	}

}
