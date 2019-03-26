//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsgrpc

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/store"
	"github.com/oracle/speedle/pkg/svcs/pmsgrpc/pb"
	"github.com/oracle/speedle/pkg/svcs/pmsimpl"

	"context"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"

	"strings"

	"github.com/oracle/speedle/pkg/logging"
)

type serviceImpl struct {
	policyStore pms.PolicyStoreManager
}

// NewServiceImpl initializes a new PMS GRPC instance
func NewServiceImpl(ps pms.PolicyStoreManager) *serviceImpl {
	return &serviceImpl{
		policyStore: ps,
	}
}

func convertRPCFunction(rpcFunction *pb.Function) *pms.Function {
	return &pms.Function{
		Name:           rpcFunction.Name,
		Description:    rpcFunction.Description,
		FuncURL:        rpcFunction.FuncUrl,
		LocalFuncURL:   rpcFunction.LocalFuncUrl,
		CA:             rpcFunction.Ca,
		ResultCachable: rpcFunction.ResultCachable,
		ResultTTL:      rpcFunction.ResultTTL,
	}
}

func convertMetaFunction(function *pms.Function) *pb.Function {
	ret := pb.Function{
		Name:           function.Name,
		Description:    function.Description,
		FuncUrl:        function.FuncURL,
		LocalFuncUrl:   function.LocalFuncURL,
		Ca:             function.CA,
		ResultCachable: function.ResultCachable,
		ResultTTL:      function.ResultTTL,
	}
	return &ret
}

func convertRPCServiceRequest(rpcService *pb.ServiceRequest) *pms.Service {
	ret := pms.Service{
		Name: rpcService.Name,
	}
	switch rpcService.Type {
	case pb.ServiceType_APPLICATION:
		ret.Type = pms.TypeApplication
		break
	case pb.ServiceType_K8S_CLUSTER:
		ret.Type = pms.TypeK8SCluster
		break
	}

	return &ret
}

func convertRPCPrincipals(principals []*pb.AndPrincipals) [][]string {
	ret := [][]string{}
	for _, andPrincipals := range principals {
		ret = append(ret, andPrincipals.Principals)
	}
	return ret
}

func convertRPCRolePolicy(rpcPolicy *pb.RolePolicy) *pms.RolePolicy {
	ret := pms.RolePolicy{
		ID:                  rpcPolicy.Id,
		Name:                rpcPolicy.Name,
		Principals:          rpcPolicy.Principals,
		Roles:               rpcPolicy.Roles,
		Resources:           rpcPolicy.Resources,
		ResourceExpressions: rpcPolicy.ResourceExpressions,
		Condition:           rpcPolicy.Condition,
	}
	switch rpcPolicy.Effect {
	case pb.Effect_GRANT:
		ret.Effect = pms.Grant
		break
	case pb.Effect_DENY:
		ret.Effect = pms.Deny
		break
	}
	return &ret
}

func convertRPCPolicy(rpcPolicy *pb.Policy) *pms.Policy {
	ret := pms.Policy{
		ID:        rpcPolicy.Id,
		Name:      rpcPolicy.Name,
		Condition: rpcPolicy.Condition,
	}
	ret.Principals = convertRPCPrincipals(rpcPolicy.Principals)
	switch rpcPolicy.Effect {
	case pb.Effect_GRANT:
		ret.Effect = pms.Grant
		break
	case pb.Effect_DENY:
		ret.Effect = pms.Deny
		break
	}
	if rpcPolicy.Permissions == nil {
		return &ret
	}

	for _, permission := range rpcPolicy.Permissions {
		ret.Permissions = append(ret.Permissions, convertRPCPermission(permission))
	}
	return &ret
}

func convertRPCPermission(perm *pb.Policy_Permission) *pms.Permission {
	ret := pms.Permission{
		Actions:            perm.Actions,
		Resource:           perm.GetResource(),
		ResourceExpression: perm.GetResourceExpression(),
	}
	return &ret
}

func convertMetaService(service *pms.Service) *pb.Service {
	ret := pb.Service{
		Name: service.Name,
	}
	switch service.Type {
	case pms.TypeApplication:
		ret.Type = pb.ServiceType_APPLICATION
		break
	case pms.TypeK8SCluster:
		ret.Type = pb.ServiceType_K8S_CLUSTER
		break
	}
	if len(service.Policies) > 0 {
		for _, policy := range service.Policies {
			ret.Policies = append(ret.Policies, convertMetaPolicy(policy))
		}
	}
	if len(service.RolePolicies) > 0 {
		for _, rolePolicy := range service.RolePolicies {
			ret.RolePolicies = append(ret.RolePolicies, convertMetaRolePolicy(rolePolicy))
		}
	}

	return &ret
}

