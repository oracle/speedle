//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

func TestRolePolicyWithoutResWithCondition(t *testing.T) {
	//role policy without resource means any resource
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",			
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["role1", "role2"],
					"principals":  [
						"user:cynthia", "group:grp3"
					],
					"condition": "rt_str_attr=='abc'"
				},
				{
					"id": "rp2",
					"effect": "grant",
					"roles": ["role1"],
					"principals": [
						"user:bill", "group:grp1"
					]
				},
				{
					"id": "rp3",
					"effect": "grant",
					"roles": ["role2"],
					"principals":  [
						"user:william", "group:grp2"
					]
				},
				{
					"id": "rp4",
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
				Name: "cynthia",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp3",
			},
		},
	}

	context := adsapi.RequestContext{Subject: &subject, ServiceName: "erp", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"rt_str_attr": "abc"}}

	var roles []string

	roles, err = evaluator.GetAllGrantedRoles(context)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 3 {
		t.Fatalf("Three roles should be returned, but returned %q.", roles)
		return
	}
	foundRole1 := false
	foundRole2 := false
	foundRole3 := false
	for _, role := range roles {
		switch role {
		case "role1":
			foundRole1 = true
			break
		case "role2":
			foundRole2 = true
			break
		case "role3":
			foundRole3 = true
			break
		}
	}
	if !foundRole1 {
		t.Fatalf("Role [role1] should be returned.")
	}
	if !foundRole2 {
		t.Fatalf("Role [role2] should be returned.")
	}
	if !foundRole3 {
		t.Fatalf("Role [role3] should be returned.")
	}
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 1 {
		t.Fatalf("1 role should be returned, but returned %q.", roles)
		return
	}
}

func TestRolePolicyWithoutResCondifiton(t *testing.T) {
	//role policy without resource means any resource
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["role1", "role2"],
					"principals": [
						"user:cynthia", "group:grp3"
					]
				},
				{
					"id": "rp2",
					"effect": "grant",
					"roles": ["role1"],
					"principals": [
						"user:bill", "group:grp1"
					]
				},
				{
					"id": "rp3",
					"effect": "grant",
					"roles": ["role2"],
					"principals": [
						"user:william", "group:grp2"
					]
				},
				{
					"id": "rp4",
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
				Name: "cynthia",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp3",
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
	if len(roles) != 3 {
		t.Fatalf("Two role should be returned, but returned %q.", roles)
		return
	}
	foundRole1 := false
	foundRole2 := false
	foundRole3 := false
	for _, role := range roles {
		switch role {
		case "role1":
			foundRole1 = true
			break
		case "role2":
			foundRole2 = true
			break
		case "role3":
			foundRole3 = true
			break
		}
	}
	if !foundRole1 {
		t.Fatalf("Role [role1] should be returned.")
	}
	if !foundRole2 {
		t.Fatalf("Role [role2] should be returned.")
	}
	if !foundRole3 {
		t.Fatalf("Role [role3] should be returned.")
	}
}

/*
func TestRolePolicyHierarchy(t *testing.T) {
	//role policy without resource means any resource
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["role1"],
					"principals": [
						"user:cynthia", "group:grp3"
					]
				},
				{
					"id": "rp2",
					"effect": "grant",
					"roles": ["role2"],
					"principals": [
						"role:role1"
					]
				},
				{
					"id": "rp3",
					"effect": "deny",
					"roles": ["role2"],
					"principals": [
						"group:grp4"
					]
				},
				{
					"id": "rp4",
					"effect": "grant",
					"roles": ["role1"],
					"principals": [
						"user:bill", "group:grp1"
					]
				},
				{
					"id": "rp5",
					"effect": "grant",
					"roles": ["role2"],
					"principals": [
						"user:william", "group:grp2"
					]
				},
				{
					"id": "rp6",
					"effect": "grant",
					"roles": ["role3"],
					"principals": [
						"role:role2"
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
			&adsapi.Principal{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "cynthia",
			},
			&adsapi.Principal{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp3",
			},
		},
		Attributes: nil,
	}

	var roles []string
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 3 {
		t.Fatalf("Three role should be returned, but returned %q.", roles)
		return
	}
	foundRole1 := false
	foundRole2 := false
	foundRole3 := false
	for _, role := range roles {
		switch role {
		case "role1":
			foundRole1 = true
			break
		case "role2":
			foundRole2 = true
			break
		case "role3":
			foundRole3 = true
			break
		}
	}
	if !foundRole1 {
		t.Fatalf("Role [role1] should be returned.")
	}
	if !foundRole2 {
		t.Fatalf("Role [role2] should be returned.")
	}
	if !foundRole3 {
		t.Fatalf("Role [role3] should be returned.")
	}

	subject = adsapi.Subject{
		Principals: []*adsapi.Principal{
			&adsapi.Principal{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "cynthia",
			},
			&adsapi.Principal{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp4",
			},
		},
		Attributes: nil,
	}

	ctx = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
	}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 2 {
		t.Fatalf("2 roles should be returned, but returned %q.", roles)
		return
	}

	foundRole1 = false
	foundRole3 = false
	for _, role := range roles {
		switch role {
		case "role1":
			foundRole1 = true
			break
		case "role3":
			foundRole3 = true
			break
		}
	}
	if foundRole1 {
		t.Fatalf("Role [role1] shouldn't be returned.")
	}
	if foundRole3 {
		t.Fatalf("Role [role3] shouldn't be returned.")
	}
}
*/

