//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	adsPB "github.com/oracle/speedle/pkg/svcs/adsgrpc/pb"
	pmsPB "github.com/oracle/speedle/pkg/svcs/pmsgrpc/pb"

	"google.golang.org/grpc"
)

//-------------------GRpcClient definition------------------------

type GRpcClient struct {
	pmsConn   *grpc.ClientConn          //pmsConnect for GRpc Client for PMS
	pmsClient pmsPB.PolicyManagerClient //policy mananger client
	adsConn   *grpc.ClientConn          //pmsConnect for GRpc Client for ADS
	adsClient adsPB.EvaluatorClient     //ads (evaluator) client
}

func NewGRpcClient() *GRpcClient {
	gc := &GRpcClient{}
	return gc
}

func (gc *GRpcClient) SetupConnection() error {
	if gc.pmsConn == nil {
		tmpConn, err := grpc.Dial("localhost:50001", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not pmsConnect: %v", err)
			return err
		}
		gc.pmsClient = pmsPB.NewPolicyManagerClient(tmpConn)
		gc.pmsConn = tmpConn
	}

	if gc.adsConn == nil {
		tmpConn, err := grpc.Dial("localhost:50002", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not adsConnect: %v", err)
			return err
		}
		gc.adsClient = adsPB.NewEvaluatorClient(tmpConn)
		gc.adsConn = tmpConn
	}
	return nil
}

func (gc *GRpcClient) CloseConnection() {
	if gc.pmsConn != nil {
		gc.pmsConn.Close()
		gc.pmsConn = nil
	}
}

//-------------------GRpcTest definition------------------------
//test data for spxctl command
type GRpcTestData struct {
	InputBody    interface{} //Optional. Request body for method
	OutputMsg    string      //Optional. Actual output Message
	OutputBody   interface{} //Optional. Actual output body for json object
	ExpectedMsg  string      //Optional. expected message for the command
	ExpectedBody interface{} //Optional. expected body about json object string
}

type GRpcTest struct {
	Client *GRpcClient
}

func NewGRpcTestExecuter() TestExecuter {
	return &GRpcTest{
		Client: NewGRpcClient(),
	}
}

//Prepare for Test Execution. Set the default func in testcase
func (test *GRpcTest) PreExecute(testcase *TestCase, tc *TestContext) error {
	test.Client.SetupConnection()
	testcase.VerifyTestFunc = VerifyGRpcTestByDefault

	switch testcase.Method {
	case METHOD_CREATE_SERVICE:
		break
	case METHOD_QUERY_SERVICE:
		testcase.SetPostTestFunc(PostListServiceTest)
		break
	case METHOD_DELETE_SERVICE:
		break
	case METHOD_CREATE_POLICY:
		testcase.SetPostTestFunc(PostCreateGetPolicyTest)
		break
	case METHOD_QUERY_POLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		testcase.SetPostTestFunc(PostListPolicyTest)
		break
	case METHOD_DELETE_POLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		break
	case METHOD_CREATE_ROLEPOLICY:
		testcase.SetPostTestFunc(PostCreateGetRolePolicyTest)
		break
	case METHOD_QUERY_ROLEPOLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		testcase.SetPostTestFunc(PostListRolePolicyTest)
		break
	case METHOD_DELETE_ROLEPOLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		break
	case METHOD_IS_ALLOWED:
		//use default verification method: VerifyGRpcTestByDefault
		break
	case METHOD_GET_GRANTED_ROLES:
		testcase.VerifyTestFunc = VerifyGRpcAllGrantedRoleResult
		break
	case METHOD_GET_GRANTED_PERMISSIONS:
		testcase.VerifyTestFunc = VerifyGRpcAllGrantedPermissions
		break
	default:
		return errors.New(ERROR_SPEEDLE_NOT_SUPPORTED)
	}

	return nil
}

