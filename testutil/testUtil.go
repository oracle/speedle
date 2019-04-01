//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package testutil

import (
	"reflect"
	"sort"
	"strings"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/svcs/pmsgrpc/pb"
)

//Test Methods for Speedle
const (
	METHOD_CREATE_SERVICE          = "createService"
	METHOD_GET_SERVICE             = "getService" //not used in grpc
	METHOD_QUERY_SERVICE           = "queryService"
	METHOD_DELETE_SERVICE          = "deleteService"
	METHOD_CREATE_POLICY           = "createPolicy"
	METHOD_GET_POLICY              = "getPolicy" //not used in grpc
	METHOD_QUERY_POLICY            = "queryPolicy"
	METHOD_DELETE_POLICY           = "deletePolicy"
	METHOD_CREATE_ROLEPOLICY       = "createRolePolicy"
	METHOD_GET_ROLEPOLICY          = "getRolePolicy" //not used in grpc
	METHOD_QUERY_ROLEPOLICY        = "queryRolePolicy"
	METHOD_DELETE_ROLEPOLICY       = "deleteRolePolicy"
	METHOD_IS_ALLOWED              = "isAllowed"
	METHOD_GET_GRANTED_PERMISSIONS = "getPerssions"
	METHOD_GET_GRANTED_ROLES       = "getRoles"
	METHOD_SLEEP                   = "sleep"
)

const (
	ERROR_SPEEDLE_NOT_SUPPORTED = "speedle doesn't support"
)

/**
  Function called after create/get policy test
  Modify the ID in outputBody and expected body
**/
func PostCreateGetPolicyTest(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postCreateGetPolicyTest_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postCreateGetPolicyTest_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
	} else if grpcTD, ok := data.(*GRpcTestData); ok {
		postCreateGetPolicyTest_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

func postCreateGetPolicyTest_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {

	if outputBody != nil {
		//check if outputPolicy is pms.Policy (used in kauctl command and REST test)
		outputPolicy, ok := outputBody.(*pms.Policy)
		if ok {
			context.NameIDMap[outputPolicy.Name] = outputPolicy.ID
			context.NameObjectMap[outputPolicy.Name] = outputPolicy
		} else {
			TestLog.Log("Fail to convert outputBody to policy object array")
			return
		}

	} else {
		return
	}

	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*pms.Policy)
		if ok {
			expectedPolicy.ID = context.NameIDMap[expectedPolicy.Name]
		}
	}
}

func postCreateGetPolicyTest_grpc(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	if outputBody != nil {
		outputPolicy, ok := outputBody.(*pb.Policy)
		if ok {
			context.NameIDMap[outputPolicy.Name] = outputPolicy.Id
		} else {
			TestLog.Log("Fail to convert outputBody to policy object array")
			return
		}

	} else {
		return
	}

	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*pb.Policy)
		if ok {
			expectedPolicy.Id = context.NameIDMap[expectedPolicy.Name]
		}
	}
}

/**
  Function called after List policy test
  Modify the ID in outputBody and expectedBody
**/

func PostListPolicyTest(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postListPolicyTest_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postListPolicyTest_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
	} else if grpcTD, ok := data.(*GRpcTestData); ok {
		postListPolicyTest_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

func postListPolicyTest_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	adjustMap := make(map[string]*pms.Policy)

	if outputBody != nil {
		outputPolicy, ok := outputBody.(*[]*pms.Policy)
		if !ok {
			TestLog.Log("Fail to convert outputBody to policy object array")
			return
		}

		for i := 0; i < len(*outputPolicy); i = i + 1 {
			context.NameIDMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i].ID
			adjustMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i]
		}
	} else {
		return
	}

	newOutputBody := []*pms.Policy{}
	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*[]*pms.Policy)
		if !ok {
			TestLog.Log("Fail to convert expectedBody to policy object")
			return
		}

		for i := 0; i < len(*expectedPolicy); i = i + 1 {
			id, ok := context.NameIDMap[(*expectedPolicy)[i].Name]
			if ok {
				(*expectedPolicy)[i].ID = id
			}
			exp, ok := adjustMap[(*expectedPolicy)[i].Name]
			if ok {
				newOutputBody = append(newOutputBody, exp)
				delete(adjustMap, (*expectedPolicy)[i].Name)
			}
		}
	}
	for _, value := range adjustMap {
		newOutputBody = append(newOutputBody, value)

	}

	outputPolicy, ok := outputBody.(*[]*pms.Policy)
	if !ok {
		TestLog.Log("Fail to convert outputBody to policy object array")
		return
	}
	*outputPolicy = newOutputBody
}

