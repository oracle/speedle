//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsrest

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab-odx.oracledx.com/wcai/speedle/api/pms"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/svcs"
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
		route{
			"CreatePolicy",
			"POST",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy",
			manager.CreatePolicy,
		},

		route{
			"DeletePolicies",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy",
			manager.DeletePolicies,
		},

		route{
			"DeletePolicy",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy/{policyID}",
			manager.DeletePolicy,
		},

		route{
			"GetPolicy",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy/{policyID}",
			manager.GetPolicy,
		},

		route{
			"ListPolicies",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/policy",
			manager.ListPolicies,
		},

		route{
			"CreateRolePolicy",
			"POST",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy",
			manager.CreateRolePolicy,
		},

		route{
			"DeleteRolePolicies",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy",
			manager.DeleteRolePolicies,
		},

		route{
			"DeleteRolePolicy",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy/{rolePolicyID}",
			manager.DeleteRolePolicy,
		},

		route{
			"GetRolePolicy",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy/{rolePolicyID}",
			manager.GetRolePolicy,
		},

		route{
			"ListRolePolicies",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}/role-policy",
			manager.ListRolePolicies,
		},
	}
	svcRoutes = append(svcRoutes, policyManagerRoutes...)

	serviceManageRoutes := []route{
		route{
			"CreateService",
			"POST",
			svcs.PolicyMgmtPath + "service",
			manager.CreateService,
		},

		route{
			"DeleteService",
			"DELETE",
			svcs.PolicyMgmtPath + "service/{serviceName}",
			manager.DeleteService,
		},

		route{
			"DeleteServices",
			"DELETE",
			svcs.PolicyMgmtPath + "service",
			manager.DeleteServices,
		},

		route{
			"GetService",
			"GET",
			svcs.PolicyMgmtPath + "service/{serviceName}",
			manager.GetService,
		},

		route{
			"ListServices",
			"GET",
			svcs.PolicyMgmtPath + "service",
			manager.ListServices,
		},

		route{
			"ListPolicyCounts",
			"GET",
			svcs.PolicyMgmtPath + "policy-counts",
			manager.ListPolicyAndRolePolicyCounts,
		},
	}
	svcRoutes = append(svcRoutes, serviceManageRoutes...)

	functionManageRoutes := []route{
		route{
			"CreateFunction",
			"POST",
			svcs.PolicyMgmtPath + "function",
			manager.CreateFunction,
		},

		route{
			"DeleteFunction",
			"DELETE",
			svcs.PolicyMgmtPath + "function/{functionName}",
			manager.DeleteFunction,
		},

		route{
			"DeleteFunctions",
			"DELETE",
			svcs.PolicyMgmtPath + "function",
			manager.DeleteFunctions,
		},

		route{
			"GetFunction",
			"GET",
			svcs.PolicyMgmtPath + "function/{functionName}",
			manager.GetFunction,
		},

		route{
			"ListFunctions",
			"GET",
			svcs.PolicyMgmtPath + "function",
			manager.ListFunctions,
		},
	}
	svcRoutes = append(svcRoutes, functionManageRoutes...)

	discoverRequestManageRoutes := []route{
		route{
			"GetAllDiscoverRequests",
			"GET",
			svcs.PolicyMgmtPath + "discover-request",
			manager.GetAllDiscoverRequests,
		},

		route{
			"GetDiscoverRequests",
			"GET",
			svcs.PolicyMgmtPath + "discover-request/{serviceName}",
			manager.GetDiscoverRequests,
		},

		route{
			"ResetDiscoverRequests",
			"DELETE",
			svcs.PolicyMgmtPath + "discover-request/{serviceName}",
			manager.ResetDiscoverRequests,
		},

		route{
			"ResetAllDiscoverRequests",
			"DELETE",
			svcs.PolicyMgmtPath + "discover-request",
			manager.ResetAllDiscoverRequests,
		},

		route{
			"GetDiscoverPolicies",
			"GET",
			svcs.PolicyMgmtPath + "discover-policy/{serviceName}",
			manager.GetDiscoverPolicies,
		},

		route{
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
