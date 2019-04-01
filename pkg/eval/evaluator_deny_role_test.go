//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

// test issue :https://gitlab-odx.oracledx.com/wcai/kauthz/issues/120
func TestGetRoles120(t *testing.T) {
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
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2"],
				"principals": ["role:role1"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role3"],
				"principals": ["role:role2"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role4"],
				"principals": ["role:role3"]
			},
			{
				"id": "rp5",
				"effect": "deny",
				"roles": ["role3"],
				"principals": ["user:bill"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"role2", "role1"}
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

func TestGetRoles120_1(t *testing.T) {
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
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2", "role3"],
				"principals": ["role:role1"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role3"],
				"principals": ["role:role2"]
			},
			{
				"id": "rp4",
				"effect": "deny",
				"roles": ["role2"],
				"principals": ["role:role1"]
			},
			{
				"id": "rp5",
				"effect": "deny",
				"roles": ["role3"],
				"principals":["role:role2"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	expectedRoles := []string{"role3", "role1"}
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

func TestGetRoles120_2(t *testing.T) {
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
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2", "role3"],
				"principals": ["role:role1"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role3"],
				"principals": ["role:role2"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role4"],
				"principals": ["role:role3"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["role1"],
				"principals": ["role:role4"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"role3", "role2", "role4", "role1"}
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

func TestGetRoles120_3(t *testing.T) {
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
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2", "role3"],
				"principals": ["role:role1"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role2"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role3", "role4"],
				"principals": ["role:role2"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["role5"],
				"principals": ["role:role4"]
			},
			{
				"id": "rp6",
				"effect": "deny",
				"roles": ["role1"],
				"principals": ["role:role5"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	expectedRoles := []string{"role3", "role2", "role4", "role5"}
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

func TestGetRoles120_4(t *testing.T) {
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
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2", "role3"],
				"principals": ["role:role1"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role4"],
				"principals": ["role:role3"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role2"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["role3", "role4"],
				"principals": ["role:role2"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["role5"],
				"principals": ["role:role4"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["role3"],
				"principals": ["role:role5"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	expectedRoles := []string{"role1", "role2", "role4", "role5"}
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

func TestGetRoles120_5(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["C", "B", "D"],
				"principals": ["role:A"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:C"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:B"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["H"],
				"principals": ["role:F"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["G"],
				"principals": ["role:D"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["B"],
				"principals": ["role:E"]
			},
			{
				"id": "rp8",
				"effect": "deny",
				"roles": ["D"],
				"principals": ["role:H"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	expectedRoles := []string{"A", "C", "D", "E", "G"}
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

func TestGetRoles120_6(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["C", "B"],
				"principals":["role:A"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["D"],
				"principals": ["role:C"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:D"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:E"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:B"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["D"],
				"principals": ["role:F"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	expectedRoles := []string{"A", "B", "C", "E", "F"}
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

func TestGetRoles120_7(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["C", "B"],
				"principals": ["role:A"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["D"],
				"principals": ["role:C"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:D"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:E"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:B"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["D"],
				"principals": ["role:F"]
			},
			{
				"id": "rp8",
				"effect": "grant",
				"roles": ["G"],
				"principals": ["role:A"]
			},
			{
				"id": "rp9",
				"effect": "deny",
				"roles": ["E"],
				"principals": ["role:G"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "B", "C", "D", "G"}
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

func TestGetRoles120_8(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["C", "B"],
				"principals": ["role:A"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["D"],
				"principals": ["role:C"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:D"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:E"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["D"],
				"principals": ["role:F"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "B", "C"}
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

func TestGetRoles120_9(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["B"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["C", "D"],
				"principals": ["role:A"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["D", "E"],
				"principals": ["role:B"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["G"],
				"principals": ["role:C"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:E"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["D"],
				"principals": ["role:G"]
			},
			{
				"id": "rp8",
				"effect": "deny",
				"roles": ["E"],
				"principals": ["role:D"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "B", "C", "E", "F", "G"}
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

func TestGetRoles120_10(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["D"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["B"],
				"principals": ["role:A"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:D"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:E"]
			},
			{
				"id": "rp6",
				"effect": "deny",
				"roles": ["B"],
				"principals": ["role:F"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["E"],
				"principals": ["role:B"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "D"}
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

func TestGetRoles120_11(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["D"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["B"],
				"principals": ["role:A"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["E"],
				"principals": ["role:D"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:E"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["G"],
				"principals": ["role:F"]
			},
			{
				"id": "rp7",
				"effect": "grant",
				"roles": ["F"],
				"principals": ["role:A"]
			},
			{
				"id": "rp8",
				"effect": "deny",
				"roles": ["B"],
				"principals": ["role:G"]
			},
			{
				"id": "rp9",
				"effect": "deny",
				"roles": ["E"],
				"principals": ["role:B"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "D", "E", "F", "G"}
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

func TestGetRoles120_12(t *testing.T) {
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
				"roles": ["A"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["B"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["D", "E"],
				"principals": ["role:A"]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["E", "F"],
				"principals": ["role:B"]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["G"],
				"principals": ["role:E"]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["G", "H"],
				"principals": ["role:F"]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["E"],
				"principals": ["user:bill"]
			},
			{
				"id": "rp8",
				"effect": "deny",
				"roles": ["H"],
				"principals": ["role:G"]
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "B", "D", "F", "G"}
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

func TestGetRoles120_13(t *testing.T) {
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
				"roles": ["role1", "role2"],
				"principals": [
					"user:bill",
					"role:A"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2", "role3"],
				"principals": [
					"role:role1",
					"role:B"

				]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role4"],
				"principals": [
					"role:role3",
					"role:C"
				]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role2"],
				"principals": [
					"user:bill",
					"role:role1"
				]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["role3", "role4"],
				"principals": [
					"role:role2",
					"role:E"
				]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["role5"],
				"principals": [
					"role:role4",
					"role:F"
				]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["role3"],
				"principals": [
					"role:role5",
					"role:A"
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"role1", "role2", "role4", "role5"}
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

func TestGetRoles120_14(t *testing.T) {
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
				"roles": ["A", "G"],
				"principals": [
					"user:bill",
					"role:role1"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["B", "D", "I"],
				"principals": [
					"role:role1",
					"role:A"

				]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["H"],
				"principals": [
					"role:role3",
					"role:G"
				]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["C"],
				"principals": [
					"role:B"
				]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["E"],
				"principals": [
					"role:role2",
					"role:C"
				]
			},
			{
				"id": "rp6",
				"effect": "grant",
				"roles": ["F"],
				"principals": [
					"role:role4",
					"role:D"
				]
			},
			{
				"id": "rp7",
				"effect": "grant",
				"roles": ["I"],
				"principals": [
					"role:role4",
					"role:H"
				]
			},
			{
				"id": "rp8",
				"effect": "deny",
				"roles": ["H"],
				"principals": [
					"role:role4",
					"role:A"
				]
			},
			{
				"id": "rp9",
				"effect": "deny",
				"roles": ["G"],
				"principals": [
					"role:role5",
					"role:F"
				]
			},
			{
				"id": "rp10",
				"effect": "deny",
				"roles": ["D"],
				"principals": [
					"role:role5",
					"role:E"
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"A", "B", "C", "E", "G", "I"}
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

func TestGetRoles120_15(t *testing.T) {
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
				"roles": ["role1","role2", "role3"],
				"principals": [
					"user:bill"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role12", "role1-12"],
				"principals": [
					"role:role1"
				]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role12-3","role12-3-3"],
				"principals": [
					"role:role3"
				]
			},
			{
				"id": "rp4",
				"effect": "grant",
				"roles": ["role1-12", "role121", "role12-3"],
				"principals": [
					"role:role12"
				]
			},
			{
				"id": "rp5",
				"effect": "grant",
				"roles": ["role12-3-3"],
				"principals": [
					"role:role12-3"
				]
			},
			{
				"id": "rp6",
				"effect": "deny",
				"roles": ["role12"],
				"principals": [
					"role:role2"
				]
			},
			{
				"id": "rp7",
				"effect": "deny",
				"roles": ["role3"],
				"principals": [
					"role:role1-12"
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

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

	expectedRoles := []string{"role1", "role2", "role1-12"}
	returnedRolesMap := make(map[string]bool)
	for _, role := range roles {
		returnedRolesMap[role] = true
	}
	if len(expectedRoles) != len(roles) {
		t.Fatalf("Expected roles: %v, but returned roles: %v", expectedRoles, roles)
	}
	for _, expectedRole := range expectedRoles {
		if _, ok := returnedRolesMap[expectedRole]; !ok {
			t.Fatalf("Expected role %s, but not returned", expectedRole)
		}
	}
}
