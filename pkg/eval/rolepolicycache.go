//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"regexp"

	"github.com/oracle/speedle/api/pms"

	log "github.com/sirupsen/logrus"

	"github.com/oracle/speedle/3rdparty/github.com/Knetic/govaluate"
)

type RolePolicyCacheData struct {
	BasePolicyCacheData
	PolicyMap map[string]*pms.RolePolicy
}

func NewRolePolicyCacheData() (p *RolePolicyCacheData) {
	return &RolePolicyCacheData{
		BasePolicyCacheData: BasePolicyCacheData{
			PrincipalToPolicies:    make(map[string]*ResourceToPolicyMap),
			NilPrincipalToPolicies: &ResourceToPolicyMap{},
			Conditions:             make(map[string]*govaluate.EvaluableExpression),
		},
		PolicyMap: make(map[string]*pms.RolePolicy),
	}
}

func nilPrincipalRolePolicy(rolePolicy *pms.RolePolicy) (result bool) {
	if rolePolicy.Principals == nil || len(rolePolicy.Principals) == 0 {
		return true
	}

	return false
}

func (p *RolePolicyCacheData) AddRolePolicyToCache(policy *pms.RolePolicy, condition *govaluate.EvaluableExpression) {
	//First add role policy to PolicyMap
	p.PolicyMap[policy.ID] = policy
	if condition != nil {
		p.Conditions[policy.ID] = condition
	}

	//No principal defined. that means the roles are granted to any user
	if nilPrincipalRolePolicy(policy) {
		p.addRolePolicyToResourceToRolePolicyMap(p.NilPrincipalToPolicies, policy)
		return
	}

	/*
		In current cache, we don't distinguish andPrincipals and orPrincipals.
		If one principal occured in one policy, we will use this principal as key to index this policy.
		That mean after quiried all related policies, need further match operation to verify the policy
	*/
	for _, principal := range policy.Principals {
		if resourceToRolePolicyMap, exist := p.PrincipalToPolicies[principal]; exist {
			p.addRolePolicyToResourceToRolePolicyMap(resourceToRolePolicyMap, policy)
		} else {
			resourceToRolePolicyMap := &ResourceToPolicyMap{
				//ResourceToPolicies:           make(map[string]map[string]pms.Element),
				//ResourceExpressionToPolicies: make(map[string]map[string]pms.Element),
				//NilResourceToPolicies:        make(map[string]pms.Element),
			}
			p.PrincipalToPolicies[principal] = resourceToRolePolicyMap
			p.addRolePolicyToResourceToRolePolicyMap(resourceToRolePolicyMap, policy)
		}
	}
}

func (p *RolePolicyCacheData) addRolePolicyToResourceToRolePolicyMap(resourceToRolePolicyMap *ResourceToPolicyMap, rolePolicy *pms.RolePolicy) {
	//in role policy, resources/resExpressions could be empty, which means any resource
	if nilResourceRolePolicy(rolePolicy) {
		if resourceToRolePolicyMap.NilResourceToPolicies == nil {
			resourceToRolePolicyMap.NilResourceToPolicies = make(map[string]bool)
		}
		resourceToRolePolicyMap.NilResourceToPolicies[rolePolicy.ID] = true
		return
	}

	for _, resource := range rolePolicy.Resources {

		if resourceToRolePolicyMap.ResourceToPolicies == nil {
			resourceToRolePolicyMap.ResourceToPolicies = make(map[string]map[string]bool)
			policyIDSet := make(map[string]bool)
			policyIDSet[rolePolicy.ID] = true
			resourceToRolePolicyMap.ResourceToPolicies[resource] = policyIDSet
		} else {
			if policyIDSet, exist := resourceToRolePolicyMap.ResourceToPolicies[resource]; exist {
				policyIDSet[rolePolicy.ID] = true
			} else {
				policyIDSet := make(map[string]bool)
				policyIDSet[rolePolicy.ID] = true
				resourceToRolePolicyMap.ResourceToPolicies[resource] = policyIDSet
			}
		}
	}

	for _, resourceExpression := range rolePolicy.ResourceExpressions {
		AddPolicyToResourceExpressionCache(resourceToRolePolicyMap, resourceExpression, rolePolicy.ID)
	}
}

func (p *RolePolicyCacheData) DeleteRolePolicyFromCache(policyID string) {

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

	if nilPrincipalRolePolicy(policy) {
		p.deleteRolePolicyFromResourceToRolePolicyMap(p.NilPrincipalToPolicies, policy)
		return
	}

	if p.PrincipalToPolicies != nil {
		for _, principal := range policy.Principals {
			if resourceToRolePolicyMap, exist := p.PrincipalToPolicies[principal]; exist {
				p.deleteRolePolicyFromResourceToRolePolicyMap(resourceToRolePolicyMap, policy)
				if resourceToRolePolicyMap.isEmpty() {
					delete(p.PrincipalToPolicies, principal)
					resourceToRolePolicyMap = nil
				}
			}
		}
	}
}

