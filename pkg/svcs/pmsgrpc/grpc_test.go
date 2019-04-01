//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsgrpc

import (
	"testing"

	"github.com/oracle/speedle/pkg/svcs/pmsgrpc/pb"
	"github.com/oracle/speedle/testutil"
)

//Test GRpc API for Service Manangement
func TestMats_GRpc_Service(t *testing.T) {
	sName := "TestMats_GRpc_Sevice"
	sName1 := "TestMats_GRpc_Sevice1"
	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceRequest{
					Name: sName,
					Type: pb.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pb.Service{},
				ExpectedBody: &pb.Service{
					Name: sName,
					Type: pb.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddService1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceRequest{
					Name: sName1,
					Type: pb.ServiceType_APPLICATION,
				},
				OutputBody: &pb.Service{},
				ExpectedBody: &pb.Service{
					Name: sName1,
					Type: pb.ServiceType_APPLICATION,
				},
			},
		},
		{
			Name:     "TestQueryServiceWithSpecifiedName",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody: &[]*pb.Service{},
				ExpectedBody: &[]*pb.Service{
					{
						Name: sName,
						Type: pb.ServiceType_K8S_CLUSTER,
					},
				},
			},
		},
		{
			Name:     "TestQueryServiceWithoutName",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody:  &pb.ServiceQueryRequest{},
				OutputBody: &[]*pb.Service{},
				ExpectedBody: &[]*pb.Service{
					{
						Name: sName,
						Type: pb.ServiceType_K8S_CLUSTER,
					},
					{
						Name: sName1,
						Type: pb.ServiceType_APPLICATION,
					},
				},
			},
			PostTestFunc: testutil.PostListServiceTest,
		},
		{
			Name:     "TestDeleteService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pb.Service{},
				ExpectedBody: &pb.Empty{},
			},
		},
		{
			Name:     "TestQueryService1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceQueryRequest{
					Name: sName1,
				},
				OutputBody: &[]*pb.Service{},
				ExpectedBody: &[]*pb.Service{
					{
						Name: sName1,
					},
				},
			},
		},
		{
			Name:     "TesDeleteService1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceQueryRequest{
					Name: sName1,
				},
				OutputBody:   &pb.Service{},
				ExpectedBody: &pb.Empty{},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Test GRpc API for Policy Manangement