func convertMetaPrincipals(principals [][]string) []*pb.AndPrincipals {
	ret := []*pb.AndPrincipals{}
	for _, andPrincipals := range principals {
		ret = append(ret, &pb.AndPrincipals{
			Principals: andPrincipals,
		})
	}
	return ret
}

func convertMetaRolePolicy(policy *pms.RolePolicy) *pb.RolePolicy {
	ret := pb.RolePolicy{
		Id:                  policy.ID,
		Name:                policy.Name,
		Principals:          policy.Principals,
		Roles:               policy.Roles,
		Resources:           policy.Resources,
		ResourceExpressions: policy.ResourceExpressions,
		Condition:           policy.Condition,
	}
	switch policy.Effect {
	case pms.Grant:
		ret.Effect = pb.Effect_GRANT
		break
	case pms.Deny:
		ret.Effect = pb.Effect_DENY
		break
	}
	return &ret
}

func convertMetaPolicy(policy *pms.Policy) *pb.Policy {
	ret := pb.Policy{
		Id:        policy.ID,
		Name:      policy.Name,
		Condition: policy.Condition,
	}
	ret.Principals = convertMetaPrincipals(policy.Principals)
	switch policy.Effect {
	case pms.Grant:
		ret.Effect = pb.Effect_GRANT
		break
	case pms.Deny:
		ret.Effect = pb.Effect_DENY
		break
	}

	if len(policy.Permissions) == 0 {
		return &ret
	}

	for _, permission := range policy.Permissions {
		ret.Permissions = append(ret.Permissions, convertMetaPermission(permission))
	}
	return &ret
}

func convertMetaPermission(perm *pms.Permission) *pb.Policy_Permission {
	ret := pb.Policy_Permission{
		Resource:           perm.Resource,
		ResourceExpression: perm.ResourceExpression,
		Actions:            perm.Actions,
	}
	return &ret
}

func toGRPCStatus(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch errors.Code(err) {
	case errors.StoreError:
		return status.Error(codes.Internal, msg)
	case errors.EntityNotFound:
		return status.Error(codes.NotFound, msg)
	case errors.EntityAlreadyExists:
		return status.Error(codes.AlreadyExists, msg)
	case errors.SerializationError:
		return status.Error(codes.Internal, msg)
	case errors.ExceedLimit:
		return status.Error(codes.ResourceExhausted, msg)
	case errors.InvalidRequest:
		return status.Error(codes.InvalidArgument, msg)
	default:
		return status.Error(codes.Unknown, msg)
	}
}

func (impl *serviceImpl) CreateFunction(ctx context.Context, in *pb.Function) (*pb.Function, error) {
	function := convertRPCFunction(in)
	if function, err := impl.policyStore.CreateFunction(function); err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]CreateFunction", function, err.Error())
		return nil, toGRPCStatus(err)
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]CreateFunction", function, nil)

	return convertMetaFunction(function), nil
}

func (impl *serviceImpl) QueryFunctions(ctx context.Context, in *pb.FunctionQueryRequest) (*pb.FunctionQueryResponse, error) {
	var functions = []*pms.Function{}
	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"name":    in.Name,
		"filters": in.Filters,
	}
	if len(in.Name) == 0 {
		if len(in.Filters) != 0 && strings.HasPrefix(in.Filters, "name") { //Query by name
			functionsMatched, err := impl.policyStore.ListAllFunctions(in.Filters)
			if err != nil {
				// Audit log
				logging.WriteFailedAuditLog("[gRPC]QueryFunctions", ctxFields, err.Error())
				return nil, toGRPCStatus(err)
			}
			functions = functionsMatched
		} else { // Query all functions
			functionsMatched, err := impl.policyStore.ListAllFunctions("")
			if err != nil {
				// Audit log
				logging.WriteFailedAuditLog("[gRPC]QueryFunctions", ctxFields, err.Error())
				return nil, toGRPCStatus(err)
			}
			functions = functionsMatched
		}
	} else {
		function, err := impl.policyStore.GetFunction(in.Name)
		if err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]QueryFunctions", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
		functions = append(functions, function)
	}

	retFunctions := pb.FunctionQueryResponse{
		Functions: make([]*pb.Function, 0),
	}
	for _, f := range functions {
		retFunctions.Functions = append(retFunctions.Functions, convertMetaFunction(f))
	}

	// Audit log
	if len(in.Name) == 0 {
		logging.WriteSucceededAuditLog("[gRPC]QueryFunctions", ctxFields, map[string]interface{}{"functionCount": len(functions)})
	} else {
		logging.WriteSucceededAuditLog("[gRPC]QueryFunctions", ctxFields, map[string]interface{}{"function": functions[0]})
	}

	return &retFunctions, nil
}

