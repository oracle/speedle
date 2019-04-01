//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oracle/speedle/3rdparty/github.com/Knetic/govaluate"
	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/eval/function"
	"github.com/oracle/speedle/pkg/subjectutils"

	"github.com/oracle/speedle/api/pms"

	log "github.com/sirupsen/logrus"
)

var builtinFunctions = map[string]govaluate.ExpressionFunction{
	"Sqrt":     function.Sqrt,
	"Max":      function.Max,
	"Min":      function.Min,
	"Sum":      function.Sum,
	"Avg":      function.Avg,
	"IsSubSet": function.IsSubSet,
}

type TokenAsserter interface {
	// set asserter func for policy evaluator
	SetAsserterFunc(f func(ctx *adsapi.RequestContext) error)
	// AssertToken assert token and generate subject to represent the identity
	AssertToken(ctx *adsapi.RequestContext) error
}

type InternalEvaluator interface {
	adsapi.PolicyEvaluator
	TokenAsserter
}

type internalRequestContext struct {
	Subject       *subject
	Service       *RuntimeService
	GlobalService *RuntimeService
	Resource      string
	Action        string
	Attributes    map[string]interface{}
}

type subject struct {
	Users      []string
	Groups     []string
	Entities   []string
	Principals []string
	TokenType  string
	Token      string
}

type PolicyEvalImpl struct {
	RuntimePolicyStore *RuntimePolicyStore //This is runtime policy store
	Store              pms.PolicyStoreManagerADS
	AsserterFunc       func(ctx *adsapi.RequestContext) error
}

func (p *PolicyEvalImpl) deleteService(serviceName string) {
	// Delete application
	p.RuntimePolicyStore.deleteService(serviceName)
}

func (p *PolicyEvalImpl) fullReloadRuntimeCache() {
	ps, err := p.Store.ReadPolicyStore()
	if err != nil {
		log.Errorf("Fail to full reload runtime cache, err:%v", err)
		return
	}
	p.RuntimePolicyStore.reloadPolicyStore(ps)
}

func (p *PolicyEvalImpl) Refresh() error {
	p.fullReloadRuntimeCache()
	return nil
}

func (p *PolicyEvalImpl) AddServiceInRuntimeCache(service *pms.Service) {
	p.RuntimePolicyStore.addService(service)
}

func (p *PolicyEvalImpl) AddPolicyInRuntimeCache(serviceName string, policy *pms.Policy) {
	p.RuntimePolicyStore.addPolicy(serviceName, policy)
}

func (p *PolicyEvalImpl) AddRolePolicyInRuntimeCache(serviceName string, rolepolicy *pms.RolePolicy) {
	p.RuntimePolicyStore.addRolePolicy(serviceName, rolepolicy)
}

func (p *PolicyEvalImpl) DeletePolicyInRuntimeCache(serviceName string, policyID string) {
	p.RuntimePolicyStore.deletePolicy(serviceName, policyID)
}
func (p *PolicyEvalImpl) DeleteRolePolicyInRuntimeCache(serviceName string, rolePolicyID string) {
	p.RuntimePolicyStore.deleteRolePolicy(serviceName, rolePolicyID)
}

func (p *PolicyEvalImpl) DeleteFunctionInRuntimeCache(funcName string) {
	p.RuntimePolicyStore.deleteFunction(funcName)
}

func (p *PolicyEvalImpl) AddFunctionInRuntimeCache(cf *pms.Function) {
	p.RuntimePolicyStore.addFunction(cf)
}

func (p *PolicyEvalImpl) CleanExpiredFunctionResult() {
	p.RuntimePolicyStore.expireFunctionResultCache()
}

func (p *PolicyEvalImpl) SetAsserterFunc(f func(ctx *adsapi.RequestContext) error) {
	p.AsserterFunc = f
}

func (p *PolicyEvalImpl) AssertToken(ctx *adsapi.RequestContext) error {
	// Assert identity token
	if ctx.Subject != nil &&
		p.AsserterFunc != nil &&
		len(ctx.Subject.TokenType) != 0 &&
		len(ctx.Subject.Token) != 0 && !ctx.Subject.Asserted {
		err := p.AsserterFunc(ctx)
		if err == nil {
			ctx.Subject.Asserted = true
		}
		return err
	}

	return nil
}

