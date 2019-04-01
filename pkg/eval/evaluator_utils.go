//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"regexp"
	"strings"

	"github.com/oracle/speedle/3rdparty/github.com/Knetic/govaluate"
	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	log "github.com/sirupsen/logrus"
)

func matchResource(requestRes string, resources, resExpressions []string) bool {
	//in role policy, resources/resExpressions could be empty, which means any resource
	if (resources == nil || len(resources) == 0) && (resExpressions == nil || len(resExpressions) == 0) {
		return true
	}
	for _, res := range resources {
		if requestRes == res {
			return true
		}
	}
	for _, resExp := range resExpressions {
		matched, err := regexp.MatchString(resExp, requestRes)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// Returns if policy is matched
func matchResourceAction(policy *pms.Policy, ctx *internalRequestContext) bool {
	//we interpret nil or empty resource/permission/action/principal etc as ANY resource/permission/action/principal
	if policy.Permissions == nil || len(policy.Permissions) == 0 { //any permissions
		return true
	}
	for _, perm := range policy.Permissions {
		resExpMatch := false
		if len(perm.ResourceExpression) != 0 {
			resExpMatch, _ = regexp.MatchString(perm.ResourceExpression, ctx.Resource)
			//TODO log error
		}
		resNameMatch := perm.Resource == ctx.Resource
		if (len(perm.Resource) == 0 && len(perm.ResourceExpression) == 0) || resExpMatch || resNameMatch {
			if perm.Actions == nil || len(perm.Actions) == 0 { //any action
				return true
			}
			for _, act := range perm.Actions {
				if act == ctx.Action {
					return true
				}
			}
		}

	}
	return false
}

func denyOverwriteCombiner(grantedPolicies []*pms.Policy, deniedPolicies []*pms.Policy,
	context *internalRequestContext, evaluationResult *adsapi.EvaluationResult) (bool, adsapi.Reason) {

	if evaluationResult != nil {
		evaluationResult.AddPolicies(grantedPolicies, deniedPolicies)
	}

	// Evaluate denied policies first
	if len(deniedPolicies) > 0 {
		// If the number of matched denied policies is bigger than 0, then return false directly
		return false, adsapi.DENY_POLICY_FOUND
	}

	if len(grantedPolicies) == 0 {
		// No granted policy defined, return false directly
		// No need to check deny policies
		return false, adsapi.NO_APPLICABLE_POLICIES
	}

	if len(grantedPolicies) > 0 {
		// If the number of matched granted policies is bigger than 0, then return true directly
		return true, adsapi.GRANT_POLICY_FOUND
	}

	//should not go here
	return false, adsapi.REASON_NOT_AVAILABLE
}

func updateSubjectWithBuiltInRoles(s *subject) {
	principals := []string{"role:" + adsapi.BuiltIn_Role_Everyone}
	if s == nil {
		principals = append(principals, "role:"+adsapi.BuiltIn_Role_Anonymous)
	} else {
		if len(s.Users) == 0 && len(s.Groups) == 0 && len(s.Entities) == 0 {
			principals = append(principals, "role:"+adsapi.BuiltIn_Role_Anonymous)
		} else {
			if len(s.Users) != 0 {
				principals = append(principals, s.Users...)
			}
			if len(s.Entities) != 0 {
				principals = append(principals, s.Entities...)
			}

			principals = append(principals, "role:"+adsapi.BuiltIn_Role_Authenticated)
		}

		principals = append(principals, s.Groups...)
	}

	s.Principals = principals
}

func convertRoleToPrincipal(name string) string {
	return "role:" + name
}

func calculatePermissions(grantedPermissions, deniedPermissions []pms.Permission) []pms.Permission {
	if len(deniedPermissions) == 0 {
		return grantedPermissions
	}
	var finalPermissions []pms.Permission
	for _, permission := range grantedPermissions {
		grantPermission := pms.Permission{
			Resource: permission.Resource,
			Actions:  permission.Actions,
		}
		isDenied := false
		for _, deniedPermission := range deniedPermissions {
			expMatched := false
			if len(deniedPermission.ResourceExpression) > 0 {
				matched, err := regexp.MatchString(deniedPermission.ResourceExpression, grantPermission.Resource)
				if err != nil || matched {
					//TODO: log err
					expMatched = true
				}
			}

			if (len(deniedPermission.Resource) == 0 && len(deniedPermission.ResourceExpression) == 0) || expMatched || strings.Compare(deniedPermission.Resource, grantPermission.Resource) == 0 {
				//if resource match, then remove denied actions
				var actions []string
				for _, grantedAction := range grantPermission.Actions {
					actionDenied := false
					for _, deniedAction := range deniedPermission.Actions {
						if grantedAction == deniedAction {
							actionDenied = true
							break
						}
					}
					if !actionDenied {
						actions = append(actions, grantedAction)
					}
				}
				if len(actions) == 0 {
					isDenied = true
					break
				} else {
					grantPermission = pms.Permission{
						Resource: grantPermission.Resource,
						Actions:  actions,
					}
				}
			}
		}
		if !isDenied {
			finalPermissions = append(finalPermissions, grantPermission)
		}
	}
	return finalPermissions

}

func evaluateCondition(condition *govaluate.EvaluableExpression, attributes map[string]interface{}) (bool, error) {
	res, err := condition.Evaluate(attributes)
	if err != nil || res != true {
		if err != nil {
			log.Errorf("Error happens in evaluating condition (%s): %v", condition.String(), err)
		}
		return false, err
	}
	return true, nil
}

func matchRolePolicyPrincipals(subjectPrincipalList []string, rolePolicyPrincipalList []string) bool {
	if subjectPrincipalList == nil || len(subjectPrincipalList) == 0 {
		return false
	}

	if rolePolicyPrincipalList == nil || len(rolePolicyPrincipalList) == 0 {
		return true
	}
	matched := false
	for _, policyPrincipal := range rolePolicyPrincipalList {
		for _, subjectPrincipal := range subjectPrincipalList {
			if policyPrincipal == subjectPrincipal {
				matched = true
				break
			}
		}
		if matched {
			break
		}
	}
	return matched

}

/**
It's regarded as matched only if all items in princs2 are included in princs1
*/
func matchPrincipals(subjectPrincipalList []string, policyPrincipalList [][]string) bool {
	if subjectPrincipalList == nil || len(subjectPrincipalList) == 0 {
		return false
	}

	if policyPrincipalList == nil || len(policyPrincipalList) == 0 {
		return true
	}

	for _, andPrincipals := range policyPrincipalList {
		// one of item in policy principals matched, returns true
		matched := true
		for _, policyPrincipal := range andPrincipals {
			matchedOnePrincipal := false
			// Check if the policy principal in subject principal list
			for _, subjectPrincipal := range subjectPrincipalList {
				if policyPrincipal == subjectPrincipal {
					// matched
					matchedOnePrincipal = true
					break
				}
			}
			// The policy principal is not found in subjectPrincipalList, not match
			if !matchedOnePrincipal {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}

	return false
}
