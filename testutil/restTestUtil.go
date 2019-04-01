//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package testutil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	pmsapi "github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/cmd/flags"
	"github.com/oracle/speedle/pkg/svcs"
)

const (
	PMS_ENDPOINT    = "PMS_ENDPOINT"    //policy management endpoint
	ADS_ENDPOINT    = "ADS_ENDPOINT"    //authorization check endpoint
	PMS_ADMIN_TOKEN = "PMS_ADMIN_TOKEN" //token when do policy management
	ADS_ADMIN_TOKEN = "ADS_ADMIN_TOKEN" //token when do authorization check
	CA_LOCATION     = "CA_LOCATION"     //CA file absoulute location
	CERT_LOCATION   = "CERT_LOCATION"   //CERT file absoulute location
	KEY_LOCATION    = "KEY_LOCATION"    //public key file location
)

const (
	URI_POLICY_MGMT   = svcs.PolicyMgmtPath
	URI_IS_ALLOWD     = svcs.PolicyAtzPath + "is-allowed"
	URI_GRANTED_ROLES = svcs.PolicyAtzPath + "all-granted-roles"
	URI_GRANTED_PERMS = svcs.PolicyAtzPath + "all-granted-permissions"

	WAIT_POLICY_Initialize = 500 //ms
)

//--------------RestClient definition-----------

type RestClient struct {
	httpClient *http.Client
	prefix     string //the URI prefix when request rest api
	token      string //token would override userName and userPwd
	userName   string //the userName for basic authentication
	userPwd    string //the userPwd for basic authentication
}

//Set prefix of RestClient
func (client *RestClient) SetPrefix(newPrefix string) {
	client.prefix = newPrefix
}

//Get prefix of RestClient
func (client *RestClient) Prefix() string {
	return client.prefix
}

//Set token of RestClient
func (client *RestClient) SetToken(newToken string) {
	client.token = newToken
	client.userName = ""
	client.userPwd = ""
}

//Get token of RestClient
func (client *RestClient) Token() string {
	return client.token
}

//Set baisc auth userName/userPwd of RestClient
func (client *RestClient) SetBasicAuth(name string, pwd string) {
	client.userName = name
	client.userPwd = pwd
	client.token = ""
}

//Do REST GET api method with specified data (RestTestData)
func (client *RestClient) Get(data interface{}) error {
	restTD, ok := data.(*RestTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to RestTestData")
	}

	url := client.prefix + restTD.URI
	TestLog.Logf("GET URI=%s \n", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if len(client.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+client.token)
	} else if len(client.userName) > 0 {
		req.SetBasicAuth(client.userName, client.userPwd)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	restTD.actualStatus = resp.StatusCode
	TestLog.Logf("GET output=%v", resp)

	defer resp.Body.Close()
	if restTD.OutputBody != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		TestLog.Logf("GET Output body=%s", body)
		if err := json.Unmarshal(body, restTD.OutputBody); err != nil {
			TestLog.Logf("Fail to unmarshal JSON format, err=%s", err.Error())
			restTD.OutputBody = string(body)
			return nil
		}
	}
	return nil
}

//Do REST DELETE api method with specified data (RestTestData)
func (client *RestClient) Delete(data interface{}) error {
	restTD, ok := data.(*RestTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to RestTestData")
	}

	url := client.prefix + restTD.URI
	TestLog.Logf("DELETE URI=%s \n", url)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if len(client.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+client.token)
	} else if len(client.userName) > 0 {
		req.SetBasicAuth(client.userName, client.userPwd)
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	TestLog.Logf("DELETE output=%v", resp)
	restTD.actualStatus = resp.StatusCode
	defer resp.Body.Close()

	if restTD.OutputBody != nil {

		body, _ := ioutil.ReadAll(resp.Body)
		TestLog.Logf("DELETE Output body=%s", body)
		if err := json.Unmarshal(body, restTD.OutputBody); err != nil {
			TestLog.Logf("Fail to unmarshal JSON format, err=%s", err.Error())
			restTD.OutputBody = string(body)
			return nil
		}

	}
	return nil
}

