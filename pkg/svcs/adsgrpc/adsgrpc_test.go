//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package adsgrpc

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	adsPB "github.com/oracle/speedle/pkg/svcs/adsgrpc/pb"
	pmsPB "github.com/oracle/speedle/pkg/svcs/pmsgrpc/pb"
	"github.com/oracle/speedle/testutil"
)

var POLICY_RELOAD_TIME = 500 //ms

//TestMats_GRpc_IsAllowed_User, policy is "grant user userA read res1"
func TestMats_GRpc_IsAllowed_User(t *testing.T) {
	sName := "TestMats_GRpc_IsAllowed_User"
	pName := "Policy_allow_user"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceRequest{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pmsPB.Service{},
				ExpectedBody: &pmsPB.Service{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddPolicy1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.PolicyRequest{
					ServiceName: sName,
					Policy: &pmsPB.Policy{
						Name:   pName,
						Effect: pmsPB.Effect_GRANT,
						Principals: []*pmsPB.AndPrincipals{
							{
								Principals: []string{"user:userA"},
							},
						},
						Permissions: []*pmsPB.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
				OutputBody: &pmsPB.Policy{},
				ExpectedBody: &pmsPB.Policy{
					Name:   pName,
					Effect: pmsPB.Effect_GRANT,
					Principals: []*pmsPB.AndPrincipals{
						{
							Principals: []string{"user:userA"},
						},
					},
					Permissions: []*pmsPB.Policy_Permission{
						{
							Resource: "res1",
							Actions:  []string{"read"},
						},
					},
				},
			},
			PostTestFunc: testutil.PostCreateGetPolicyTest,
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "TestIsAllowed=true (userA read res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res1",
					Action:      "read",
				},
				OutputBody: &adsPB.IsAllowedResponse{},
				ExpectedBody: &adsPB.IsAllowedResponse{
					Allowed: true,
					Reason:  int32(adsapi.GRANT_POLICY_FOUND),
				},
			},
		},
		{
			Name:     "TestIsAllowed=false (invalid-userA read res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "invaliduserA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res1",
					Action:      "read",
				},
				OutputBody: &adsPB.IsAllowedResponse{},
				ExpectedBody: &adsPB.IsAllowedResponse{
					Allowed: false,
					Reason:  int32(adsapi.NO_APPLICABLE_POLICIES),
				},
			},
		},
		{
			Name:     "TestIsAllowed=false (userA invalid-read res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res1",
					Action:      "invalid-read",
				},
				OutputBody: &adsPB.IsAllowedResponse{},
				ExpectedBody: &adsPB.IsAllowedResponse{
					Allowed: false,
					Reason:  int32(adsapi.NO_APPLICABLE_POLICIES),
				},
			},
		},
		{
			Name:     "TestIsAllowed=false (userA read invalid-res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "invalid-res1",
					Action:      "read",
				},
				OutputBody: &adsPB.IsAllowedResponse{},
				ExpectedBody: &adsPB.IsAllowedResponse{
					Allowed: false,
					Reason:  int32(adsapi.NO_APPLICABLE_POLICIES),
				},
			},
		},
		{
			Name:     "TesDeleteService1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pmsPB.Service{},
				ExpectedBody: &pmsPB.Empty{},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
	}
	testutil.RunTestCases(t, data, nil)
}

