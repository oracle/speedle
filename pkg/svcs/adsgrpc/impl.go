//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package adsgrpc

import (
	"context"
	"fmt"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/eval"
	"github.com/oracle/speedle/pkg/svcs/adsgrpc/pb"

	"github.com/oracle/speedle/pkg/logging"
)

// GRPCService is the ADS GRPC implementation
type GRPCService struct {
	evaluator eval.InternalEvaluator
}

// NewGRPCService constructs a new ADS GRPC service instance
func NewGRPCService(evaluator eval.InternalEvaluator) (*GRPCService, error) {

	return &GRPCService{
		evaluator: evaluator,
	}, nil
}

func convertGRPCContextRequest(context *pb.ContextRequest) *adsapi.RequestContext {
	ret := adsapi.RequestContext{
		Subject:     convertGRPCSubject(context.Subject),
		ServiceName: context.ServiceName,
		Resource:    context.Resource,
		Action:      context.Action,
	}

	if context.Attributes == nil {
		return &ret
	}
	ret.Attributes = make(map[string]interface{})
	for k, v := range context.Attributes {
		ret.Attributes[k] = v
	}
	return &ret
}

func convertGRPCPrincipals(principals []*pb.Principal) []*adsapi.Principal {
	if principals == nil {
		return nil
	}

	ret := []*adsapi.Principal{}
	for _, princ := range principals {
		ret = append(ret, &adsapi.Principal{
			Type: princ.Type,
			Name: princ.Name,
			IDD:  princ.Idd,
		})
	}
	return ret
}

func convertGRPCSubject(subject *pb.Subject) *adsapi.Subject {
	if subject == nil {
		return nil
	}
	ret := adsapi.Subject{
		Principals: convertGRPCPrincipals(subject.Principals),
		TokenType:  subject.TokenType,
		Token:      subject.Token,
	}
	return &ret
}

func (impl *GRPCService) IsAllowed(ctx context.Context, in *pb.ContextRequest) (*pb.IsAllowedResponse, error) {
	reqCtx := convertGRPCContextRequest(in)

	// assert token
	impl.evaluator.AssertToken(reqCtx)

	allowed, reason, err := impl.evaluator.IsAllowed(*reqCtx)
	if err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]IsAllowed", reqCtx, err.Error())
		return nil, err
	}

	response := pb.IsAllowedResponse{
		Allowed: allowed,
		Reason:  int32(reason),
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]IsAllowed", reqCtx, response)

	return &response, nil
}

func (impl *GRPCService) GetAllGrantedRoles(ctx context.Context, in *pb.ContextRequest) (*pb.AllRoleResponse, error) {
	reqCtx := convertGRPCContextRequest(in)

	// assert token
	impl.evaluator.AssertToken(reqCtx)

	roles, err := impl.evaluator.GetAllGrantedRoles(*reqCtx)
	if err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]GetAllGrantedRoles", reqCtx, err.Error())
		return nil, err
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]GetAllGrantedRoles", reqCtx, roles)

	return &pb.AllRoleResponse{
		Roles: roles,
	}, nil
}

func (impl *GRPCService) GetAllPermissions(ctx context.Context, in *pb.ContextRequest) (*pb.AllPermissionResponse, error) {
	reqCtx := convertGRPCContextRequest(in)

	// assert token
	impl.evaluator.AssertToken(reqCtx)

	perms, err := impl.evaluator.GetAllGrantedPermissions(*reqCtx)
	if err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]GetAllGrantedPermissions", reqCtx, err.Error())
		return nil, err
	}

	ret := pb.AllPermissionResponse{
		Permissions: make([]*pb.AllPermissionResponse_Permission, 0),
	}
	for _, perm := range perms {
		ret.Permissions = append(ret.Permissions, &pb.AllPermissionResponse_Permission{
			Resource: perm.Resource,
			Actions:  perm.Actions,
		})
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]GetAllGrantedPermissions", reqCtx, ret.Permissions)

	return &ret, nil
}