//Do REST POST api method with specified data (RestTestData)
func (client *RestClient) Post(data interface{}) error {
	restTD, ok := data.(*RestTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to RestTestData")
	}

	absUri := client.prefix + restTD.URI
	var payload []byte
	isFormBody := false

	var req *http.Request = nil
	var err error

	TestLog.Logf("POST URI=%s \n", absUri)

	if restTD.InputBody != nil {
		if reflect.ValueOf(restTD.InputBody).Kind() == reflect.Map {
			isFormBody = true
			data := url.Values{}
			inputMap := restTD.InputBody.(map[string]string)
			for k, v := range inputMap {
				data.Add(k, v)
			}
			req, err = http.NewRequest(http.MethodPost, absUri, strings.NewReader(data.Encode()))
			TestLog.Logf("POST Form Input == %s \r\n", data.Encode())
		} else if reflect.ValueOf(restTD.InputBody).Kind() == reflect.String {
			payload := restTD.InputBody.(string)
			req, err = http.NewRequest(http.MethodPut, absUri, strings.NewReader(payload))
			TestLog.Logf("Put Input == %s \r\n", payload)
		} else {
			payload, err = json.Marshal(restTD.InputBody)
			if err != nil {
				return err
			}
			req, err = http.NewRequest(http.MethodPost, absUri, bytes.NewBuffer(payload))
			TestLog.Logf("POST Input == %s \r\n", payload)
		}
	} else {
		req, err = http.NewRequest(http.MethodPost, absUri, nil)
	}

	if err != nil {
		return err
	}

	if len(client.token) > 0 {
		TestLog.Log("set token Authorization header.")
		req.Header.Set("Authorization", "Bearer "+client.token)
	} else if len(client.userName) > 0 {
		TestLog.Log("set basic authz")
		req.SetBasicAuth(client.userName, client.userPwd)
	}

	req.Header.Set("Accept", "application/json")

	//Set content-type
	if isFormBody {
		TestLog.Log("set content-type : form.")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		TestLog.Log("set content-type : json.")
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return err
	}
	TestLog.Logf("POST output=%v", resp)
	restTD.actualStatus = resp.StatusCode
	defer resp.Body.Close()

	if restTD.OutputBody != nil {

		body, _ := ioutil.ReadAll(resp.Body)
		TestLog.Logf("POST Output body=%s", body)
		if err := json.Unmarshal(body, restTD.OutputBody); err != nil {
			TestLog.Logf("Fail to unmarshal JSON format, err=%s", err.Error())
			restTD.OutputBody = string(body)
			return nil
		}

	}
	return nil
}

//Do REST PUT api method with specified data (RestTestData)
func (client *RestClient) Put(data interface{}) error {
	restTD, ok := data.(*RestTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to RestTestData")
	}

	url := client.prefix + restTD.URI
	TestLog.Logf("PUT URI=%s \n", url)

	var req *http.Request = nil
	var err error

	if reflect.ValueOf(restTD.InputBody).Kind() == reflect.String {
		payload := restTD.InputBody.(string)
		req, err = http.NewRequest(http.MethodPut, url, strings.NewReader(payload))
		TestLog.Logf("Put Input == %s \r\n", payload)
	} else {

		payload, err := json.Marshal(restTD.InputBody)
		if err != nil {
			return err
		}
		TestLog.Logf("PUT Input == %s \r\n", payload)
		req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	}

	if err != nil {
		return err
	}
	if len(client.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+client.token)
		req.Header.Set("Accept", "application/json")

	} else if len(client.userName) > 0 {
		req.SetBasicAuth(client.userName, client.userPwd)
	}

	resp, err := client.httpClient.Do(req)
	restTD.actualStatus = resp.StatusCode
	if err != nil {
		return err
	}
	TestLog.Logf("PUT output=%v", resp)
	defer resp.Body.Close()

	if restTD.OutputBody != nil {

		body, _ := ioutil.ReadAll(resp.Body)
		TestLog.Logf("PUT Output body=%s", body)
		if err := json.Unmarshal(body, restTD.OutputBody); err != nil {
			TestLog.Logf("Fail to unmarshal output body, err=%s", err.Error())
			return nil
		}
	}
	return nil
}

//Get Rest Client
//if both curToken and basicAuth are specified, curToken would do effect
func NewRestClient(endpoint string, curToken string, basicAuthName string, basicAuthPwd string) (*RestClient, error) {
	caLoc := GetOSEnv(CA_LOCATION, "")
	certLoc := GetOSEnv(CERT_LOCATION, "")
	keyLoc := GetOSEnv(KEY_LOCATION, "")

	if len(caLoc) == 0 {
		return &RestClient{
			httpClient: &http.Client{},
			//prefix:     "http://" + endpoint,
			prefix:   endpoint,
			token:    curToken,
			userName: basicAuthName,
			userPwd:  basicAuthPwd,
		}, nil
	}
	fmt.Printf("caLoc: %s, certLoc: %s, keyLoc: %s\n", caLoc, certLoc, keyLoc)
	caCert, err := ioutil.ReadFile(caLoc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certLoc, keyLoc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	tlsConfig := &tls.Config{
		// TODO Remove InsecureSkipVerify, so that client verify server's cert
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
	}
	return &RestClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		prefix:   endpoint,
		token:    curToken,
		userName: basicAuthName,
		userPwd:  basicAuthPwd,
	}, nil
}