func (p *PolicyEvalImpl) populateContext(ctx *adsapi.RequestContext) (*internalRequestContext, error) {
	service, err := p.getService(ctx.ServiceName)
	if err != nil {
		return nil, err
	}

	// Assert identity token
	err = p.AssertToken(ctx)
	if err != nil {
		return nil, err
	}

	var globalService *RuntimeService
	if ctx.ServiceName != pms.GlobalService {
		globalService, _ = p.getService(pms.GlobalService)
	}

	newCtx := internalRequestContext{
		Resource:      ctx.Resource,
		Action:        ctx.Action,
		Service:       service,
		GlobalService: globalService,
		Attributes:    make(map[string]interface{}),
	}

	now := time.Now()
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestTime] = now.Unix()
	year, month, day := now.Date()
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestYear] = year
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestMonth] = int(month)
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestDay] = day
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestWeekday] = now.Weekday().String()
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestHour] = now.Hour()

	newCtx.Subject = &subject{
		Users:    []string{},
		Groups:   []string{},
		Entities: []string{},
	}
	if ctx.Subject != nil {
		groups := []interface{}{}
		var user, entity interface{}
		for _, principal := range ctx.Subject.Principals {
			encodedPrincipal := subjectutils.EncodePrincipal(principal)
			principalWithoutIDD := ""
			if len(principal.IDD) != 0 {
				principalWithoutIDD = subjectutils.EncodePrincipal(&adsapi.Principal{
					Type: principal.Type,
					Name: principal.Name,
				})
			}
			switch principal.Type {
			case adsapi.PRINCIPAL_TYPE_USER:
				newCtx.Subject.Users = append(newCtx.Subject.Users, encodedPrincipal)
				if len(principalWithoutIDD) != 0 {
					newCtx.Subject.Users = append(newCtx.Subject.Users, principalWithoutIDD)
				}
				if user == nil {
					user = principal.Name
				}
				break
			case adsapi.PRINCIPAL_TYPE_GROUP:
				newCtx.Subject.Groups = append(newCtx.Subject.Groups, encodedPrincipal)
				groups = append(groups, principal.Name)
				if len(principalWithoutIDD) != 0 {
					newCtx.Subject.Groups = append(newCtx.Subject.Groups, principalWithoutIDD)
				}
				break
			case adsapi.PRINCIPAL_TYPE_ENTITY:
				newCtx.Subject.Entities = append(newCtx.Subject.Entities, encodedPrincipal)
				if len(principalWithoutIDD) != 0 {
					newCtx.Subject.Entities = append(newCtx.Subject.Entities, principalWithoutIDD)
				}
				if entity == nil {
					entity = principal.Name
				}
				break
			}
		}
		if user != nil {
			newCtx.Attributes[adsapi.BuiltIn_Attr_RequestUser] = user
		}
		newCtx.Attributes[adsapi.BuiltIn_Attr_RequestGroups] = groups
		if entity != nil {
			newCtx.Attributes[adsapi.BuiltIn_Attr_RequestEntity] = entity
		}
	}
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestResource] = ctx.Resource
	newCtx.Attributes[adsapi.BuiltIn_Attr_RequestAction] = ctx.Action
	for key, value := range ctx.Attributes {
		newCtx.Attributes[key] = value
	}

	updateSubjectWithBuiltInRoles(newCtx.Subject)

	return &newCtx, nil
}

func (p *PolicyEvalImpl) IsAllowed(ctx adsapi.RequestContext) (bool, adsapi.Reason, error) {
	//IsAllowed don't need return EvaluationResult, so pass nil
	return p.InternalIsAllowed(&ctx, nil)
}

func (p *PolicyEvalImpl) InternalIsAllowed(ctx *adsapi.RequestContext, evaluationResult *adsapi.EvaluationResult) (bool, adsapi.Reason, error) {
	p.RuntimePolicyStore.RLock()
	defer p.RuntimePolicyStore.RUnlock()
	newCtx, err := p.populateContext(ctx)
	if err != nil {
		return false, adsapi.SERVICE_NOT_FOUND, err
	}
	newCtx.Service.RLock()
	defer newCtx.Service.RUnlock()
	if newCtx.Service.PoliciesCache.isEmpty() {
		return false, adsapi.NO_APPLICABLE_POLICIES, nil
	}

	if evaluationResult != nil {
		evaluationResult.Attributes = newCtx.Attributes
	}

	if err := p.resolveSubject(newCtx, evaluationResult); err != nil {
		return false, adsapi.ERROR_IN_EVALUATION, err
	}

	grantedPolicies, deniedPolicies, err := p.getPolicyList(newCtx, true, true, evaluationResult)
	if err != nil {
		return false, adsapi.ERROR_IN_EVALUATION, err
	}

	allowed, reason := denyOverwriteCombiner(grantedPolicies, deniedPolicies, newCtx, evaluationResult)
	return allowed, reason, nil
}

// Return all the policies related to a subject
func (p *PolicyEvalImpl) Diagnose(ctx adsapi.RequestContext) (*adsapi.EvaluationResult, error) {
	// Construct the evaluation result
	retCtx := ctx
	evaResult := adsapi.EvaluationResult{
		Allowed:      false,
		RequestCtx:   &retCtx,
		Attributes:   nil,
		GrantedRoles: make([]string, 0),
		RolePolicies: make([]*adsapi.EvaluatedRolePolicy, 0),
		Policies:     make([]*adsapi.EvaluatedPolicy, 0),
	}
	allowed, reason, err := p.InternalIsAllowed(&ctx, &evaResult)
	evaResult.Allowed = allowed
	evaResult.Reason = reason

	return &evaResult, err
}

func (p *PolicyEvalImpl) GetAllGrantedRoles(ctx adsapi.RequestContext) ([]string, error) {
	p.RuntimePolicyStore.RLock()
	defer p.RuntimePolicyStore.RUnlock()
	newCtx, err := p.populateContext(&ctx)
	if err != nil {
		return nil, err
	}
	newCtx.Service.RLock()
	defer newCtx.Service.RUnlock()

	ret, err := p.getGrantedRolesFromService(newCtx, nil)
	return ret, err
}

