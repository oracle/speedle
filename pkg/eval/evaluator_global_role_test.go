//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

// test issue :https://gitlab-odx.oracledx.com/wcai/kauthz/issues/232

// a single global role
// global: grant role1
// local: no rolepolicy
func TestGetRoles232_1(t *testing.T) {
	const appStream string = `
	{
		"services": [
		{
			"name": "global",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1"],
				"principals": [
					"user:bill"
				]
			}
			]
		},
		{
			"name": "erp",
			"type": "Applications",
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
					["role:role1"]
				]
			}
			]
		}
		]
	}
	`
	expectedRoles := []string{"role1"}

	evaluate(t, appStream, expectedRoles)

	// Global role should take effect in isAllowed API too.
	// Bill should be allowed to get /node1.
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}

	request := adsapi.RequestContext{
		Subject: &adsapi.Subject{
			Principals: []*adsapi.Principal{
				{
					Type: adsapi.PRINCIPAL_TYPE_USER,
					Name: "bill",
				},
			},
		},
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "get",
	}
	allowed, _, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}

}

// combine global and local role
// global: grant role1 role2
// local: grant role2 role3
func TestGetRoles232_2(t *testing.T) {
	const appStream string = `
	{
		"services": [
		{
			"name": "global",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1","role2"],
				"principals": [
					"user:bill"
				]
			}
			]
		},
		{
			"name": "erp",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2","role3"],
				"principals": [
					"user:bill"
				]
			}
			]
		}
		]
	}
	`
	expectedRoles := []string{"role1", "role2", "role3"}

	evaluate(t, appStream, expectedRoles)
}

// transitivity
// global: grant role1
// local: grant role:role1 role2
func TestGetRoles232_3(t *testing.T) {
	const appStream string = `
	{
		"services": [
		{
			"name": "global",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1"],
				"principals": [
					"user:bill"
				]
			}
			]
		},
		{
			"name": "erp",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["role2"],
				"principals": [
					"role:role1"
				]
			}
			]
		}
		]
	}
	`
	expectedRoles := []string{"role1", "role2"}

	evaluate(t, appStream, expectedRoles)
}

// deny transitivity
// global: grant role1 role2
//			deny role3
// local: grant role2 role3
//			deny role1
func TestGetRoles232_4(t *testing.T) {
	const appStream string = `
	{
		"services": [
		{
			"name": "global",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1","role2"],
				"principals": [
					"user:bill"
				]
			},
			{
				"id": "rp2",
				"effect": "deny",
				"roles": ["role3"],
				"principals": [
					"user:bill"
				]
			}
			]
		},
		{
			"name": "erp",
			"type": "Applications",
			"rolePolicies": [
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["role2","role3"],
				"principals": [
					"user:bill"
				]
			},
			{
				"id": "rp4",
				"effect": "deny",
				"roles": ["role1"],
				"principals": [
					"user:bill"
				]
			}
			]
		}
		]
	}
	`
	expectedRoles := []string{"role2"}

	evaluate(t, appStream, expectedRoles)
}

func evaluate(t *testing.T, appStream string, expectedRoles []string) {

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

	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err := evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}

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