func (impl *serviceImpl) DeleteFunctions(ctx context.Context, in *pb.FunctionQueryRequest) (*pb.Empty, error) {
	if len(in.Name) == 0 && len(in.Filters) == 0 {
		return nil, status.Error(codes.InvalidArgument, "both Name and Filters are not passed")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"name":    in.Name,
		"filters": in.Filters,
	}

	//TODO: revisit the query related APIs, currently filter does not work for delete API.
	if len(in.Name) == 0 {
		if err := impl.policyStore.DeleteFunctions(); err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]DeleteFunctions", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
	} else {
		if err := impl.policyStore.DeleteFunction(in.Name); err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]DeleteFunctions", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
	}

	// Audit log
	logging.WriteSucceededAuditLog("[gRPC]DeleteFunctions", ctxFields, nil)

	return &pb.Empty{}, nil
}

func (impl *serviceImpl) CreateService(ctx context.Context, in *pb.ServiceRequest) (*pb.Service, error) {
	service := convertRPCServiceRequest(in)

	err := pmsimpl.CheckService(service, impl.policyStore)
	if err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]CreateService", &service, err.Error())
		return nil, toGRPCStatus(err)
	}

	if err := impl.policyStore.CreateService(service); err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]CreateService", service, err.Error())
		return nil, toGRPCStatus(err)
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]CreateService", service, nil)

	return convertMetaService(service), nil
}

func (impl *serviceImpl) QueryServices(ctx context.Context, in *pb.ServiceQueryRequest) (*pb.ServiceQueryResponse, error) {
	var ss []*pms.Service
	if len(in.Name) == 0 {
		// Get all services
		var err error
		if ss, err = impl.policyStore.ListAllServices(); err != nil {
			// Audit log
			logging.WriteSimpleFailedAuditLog("[gRPC]QueryServices", in.Name, err.Error())
			return nil, toGRPCStatus(err)
		}
	} else {
		svc, err := impl.policyStore.GetService(in.Name)
		if err != nil {
			// Audit log
			logging.WriteSimpleFailedAuditLog("[gRPC]QueryServices", in.Name, err.Error())
			return nil, toGRPCStatus(err)
		}
		ss = append(ss, svc)
	}
	ret := pb.ServiceQueryResponse{
		Services: make([]*pb.Service, 0),
	}

	for _, svc := range ss {
		ret.Services = append(ret.Services, convertMetaService(svc))
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("[gRPC]QueryServices", in.Name, len(ss))

	return &ret, nil
}

func (impl *serviceImpl) DeleteServices(ctx context.Context, in *pb.ServiceQueryRequest) (*pb.Empty, error) {
	if len(in.Name) == 0 {
		if err := impl.policyStore.DeleteServices(); err != nil {
			// Audit log
			logging.WriteSimpleFailedAuditLog("[gRPC]DeleteServices", in.Name, err.Error())
			return nil, toGRPCStatus(err)
		}

		// Audit log
		logging.WriteSimpleSucceededAuditLog("[gRPC]DeleteServices", in.Name, nil)
		return &pb.Empty{}, nil
	}

	if err := impl.policyStore.DeleteService(in.Name); err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]DeleteServices", in.Name, err.Error())
		return nil, toGRPCStatus(err)
	}

	// Audit log for response
	logging.WriteSimpleSucceededAuditLog("[gRPC]DeleteServices", in.Name, nil)
	return &pb.Empty{}, nil
}