func nilResourceRolePolicy(rolePolicy *pms.RolePolicy) (result bool) {
	if (rolePolicy.Resources == nil || len(rolePolicy.Resources) == 0) &&
		(rolePolicy.ResourceExpressions == nil || len(rolePolicy.ResourceExpressions) == 0) {
		return true
	}

	return false
}

func (p *RolePolicyCacheData) deleteRolePolicyFromResourceToRolePolicyMap(resourceToRolePolicyMap *ResourceToPolicyMap, policy *pms.RolePolicy) {
	//in role policy, resources/resExpressions could be empty, which means any resource
	if nilResourceRolePolicy(policy) {
		if resourceToRolePolicyMap.NilResourceToPolicies != nil {
			delete(resourceToRolePolicyMap.NilResourceToPolicies, policy.ID)
		}

		return
	}

	if resourceToRolePolicyMap.ResourceToPolicies != nil {
		for _, resource := range policy.Resources {
			if policyIDSet, exist := resourceToRolePolicyMap.ResourceToPolicies[resource]; exist {
				delete(policyIDSet, policy.ID)
				if len(policyIDSet) == 0 {
					delete(resourceToRolePolicyMap.ResourceToPolicies, resource)
				}
			}
		}
	}

	for _, resourceExpression := range policy.ResourceExpressions {
		DeletePolicyFromResourceExpressionCache(resourceToRolePolicyMap, resourceExpression, policy.ID)
	}
}

func (p *RolePolicyCacheData) GetRelatedRolePolicyMap(subjectPrincipals []string, resource string) map[string]*pms.RolePolicy {
	resultRolePolicyMap := make(map[string]*pms.RolePolicy)

	//First add null princial role policies into result map
	p.getRolePolicyFromResourceToRolePolicyMap(p.NilPrincipalToPolicies, resultRolePolicyMap, resource)

	//Only return the role policy grant to nil principal if subject principals are nil
	if subjectPrincipals == nil || len(subjectPrincipals) == 0 {
		return resultRolePolicyMap
	}

	if p.PrincipalToPolicies != nil {
		for _, principal := range subjectPrincipals {
			if resourceToRolePolicyMap, exist := p.PrincipalToPolicies[principal]; exist {
				p.getRolePolicyFromResourceToRolePolicyMap(resourceToRolePolicyMap, resultRolePolicyMap, resource)
			}
		}
	}

	return resultRolePolicyMap
}

func (p *RolePolicyCacheData) getRolePolicyFromResourceToRolePolicyMap(resourceToRolePolicyMap *ResourceToPolicyMap, resultRolePolicyMap map[string]*pms.RolePolicy, resource string) {

	//Add all of the nil resource policies to result first
	if resourceToRolePolicyMap.NilResourceToPolicies != nil {
		for id := range resourceToRolePolicyMap.NilResourceToPolicies {
			resultRolePolicyMap[id] = p.PolicyMap[id]
		}
	}

	//Check resource policy map
	if resourceToRolePolicyMap.ResourceToPolicies != nil {
		if policyIDSet, ok := resourceToRolePolicyMap.ResourceToPolicies[resource]; ok {
			//Add all related policies to result policy map
			for id := range policyIDSet {
				resultRolePolicyMap[id] = p.PolicyMap[id]
			}
		}
	}

	//Check resource expression
	p.getRolePoliciesFromResourceExpressionMap(resourceToRolePolicyMap, resultRolePolicyMap, resource)
}

func (p *RolePolicyCacheData) getRolePoliciesFromResourceExpressionMap(resourceToRolePolicyMap *ResourceToPolicyMap, resultPolicyMap map[string]*pms.RolePolicy, resource string) {

	fn := func(s string, v interface{}) bool {
		policyIDSet := v.(map[string]bool)
		for id := range policyIDSet {
			resultPolicyMap[id] = p.PolicyMap[id]
		}

		return false
	}
	if resourceToRolePolicyMap.PrefixResourceExpressionTree != nil {
		resourceToRolePolicyMap.PrefixResourceExpressionTree.WalkPath(resource, fn)
	}

	if resourceToRolePolicyMap.SuffixResourceExpressionTree != nil {
		resourceToRolePolicyMap.SuffixResourceExpressionTree.WalkPath(ReverseString(resource), fn)
	}

	if resourceToRolePolicyMap.ResourceExpressionToPolicies != nil {
		for resExp, policyIDSet := range resourceToRolePolicyMap.ResourceExpressionToPolicies {
			matched, err := regexp.MatchString(resExp, resource)
			if err != nil {
				log.Errorf("Meet error when match the resource expression in role poliy. err: %s", err)
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
}