//Execute current test with test data and context
func (test *GRpcTest) Execute(testcase *TestCase, tc *TestContext) error {

	grpcTD := testcase.Data.(*GRpcTestData)

	var resp interface{}
	var err error

	switch testcase.Method {
	case METHOD_CREATE_SERVICE:
		resp, err = test.Client.pmsClient.CreateService(context.Background(), grpcTD.InputBody.(*(pmsPB.ServiceRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_QUERY_SERVICE:
		resp, err = test.Client.pmsClient.QueryServices(context.Background(), grpcTD.InputBody.(*(pmsPB.ServiceQueryRequest)))
		serviceResp, ok := resp.(*pmsPB.ServiceQueryResponse)
		if err == nil && ok {
			grpcTD.OutputBody = &serviceResp.Services
		}
		break
	case METHOD_DELETE_SERVICE:
		resp, err = test.Client.pmsClient.DeleteServices(context.Background(), grpcTD.InputBody.(*(pmsPB.ServiceQueryRequest)))
		grpcTD.OutputBody = resp
		break
	case METHOD_CREATE_POLICY:
		resp, err = test.Client.pmsClient.CreatePolicy(context.Background(), grpcTD.InputBody.(*(pmsPB.PolicyRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_QUERY_POLICY:
		resp, err = test.Client.pmsClient.QueryPolicies(context.Background(), grpcTD.InputBody.(*(pmsPB.PolicyQueryRequest)))
		policyRsp, ok := resp.(*pmsPB.PolicyQueryResponse)
		if err == nil && ok {
			grpcTD.OutputBody = &policyRsp.Policies
		}
		break
	case METHOD_DELETE_POLICY:
		resp, err = test.Client.pmsClient.DeletePolicies(context.Background(), grpcTD.InputBody.(*(pmsPB.PolicyQueryRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_CREATE_ROLEPOLICY:
		resp, err = test.Client.pmsClient.CreateRolePolicy(context.Background(), grpcTD.InputBody.(*(pmsPB.RolePolicyRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_QUERY_ROLEPOLICY:
		resp, err = test.Client.pmsClient.QueryRolePolicies(context.Background(), grpcTD.InputBody.(*(pmsPB.RolePolicyQueryRequest)))
		rolePolicyRsp, ok := resp.(*pmsPB.RolePolicyQueryResponse)
		if err == nil && ok {
			grpcTD.OutputBody = &rolePolicyRsp.RolePolicies
		}
		break
	case METHOD_DELETE_ROLEPOLICY:
		resp, err = test.Client.pmsClient.DeleteRolePolicies(context.Background(), grpcTD.InputBody.(*(pmsPB.RolePolicyQueryRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_IS_ALLOWED:
		resp, err = test.Client.adsClient.IsAllowed(context.Background(), grpcTD.InputBody.(*(adsPB.ContextRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_GET_GRANTED_ROLES:
		resp, err = test.Client.adsClient.GetAllGrantedRoles(context.Background(), grpcTD.InputBody.(*(adsPB.ContextRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	case METHOD_GET_GRANTED_PERMISSIONS:
		resp, err = test.Client.adsClient.GetAllPermissions(context.Background(), grpcTD.InputBody.(*(adsPB.ContextRequest)))
		if err == nil {
			grpcTD.OutputBody = resp
		}
		break
	}

	test.Client.CloseConnection()
	if err != nil {
		grpcTD.OutputMsg = err.Error()
	}

	return nil
}

//-------------Common util func for REST------------------------

//Verify RestTestData By Default
func VerifyGRpcTestByDefault(data interface{}, context *TestContext) bool {
	grpcTD, ok := data.(*GRpcTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to GRpcTestData")
		return false
	}

	//verify the error message code in response
	if len(grpcTD.ExpectedMsg) > 0 {
		if strings.Contains(grpcTD.OutputMsg, grpcTD.ExpectedMsg) {
			return true
		} else {
			TestLog.Fatalf("Expected ErrMsg=%s, but returned status ErrMsg=%s \n", grpcTD.ExpectedMsg, grpcTD.OutputMsg)
			return false
		}
	} else {
		//verify the response body
		equal := reflect.DeepEqual(grpcTD.ExpectedBody, grpcTD.OutputBody)
		if !equal {
			epc, _ := json.Marshal(grpcTD.ExpectedBody)
			act, _ := json.Marshal(grpcTD.OutputBody)

			//workaroud for issue 117.
			if string(act) == "null" && string(epc) == "[]" {
				TestLog.Log("===Nothing returned")
				equal = true
			} else {
				TestLog.Fatalf("Expected response body %s,\n but returned response body %s \n", epc, act)
			}
		}
		return equal
	}

}

//Verify RestTestData By Default
func VerifyGRpcAllGrantedRoleResult(data interface{}, context *TestContext) bool {
	grpcTD, ok := data.(*GRpcTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to GRpcTestData")
		return false
	}
	epcRoleResp, _ := grpcTD.ExpectedBody.(*adsPB.AllRoleResponse)
	actRoleResp, _ := grpcTD.OutputBody.(*adsPB.AllRoleResponse)
	sort.Strings(epcRoleResp.Roles)
	sort.Strings(actRoleResp.Roles)
	return reflect.DeepEqual(epcRoleResp.Roles, actRoleResp.Roles)
	//return CompareStringArray_NoOrder(epcRoleResp.Roles, actRoleResp.Roles)
}

//VerifyGRpcAllGrantedPermissions verify allGrantedPermissions response
func VerifyGRpcAllGrantedPermissions(data interface{}, context *TestContext) bool {
	grpcTD, ok := data.(*GRpcTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to GRpcTestData")
		return false
	}

	epcPermResp, _ := grpcTD.ExpectedBody.(*adsPB.AllPermissionResponse)
	actPermResp, _ := grpcTD.OutputBody.(*adsPB.AllPermissionResponse)

	if len(epcPermResp.Permissions) != len(actPermResp.Permissions) {
		TestLog.Fatalf("The length of permission in GRPC AllPermissionResponse is not equal. Actual/Expct=%v/%v", len(epcPermResp.Permissions), len(actPermResp.Permissions))
		return false
	}

	for i := 0; i < len(epcPermResp.Permissions); i++ {
		sort.Strings(epcPermResp.Permissions[i].Actions)
		sort.Strings(actPermResp.Permissions[i].Actions)
	}

	TestLog.Logf("VerifyGRpcAllGrantedPermissions,epcPermsResp=%v", epcPermResp)
	TestLog.Logf("VerifyGRpcAllGrantedPermissions,actPermsResp=%v", actPermResp)

	epcPerms := []string{}
	actPerms := []string{}

	for j := 0; j < len(epcPermResp.Permissions); j++ {
		epcPerms = append(epcPerms, fmt.Sprintf("%v-%v", epcPermResp.Permissions[j].Resource, epcPermResp.Permissions[j].Actions))
		actPerms = append(actPerms, fmt.Sprintf("%v-%v", actPermResp.Permissions[j].Resource, actPermResp.Permissions[j].Actions))
	}

	sort.Strings(epcPerms)
	sort.Strings(actPerms)

	if !reflect.DeepEqual(epcPerms, actPerms) {
		TestLog.Logf("VerifyGRpcAllGrantedPermissions,after change to string slice, epcPermsResp=%v", epcPerms)
		TestLog.Logf("VerifyGRpcAllGrantedPermissions,after change to string slice, actPermsResp=%v", actPerms)
		return false
	}
	return true
}
