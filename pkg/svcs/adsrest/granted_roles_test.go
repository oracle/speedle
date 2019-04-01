//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

//+build runtime_test

package adsrest

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/testutil"
)

var URI_GRANTED_ROLES = "/authz-check/v1/all-granted-roles"

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
func TestMats_GrantedRole_Simple(t *testing.T) {
	data := &[]testutil.TestCase{
		//
		{
			Name:     "NoRole",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{},
			},
		},
		{
			Name:     "GrantUserToRole",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userWithRole1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role1"},
			},
		},
		{
			Name:     "DenyGroupToRole",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "group2WithRole2"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_deny",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{},
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-complex-role"
func TestLrg_GrantedRole_RoleEmbeded_bug120(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "user1-in-allroles",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role1", "role2", "role3", "role4", "role5", "role6", "role7", "role8", "role9"},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
		{
			Name:     "user2-not-in-role1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role2", "role3", "role4", "role5", "role6", "role7", "role8", "role9"},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
		{
			Name:     "user1-in-allroles",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role1", "role2", "role3", "role4", "role5", "role6", "role7", "role8", "role9"},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
		{
			Name:     "user22-denied-in-middle-role",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user22"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role-denined", "role-denined1", "role2", "role3", "role4"},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Test granted roles are cprrect if the rolepolicy is only for specified resource
/*
	Policies are defined in check_prepare_test.go : "service-complex-role"
	"grant user userRes1 role10  on res1",
	"grant user userRes1 role11 if request_action == 'get'",
	"grant role role10 role12 on res1"
	"grant role role10 role12 on res2"
*/
func TestLrg_GrantedRole_OnResource(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "userRes1-on-allowed-res",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userRes1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role10", "role11", "role12"},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
		{
			Name:     "userRes1-denied-by-Condition",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userRes1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role10", "role12"},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
		{
			Name:     "userRes1-on-NOT-allowed-res",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userRes1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res2",
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{},
			},
			PostTestFunc: testutil.PostGetAllGrantedRoles,
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Test granted roles are correct if the policy/rolepolicy containing condition
/*
	Policies are defined in check_prepare_test.go : "service-simple"
	"grant user user_condtion1 role_condition1 if age > 10",
*/
func TestLrg_GrantedRole_WithCondition(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "user_condtion1-with-condition(age>10), request age=15",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_condtion1"}}},
					ServiceName: SERVICE_SIMPLE,
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 15}},
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{"role_condition1"},
			},
		},
		{
			Name:     "user_condtion1-with-condition(age>10), request age=10",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_GET_GRANTED_ROLES,
			Data: &testutil.RestTestData{
				URI: URI_GRANTED_ROLES,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_condtion1"}}},
					ServiceName: SERVICE_SIMPLE,
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 10}},
				},
				ExpectedStatus: 200,
				OutputBody:     &[]string{},
				ExpectedBody:   &[]string{},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
