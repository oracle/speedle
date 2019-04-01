//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsrest

import (
	"testing"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/testutil"
)

//Test REST API for Service Manangement
func TestMats_PMSRest_Service(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "TestGetEmptyService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Service{},
				ExpectedBody:   &[]pms.Service{},
			},
		},
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: "app1",
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: "app1",
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "TestAddSameService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 400,
				InputBody: &pms.Service{
					Name: "app1",
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "TestAddK8sService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: "k8s",
					Type: pms.TypeK8SCluster,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: "k8s",
					Type: pms.TypeK8SCluster,
				},
			},
		},
		{
			Name:     "TestGetAllService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_QUERY_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Service{},
				ExpectedBody: &[]pms.Service{
					{
						Name: "app1",
						Type: pms.TypeApplication,
					},
					{
						Name: "k8s",
						Type: pms.TypeK8SCluster,
					},
				},
			},
		},
		{
			Name:     "TestGetK8sService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/k8s",
				ExpectedStatus: 200,
				OutputBody:     &pms.Service{},
				ExpectedBody: &pms.Service{
					Name:         "k8s",
					Type:         pms.TypeK8SCluster,
					Policies:     []*pms.Policy{},
					RolePolicies: []*pms.RolePolicy{},
				},
			},
		},
		{
			Name:     "TestGetDelK8sService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/k8s",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestGetK8sServiceNeg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/k8s",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelK8sServiceNeg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/k8s",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelAllService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestDelAllService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestGetApp1ServiceNeg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/app1",
				ExpectedStatus: 404,
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Test REST API for Policy manangent
func TestMats_PMSRest_Policy(t *testing.T) {
	srvName := "TestMats_PMSRest_Policy"
	policyURI := testutil.URI_POLICY_MGMT + "service/" + srvName + "/policy"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: srvName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: srvName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "TestGetEmptyPolicy",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_QUERY_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI,
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Policy{},
				ExpectedBody:   &[]pms.Policy{},
			},
		},
		{
			Name:     "TestAddPolicy1",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI,
				ExpectedStatus: 201,
				InputBody: &pms.Policy{
					Name:   "policy1",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"user:bill", "user:cynthia"}},
				},
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy1",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"user:bill", "user:cynthia"}},
				},
			},
		},
		{
			Name:     "TestAddPolicy2",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI,
				ExpectedStatus: 201,
				InputBody: &pms.Policy{
					Name:   "policy2",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"read"},
						},
					},
					Principals: [][]string{{"user:bill"}, {"user:cynthia"}},
				},
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy2",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"read"},
						},
					},
					Principals: [][]string{{"user:bill"}, {"user:cynthia"}},
				},
			},
		},
		{
			Name:     "TestGetAllPolicies",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_QUERY_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI,
				ExpectedStatus: 200,
				OutputBody:     &[]*pms.Policy{},
				ExpectedBody: &[]*pms.Policy{
					{
						Name:   "policy1",
						Effect: pms.Grant,
						Permissions: []*pms.Permission{
							{
								Resource: "res1",
								Actions:  []string{"read", "write"},
							},
						},
						Principals: [][]string{{"user:bill", "user:cynthia"}},
					},
					{
						Name:   "policy2",
						Effect: pms.Grant,
						Permissions: []*pms.Permission{
							{
								Resource: "res2",
								Actions:  []string{"read"},
							},
						},
						Principals: [][]string{{"user:bill"}, {"user:cynthia"}},
					},
				},
			},
		},
		{
			Name:     "TestGetPolicy2",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI + "/policy2",
				ExpectedStatus: 200,
				OutputBody:     &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy2",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"read"},
						},
					},
					Principals: [][]string{{"user:bill"}, {"user:cynthia"}},
				},
			},
		},
		{
			Name:     "TestDeletePolicy2",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI + "/policy2",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestGetPolicy2Neg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI + "/policy2",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelPolicy2Neg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI + "/policy2",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelAllPolicy",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI,
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestGetPolicy1Neg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI + "/policy1",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestAddPolicyWithPrincpleContainingTwoSpiffeID",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI,
				ExpectedStatus: 201,
				InputBody: &pms.Policy{
					Name:   "policy-2-entity",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"entity:spiffe:staging.acme.com/payments/mysql"},
						{"entity:spiffe:staging.acme.com/payments/web-fe"}},
				},
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy-2-entity",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"entity:spiffe:staging.acme.com/payments/mysql"},
						{"entity:spiffe:staging.acme.com/payments/web-fe"}},
				},
			},
		},
		{
			Name:     "TestGetPolicyWithPrincpleContainingTwoSpiffeID",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.RestTestData{
				URI:            policyURI + "/policy-2-entity",
				ExpectedStatus: 200,
				OutputBody:     &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy-2-entity",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"entity:spiffe:staging.acme.com/payments/mysql"},
						{"entity:spiffe:staging.acme.com/payments/web-fe"}},
				},
			},
		},
		{
			Name:     "TestDelAllService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Test REST API for RolePolicy Manangement
func TestMats_PMSRest_RolePolicy(t *testing.T) {
	srvName := "TestMats_PMSRest_RolePolicy"
	rolePolicyURI := testutil.URI_POLICY_MGMT + "service/" + srvName + "/role-policy"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: srvName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: srvName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "TestGetEmptyRolePolicy",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI,
				ExpectedStatus: 200,
				OutputBody:     &[]pms.RolePolicy{},
				ExpectedBody:   &[]pms.RolePolicy{},
			},
		},
		{
			Name:     "TestAddRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI,
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"manager", "staff"},
					Principals: []string{"user:bill", "user:cynthia"},
					Resources:  []string{"res1"},
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"manager", "staff"},
					Principals: []string{"user:bill", "user:cynthia"},
					Resources:  []string{"res1"},
				},
			},
		},
		{
			Name:     "TestAddRolePolicy2",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI,
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy2",
					Effect:     pms.Deny,
					Roles:      []string{"manager"},
					Principals: []string{"user:bill", "user:cynthia"},
					Resources:  []string{"res2"},
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy2",
					Effect:     pms.Deny,
					Roles:      []string{"manager"},
					Principals: []string{"user:bill", "user:cynthia"},
					Resources:  []string{"res2"},
				},
			},
		},
		{
			Name:     "TestGetAllRolePolicies",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI,
				ExpectedStatus: 200,
				OutputBody:     &[]*pms.RolePolicy{},
				ExpectedBody: &[]*pms.RolePolicy{
					{
						Name:       "role-policy1",
						Effect:     pms.Grant,
						Roles:      []string{"manager", "staff"},
						Principals: []string{"user:bill", "user:cynthia"},
						Resources:  []string{"res1"},
					},
					{
						Name:       "role-policy2",
						Effect:     pms.Deny,
						Roles:      []string{"manager"},
						Principals: []string{"user:bill", "user:cynthia"},
						Resources:  []string{"res2"},
					},
				},
			},
		},
		{
			Name:     "TestGetRolePolicy2",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI + "/role-policy2",
				ExpectedStatus: 200,
				OutputBody:     &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy2",
					Effect:     pms.Deny,
					Roles:      []string{"manager"},
					Principals: []string{"user:bill", "user:cynthia"},
					Resources:  []string{"res2"},
				},
			},
		},
		{
			Name:     "TestDeleteRolePolicy2",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI + "/role-policy2",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestGetRolePolicy2Neg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI + "/role-policy2",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelRolePolicy2Neg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI + "/role-policy2",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelAllRolePolicy",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI,
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "TestGetRolePolicy1Neg",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            rolePolicyURI + "/role-policy1",
				ExpectedStatus: 404,
			},
		},
		{
			Name:     "TestDelAllService",
			Enabled:  true,
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