//Limitations: This function only calculate granted permissions with resource, will not calculate granted permissions with resource expression.
func (p *PolicyEvalImpl) GetAllGrantedPermissions(ctx adsapi.RequestContext) ([]pms.Permission, error) {
	p.RuntimePolicyStore.RLock()
	defer p.RuntimePolicyStore.RUnlock()
	newCtx, err := p.populateContext(&ctx)
	if err != nil {
		return nil, err
	}

	newCtx.Service.RLock()
	defer newCtx.Service.RUnlock()
	if newCtx.Service.PoliciesCache.isEmpty() {
		return []pms.Permission{}, nil
	}

	if err := p.resolveSubject(newCtx, nil); err != nil {
		return nil, err
	}

	grantedPolicies, deniedPolicies, err := p.getPolicyList(newCtx, false, true, nil)
	if err != nil {
		return nil, err
	}

	var grantedPermissionList, deniedPermissionList []pms.Permission
	for _, policy := range grantedPolicies {
		permissions := policy.Permissions
		if permissions == nil { //means grant any permissions, ignore here
			continue
		}
		for _, permission := range permissions {
			if len(permission.Resource) != 0 {
				grantedPermissionList = append(grantedPermissionList, pms.Permission{
					Resource: permission.Resource,
					Actions:  permission.Actions,
				})
			}
		}
	}
	for _, policy := range deniedPolicies {
		permissions := policy.Permissions
		if permissions == nil { //means deny any permission
			return []pms.Permission{}, nil
		}
		for _, permission := range permissions {
			deniedPermissionList = append(deniedPermissionList, pms.Permission{
				Resource:           permission.Resource,
				Actions:            permission.Actions,
				ResourceExpression: permission.ResourceExpression,
			})
		}
	}

	ret := calculatePermissions(grantedPermissionList, deniedPermissionList)
	return ret, nil
}

func (p *PolicyEvalImpl) getService(serviceName string) (*RuntimeService, error) {
	if runtimeService, exist := p.RuntimePolicyStore.RuntimeServices[serviceName]; exist {
		return runtimeService, nil
	}
	return nil, errors.Errorf(errors.EvalEngineError, "Application %s is not found ", serviceName)
}

func (p *PolicyEvalImpl) resolveSubject(ctx *internalRequestContext, evaluationResult *adsapi.EvaluationResult) error {
	roles, err := p.getGrantedRolesFromService(ctx, evaluationResult)
	if err != nil {
		return err
	}
	for _, role := range roles {
		ctx.Subject.Principals = append(ctx.Subject.Principals, convertRoleToPrincipal(role))
	}

	//Set EvalutionResult
	if evaluationResult != nil {
		evaluationResult.GrantedRoles = roles
	}

	return nil
}

// The firstly returned is granted rolePolicies.
// The second returned value is denied rolePolicies.
func (p *PolicyEvalImpl) getDirectRolePolices(principals []string,
	ctx *internalRequestContext, policyIDMap map[string]bool, evaluationResult *adsapi.EvaluationResult) ([]*pms.RolePolicy, []*pms.RolePolicy, error) {

	grantedRolePolicies := make([]*pms.RolePolicy, 0)
	deniedRolePolicies := make([]*pms.RolePolicy, 0)
	grantedRolePolicies, deniedRolePolicies, err := p.getDirectRolePolicesInService(principals, ctx.Service, ctx.Resource, ctx.Attributes, policyIDMap, evaluationResult, grantedRolePolicies, deniedRolePolicies)
	if err != nil {
		return nil, nil, err
	}
	if ctx.GlobalService != nil {
		grantedRolePolicies, deniedRolePolicies, err = p.getDirectRolePolicesInService(principals, ctx.GlobalService, ctx.Resource, ctx.Attributes, policyIDMap, evaluationResult, grantedRolePolicies, deniedRolePolicies)
		if err != nil {
			return nil, nil, err
		}
	}
	return grantedRolePolicies, deniedRolePolicies, nil
}

func (p *PolicyEvalImpl) getDirectRolePolicesInService(principals []string,
	service *RuntimeService, resource string, attributes map[string]interface{}, policyIDMap map[string]bool, evaluationResult *adsapi.EvaluationResult, grantedRolePolicies []*pms.RolePolicy, deniedRolePolicies []*pms.RolePolicy) ([]*pms.RolePolicy, []*pms.RolePolicy, error) {
	for _, policy := range service.GetRelatedRolePolicyMap(principals, resource) {

		if policyIDMap[policy.ID] {
			continue
		}

		// No principal defined. that means the roles are granted to any user
		if (policy.Principals == nil || len(policy.Principals) == 0 || matchRolePolicyPrincipals(principals, policy.Principals)) && matchResource(resource, policy.Resources, policy.ResourceExpressions) {
			// Evaluate conditions
			condition, ok := service.RolePoliciesCache.Conditions[policy.ID]
			// If no conditions defined, the condition evaluation result is true
			result := true
			if !ok && len(policy.Condition) != 0 {
				if cond, err := p.RuntimePolicyStore.recompileRolePolicyConditionAtRuntime(service.Name, policy); err != nil {
					return nil, nil, err
				} else {
					condition = cond
				}
			}
			if condition != nil {
				result, _ = evaluateCondition(condition, attributes)
			}

			if evaluationResult != nil {
				evaluationResult.AddRolePolicy(policy, result)
			}
			if result {
				switch policy.Effect {
				case pms.Grant:
					grantedRolePolicies = append(grantedRolePolicies, policy)
					break
				case pms.Deny:
					deniedRolePolicies = append(deniedRolePolicies, policy)
					break
				default:
					// TODO: Log a warning, currently do nothing
				}
			}
		}
	}
	return grantedRolePolicies, deniedRolePolicies, nil
}

type Role struct {
	Name               string
	ParentPrincipals   map[string]bool
	ParentRoles        map[string]bool
	ChildRoles         map[string]bool
	DeniedRoles        map[string]bool
	DeniedByRoles      map[string]bool
	DeniedByPrincipals map[string]bool //TODO this could be removed?
}

