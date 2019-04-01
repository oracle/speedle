//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package testutil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	pmsapi "github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/testutil/msg"
	"github.com/oracle/speedle/testutil/param"
)

const (
	SP_APP_NAME = "SP_APP_NAME"
)

//Test Methods for Speedle
const (
	METHOD_CONFIG = "config"
)

//--------------CmdClient definition-----------

type CmdClient struct {
	cmd    string
	output string
}

//execute spxctl command with specified parameter and test data
func (client *CmdClient) ExecuteCmd(parameter string, data *CmdTestData) error {
	cmdStr := client.cmd + parameter

	TestLog.Log("\r\n\r\n-----------------execute command-------------------")
	TestLog.Log("Command string= " + cmdStr)

	c := Command("/bin/bash", "-c", cmdStr)
	c.Run()
	client.output = c.Stdout()
	TestLog.Log("The output of command is:-------")
	TestLog.Log(client.output)

	// When error occurs, the exit code should not be 0, so comment out this
	// if !c.Success() {
	// 	TestLog.Log("Expected to succeed, but failed")
	// 	return errors.New("Fail to execute command")
	// }

	data.OutputMsg = client.output
	return nil
}

//Get the test client for Cmd
func NewCmdClient() *CmdClient {
	var cmdStr string
	binPath := os.Getenv("GOPATH")
	if strings.Compare(string(binPath[len(binPath)-1:]), "/") != 0 {
		binPath = binPath + "/"
	}

	cmdStr = binPath + "bin/" + os.Getenv(SP_APP_NAME) + " "
	token := GetOSEnv(PMS_ADMIN_TOKEN, "")
	if len(token) > 0 {
		cmdStr = cmdStr + " --token " + token + " "
	}

	return &CmdClient{
		cmd:    cmdStr,
		output: "",
	}
}

//Get the test client for Cmd
func NewCmdClient_token(token string) *CmdClient {
	var cmdStr string
	binPath := os.Getenv("GOPATH")
	if strings.Compare(string(binPath[len(binPath)-1:]), "/") != 0 {
		binPath = binPath + "/"
	}

	cmdStr = binPath + "bin/" + os.Getenv(SP_APP_NAME) + " "
	if len(token) > 0 {
		cmdStr = cmdStr + " --token " + token + " "
	}

	return &CmdClient{
		cmd:    cmdStr,
		output: "",
	}
}

//--------------CmdTest and CmdTestData definition-----------

//test data for spctl command
type CmdTestData struct {
	Param        string      //parameter string after spxctl command
	FileContent  interface{} //The content of File (JSON or PDL )
	OutputMsg    string      //Actual output Message
	OutputBody   interface{} //Actual output body for json object
	ExpectedMsg  string      //expected message for the command
	ExpectedBody interface{} //expected body about json object string
}

type CmdTest struct {
	Client *CmdClient
}

func NewCmdTest() TestExecuter {
	return &CmdTest{}
}

//Prepare for Test Execution. Set the default func in testcase
func (test *CmdTest) PreExecute(testcase *TestCase, ctx *TestContext) error {
	test.Client = NewCmdClient()
	testcase.SetVerifyTestFunc(VerifyCmdTestByDefault)

	switch testcase.Method {
	case METHOD_CREATE_SERVICE:
		testcase.SetPostTestFunc(PostCreateGetServiceTest)
		break
	case METHOD_GET_SERVICE:
		testcase.SetPostTestFunc(PostCreateGetServiceTest)
		break
	case METHOD_QUERY_SERVICE:
		testcase.SetPostTestFunc(PostListServiceTest)
		break
	case METHOD_CONFIG:
		test.Client = NewCmdClient_token("")
		break
	case METHOD_DELETE_SERVICE:
		break
	case METHOD_CREATE_POLICY:
		testcase.SetPostTestFunc(PostCreateGetPolicyTest)
		break
	case METHOD_GET_POLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		testcase.SetPostTestFunc(PostCreateGetPolicyTest)
		break
	case METHOD_QUERY_POLICY:
		testcase.SetPostTestFunc(PostListPolicyTest)
		break
	case METHOD_DELETE_POLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		break
	case METHOD_CREATE_ROLEPOLICY:
		testcase.SetPostTestFunc(PostCreateGetRolePolicyTest)
		break
	case METHOD_GET_ROLEPOLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		testcase.SetPostTestFunc(PostCreateGetRolePolicyTest)
		break
	case METHOD_QUERY_ROLEPOLICY:
		testcase.SetPostTestFunc(PostListRolePolicyTest)
		break
	case METHOD_DELETE_ROLEPOLICY:
		testcase.SetPreTestFunc(PreGetDeletePolicyTest)
		break
	default:
		return errors.New(ERROR_SPEEDLE_NOT_SUPPORTED)
	}

	return nil
}

//Execute current test with test data and context
func (test *CmdTest) Execute(testcase *TestCase, ctx *TestContext) error {
	testData := testcase.Data
	cmdTD := testData.(*CmdTestData)

	err := test.Client.ExecuteCmd(cmdTD.Param, cmdTD)
	if err != nil {
		return err
	}

	switch testcase.Method {
	case METHOD_CREATE_SERVICE, METHOD_CREATE_POLICY, METHOD_CREATE_ROLEPOLICY:
		return test.ParseOutputForCreate(testData)
	case METHOD_GET_SERVICE, METHOD_QUERY_SERVICE, METHOD_GET_POLICY, METHOD_QUERY_POLICY,
		METHOD_GET_ROLEPOLICY, METHOD_QUERY_ROLEPOLICY, METHOD_CONFIG:
		return test.ParseOutputForGet(testData)
	}
	return nil
}

