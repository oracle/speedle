//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"testing"

	"github.com/oracle/speedle/api/pms"
)

func TestRolePolicyCache(t *testing.T) {
	cache := NewRolePolicyCacheData()

	policy1 := &pms.RolePolicy{
		ID:                  "policy1",
		Name:                "policy1",
		Effect:              "grant",
		Roles:               []string{"role1"},
		Principals:          []string{"user:william.cai", "group:Dev"},
		ResourceExpressions: []string{"k8s:.*:(dev|qa)/core/pods/.*"},
	}
	policy2 := &pms.RolePolicy{
		ID:         "policy2",
		Name:       "policy2",
		Effect:     "grant",
		Roles:      []string{"role2"},
		Principals: []string{"user:jianz", "group:Dev"},
		Resources:  []string{"/tmp/abc/test"},
	}

	policy3 := &pms.RolePolicy{
		ID:         "policy3",
		Name:       "policy3",
		Effect:     "grant",
		Roles:      []string{"role3"},
		Principals: []string{"user:Bill", "group:Dev"},
	}

	policy4 := &pms.RolePolicy{
		ID:         "policy4",
		Name:       "policy4",
		Effect:     "grant",
		Roles:      []string{"role4"},
		Principals: []string{"group:Dev"},
	}

	cache.AddRolePolicyToCache(policy1, nil)
	cache.AddRolePolicyToCache(policy2, nil)
	cache.AddRolePolicyToCache(policy3, nil)
	cache.AddRolePolicyToCache(policy4, nil)

	results := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "")
	if results == nil || len(results) != 2 {
		t.Errorf("There should be 2 role policies matched in cache.")
	}

	results2 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "/tmp/abc/test")
	if results2 == nil || len(results2) != 3 {
		t.Errorf("There should be 3 role policies matched in cache.")
	}

	results3 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "k8s:.123:dev/core/pods/.abc")
	if results3 == nil || len(results3) != 3 {
		t.Errorf("There should be 3 role policies matched in cache.")
	}

	cache.DeleteRolePolicyFromCache(policy2.ID)
	results4 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "/tmp/abc/test")
	if results4 == nil || len(results4) != 2 {
		t.Errorf("There should be 2 role policies matched in cache.")
	}

	cache.DeleteRolePolicyFromCache(policy3.ID)
	results5 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "/tmp/abc/test")
	if results5 == nil || len(results5) != 1 {
		t.Errorf("There should be 1 role policies matched in cache.")
	}
}

func TestRolePolicyCacheAboutResourceExpression(t *testing.T) {
	cache := NewRolePolicyCacheData()

	policy1 := &pms.RolePolicy{
		ID:                  "policy1",
		Name:                "policy1",
		Effect:              "grant",
		Roles:               []string{"role1"},
		Principals:          []string{"group:Dev"},
		ResourceExpressions: []string{"k8s/core/pods/.*"},
	}
	policy2 := &pms.RolePolicy{
		ID:                  "policy2",
		Name:                "policy2",
		Effect:              "grant",
		Roles:               []string{"role2"},
		Principals:          []string{"group:Dev"},
		ResourceExpressions: []string{`^k8s/core/pods/.*$`},
	}

	policy3 := &pms.RolePolicy{
		ID:                  "policy3",
		Name:                "policy3",
		Effect:              "grant",
		Roles:               []string{"role3"},
		Principals:          []string{"group:Dev"},
		ResourceExpressions: []string{`.*k8s/core/pods/`},
	}

	policy4 := &pms.RolePolicy{
		ID:                  "policy4",
		Name:                "policy4",
		Effect:              "grant",
		Roles:               []string{"role4"},
		Principals:          []string{"group:Dev"},
		ResourceExpressions: []string{`^.*k8s/core/pods/$`},
	}

	policy5 := &pms.RolePolicy{
		ID:                  "policy5",
		Name:                "policy5",
		Effect:              "grant",
		Roles:               []string{"role5"},
		Principals:          []string{"group:Dev"},
		ResourceExpressions: []string{`.*`},
	}

	policy6 := &pms.RolePolicy{
		ID:                  "policy6",
		Name:                "policy6",
		Effect:              "grant",
		Roles:               []string{"role6"},
		Principals:          []string{"group:Dev"},
		ResourceExpressions: []string{`^.*$`},
	}

	cache.AddRolePolicyToCache(policy1, nil)
	cache.AddRolePolicyToCache(policy2, nil)
	cache.AddRolePolicyToCache(policy3, nil)
	cache.AddRolePolicyToCache(policy4, nil)
	cache.AddRolePolicyToCache(policy5, nil)
	cache.AddRolePolicyToCache(policy6, nil)

	results := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "k8s/core/pods/resource1")
	if results == nil || len(results) != 4 {
		t.Errorf("There should be 4 role policies matched in cache, but meet %d", len(results))
	}

	results2 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "/tmp/abc/test")
	if results2 == nil || len(results2) != 2 {
		t.Errorf("There should be 2 role policies matched in cache, but meet %d", len(results2))
	}

	results3 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "k8s/core/pods/.abc")
	if results3 == nil || len(results3) != 4 {
		t.Errorf("There should be 4 role policies matched in cache, but meet %d", len(results3))
	}

	results4 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "site1/k8s/core/pods/")
	if results4 == nil || len(results4) != 4 {
		t.Errorf("There should be 4 role policies matched in cache, but meet %d", len(results4))
	}

	results5 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "site2/k8s/core/pods/")
	if results5 == nil || len(results5) != 4 {
		t.Errorf("There should be 4 role policies matched in cache, but meet %d", len(results5))
	}

	cache.DeleteRolePolicyFromCache(policy2.ID)
	cache.DeleteRolePolicyFromCache(policy5.ID)
	results6 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "k8s/core/pods/.abc")
	if results6 == nil || len(results6) != 2 {
		t.Errorf("There should be 2 role policies matched in cache, but meet %d", len(results6))
	}

	cache.DeleteRolePolicyFromCache(policy3.ID)
	cache.DeleteRolePolicyFromCache(policy6.ID)
	results7 := cache.GetRelatedRolePolicyMap([]string{"group:Dev"}, "site2/k8s/core/pods/")
	if results7 == nil || len(results7) != 1 {
		t.Errorf("There should be 1 role policies matched in cache, but meet %d", len(results7))
	}

	cache.DeleteRolePolicyFromCache(policy1.ID)
	cache.DeleteRolePolicyFromCache(policy4.ID)
	if len(cache.PolicyMap) != 0 ||
		len(cache.PrincipalToPolicies) != 0 ||
		!cache.NilPrincipalToPolicies.isEmpty() {
		t.Errorf("The cache should be empty after delete all the role policies")
	}
}