//assume role policy does not support AND Principal
//assume ctx.Subject.Principals does not contain user defined roles,
//assume built-in role like anonymous role and authenticated role can't be used in role policy
func (p *PolicyEvalImpl) getGrantedRolesFromService(ctx *internalRequestContext, evaluationResult *adsapi.EvaluationResult) ([]string, error) {
	if ctx.GlobalService != nil {
		ctx.GlobalService.RLock()
		defer ctx.GlobalService.RUnlock()
	}

	relatedRolesMap := make(map[string]*Role) //contain all role info related to the role calculation
	policyIDMap := make(map[string]bool)      // this is to avoid repeat processing of same policy.
	grantedRoleMap := make(map[string]bool)   //this is to keep all possiblely granted roles

	subjectPrincipalMap := make(map[string]bool) //this contains non-role principals of ctx.Subject
	for _, principal := range ctx.Subject.Principals {
		if !strings.HasPrefix(principal, "role:") {
			subjectPrincipalMap[principal] = true
		}
	}
	directGrantedRolePolicies, directDeniedRolePolicies, err := p.getDirectRolePolices(ctx.Subject.Principals, ctx, policyIDMap, evaluationResult)
	if err != nil {
		return nil, err
	}

	directDeniedRoleMap := make(map[string]bool) //roles directly denied by user/group in ctx.subject , which always could be denied safely
	deniedRoleMap := make(map[string]bool)       //this contains roles possiblely denied by another role.
	var newlyGrantedRoles []string

	for _, rolePolicy := range directDeniedRolePolicies {
		if _, ok := policyIDMap[rolePolicy.ID]; !ok {
			policyIDMap[rolePolicy.ID] = true
			updateRelatedRoleMapWithDenyRolePolicy(rolePolicy, relatedRolesMap, subjectPrincipalMap, directDeniedRoleMap, grantedRoleMap, deniedRoleMap)
		}
	}

	for _, rolePolicy := range directGrantedRolePolicies {
		if _, ok := policyIDMap[rolePolicy.ID]; !ok {
			policyIDMap[rolePolicy.ID] = true
			newlyGrantedRoles = append(newlyGrantedRoles, updateRelatedRoleMapWithGrantRolePolicy(rolePolicy, relatedRolesMap, subjectPrincipalMap, directDeniedRoleMap, grantedRoleMap)...)
		}
	}

	for len(newlyGrantedRoles) != 0 {
		newSubjectPrincipals := []string{}
		for _, role := range newlyGrantedRoles {
			newSubjectPrincipals = append(newSubjectPrincipals, convertRoleToPrincipal(role))
		}
		indirectGrantedRolePolicies, _, err := p.getDirectRolePolices(newSubjectPrincipals, ctx, policyIDMap, evaluationResult)
		if err != nil {
			return nil, err
		}

		newlyGrantedRoles = []string{}
		for _, rolePolicy := range indirectGrantedRolePolicies {
			if _, ok := policyIDMap[rolePolicy.ID]; !ok {
				policyIDMap[rolePolicy.ID] = true
				newlyGrantedRoles = append(newlyGrantedRoles, updateRelatedRoleMapWithGrantRolePolicy(rolePolicy, relatedRolesMap, subjectPrincipalMap, directDeniedRoleMap, grantedRoleMap)...)
			}
		}
	}

	newSubjectPrincipals := []string{}
	for role := range grantedRoleMap {
		newSubjectPrincipals = append(newSubjectPrincipals, convertRoleToPrincipal(role))
	}
	_, DeniedRolePolicies, err := p.getDirectRolePolices(newSubjectPrincipals, ctx, policyIDMap, evaluationResult)
	if err != nil {
		return nil, err
	}

	for _, rolePolicy := range DeniedRolePolicies {
		if _, ok := policyIDMap[rolePolicy.ID]; !ok {
			policyIDMap[rolePolicy.ID] = true
			updateRelatedRoleMapWithDenyRolePolicy(rolePolicy, relatedRolesMap, subjectPrincipalMap, directDeniedRoleMap, grantedRoleMap, deniedRoleMap)
		}
	}
	//only keep role node that is in possible granted roles
	cleanRelatedRoleMap(relatedRolesMap, grantedRoleMap)

	for {
		//find safely denied role
		safelyDeniedRoles := []string{}
		for deniedRole := range deniedRoleMap {
			if couldRoleSafelyBeDenied(deniedRole, relatedRolesMap, deniedRoleMap) {
				safelyDeniedRoles = append(safelyDeniedRoles, deniedRole)
			}
		}

		if len(safelyDeniedRoles) > 0 {
			for _, deniedRole := range safelyDeniedRoles {
				denyRoleAndDescendants(deniedRole, relatedRolesMap, grantedRoleMap, deniedRoleMap)
				//printRelatedRoleMap(relatedRolesMap)
			}
			//get deniedRoles based on the left role nodes.
			deniedRoleMap = getDeniedRoles(relatedRolesMap, grantedRoleMap)

			continue
		} else {
			//when no safely denied roles found, just delete all denied roles one by one
			for deniedRole := range deniedRoleMap {
				denyRoleAndDescendants(deniedRole, relatedRolesMap, grantedRoleMap, deniedRoleMap)
			}
			break
		}

	}
	finalGrantedRoles := []string{}
	for role := range grantedRoleMap {
		finalGrantedRoles = append(finalGrantedRoles, role)
	}

	return finalGrantedRoles, nil
}

func printRelatedRoleMap(relatedRoleMap map[string]*Role) {
	fmt.Println("----related role map start----")
	for roleName, roleNode := range relatedRoleMap {
		fmt.Printf("%s :\n    parentPrincipals=%v\n    parentRoles=%v\n    deniedByRoles=%v\n    deniedByPrinncipals=%v\n    childRoles=%v\n    denedRoles=%v\n",
			roleName, roleNode.ParentPrincipals, roleNode.ParentRoles, roleNode.DeniedByRoles, roleNode.DeniedByPrincipals, roleNode.ChildRoles, roleNode.DeniedRoles)
	}
	fmt.Println("----related role map end----")
}

