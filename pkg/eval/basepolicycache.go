//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"regexp"
	"strings"

	radix "github.com/armon/go-radix"
	"github.com/oracle/speedle/3rdparty/github.com/Knetic/govaluate"
)

// Patterns used to match resource expression.
// Now support prefix, suffix and all three patterns.
var /*const*/ Prefix_Pattern = regexp.MustCompile(`^\^?[\w/]+\.\*\$?$`)
var /*const*/ Suffix_Pattern = regexp.MustCompile(`^\^?\.\*[\w/]+\$?$`)
var /*const*/ All_Pattern = regexp.MustCompile(`^\^?\.\*\$?$`)

type ResourceToPolicyMap struct {
	//{resource:{policyID: bool}}
	ResourceToPolicies           map[string]map[string]bool
	PrefixResourceExpressionTree *radix.Tree
	SuffixResourceExpressionTree *radix.Tree
	//{resourceExpression:{policyID: bool}}
	//This map contains the resource expressions not match prefix, suffix and all patterns.
	//That mean the incoming resource in isAllowed need match these resource expressions one by one.
	ResourceExpressionToPolicies map[string]map[string]bool
	//resources/resExpressions could be empty, which means any resource
	//Also use this map to store the ".*" resourceexpression policy. which also
	//means any resource.
	NilResourceToPolicies map[string]bool
}

func (p *ResourceToPolicyMap) isEmpty() bool {
	if p.NilResourceToPolicies != nil &&
		len(p.NilResourceToPolicies) > 0 {
		return false
	}

	if p.ResourceExpressionToPolicies != nil &&
		len(p.ResourceExpressionToPolicies) > 0 {
		return false
	}

	if p.ResourceToPolicies != nil &&
		len(p.ResourceToPolicies) > 0 {
		return false
	}

	if p.PrefixResourceExpressionTree != nil &&
		p.PrefixResourceExpressionTree.Len() > 0 {
		return false
	}

	if p.SuffixResourceExpressionTree != nil &&
		p.SuffixResourceExpressionTree.Len() > 0 {
		return false
	}

	return true
}

type BasePolicyCacheData struct {
	/*
		In current cache, we don't distinguish andPrincipals and orPrincipals.
		If one principal occured in one policy, we will use this principal as key to index this policy.
		That mean after quiried all related policies, need further match operation to verify the policy
	*/
	//{principal:ResourceToPolicyMap}
	PrincipalToPolicies map[string]*ResourceToPolicyMap
	//No principal defined in policy, mean match any principal
	NilPrincipalToPolicies *ResourceToPolicyMap
	Conditions             map[string]*govaluate.EvaluableExpression
}

func (p *BasePolicyCacheData) isEmpty() bool {
	if p.PrincipalToPolicies != nil && len(p.PrincipalToPolicies) > 0 {
		return false
	}

	if p.NilPrincipalToPolicies != nil {
		if !p.NilPrincipalToPolicies.isEmpty() {
			return false
		}
	}

	return true
}

func (p *BasePolicyCacheData) clearConditions() {
	p.Conditions = make(map[string]*govaluate.EvaluableExpression)
}

func ReverseString(s string) string {
	bytes := []byte(s)
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return string(bytes)
}