func TestMats_GRpc_Policy(t *testing.T) {
	sName := "TestMats_GRpc_Policy"

	pName := "policy1"
	pName1 := "policy2"
	pName_entity := "policy_with_entity"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceRequest{
					Name: sName,
					Type: pb.ServiceType_APPLICATION,
				},
				OutputBody: &pb.Service{},
				ExpectedBody: &pb.Service{
					Name: sName,
					Type: pb.ServiceType_APPLICATION,
				},
			},
		},
		{
			Name:     "TestAddPolicy1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyRequest{
					ServiceName: sName,
					Policy: &pb.Policy{
						Name:   pName,
						Effect: pb.Effect_GRANT,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"user:userA", "group:groupA"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
				OutputBody: &pb.Policy{},
				ExpectedBody: &pb.Policy{
					Name:   pName,
					Effect: pb.Effect_GRANT,
					Principals: []*pb.AndPrincipals{
						{
							Principals: []string{"user:userA", "group:groupA"},
						},
					},
					Permissions: []*pb.Policy_Permission{
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
			Name:     "TestAddPolicy2",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyRequest{
					ServiceName: sName,
					Policy: &pb.Policy{
						Name:   pName1,
						Effect: pb.Effect_DENY,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"user:userA", "user:userB"},
							},
							{
								Principals: []string{"user:userA", "group:groupB"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res2",
								Actions:  []string{"write"},
							},
						},
					},
				},
				OutputBody: &pb.Policy{},
				ExpectedBody: &pb.Policy{
					Name:   pName1,
					Effect: pb.Effect_DENY,
					Principals: []*pb.AndPrincipals{
						{
							Principals: []string{"user:userA", "user:userB"},
						},
						{
							Principals: []string{"user:userA", "group:groupB"},
						},
					},
					Permissions: []*pb.Policy_Permission{
						{
							Resource: "res2",
							Actions:  []string{"write"},
						},
					},
				},
			},
			PostTestFunc: testutil.PostCreateGetPolicyTest,
		},
		{
			Name:     "TestGetPolicy",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyQueryRequest{
					ServiceName: sName,
					PolicyID:    pName, //would be updated with ID in PreTestFunc
				},
				OutputBody: &[]*pb.Policy{},
				ExpectedBody: &[]*pb.Policy{
					{
						Name:   pName,
						Effect: pb.Effect_GRANT,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"user:userA", "group:groupA"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
			},
		},
		{
			Name:     "TestListAllPolicies",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyQueryRequest{
					ServiceName: sName,
				},
				OutputBody: &[]*pb.Policy{},
				ExpectedBody: &[]*pb.Policy{
					{
						Name:   pName,
						Effect: pb.Effect_GRANT,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"user:userA", "group:groupA"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
					{
						Name:   pName1,
						Effect: pb.Effect_DENY,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"user:userA", "user:userB"},
							},
							{
								Principals: []string{"user:userA", "group:groupB"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res2",
								Actions:  []string{"write"},
							},
						},
					},
				},
			},
			PostTestFunc: testutil.PostListPolicyTest,
		},
		{
			Name:     "TestListPoliciesByName",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyQueryRequest{
					ServiceName: sName,
					Filters:     "name eq " + pName,
				},
				OutputBody: &[]*pb.Policy{},
				ExpectedBody: &[]*pb.Policy{
					{
						Name:   pName,
						Effect: pb.Effect_GRANT,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"user:userA", "group:groupA"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
			},
			PostTestFunc: testutil.PostListPolicyTest,
		},
		{
			Name:     "TesDeletePolicy",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyQueryRequest{
					ServiceName: sName,
					PolicyID:    pName, //would be updated with ID in PreTestFunc
				},
				OutputBody:   &pb.Service{},
				ExpectedBody: &pb.Empty{},
			},
		},
		{
			Name:     "TestGetPolicy-Non-existing",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyQueryRequest{
					ServiceName: sName,
					PolicyID:    pName, //would be updated with ID in PreTestFunc
				},
				OutputBody:   &pb.Policy{},
				ExpectedMsg:  "rpc error: code = NotFound",
				ExpectedBody: &pb.Empty{},
			},
		},
		{
			Name:     "TestAddPolicyContainingEntityPrinciple",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.PolicyRequest{
					ServiceName: sName,
					Policy: &pb.Policy{
						Name:   pName_entity,
						Effect: pb.Effect_GRANT,
						Principals: []*pb.AndPrincipals{
							{
								Principals: []string{"entity:schema://domain.name/path1"},
							},
							{
								Principals: []string{"entity:schema://domain.name/path2"},
							},
						},
						Permissions: []*pb.Policy_Permission{
							{
								Resource: "res1",
								Actions:  []string{"read"},
							},
						},
					},
				},
				OutputBody: &pb.Policy{},
				ExpectedBody: &pb.Policy{
					Name:   pName_entity,
					Effect: pb.Effect_GRANT,
					Principals: []*pb.AndPrincipals{
						{
							Principals: []string{"entity:schema://domain.name/path1"},
						},
						{
							Principals: []string{"entity:schema://domain.name/path2"},
						},
					},
					Permissions: []*pb.Policy_Permission{
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
			Name:     "TesDeleteService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pb.Service{},
				ExpectedBody: &pb.Empty{},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Failed by issue:https://gitlab-odx.oracledx.com/wcai/kauthz/issues/42
//Test GRpc API for RolePolicy Manangement
func TestMats_GRpc_RolePolicy(t *testing.T) {
	sName := "TestMats_GRpc_Policy"

	rpName := "role-policy1"
	rpName1 := "role-policy2"

	data := &[]testutil.TestCase{
		{
			Name:     "TestAddService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceRequest{
					Name: sName,
					Type: pb.ServiceType_K8S_CLUSTER,
				},
				OutputBody: &pb.Service{},
				ExpectedBody: &pb.Service{
					Name: sName,
					Type: pb.ServiceType_K8S_CLUSTER,
				},
			},
		},
		{
			Name:     "TestAddRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pb.RolePolicy{
						Name:       rpName,
						Effect:     pb.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA", "group:groupA"},
						Resources:  []string{"res1"},
					},
				},
				OutputBody: &pb.RolePolicy{},
				ExpectedBody: &pb.RolePolicy{
					Name:       rpName,
					Effect:     pb.Effect_GRANT,
					Roles:      []string{"role1", "role2"},
					Principals: []string{"user:userA", "group:groupA"},
					Resources:  []string{"res1"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestAddRolePolicy2",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyRequest{
					ServiceName: sName,
					RolePolicy: &pb.RolePolicy{
						Name:       rpName1,
						Effect:     pb.Effect_DENY,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA", "user:userB", "user:userA", "group:groupA"},
						Resources:  []string{"res2"},
					},
				},
				OutputBody: &pb.RolePolicy{},
				ExpectedBody: &pb.RolePolicy{
					Name:       rpName1,
					Effect:     pb.Effect_DENY,
					Roles:      []string{"role1", "role2"},
					Principals: []string{"user:userA", "user:userB", "user:userA", "group:groupA"},
					Resources:  []string{"res2"},
				},
			},
			PostTestFunc: testutil.PostCreateGetRolePolicyTest,
		},
		{
			Name:     "TestGetRolePolicy",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyQueryRequest{
					ServiceName:  sName,
					RolePolicyID: rpName, //would be updated with ID in PreTestFunc
				},
				OutputBody: &[]*pb.RolePolicy{},
				ExpectedBody: &[]*pb.RolePolicy{
					{
						Name:       rpName,
						Effect:     pb.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA", "group:groupA"},
						Resources:  []string{"res1"},
					},
				},
			},
		},
		{
			Name:     "TestListAllRolePolicies",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyQueryRequest{
					ServiceName: sName,
				},
				OutputBody: &[]*pb.RolePolicy{},
				ExpectedBody: &[]*pb.RolePolicy{
					{
						Name:       rpName,
						Effect:     pb.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA", "group:groupA"},
						Resources:  []string{"res1"},
					},
					{
						Name:       rpName1,
						Effect:     pb.Effect_DENY,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA", "user:userB", "user:userA", "group:groupA"},
						Resources:  []string{"res2"},
					},
				},
			},
			PostTestFunc: testutil.PostListRolePolicyTest,
		},
		{
			Name:     "TestListAllRolePoliciesByName",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyQueryRequest{
					ServiceName: sName,
					Filters:     "name eq " + rpName,
				},
				OutputBody: &[]*pb.RolePolicy{},
				ExpectedBody: &[]*pb.RolePolicy{
					{
						Name:       rpName,
						Effect:     pb.Effect_GRANT,
						Roles:      []string{"role1", "role2"},
						Principals: []string{"user:userA", "group:groupA"},
						Resources:  []string{"res1"},
					},
				},
			},
			PostTestFunc: testutil.PostListRolePolicyTest,
		},
		{
			Name:     "TesDeleteRolePolicy",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyQueryRequest{
					ServiceName:  sName,
					RolePolicyID: rpName, //would be updated with ID in PreTestFunc
				},
				OutputBody:   &pb.Empty{},
				ExpectedBody: &pb.Empty{},
			},
		},
		{
			Name:     "TestGetPolicy-Non-existing",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.RolePolicyQueryRequest{
					ServiceName:  sName,
					RolePolicyID: rpName, //would be updated with ID in PreTestFunc
				},
				ExpectedMsg: "rpc error: code = NotFound",
			},
		},
		{
			Name:     "TesDeleteService",
			Enabled:  true,
			Executer: testutil.NewGRpcTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.GRpcTestData{
				InputBody: &pb.ServiceQueryRequest{
					Name: sName,
				},
				OutputBody:   &pb.Service{},
				ExpectedBody: &pb.Empty{},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