func updateRelatedRoleMapWithGrantRolePolicy(rolePolicy *pms.RolePolicy, relatedRolesMap map[string]*Role,
	subjectPrincipalMap map[string]bool, directDeniedRoleMap map[string]bool, grantedRoleMap map[string]bool) []string {

	newlyGrantedRoles := []string{}

	parentRoles := []string{}
	parentPrincipals := []string{}

	for _, principal := range rolePolicy.Principals {
		if strings.HasPrefix(principal, "role:") {
			parentRoles = append(parentRoles, strings.TrimPrefix(principal, "role:"))
		} else {
			if _, ok := subjectPrincipalMap[principal]; ok {
				parentPrincipals = append(parentPrincipals, principal)
			}
		}
	}

	for _, role := range rolePolicy.Roles {
		if _, ok := directDeniedRoleMap[role]; !ok {
			roleNode, ok := relatedRolesMap[role]
			if !ok {
				roleNode = &Role{
					Name:               role,
					ParentPrincipals:   make(map[string]bool),
					ParentRoles:        make(map[string]bool),
					ChildRoles:         make(map[string]bool),
					DeniedRoles:        make(map[string]bool),
					DeniedByRoles:      make(map[string]bool),
					DeniedByPrincipals: make(map[string]bool),
				}
				relatedRolesMap[role] = roleNode
			}
			for _, proleName := range parentRoles {
				roleNode.ParentRoles[proleName] = true
			}
			for _, pprincipalName := range parentPrincipals {
				roleNode.ParentPrincipals[pprincipalName] = true
			}
			if _, ok := grantedRoleMap[role]; !ok {
				grantedRoleMap[role] = true
				newlyGrantedRoles = append(newlyGrantedRoles, role)
			}

		}
	}
	for _, proleName := range parentRoles {
		parentRoleNode, ok := relatedRolesMap[proleName]
		if !ok {
			parentRoleNode = &Role{
				Name:               proleName,
				ParentPrincipals:   make(map[string]bool),
				ParentRoles:        make(map[string]bool),
				ChildRoles:         make(map[string]bool),
				DeniedRoles:        make(map[string]bool),
				DeniedByRoles:      make(map[string]bool),
				DeniedByPrincipals: make(map[string]bool),
			}
			relatedRolesMap[proleName] = parentRoleNode
		}
		for _, childRole := range rolePolicy.Roles {
			if _, ok := directDeniedRoleMap[childRole]; !ok {
				parentRoleNode.ChildRoles[childRole] = true
			}
		}
	}

	return newlyGrantedRoles
}

func updateRelatedRoleMapWithDenyRolePolicy(rolePolicy *pms.RolePolicy, relatedRolesMap map[string]*Role,
	subjectPrincipalMap map[string]bool, directDeniedRoleMap map[string]bool, grantedRoleMap map[string]bool,
	deniedRoleMap map[string]bool) {

	deniedByRoles := []string{}
	deniedByPrincipals := []string{}

	for _, principal := range rolePolicy.Principals {
		if strings.HasPrefix(principal, "role:") {
			deniedByRoles = append(deniedByRoles, strings.TrimPrefix(principal, "role:"))
		} else {
			if _, ok := subjectPrincipalMap[principal]; ok {
				deniedByPrincipals = append(deniedByPrincipals, principal)
			}
		}
	}

	for _, role := range deniedByRoles {
		roleNode, ok := relatedRolesMap[role]
		if !ok {
			roleNode = &Role{
				Name:               role,
				ParentPrincipals:   make(map[string]bool),
				ParentRoles:        make(map[string]bool),
				ChildRoles:         make(map[string]bool),
				DeniedRoles:        make(map[string]bool),
				DeniedByRoles:      make(map[string]bool),
				DeniedByPrincipals: make(map[string]bool),
			}
			relatedRolesMap[role] = roleNode
		}
		for _, deniedRole := range rolePolicy.Roles {
			roleNode.DeniedRoles[deniedRole] = true
		}
	}
	for _, role := range rolePolicy.Roles {
		roleNode, ok := relatedRolesMap[role]
		if !ok {
			roleNode = &Role{
				Name:               role,
				ParentPrincipals:   make(map[string]bool),
				ParentRoles:        make(map[string]bool),
				ChildRoles:         make(map[string]bool),
				DeniedRoles:        make(map[string]bool),
				DeniedByRoles:      make(map[string]bool),
				DeniedByPrincipals: make(map[string]bool),
			}
			relatedRolesMap[role] = roleNode
		}
		for _, deniedByRole := range deniedByRoles {
			roleNode.DeniedByRoles[deniedByRole] = true
		}
		for _, deniedByPrincipal := range deniedByPrincipals {
			roleNode.DeniedByPrincipals[deniedByPrincipal] = true
		}

		if len(deniedByPrincipals) > 0 {
			directDeniedRoleMap[role] = true
		} else {
			if _, ok := grantedRoleMap[role]; ok {
				if _, ok := deniedRoleMap[role]; !ok {
					deniedRoleMap[role] = true
				}
			}
		}
	}
}

