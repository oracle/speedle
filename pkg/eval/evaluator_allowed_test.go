//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

func TestIsAllowedWithoutApp(t *testing.T) {
	const appStream = `
	{
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	//evaluator, err := New(configFile)
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
	_, reason, err := evaluator.IsAllowed(adsapi.RequestContext{
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
	if reason != adsapi.SERVICE_NOT_FOUND {
		t.Fatal("reason should be SERVICE_NOT_FOUND")
		return
	}
}

func TestIsAllowedNotMatchApp(t *testing.T) {
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

	//evaluator, err := New(configFile)
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
	_, reason, err := evaluator.IsAllowed(adsapi.RequestContext{
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
	if reason != adsapi.SERVICE_NOT_FOUND {
		t.Fatal("reason should be SERVICE_NOT_FOUND")
		return
	}
}

func TestIsAllowedNoPolicy(t *testing.T) {
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

	//evaluator, err := New(configFile)
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
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Fatal("reason should be NO_APPLICABLE_POLICIES")
		return
	}
}

func TestIsAllowedNoPrincDef(t *testing.T) {
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

	//evaluator, err := New(configFile)
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
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
	if reason != adsapi.GRANT_POLICY_FOUND {
		t.Fatal("reason should be GRANT_POLICY_FOUND")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "post",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Fatal("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node2",
		Action:      "get",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Fatal("reason should be NO_APPLICABLE_POLICIES")
		return
	}

}

func TestIsAllowedMatch(t *testing.T) {
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
					"user:bill",
					"group:grp1"
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
					["user:william"],
					["user:bill"]
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
					["user:william"],
					["role:role1"]
				]
			},
			{
				"id": "policy3",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node3",
					"actions": ["delete"]
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
					"actions": ["delete"]
				}
				],
				"principals": [
					["group:grp1", "user:cynthia"]
				]
			},
			{
				"id": "policy5",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node5",
					"actions": ["list"]
				}
				],
				"principals": [
					["group:grp1"],
					["role:role2"]
				]
			},
			{
				"id": "policy6",
				"effect": "grant",
				"permissions": [
				{
					"resourceExpression": "/node/res_.*",
					"actions": ["list"]
				}
				],
				"principals": [
					["group:grp1"],
					["role:role2"]
				]
			},
			{
				"id": "policy7",
				"effect": "deny",
				"permissions": [
				{
					"resourceExpression": "/node/res_9999",
					"actions": ["list"]
				}
				],
				"principals": [
					["group:grp1"],
					["role:role2"]
				]
			},
			{
				"id": "policy8",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node8",
					"actions": ["list"]
				}
				],
				"principals": [
					["group:grp10", "user:cynthia"],
					["group:grp2", "user:cynthia"]
				]
			},
			{
				"id": "policy9",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node9",
					"actions": ["list"]
				}
				],
				"principals": [
					["group:grp10", "user:cynthia"],
					["group:grp11", "user:cynthia"]
				]
			}
			]
		}
		]
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	//evaluator, err := New(configFile)
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
				Name: "grp2",
			},
		},
	}

	var allowed bool
	request := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "get",
	}
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should be denied.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node2",
		Action:      "post",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should be denied.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node3",
		Action:      "delete",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
	if reason != adsapi.GRANT_POLICY_FOUND {
		t.Error("reason should be GRANT_POLICY_FOUND")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node4",
		Action:      "delete",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node5",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should be denied.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}
	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node8",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
	if reason != adsapi.GRANT_POLICY_FOUND {
		t.Error("reason should be GRANT_POLICY_FOUND")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node9",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should be denied.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	subject = adsapi.Subject{
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

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node/res_1",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
	if reason != adsapi.GRANT_POLICY_FOUND {
		t.Error("reason should be GRANT_POLICY_FOUND")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node/res_1000",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
	if reason != adsapi.GRANT_POLICY_FOUND {
		t.Error("reason should be GRANT_POLICY_FOUND")
		return
	}

	request = adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "/node/res_9999",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should not be allowed.", request)
		return
	}
	if reason != adsapi.DENY_POLICY_FOUND {
		t.Error("reason should be DENY_POLICY_FOUND")
		return
	}
}

//and principal is not allowed in role policy
func ANDPrincpal_IN_RolePolicy_TestIsAllowedMatch_1(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["roleD"],
				"principals": [
					"user:userA", "group:groupB", "role:roleC"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["roleC"],
				"principals": [
					"user:userA"
				]
			}
			],
			"policies": [
			{
				"id": "policy1",
				"effect": "grant",
				"permissions": [
				{
					"resource": "book",
					"actions": ["read"]
				}
				],
				"principals": [
					["role:roleD"]
				]
			}
			]
		}
		]
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	//evaluator, err := New(configFile)
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "userA",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "groupB",
			},
		},
	}

	var allowed bool
	request := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "book",
		Action:      "read",
	}
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed. reason=%v", request, reason)
		return
	}
}

//and principal is not allowed in role policy
func ANDPrincpal_IN_RolePolicy_TestIsAllowedMatch_2(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["roleC"],
				"principals": [
					"user:userA"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["roleD"],
				"principals": [
					"user:userA", "group:groupB", "role:roleC"
				]
			}
			],
			"policies": [
			{
				"id": "policy1",
				"effect": "grant",
				"permissions": [
				{
					"resource": "book",
					"actions": ["read"]
				}
				],
				"principals": [
					["role:roleD"]
				]
			}
			]
		}
		]
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	//evaluator, err := New(configFile)
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "userA",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "groupB",
			},
		},
	}

	var allowed bool
	request := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "book",
		Action:      "read",
	}
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed. reason=%v", request, reason)
		return
	}
}

func TestIsAllowedMatch_3(t *testing.T) {
	const appStream = `
	{
		"services": [
		{
			"name": "erp",
			"rolePolicies": [
			{
				"id": "rp1",
				"effect": "grant",
				"roles": ["roleC"],
				"principals": [
					"role:roleB"
				]
			},
			{
				"id": "rp2",
				"effect": "grant",
				"roles": ["roleA"],
				"principals": [
					"user:userA"
				]
			},
			{
				"id": "rp3",
				"effect": "grant",
				"roles": ["roleB"],
				"principals": [
					"user:userA"
				]
			}
			],
			"policies": [
			{
				"id": "policy1",
				"effect": "grant",
				"permissions": [
				{
					"resource": "book",
					"actions": ["read"]
				}
				],
				"principals": [
					["role:roleC"]
				]
			}
			]
		}
		]
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	//evaluator, err := New(configFile)
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "userA",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "groupB",
			},
		},
	}

	var allowed bool
	request := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: "erp",
		Resource:    "book",
		Action:      "read",
	}
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed. reason=%v", request, reason)
		return
	}
}

