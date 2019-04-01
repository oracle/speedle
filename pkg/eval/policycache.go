//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/oracle/speedle/3rdparty/github.com/Knetic/govaluate"
	"github.com/oracle/speedle/api/pms"
)

type PolicyCacheData struct {
	BasePolicyCacheData
	PolicyMap map[string]*pms.Policy
}

func NewPolicyCacheData() (p *PolicyCacheData) {
	return &PolicyCacheData{
		BasePolicyCacheData: BasePolicyCacheData{
			PrincipalToPolicies:    make(map[string]*ResourceToPolicyMap),
			NilPrincipalToPolicies: &ResourceToPolicyMap{},
			Conditions:             make(map[string]*govaluate.EvaluableExpression),
		},
		PolicyMap: make(map[string]*pms.Policy),
	}
}

func nilPrincipalPolicy(policy *pms.Policy) (result bool) {
	if policy.Principals == nil || len(policy.Principals) == 0 {
		return true
	}

	return false
}

func (p *PolicyCacheData) AddPolicyToCache(policy *pms.Policy, condition *govaluate.EvaluableExpression) {
	//First add to PolicyMap
	p.PolicyMap[policy.ID] = policy
	if condition != nil {
		p.Conditions[policy.ID] = condition
	}

	//No principal defined. that means the permissions are granted to any principal
	if nilPrincipalPolicy(policy) {
		p.addPolicyToResourceToPolicyMap(p.NilPrincipalToPolicies, policy)
		return
	}

	/*
		In current cache, we don't distinguish andPrincipals and orPrincipals.
		If one principal occured in one policy, we will use this principal as key to index this policy.
		That mean after quiried all related policies, need further match operation to verify the policy
	*/
	for _, andPrincipals := range policy.Principals {
		for _, principal := range andPrincipals {
			if resourceToPolicyMap, exist := p.PrincipalToPolicies[principal]; exist {
				p.addPolicyToResourceToPolicyMap(resourceToPolicyMap, policy)
			} else {
				resourceToPolicyMap := &ResourceToPolicyMap{
					//ResourceToPolicies:           make(map[string]map[string]pms.Element),
					//ResourceExpressionToPolicies: make(map[string]map[string]pms.Element),
					//NilResourceToPolicies:        make(map[string]pms.Element),
				}
				p.PrincipalToPolicies[principal] = resourceToPolicyMap
				p.addPolicyToResourceToPolicyMap(resourceToPolicyMap, policy)
			}
		}
	}
}

func (p *PolicyCacheData) addPolicyToResourceToPolicyMap(resourceToPolicyMap *ResourceToPolicyMap, policy *pms.Policy) {

	//Nil permissions
	if policy.Permissions == nil || len(policy.Permissions) == 0 {
		if resourceToPolicyMap.NilResourceToPolicies == nil {
			resourceToPolicyMap.NilResourceToPolicies = make(map[string]bool)
		}
		resourceToPolicyMap.NilResourceToPolicies[policy.ID] = true
		return
	}

	for _, permission := range policy.Permissions {

		if permission.Resource == "" && permission.ResourceExpression == "" {
			if resourceToPolicyMap.NilResourceToPolicies == nil {
				resourceToPolicyMap.NilResourceToPolicies = make(map[string]bool)
			}
			resourceToPolicyMap.NilResourceToPolicies[policy.ID] = true
		}

		if permission.Resource != "" {
			if resourceToPolicyMap.ResourceToPolicies == nil {
				resourceToPolicyMap.ResourceToPolicies = make(map[string]map[string]bool)
				policyIDSet := make(map[string]bool)
				policyIDSet[policy.ID] = true
				resourceToPolicyMap.ResourceToPolicies[permission.Resource] = policyIDSet
			} else {
				if policyIDSet, exist := resourceToPolicyMap.ResourceToPolicies[permission.Resource]; exist {
					policyIDSet[policy.ID] = true
				} else {
					policyIDSet := make(map[string]bool)
					policyIDSet[policy.ID] = true
					resourceToPolicyMap.ResourceToPolicies[permission.Resource] = policyIDSet
				}
			}
		}

		if permission.ResourceExpression != "" {
			AddPolicyToResourceExpressionCache(resourceToPolicyMap, permission.ResourceExpression, policy.ID)
		}
	}
}

func (p *PolicyCacheData) DeletePolicyFromCache(policyID string) {

	//First delete from PolicyMap
	policy := p.PolicyMap[policyID]
	if policy == nil {
		//This policy does not exist in cache
		return
	}
	delete(p.PolicyMap, policyID)
	if len(policy.Condition) > 0 { //remove related condition cache
		delete(p.Conditions, policyID)
	}

	if nilPrincipalPolicy(policy) {
		p.deletePolicyFromResourceToPolicyMap(p.NilPrincipalToPolicies, policy)
		return
	}

	if p.PrincipalToPolicies != nil {
		for _, andPrincipals := range policy.Principals {
			for _, principal := range andPrincipals {
				if resourceToPolicyMap, exist := p.PrincipalToPolicies[principal]; exist {
					p.deletePolicyFromResourceToPolicyMap(resourceToPolicyMap, policy)
					if resourceToPolicyMap.isEmpty() {
						delete(p.PrincipalToPolicies, principal)
						resourceToPolicyMap = nil
					}
				}
			}
		}
	}
}