func cleanRelatedRoleMap(relatedRoleMap map[string]*Role, grantedRoleMap map[string]bool) {
	for role, roleNode := range relatedRoleMap {
		if _, ok := grantedRoleMap[role]; !ok {
			delete(relatedRoleMap, role)
		} else {
			for pRole := range roleNode.ParentRoles {
				if _, ok := grantedRoleMap[pRole]; !ok {
					delete(roleNode.ParentRoles, pRole)
				}
			}
			for pRole := range roleNode.DeniedByRoles {
				if _, ok := grantedRoleMap[pRole]; !ok {
					delete(roleNode.DeniedByRoles, pRole)
				}
			}
			for cRole := range roleNode.ChildRoles {
				if _, ok := grantedRoleMap[cRole]; !ok {
					delete(roleNode.ChildRoles, cRole)
				}
			}
			for cRole := range roleNode.DeniedRoles {
				if _, ok := grantedRoleMap[cRole]; !ok {
					delete(roleNode.DeniedRoles, cRole)
				}
			}
		}
	}
}

func getDeniedRoles(relatedRoleMap map[string]*Role, grantedRoleMap map[string]bool) map[string]bool {
	deniedRoleMap := make(map[string]bool)
	for role, roleNode := range relatedRoleMap {
		if len(roleNode.DeniedByRoles) > 0 {
			if _, ok := grantedRoleMap[role]; ok {
				deniedRoleMap[role] = true
			}
		}
	}
	return deniedRoleMap
}

func denyRoleAndDescendants(role string, relatedRoleMap map[string]*Role, grantedRoleMap map[string]bool, deniedRoleMap map[string]bool) {
	deletedRoles := make(map[string]bool)
	deletedRoles[role] = true
	descendants := getDeniableDescendantRoles(role, relatedRoleMap)
	//fmt.Printf("deniable descendants for %s are %v \n", role, descendants)
	for _, d := range descendants {
		deletedRoles[d] = true
	}
	for role, roleNode := range relatedRoleMap {
		if _, ok := deletedRoles[role]; !ok {
			for deletedRole := range deletedRoles {
				delete(roleNode.ChildRoles, deletedRole)
				delete(roleNode.ParentRoles, deletedRole)
				delete(roleNode.DeniedRoles, deletedRole)
				delete(roleNode.DeniedByRoles, deletedRole)
				if len(roleNode.ParentRoles) == 0 && len(roleNode.ParentPrincipals) == 0 {
					delete(relatedRoleMap, role)
					delete(grantedRoleMap, role)
				}
			}
		}
	}
	for deletedRole := range deletedRoles {
		delete(relatedRoleMap, deletedRole)
		delete(grantedRoleMap, deletedRole)
	}
	delete(deniedRoleMap, role)
}