//TestMats_GRpc_IsAllowed_Role
//policy is: grant role1 read res1
//rolepolicy is: grant user userA role1,role2
func TestMats_GRpc_IsAllowed_Role(t *testing.T) {
	sName := "TestMats_GRpc_IsAllowed_Role"

	rpName := "role-policy1"
	pName := "policy1"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceRequest{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pmsPB.Service{},
				ExpectedBody: &pmsPB.Service{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA"},
						Resources:  []string{"res1"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role1", "role2"},
					Principals: []string{"user:userA"},
					Resources:  []string{"res1"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestAddPolicy1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.PolicyRequest{
					ServiceName: sName,
					Policy: &pmsPB.Policy{
						Name:   pName,
						Effect: pmsPB.Effect_GRANT,
						Principals: []*pmsPB.AndPrincipals{
							{
								Principals: []string{"role:role1"},
							},
						},
						Permissions: []*pmsPB.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
				OutputBody: &pmsPB.Policy{},
				ExpectedBody: &pmsPB.Policy{
					Name:   pName,
					Effect: pmsPB.Effect_GRANT,
					Principals: []*pmsPB.AndPrincipals{
						{
							Principals: []string{"role:role1"},
						},
					},
					Permissions: []*pmsPB.Policy_Permission{
						{
							Resource: "res1",
							Actions:  []string{"read"},
						},
					},
				},
			},
			PostTestFunc: testutil.PostCreateGetPolicyTest,
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "TestIsAllowed=true (userA read res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res1",
					Action:      "read",
				},
				OutputBody: &adsPB.IsAllowedResponse{},
				ExpectedBody: &adsPB.IsAllowedResponse{
					Allowed: true,
					Reason:  int32(adsapi.GRANT_POLICY_FOUND),
				},
			},
		},
		{
			Name:     "TesDeleteService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pmsPB.Service{},
				ExpectedBody: &pmsPB.Empty{},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
	}
	testutil.RunTestCases(t, data, nil)
}

//TestMats_GRpc_GetAllGrantedRoles_basic
//rolepolicy is:
// grant user userA role1,role2
// grant user userA role3,role4 on res1
func TestMats_GRpc_GetAllGrantedRoles_basic(t *testing.T) {
	sName := "TestMats_GRpc_GetAllGrantedRoles_basic"

	rpName1 := "role-policy1"
	rpName2 := "role-policy2"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceRequest{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pmsPB.Service{},
				ExpectedBody: &pmsPB.Service{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddRolePolicy1(grant user userA role1,role2)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName1,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName1,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role1", "role2"},
					Principals: []string{"user:userA"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestAddRolePolicy2(grant user userA role3,role4 on res2)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName2,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role3", "role4"},
						Principals: []string{"user:userA"},
						Resources:  []string{"res1"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName2,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role3", "role4"},
					Principals: []string{"user:userA"},
					Resources:  []string{"res1"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "TestGetAllGrantedRoles with userA on res1, result is role1,role2,role3,role4",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res1",
				},
				OutputBody: &adsPB.AllRoleResponse{},
				ExpectedBody: &adsPB.AllRoleResponse{
					Roles: []string{"role1", "role2", "role3", "role4"},
				},
			},
		},
		{
			Name:     "TestGetAllGrantedRoles with userA, result is role1,role2",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
				},
				OutputBody: &adsPB.AllRoleResponse{},
				ExpectedBody: &adsPB.AllRoleResponse{
					Roles: []string{"role1", "role2"},
				},
			},
		},
		{
			Name:     "TesDeleteService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pmsPB.Service{},
				ExpectedBody: &pmsPB.Empty{},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
	}
	testutil.RunTestCases(t, data, nil)
}

//TestMats_GRpc_IsAllowed_Role
//rolepolicy is:
// grant user userA role1,role2 on res1
// grant user userA role3,role4 on res2
// grant role role1 role5
func TestMats_GRpc_GetAllGrantedRoles_RoleToRole(t *testing.T) {
	sName := "TestMats_GRpc_GetAllGrantedRoles_RoleToRole"

	rpName1 := "role-policy1"
	rpName2 := "role-policy2"
	rpName3 := "role-policy3"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceRequest{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pmsPB.Service{},
				ExpectedBody: &pmsPB.Service{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddRolePolicy1(grant user userA role1,role2 on res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName1,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA"},
						Resources:  []string{"res1"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName1,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role1", "role2"},
					Principals: []string{"user:userA"},
					Resources:  []string{"res1"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestAddRolePolicy2(grant user userA role3,role4 on res2)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName2,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role3", "role4"},
						Principals: []string{"user:userA"},
						Resources:  []string{"res2"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName2,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role3", "role4"},
					Principals: []string{"user:userA"},
					Resources:  []string{"res2"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestAddRolePolicy3(grant role role1 role5)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName3,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role5"},
						Principals: []string{"role:role1"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName3,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role5"},
					Principals: []string{"role:role1"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "TestGetAllGrantedRoles with userA on res1, result is role1,role2,role5",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res1",
				},
				OutputBody: &adsPB.AllRoleResponse{},
				ExpectedBody: &adsPB.AllRoleResponse{
					Roles: []string{"role1", "role2", "role5"},
				},
			},
		},
		{
			Name:     "TestGetAllGrantedRoles with userA on res2, result is role3,role4",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
					Resource:    "res2",
				},
				OutputBody: &adsPB.AllRoleResponse{},
				ExpectedBody: &adsPB.AllRoleResponse{
					Roles: []string{"role3", "role4"},
				},
			},
		},
		{
			Name:     "TestGetAllGrantedRoles with userA, result is empty",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
				},
				OutputBody:   &adsPB.AllRoleResponse{},
				ExpectedBody: &adsPB.AllRoleResponse{},
			},
		},
		{
			Name:     "TesDeleteService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pmsPB.Service{},
				ExpectedBody: &pmsPB.Empty{},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
	}
	testutil.RunTestCases(t, data, nil)
}

//TestMats_GRpc_GetAllGrantedPermission, policy is
// "grant user userA read res1"
// "grant user userA write res1"
// "grant user userA role1"
// "grant role role1 read res2"
func TestMats_GRpc_GetAllGrantedPermission(t *testing.T) {
	sName := "TestMats_GRpc_GetAllGrantedPermission"
	pName1 := "Policy4Permission1"
	pName2 := "Policy4Permission2"
	pName3 := "Policy4Permission3"
	rpName1 := "rolePolicy4Permission"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceRequest{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pmsPB.Service{},
				ExpectedBody: &pmsPB.Service{
					Name: sName,
					Type: pmsPB.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddPolicy1 (grant user userA read res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.PolicyRequest{
					ServiceName: sName,
					Policy: &pmsPB.Policy{
						Name:   pName1,
						Effect: pmsPB.Effect_GRANT,
						Principals: []*pmsPB.AndPrincipals{
							{
								Principals: []string{"user:userA"},
							},
						},
						Permissions: []*pmsPB.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
				OutputBody: &pmsPB.Policy{},
				ExpectedBody: &pmsPB.Policy{
					Name:   pName1,
					Effect: pmsPB.Effect_GRANT,
					Principals: []*pmsPB.AndPrincipals{
						{
							Principals: []string{"user:userA"},
						},
					},
					Permissions: []*pmsPB.Policy_Permission{
						{
							Resource: "res1",
							Actions:  []string{"read"},
						},
					},
				},
			},
			PostTestFunc: testutil.PostCreateGetPolicyTest,
		},
		{
			Name:     "TestAddPolicy2(grant user userA write res1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.PolicyRequest{
					ServiceName: sName,
					Policy: &pmsPB.Policy{
						Name:   pName2,
						Effect: pmsPB.Effect_GRANT,
						Principals: []*pmsPB.AndPrincipals{
							{
								Principals: []string{"user:userA"},
							},
						},
						Permissions: []*pmsPB.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"write"},
							},
						},
					},
				},
				OutputBody: &pmsPB.Policy{},
				ExpectedBody: &pmsPB.Policy{
					Name:   pName2,
					Effect: pmsPB.Effect_GRANT,
					Principals: []*pmsPB.AndPrincipals{
						{
							Principals: []string{"user:userA"},
						},
					},
					Permissions: []*pmsPB.Policy_Permission{
						{
							Resource: "res1",
							Actions:  []string{"write"},
						},
					},
				},
			},
			PostTestFunc: testutil.PostCreateGetPolicyTest,
		},
		{
			Name:     "TestAddRolePolicy(grant user userA role1)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pmsPB.RolePolicy{
						Name:       rpName1,
						Effect:     pmsPB.Effect_GRANT,
						Roles:      []string{"role1"},
						Principals: []string{"user:userA"},
					},
				},
				OutputBody: &pmsPB.RolePolicy{},
				ExpectedBody: &pmsPB.RolePolicy{
					Name:       rpName1,
					Effect:     pmsPB.Effect_GRANT,
					Roles:      []string{"role1"},
					Principals: []string{"user:userA"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestAddPolicy3(grant role role1 read res2)",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.PolicyRequest{
					ServiceName: sName,
					Policy: &pmsPB.Policy{
						Name:   pName3,
						Effect: pmsPB.Effect_GRANT,
						Principals: []*pmsPB.AndPrincipals{
							{
								Principals: []string{"role:role1"},
							},
						},
						Permissions: []*pmsPB.Policy_Permission{
							{
								Resource: "res2",
								Actions:  []string{"read"},
							},
						},
					},
				},
				OutputBody: &pmsPB.Policy{},
				ExpectedBody: &pmsPB.Policy{
					Name:   pName3,
					Effect: pmsPB.Effect_GRANT,
					Principals: []*pmsPB.AndPrincipals{
						{
							Principals: []string{"role:role1"},
						},
					},
					Permissions: []*pmsPB.Policy_Permission{
						{
							Resource: "res2",
							Actions:  []string{"read"},
						},
					},
				},
			},
			PostTestFunc: testutil.PostCreateGetPolicyTest,
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "TestGetPermissions with userA",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.GRpcTestData{
				InputBody: &adsPB.ContextRequest{
					Subject: &adsPB.Subject{
						Principals: []*adsPB.Principal{
							{
								Type: "user",
								Name: "userA",
							},
						},
					},
					ServiceName: sName,
				},
				OutputBody: &adsPB.AllPermissionResponse{},
				ExpectedBody: &adsPB.AllPermissionResponse{
					Permissions: []*adsPB.AllPermissionResponse_Permission{
						{
							Resource: "res1",
							Actions:  []string{"read"},
						},
						{
							Resource: "res1",
							Actions:  []string{"write"},
						},
						{
							Resource: "res2",
							Actions:  []string{"read"},
						},
					}},
			},
		},
		{
			Name:     "TesDeleteService1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pmsPB.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pmsPB.Service{},
				ExpectedBody: &pmsPB.Empty{},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
	}
	testutil.RunTestCases(t, data, nil)
}
