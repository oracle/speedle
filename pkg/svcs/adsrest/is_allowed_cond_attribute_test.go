//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

//+build runtime_test

package adsrest

import (
	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/testutil"
	"testing"
)

//Policies are defined in check_prepare_test.go : "service-condition-attribute"
//Policy's condition contains bool attribute
func TestMats_IsAllowed_Attri_Bool(t *testing.T) {

	data := &[]testutil.TestCase{
		//x=false allow; x=true deny; x=abc error
		//"grant user  user_bool_equal1   get,del res_equal1 if x == false",
		{
			Name:     "Condition(x==false), x=false: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_equal1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(x==false), x=true: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_equal1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		//x=false deny; x=true allow;
		//"grant user  user_bool_equal2   get,del res_equal2 if !x == false",
		{
			Name:     "Condition(!x==false), x=false: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(!x==false), x=true: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//x=false allow; x=true allow; x=abc error
		//"grant user  user_bool_notequal   get,del res_equal2 if x != true",
		{
			Name:     "Condition(x!=true), x=true: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(x!=true), x=false: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//s1=true,s2=false,s3=true allow
		//s1=true,s2=true,s3=true allow
		//s1=false,s2=false,s3=true allow
		//s1=true,s2=true,s3=false deny
		//s1=false,s2=true,s3=false deny
		//"grant user  user_bool_complex   get,del res_complex if s1 && !s2 || s3 == true",
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=true,s2=false,s3=true: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_complex"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_complex",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: false},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=true,s2=true,s3=true: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_complex"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_complex",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: true},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=false,s2=false,s3=true: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_complex"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_complex",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: false},
						{Name: "s2", Value: false},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=true,s2=true,s3=false: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_complex"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_complex",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: true},
						{Name: "s3", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},
			},
		},
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=false,s2=true,s3=false: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_complex"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_complex",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: false},
						{Name: "s2", Value: true},
						{Name: "s3", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Failed due to issue#162 in Gitlab
//Policy's condition contains bool attribute, and no space between the operators in bool expression
func TestLrg_IsAllowed_Attri_Bool_NoSpace_Between_Operator_bug162(t *testing.T) {

	data := &[]testutil.TestCase{
		//s1=true,s2=false,s3=true allow
		//"grant user  user_bool1   get,del res1 if s1&&!s2||s3==true",
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=true,s2=false,s3=true: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: false},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=true,s2=false,s3=true allow
		//"grant user  user_bool2   get,del res1 if s1&&(!s2)||s3==true",
		{
			Name:     "Condition(s1&&(!s2)||s3==true),s1=true,s2=false,s3=true: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: false},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=true,s2=false,s3=false allow
		//"grant user  user_bool3   get,del res1 if s1&&(!s2||s3)==true",
		{
			Name:     "Condition(s1&&(!s2||s3)==true),s1=true,s2=false,s3=false: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool3"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: false},
						{Name: "s3", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policy's condition contains Numeric attribute
func TestMats_IsAllowed_Attri_Num(t *testing.T) {

	data := &[]testutil.TestCase{
		//s1=20, s2=2, s3=3  allow
		//s1=20, s2=2, s3=2  deny
		//s1=20, s2=2, s3=-3  deny
		//s1=-22, s2=2, s3=-3  allow
		//s1=-22, s2=2, s3=3  deny
		//s1=12, s2=-2, s3=3  allow
		//s1=abc, s2=-2, s3=3  error  //TODO:
		//"grant user  user_num_equal   get,del res_num1 if (s1+5-s2*2)/3%4 == s3",
		{
			Name:     "Condition((s1+5-s2*2)/3%4 == s3),s1=20,s2=2,s3=3: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 20},
						{Name: "s2", Value: 2},
						{Name: "s3", Value: 3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition((s1+5-s2*2)/3%4 == s3),s1=20,s2=2,s3=2: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 20},
						{Name: "s2", Value: 2},
						{Name: "s3", Value: 2}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition((s1+5-s2*2)/3%4 == s3),s1=20,s2=2,s3=-3: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 20},
						{Name: "s2", Value: 2},
						{Name: "s3", Value: -3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition((s1+5-s2*2)/3%4 == s3),s1=-22,s2=2,s3=-3: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -22},
						{Name: "s2", Value: 2},
						{Name: "s3", Value: -3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition((s1+5-s2*2)/3%4 == s3),s1=-22,s2=2,s3=3: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -22},
						{Name: "s2", Value: 2},
						{Name: "s3", Value: 3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition((s1+5-s2*2)/3%4 == s3),s1=12,s2=-2,s3=3: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 12},
						{Name: "s2", Value: -2},
						{Name: "s3", Value: 3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=-10, s2=3, s3=-1 allow
		//s1=-10, s2=3, s3=1 deny
		//s1=0, s2=-3, s3=1 allow
		//s1=0, s2=-3, s3=-1 deny
		//s1=1.2, s2=-3.0, s3=1 allow
		//"grant user  user_num_equal2   get,del res_num2 if 5*(s1+5)%s2 == s3",
		{
			Name:     "Condition(5*(s1+5)%s2 == s3),s1=-10,s2=3,s3=-1: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -10},
						{Name: "s2", Value: 3},
						{Name: "s3", Value: -1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(5*(s1+5)%s2 == s3),s1=-10,s2=3,s3=1: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -10},
						{Name: "s2", Value: 3},
						{Name: "s3", Value: 1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(5*(s1+5)%s2 == s3),s1=0,s2=-3,s3=1: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 0},
						{Name: "s2", Value: -3},
						{Name: "s3", Value: 1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(5*(s1+5)%s2 == s3),s1=0,s2=-3,s3=-1: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 0},
						{Name: "s2", Value: -3},
						{Name: "s3", Value: -1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(5*(s1+5)%s2 == s3),s1=1.2,s2=-3.0,s3=1: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1.2},
						{Name: "s2", Value: -3.0},
						{Name: "s3", Value: 1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=0, s2=10, deny
		//s1=1, s2=10, allow
		//s1=abc, s2=10, erro //TODO:
		//"grant user  user_num_notequal  get,del res_num3 if (s1+5)*2 != s2",
		{
			Name:     "Condition((s1+5)*2 != s2),s1=0,s2=10: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num3",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 0},
						{Name: "s2", Value: 10}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition((s1+5)*2 != s2),s1=1,s2=10: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num3",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1},
						{Name: "s2", Value: 10}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=-5, s2=0, deny
		//s1=-2, s2=0, allow
		//s1=0, s2=0, allow
		//s1=-10, s2=-10, allow
		//"grant user  user_num_greater   get,del res_num4 if (s1+5)%2 > s2",
		{
			Name:     "Condition((s1+5)%2 > s2),s1=-5,s2=0: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -5},
						{Name: "s2", Value: 0}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition((s1+5)%2 > s2),s1=-2,s2=0: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -2},
						{Name: "s2", Value: 0}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition((s1+5)%2 > s2),s1=0,s2=0: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 0},
						{Name: "s2", Value: 0}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition((s1+5)%2 > s2),s1=-10,s2=-10: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -10},
						{Name: "s2", Value: -10}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=-5, s2=0, allow
		//s1=-20, s2=-20, allow
		//s1=30, s2=10.02, allow
		//"grant user  user_num_greater_equal   get,del res_num5 if (s1+5)/2.0 >= s2",
		{
			Name:     "Condition((s1+5)/2.0 >= s2),s1=-5,s2=0: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num5",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -5},
						{Name: "s2", Value: 0}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition((s1+5)/2.0 >= s2),s1=-20,s2=-20: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num5",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -20},
						{Name: "s2", Value: -20}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition((s1+5)/2.0 >= s2),s1=30,s2=10.02: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num5",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 30},
						{Name: "s2", Value: 10.02}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1=1, s2=0.4 deny
		//s1=1, s2=0.5 allow
		//s1="1", s2=".5" error //TODO
		//s1="a", s2=".5" error //TODO
		//"grant user  user_num_greater_equal2  get,del res_num6 if s1+s2 >= 1.41421356237309504880168872420969807856967187537694807317667974",
		{
			Name:     "Condition(s1+s2 >= 1.414...4),s1=1,s2=0.4: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num6",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1},
						{Name: "s2", Value: 0.4}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},

		{
			Name:     "Condition(s1+s2 >= 1.414...4),s1=1,s2=0.5: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_num_greater_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_num6",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1},
						{Name: "s2", Value: 0.5}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Failed due to issue#163 in Gitlab
//Policy's condition contains String attribute
func TestMats_IsAllowed_Attri_String(t *testing.T) {

	data := &[]testutil.TestCase{

		//s1='1', s2='23' allow
		//s1='1', s2='22' deny
		//s1='', s2='123' allow
		//s1=null, s2='123' allow
		//s1=1, s2=23 error //TODO
		// "grant user  user_str_equal1  get,del res_str1 if s1+s2 == '123'",
		{
			Name:     "Condition(s1+s2 == '123'),s1='1',s2='23': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_equal1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "1"},
						{Name: "s2", Value: "23"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1+s2 == '123'),s1='1',s2='22': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_equal1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "1"},
						{Name: "s2", Value: "22"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1+s2 == '123'),s1='',s2='123': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_equal1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: ""},
						{Name: "s2", Value: "123"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//Failed with issue#166 in Gitlab
		/*{
				Name:       "Condition(s1+s2 == '123'),s1=null,s2='123': allow",
				Executer: testutil.NewRestTestExecuter,
		Method: testutil.METHOD_IS_ALLOWED,
				Data: &testutil.RestTestData{
		URI:        URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_str_equal1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						//&JsonAttribute{Name: "s1", Value: nil},
						&JsonAttribute{Name: "s2", Value: "123"}},
				},
				ExpectedStatus: 200, //Null value is not allowed
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false},},
			},*/

		//s1='1', s2='1' deny
		//s1='1', s2='-1' allow
		//s1="1", s2=1 error  //TODO
		//s1=nil, s2='nil' error  //TODO. Block by issue166
		//s1="", s2='nil' allow
		//"grant user  user_str_notequal  get,del res_str2 if s1 != s2",
		{
			Name:     "Condition(s1 != s2),s1='1',s2='1': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "1"},
						{Name: "s2", Value: "1"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1 != s2),s1='1',s2='-1': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "1"},
						{Name: "s2", Value: "-1"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1 != s2),s1='',s2='nil': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: ""},
						{Name: "s2", Value: "nil"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//s1=abc, s2=abc allow
		//s1=abc, s2="123" deny
		//"grant user  user_str_regx  get,del res_str3 if s1+s2=~'.*abc'",
		{
			Name:     "Condition(s1+s2=~'.*abc$'),s1='abc',s2='abc': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_regx"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str3",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1+s2=~'.*abc$'),s1='abc',s2='123': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_regx"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str3",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "123"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		//s1='abc', s2='Abc' allow
		//s1='abc', s2='abc' deny
		//s1='abc', s2='' allow
		//"grant user  user_str_greater  get,del res_str4 if s1>s2",
		{
			Name:     "Condition(s1>s2),s1='abc',s2='Abc': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "Abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1>s2),s1='abc',s2='abc': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1>s2),s1='abc',s2='': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str4",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: ""}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//s1='abc', s2='Abc' allow
		//s1='abc', s2='abc' allow
		//s1='null', s2='' allow
		//"grant user  user_str_greater_equal  get,del res_str5 if s1>=s2",
		{
			Name:     "Condition(s1>=s2),s1='abc',s2='Abc': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str5",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "Abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1>=s2),s1='abc',s2='abc': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str5",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1>=s2),s1='null',s2='': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str5",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "null"},
						{Name: "s2", Value: ""}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//s1='abc', s2='Abc' deny
		//s1='abc', s2='abc' allow
		//s1='null', s2='' deny  //TODO, use nil or "null"
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-02T15:04:05-07:00' allow
		//"grant user  user_str_less_equal  get,del res_str6 if s1<=s2",
		{
			Name:     "Condition(s1<=s2),s1='abc',s2='Abc': Deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_less_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str6",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "Abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1<=s2),s1='abc',s2='abc': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_less_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str6",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},

		//s1='abc', s2='Abc' deny
		//s1='abc', s2='abc' deny
		//s1='null', s2='' deny   //TODO: should add this test cases once issue#166 fixed
		//s1='a bc', s2='a1bc' allow
		//"grant user  user_str_less  get,del res_str7 if s1<s2",
		{
			Name:     "Condition(s1<s2),s1='abc',s2='Abc': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_less"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str7",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "Abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1<s2),s1='abc',s2='abc': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_less"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str7",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1<s2),s1='a bc',s2='a1bc': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_str_less"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_str7",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "a bc"},
						{Name: "s2", Value: "a1bc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policy's condition contains Datetime attribute
func TestMats_IsAllowed_Attri_Datetime(t *testing.T) {

	data := &[]testutil.TestCase{

		//------------------------Datetime-------------------------------
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-02T15:04:05-09:00' deny
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-22T15:04:05-07:00' deny
		//"grant user  user_date_equal  get,del res1 if s1=s2",
		{
			Name:     "Condition(s1==s2),s1='2017-01-02T15:04:05-07:00',s2='2017-01-02T15:04:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1==s2),s1='2017-01-02T15:04:05-07:00',s2='2017-01-02T15:04:05-09:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1==s2),s1='2017-01-02T15:04:05-07:00',s2='2017-01-22T15:04:05-07:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-22T15:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},

		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' deny
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' allow
		//"grant user  user_date_notequal  get,del res1 if s1!=s2",
		{
			Name:     "Condition(s1!=s2),s1='2017-01-02T15:04:05+07:00',s2='2017-01-02T15:04:05+07:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05+07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05+07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1!=s2),s1='2017-01-02T15:04:05-07:00',s2='2017-01-02T15:04:05+07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05+07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' deny
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' deny
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00' allow
		//"grant user  user_date_greater  get,del res1 if s1>s2",
		{
			Name:     "Condition(s1>s2),s1='2017-01-02T15:04:05-07:00',s2='2017-01-02T15:04:05-07:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1>s2),s1='2017-01-02T15:04:05-07:00',s2='2017-01-02T15:04:05-09:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1>s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:04:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1>s2),s1='2017-01-02T05:05:05-09:00',s2='2017-01-02T05:05:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1='2017-01-02T15:05:05-07:00', s2='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:05:05-09:00', s2='2017-01-02T15:05:05-07:00' allow
		//s1='2017-01-02T15:05:04-07:00', s2='2017-01-02T15:04:05-07:00' deny
		//"grant user  user_date_greater_equal  get,del res1 if s1>=s2",
		{
			Name:     "Condition(s1>=s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:05:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1>=s2),s1='2017-01-02T05:05:05-09:00',s2='2017-01-02T05:05:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-09:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1>=s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:05:04-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_greater_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:04-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' allow
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00' deny
		//"grant user  user_date_less_equal  get,del res1 if s1<=s2",
		{
			Name:     "Condition(s1<=s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:05:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_less_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1<=s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:05:05-09:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_less_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1<=s2),s1='2017-01-02T05:05:04-07:00',s2='2017-01-02T05:05:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_less_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:04-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' deny
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00'allow
		//s1='2017-01-02T15:04:04-07:00', s2='2017-01-02T15:04:05-07:00' allow
		//"grant user  user_date_less  get,del res1 if s1<s2",
		{
			Name:     "Condition(s1<=s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:05:05-07:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_less"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(s1<=s2),s1='2017-01-02T05:05:05-07:00',s2='2017-01-02T05:05:05-09:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_less"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(s1<=s2),s1='2017-01-02T05:05:04-07:00',s2='2017-01-02T05:05:05-07:00': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_less"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T05:05:04-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T05:05:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Datetime is not in valid format
func TestLrg_IsAllowed_Negative_Datetime_InvalidFormat(t *testing.T) {

	data := &[]testutil.TestCase{
		//s1='2017-02-29T15:04:05-07:00', s1='2017-01-22T15:04:05-07:00' error
		//"grant user  user_date_equal  get,del res1 if s1=s2",
		{
			Name:     "Condition(s1=s2),s1='2017-02-29T15:04:05-07:00',s2='2017-01-02T15:04:05-07:00': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-02-29T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-22T15:04:05-07:00"}},
				},
				ExpectedStatus: 400,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.ERROR_IN_EVALUATION), ErrorMessage: "Invalid datetime value"}},
		},
		{
			Name:     "Condition(s1=s2),invalid Attribute type : error",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "date-time", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "date-time", Value: "2017-01-02T15:04:05-07:00"}},
				},
				ExpectedStatus: 400,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.ERROR_IN_EVALUATION), ErrorMessage: "Invalid attribute type"}},
		},
		{
			Name:     "Condition(s1=s2),invalid datetime value ': error",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_date_equal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "date-time", Value: "2017-01-02T15:04:05"},
						{Name: "s2", Type: "date-time", Value: "2017-01-02T15:04:05"}},
				},
				ExpectedStatus: 400,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.ERROR_IN_EVALUATION), ErrorMessage: "Invalid datetime value"}},
		},
	}
	testutil.RunTestCases(t, data, nil)

}

