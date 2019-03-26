//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"fmt"
	"testing"

	"github.com/oracle/speedle/api/pms"
)

func TestPolicyCacheBasic(t *testing.T) {
	cache := NewPolicyCacheData()

	policy1 := &pms.Policy{
		ID:         "policy1",
		Name:       "policy1",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai", "group:Dev"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: "k8s:.*:(dev|qa)/core/pods/.*",
				Actions: []string{
					"*",
				},
			},
		},
	}
	policy2 := &pms.Policy{
		ID:         "policy2",
		Name:       "policy2",
		Effect:     "grant",
		Principals: [][]string{{"user:jianz", "group:Dev"}},
		Permissions: []*pms.Permission{
			{
				Resource: "/tmp/abc/test",
				Actions: []string{
					"*",
				},
			},
		},
	}

	policy3 := &pms.Policy{
		ID:         "policy3",
		Name:       "policy3",
		Effect:     "grant",
		Principals: [][]string{{"user:Bill", "group:Dev"}},
		Permissions: []*pms.Permission{
			{
				Actions: []string{
					"*",
				},
			},
		},
	}

	cache.AddPolicyToCache(policy1, nil)
	cache.AddPolicyToCache(policy2, nil)
	cache.AddPolicyToCache(policy3, nil)

	results := cache.GetRelatedPolicyMap([]string{"group:Dev"}, "", false)
	if results == nil || len(results) != 3 {
		t.Errorf("There should be 3 policies matched in cache.")
	}

	results2 := cache.GetRelatedPolicyMap([]string{"group:Dev"}, "/tmp/abc/test", true)
	if results2 == nil || len(results2) != 2 {
		t.Errorf("There should be 2 policies matched in cache.")
	}

	results3 := cache.GetRelatedPolicyMap([]string{"group:Dev"}, "k8s:.123:dev/core/pods/.abc", true)
	fmt.Println(len(results3))
	if results3 == nil || len(results3) != 2 {
		t.Errorf("There should be 2 policies matched in cache.")
	}

	cache.DeletePolicyFromCache(policy1.ID)
	results4 := cache.GetRelatedPolicyMap([]string{"group:Dev"}, "", false)
	if results4 == nil || len(results4) != 2 {
		t.Errorf("There should be 2 policies matched in cache.")
	}

	cache.DeletePolicyFromCache(policy2.ID)
	results5 := cache.GetRelatedPolicyMap([]string{"group:Dev"}, "/tmp/abc/test", true)
	if results5 == nil || len(results5) != 1 {
		t.Errorf("There should be 1 policies matched in cache.")
	}

	cache.DeletePolicyFromCache(policy3.ID)
}

func TestPolicyCacheAboutResourceExpression(t *testing.T) {
	cache := NewPolicyCacheData()

	policy1 := &pms.Policy{
		ID:         "policy1",
		Name:       "policy1",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: "k8s/core/pods/.*",
				Actions: []string{
					"*",
				},
			},
		},
	}

	policy2 := &pms.Policy{
		ID:         "policy2",
		Name:       "policy2",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: `^k8s/core/pods/.*$`,
				Actions: []string{
					"*",
				},
			},
		},
	}

	policy3 := &pms.Policy{
		ID:         "policy3",
		Name:       "policy3",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: `^.*k8s/core/pods/$`,
				Actions: []string{
					"*",
				},
			},
		},
	}

	policy4 := &pms.Policy{
		ID:         "policy4",
		Name:       "policy4",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: `.*k8s/core/pods/`,
				Actions: []string{
					"*",
				},
			},
		},
	}

	policy5 := &pms.Policy{
		ID:         "policy5",
		Name:       "policy5",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: `.*`,
				Actions: []string{
					"*",
				},
			},
		},
	}

	policy6 := &pms.Policy{
		ID:         "policy6",
		Name:       "policy6",
		Effect:     "grant",
		Principals: [][]string{{"user:william.cai"}},
		Permissions: []*pms.Permission{
			{
				ResourceExpression: `^.*$`,
				Actions: []string{
					"*",
				},
			},
		},
	}

	cache.AddPolicyToCache(policy1, nil)
	cache.AddPolicyToCache(policy2, nil)
	cache.AddPolicyToCache(policy3, nil)
	cache.AddPolicyToCache(policy4, nil)
	cache.AddPolicyToCache(policy5, nil)
	cache.AddPolicyToCache(policy6, nil)

	results1 := cache.GetRelatedPolicyMap([]string{"user:william.cai"}, "k8s/core/pods/resource1", true)
	if results1 == nil || len(results1) != 4 {
		t.Errorf("There should be 4 policies matched in cache, but meet %d", len(results1))
	}

	results2 := cache.GetRelatedPolicyMap([]string{"user:william.cai"}, "k8s/core/pods/resource2", true)
	if results2 == nil || len(results2) != 4 {
		t.Errorf("There should be 4 policies matched in cache, but meet %d", len(results2))
	}

	results3 := cache.GetRelatedPolicyMap([]string{"user:william.cai"}, "abcdef/k8s/core/pods/", true)
	fmt.Println(len(results3))
	if results3 == nil || len(results3) != 4 {
		t.Errorf("There should be 4 policies matched in cache, but meet %d", len(results3))
	}

	results4 := cache.GetRelatedPolicyMap([]string{"user:william.cai"}, "abcdef2/k8s/core/pods/", true)
	fmt.Println(len(results4))
	if results4 == nil || len(results4) != 4 {
		t.Errorf("There should be 4 policies matched in cache, but meet %d", len(results4))
	}

	cache.DeletePolicyFromCache(policy1.ID)
	cache.DeletePolicyFromCache(policy5.ID)
	results5 := cache.GetRelatedPolicyMap([]string{"user:william.cai"}, "k8s/core/pods/resource1", true)
	if results5 == nil || len(results5) != 2 {
		t.Errorf("There should be 2 policies matched in cache, but meet %d", len(results5))
	}

	cache.DeletePolicyFromCache(policy3.ID)
	cache.DeletePolicyFromCache(policy6.ID)
	results6 := cache.GetRelatedPolicyMap([]string{"user:william.cai"}, "abcdef/k8s/core/pods/", true)
	if results6 == nil || len(results6) != 1 {
		t.Errorf("There should be 1 policies matched in cache, but meet %d", len(results6))
	}

	cache.DeletePolicyFromCache(policy2.ID)
	cache.DeletePolicyFromCache(policy4.ID)
	if len(cache.PolicyMap) != 0 ||
		len(cache.PrincipalToPolicies) != 0 ||
		!cache.NilPrincipalToPolicies.isEmpty() {
		t.Errorf("The cache should be empty after delete all the policies")
	}
}