func (impl *GRPCService) Discover(ctx context.Context, in *pb.ContextRequest) (*pb.IsAllowedResponse, error) {
	reqCtx := convertGRPCContextRequest(in)

	// assert token
	impl.evaluator.AssertToken(reqCtx)

	allowed, reason, err := impl.evaluator.Discover(*reqCtx)
	if err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]Discovery", reqCtx, err.Error())
		return nil, err
	}
	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]Discovery", reqCtx, nil)

	return &pb.IsAllowedResponse{
		Allowed: allowed,
		Reason:  int32(reason),
	}, nil

}

func convertAttributes(in map[string]interface{}) map[string]string {
	var out map[string]string
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out

}

func convertToGRPCPrincipals(principals [][]string) []*pb.AndPrincipals {
	ret := []*pb.AndPrincipals{}
	for _, andPrincipals := range principals {
		andPrinc := &pb.AndPrincipals{
			Principals: andPrincipals,
		}

		ret = append(ret, andPrinc)
	}
	return ret
}

func convertAPIPrincipals(principals []*adsapi.Principal) []*pb.Principal {
	if principals == nil {
		return nil
	}

	ret := []*pb.Principal{}
	for _, princ := range principals {
		ret = append(ret, &pb.Principal{
			Type: princ.Type,
			Name: princ.Name,
			Idd:  princ.IDD,
		})
	}
	return ret
}

func convertAPISubject(subject *adsapi.Subject) *pb.Subject {
	if subject == nil {
		return nil
	}
	return &pb.Subject{
		Principals: convertAPIPrincipals(subject.Principals),
		Token:      subject.Token,
		TokenType:  subject.TokenType,
	}
}

func convertAPIRequestContext(req *adsapi.RequestContext) *pb.ContextRequest {
	return &pb.ContextRequest{
		Subject:     convertAPISubject(req.Subject),
		ServiceName: req.ServiceName,
		Resource:    req.Resource,
		Action:      req.Action,
		Attributes:  convertAttributes(req.Attributes),
	}
}

func convertAPIPolicy2PolicyResponse(apiPolicy *pms.Policy, policyResp *pb.Policy) {
	if apiPolicy == nil || policyResp == nil {
		// It shouldn't happen
		return
	}

	retPermission := make([]*pb.Policy_Permission, 0)

	for _, permission := range apiPolicy.Permissions {
		retPermission = append(retPermission, &pb.Policy_Permission{
			Resource:           permission.Resource,
			Actions:            permission.Actions,
			ResourceExpression: permission.ResourceExpression,
		})
	}

	policyResp.ID = apiPolicy.ID
	policyResp.Name = apiPolicy.Name
	policyResp.Effect = apiPolicy.Effect
	policyResp.Permissions = retPermission
	policyResp.Principals = convertToGRPCPrincipals(apiPolicy.Principals)
	policyResp.Condition = apiPolicy.Condition
}

func convertAPIPolicy2EvaluatedPolicyResponse(apiPolicy *adsapi.EvaluatedPolicy, policyResp *pb.EvaluatedPolicy) {
	if apiPolicy == nil || policyResp == nil {
		// It shouldn't happen
		return
	}

	retPermission := make([]*pb.EvaluatedPolicy_Permission, 0)
	for _, permission := range apiPolicy.Permissions {
		retPermission = append(retPermission, &pb.EvaluatedPolicy_Permission{
			Resource:           permission.Resource,
			Actions:            permission.Actions,
			ResourceExpression: permission.ResourceExpression,
		})
	}

	policyResp.Status = apiPolicy.Status
	policyResp.ID = apiPolicy.ID
	policyResp.Name = apiPolicy.Name
	policyResp.Effect = apiPolicy.Effect
	policyResp.Permissions = retPermission
	if apiPolicy.Principals != nil && len(apiPolicy.Principals) > 0 {
		policyResp.Principals = apiPolicy.Principals[0]
	}
	if apiPolicy.Condition != nil {
		policyResp.Condition = &pb.EvaluatedCondition{
			ConditionExpression: apiPolicy.Condition.ConditionExpression,
			EvaluationResult:    apiPolicy.Condition.EvaluationResult,
		}
	}
}

