//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"testing"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/testutil"
	"github.com/oracle/speedle/testutil/msg"
	"github.com/oracle/speedle/testutil/param"
)

//Create/Get/Delete service
func TestMats_Service(t *testing.T) {

	sName := "k8s"
	sType := pms.TypeApplication

	sName1 := "k8s1"
	sType1 := pms.TypeApplication

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateService1",
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
		},
		{
			Name:     "TestCreateService2",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_SERVICE(sName1, sType1),
				ExpectedMsg: msg.OUTPUT_SERVICE_CREATED(),
				OutputBody:  &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName1,
					Type: sType1,
				},
			},
		},
		{
			Name:     "TestGetService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.CmdTestData{
				Param:      param.GET_SERVICE(sName),
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName,
					Type: sType,
				},
			},
		},
		{
			Name:     "TestGetService2",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.CmdTestData{
				Param:      param.GET_SERVICE(sName1),
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName1,
					Type: sType1,
				},
			},
		},
		{
			Name:     "TestGetServiceAll",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_QUERY_SERVICE,
			Data: &testutil.CmdTestData{
				Param:      param.GET_SERVICE_ALL(),
				OutputBody: &[]pms.Service{},
				ExpectedBody: &[]pms.Service{
					{
						Name: sName,
						Type: sType,
					},
					{
						Name: sName1,
						Type: sType1,
					},
				},
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
		{
			Name:     "TestDeleteService2",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName1),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName1),
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Create service with JSON file in which only service is specified
func TestMats_CreateServiceOnlyByJsonFile(t *testing.T) {
	sName := "serviceOnly"
	sType := pms.TypeApplication
	jsonFile := "/tmp/service.json"
	//pdlFile := "/tmp/pdl.txt"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
		FileName:      jsonFile,
	}

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateServiceOnly",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.CmdTestData{
				Param: param.CREATE_SERVICE_WITH_JSONFILE(jsonFile),
				FileContent: &pms.Service{
					Name: sName,
					Type: sType,
				},
				ExpectedMsg: msg.OUTPUT_SERVICE_CREATED(),
				OutputBody:  &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName,
					Type: sType,
				},
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				tmpService := cmdTD.FileContent.(*pms.Service)
				testutil.GenerateJsonFileWithService(context.FileName, tmpService)
			},
		},
		{
			Name:     "TestDeleteServiceOnly",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE("serviceOnly"),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED("serviceOnly"),
			},
		},
	}

	testutil.RunTestCases(t, data, context)
}

//Create service with JSON file in which policy/rolepolicies are specified within service
func TestMats_CreateServiceWithPolicyByJsonFile(t *testing.T) {
	sName := "serviceWithPolicy"
	sType := pms.TypeApplication
	jsonFile := "/tmp/service.json"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
		FileName:      jsonFile,
	}

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateServiceWithPolicy",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.CmdTestData{
				Param: param.CREATE_SERVICE_WITH_JSONFILE(jsonFile),
				FileContent: &pms.Service{
					Name: sName,
					Type: sType,
					Policies: []*pms.Policy{
						{
							Name:   "policy1",
							Effect: pms.Grant,
							Permissions: []*pms.Permission{
								{
									Resource: "res1",
									Actions:  []string{"get", "delete", "add"},
								},
							},
							Principals: [][]string{{"group:Administrators"}},
						},
						{
							Name:   "policy2",
							Effect: pms.Grant,
							Permissions: []*pms.Permission{
								{
									Resource: "res2",
									Actions:  []string{"read", "write"},
								},
							},
							Principals: [][]string{{"role:role1"}},
						},
					},
					RolePolicies: []*pms.RolePolicy{
						{
							Name:       "rolePolicy1",
							Effect:     pms.Deny,
							Roles:      []string{"role1"},
							Principals: []string{"user:userA"},
							Resources:  []string{"res1"},
						},
					},
				},
				ExpectedMsg: msg.OUTPUT_SERVICE_CREATED(),
				OutputBody:  &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName,
					Type: sType,
				},
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				tmpService := cmdTD.FileContent.(*pms.Service)
				testutil.GenerateJsonFileWithService(context.FileName, tmpService)
				cmdTD.ExpectedBody = cmdTD.FileContent
			},
		},
		{
			Name:     "TestDeleteServiceWithPolicy",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, context)
}