func (impl *serviceImpl) CreatePolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.Policy, error) {
	if len(in.ServiceName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "service name is not passed")
	}
	if in.Policy == nil {
		return nil, status.Error(codes.InvalidArgument, "policy is not passed")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName": in.ServiceName,
		"policy":      in.Policy,
	}

	metaPolicy := convertRPCPolicy(in.Policy)

	if err := pmsimpl.CheckPolicy(in.ServiceName, metaPolicy, impl.policyStore); err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]CreatePolicy", ctxFields, err.Error())
		return nil, toGRPCStatus(err)
	}

	retPolicy, err := impl.policyStore.CreatePolicy(in.ServiceName, metaPolicy)
	if err != nil {
		// Audit log
		logging.WriteFailedAuditLog("[gRPC]CreatePolicy", ctxFields, err.Error())
		return nil, toGRPCStatus(err)
	}

	// Audit log
	logging.WriteSucceededAuditLog("[gRPC]CreatePolicy", ctxFields, nil)

	return convertMetaPolicy(retPolicy), nil
}

func (impl *serviceImpl) QueryPolicies(ctx context.Context, in *pb.PolicyQueryRequest) (*pb.PolicyQueryResponse, error) {
	if len(in.ServiceName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "service name is not passed")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName": in.ServiceName,
		"policyId":    in.PolicyID,
	}

	var policies = []*pms.Policy{}
	if len(in.PolicyID) == 0 {
		if len(in.Filters) != 0 && strings.HasPrefix(in.Filters, "name") { //Query by name
			policiesMatched, err := impl.policyStore.ListAllPolicies(in.ServiceName, in.Filters)
			if err != nil {
				// Audit log
				logging.WriteFailedAuditLog("[gRPC]QueryPolicies", ctxFields, err.Error())
				return nil, toGRPCStatus(err)
			}
			policies = policiesMatched
		} else { // Query all policies
			service, err := impl.policyStore.GetService(in.ServiceName)
			if err != nil {
				// Audit log
				logging.WriteFailedAuditLog("[gRPC]QueryPolicies", ctxFields, err.Error())
				return nil, toGRPCStatus(err)
			}
			policies = service.Policies
		}
	} else {
		policy, err := impl.policyStore.GetPolicy(in.ServiceName, in.PolicyID)
		if err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]QueryPolicies", ctxFields, err.Error())
			return nil, err
		}
		policies = append(policies, policy)
	}

	retPolicies := pb.PolicyQueryResponse{
		Policies: make([]*pb.Policy, 0),
	}
	for _, policy := range policies {
		retPolicies.Policies = append(retPolicies.Policies, convertMetaPolicy(policy))
	}

	// Audit log
	if len(in.PolicyID) == 0 {
		logging.WriteSucceededAuditLog("[gRPC]QueryPolicies", ctxFields, map[string]interface{}{"policyCount": len(policies)})
	} else {
		logging.WriteSucceededAuditLog("[gRPC]QueryPolicies", ctxFields, map[string]interface{}{"policy": policies[0]})
	}

	return &retPolicies, nil
}

func (impl *serviceImpl) DeletePolicies(ctx context.Context, in *pb.PolicyQueryRequest) (*pb.Empty, error) {
	if len(in.ServiceName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "service name is not passed")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName": in.ServiceName,
		"policyId":    in.PolicyID,
	}

	if len(in.PolicyID) == 0 {
		if err := impl.policyStore.DeletePolicies(in.ServiceName); err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]DeletePolicies", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
	} else {
		if err := impl.policyStore.DeletePolicy(in.ServiceName, in.PolicyID); err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]DeletePolicies", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
	}

	// Audit log
	logging.WriteSucceededAuditLog("[gRPC]DeletePolicies", ctxFields, nil)

	return &pb.Empty{}, nil
}

func (impl *serviceImpl) CreateRolePolicy(ctx context.Context, in *pb.RolePolicyRequest) (*pb.RolePolicy, error) {
	if len(in.ServiceName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "service name is not passed")
	}
	if in.RolePolicy == nil {
		return nil, status.Error(codes.InvalidArgument, "Policy is not passed")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName": in.ServiceName,
		"rolePolicy":  in.RolePolicy,
	}

	metaRolePolicy := convertRPCRolePolicy(in.RolePolicy)

	if err := pmsimpl.CheckRolePolicy(in.ServiceName, metaRolePolicy, impl.policyStore); err != nil {
		// Audit log
		logging.WriteSimpleFailedAuditLog("[gRPC]CreateRolePolicy", ctxFields, err.Error())
		return nil, toGRPCStatus(err)
	}

	retPolicy, err := impl.policyStore.CreateRolePolicy(in.ServiceName, metaRolePolicy)
	if err != nil {
		// Audit log
		logging.WriteFailedAuditLog("[gRPC]CreateRolePolicy", ctxFields, err.Error())
		return nil, toGRPCStatus(err)
	}

	// Audit log
	logging.WriteSucceededAuditLog("[gRPC]CreateRolePolicy", ctxFields, nil)

	return convertMetaRolePolicy(retPolicy), nil
}