func TestIsAllowedWithIDD(t *testing.T) {
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
					"actions": ["list"]
				}
				],
				"principals": [
					["idd=cisco:user:bill"]
				]
			},
			{
				"id": "policy2",
				"effect": "grant",
				"permissions": [
				{
					"resource": "/node2",
					"actions": ["list"]
				}
				],
				"principals": [
					["user:bill"]
				]
			}
			]
		}
		]
	}
	`
	preparePolicyDataInStore([]byte(appStream), t)

	//evaluator, err := New(configFile)
	evaluator, err := NewWithStore(conf, testPS)
	if err != nil {
		t.Errorf("Unable to initialize evaluator due to error [%v].", err)
		return
	}

	// Evaluate policy 1 without idd, denied is expected
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
		Action:      "list",
	}
	allowed, reason, err := evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should be denied.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	// Evaluate policy 1 with wrong idd, denied is expected
	request = adsapi.RequestContext{
		Subject: &adsapi.Subject{
			Principals: []*adsapi.Principal{
				{
					Type: adsapi.PRINCIPAL_TYPE_USER,
					Name: "bill",
					IDD:  "intel",
				},
			},
		},
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "list",
	}
	allowed, reason, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if allowed {
		t.Fatalf("Request %v should be denied.", request)
		return
	}
	if reason != adsapi.NO_APPLICABLE_POLICIES {
		t.Error("reason should be NO_APPLICABLE_POLICIES")
		return
	}

	// Evaluate policy 1 with correct idd, denied is expected
	request = adsapi.RequestContext{
		Subject: &adsapi.Subject{
			Principals: []*adsapi.Principal{
				{
					Type: adsapi.PRINCIPAL_TYPE_USER,
					Name: "bill",
					IDD:  "cisco",
				},
			},
		},
		ServiceName: "erp",
		Resource:    "/node1",
		Action:      "list",
	}
	allowed, _, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}

	// Evaluate policy 2 with idd, allowed is expected
	request = adsapi.RequestContext{
		Subject: &adsapi.Subject{
			Principals: []*adsapi.Principal{
				{
					Type: adsapi.PRINCIPAL_TYPE_USER,
					Name: "bill",
					IDD:  "cisco",
				},
			},
		},
		ServiceName: "erp",
		Resource:    "/node2",
		Action:      "list",
	}
	allowed, _, err = evaluator.IsAllowed(request)
	if err != nil {
		t.Errorf("Unexcepted error happened [%v].", err)
		return
	}
	if !allowed {
		t.Fatalf("Request %v should be allowed.", request)
		return
	}
}