func getDeniableDescendantRoles(role string, relatedRoleMap map[string]*Role) []string {
	descendants := []string{}
	if roleNode, ok := relatedRoleMap[role]; ok {
		//get descendant nodes
		for childrole := range roleNode.ChildRoles {
			if childRoleNode, ok := relatedRoleMap[childrole]; ok {
				allParentRolesDenied := true
				for parentRole := range childRoleNode.ParentRoles {
					if parentRole != role && !contains(descendants, parentRole) {
						allParentRolesDenied = false
						break
					}
				}
				if allParentRolesDenied && len(childRoleNode.ParentPrincipals) == 0 {
					descendants = append(descendants, childrole)
					descendants = append(descendants, getDeniableDescendantRoles(childrole, relatedRoleMap)...)
				}
			}
		}
	}
	return descendants
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

//If any of the deniedByRole and all its ancestors are not denied, we take it as could be safely denied.
func couldRoleSafelyBeDenied(role string, relatedRoleMap map[string]*Role, deniedRoleMap map[string]bool) bool {
	allDeniedByRoleBeDenied := true
	if roleNode, ok := relatedRoleMap[role]; ok {
		for deniedByRole := range roleNode.DeniedByRoles {
			if !selfOrAncestorsBeDenied(deniedByRole, relatedRoleMap, deniedRoleMap) {
				allDeniedByRoleBeDenied = false
				break
			}
		}
	} else {
		fmt.Println("error in couldRoleSafelyBeDenied")
	}
	return !allDeniedByRoleBeDenied
}

func selfOrAncestorsBeDenied(deniedByRole string, relatedRoleMap map[string]*Role, deniedRoleMap map[string]bool) bool {
	if _, ok := deniedRoleMap[deniedByRole]; ok { //self be denied
		return true
	}

	if deniedByRoleNode, ok := relatedRoleMap[deniedByRole]; ok {
		if len(deniedByRoleNode.ParentPrincipals) > 0 { //first level role, we could end up here
			return false
		} else {
			for parentRole := range deniedByRoleNode.ParentRoles {
				if !selfOrAncestorsBeDenied(parentRole, relatedRoleMap, deniedRoleMap) {
					return false
				}
			}
		}

	} else {
		fmt.Println("error in selfOrAncestorsBeDenied")
	}
	return true
}

// Returns granted and denied policies
// The first returned value is granted policies
// The second returned value is denied policies
func (p *PolicyEvalImpl) getPolicyList(ctx *internalRequestContext, matchResource bool, matchCondition bool, evaluationResult *adsapi.EvaluationResult) ([]*pms.Policy, []*pms.Policy, error) {
	var grantedPolicyList []*pms.Policy
	var deniedPolicyList []*pms.Policy

	principals := ctx.Subject.Principals
	for _, policy := range ctx.Service.GetRelatedPolicyMap(principals, ctx.Resource, matchResource) {
		// No principal defined. that means the resource actions are granted to any user
		if policy.Principals == nil || len(policy.Principals) == 0 || matchPrincipals(principals, policy.Principals) {
			// Check the resource and action
			if !matchResource || (matchResource && matchResourceAction(policy, ctx)) {
				// Evaluate conditions
				condition, ok := ctx.Service.PoliciesCache.Conditions[policy.ID]
				// If no conditions defined, the condition evaluation result is true
				result := true
				if !ok && len(policy.Condition) != 0 {
					if cond, err := p.RuntimePolicyStore.recompilePolicyConditionAtRuntime(ctx.Service.Name, policy); err != nil {
						return nil, nil, err
					} else {
						condition = cond
					}
				}
				if condition != nil {
					result, _ = evaluateCondition(condition, ctx.Attributes)
				}

				if result {
					switch policy.Effect {
					case pms.Grant:
						grantedPolicyList = append(grantedPolicyList, policy)
						break
					case pms.Deny:
						deniedPolicyList = append(deniedPolicyList, policy)
						break
					default:
						// TODO: Log a warning, currently do nothing.
					}
				} else if evaluationResult != nil {
					//addConditionFailedPolicyToEvaluationResult(policy, result, evaluationResult)
					evaluationResult.AddPolicy(policy, adsapi.Evaluation_ConditionFailed, result)
				}
			}
		}
	}
	return grantedPolicyList, deniedPolicyList, nil
}

//dataSet should be this:
// {
//   [] string, //service name slice
//   map[int]string, //map of policy ID to it's serviceName
//   map[int]string, //map of rolePolicy ID to it's serviceName
//   [] string, //custom function slice
// }
func (p *PolicyEvalImpl) syncRuntimeCache(dataSet []interface{}) error {
	log.Info("start to sync runtime cache data.")
	if len(dataSet) != 4 {
		return errors.New(errors.EvalCacheError, "invalid data set for cache reloading")
	}
	policyStore := p.Store
	log.Info("start to sync service data.")
	if serviceArr, ok := dataSet[0].([]string); ok {
		for _, serviceInDB := range serviceArr { //add new services
			if p.RuntimePolicyStore.RuntimeServices[serviceInDB] == nil {
				service, err := policyStore.GetService(serviceInDB)
				if err != nil {
					return err
				}
				p.AddServiceInRuntimeCache(service)
			}
		}

		for serviceInCache := range p.RuntimePolicyStore.RuntimeServices {
			isInDB := false
			for _, serviceInDB := range serviceArr {
				if serviceInDB == serviceInCache {
					isInDB = true
					break
				}
			}
			if !isInDB {
				p.deleteService(serviceInCache)
			}
		}
	} else {
		return errors.New(errors.EvalCacheError, "invalid service data set for cache reloading")
	}

	log.Info("start to sync policy data.")
	if policySetInDB, ok := dataSet[1].(map[int]string); ok {
		idInDB := getSortedIndexOfMap(policySetInDB)
		policySetInCache, idInCache, err := p.getPolicySetInCache()
		if err != nil {
			return err
		}
		missed, removed := difPolicySets(idInCache, idInDB)
		for _, id := range missed {
			serviceName := policySetInDB[id]
			policy, err := policyStore.GetPolicy(serviceName, strconv.Itoa(id))
			if err != nil {
				return err
			} else {
				p.AddPolicyInRuntimeCache(serviceName, policy)
			}
		}
		for _, id := range removed {
			serviceName := policySetInCache[id]
			p.DeletePolicyInRuntimeCache(serviceName, strconv.Itoa(id))
		}
	} else {
		return errors.New(errors.EvalCacheError, "invalid policy data set for cache reloading")
	}

	log.Info("start to sync rolePolicy data.")
	if rolePolicySetInDB, ok := dataSet[2].(map[int]string); ok {
		idInDB := getSortedIndexOfMap(rolePolicySetInDB)
		rolePolicySetInCache, idInCache, err := p.getRolePolicySetInCache()
		if err != nil {
			return err
		}
		missed, removed := difPolicySets(idInCache, idInDB)
		for _, id := range missed {
			serviceName := rolePolicySetInDB[id]
			policy, err := policyStore.GetRolePolicy(serviceName, strconv.Itoa(id))
			if err != nil {
				return err
			} else {
				p.AddRolePolicyInRuntimeCache(serviceName, policy)
			}
		}
		for _, id := range removed {
			serviceName := rolePolicySetInCache[id]
			p.DeleteRolePolicyInRuntimeCache(serviceName, strconv.Itoa(id))
		}
	} else {
		return errors.New(errors.EvalCacheError, "invalid rolePolicy data set for cache reloading")
	}

	log.Info("start to sync custom functions.")
	if CustFuncSetInDB, ok := dataSet[3].([]string); ok {
		custFuncSetInCache, err := p.getCustFunctionSetInCache()
		if err != nil {
			return err
		}
		missed, removed := difFuncSets(custFuncSetInCache, CustFuncSetInDB)
		for _, name := range missed {
			function, err := policyStore.GetFunction(name)
			if err != nil {
				return err
			} else {
				p.AddFunctionInRuntimeCache(function)
			}
		}
		for _, name := range removed {
			p.DeleteFunctionInRuntimeCache(name)
		}
	} else {
		return errors.New(errors.EvalCacheError, "invalid custom function data set for cache reloading")
	}

	log.Info("finished sync of runtime cache.")
	return nil
}

func getSortedIndexOfMap(m map[int]string) []int {
	var idx []int
	for i := range m {
		idx = append(idx, i)
	}
	sort.Ints(idx)
	return idx
}

func (p *PolicyEvalImpl) getPolicySetInCache() (map[int]string, []int, error) {
	resultSet := make(map[int]string)
	var idx []int
	for _, rService := range p.RuntimePolicyStore.RuntimeServices {
		for pId := range rService.PoliciesCache.PolicyMap {
			id, err := strconv.Atoi(pId)
			if err != nil {
				return nil, nil, errors.Wrapf(err, errors.EvalCacheError, "unable to convert policy ID %q", pId)
			}
			idx = append(idx, id)
			resultSet[id] = rService.Name
		}
	}
	sort.Ints(idx)
	return resultSet, idx, nil
}

func (p *PolicyEvalImpl) getRolePolicySetInCache() (map[int]string, []int, error) {
	resultSet := make(map[int]string)
	var idx []int
	for _, rService := range p.RuntimePolicyStore.RuntimeServices {
		for pId := range rService.RolePoliciesCache.PolicyMap {
			id, err := strconv.Atoi(pId)
			if err != nil {
				return nil, nil, errors.Wrapf(err, errors.EvalCacheError, "unable to convert policy ID %q", pId)
			}
			idx = append(idx, id)
			resultSet[id] = rService.Name
		}
	}
	sort.Ints(idx)
	return resultSet, idx, nil
}

func (p *PolicyEvalImpl) getCustFunctionSetInCache() ([]string, error) {
	resultSet := []string{}
	for funcName := range p.RuntimePolicyStore.Functions {
		if _, ok := builtinFunctions[funcName]; !ok {
			resultSet = append(resultSet, funcName)
		}
	}
	return resultSet, nil
}

func (p *PolicyEvalImpl) updateRuntimeCacheWithStoreChange(updateChan pms.StorageChangeChannel) {
	for e := range updateChan {
		switch e.Type {
		case pms.SERVICE_ADD: ///Event content: StoreUpdateData{ParentID:serviceName, Data:*service}
			serviceGot := e.Content.(*pms.Service)
			p.AddServiceInRuntimeCache(serviceGot)
		case pms.SERVICE_DELETE: //Event content:[]StoreUpdateData{ParentID:serviceName, Data:servieName}
			services := e.Content.([]string)
			for _, s := range services {
				p.deleteService(s)
			}
		case pms.POLICY_ADD: //Event content :[]StoreUpdateData{ParentID:serviceName, Data:*policy}
			data := e.Content.([]pms.StoreUpdateData)
			for _, s := range data {
				policy := s.Data.(*pms.Policy)
				p.AddPolicyInRuntimeCache(s.ServiceName, policy)
			}
		case pms.POLICY_DELETE: // Event content:[]StoreUpdateData{ParentID:serviceName, Data:*pms.Policy}
			data := e.Content.([]pms.StoreUpdateData)
			for _, s := range data {
				policy := s.Data.(*pms.Policy)
				p.DeletePolicyInRuntimeCache(s.ServiceName, policy.ID)
			}
		case pms.ROLEPOLICY_ADD: //Event content :[]StoreUpdateData{ParentID:serviceName, Data:*rolepolicy}
			data := e.Content.([]pms.StoreUpdateData)
			for _, s := range data {
				rolepolicy := s.Data.(*pms.RolePolicy)
				p.AddRolePolicyInRuntimeCache(s.ServiceName, rolepolicy)
			}
		case pms.ROLEPOLICY_DELETE: //Event content:[]StoreUpdateData{ParentID:serviceName, Data:*pms.RolePolicy}
			data := e.Content.([]pms.StoreUpdateData)
			for _, s := range data {
				rolePolicy := s.Data.(*pms.RolePolicy)
				p.DeleteRolePolicyInRuntimeCache(s.ServiceName, rolePolicy.ID)
			}
		case pms.SYNC_RELOAD:
			data := e.Content.([]interface{})
			err := p.syncRuntimeCache(data)
			if err != nil {
				log.Error("failed to reload cache data. ", err)
			}
		case pms.FUNCTION_ADD:
			f := e.Content.(*pms.Function)
			p.AddFunctionInRuntimeCache(f)
		case pms.FUNCTION_DELETE:
			fs := e.Content.([]string)
			for _, f := range fs {
				p.DeleteFunctionInRuntimeCache(f)
			}
		case pms.FULL_RELOAD:
			p.fullReloadRuntimeCache()
		}
	}
}

func (p *PolicyEvalImpl) cleanExpiredFunctionResultPeriodically() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				p.CleanExpiredFunctionResult()
			}
		}
	}()
}