func TestRolePolicyWithResWithCondition(t *testing.T) {
	//role policy without resource means any resource
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["role1"],
					"principals": [
						"user:cynthia", "group:grp3"
					],
					"resources": ["/node1","/node2"],	
					"condition": "rt_str_attr=='abc'"
				},
				{
					"id": "rp2",
					"effect": "grant",
					"roles": ["role2"],
					"principals": [
					 	"role:role1"
					],
					"resources": ["/node1"],
					"condition": "rt_str_attr=='abc'"
				},
				{
					"id": "rp3",
					"effect": "grant",
					"roles": ["role1"],
					"principals": [
						"user:bill", "group:grp1"
					]
				},
				{
					"id": "rp4",
					"effect": "grant",
					"roles": ["role2"],
					"principals": [
						"user:william", "group:grp2"
					]
				},
				{
					"id": "rp5",
					"effect": "grant",
					"roles": ["role3"],
					"principals": [
						"role:role2"
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
				},
				{
					"res": "/node2",
					"acts": ["get"]
				}
				],
				"principals": [
					["role:role3"]
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
				Name: "cynthia",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "grp3",
			},
		},
	}

	context := adsapi.RequestContext{Subject: &subject, ServiceName: "erp", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"rt_str_attr": "abc"}}

	var roles []string

	roles, err = evaluator.GetAllGrantedRoles(context)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 3 {
		t.Fatalf("Three roles should be returned, but returned %q.", roles)
		return
	}
	foundRole1 := false
	foundRole2 := false
	foundRole3 := false
	for _, role := range roles {
		switch role {
		case "role1":
			foundRole1 = true
			break
		case "role2":
			foundRole2 = true
			break
		case "role3":
			foundRole3 = true
			break
		}
	}
	if !foundRole1 {
		t.Fatalf("Role [role1] should be returned.")
	}
	if !foundRole2 {
		t.Fatalf("Role [role2] should be returned.")
	}
	if !foundRole3 {
		t.Fatalf("Role [role3] should be returned.")
	}

	context = adsapi.RequestContext{Subject: &subject, ServiceName: "erp", Resource: "/node2", Action: "get", Attributes: map[string]interface{}{"rt_str_attr": "abc"}}
	roles, err = evaluator.GetAllGrantedRoles(context)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 1 {
		t.Fatalf("One roles should be returned, but returned %q.", roles)
		return
	}
	ctx := adsapi.RequestContext{Subject: &subject, ServiceName: "erp"}
	roles, err = evaluator.GetAllGrantedRoles(ctx)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if len(roles) != 0 {
		t.Fatalf("0 role should be returned, but returned %q.", roles)
		return
	}

	var allowed bool
	context = adsapi.RequestContext{Subject: &subject, ServiceName: "erp", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"rt_str_attr": "abc"}}
	allowed, _, err = evaluator.IsAllowed(context)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", context)
		return
	}

	context = adsapi.RequestContext{Subject: &subject, ServiceName: "erp", Resource: "/node2", Action: "get", Attributes: map[string]interface{}{"rt_str_attr": "abc"}}
	allowed, _, err = evaluator.IsAllowed(context)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should not be allowed.", context)
		return
	}

}
