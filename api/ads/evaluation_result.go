//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package ads

import (
	"strconv"

	"github.com/oracle/speedle/api/pms"
)

func (p *EvaluationResult) AddRolePolicy(rolePolicy *pms.RolePolicy, result bool) {
	var apiEvaluatedRolePolicy EvaluatedRolePolicy
	convertMetaRolePolicy2ApiEvaluatedRolePolicy(rolePolicy, &apiEvaluatedRolePolicy, result)
	p.RolePolicies = append(p.RolePolicies, &apiEvaluatedRolePolicy)
}

func (p *EvaluationResult) AddPolicy(policy *pms.Policy, policyStatus string, result bool) {
	var apiEvaluatedPolicy EvaluatedPolicy
	convertMetaPolicy2ApiEvaluatedPolicy(policy, &apiEvaluatedPolicy, policyStatus, strconv.FormatBool(result))
	p.Policies = append(p.Policies, &apiEvaluatedPolicy)
}

func (p *EvaluationResult) AddPolicies(grantedPolicies []*pms.Policy, deniedPolicies []*pms.Policy) {
	needIgnore := false
	for _, metaPolicy := range deniedPolicies {
		var apiEvaluatedPolicy EvaluatedPolicy
		if needIgnore {
			convertMetaPolicy2ApiEvaluatedPolicy(metaPolicy, &apiEvaluatedPolicy, Evaluation_Ignored, "")
		} else {
			needIgnore = true
			convertMetaPolicy2ApiEvaluatedPolicy(metaPolicy, &apiEvaluatedPolicy, Evaluation_TakeEffect, strconv.FormatBool(true))
		}
		p.Policies = append(p.Policies, &apiEvaluatedPolicy)
	}
	for _, metaPolicy := range grantedPolicies {
		var apiEvaluatedPolicy EvaluatedPolicy
		if needIgnore {
			convertMetaPolicy2ApiEvaluatedPolicy(metaPolicy, &apiEvaluatedPolicy, Evaluation_Ignored, "")
		} else {
			needIgnore = true
			convertMetaPolicy2ApiEvaluatedPolicy(metaPolicy, &apiEvaluatedPolicy, Evaluation_TakeEffect, strconv.FormatBool(true))
		}
		p.Policies = append(p.Policies, &apiEvaluatedPolicy)
	}
}

// 	This function needs to be updated once the "Strategy" is removed from Policy
func convertMetaPolicy2ApiEvaluatedPolicy(metaPolicy *pms.Policy, apiPolicy *EvaluatedPolicy, policyStatus string, evaluationResult string) {
	if metaPolicy == nil || apiPolicy == nil {
		// It shouldn't happen
		return
	}

	var retPermission []pms.Permission
	for _, permission := range metaPolicy.Permissions {
		retPermission = append(retPermission, pms.Permission{
			Resource:           permission.Resource,
			Actions:            permission.Actions,
			ResourceExpression: permission.ResourceExpression,
		})
	}

	apiPolicy.Status = policyStatus
	apiPolicy.ID = metaPolicy.ID
	apiPolicy.Name = metaPolicy.Name
	apiPolicy.Effect = metaPolicy.Effect
	apiPolicy.Permissions = retPermission
	apiPolicy.Principals = metaPolicy.Principals

	if len(metaPolicy.Condition) > 0 {
		apiPolicy.Condition = &EvaluatedCondition{
			ConditionExpression: metaPolicy.Condition,
			EvaluationResult:    evaluationResult,
		}
	}
}

// 	This function needs to be updated once the "Strategy" is removed from RolePolicy
func convertMetaRolePolicy2ApiEvaluatedRolePolicy(metaRolePolicy *pms.RolePolicy, apiRolePolicy *EvaluatedRolePolicy, evaluationResult bool) {
	if metaRolePolicy == nil || apiRolePolicy == nil {
		// It shouldn't happen
		return
	}

	if evaluationResult {
		apiRolePolicy.Status = Evaluation_TakeEffect
	} else {
		apiRolePolicy.Status = Evaluation_ConditionFailed
	}
	apiRolePolicy.ID = metaRolePolicy.ID
	apiRolePolicy.Name = metaRolePolicy.Name
	apiRolePolicy.Effect = metaRolePolicy.Effect
	apiRolePolicy.Roles = metaRolePolicy.Roles
	apiRolePolicy.Principals = metaRolePolicy.Principals
	apiRolePolicy.Resources = metaRolePolicy.Resources
	apiRolePolicy.ResourceExpressions = metaRolePolicy.ResourceExpressions

	if len(metaRolePolicy.Condition) > 0 {
		apiRolePolicy.Condition = &EvaluatedCondition{
			ConditionExpression: metaRolePolicy.Condition,
			EvaluationResult:    strconv.FormatBool(evaluationResult),
		}
	}
}