//Create service with PDL file
func testMats_ServiceWithPDLFile(t *testing.T) {

	sName := "TestMats_PolicyWithPDLFile"
	sType := pms.TypeApplication

	pdlFile := "/tmp/pdl.txt"

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateServiceWithPdlFile",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.CmdTestData{
				Param: param.CREATE_SERVICE_WITH_PDLFILE(sName, sType, pdlFile),
				FileContent: []string{
					"grant user User1 Role1 on res1",
					"grant group Group1 Role2 on res2",
					"grant group Administrators GET,POST,DELETE expr:/service/*",
				},
				ExpectedMsg: msg.OUTPUT_SERVICE_CREATED(),
				OutputBody:  &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName,
					Type: sType},
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, nil)
}

//Create/Get/Delete policy
func TestMats_Policy(t *testing.T) {

	sName := "TestMats_Policy"
	sType := pms.TypeApplication

	pName := "policy1"
	pName1 := "policy2"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateService1",
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
		},
		{
			Name:     "TestGetService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.CmdTestData{
				Param:      param.GET_SERVICE(sName),
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName,
					Type: sType,
				},
			},
		},
		{
			Name:     "TestCreatePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_POLICY(sName, pName, "grant group Administrators list,watch,get c1:default/core/pods/*"),
				ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
				OutputBody:  &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   pName,
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "c1:default/core/pods/*",
							Actions:  []string{"list", "watch", "get"},
						},
					},
					Principals: [][]string{{"group:Administrators"}},
				},
			},
		},
		{
			Name:     "TestCreatePolicy2",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_POLICY(sName, pName1, "deny user userA get,delete c1:default/core/pods/*"),
				ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
				OutputBody:  &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   pName1,
					Effect: pms.Deny,
					Permissions: []*pms.Permission{
						{
							Resource: "c1:default/core/pods/*",
							Actions:  []string{"get", "delete"},
						},
					},
					Principals: [][]string{{"user:userA"}},
				},
			},
		},
		{
			Name:     "TestCreatePolicy3",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_POLICY(sName, "policy3", "grant group Administrators from oracle list,watch,get c1:default/core/nodes/*"),
				ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
				OutputBody:  &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   "policy3",
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "c1:default/core/nodes/*",
							Actions:  []string{"list", "watch", "get"},
						},
					},
					Principals: [][]string{{"idd=oracle:group:Administrators"}},
				},
			},
		},
		{
			Name:     "TestGetPolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.CmdTestData{
				Param:      "TO be init in PreTestFunc with id",
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   pName,
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "c1:default/core/pods/*",
							Actions:  []string{"list", "watch", "get"},
						},
					},
					Principals: [][]string{{"group:Administrators"}},
				},
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[pName]
				if ok {
					cmdTD.Param = param.GET_POLICY(sName, id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestGetPolicyAll",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_QUERY_POLICY,
			Data: &testutil.CmdTestData{
				Param:      param.GET_POLICY_ALL(sName),
				OutputBody: &[]*pms.Policy{},
				ExpectedBody: &[]*pms.Policy{
					{
						Name:   pName,
						Effect: pms.Grant,
						Permissions: []*pms.Permission{
							{
								Resource: "c1:default/core/pods/*",
								Actions:  []string{"list", "watch", "get"},
							},
						},
						Principals: [][]string{{"group:Administrators"}},
					},
					{
						Name:   pName1,
						Effect: pms.Deny,
						Permissions: []*pms.Permission{
							{
								Resource: "c1:default/core/pods/*",
								Actions:  []string{"get", "delete"},
							},
						},
						Principals: [][]string{{"user:userA"}},
					},
					{
						Name:   "policy3",
						Effect: pms.Grant,
						Permissions: []*pms.Permission{
							{
								Resource: "c1:default/core/nodes/*",
								Actions:  []string{"list", "watch", "get"},
							},
						},
						Principals: [][]string{{"idd=oracle:group:Administrators"}},
					},
				},
			},
		},
		{
			Name:     "TestDeletePolicy2",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.CmdTestData{
				Param:       "to be init in preTestFunc",
				ExpectedMsg: "to be init in preTestFunc",
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[pName1]
				if ok {
					cmdTD.Param = param.DELETE_POLICY(sName, id)
					cmdTD.ExpectedMsg = msg.OUTPUT_POLICY_DELETED(id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestGetPolicy2-NotExisting",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.CmdTestData{
				Param:       "TO be init in PreTestFunc with id",
				ExpectedMsg: "TO be init in PreTestFunc with id",
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[pName1]
				if ok {
					cmdTD.Param = param.GET_POLICY(sName, id)
					cmdTD.ExpectedMsg = msg.OUTPUT_POLICY_NOTFOUND(sName, id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, context)
}

//Create/Get/Delete role-policy
func TestMats_RolePolicy(t *testing.T) {

	sName := "TestMats_RolePolicy"
	sType := pms.TypeApplication

	rpName := "role-policy1"
	rpName1 := "role-policy2"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateService1",
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
		},
		{
			Name:     "TestGetService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_SERVICE,
			Data: &testutil.CmdTestData{
				Param:      param.GET_SERVICE(sName),
				OutputBody: &pms.Service{},
				ExpectedBody: &pms.Service{
					Name: sName,
					Type: sType,
				},
			},
		},
		{
			Name:     "TestCreateRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_ROLEPOLICY(sName, rpName, "grant user user1 role1 on res1"),
				ExpectedMsg: msg.OUTPUT_ROLEPOLICY_CREATED(),
				OutputBody:  &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       rpName,
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
			},
		},
		{
			Name:     "TestCreateRolePolicy2",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_ROLEPOLICY(sName, rpName1, "grant user User1,role Role2 Role1,Role2 on res2"),
				ExpectedMsg: msg.OUTPUT_ROLEPOLICY_CREATED(),
				OutputBody:  &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       rpName1,
					Effect:     pms.Grant,
					Roles:      []string{"Role1", "Role2"},
					Principals: []string{"user:User1", "role:Role2"},
					Resources:  []string{"res2"},
				},
			},
		},
		{
			Name:     "TestGetRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:      "TO be init in PreTestFunc with id",
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       rpName,
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"user:user1"},
					Resources:  []string{"res1"},
				},
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[rpName]
				if ok {
					cmdTD.Param = param.GET_ROLEPOLICY(sName, id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestGetRolePolicyAll",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_QUERY_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:      param.GET_ROLEPOLICY_ALL(sName),
				OutputBody: &[]*pms.RolePolicy{},
				ExpectedBody: &[]*pms.RolePolicy{
					{
						Name:       rpName,
						Effect:     pms.Grant,
						Roles:      []string{"role1"},
						Principals: []string{"user:user1"},
						Resources:  []string{"res1"},
					},
					{
						Name:       rpName1,
						Effect:     pms.Grant,
						Roles:      []string{"Role1", "Role2"},
						Principals: []string{"user:User1", "role:Role2"},
						Resources:  []string{"res2"},
					},
				},
			},
		},
		{
			Name:     "TestDeleteRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:       "to be init in preTestFun",
				ExpectedMsg: "to be init in preTestFun",
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[rpName]
				if ok {
					cmdTD.Param = param.DELETE_ROLEPOLICY(sName, id)
					cmdTD.ExpectedMsg = msg.OUTPUT_ROLEPOLICY_DELETED(id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestGetRolePolicy1-NotExisting",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:       "TO be init in PreTestFunc with id",
				ExpectedMsg: "TO be init in PreTestFunc with id",
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[rpName]
				if ok {
					cmdTD.Param = param.GET_ROLEPOLICY(sName, id)
					cmdTD.ExpectedMsg = msg.OUTPUT_ROLEPOLICY_NOTFOUND(sName, id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, context)
}

//Create/Get/Delete policy/rolepolicy with principle containing single entity
func TestMats_PolicyWithSingleEntity(t *testing.T) {

	sName := "TestMats_PolicyWithEntity"
	sType := pms.TypeApplication

	pName := "policy1"
	rpName := "role-policy1"

	context := &testutil.TestContext{
		NameIDMap:     make(map[string]string),
		NameObjectMap: make(map[string]interface{}),
	}

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateService1",
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
		},
		{
			Name:     "TestCreatePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_POLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_POLICY(sName, pName, "grant entity spiffe://staging.acme.com/payments/mysql list,watch,get res1"),
				ExpectedMsg: msg.OUTPUT_POLICY_CREATED(),
				OutputBody:  &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   pName,
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"list", "watch", "get"},
						},
					},
					Principals: [][]string{{"entity:spiffe://staging.acme.com/payments/mysql"}},
				},
			},
		},
		{
			Name:     "TestGetPolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_POLICY,
			Data: &testutil.CmdTestData{
				Param:      "TO be init in PreTestFunc with id",
				OutputBody: &pms.Policy{},
				ExpectedBody: &pms.Policy{
					Name:   pName,
					Effect: pms.Grant,
					Permissions: []*pms.Permission{
						{
							Resource: "res1",
							Actions:  []string{"list", "watch", "get"},
						},
					},
					Principals: [][]string{{"entity:spiffe://staging.acme.com/payments/mysql"}},
				},
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[pName]
				if ok {
					cmdTD.Param = param.GET_POLICY(sName, id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestDeletePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_POLICY,
			Data: &testutil.CmdTestData{
				Param:       "to be init in preTestFunc",
				ExpectedMsg: "to be init in preTestFunc",
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[pName]
				if ok {
					cmdTD.Param = param.DELETE_POLICY(sName, id)
					cmdTD.ExpectedMsg = msg.OUTPUT_POLICY_DELETED(id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestCreateRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:       param.CREATE_ROLEPOLICY(sName, rpName, "grant entity spiffe://staging.acme.com/payments/mysql role1 on res1"),
				ExpectedMsg: msg.OUTPUT_ROLEPOLICY_CREATED(),
				OutputBody:  &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       rpName,
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"entity:spiffe://staging.acme.com/payments/mysql"},
					Resources:  []string{"res1"},
				},
			},
		},
		{
			Name:     "TestGetRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_GET_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:      "TO be init in PreTestFunc with id",
				OutputBody: &pms.RolePolicy{},
				ExpectedBody: &pms.RolePolicy{
					Name:       rpName,
					Effect:     pms.Grant,
					Roles:      []string{"role1"},
					Principals: []string{"entity:spiffe://staging.acme.com/payments/mysql"},
					Resources:  []string{"res1"},
				},
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[rpName]
				if ok {
					cmdTD.Param = param.GET_ROLEPOLICY(sName, id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestDeleteRolePolicy1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_ROLEPOLICY,
			Data: &testutil.CmdTestData{
				Param:       "to be init in preTestFun",
				ExpectedMsg: "to be init in preTestFun",
			},
			PreTestFunc: func(data interface{}, context *testutil.TestContext) {
				cmdTD := data.(*testutil.CmdTestData)
				id, ok := context.NameIDMap[rpName]
				if ok {
					cmdTD.Param = param.DELETE_ROLEPOLICY(sName, id)
					cmdTD.ExpectedMsg = msg.OUTPUT_ROLEPOLICY_DELETED(id)
					testutil.TestLog.Log(id)
				}
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, context)
}

//Create service with special char
func TestLrg_Service_SpecialCharInName_bug66_bug52(t *testing.T) {

	sName := "k8s-*=+/_."
	sType := "k8s-type"

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateService1",
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
		},
		{
			Name:     "TestCreateService1_repeat",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:        param.CREATE_SERVICE(sName, sType),
				ExpectedMsg:  "Already Exist",
				OutputBody:   &pms.Service{},
				ExpectedBody: &pms.Service{},
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, nil)

}

//Negative test case to create service

func TestLrg_Service_Neg_CreateTwice_bug66(t *testing.T) {
	sName := "k8s"
	sType := pms.TypeApplication

	data := &[]testutil.TestCase{
		{
			Name:     "TestCreateService1",
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
		},
		{
			Name:     "TestCreateService1_repeat",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_CREATE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:        param.CREATE_SERVICE(sName, sType),
				ExpectedMsg:  "Already Exist",
				OutputBody:   &pms.Service{},
				ExpectedBody: &pms.Service{},
			},
		},
		{
			Name:     "TestDeleteService1",
			Enabled:  true,
			Executer: testutil.NewCmdTest,
			Method:   testutil.METHOD_DELETE_SERVICE,
			Data: &testutil.CmdTestData{
				Param:       param.DELETE_SERVICE(sName),
				ExpectedMsg: msg.OUTPUT_SERVICE_DELETED(sName),
			},
		},
	}

	testutil.RunTestCases(t, data, nil)

}
