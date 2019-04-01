//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

//+build runtime_test

package adsrest

import (
	"testing"
	"time"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/svcs"
	"github.com/oracle/speedle/testutil"
)

const (
	SERVICE_SIMPLE                           = "service-simple"
	SERVICE_BOTH_GRANT_DENY                  = "service-both-grant-deny"
	SERVICE_COND_FUNC                        = "service-condition-func"
	SERVICE_COND_ATTRIBUTE                   = "service-condition-attribute"
	SERVICE_BUILTIN_ATTRIBUTE                = "service-builtin-attribute"
	SERVICE_COMPLEX                          = "service-complex"
	SERVICE_COMPLEX_ROLE                     = "service-complex-role"
	SERVICE_COMPLEX_RESEXPR                  = "service-with-resexpr"
	SERVICE_COMPLEX_PRINCIPLE_IN_POILICY     = "service-with-complex-principle-policy"
	SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY = "service-with-complex-principle_rolePolicy"
	SERVICE_WITH_ENTITY_PRINCIPLE            = "service-with-entity-principle"
)

var URI_POLICY_MGMT = svcs.PolicyMgmtPath
var URI_IS_ALLOWD = svcs.PolicyAtzPath + "is-allowed"

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
//Policy is simple to Grant user/group to do action on resource. Including invalid res,action,subject,group cases
func TestMats_IsAllowed_SimpleGrant(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "GrantUser",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantUser-invalidAction",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "get1",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "GrantUser-invalidResource",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "GrantUser_subjectContainsUser",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1, user2"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "GrantUser_withAnyGroup",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "groupAny"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantGroup",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantGroup_withAnyUser",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantRoleToUser-RoleAllowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userWithRole1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_allow",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantRoleToUser-RoleDenied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userWithRole1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_deny",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-simple"
//Policy is simple to Deny user/group/role to do action on resource.
func TestMats_IsAllowed_SimpleDeny(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "DenyUser",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_deny",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
		{
			Name:     "DenyGroup",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_deny",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantRoleToUser_RoleDenied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userWithRole1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_deny",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
		{
			Name:     "DenyRoleToGroup_RoleAllowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "groupWithRole2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_SIMPLE,
					Resource:    "res_deny",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-both-grant-deny"
/*
	"role-policies:",
	"grant user  userWithRole1  role1 on res1",
	"deny  group groupWithRole1 role1 on res1",
	"policies:",
	"deny  user  userWithRole1  get,del res1",
	"grant role  role1          get,del res1",
	"grant group groupWithRole1 get,del res1",

	"grant user  user_allowed,group group_allowed   get,del res1",
	//"grant group group_allowed  get,del res1",
	"deny user  user_denied, group group_denied    get,del res1",
*/
//Multi Policies exist, both grant and deny user/group/role exist
func TestMats_IsAllowed_BothGrantAndDeny(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "DenyUser_GrantRole",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userWithRole1"}}},
					ServiceName: SERVICE_BOTH_GRANT_DENY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
		{
			Name:     "GrantGroup_DenyRole",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "groupWithRole1"}}},
					ServiceName: SERVICE_BOTH_GRANT_DENY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "GrantUser_DenyGroup",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_allowed"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group_denied"}}},
					ServiceName: SERVICE_BOTH_GRANT_DENY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
		{
			Name:     "DenyUser_GrantGroup",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_denied"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group_allowed"}}},
					ServiceName: SERVICE_BOTH_GRANT_DENY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-builtin-attribute"
