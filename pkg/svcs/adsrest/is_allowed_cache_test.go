//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
//+build runtime_cache_test

package adsrest

import (
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/testutil"

	"github.com/oracle/speedle/api/pms"
)

var POLICY_RELOAD_TIME = 500 //ms
/**
Allowed result changed due to deny policy added/removed
  test1:Policy1 exist, allow=true;
  test2:Policy2 Added, allowed=false
  test3:Policy2 removed, allow=true
  test4:Policy1 removed, allow=false
*/
func TestMats_Cache_With_Policy_Added(t *testing.T) {
	appName := "TestCacheWithPolicyAdded"

	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user1"}},
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
					Principals: [][]string{{"user:user1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step4-AddPolicy2 to deny ",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
				ExpectedStatus: 201,
				InputBody: &pms.Policy{
					Name:   "policy2",
					Effect: pms.Deny,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"group:group1"}},
				},
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy2",
					Effect: pms.Deny,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"group:group1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step5-Check user is denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.DENY_POLICY_FOUND)}},
		},
		{
			Name:     "Step5-DeletePolicy2",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy/policy2",
				ExpectedStatus: 204,
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step6-Check user is allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step7-DeletePolicy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy/policy1",
				ExpectedStatus: 204,
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step8-Check user is denied",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}, {Type: adsapi.PRINCIPAL_TYPE_GROUP, Name: "group1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step9-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to policy Condition changed
  test1:Policy1 exist, allow=true;
  test2:Policy1 Removed, and create Policy1 again with Condition changed. allowed=false
  test3:Policy1 removed, allow=false
*/
func TestMats_Cache_With_Policy_Condition_Changed(t *testing.T) {
	appName := "TestCacheWithPolicyConditionChanged"

	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(Condition:age>10) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user1"}},
					Condition:  "age>10",
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
					Principals: [][]string{{"user:user1"}},
					Condition:  "age>10",
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is allowed (age=15)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 15}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step4-DeletePolicy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy/policy1",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "Step5-AddPolicy1(Condition:age>20) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user1"}},
					Condition:  "age>20",
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
					Principals: [][]string{{"user:user1"}},
					Condition:  "age>20",
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step6-Check user is Allowed (age=25)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 25}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step7-Check user is Denied (age=15)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 15}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "Step8-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to policy Principle changed
*/
func TestMats_Cache_With_Policy_Subject_Changed(t *testing.T) {
	appName := "TestCacheWithPolicyConditionChanged"

	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(user:user1) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user1"}},
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
					Principals: [][]string{{"user:user1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is allowed (user1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step4-DeletePolicy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy/policy1",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "Step5-AddPolicy1(principle=user2) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user2"}},
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
					Principals: [][]string{{"user:user2"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step6-Check user is Allowed (user2)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user2"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step7-Check user is Denied (user1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step8-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to policy Permission resource changed
*/
func TestMats_Cache_With_Policy_Permission_Res_Changed(t *testing.T) {
	appName := "TestCacheWithPolicyPermissionResChanged"

	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(res1) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user1"}},
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
					Principals: [][]string{{"user:user1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is allowed (user1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step4-DeletePolicy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy/policy1",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "Step5-AddPolicy1(res2) again",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
				ExpectedStatus: 201,
				InputBody: &pms.Policy{
					Name:   "policy1",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"user:user1"}},
				},
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy1",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"read", "write"},
						},
					},
					Principals: [][]string{{"user:user1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step6-Check user is Allowed (res2)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res2",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step7-Check user is Denied (res1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step8-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to policy Permission Actions changed
*/
func TestMats_Cache_With_Policy_Permission_Action_Changed(t *testing.T) {
	appName := "TestCacheWithPolicyPermissionActionChanged"

	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(res1) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"user:user1"}},
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
					Principals: [][]string{{"user:user1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is allowed (action:read)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step4-DeletePolicy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy/policy1",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "Step5-AddPolicy1(action:read removed) again",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
				ExpectedStatus: 201,
				InputBody: &pms.Policy{
					Name:   "policy1",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"write"},
						},
					},
					Principals: [][]string{{"user:user1"}},
				},
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy1",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res2",
							Actions:  []string{"write"},
						},
					},
					Principals: [][]string{{"user:user1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step6-Check user is Allowed (action:write)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res2",
					Action:      "write",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step7-Check user is Denied (action:read)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step8-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to role-policy added/removed
  test1:Policy1(role granted) exist, allow=false;
  test2:Role-Policy1 Added, allowed=true
  test3:Role-Policy2 Added, allowed=false
  test4:Role-Policy2 Removed, allowed=true
  test5:Role-Policy1 Removed, allowed=false
*/
func TestMats_Cache_With_RolePolicy_Added(t *testing.T) {
	appName := "TestCacheWithPolicyPermissionResChanged"

	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(role1) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"role:role1"}},
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
					Principals: [][]string{{"role:role1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is denied (user1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step4-Add RolePolicy (allowed)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy",
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step5-Check user is allowed (user1) with grant rolePolicy added",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step6-Add RolePolicy (denied)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy",
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy2",
					Effect:     pms.Deny,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy2",
					Effect:     pms.Deny,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step7-Check user is denied (user1) with denied rolePolicy added",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step8-Delete the denied rolePolicy",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy/role-policy2",
				ExpectedStatus: 204,
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step9-Check user is allowed (user1) with denied rolePolicy deleted",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step10-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to role-policy subject changed
*/
func TestMats_Cache_With_RolePolicy_Subject_Changed(t *testing.T) {
	appName := "TestMatsWithRolePolicySubjectChanged"
	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(role1) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"role:role1"}},
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
					Principals: [][]string{{"role:role1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is denied (user1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step4-Add RolePolicy (allowed)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy",
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step5-Check user is allowed (user1) with grant rolePolicy added",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step6-Delete the rolePolicy",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy/role-policy1",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "Step6-Add RolePolicy1 again with subject changed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy",
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Deny,
					Roles:      []string{"role1"},
					Principals: []string{"user:user2"},
					Resources:  []string{"res1"},
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Deny,
					Roles:      []string{"role1"},
					Principals: []string{"user:user2"},
					Resources:  []string{"res1"},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step7-Check user is denied (user1) with changed role-policy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step8-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

/**
"Allowed" result changed due to role-policy Condition changed
*/
func TestMats_Cache_With_RolePolicy_Condition_Changed(t *testing.T) {
	appName := "TestMatsWithRolePolicyConditionChanged"
	data := &[]testutil.TestCase{
		{
			Name:     "Step1-AddService",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service",
				ExpectedStatus: 201,
				InputBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: appName,
					Type: pms.TypeApplication,
				},
			},
		},
		{
			Name:     "Step2-AddPolicy1(role1) allowed",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/policy",
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
					Principals: [][]string{{"role:role1"}},
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
					Principals: [][]string{{"role:role1"}},
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step3-Check user is denied (user1)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 15}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step4-Add RolePolicy (condition: age>10)",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy",
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
					Condition:  "age>10",
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
					Condition:  "age>10",
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step5-Check user is allowed (age=15) with grant rolePolicy added",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 15}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step6-Delete the rolePolicy",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy/role-policy1",
				ExpectedStatus: 204,
			},
		},
		{
			Name:     "Step7-Add RolePolicy1 again with condition:age>20",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName + "/role-policy",
				ExpectedStatus: 201,
				InputBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
					Condition:  "age>20",
				},
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       "role-policy1",
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
					Condition:  "age>20",
				},
			},
		},
		testutil.GetTestData_Sleep(POLICY_RELOAD_TIME),
		{
			Name:     "Step7-Check user is denied (user1) with changed role-policy1",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 15}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Step7-Check user is allowed (user1) with age=25",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: testutil.URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user1"}}},
					ServiceName: appName,
					Resource:    "res1",
					Action:      "read",
					Attributes: []*JsonAttribute{
						{Name: "age", Value: 25}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Step9-Delete Service",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.RestTestData{
				URI:            testutil.URI_POLICY_MGMT + "service/" + appName,
				ExpectedStatus: 204,
			},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
