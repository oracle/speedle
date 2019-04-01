//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"fmt"
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
)

func TestGetPermissions(t *testing.T) {
	alice := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "alice",
			},
		},
	}
	bill := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "bill",
			},
		},
	}
	testCases := []struct {
		stream  string
		request adsapi.RequestContext
		results []pms.Permission
	}{
		{
			stream: `
			{
		"services": [
		{
			"name": "erp",
			"policies": [
			{
				"id": "policy1",
				"name": "policy1",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["get"]
				}
				],
				"principals": [
					["user:alice"]
				]
			},
			{
				"id": "policy2",
				"name": "policy2",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["get"]
				}
				],
				"principals": [
					["user:bill"]
				]
			},			
			{
				"id": "policy3",
				"name": "policy3",
				"effect": "deny",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["user:alice"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Attributes: map[string]interface{}{}},
			results: []pms.Permission{{Resource: "/node1", Actions: []string{"get"}}, {Resource: "/node2", Actions: []string{"get"}}},
		},
		{
			stream: `
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
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["get"]
				}
				],
				"principals": [
					["user:alice"]
				]
			},
			{
				"id": "policy2",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["get"]
				}
				],
				"principals": [
					["user:bill"]
				]
			},			
			{
				"id": "policy3",
				"effect": "deny",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["user:alice"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &bill, ServiceName: "erp", Attributes: map[string]interface{}{}},
			results: []pms.Permission{{Resource: "/node1", Actions: []string{"get", "create", "delete"}}, {Resource: "/node2", Actions: []string{"get"}}},
		},
		{
			stream: `
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
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["create"]
				}
				],
				"principals": [
					["user:alice"]
				]
			},
			{
				"id": "policy2",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["create"]
				}
				],
				"principals": [
					["user:bill"]
				]
			},			
			{
				"id": "policy3",
				"effect": "deny",
				"permissions": [
				{
					"resourceExpression": "/node.*",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["user:alice"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Attributes: map[string]interface{}{}},
			results: []pms.Permission{{Resource: "/node1", Actions: []string{"get"}}},
		},
		{
			stream: `
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
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["create"]
				}
				],
				"principals": [
					["user:bill", "group:grp1"]
				]
			},
			{
				"id": "policy2",
				"effect": "deny",
				"permissions": [
				{
					"resourceExpression": "/node.*",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["user:alice"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &bill, ServiceName: "erp", Attributes: map[string]interface{}{}},
			results: []pms.Permission{},
		},
		{
			stream: `
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
					"actions": ["get","create","delete"]
				},
				{
					"resource": "/node2",
					"actions": ["create"]
				}
				],
				"principals": [
					["user:alice", "user:bill"]
				]
			},
			{
				"id": "policy2",
				"effect": "deny",
				"principals": [
					["user:alice"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Attributes: map[string]interface{}{}},
			results: []pms.Permission{},
		},
	}

	for i, tc := range testCases {
		preparePolicyDataInStore([]byte(tc.stream), t)
		eval, err := NewWithStore(conf, testPS)
		if err != nil {
			t.Errorf("case number %d, error creating evaluator : %v", i, err)
			continue
		}
		got, err := eval.GetAllGrantedPermissions(tc.request)
		if err != nil {
			t.Errorf("case number %d, Error found in getAllGrantedPermissions: %v", i, err)
		}
		if len(got) != len(tc.results) {
			t.Errorf("case number %d, expect %d permissions, but returned %d permissions", i, len(tc.results), len(got))
		}
		for _, gotPerm := range got {
			permFound := false
			for _, resultPerm := range tc.results {
				if (gotPerm.Resource == resultPerm.Resource) && (len(gotPerm.Actions) == len(resultPerm.Actions)) {
					var actionFound bool
					for _, gotAction := range gotPerm.Actions {
						actionFound = false
						for _, resultAction := range resultPerm.Actions {
							if gotAction == resultAction {
								actionFound = true
								break
							}
						}
						if !actionFound {
							break
						}
					}
					if actionFound {
						permFound = true
					}

				}
			}
			if !permFound {
				t.Errorf("case number %d, Permission is unexpected:%v", i, gotPerm)
			}
		}
	}
}

//empty resource and empty resource expression means any resource
//empty principals means any principal
func TestGetRoles_empty(t *testing.T) {
	alice := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "alice",
			},
		},
	}
	bill := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "bill",
			},
		},
	}
	testCases := []struct {
		stream        string
		request       adsapi.RequestContext
		expectedRoles []string
	}{
		{ // any principal, any resource
			stream: `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",			
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["anyuser"]
				}
			]
		}
		]
	}
	`,
			request:       adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Attributes: map[string]interface{}{}},
			expectedRoles: []string{"anyuser"},
		},
		{ //make sure role policy support resource exp
			stream: `
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
						"user:alice"
					],
					"ResourceExpressions": ["/node.*"]
				}
			]
		}
		]
	}
	`,
			request:       adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Resource: "/node1213131/wqwq", Attributes: map[string]interface{}{}},
			expectedRoles: []string{"role1"},
		},
		{ //any principal
			stream: `
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
					"ResourceExpressions": ["/node.*"]
				}
			]
		}
		]
	}
	`,
			request:       adsapi.RequestContext{Subject: &bill, ServiceName: "erp", Resource: "/node1213131/wqwq", Attributes: map[string]interface{}{}},
			expectedRoles: []string{"role1"},
		},
		{ //any principal
			stream: `
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
					"ResourceExpressions": ["/node.*"]
				}
			]
		}
		]
	}
	`,
			request:       adsapi.RequestContext{Subject: &bill, ServiceName: "erp", Resource: "/aaaanode1213131/wqwq", Attributes: map[string]interface{}{}},
			expectedRoles: []string{},
		},
		{ //any resource & built-in attr
			stream: `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",			
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["auser"],
					"condition": "request_user =~'al*?'"
				}
			]
		}
		]
	}
	`,
			request:       adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Attributes: map[string]interface{}{}},
			expectedRoles: []string{"auser"},
		},
		{ //any resource & built-in attr
			stream: `
	{
		"services": [
		{
			"name": "erp",
			"type": "applications",			
			"rolePolicies": [
				{
					"id": "rp1",
					"effect": "grant",
					"roles": ["auser"],
					"condition": "request_user =~'al*?'"
				}
			]
		}
		]
	}
	`,
			request:       adsapi.RequestContext{Subject: &bill, ServiceName: "erp", Attributes: map[string]interface{}{}},
			expectedRoles: []string{},
		},
	}
	for i, tc := range testCases {
		preparePolicyDataInStore([]byte(tc.stream), t)
		eval, err := NewWithStore(conf, testPS)
		if err != nil {
			t.Errorf("case number %d, error creating evaluator : %v", i, err)
			continue
		}
		got, err := eval.GetAllGrantedRoles(tc.request)
		if err != nil {
			t.Errorf("case number %d, Error found in getAllGrantedRoles: %v", i, err)
		}
		if len(got) != len(tc.expectedRoles) {
			fmt.Printf("Case %v.\n", tc)
			t.Errorf("case number %d, expect %d roles, but returned %d roles", i, len(tc.expectedRoles), len(got))
		}
		for _, role := range got {
			roleFound := false
			for _, expectRole := range tc.expectedRoles {
				if expectRole == role {
					roleFound = true
					break
				}
			}
			if !roleFound {
				t.Errorf("case number %d, Role %s is unexpected:", i, role)
			}
		}
	}
}

//empty Permissions means any permission
//empty resource and empty resource expression means any resource
//empty principals means any principal
func TestIsAllowed_empty(t *testing.T) {
	alice := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "alice",
			},
		},
	}
	anyUser := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "anyUser",
			},
		},
	}
	anonymousUser := adsapi.Subject{}
	testCases := []struct {
		stream  string
		request adsapi.RequestContext
		result  bool
	}{
		{ //any permission
			stream: `
				{
			"services": [
			{
				"name": "erp",
				"policies": [
				{
					"id": "grant any permission to principal",
					"effect": "grant",
					"principals": [
						["user:alice"]
					]
				}
				]
			}
			]
		}
				`,
			request: adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Resource: "/any resource", Action: "any action", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //any principal
			stream: `
				{
			"services": [
			{
				"name": "erp",
				"policies": [
				{
					"id": "grant permissions to any principal",
					"effect": "grant",
					"permissions": [
					{
						"resource": "/node1",
						"actions": ["create","delete"]
					}
					]
				}
				]
			}
			]
		}
				`,
			request: adsapi.RequestContext{Subject: &anyUser, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //any principal
			stream: `
				{
			"services": [
			{
				"name": "erp",
				"policies": [
				{
					"id": "grant permissions to any principal",
					"effect": "grant",
					"permissions": [
					{
						"resource": "/node1",
						"actions": ["create","delete"]
					}
					]
				}
				]
			}
			]
		}
				`,
			request: adsapi.RequestContext{Subject: nil, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //any resource
			stream: `
				{
			"services": [
			{
				"name": "erp",
				"policies": [
				{
					"id": "grant permissions to any principal",
					"effect": "grant",
					"permissions": [
					{
						"actions": ["create","delete"]
					}
					],
					"principals": [
						["user:alice"]
					]
				}
				]
			}
			]
		}
				`,
			request: adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Resource: "any resource", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //any action
			stream: `
				{
			"services": [
			{
				"name": "erp",
				"policies": [
				{
					"id": "grant permissions to any principal",
					"effect": "grant",
					"permissions": [
					{
						"resource": "/node1"
					}
					],
					"principals": [
						["user:alice"]
					]
				}
				]
			}
			]
		}
				`,
			request: adsapi.RequestContext{Subject: &alice, ServiceName: "erp", Resource: "/node1", Action: "any action", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //anonymous role
			stream: `
			{
		"services": [
		{
			"name": "erp",
			"policies": [			
			{
				"id": "grant permissions to anonymous role",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["role:anonymous_role"]

				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: nil, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //anonymous role
			stream: `
			{
		"services": [
		{
			"name": "erp",
			"policies": [			
			{
				"id": "grant permissions to anonymous role",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["role:anonymous_role"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &anonymousUser, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //authenticated role
			stream: `
			{
		"services": [
		{
			"name": "erp",
			"policies": [			
			{
				"id": "grant permissions to anonymous role",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["role:authenticated_role"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &anyUser, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //everyone role
			stream: `
			{
		"services": [
		{
			"name": "erp",
			"policies": [			
			{
				"id": "grant permissions to everyone role",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["role:everyone_role"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &anyUser, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
		{ //everyone role
			stream: `
			{
		"services": [
		{
			"name": "erp",
			"policies": [			
			{
				"id": "grant permissions to everyone role",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node1",
					"actions": ["create","delete"]
				}
				],
				"principals": [
					["role:everyone_role"]
				]
			}
			]
		}
		]
	}
			`,
			request: adsapi.RequestContext{Subject: &anonymousUser, ServiceName: "erp", Resource: "/node1", Action: "delete", Attributes: map[string]interface{}{}},
			result:  true,
		},
	}

	for i, tc := range testCases {
		preparePolicyDataInStore([]byte(tc.stream), t)
		eval, err := NewWithStore(conf, testPS)
		if err != nil {
			t.Errorf("case number %d, error creating evaluator : %v", i, err)
			continue
		}
		isAllowed, _, err := eval.IsAllowed(tc.request)
		if err != nil {
			t.Errorf("case number %d, Error found in isAllowed: %v", i, err)
		}
		if isAllowed != tc.result {
			t.Errorf("case number %d, expect %t , but returned %t ", i, tc.result, isAllowed)
		}

	}
}

func TestEntityPrincipal(t *testing.T) {
	mysql := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_ENTITY,
				Name: "spiffe://staging.acme.com/payments/mysql",
			},
		},
	}
	etcd := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_ENTITY,
				Name: "spiffe://staging.acme.com/payments/etcd",
			},
		},
	}
	docker := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_ENTITY,
				Name: "spiffe://staging.acme.com/payments/docker",
			},
		},
	}
	ps := `
				{
		"services": [
		{
			"name": "erp",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["role1"],
				"principals": ["entity:spiffe://staging.acme.com/payments/mysql"]
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
					["role:role1"],
					["entity:spiffe://staging.acme.com/payments/etcd"]
				]
			},
			{
				"id": "policy2",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node2",
					"actions": ["post"]
				}
				],
				"principals": [
					["entity:spiffe://staging.acme.com/payments/docker"]
				]
			}

			]
		}
		]
	}`

	preparePolicyDataInStore([]byte(ps), t)
	eval, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Error("Fail to new evaluator")
	}

	request := adsapi.RequestContext{Subject: &etcd, ServiceName: "erp", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{}}

	isAllowed, reason, err := eval.IsAllowed(request)
	fmt.Println(isAllowed, reason, err)
	if err != nil {
		t.Error("Error happens in evaluator")
	}
	if !isAllowed {
		t.Error("expect allowed, but denied")
	}

	request = adsapi.RequestContext{Subject: &mysql, ServiceName: "erp", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{}}
	isAllowed, reason, err = eval.IsAllowed(request)
	fmt.Println(isAllowed, reason, err)
	if err != nil {
		t.Error("Error happens in evaluator")
	}
	if !isAllowed {
		t.Error("expect allowed, but denied")
	}

	request = adsapi.RequestContext{Subject: &docker, ServiceName: "erp", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{}}
	isAllowed, reason, err = eval.IsAllowed(request)
	fmt.Println(isAllowed, reason, err)
	if err != nil {
		t.Error("Error happens in evaluator")
	}
	if isAllowed {
		t.Error("expect deny, but allow")
	}

	request = adsapi.RequestContext{Subject: &docker, ServiceName: "erp", Resource: "/node2", Action: "post", Attributes: map[string]interface{}{}}
	isAllowed, reason, err = eval.IsAllowed(request)
	fmt.Println(isAllowed, reason, err)
	if err != nil {
		t.Error("Error happens in evaluator")
	}
	if !isAllowed {
		t.Error("expect allowed, but denied")
	}

}