//Parse output result for creating service/policy/rolepolicy
func (test *CmdTest) ParseOutputForCreate(testData interface{}) error {
	cmdTD := testData.(*CmdTestData)
	tmp := test.Client.output
	//The output message always ends with \n"
	start := strings.Index(tmp, "\n")
	if start > 0 && len(tmp[start+1:]) > 0 && strings.Index(tmp[start+1:], "{") >= 0 {
		cmdTD.OutputMsg = tmp[:start]
		err := json.Unmarshal([]byte(tmp[start+1:]), cmdTD.OutputBody)
		if err != nil {
			TestLog.Logf("Fail to unmarshall output into json object %s\n", cmdTD.OutputBody)
			return err
		}
	} else {
		cmdTD.OutputMsg = tmp
	}
	return nil
}

//Parse output result for getting service/policy/rolepolicy
func (test *CmdTest) ParseOutputForGet(data interface{}) error {
	cmdTD := data.(*CmdTestData)

	tmp := test.Client.output
	if len(tmp) > 0 && strings.Index(tmp, "{") >= 0 {
		err := json.Unmarshal([]byte(tmp), &cmdTD.OutputBody)
		if err != nil {
			TestLog.Log("Fail to unmarshall output into json object. Err=" + err.Error())
			return err
		}
	} else {
		//cmdTD.OutputMsg = tmp
		TestLog.Log("output is empty or not in json format")
	}
	return nil
}

//-------------Common util func for REST------------------------

//Verify CmdTestData By Default
func VerifyCmdTestByDefault(data interface{}, context *TestContext) bool {
	cmdTD, ok := data.(*CmdTestData)
	if !ok {
		TestLog.Fatalf("Fail to convert data to CmdTestData")
		return false
	}

	if !strings.Contains(cmdTD.OutputMsg, cmdTD.ExpectedMsg) {
		TestLog.Fatalf("Actual ouput message is %s, doesn't contain expected string %s", cmdTD.OutputMsg, cmdTD.ExpectedMsg)
	}
	removeMetaData(reflect.ValueOf(cmdTD.OutputBody))
	//verify the response body
	equal := reflect.DeepEqual(cmdTD.ExpectedBody, cmdTD.OutputBody)
	if !equal {

		epc, _ := json.Marshal(cmdTD.ExpectedBody)
		act, _ := json.Marshal(cmdTD.OutputBody)

		//workaroud for issue 117.
		if string(act) == "null" && string(epc) == "[]" {
			TestLog.Log("===Nothing returned")
			equal = true
		} else {
			TestLog.Fatalf("Expected response body %s,\r\n but returned response body %s for cmd %s.", epc, act, cmdTD.Param)
		}
	}

	return equal
}

func checkFile(fileNameWithPath string) {
	if _, err := os.Stat(fileNameWithPath); !os.IsNotExist(err) {
		TestLog.Log("File " + fileNameWithPath + " exist, remove it firstly")
		os.Remove(fileNameWithPath)
	}

}

func createFile(fileNameWithPath string, content []byte) {
	err := ioutil.WriteFile(fileNameWithPath, content, 0644)
	if err != nil {
		TestLog.Log("Fail to write to file. Err=" + err.Error())
	}
}

//generate a file according to pdl list
func GeneratePdlFile(fileNameWithPath string, pdlList []string) {
	var tmp = ""
	if pdlList != nil {
		for _, v := range pdlList {
			tmp += v + "\n"
		}
	}
	buf := []byte(tmp)
	createFile(fileNameWithPath, buf)
}

//generate a json file according to service object
func GenerateJsonFileWithService(fileNameWithPath string, service *pmsapi.Service) {
	tmp, err := json.Marshal(service)
	if err != nil {
		TestLog.Log("Fail to Marshal pmsapi.Service object")
	}
	createFile(fileNameWithPath, tmp)
}

//generate a json file according to policy object
func GenerateJsonFileWithPolicy(fileNameWithPath string, policy *pmsapi.Policy) {
	tmp, err := json.Marshal(policy)
	if err != nil {
		TestLog.Log("Fail to Marshal pmsapi.Policy object")
	}
	createFile(fileNameWithPath, tmp)
}

//generate a json file according to role policy object
func GenerateJsonFileWithRolePolicy(fileNameWithPath string, service *pmsapi.RolePolicy) {
	tmp, err := json.Marshal(service)
	if err != nil {
		TestLog.Log("Fail to Marshal pmsapi.RolePolicy object")
	}
	createFile(fileNameWithPath, tmp)
}

//Create service with PDL file
func (client *CmdClient) CreateServiceWithPDL(serviceName string, pdl []string) bool {
	//delete service firstly
	cmdStr := client.cmd + param.DELETE_SERVICE(serviceName)
	c := Command("/bin/bash", "-c", cmdStr)
	c.Run()
	log.Println("DeleteService firstly. Output is: " + c.stdout)

	//create service
	tmpPDLFile := "/tmp/pdl-file"
	GeneratePdlFile(tmpPDLFile, pdl)
	cmdStr = client.cmd + param.CREATE_SERVICE_WITH_PDLFILE(serviceName, pmsapi.TypeApplication, tmpPDLFile)
	log.Println("CreateServiceWithPDL cmdstr= " + cmdStr)
	c = Command("/bin/bash", "-c", cmdStr)
	c.Run()
	client.output = c.Stdout()
	log.Println("CreateServiceWithPDL cmd output-" + client.output)
	if !strings.Contains(client.output, msg.OUTPUT_SERVICE_CREATED()) {
		log.Println("Service created failed.")
		return false
	}
	return true
}
