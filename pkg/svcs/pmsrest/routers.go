//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsrest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/svcs"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func initRouters(ps pms.PolicyStoreManager) (*[]route, error) {

	manager, err := NewRestService(ps)
	if err != nil {
		return nil, err
	}

	svcRoutes := []route{}

	policyManagerRoutes := []route{
		{
			"CreatePolicy",
			"POST",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy",
			manager.CreatePolicy,
		},

		{
			"DeletePolicies",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy",
			manager.DeletePolicies,
		},

		{
			"DeletePolicy",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy/{policyID}",
			manager.DeletePolicy,
		},

		{
			"GetPolicy",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy/{policyID}",
			manager.GetPolicy,
		},

		{
			"ListPolicies",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy",
			manager.ListPolicies,
		},

		{
			"CreateRolePolicy",
			"POST",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy",
			manager.CreateRolePolicy,
		},

		{
			"DeleteRolePolicies",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy",
			manager.DeleteRolePolicies,
		},

		{
			"DeleteRolePolicy",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy/{rolePolicyID}",
			manager.DeleteRolePolicy,
		},

		{
			"GetRolePolicy",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy/{rolePolicyID}",
			manager.GetRolePolicy,
		},

		{
			"ListRolePolicies",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy",
			manager.ListRolePolicies,
		},
	}
	svcRoutes = append(svcRoutes, policyManagerRoutes...)

	serviceManageRoutes := []route{
		{
			"CreateService",
			"POST",
			svcs.PolicyMgmtPath + "service",
			manager.CreateService,
		},

		{
			"DeleteService",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}",
			manager.DeleteService,
		},

		{
			"DeleteServices",
			"DELETE",
			svcs.PolicyMgmtPath + "service",
			manager.DeleteServices,
		},

		{
			"GetService",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}",
			manager.GetService,
		},

		{
			"ListServices",
			"GET",
			svcs.PolicyMgmtPath + "service",
			manager.ListServices,
		},

		{
			"ListPolicyCounts",
			"GET",
			svcs.PolicyMgmtPath + "policy-counts",
			manager.ListPolicyAndRolePolicyCounts,
		},
	}
	svcRoutes = append(svcRoutes, serviceManageRoutes...)

	functionManageRoutes := []route{
		{
			"CreateFunction",
			"POST",
			svcs.PolicyMgmtPath + "function",
			manager.CreateFunction,
		},

		{
			"DeleteFunction",
			"DELETE",
			svcs.PolicyMgmtPath + "function/{functionName}",
			manager.DeleteFunction,
		},

		{
			"DeleteFunctions",
			"DELETE",
			svcs.PolicyMgmtPath + "function",
			manager.DeleteFunctions,
		},

		{
			"GetFunction",
			"GET",
			svcs.PolicyMgmtPath + "function/{functionName}",
			manager.GetFunction,
		},

		{
			"ListFunctions",
			"GET",
			svcs.PolicyMgmtPath + "function",
			manager.ListFunctions,
		},
	}
	svcRoutes = append(svcRoutes, functionManageRoutes...)

	discoverRequestManageRoutes := []route{
		{
			"GetAllDiscoverRequests",
			"GET",
			svcs.PolicyMgmtPath + "discover-request",
			manager.GetAllDiscoverRequests,
		},

		{
			"GetDiscoverRequests",
			"GET",
			svcs.PolicyMgmtPath + "discover-request/{serviceName}",
			manager.GetDiscoverRequests,
		},

		{
			"ResetDiscoverRequests",
			"DELETE",
			svcs.PolicyMgmtPath + "discover-request/{serviceName}",
			manager.ResetDiscoverRequests,
		},

		{
			"ResetAllDiscoverRequests",
			"DELETE",
			svcs.PolicyMgmtPath + "discover-request",
			manager.ResetAllDiscoverRequests,
		},

		{
			"GetDiscoverPolicies",
			"GET",
			svcs.PolicyMgmtPath + "discover-policy/{serviceName}",
			manager.GetDiscoverPolicies,
		},

		{
			"GetAllDiscoverPolicies",
			"GET",
			svcs.PolicyMgmtPath + "discover-policy",
			manager.GetDiscoverPolicies,
		},
	}
	svcRoutes = append(svcRoutes, discoverRequestManageRoutes...)

	return &svcRoutes, nil

}

func NewRouter(ps pms.PolicyStoreManager) (*mux.Router, error) {
	routes, err := initRouters(ps)
	if err != nil {
		return nil, err
	}
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range *routes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router, nil
}