//Policy's condition is string in (string_array)
func TestMats_IsAllowed_Attri_InArray_String(t *testing.T) {

	data := &[]testutil.TestCase{

		//------------------------Array-------------------------------
		//x='a' allow
		//x='b' deny
		//x='1' allow
		//x='' deny
		//x=1 error
		//"grant user  user_array_str  get,del res1 if x in ('a','c','1')",
		{
			Name:     "Condition(x in ('a','c','1')),x='a': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: "a"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(x in ('a','c','1')),x='b': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: "b"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(x in ('a','c','1')),x='1': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: "1"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(x in ('a','c','1')),x='': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: ""}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(x in ('a','c','1')),x=1: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: 1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		//x='a' allow
		//x='b' deny
		//"grant user  user_array_str_single  get,del res1 if x in ('a')",
		/* Comment below test cases due to Issue#176 is not a valid bug
			{
				Name:       "Condition(x in ('a')),x='a': allow",
				Executer: testutil.NewRestTestExecuter,
		Method: testutil.METHOD_IS_ALLOWED,
				Data: &testutil.RestTestData{
		URI:        URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_array_str_single"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						&JsonAttribute{Name: "x", Value: "a"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},},
			},
			{
				Name:       "Condition(x in ('a')),x='b': deny",
				Executer: testutil.NewRestTestExecuter,
		Method: testutil.METHOD_IS_ALLOWED,
				Data: &testutil.RestTestData{
		URI:        URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_array_str_single"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						&JsonAttribute{Name: "x", Value: "b"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},},
			},
		*/
		//"grant user  user_array_str_single1  get,del res1 if s1 in s2",
		{
			Name:     "Condition( s1 in s2)),s1='a',s2='a': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_single1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "a"},
						{Name: "s2", Value: []string{"a"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in s2)),s1='',s2=Empty: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_single1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: ""},
						{Name: "s2", Value: []string{}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition( s1 in s2)),s1='',s2='a': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_single1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: ""},
						{Name: "s2", Value: []string{"a"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
	}
	testutil.RunTestCases(t, data, nil)

}

//Policy's condition is num in (num_array)
func TestMats_IsAllowed_Attri_InArray_Num(t *testing.T) {

	data := &[]testutil.TestCase{

		//x=1 allow
		//x=0 deny
		//"grant user  user_array_num_single  get,del res1 if x in (1)",
		//Comment below test cases due to Issue#176 is not a valid bug
		/*
					{
						Name:       "Condition(x in (1)),x=1: allow",
						Executer: testutil.NewRestTestExecuter,
			Method: testutil.METHOD_IS_ALLOWED,
						Data: &testutil.RestTestData{
			URI:        URI_IS_ALLOWD,
						InputBody: &JsonContext{
							Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_array_num_single"}}},
							ServiceName: SERVICE_COND_ATTRIBUTE,
							Resource:    "res1",
							Action:      "get",
							Attributes: []*JsonAttribute{
								&JsonAttribute{Name: "x", Value: 1}},
						},
						ExpectedStatus: 200,
						OutputBody:     &IsAllowedResponse{},
						ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},},
					},
					{
						Name:       "Condition(x in (1)),x=2: deny",
						Executer: testutil.NewRestTestExecuter,
			Method: testutil.METHOD_IS_ALLOWED,
						Data: &testutil.RestTestData{
			URI:        URI_IS_ALLOWD,
						InputBody: &JsonContext{
							Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_array_num_single"}}},
							ServiceName: SERVICE_COND_ATTRIBUTE,
							Resource:    "res1",
							Action:      "get",
							Attributes: []*JsonAttribute{
								&JsonAttribute{Name: "x", Value: 2}},
						},
						ExpectedStatus: 200,
						OutputBody:     &IsAllowedResponse{},
						ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},},
					},*/
		//x=1 allow
		//x=2 allow
		//x=2.0 allow
		//x=-2.1 allow
		//x=3.5 deny
		//x=3.567 allow
		//x='1' error
		//"grant user  user_array_num  get,del res_in2 if x in (1,2.0,-2.1,3.567)",
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x=1, allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: 1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x=2, allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: 2}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x=2.0, allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: 2.0}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x=-2.1, allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: -2.1}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x=3.5,deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: 3.5}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x=3.567,allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: 3.567}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( x in (1,2.0,-2.1,3.567)),x='1',deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_num"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: '1'}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
	}
	testutil.RunTestCases(t, data, nil)

}