func AddPolicyToResourceExpressionCache(resourceToPolicyMap *ResourceToPolicyMap, resourceExpression string, policyID string) {
	if Prefix_Pattern.MatchString(resourceExpression) {

		resourceExpression := trimResourceExpressionSuffix(resourceExpression)

		if resourceToPolicyMap.PrefixResourceExpressionTree == nil {
			resourceToPolicyMap.PrefixResourceExpressionTree = radix.New()
			policyIDSet := make(map[string]bool)
			policyIDSet[policyID] = true
			resourceToPolicyMap.PrefixResourceExpressionTree.Insert(resourceExpression, policyIDSet)
		} else {
			if value, exist := resourceToPolicyMap.PrefixResourceExpressionTree.Get(resourceExpression); exist {
				policyIDSet := value.(map[string]bool)
				policyIDSet[policyID] = true
			} else {
				policyIDSet := make(map[string]bool)
				policyIDSet[policyID] = true
				resourceToPolicyMap.PrefixResourceExpressionTree.Insert(resourceExpression, policyIDSet)
			}
		}
	} else if Suffix_Pattern.MatchString(resourceExpression) {
		resourceExpression := ReverseString(trimResourceExpressionPrefix(resourceExpression))

		if resourceToPolicyMap.SuffixResourceExpressionTree == nil {
			resourceToPolicyMap.SuffixResourceExpressionTree = radix.New()
			policyIDSet := make(map[string]bool)
			policyIDSet[policyID] = true
			resourceToPolicyMap.SuffixResourceExpressionTree.Insert(resourceExpression, policyIDSet)
		} else {
			if value, exist := resourceToPolicyMap.SuffixResourceExpressionTree.Get(resourceExpression); exist {
				policyIDSet := value.(map[string]bool)
				policyIDSet[policyID] = true
			} else {
				policyIDSet := make(map[string]bool)
				policyIDSet[policyID] = true
				resourceToPolicyMap.SuffixResourceExpressionTree.Insert(resourceExpression, policyIDSet)
			}
		}
	} else if All_Pattern.MatchString(resourceExpression) {
		if resourceToPolicyMap.NilResourceToPolicies == nil {
			resourceToPolicyMap.NilResourceToPolicies = make(map[string]bool)
		}
		resourceToPolicyMap.NilResourceToPolicies[policyID] = true
	} else {
		//No perfix and no suffix and no all pattern matched
		if resourceToPolicyMap.ResourceExpressionToPolicies == nil {
			resourceToPolicyMap.ResourceExpressionToPolicies = make(map[string]map[string]bool)
			policyIDSet := make(map[string]bool)
			policyIDSet[policyID] = true
			resourceToPolicyMap.ResourceExpressionToPolicies[resourceExpression] = policyIDSet
		} else {
			if policyIDSet, exist := resourceToPolicyMap.ResourceExpressionToPolicies[resourceExpression]; exist {
				policyIDSet[policyID] = true
			} else {
				policyIDSet := make(map[string]bool)
				policyIDSet[policyID] = true
				resourceToPolicyMap.ResourceExpressionToPolicies[resourceExpression] = policyIDSet
			}
		}
	}
}

func trimResourceExpressionSuffix(resourceExpression string) string {
	var suffix string
	if strings.HasSuffix(resourceExpression, `.*$`) {
		suffix = `.*$`
	} else {
		suffix = `.*`
	}
	//Also need trim the `^` if it exist
	return strings.TrimPrefix(strings.TrimSuffix(resourceExpression, suffix), `^`)
}

func trimResourceExpressionPrefix(resourceExpression string) string {
	var prefix string
	if strings.HasPrefix(resourceExpression, `^.*`) {
		prefix = `^.*`
	} else {
		prefix = `.*`
	}
	//Also need trim the `$` if it exist
	return strings.TrimSuffix(strings.TrimPrefix(resourceExpression, prefix), `$`)
}

func DeletePolicyFromResourceExpressionCache(resourceToPolicyMap *ResourceToPolicyMap, resourceExpression string, policyID string) {
	if Prefix_Pattern.MatchString(resourceExpression) {
		if resourceToPolicyMap.PrefixResourceExpressionTree == nil {
			return
		}

		resourceExpression := trimResourceExpressionSuffix(resourceExpression)
		if value, exist := resourceToPolicyMap.PrefixResourceExpressionTree.Get(resourceExpression); exist {
			policyIDSet := value.(map[string]bool)
			delete(policyIDSet, policyID)
			if len(policyIDSet) == 0 {
				resourceToPolicyMap.PrefixResourceExpressionTree.Delete(resourceExpression)
			}
		}
	} else if Suffix_Pattern.MatchString(resourceExpression) {
		if resourceToPolicyMap.SuffixResourceExpressionTree == nil {
			return
		}

		resourceExpression := ReverseString(trimResourceExpressionPrefix(resourceExpression))
		if value, exist := resourceToPolicyMap.SuffixResourceExpressionTree.Get(resourceExpression); exist {
			policyIDSet := value.(map[string]bool)
			delete(policyIDSet, policyID)
			if len(policyIDSet) == 0 {
				resourceToPolicyMap.SuffixResourceExpressionTree.Delete(resourceExpression)
			}
		}
	} else if All_Pattern.MatchString(resourceExpression) {
		if resourceToPolicyMap.NilResourceToPolicies != nil {
			delete(resourceToPolicyMap.NilResourceToPolicies, policyID)
		}
	} else {
		//No perfix and no suffix and no all pattern matched
		if resourceToPolicyMap.ResourceExpressionToPolicies == nil {
			return
		}

		if policyIDSet, exist := resourceToPolicyMap.ResourceExpressionToPolicies[resourceExpression]; exist {
			delete(policyIDSet, policyID)
			if len(policyIDSet) == 0 {
				delete(resourceToPolicyMap.ResourceExpressionToPolicies, resourceExpression)
			}
		}
	}
}