//Policies' conditions contain builtin attribute
func TestLrg_IsAllowed_ConditionWithBuiltinAttribute(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "Request_time > 2017-09-04:grant",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_time1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		//"grant group  group_attri1    get,del res_request_user1 if request_user == 'admin'",
		{
			Name:     "Request_user==admin:grant",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "admin"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_user1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "Request_user==user_attri1:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group_func1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_user1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		//"grant group  group1    get,del res_request_groups1 if IsSubSet(request_groups,('group1','group2','group3')",
		{
			Name:     "Request_groups isSubSet: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_groups1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "Request_groups isSubSet: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group11"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_groups1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "Request_action==get:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_action1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "Request_action==del:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_action1",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "Request_resource==resource1:grant",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_resource1",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "Request_resource==resource2:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_resource2",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "Request_weekday==Monday",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_weekday1",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: time.Now().Weekday().String() == "Monday"},
			},
		},
		{
			Name:     "Request_year==2017:",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_year_equal_2017",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: time.Now().Year() == 2017},
			},
		},
		{
			Name:     "Request_year>2017:",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_year_greater_2017",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "Request_month==Nov",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_month_nov",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: time.Now().Month().String() == "November"},
			},
		},

		{
			Name:     "Request_day==14",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_attri1"}}},
					ServiceName: SERVICE_BUILTIN_ATTRIBUTE,
					Resource:    "res_request_day_14",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: time.Now().Day() == 14},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-complex"
//Policies are complex (Policy, rolePolicy,Condition, Grant,Deny…)
func TestLrg_IsAllowed_Complex_bug214(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{

		//"role-policies:",
		//"grant user user_complex1, user user_complex1A,user user_complex1B role_complex1",
		//"policies:",
		//"grant role role_complex1 get,del res_complex1 if request_user != 'user_complex1'",
		//"deny user user_complex1A del res_complex1",
		{
			Name:     "user_complex1_denied_in_condition",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_complex1"}}},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user_complex1B_allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_complex1B"}}},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user_complex1A_get_allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_complex1A"}}},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user_complex1A_del_denied_res_complex1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_complex1A"}}},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user_complex1A_get_allowed_res_complexAny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_complex1A"}}},
					ServiceName: SERVICE_COMPLEX,
					Resource:    "res_complex1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-complex-role