//Policy's condition is bool in (boolean_array)
func TestMats_IsAllowed_Attri_InArray_Boolean(t *testing.T) {

	data := &[]testutil.TestCase{
		//"grant user  user_array_str_multi  get,del res1 if s1 in (s2, s3)",
		{
			Name:     "Condition( s1 in (s2, s3)),s1='true', s2='true', s3='true': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: true},
						{Name: "s2", Value: true},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1='false', s2='false', s3='false': Allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: false},
						{Name: "s2", Value: false},
						{Name: "s3", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1='false', s2='true', s3='false': Allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: false},
						{Name: "s2", Value: true},
						{Name: "s3", Value: false}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1='false', s2='true', s3='true': Deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: false},
						{Name: "s2", Value: true},
						{Name: "s3", Value: true}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//Policy's condition is S1 in (S2, S3), and S* may be any type
func TestMats_IsAllowed_Attri_InArray_All(t *testing.T) {

	data := &[]testutil.TestCase{

		//s1='1', s2='1', s3='2' allow
		//s1='1', s2='2', s3='2' deny
		//s1=1, s2=1, s3=2 allow
		//s1=1, s2=2, s3=2 deny
		//s1=1.1, s2=1.1, s3=2 allow
		//s1=-2.0, s2=-2, s3=3 allow
		//s1=-2, s2=-2.0, s3=3 allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00', s3='2017-01-02T15:04:05-09:00' allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-10:00', s3='2017-01-02T15:04:05-09:00' deny
		//"grant user  user_array_str_multi  get,del res1 if s1 in (s2, s3)",
		{
			Name:     "Condition( s1 in (s2, s3)),s1='1', s2='1', s3='2': allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "1"},
						{Name: "s2", Value: "1"},
						{Name: "s3", Value: "2"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1='1', s2='2', s3='2': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "1"},
						{Name: "s2", Value: "2"},
						{Name: "s3", Value: "2"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1=1, s2=1, s3=2: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1},
						{Name: "s2", Value: 1},
						{Name: "s3", Value: 2}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1=1, s2=2, s3=2: deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1},
						{Name: "s2", Value: 2},
						{Name: "s3", Value: 2}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1=1.1, s2=1.1, s3=2: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: 1.1},
						{Name: "s2", Value: 1.1},
						{Name: "s3", Value: 2}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1=-2.0, s2=-2, s3=3: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -2.0},
						{Name: "s2", Value: -2},
						{Name: "s3", Value: 3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1=-2, s2=-2.0, s3=3: allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: -2},
						{Name: "s2", Value: -2.0},
						{Name: "s3", Value: 3}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//Failed by Issue#179
		{
			Name:     "Condition( s1 in (s2, s3)),s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00', s3='2017-01-02T15:04:05-09:00', allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s3", Type: "datetime", Value: "2017-01-02T15:04:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition( s1 in (s2, s3)),s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-10:00', s3='2017-01-02T15:04:05-09:00', deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"},
						{Name: "s2", Type: "datetime", Value: "2017-01-02T15:04:05-10:00"},
						{Name: "s3", Type: "datetime", Value: "2017-01-02T15:04:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
	}
	testutil.RunTestCases(t, data, nil)

}