func (impl *serviceImpl) QueryRolePolicies(ctx context.Context, in *pb.RolePolicyQueryRequest) (*pb.RolePolicyQueryResponse, error) {
	if len(in.ServiceName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "service name is not passed.")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName":  in.ServiceName,
		"rolePolicyId": in.RolePolicyID,
	}

	var policies = []*pms.RolePolicy{}
	if len(in.RolePolicyID) == 0 {
		if len(in.Filters) != 0 && strings.HasPrefix(in.Filters, "name") { //Query by name
			policiesMatched, err := impl.policyStore.ListAllRolePolicies(in.ServiceName, in.Filters)
			if err != nil {
				// Audit log
				logging.WriteFailedAuditLog("[gRPC]QueryRolePolicies", ctxFields, err.Error())
				return nil, toGRPCStatus(err)
			}
			policies = policiesMatched
		} else { // Query all policies
			service, err := impl.policyStore.GetService(in.ServiceName)
			if err != nil {
				// Audit log
				logging.WriteFailedAuditLog("[gRPC]QueryRolePolicies", ctxFields, err.Error())
				return nil, toGRPCStatus(err)
			}
			policies = service.RolePolicies
		}
		// Audit log
		logging.WriteSucceededAuditLog("[gRPC]QueryRolePolicies", ctxFields, map[string]interface{}{"rolePolicyCount": len(policies)})
	} else {
		policy, err := impl.policyStore.GetRolePolicy(in.ServiceName, in.RolePolicyID)
		if err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]QueryRolePolicies", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
		policies = append(policies, policy)

		// Audit log
		logging.WriteSucceededAuditLog("[gRPC]QueryRolePolicies", ctxFields, map[string]interface{}{"rolePolicy": policy})
	}

	retPolicies := pb.RolePolicyQueryResponse{
		RolePolicies: make([]*pb.RolePolicy, 0),
	}
	for _, policy := range policies {
		retPolicies.RolePolicies = append(retPolicies.RolePolicies, convertMetaRolePolicy(policy))
	}

	return &retPolicies, nil
}