func postListPolicyTest_grpc(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	adjustMap := make(map[string]*pb.Policy)

	if outputBody != nil {
		outputPolicy, ok := outputBody.(*[]*pb.Policy)
		if !ok {
			TestLog.Log("Fail to convert outputBody to pb.policy object array")
			return
		}

		for i := 0; i < len(*outputPolicy); i = i + 1 {
			context.NameIDMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i].Id
			adjustMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i]
		}
	} else {
		return
	}

	newOutputBody := []*pb.Policy{}
	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*[]*pb.Policy)
		if !ok {
			TestLog.Log("Fail to convert expectedBody to pb.policy object")
			return
		}

		for i := 0; i < len(*expectedPolicy); i = i + 1 {
			id, ok := context.NameIDMap[(*expectedPolicy)[i].Name]
			if ok {
				(*expectedPolicy)[i].Id = id
			}
			exp, ok := adjustMap[(*expectedPolicy)[i].Name]
			if ok {
				newOutputBody = append(newOutputBody, exp)
				delete(adjustMap, (*expectedPolicy)[i].Name)
			}
		}
	}
	for _, value := range adjustMap {
		newOutputBody = append(newOutputBody, value)

	}

	outputPolicy, ok := outputBody.(*[]*pb.Policy)
	if !ok {
		TestLog.Log("Fail to convert outputBody to pb.policy object array")
		return
	}
	*outputPolicy = newOutputBody
}

/**
  Function called after create/get role policy test
  Modify the ID in outputBody and expected body
**/
func PostCreateGetRolePolicyTest(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postCreateGetRolePolicyTest_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postCreateGetRolePolicyTest_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
	} else if grpcTD, ok := data.(*GRpcTestData); ok {
		postCreateGetRolePolicyTest_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

func postCreateGetRolePolicyTest_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	var outputPolicy *pms.RolePolicy
	if outputBody != nil {
		var ok bool
		outputPolicy, ok = outputBody.(*pms.RolePolicy)
		if !ok {
			return
		}
		context.NameIDMap[outputPolicy.Name] = outputPolicy.ID
		context.NameObjectMap[outputPolicy.Name] = outputPolicy
	} else {
		return
	}

	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*pms.RolePolicy)
		if !ok {
			return
		}
		expectedPolicy.ID = outputPolicy.ID
	}
}

func postCreateGetRolePolicyTest_grpc(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	var outputPolicy *pb.RolePolicy
	if outputBody != nil {
		ok := false
		outputPolicy, ok = outputBody.(*pb.RolePolicy)
		if !ok {
			return
		}
		context.NameIDMap[outputPolicy.Name] = outputPolicy.Id
	} else {
		return
	}

	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*pb.RolePolicy)
		if !ok {
			return
		}
		expectedPolicy.Id = context.NameIDMap[expectedPolicy.Name]
	}
}

