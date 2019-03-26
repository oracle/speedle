//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"fmt"
	"testing"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/testutil"
	"github.com/oracle/speedle/testutil/msg"
	"github.com/oracle/speedle/testutil/param"
)

func getCmdTestDataForCreateService(sName string, sType string) testutil.TestCase {

	tmpData := testutil.TestCase{
		Name:     "TestCreateService-" + serviceName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_SERVICE,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_SERVICE(sName, sType),
			ExpectedMsg: msg.OUTPUT_SERVICE_CREATED(),
			OutputBody:  &pms.Service{},
			ExpectedBody: &pms.Service{
				Name: sName,
				Type: sType,
			},
		},
	}
	return tmpData
}

func getCmdTestDataForDeleteService(sName string) testutil.TestCase {

	tmpData := testutil.TestCase{
		Name:     "TestDeleteService",
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_DELETE_SERVICE,
		Data: &testutil.CmdTestData{
			Param:       param.DELETE_SERVICE(sName),
			ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
		},
	}
	return tmpData
}

func getCmdTestDataForCreatePolicyWthGroupName(sName string, policyName string, groupName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant group %s add,del res1", groupName)
	tmpData := testutil.TestCase{
		Name:     "TestCreatePolicy-" + policyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_POLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_POLICY(sName, policyName, pdl),
			ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
			OutputBody:  &pms.Policy{},
			ExpectedBody: &pms.Policy{
				Name:   policyName,
				Effect: pms.Grant,
				Permissions: []*pms.Permission{
					{
						Resource: "res1",
						Actions:  []string{"add", "del"},
					},
				},
				Principals: [][]string{{"group:" + groupName}},
			},
		},
		PostTestFunc: testutil.PostCreateGetPolicyTest,
	}
	return tmpData
}

func getCmdTestDataForCreatePolicyWthUserName(sName string, policyName string, userName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant user %s add,del res1", userName)
	tmpData := testutil.TestCase{
		Name:     "TestCreatePolicy-" + policyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_POLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_POLICY(sName, policyName, pdl),
			ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
			OutputBody:  &pms.Policy{},
			ExpectedBody: &pms.Policy{
				Name:   policyName,
				Effect: pms.Grant,
				Permissions: []*pms.Permission{
					{
						Resource: "res1",
						Actions:  []string{"add", "del"},
					},
				},
				Principals: [][]string{{"user:" + userName}},
			},
		},
		PostTestFunc: testutil.PostCreateGetPolicyTest,
	}
	return tmpData
}

func getCmdTestDataForCreatePolicyWthRoleName(sName string, policyName string, roleName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant role %s add,del res1", roleName)
	tmpData := testutil.TestCase{
		Name:     "TestCreatePolicy-" + policyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_POLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_POLICY(sName, policyName, pdl),
			ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
			OutputBody:  &pms.Policy{},
			ExpectedBody: &pms.Policy{
				Name:   policyName,
				Effect: pms.Grant,
				Permissions: []*pms.Permission{
					{
						Resource: "res1",
						Actions:  []string{"add", "del"},
					},
				},
				Principals: [][]string{{"role:" + roleName}},
			},
		},
		PostTestFunc: testutil.PostCreateGetPolicyTest,
	}
	return tmpData
}

func getCmdTestDataForCreatePolicyWthResName(sName string, policyName string, resName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant user user1 add,del %s", resName)
	tmpData := testutil.TestCase{
		Name:     "TestCreatePolicy-" + policyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_POLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_POLICY(sName, policyName, pdl),
			ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
			OutputBody:  &pms.Policy{},
			ExpectedBody: &pms.Policy{
				Name:   policyName,
				Effect: pms.Grant,
				Permissions: []*pms.Permission{
					{
						Resource: resName,
						Actions:  []string{"add", "del"},
					},
				},
				Principals: [][]string{{"user:user1"}},
			},
		},
		PostTestFunc: testutil.PostCreateGetPolicyTest,
	}
	return tmpData
}

func getCmdTestDataForCreateRolePolicyWthResName(sName string, rpolicyName string, resName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant user user1 role1 on %s", resName)
	tmpData := testutil.TestCase{
		Name:     "TestCreateRolePolicy-" + rpolicyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_ROLEPOLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_ROLEPOLICY(sName, rpolicyName, pdl),
			ExpectedMsg: msg.OUTPUT_ROLEPOLICY_CREATED(),
			OutputBody:  &pms.RolePolicy{},
			ExpectedBody: &pms.RolePolicy{
				Name:       rpolicyName,
				Roles:      []string{"role1"},
				Effect:     pms.Grant,
				Principals: []string{"user:user1"},
				Resources:  []string{resName},
			},
		},
		PostTestFunc: testutil.PostCreateGetRolePolicyTest,
	}
	return tmpData
}

func getCmdTestDataForCreateRolePolicyWthRoleName(sName string, rolePolicyName string, roleName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant user user1 %s on res1", roleName)
	tmpData := testutil.TestCase{
		Name:     "TestCreateRolePolicy-" + rolePolicyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_ROLEPOLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_ROLEPOLICY(sName, rolePolicyName, pdl),
			ExpectedMsg: msg.OUTPUT_ROLEPOLICY_CREATED(),
			OutputBody:  &pms.RolePolicy{},
			ExpectedBody: &pms.RolePolicy{
				Name:       rolePolicyName,
				Roles:      []string{roleName},
				Effect:     pms.Grant,
				Principals: []string{"user:user1"},
				Resources:  []string{"res1"},
			},
		},
		PostTestFunc: testutil.PostCreateGetRolePolicyTest,
	}
	return tmpData
}