//Multi grant/deny role policies (role embeded) exist
func TestLrg_IsAllowed_Complex_Role_bug120(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "user1 get res1: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user1 del res2: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res2",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user1 del res3: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user11 del res3: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user11"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			}},
		{
			Name:     "user11,group1 del res3: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user11"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user11,group1 get res3: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user11"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			}},
		{
			Name:     "userAny,group1 get res3: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get res2: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get res1: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 del res1: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user2 get res3: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user22 get res1: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user22"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user22 del res3: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user22"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "del",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user22 get res3: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user22"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res3",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user1 get res9: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_ROLE,
					Resource:    "res9",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-with-resexpr"
//Multi grant/deny policies on resource expression
func TestLrg_IsAllowed_Complex_ResExpr_bug113(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "user1 get res: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user1 get res 2: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res 2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user1 get res*: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res*",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get res*: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res*",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get res2*: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res2&",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get *res2*: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "*res*",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get 22res222: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "22res222",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user2 get ?res?: allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "?res?",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "user1 get res-denied*: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res-denied*",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user1 get res-denied: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "res-denied",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user1 get 11res-denied11: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "11res-denied11",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "user2 get 11res-denied11: denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: SERVICE_COMPLEX_RESEXPR,
					Resource:    "11res-denied11",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-with-complex-principle-policy"
func TestMats_IsAllowed_ComplexPrincipleInPolicy(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "TwoUsersOneGroupInPolicyPrinciple-Deny-With-OneUser",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "OneUserAndTwoGroupsInPolicyPrinciple-Allow-with-userAndGroups",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group22"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "OneUserAndTwoGroupsInPolicyPrinciple-Deny-with-one-group",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "OnlyTwoGroupsInPolicyPrinciple-Allow-with-two-groups",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group3"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group33"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "OnlyTwoGroupsInPolicyPrinciple-Allow-with-two-groups",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group3"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group33"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "OneUserAndOneGroupInPolicyPrinciple-Allow-with-all",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "OneUserAndOneGroupInPolicyPrinciple-Allow-with-all-and-another-group",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "Two-OR-PrincipleInPolicyPrinciple-Deny-with-neither-matched",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group22"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_POILICY,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//comment this out as role policy does not support AND-Principal
//Policies are defined in check_prepare_test.go : "service-with-complex-principle_rolePolicy"
func _TestMats_IsAllowed_ComplexPrincipleInRolePolicy(t *testing.T) {
	URI_IS_ALLOWD := svcs.PolicyAtzPath + "is-allowed"
	data := &[]testutil.TestCase{
		{
			Name:     "TwoUsersOneGroupInRolePolicyPrinciple-Deny-With-OneUser",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "OneUserTwoGroupsInRolePolicyPrinciple-Allow-with-userAndGroups",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group22"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "OneUserTwoGroupsInRolePolicyPrinciple-Deny-with-one-group",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "OnlyTwoGroupsInRolePolicyPrinciple-Allow-with-two-groups",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group3"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group33"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GroupAndRoleInPolicyPrinciple-Allow-with-all",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group22"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group4"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res4",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GroupAndRoleInPolicyPrinciple-Deny-with-role",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group22"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res4",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "GroupAndRoleInPolicyPrinciple-Deny-with-group",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group4"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res4",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "UserAndRoleInPolicyPrinciple-Allow-with-all",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user4"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group3"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group33"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res4",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "UserAndRoleInPolicyPrinciple-Deny-with-role",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group3"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group33"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res4",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "UserAndRoleInPolicyPrinciple-Deny-with-user",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user4"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res4",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "GroupOrRoleInPolicyPrinciple-Allow-with-group",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "userAny"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group5"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res5",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
		{
			Name:     "GroupOrRoleInPolicyPrinciple-Allow-with-Role",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group2"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group22"}}},
					ServiceName: SERVICE_COMPLEX_PRINCIPLE_IN_ROLEPOILICY,
					Resource:    "res5",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policies are defined in check_prepare_test.go : "service-with-entity-principle"
//Test Entity is specified in principle in policy/rolepolicy
/*
	"role-policies:",
	"grant (entity schema://domain.name/path1), (entity schema://domain.name/path2) role1",
	//"grant (entity spiffe://domain.name/path1, entity spiffe://domain.name/ns/user1) role2", //multi entities is not supported when do adsrest

	"policies:",
	"grant role role1  get,del res1",
	"grant (group group1, entity spiffe://acme.com/9eebccd2-12bf-40a6-b262-65fe0487d453), role role1 get,del res2",
	"deny entity schema://domain.name/path2 get,del res2",
*/
func TestMats_IsAllowed_EntityPrinciple(t *testing.T) {
	data := &[]testutil.TestCase{
		{
			Name:     "GeneralEntity1_AccessRes1_Allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_ENTITY, Name: "schema://domain.name/path1"}}},
					ServiceName: SERVICE_WITH_ENTITY_PRINCIPLE,
					Resource:    "res1",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true},
			},
		},
		{
			Name:     "GeneralEntity1_AccessRes2_Allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_ENTITY, Name: "schema://domain.name/path1"}}},
					ServiceName: SERVICE_WITH_ENTITY_PRINCIPLE,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true},
			},
		},
		{
			Name:     "GeneralEntity2_AccessRes2_Denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_ENTITY, Name: "schema://domain.name/path2"}}},
					ServiceName: SERVICE_WITH_ENTITY_PRINCIPLE,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: 1},
			},
		},
		{
			Name:     "SpiffeEntity_AccessRes2_Denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_ENTITY, Name: "spiffe://acme.com/9eebccd2-12bf-40a6-b262-65fe0487d453"}}},
					ServiceName: SERVICE_WITH_ENTITY_PRINCIPLE,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: 3},
			},
		},
		{
			Name:     "SpiffeEntityAndGroup_AccessRes2_Allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}, {Type: adsapi.PRINCIPAL_TYPE_ENTITY, Name: "spiffe://acme.com/9eebccd2-12bf-40a6-b262-65fe0487d453"}}},
					ServiceName: SERVICE_WITH_ENTITY_PRINCIPLE,
					Resource:    "res2",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true},
			},
		},
		{
			Name:     "SpecialEntity_AccessRes3_Allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_ENTITY, Name: "special-schema.1+2://user1:pwd@domain1/path1/path-2/a"}}},
					ServiceName: SERVICE_WITH_ENTITY_PRINCIPLE,
					Resource:    "res3",
					Action:      "get",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true},
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