/**
  Function called after List role  policy test
  Modify the ID in outputBody and expected body
**/
func PostListRolePolicyTest(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postListRolePolicyTest_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postListRolePolicyTest_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
	} else if grpcTD, ok := data.(*GRpcTestData); ok {
		postListRolePolicyTest_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

func postListRolePolicyTest_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	adjustMap := make(map[string]*pms.RolePolicy)
	if outputBody != nil {
		outputPolicy, ok := outputBody.(*[]*pms.RolePolicy)
		if !ok {
			return
		}
		for i := 0; i < len(*outputPolicy); i = i + 1 {
			context.NameIDMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i].ID
			adjustMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i]
		}
	} else {
		return
	}

	newOutputBody := []*pms.RolePolicy{}
	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*[]*pms.RolePolicy)
		if !ok {
			return
		}
		for i := 0; i < len(*expectedPolicy); i = i + 1 {
			id, ok := context.NameIDMap[(*expectedPolicy)[i].Name]
			if ok {
				(*expectedPolicy)[i].ID = id
			}

			exp, ok := adjustMap[(*expectedPolicy)[i].Name]
			if ok {
				newOutputBody = append(newOutputBody, exp)
				delete(adjustMap, (*expectedPolicy)[i].Name)
			}
		}
	}
	for _, value := range adjustMap {
		newOutputBody = append(newOutputBody, value)
	}

	outputPolicy, ok := outputBody.(*[]*pms.RolePolicy)
	if !ok {
		return
	}
	*outputPolicy = newOutputBody
}

func postListRolePolicyTest_grpc(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	adjustMap := make(map[string]*pb.RolePolicy)
	if outputBody != nil {
		outputPolicy, ok := outputBody.(*[]*pb.RolePolicy)
		if !ok {
			return
		}
		for i := 0; i < len(*outputPolicy); i = i + 1 {
			context.NameIDMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i].Id
			adjustMap[(*outputPolicy)[i].Name] = (*outputPolicy)[i]
		}
	} else {
		return
	}

	newOutputBody := []*pb.RolePolicy{}
	if expectedBody != nil {
		expectedPolicy, ok := expectedBody.(*[]*pb.RolePolicy)
		if !ok {
			return
		}
		for i := 0; i < len(*expectedPolicy); i = i + 1 {
			id, ok := context.NameIDMap[(*expectedPolicy)[i].Name]
			if ok {
				(*expectedPolicy)[i].Id = id
			}

			exp, ok := adjustMap[(*expectedPolicy)[i].Name]
			if ok {
				newOutputBody = append(newOutputBody, exp)
				delete(adjustMap, (*expectedPolicy)[i].Name)
			}
		}
	}
	for _, value := range adjustMap {
		newOutputBody = append(newOutputBody, value)
	}

	outputPolicy, ok := outputBody.(*[]*pb.RolePolicy)
	if !ok {
		return
	}
	*outputPolicy = newOutputBody
}

/**
  Function called after Create or Get service test
  Modify the ID in outputBody and expected body
**/
func PostCreateGetServiceTest(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postCreateGetServiceTest_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postCreateGetServiceTest_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
	} else if grpcTD, ok := data.(*GRpcTestData); ok {
		postCreateGetServiceTest_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

func postCreateGetServiceTest_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	if outputBody != nil {
		outputService, ok := outputBody.(*pms.Service)
		if !ok {
			TestLog.Log("Fail to convert outputBody to pms.Service")
			return
		}
		expectService, _ := expectedBody.(*pms.Service)

		//The policies/rolepolicies is nil in expect/output sometimes
		if outputService.Policies == nil {
			TestLog.Log("outputService.Policies is nil. Set it to Empty")
			outputService.Policies = []*pms.Policy{}
		}
		if outputService.RolePolicies == nil {
			TestLog.Log("outputService.RolePolicies is nil. Set it to Empty")
			outputService.RolePolicies = []*pms.RolePolicy{}
		}

		if expectService.Policies == nil {
			TestLog.Log("expectService.Policies is nil. Set it to Empty")
			expectService.Policies = []*pms.Policy{}
		}
		if expectService.RolePolicies == nil {
			TestLog.Log("expectService.RolePolicies is nil. Set it to Empty")
			expectService.RolePolicies = []*pms.RolePolicy{}
		}

		postListPolicyTest_rest(&outputService.Policies, &expectService.Policies, context)
		postListRolePolicyTest_rest(&outputService.RolePolicies, &expectService.RolePolicies, context)
	} else {
		return
	}
}