func getCmdTestDataForCreatePolicyWthActionName(sName string, policyName string, actionName string) testutil.TestCase {

	pdl := fmt.Sprintf("grant user user1 %s,del res1", actionName)
	tmpData := testutil.TestCase{
		Name:     "TestCreatePolicy-" + policyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_CREATE_ROLEPOLICY,
		Data: &testutil.CmdTestData{
			Param:       param.CREATE_POLICY(sName, policyName, pdl),
			ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
			OutputBody:  &pms.Policy{},
			ExpectedBody: &pms.Policy{
				Name:   policyName,
				Effect: pms.Grant,
				Permissions: []*pms.Permission{
					{
						Resource: "res1",
						Actions:  []string{actionName, "del"},
					},
				},
				Principals: [][]string{{"user:user1"}},
			},
		},
		PostTestFunc: testutil.PostCreateGetPolicyTest,
	}
	return tmpData
}

func getCmdTestDataForGetPolicyWithGroupName(sName string, policyName string, groupName string) testutil.TestCase {
	tmpData := testutil.TestCase{
		Name:     "TestGetPolicy-" + policyName,
		Enabled:  true,
		Executer: testutil.NewCmdTest,
		Method:   testutil.METHOD_GET_POLICY,
		Data: &testutil.CmdTestData{
			Param:      "TO be init in PreTestFunc with id",
			OutputBody: &pms.Policy{},
			ExpectedBody: &pms.Policy{
				Name:   policyName,
				Effect: pms.Grant,
				Permissions: []*pms.Permission{
					{
						Resource: "res1",
						Actions:  []string{"add", "del"},
					},
				},
				Principals: [][]string{{"group:" + groupName}},
			},
		},
		PreTestFunc: func(data interface{}, context *testutil.TestContext) {
			cmdTD := data.(*testutil.CmdTestData)
			id, ok := context.NameIDMap[policyName]
			if ok {
				cmdTD.Param = param.GET_POLICY(sName, id)
				testutil.TestLog.Log(id)
			}
		},
		PostTestFunc: testutil.PostCreateGetPolicyTest,
	}
	return tmpData
}

//Create policy with group name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedGroupName(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedGroupName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, gName := range testutil.SpecialNames {
		policyName := fmt.Sprintf("policyWithGroupName_%d", i)
		data = append(data, getCmdTestDataForCreatePolicyWthGroupName(sName, policyName, gName))
		data = append(data, getCmdTestDataForGetPolicyWithGroupName(sName, policyName, gName))
	}
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}

//Create policy with user name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedUserName(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedUserName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, usrName := range testutil.SpecialNames {
		policyName := fmt.Sprintf("policyWithUserName_%d", i)
		data = append(data, getCmdTestDataForCreatePolicyWthUserName(sName, policyName, usrName))
	}
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}

//Create policy with role name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedRoleName(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedRoleName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, gName := range testutil.SpecialNames {
		policyName := fmt.Sprintf("policyWithRoleName_%d", i)
		data = append(data, getCmdTestDataForCreatePolicyWthRoleName(sName, policyName, gName))
	}
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}

//Create policy with action name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedActionName(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedActionName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, gName := range testutil.SpecialNames {
		policyName := fmt.Sprintf("policyWithAction_%d", i)
		data = append(data, getCmdTestDataForCreatePolicyWthActionName(sName, policyName, gName))
	}
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}

//Create policy with resource name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedResourceName(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedResourceName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, resName := range testutil.SpecialNames {
		policyName := fmt.Sprintf("policyWithResource_%d", i)
		data = append(data, getCmdTestDataForCreatePolicyWthResName(sName, policyName, resName))
	}
	data = append(data, getCmdTestDataForCreatePolicyWthResName(sName, "policyWithResource_comma", " res1,,res2"))
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}

//Create role-policy with resource name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedResourceNameInRolePolicy(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedResourceName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, resName := range testutil.SpecialNames {
		rpolicyName := fmt.Sprintf("rolePolicyWithResource_%d", i)
		data = append(data, getCmdTestDataForCreateRolePolicyWthResName(sName, rpolicyName, resName))
	}
	data = append(data, getCmdTestDataForCreatePolicyWthResName(sName, "rolepolicyWithResource_comma", " res1,,res2"))
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}

//Create role-policy with role name containing all kinds of special chars
func TestLrg_PDLWithSpecifiedRoleNameInRolePolicy(t *testing.T) {
	sName := "TestLrg_PDLWithSpecifiedResourceName"
	sType := "application"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	testutil.InitSpecialNames()
	data := []testutil.TestCase{}
	data = append(data, getCmdTestDataForCreateService(sName, sType))
	for i, roleName := range testutil.SpecialNames {
		rpolicyName := fmt.Sprintf("rolePolicyWithRole_%d", i)
		data = append(data, getCmdTestDataForCreateRolePolicyWthRoleName(sName, rpolicyName, roleName))
	}
	data = append(data, getCmdTestDataForDeleteService(sName))
	testutil.RunTestCases(t, &data, context)
}