//Get RestClient for policy management
func NewRestClient_PMS() *RestClient {

	pmsEndpoint := GetOSEnv(PMS_ENDPOINT, "http://"+flags.DefaultPolicyMgmtEndPoint)
	pmsAdminToken := GetOSEnv(PMS_ADMIN_TOKEN, "")

	TestLog.Logf("PMSEndpoint=%s, token=%s \n", pmsEndpoint, pmsAdminToken)

	client, err := NewRestClient(pmsEndpoint, pmsAdminToken, "", "")

	if err != nil {
		TestLog.Fatalf("Init RestClient_pms failed due to error %v.", err)
		return nil
	}
	return client
}

//Get RestClient for Authz check
func NewRestClient_ADS() *RestClient {

	adsEndpoint := GetOSEnv(ADS_ENDPOINT, "http://"+flags.DefaultAuthzCheckEndPoint)
	adsToken := GetOSEnv(ADS_ADMIN_TOKEN, "")

	TestLog.Logf("ADSEndpoint=%s, token=%s", adsEndpoint, adsToken)

	client, err := NewRestClient(adsEndpoint, adsToken, "", "")

	if err != nil {
		TestLog.Fatalf("Init RestClient_ads failed due to error %v.", err)
		return nil
	}
	return client
}

//--------------RestTest and RestTestData definition-----------

//TestData for REST API testing
type RestTestData struct {
	URI            string
	actualStatus   int
	ExpectedStatus int
	InputBody      interface{} //input data for test
	OutputBody     interface{} //Actual output data
	ExpectedBody   interface{} //Expected output data
}

//Rest Test including RestClient and implementing Executer interface
type RestTest struct {
	Client *RestClient
}

//Prepare for Test Execution. Set the default func in testcase
func (test *RestTest) PreExecute(testcase *TestCase, ctx *TestContext) error {

	testcase.SetVerifyTestFunc(VerifyRestTestByDefault)

	switch testcase.Method {
	case METHOD_CREATE_SERVICE:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostCreateGetServiceTest)
		break
	case METHOD_GET_SERVICE:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostCreateGetServiceTest)
		break
	case METHOD_QUERY_SERVICE:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostListServiceTest)
		break
	case METHOD_DELETE_SERVICE:
		test.Client = NewRestClient_PMS()
		break
	case METHOD_CREATE_POLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostCreateGetPolicyTest)
		break
	case METHOD_GET_POLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		testcase.SetPostTestFunc(PostCreateGetPolicyTest)
		break
	case METHOD_QUERY_POLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostListPolicyTest)
		break
	case METHOD_DELETE_POLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		break
	case METHOD_CREATE_ROLEPOLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostCreateGetRolePolicyTest)
		break
	case METHOD_GET_ROLEPOLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		testcase.SetPostTestFunc(PostCreateGetRolePolicyTest)
		break
	case METHOD_QUERY_ROLEPOLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPostTestFunc(PostListRolePolicyTest)
		break
	case METHOD_DELETE_ROLEPOLICY:
		test.Client = NewRestClient_PMS()
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		break
	case METHOD_IS_ALLOWED:
		test.Client = NewRestClient_ADS()
		break
	case METHOD_GET_GRANTED_PERMISSIONS:
		test.Client = NewRestClient_ADS()
		break
	case METHOD_GET_GRANTED_ROLES:
		test.Client = NewRestClient_ADS()
		break
	default:
		return errors.New(ERROR_SPEEDLE_NOT_SUPPORTED)
	}

	return nil
}

//Execute current test with test data and context
func (test *RestTest) Execute(testcase *TestCase, ctx *TestContext) error {
	testData := testcase.Data

	switch testcase.Method {
	case METHOD_CREATE_SERVICE, METHOD_CREATE_POLICY, METHOD_CREATE_ROLEPOLICY,
		METHOD_IS_ALLOWED, METHOD_GET_GRANTED_ROLES, METHOD_GET_GRANTED_PERMISSIONS:
		return test.Client.Post(testData)
	case METHOD_GET_SERVICE, METHOD_QUERY_SERVICE, METHOD_GET_POLICY, METHOD_QUERY_POLICY,
		METHOD_GET_ROLEPOLICY, METHOD_QUERY_ROLEPOLICY:
		return test.Client.Get(testData)
	case METHOD_DELETE_SERVICE, METHOD_DELETE_POLICY, METHOD_DELETE_ROLEPOLICY:
		return test.Client.Delete(testData)
	default:
		return errors.New(ERROR_SPEEDLE_NOT_SUPPORTED)
	}
}

//New RestTest
func NewRestTestExecuter() TestExecuter {
	return &RestTest{
		Client: nil,
	}
}

//-------------Common util func for REST------------------------