func postCreateGetServiceTest_grpc(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	if outputBody != nil {
		outputService, ok := outputBody.(*pb.Service)
		if !ok {
			TestLog.Log("Fail to convert outputBody to pb.Service")
			return
		}
		expectService, _ := expectedBody.(*pb.Service)
		postListPolicyTest_rest(&outputService.Policies, &expectService.Policies, context)
		postListRolePolicyTest_rest(&outputService.RolePolicies, &expectService.RolePolicies, context)
	} else {
		return
	}
}

/**
  Function called after listing service test
  Modify the ID in outputBody and expected body
**/
func PostListServiceTest(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postListServiceTest_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postListServiceTest_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
	} else if grpcTD, ok := data.(*GRpcTestData); ok {
		postListServiceTest_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

func postListServiceTest_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {

	if outputBody != nil {
		adjustMap := make(map[string]*pms.Service)

		outputService, ok := outputBody.(*[]*pms.Service)
		if !ok {
			TestLog.Log("Fail to convert outputBody to pms.Service")
			return
		}
		for i := 0; i < len(*outputService); i = i + 1 {
			adjustMap[(*outputService)[i].Name] = (*outputService)[i]
		}

		newOutputBody := []*pms.Service{}
		if expectedBody != nil {
			expectedService, ok := expectedBody.(*[]*pms.Service)
			if !ok {
				return
			}
			for i := 0; i < len(*expectedService); i = i + 1 {
				exp, ok := adjustMap[(*expectedService)[i].Name]
				if ok {
					newOutputBody = append(newOutputBody, exp)
					delete(adjustMap, (*expectedService)[i].Name)
				}
			}
		}
		for _, value := range adjustMap {
			newOutputBody = append(newOutputBody, value)
		}

		*outputService = newOutputBody

	} else {
		return
	}
}

func postListServiceTest_grpc(outputBody interface{}, expectedBody interface{}, context *TestContext) {

	if outputBody != nil {
		adjustMap := make(map[string]*pb.Service)

		outputService, ok := outputBody.(*[]*pb.Service)
		if !ok {
			TestLog.Log("Fail to convert outputBody to pb.Service")
			return
		}
		for i := 0; i < len(*outputService); i = i + 1 {
			adjustMap[(*outputService)[i].Name] = (*outputService)[i]
		}

		newOutputBody := []*pb.Service{}
		if expectedBody != nil {
			expectedService, ok := expectedBody.(*[]*pb.Service)
			if !ok {
				return
			}
			for i := 0; i < len(*expectedService); i = i + 1 {
				exp, ok := adjustMap[(*expectedService)[i].Name]
				if ok {
					newOutputBody = append(newOutputBody, exp)
					delete(adjustMap, (*expectedService)[i].Name)
				}
			}
		}
		for _, value := range adjustMap {
			newOutputBody = append(newOutputBody, value)
		}

		*outputService = newOutputBody

	} else {
		return
	}
}

//Sort the output by role name string
func PostGetAllGrantedRoles(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postGetAllGrantedRoles_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postGetAllGrantedRoles_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
		//} else if grpcTD, ok := data.(*GRpcTestData); ok {
		//postGetAllGrantedRoles_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, are supported by now!")
	}
}

func postGetAllGrantedRoles_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	if outputBody != nil {
		outputRoles, _ := outputBody.(*[]string)
		sort.Strings(*outputRoles)
	}

}

//Sort the output by resource name string
func PostGetAllGrantedPermissions(data interface{}, context *TestContext) {
	if restTD, ok := data.(*RestTestData); ok {
		postGetAllGrantedPermissions_rest(restTD.OutputBody, restTD.ExpectedBody, context)
	} else if cmdTD, ok := data.(*CmdTestData); ok {
		postGetAllGrantedPermissions_rest(cmdTD.OutputBody, cmdTD.ExpectedBody, context)
		//} else if grpcTD, ok := data.(*GRpcTestData); ok {
		//postGetAllGrantedPermissions_grpc(grpcTD.OutputBody, grpcTD.ExpectedBody, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, are supported by now!")
	}
}

