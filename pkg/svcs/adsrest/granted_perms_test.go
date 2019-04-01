//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

//+build runtime_test

package adsrest

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/testutil"
)

var URI_GRANTED_PERMS = "/authz-check/v1/all-granted-permissions"

//Policies are defined in check_prepare_test.go : "service-simple"
/*=== Policies used by test cases===
SERVICE_SIMPLE
	"role-policies:",
		"grant user  userWithRole1 role1 on res_allow",
		"deny  group groupWithRole2 role2 on res_deny",
		"policies:",
		"grant user  user1  get,del res_allow",
		"grant group group1 get,del res_allow",
		"grant role  role1  get,del res_allow",

		"deny  user  user1  get,del res_deny",
		"deny  group group1 get,del res_deny",
		"deny  role  role1  get,del res_deny",
		"grant role  role2  get,del res_deny",
*/
func TestMats_GrantPerms_Simple(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "OnePolicy",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user1",
							},
						},
					},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{Resource: "res_allow", Actions: []string{"get", "del"}},
				},
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-complex"
func TestLrg_GrantedPerms_Complex_bug214(t *testing.T) {
	data := &[]testutil.TestCase{

		//"role-policies:",
		//"grant user user_complex1, user user_complex1A,user user_complex1B role_complex1",
		//"policies:",
		//"grant role role_complex1 get,del res_complex1 if request_user != 'user_complex1'",
		//"deny user user_complex1A del res_complex1",
		{
			Name:     "Nopms.Permission-excludedByCondition",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user_complex1",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody:   &[]pms.Permission{},
			},
		},
		{
			Name:     "OneActionDenied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user_complex1A",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{Resource: "res_complex1", Actions: []string{"get"}},
				},
			},
		},
		{
			Name:     "AllActionAllowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user_complex1B",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{Resource: "res_complex1", Actions: []string{"get", "del"}},
				},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-complex-role
func TestLrg_GrantedPerms_ComplexRole_bug214(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "MultiPolicies",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user1",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{Resource: "res1", Actions: []string{"get", "del"}},
					{Resource: "res2", Actions: []string{"get", "del"}},
					{Resource: "res3", Actions: []string{"get", "del"}},
					{Resource: "res9", Actions: []string{"get", "del"}},
				},
			},
		},
		{
			Name:     "MultiPolicies-OneActionDeined",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user11",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{Resource: "res1", Actions: []string{"get", "del"}},
					{Resource: "res2", Actions: []string{"get", "del"}},
					{Resource: "res3", Actions: []string{"get"}},
					{Resource: "res9", Actions: []string{"get", "del"}},
				},
			},
			//PostTestFunc: testutil.PostGetAllGrantedpms.Permissions,
		},
		{
			Name:     "MultiPolicies-MoreOneActionAllowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user2",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{Resource: "res2", Actions: []string{"get", "del"}},
					{Resource: "res3", Actions: []string{"get", "del"}},
					{Resource: "res9", Actions: []string{"get", "del"}},
					{Resource: "res1", Actions: []string{"get"}},
				},
			},
			//PostTestFunc: testutil.PostGetAllGrantedpms.Permissions,
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-with-resexpr"
//Bug#123: Due to limitation of implementation, Can't get grantedpms.Permissions if policy is with resourceexpression
func testLrg_GrantedPerms_ResExpr(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "Case1:",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_PERMISSIONS,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_PERMS,
				InputBody: &JsonContext{
					Subject: &JsonSubject{
						Principals: []*JsonPrincipal{
							{
								Type: adsapi.PRINCIPAL_TYPE_USER,
								Name: "user1",
							},
						},
					},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]pms.Permission{},
				ExpectedBody: &[]pms.Permission{
					{ResourceExpression: "res.*", Actions: []string{"get", "del"}},
				},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