//Verify RestTestData By Default
func VerifyRestTestByDefault(data interface{}, context *TestContext) bool {
	restTD, ok := data.(*RestTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to RestTestData")
		return false
	}

	//verify the status code in response
	if restTD.actualStatus != restTD.ExpectedStatus {
		TestLog.Fatalf("Expected status %d, but returned status %d for uri %s.", restTD.ExpectedStatus, restTD.actualStatus, restTD.URI)
		return false
	}

	removeMetaData(reflect.ValueOf(restTD.OutputBody))
	//verify the response body
	equal := reflect.DeepEqual(restTD.ExpectedBody, restTD.OutputBody)
	if !equal {
		epc, _ := json.Marshal(restTD.ExpectedBody)
		act, _ := json.Marshal(restTD.OutputBody)

		//workaroud for issue 117.
		if string(act) == "null" && string(epc) == "[]" {
			TestLog.Log("===Nothing returned")
			equal = true
		} else {
			TestLog.Fatalf("Expected response body %s,\r\n but returned response body %s for uri %s.", epc, act, restTD.URI)
		}
	}

	return equal
}

//GetRestTestData for creating service
func GetRestTestData_CreateService(step string, srvName string) TestCase {
	return TestCase{
		Name:     fmt.Sprintf("Step%s-AddService", step),
		Enabled:  true,
		Executer: NewRestTestExecuter,
		Method:   METHOD_CREATE_SERVICE,
		Data: &RestTestData{
			URI:            URI_POLICY_MGMT + "service",
			ExpectedStatus: 201,
			InputBody: &pmsapi.Service{
				Name: srvName,
				Type: pmsapi.TypeApplication,
			},
			OutputBody: &pmsapi.Service{},
			ExpectedBody: &pmsapi.Service{
				Name: srvName,
				Type: pmsapi.TypeApplication,
			},
		},
		PreTestFunc:    nil,
		PostTestFunc:   nil,
		VerifyTestFunc: VerifyRestTestByDefault,
	}
}

//GetRestTestData for deleting service
func GetRestTestData_DeleteService(step string, srvName string) TestCase {
	return TestCase{
		Name:     fmt.Sprintf("Step%s-Delete Service", step),
		Executer: NewRestTestExecuter,
		Method:   METHOD_DELETE_SERVICE,
		Data: &RestTestData{
			URI:            URI_POLICY_MGMT + "service/" + srvName,
			ExpectedStatus: 204,
		},
	}
}

//GetRestTestData for creating rolepolicy
func GetRestTestData_CreateRolePolicy(step string, srvName string, rolePolicyName string, effect string, roles []string, principals []string, resources []string) TestCase {
	return TestCase{
		Name:     fmt.Sprintf("Step%s-Add RolePolicy", step),
		Executer: NewRestTestExecuter,
		Method:   METHOD_CREATE_ROLEPOLICY,
		Data: &RestTestData{
			URI:            URI_POLICY_MGMT + "service/" + srvName + "/role-policy",
			ExpectedStatus: 201,
			InputBody: &pmsapi.RolePolicy{
				Name:       rolePolicyName,
				Effect:     effect,
				Roles:      roles,
				Principals: principals,
				Resources:  resources,
			},
			OutputBody: &pmsapi.RolePolicy{},
			ExpectedBody: &pmsapi.RolePolicy{
				Name:       rolePolicyName,
				Effect:     effect,
				Roles:      roles,
				Principals: principals,
				Resources:  resources,
			},
		},
		PostTestFunc: PostCreateGetRolePolicyTest,
	}
}

//GetRestTestData for creating policy
func GetRestTestData_CreatePolicy(step string, srvName string, policyName string, effect string, principals [][]string, resource string, actions []string) TestCase {
	return TestCase{
		Name:     fmt.Sprintf("Step%s-Add Policy", step),
		Executer: NewRestTestExecuter,
		Method:   METHOD_CREATE_POLICY,
		Data: &RestTestData{
			URI:            URI_POLICY_MGMT + "service/" + srvName + "/policy",
			ExpectedStatus: 201,
			InputBody: &pmsapi.Policy{
				Name:   policyName,
				Effect: effect,
				Permissions: []*pmsapi.Permission{
					{
						Resource: resource,
						Actions:  actions,
					},
				},
				Principals: principals,
			},
			OutputBody: &pmsapi.Policy{},
			ExpectedBody: &pmsapi.Policy{
				Name:   policyName,
				Effect: effect,
				Permissions: []*pmsapi.Permission{
					{
						Resource: resource,
						Actions:  actions,
					},
				},
				Principals: principals,
			},
		},
		PostTestFunc: PostCreateGetPolicyTest,
	}
}

//Get TestData for sleeping a while (in ms)
func GetTestData_Sleep(sleep int) TestCase {
	sleepStr := fmt.Sprintf("sleep_%d", sleep)
	return TestCase{
		Name:   sleepStr,
		Method: sleepStr,
	}
}