func postGetAllGrantedPermissions_rest(outputBody interface{}, expectedBody interface{}, context *TestContext) {
	if outputBody != nil {
		adjustMap := make(map[string]pms.Permission)
		outputPerms, ok := outputBody.(*[]pms.Permission)
		if !ok {
			return
		}
		for i := 0; i < len(*outputPerms); i = i + 1 {
			adjustMap[(*outputPerms)[i].Resource] = (*outputPerms)[i]
		}

		res := make([]string, len(*outputPerms))
		j := 0
		for k := range adjustMap {
			res[j] = k
			j++
		}
		sort.Strings(res)

		newOutputBody := []pms.Permission{}
		for _, v := range adjustMap {
			newOutputBody = append(newOutputBody, v)
		}

		*outputPerms = newOutputBody

	} else {
		return
	}

}

//Function called before Get/Delete policy/rolepolicy
//Update the policy/rolepolicy name with ID in the request
// since we only support get/delete with ID
func PreGetDeletePolicyTest(data interface{}, context *TestContext) {
	if _, ok := data.(*RestTestData); ok {
		preGetDeletePolicy_rest(data, context)
		//} else if cmdTD, ok := data.(*CmdTestData); ok {
		//	preGetDeletePolicy_cmd(data, context)
	} else if _, ok := data.(*GRpcTestData); ok {
		preGetDeletePolicy_grpc(data, context)
	} else {
		TestLog.Fatalf("Only RestTestData, CmdTestData, GRpcTestData are supported by now!")
	}
}

//update policy/rolepolicy's name with ID in URI
//the uri should be like http://127.0.0.1:6733/policy-mgmt/v1/service/srv1/policy/policy1
//update policy1 with its's ID
func preGetDeletePolicy_rest(data interface{}, context *TestContext) {
	restTD, _ := data.(*RestTestData)
	pos := strings.LastIndex(restTD.URI, "/")
	policyName := restTD.URI[pos+1 : len(restTD.URI)]
	id, ok := context.NameIDMap[policyName]
	if ok {
		restTD.URI = restTD.URI[0:pos+1] + id
	}
}

//Set policyID in grpc request. Name is specified in request.PolicyID or request.RolePolicyID
//And we should update it with ID of policy/rolePolicy
func preGetDeletePolicy_grpc(data interface{}, context *TestContext) {
	grpcTD := data.(*GRpcTestData)
	request, ok := grpcTD.InputBody.(*pb.PolicyQueryRequest)
	if ok {
		id, ok1 := context.NameIDMap[request.PolicyID]
		if ok1 {
			request.PolicyID = id
		}
	} else {
		TestLog.Log("preGetDeletePolicy_grpc----2")
		request, ok := grpcTD.InputBody.(*pb.RolePolicyQueryRequest)
		if ok {
			TestLog.Log("preGetDeletePolicy_grpc----2-RolePolicy")
			id, ok1 := context.NameIDMap[request.RolePolicyID]
			if ok1 {
				request.RolePolicyID = id
			}
		} else {
			TestLog.Fatalf("The request is not PolicyQueryRequest/RolePolicyQueryRequest")
		}
	}
}

func removeMetaData(data reflect.Value) {
	v := reflect.Indirect(data)
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			removeMetaData(v.Index(i))
		}
	case reflect.Struct:
		t := v.Type()
		if "Service" == t.Name() {
			s := v.Interface().(pms.Service)
			removeMetaData(reflect.ValueOf(s.Policies))
			removeMetaData(reflect.ValueOf(s.RolePolicies))
		}
		f := v.FieldByName("Metadata")
		if f.IsValid() {
			//TODO: Meta data should be not be removed if we need to verify meta data in future.
			//Method to remove createby and createtime entries:
			//	meta := f.Interface().(map[string]string)
			//	delete(meta, "createby")
			//	delete(meta, "createtime")
			f.Set(reflect.Zero(f.Type()))
		}
	}
	return
}

var SpecialNames []string