func (impl *serviceImpl) DeleteRolePolicies(ctx context.Context, in *pb.RolePolicyQueryRequest) (*pb.Empty, error) {
	if len(in.ServiceName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "service name is not passed.")
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName":  in.ServiceName,
		"rolePolicyId": in.RolePolicyID,
	}

	if len(in.RolePolicyID) == 0 {
		if err := impl.policyStore.DeleteRolePolicies(in.ServiceName); err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]DeleteRolePolicies", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
	} else {
		if err := impl.policyStore.DeleteRolePolicy(in.ServiceName, in.RolePolicyID); err != nil {
			// Audit log
			logging.WriteFailedAuditLog("[gRPC]DeleteRolePolicies", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
	}

	// Audit log
	logging.WriteSucceededAuditLog("[gRPC]DeleteRolePolicies", ctxFields, nil)

	return &pb.Empty{}, nil
}

func (impl *serviceImpl) ListPolicyCounts(ctx context.Context, in *pb.Empty) (*pb.PolicyCountsMap, error) {
	countsMap, err := impl.policyStore.GetPolicyAndRolePolicyCounts()
	if err != nil {
		// Audit log
		logging.WriteFailedAuditLog("[gRPC]ListPolicyCounts", nil, err.Error())
		return nil, toGRPCStatus(err)
	}

	retCountsMap := pb.PolicyCountsMap{
		CountMap: make(map[string]*pb.PolicyAndRolePolicyCounts, 0),
	}

	for k, v := range countsMap {
		retCountsMap.CountMap[k] = &pb.PolicyAndRolePolicyCounts{
			PolicyCount:     v.PolicyCount,
			RolePolicyCount: v.RolePolicyCount,
		}
	}

	return &retCountsMap, nil
}

func (impl *serviceImpl) GetDiscoverRequests(ctx context.Context, in *pb.DiscoverRequestsRequest) (*pb.DiscoverRequestsResponse, error) {
	discoverRequestMgr, _ := impl.policyStore.(store.DiscoverRequestManager)
	last := in.Last
	revision := in.Revision
	serviceName := in.ServiceName
	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serverName": serviceName,
		"last":       last,
		"revision":   revision,
	}

	requests := []*pb.ContextRequest{}
	if last {
		req, revision, err := discoverRequestMgr.GetLastDiscoverRequest(serviceName)
		if err != nil {
			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
		requests = append(requests, convertAPIRequestContext(req))

		// Audit log
		logging.WriteSucceededAuditLog("GetDiscoverRequests", ctxFields, map[string]interface{}{"lastRequest": req})

		return &pb.DiscoverRequestsResponse{Requests: requests, Revision: revision}, nil
	} else if revision > 0 {
		reqs, revision, err := discoverRequestMgr.GetDiscoverRequestsSinceRevision(serviceName, revision)
		if err != nil {
			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
		for _, req := range reqs {
			requests = append(requests, convertAPIRequestContext(req))
		}

		// Audit log
		logging.WriteSucceededAuditLog("GetDiscoverRequests", ctxFields, map[string]interface{}{"requestCount": len(requests)})

		return &pb.DiscoverRequestsResponse{Requests: requests, Revision: revision}, nil
	} else {
		reqs, revision, err := discoverRequestMgr.GetDiscoverRequests(serviceName)
		if err != nil {
			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return nil, toGRPCStatus(err)
		}
		for _, req := range reqs {
			requests = append(requests, convertAPIRequestContext(req))
		}

		// Audit log
		logging.WriteSucceededAuditLog("GetDiscoverRequests", ctxFields, map[string]interface{}{"requestCount": len(requests)})

		return &pb.DiscoverRequestsResponse{Requests: requests, Revision: revision}, nil
	}

}
func (impl *serviceImpl) ResetDiscoverRequests(ctx context.Context, in *pb.ResetRequestsRequest) (*pb.ResetRequestsResponse, error) {
	discoverRequestMgr, _ := impl.policyStore.(store.DiscoverRequestManager)
	err := discoverRequestMgr.ResetDiscoverRequests(in.ServiceName)

	// Audit log
	if err != nil {
		logging.WriteSimpleFailedAuditLog("ResetDiscoverRequests", in.ServiceName, err.Error())
	} else {
		logging.WriteSimpleSucceededAuditLog("ResetDiscoverRequests", in.ServiceName, nil)
	}

	return &pb.ResetRequestsResponse{}, toGRPCStatus(err)
}

func (impl *serviceImpl) GetDiscoverPolicies(ctx context.Context, in *pb.DiscoverPoliciesRequest) (*pb.DiscoverPoliciesResponse, error) {
	discoverRequestMgr, _ := impl.policyStore.(store.DiscoverRequestManager)
	serviceMap, revision, err := discoverRequestMgr.GeneratePolicies(in.ServiceName, in.PrincipalType, in.PrincipalName, in.PrincipalIdd)

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serverName":    in.ServiceName,
		"principalType": in.PrincipalType,
		"principalName": in.PrincipalName,
		"principalIdd":  in.PrincipalIdd,
	}

	// Audit log
	if err != nil {
		logging.WriteFailedAuditLog("GetDiscoverPolicies", ctxFields, err.Error())
		return nil, toGRPCStatus(err)
	}
	services := []*pb.Service{}
	for _, service := range serviceMap {
		services = append(services, convertMetaService(service))
	}

	logging.WriteSucceededAuditLog("GetDiscoverPolicies", ctxFields, map[string]interface{}{
		"revision":     revision,
		"serviceCount": len(services),
	})

	return &pb.DiscoverPoliciesResponse{Services: services, Revision: revision}, nil
}

func convertAPIPrincipals(principals []*ads.Principal) []*pb.Principal {
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

func convertAttributes(in map[string]interface{}) map[string]string {
	var out map[string]string
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out

}

func convertAPISubject(subject *ads.Subject) *pb.Subject {
	if subject == nil {
		return nil
	}
	return &pb.Subject{
		Principals: convertAPIPrincipals(subject.Principals),
		Token:      subject.Token,
		TokenType:  subject.TokenType,
	}
}

func convertAPIRequestContext(req *ads.RequestContext) *pb.ContextRequest {
	return &pb.ContextRequest{
		Subject:     convertAPISubject(req.Subject),
		ServiceName: req.ServiceName,
		Resource:    req.Resource,
		Action:      req.Action,
		Attributes:  convertAttributes(req.Attributes),
	}
}