// StopWatch stops watching policy store.
// After stopping watching, policy changes will not be updated automatically
func (p *PolicyEvalImpl) StopWatch() {
	p.Store.StopWatch()
}

func difPolicySets(oldSet, newSet []int) (missed, removed []int) {
	var missedData, removedData []int
	i := 0
	j := 0
	for i <= len(oldSet) && j <= len(newSet) {

		if i == len(oldSet) { //add all data in newSet to missed
			missedData = append(missedData, newSet[j:]...)
			break
		}
		if j == len(newSet) { //add all data in oldSet to removed
			removedData = append(removedData, oldSet[i:]...)
			break
		}

		if oldSet[i] == newSet[j] {
			i++
			j++
		} else {
			removedData = append(removedData, oldSet[i])
			i++
		}
	}
	return missedData, removedData
}

func difFuncSets(oldSet, newSet []string) (missed, removed []string) {
	if len(oldSet) == 0 {
		return newSet, nil
	}
	if len(newSet) == 0 {
		return nil, oldSet
	}
	missed = []string{}
	removed = []string{}
	newMap := make(map[string]bool)
	for _, v := range newSet {
		newMap[v] = true
	}
	oldMap := make(map[string]bool)
	for _, v := range newSet {
		oldMap[v] = true
	}
	for _, v := range newSet {
		if !oldMap[v] {
			missed = append(missed, v)
		}
	}
	for _, v := range oldSet {
		if !newMap[v] {
			removed = append(removed, v)
		}
	}
	return missed, removed

}