func InitSpecialNames() {
	if len(SpecialNames) == 0 {
		SpecialNames = append(SpecialNames, "1")
		SpecialNames = append(SpecialNames, "1_")
		SpecialNames = append(SpecialNames, "1*")
		SpecialNames = append(SpecialNames, "_Name1_")
		SpecialNames = append(SpecialNames, "#Name1##Name1#")
		SpecialNames = append(SpecialNames, "@Name1@@Name1@@")
		SpecialNames = append(SpecialNames, "$Name1$$Name1$$")
		SpecialNames = append(SpecialNames, "%Name1%%Name1%%")
		SpecialNames = append(SpecialNames, "&Name1&&Name1&&")
		SpecialNames = append(SpecialNames, "'Name1''Name1''")
		SpecialNames = append(SpecialNames, "(Group1))Name1))")
		SpecialNames = append(SpecialNames, "*Name1*Name1**")
		SpecialNames = append(SpecialNames, "\"Name1\"\"Name1\"\"") //must use \" instead of "
		SpecialNames = append(SpecialNames, "+Name1++Name1++")
		SpecialNames = append(SpecialNames, "-Name1--Name1--")
		SpecialNames = append(SpecialNames, ".Name3..Name1..")
		SpecialNames = append(SpecialNames, ":Name4::Name1::")
		SpecialNames = append(SpecialNames, ";Name5;;Name1;;")
		SpecialNames = append(SpecialNames, "=Name6==Name1==")
		SpecialNames = append(SpecialNames, "<Name1>>Name1>>")
		SpecialNames = append(SpecialNames, "?Name1??Name1??")
		SpecialNames = append(SpecialNames, "[Name2]]Name1]]")
		SpecialNames = append(SpecialNames, "\\\\Name3\\\\\\\\Name1\\\\") //must use \\ instead of \
		SpecialNames = append(SpecialNames, "^Name4^^Name1^^")
		SpecialNames = append(SpecialNames, "{Name6}}Name1}}")
		SpecialNames = append(SpecialNames, "\\|Name7\\|\\|Name1\\|\\|") //must use \| instead of |
		SpecialNames = append(SpecialNames, "~Name1~~Name1~~")
		SpecialNames = append(SpecialNames, "!Name1!!Name1!!")
		SpecialNames = append(SpecialNames, "/Name1//Name1//")
		SpecialNames = append(SpecialNames, "\\`Name1\\`\\`Name1\\`\\`") //must use \` instead of `

	}
}

//CompareStringArray_NoOrder check two string arrays containing same elements ignoring the order
func CompareStringArray_NoOrder(strArray1 []string, strArray2 []string) bool {
	TestLog.Logf("CompareStringArray_NoOrder: Array1=%v", strArray1)
	TestLog.Logf("CompareStringArray_NoOrder: Array2=%v", strArray2)
	//init map[string]int, key is the string in array, the value is the occurrence of string
	//Then no repeat string in map any more
	strMap1 := make(map[string]int)
	strMap2 := make(map[string]int)

	for _, str1 := range strArray1 {
		if v, contained := strMap1[str1]; contained {
			strMap1[str1] = v + 1
		} else {
			strMap1[str1] = 1
		}
	}

	for _, str2 := range strArray2 {
		if v, contained := strMap2[str2]; contained {
			strMap2[str2] = v + 1
		} else {
			strMap2[str2] = 1
		}
	}
	// TestLog.Logf("CompareStringArray_NoOrder: map1=%v, map2=%v", strMap1, strMap2)

	//Compare the length of two maps
	if len(strMap1) != len(strMap2) {
		TestLog.Errorf("The map length is not equal. Please check!")
		return false
	}
	//compare the key*value in two maps
	for k1, v1 := range strMap1 {
		if v2, contained := strMap2[k1]; contained {
			if v1 != v2 {
				TestLog.Errorf("The count of string %s is not same. one is %d another is %d", k1, v1, v2)
				return false
			}
		} else {
			TestLog.Errorf("The string:\"%s\" in first array is not in the second array", k1)
			return false
		}
	}

	return true
}