func (p *PolicyCacheData) deletePolicyFromResourceToPolicyMap(resourceToPolicyMap *ResourceToPolicyMap, policy *pms.Policy) {

	//Nil permissions policy
	if policy.Permissions == nil || len(policy.Permissions) == 0 {
		if resourceToPolicyMap.NilResourceToPolicies != nil {
			delete(resourceToPolicyMap.NilResourceToPolicies, policy.ID)
		}
		return
	}

	for _, permission := range policy.Permissions {

		if permission.Resource == "" && permission.ResourceExpression == "" &&
			resourceToPolicyMap.NilResourceToPolicies != nil {
			delete(resourceToPolicyMap.NilResourceToPolicies, policy.ID)
		}

		if permission.Resource != "" &&
			resourceToPolicyMap.ResourceToPolicies != nil {
			if policyIDSet, exist := resourceToPolicyMap.ResourceToPolicies[permission.Resource]; exist {
				delete(policyIDSet, policy.ID)
				if len(policyIDSet) == 0 {
					delete(resourceToPolicyMap.ResourceToPolicies, permission.Resource)
				}
			}
		}

		if permission.ResourceExpression != "" {
			DeletePolicyFromResourceExpressionCache(resourceToPolicyMap, permission.ResourceExpression, policy.ID)
		}
	}
}

func (p *PolicyCacheData) GetRelatedPolicyMap(subjectPrincipals []string, resource string, matchResource bool) map[string]*pms.Policy {
	resultPolicyMap := make(map[string]*pms.Policy)

	//First add nil principal policies
	p.getPolicyFromResourceToPolicyMap(p.NilPrincipalToPolicies, resultPolicyMap, resource, matchResource)

	if subjectPrincipals == nil || len(subjectPrincipals) == 0 {
		return resultPolicyMap
	}

	if p.PrincipalToPolicies != nil {
		for _, principal := range subjectPrincipals {
			if resourceToPolicyMap, exist := p.PrincipalToPolicies[principal]; exist {
				p.getPolicyFromResourceToPolicyMap(resourceToPolicyMap, resultPolicyMap, resource, matchResource)
			}
		}
	}

	return resultPolicyMap
}

func (p *PolicyCacheData) getPolicyFromResourceToPolicyMap(resourceToPolicyMap *ResourceToPolicyMap, resultPolicyMap map[string]*pms.Policy, resource string, matchResource bool) {

	//First add nil resource policies
	if resourceToPolicyMap.NilResourceToPolicies != nil {
		for id := range resourceToPolicyMap.NilResourceToPolicies {
			resultPolicyMap[id] = p.PolicyMap[id]
		}
	}

	if matchResource {
		//Check resource policy map
		if resourceToPolicyMap.ResourceToPolicies != nil {
			if policyIDSet, ok := resourceToPolicyMap.ResourceToPolicies[resource]; ok {
				//Add all related policies to result policy map
				for id := range policyIDSet {
					resultPolicyMap[id] = p.PolicyMap[id]
				}
			}
		}

		//Check resource expression policy map
		p.getPoliciesFromResourceExpressionMap(resourceToPolicyMap, resultPolicyMap, resource, true)
	} else {
		//Do not neet match resouce, that mean return all the policies already matched the principal
		if resourceToPolicyMap.ResourceToPolicies != nil {
			for _, policyIDSet := range resourceToPolicyMap.ResourceToPolicies {
				//Add all related policies to result policy map
				for id := range policyIDSet {
					resultPolicyMap[id] = p.PolicyMap[id]
				}
			}
		}

		p.getPoliciesFromResourceExpressionMap(resourceToPolicyMap, resultPolicyMap, resource, false)
	}
}

func (p *PolicyCacheData) getPoliciesFromResourceExpressionMap(resourceToPolicyMap *ResourceToPolicyMap, resultPolicyMap map[string]*pms.Policy, resource string, matchResource bool) {

	if matchResource {

		fn := func(s string, v interface{}) bool {
			policyIDSet := v.(map[string]bool)
			for id := range policyIDSet {
				resultPolicyMap[id] = p.PolicyMap[id]
			}

			return false
		}
		if resourceToPolicyMap.PrefixResourceExpressionTree != nil {
			resourceToPolicyMap.PrefixResourceExpressionTree.WalkPath(resource, fn)
		}

		if resourceToPolicyMap.SuffixResourceExpressionTree != nil {
			resourceToPolicyMap.SuffixResourceExpressionTree.WalkPath(ReverseString(resource), fn)
		}

		if resourceToPolicyMap.ResourceExpressionToPolicies != nil {
			for resExp, policyIDSet := range resourceToPolicyMap.ResourceExpressionToPolicies {
				matched, err := regexp.MatchString(resExp, resource)
				if err != nil {
					log.Errorf("Meet error when match the resource expression in poliy. err: %s", err)
					continue
				}
				if matched {
					//Add all related policies to result policy map
					for id := range policyIDSet {
						resultPolicyMap[id] = p.PolicyMap[id]
					}
				}
			}
		}
	} else {

		fn := func(s string, v interface{}) bool {
			policyIDSet := v.(map[string]bool)
			for id := range policyIDSet {
				resultPolicyMap[id] = p.PolicyMap[id]
			}

			return false
		}
		if resourceToPolicyMap.PrefixResourceExpressionTree != nil {
			resourceToPolicyMap.PrefixResourceExpressionTree.Walk(fn)
		}

		if resourceToPolicyMap.SuffixResourceExpressionTree != nil {
			resourceToPolicyMap.SuffixResourceExpressionTree.Walk(fn)
		}

		if resourceToPolicyMap.ResourceExpressionToPolicies != nil {
			for _, policyIDSet := range resourceToPolicyMap.ResourceExpressionToPolicies {
				//Add all related policies to result policy map
				for id := range policyIDSet {
					resultPolicyMap[id] = p.PolicyMap[id]
				}
			}
		}
	}

}