func convertAPIRolePolicy2RolePolicyResponse(apiRolePolicy *pms.RolePolicy, rolePolicyResp *pb.RolePolicy) {
	if apiRolePolicy == nil || rolePolicyResp == nil {
		// It shouldn't happen
		return
	}

	rolePolicyResp.ID = apiRolePolicy.ID
	rolePolicyResp.Name = apiRolePolicy.Name
	rolePolicyResp.Effect = apiRolePolicy.Effect
	rolePolicyResp.Roles = apiRolePolicy.Roles
	rolePolicyResp.Principals = apiRolePolicy.Principals
	rolePolicyResp.Resources = apiRolePolicy.Resources
	rolePolicyResp.ResourceExpressions = apiRolePolicy.ResourceExpressions
	rolePolicyResp.Condition = apiRolePolicy.Condition
}

func convertAPIRolePolicy2EvaluatedRolePolicyResponse(apiRolePolicy *adsapi.EvaluatedRolePolicy, rolePolicyResp *pb.EvaluatedRolePolicy) {
	if apiRolePolicy == nil || rolePolicyResp == nil {
		// It shouldn't happen
		return
	}

	rolePolicyResp.Status = apiRolePolicy.Status
	rolePolicyResp.ID = apiRolePolicy.ID
	rolePolicyResp.Name = apiRolePolicy.Name
	rolePolicyResp.Effect = apiRolePolicy.Effect
	rolePolicyResp.Roles = apiRolePolicy.Roles
	if apiRolePolicy.Principals != nil && len(apiRolePolicy.Principals) > 0 {
		rolePolicyResp.Principals = apiRolePolicy.Principals
	}
	rolePolicyResp.Resources = apiRolePolicy.Resources
	rolePolicyResp.ResourceExpressions = apiRolePolicy.ResourceExpressions
	if apiRolePolicy.Condition != nil {
		rolePolicyResp.Condition = &pb.EvaluatedCondition{
			ConditionExpression: apiRolePolicy.Condition.ConditionExpression,
			EvaluationResult:    apiRolePolicy.Condition.EvaluationResult,
		}
	}
}

func (impl *GRPCService) Diagnose(ctx context.Context, in *pb.ContextRequest) (*pb.EvaluationDebugResponse, error) {
	reqCtx := convertGRPCContextRequest(in)

	// assert token
	impl.evaluator.AssertToken(reqCtx)

	evaResult, err := impl.evaluator.Diagnose(*reqCtx)
	if err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]Diagnose", reqCtx, err.Error())
		return nil, err
	}

	// convert all the role policies
	retRolePolicies := make([]*pb.EvaluatedRolePolicy, 0)
	for _, rolePolicy := range evaResult.RolePolicies {
		var rolePolicyResp pb.EvaluatedRolePolicy
		convertAPIRolePolicy2EvaluatedRolePolicyResponse(rolePolicy, &rolePolicyResp)
		retRolePolicies = append(retRolePolicies, &rolePolicyResp)
	}

	// convert all the policies
	retPolicies := make([]*pb.EvaluatedPolicy, 0)
	for _, policy := range evaResult.Policies {
		var policyResp pb.EvaluatedPolicy
		convertAPIPolicy2EvaluatedPolicyResponse(policy, &policyResp)
		retPolicies = append(retPolicies, &policyResp)
	}

	// Construct & return the response
	response := pb.EvaluationDebugResponse{
		Allowed:        evaResult.Allowed,
		Reason:         evaResult.Reason.String(),
		RequestContext: in,
		GrantedRoles:   evaResult.GrantedRoles,
		RolePolicies:   retRolePolicies,
		Policies:       retPolicies,
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]Diagnose", reqCtx, &response)

	return &response, nil
}