//Policy's condition is isSubSet(S1,S2),  and S* may be any type
func TestMats_IsAllowed_Attri_IsSubset(t *testing.T) {

	data := &[]testutil.TestCase{
		//"grant user  user_array_issubset1  get,del res_subset if IsSubSet(s1,s2)",
		{
			Name:     "Condition(IsSubSet(s1,s2))-Num:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []int{1, 2}},
						{Name: "s2", Value: []int{1, 5, 2, 3}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-Num:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []int{1, 2}},
						{Name: "s2", Value: []int{1, 5, -2, 3}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-Single-Num:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []int{-2}},
						{Name: "s2", Value: []int{1, 5, -2, 3}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-String:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{"b", "a"}},
						{Name: "s2", Value: []string{"a", "b", "c"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-Single string:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{"b"}},
						{Name: "s2", Value: []string{"a", "b", "c"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-Single string-2:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{"b"}},
						{Name: "s2", Value: []string{"b"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-S1 is empty:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{}},
						{Name: "s2", Value: []string{"b"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-S2 is empty:deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{"a", "b"}},
						{Name: "s2", Value: []string{}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-Date:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: []string{"2017-01-02T15:04:05-07:00", "2017-01-02T15:04:05-09:00"}},
						{Name: "s2", Type: "datetime", Value: []string{"2017-01-02T15:04:05-07:00", "2017-01-02T15:04:05-09:00"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(IsSubSet(s1,s2))-Single Date:allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset1"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Type: "datetime", Value: []string{"2017-01-02T15:04:05-07:00"}},
						{Name: "s2", Type: "datetime", Value: []string{"2017-01-02T15:04:05-07:00", "2017-01-02T15:04:05-09:00"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		//"grant user  user_array_issubset2  get,del res_subset if IsSubSet(s1,('a1','b1'))",		{
		{
			Name:     "Condition(IsSubSet(s1,('a1','b1','c1'))),s1=('a1','b1'):allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_issubset2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_subset",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{"a1", "b1"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//An array can't be used in left hand of in operator
func TestLrg_IsAllowed_Negative_LeftHandOfInIsArray(t *testing.T) {

	data := &[]testutil.TestCase{
		//s1=('a','b'),s2=('a','c'),s3=('b','d'),allow
		//"grant user  user_array_str_multi  get,del res1 if s1 in (s2, s3)",

		{
			Name:     "Condition( s1 in (s2, s3)),s1=('a','b'),s2=('a','c'),s3=('b','d'),allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_str_multi"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res1",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: []string{"a", "b"}},
						{Name: "s2", Value: []string{"a", "c"}},
						{Name: "s3", Value: []string{"b", "d"}}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}

//disable temporarily for bug#179.
//Policy's condition is datetime in (datetime_array)
func TestMats_IsAllowed_Attri_In_Datetime(t *testing.T) {

	data := &[]testutil.TestCase{

		//x='2017-01-02T15:04:05-07:00' allow
		//x='2017-01-02T15:04:05-17:00' deny
		//"grant user  user_array_date  get,del res_in2 if x in ('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00')",
		{
			Name:     "Condition(user_array_date),x in arary,allow",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_date"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Type: "datetime", Value: "2017-01-02T15:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)}},
		},
		{
			Name:     "Condition(user_array_date),x NOT in arary,deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_array_date"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Type: "datetime", Value: "2017-01-02T15:04:05-17:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		//x='2017-01-02T15:04:05-07:00' allow
		//x='2017-01-02T15:04:05-17:00' deny
		//"grant user  user_array_date_single  get,del res_in2 if x in ('2017-01-02T15:04:05-07:00')",
		/* Comment below two test cases due to issue#176 is not a valid bug
			{
				Name:       "Condition(user_array_date),x in arary with 1 element,allow",
				Executer: testutil.NewRestTestExecuter,
		Method: testutil.METHOD_IS_ALLOWED,
				Data: &testutil.RestTestData{
		URI:        URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_array_date_single"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						&JsonAttribute{Name: "x",Type: "datetime", Value: "2017-01-02T15:04:05-07:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: true, Reason: int32(adsapi.GRANT_POLICY_FOUND)},},
			},
			{
				Name:       "Condition(user_array_date),x NOT in arary with 1 element,deny",
				Executer: testutil.NewRestTestExecuter,
		Method: testutil.METHOD_IS_ALLOWED,
				Data: &testutil.RestTestData{
		URI:        URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject: &JsonSubject{Principals: []*JsonPrincipal{&JsonPrincipal{Type: adsapi.PRINCIPAL_TYPE_USER, Name:  "user_array_date_single"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_in2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						&JsonAttribute{Name: "x",Type: "datetime",  Value: "2017-01-02T15:04:05-09:00"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)},},
			},*/

	}
	testutil.RunTestCases(t, data, nil)
}

//Failed due to issue#163 in Gitlab
//Boolean attribute is invalid (string or num)
func TestLrg_IsAllowed_Negative_Attri_Bool_bug163(t *testing.T) {

	data := &[]testutil.TestCase{

		//x=abc error
		//"grant user  user_bool_equal2   get,del res_equal2 if !x == false",
		{
			Name:     "Condition(!x==false), x='abc': Error",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_equal2"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.ERROR_IN_EVALUATION), ErrorMessage: "Invalid attribute value"}},
		},
		// x='0' error
		//"grant user  user_bool_notequal   get,del res_equal2 if x != true",
		{
			Name:     "Condition(x!=true), x='0': deny",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_notequal"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_equal2",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "x", Value: '0'}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.NO_APPLICABLE_POLICIES)}},
		},
		//s1=abc,s2=true,s3=abc error
		//s1=False,s2=True,s3=False error
		//"grant user  user_bool_complex   get,del res_complex if s1 && !s2 || s3 == true",
		{
			Name:     "Condition(s1&&!s2||s3==true),s1=abc,s2=true,s3=abc: error",
			Executer: testutil.NewRestTestExecuter,
			Method:   testutil.METHOD_IS_ALLOWED,
			Data: &testutil.RestTestData{
				URI: URI_IS_ALLOWD,
				InputBody: &JsonContext{
					Subject:     &JsonSubject{Principals: []*JsonPrincipal{{Type: adsapi.PRINCIPAL_TYPE_USER, Name: "user_bool_complex"}}},
					ServiceName: SERVICE_COND_ATTRIBUTE,
					Resource:    "res_complex",
					Action:      "get",
					Attributes: []*JsonAttribute{
						{Name: "s1", Value: "abc"},
						{Name: "s2", Value: true},
						{Name: "s3", Value: "abc"}},
				},
				ExpectedStatus: 200,
				OutputBody:     &IsAllowedResponse{},
				ExpectedBody:   &IsAllowedResponse{Allowed: false, Reason: int32(adsapi.ERROR_IN_EVALUATION), ErrorMessage: "Invalid parameter value"}},
		},
	}
	testutil.RunTestCases(t, data, nil)
}
