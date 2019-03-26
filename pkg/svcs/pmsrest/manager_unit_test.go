//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package pmsrest

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"os"
	"testing"

	"time"

	"encoding/json"
	"io/ioutil"

	"log"

	pmsapi "github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/cfg"
	"github.com/oracle/speedle/pkg/store"
	_ "github.com/oracle/speedle/pkg/store/file"
	"github.com/oracle/speedle/pkg/svcs"
)

var storeFile = "./fakestore.json"
var creator = "creator"
var testserver *httptest.Server

func NewTestServer() (*httptest.Server, error) {
	conf := GenerateServerConfig()
	ps, err := store.NewStore(conf.StoreConfig.StoreType, conf.StoreConfig.StoreProps)
	if err != nil {
		return nil, err
	}
	routers, err := NewRouter(ps)
	if err != nil {
		return nil, err
	}
	server := httptest.NewUnstartedServer(routers)
	server.Start()

	return server, nil
}

func GenerateServerConfig() *cfg.Config {
	var conf cfg.Config
	var storeConf cfg.StoreConfig
	storeConf.StoreType = cfg.StorageTypeFile
	storeConf.StoreProps = make(map[string]interface{})
	storeConf.StoreProps["FileLocation"] = storeFile
	conf.StoreConfig = &storeConf
	conf.EnableWatch = false
	return &conf
}

func checkCreateMetaData(metaData map[string]string, t *testing.T) {
	createTimeStr := metaData["createtime"]
	if len(createTimeStr) > 0 {
		createTime, err := time.Parse(time.RFC3339, createTimeStr)
		if err != nil {
			t.Fatal("failed to parse createtime in response. error:", err)
		} else if time.Now().Sub(createTime).Seconds() > 5 {
			t.Fatal("createtime is not reasonable. createtime:", createTime)
		}
	}
	if len(metaData["createby"]) == 0 || metaData["createby"] != creator {
		t.Fatal("createby field is not expected. createby:", metaData["createby"])
	}
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	err := ioutil.WriteFile(storeFile, []byte(`{"services":[{"name":"fakeservice","type":"app"}]}`), 0644)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	defer os.Remove(storeFile)
	testserver, err = NewTestServer()
	if err != nil {
		log.Fatal("failed to start test server. error:", err)
		return 1
	}
	defer testserver.Close()
	return m.Run()
}

func TestCreateServicePrincipalHeader(t *testing.T) {
	var service = pmsapi.Service{
		Name: "service1",
		Type: "app"}
	var policy = pmsapi.Policy{
		Name:   "p1",
		Effect: "deny",
	}
	var rpolicy = pmsapi.RolePolicy{
		Name:   "rp1",
		Effect: "deny",
	}
	service.Policies = []*pmsapi.Policy{&policy}
	service.RolePolicies = []*pmsapi.RolePolicy{&rpolicy}

	serviceData, err := json.Marshal(service)
	if err != nil {
		t.Fatal("failed to marsh service data")
	}
	req, err := http.NewRequest("POST", testserver.URL+svcs.PolicyMgmtPath+"service", bytes.NewBuffer(serviceData))
	if err != nil {
		t.Fatal("failed to make test request")
	}

	addPrincipalHeader(req)
	var client *http.Client
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	if resp.StatusCode == http.StatusCreated {
		serviceGot := pmsapi.Service{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read create service response.")
		}
		err = json.Unmarshal(body, &serviceGot)
		if err != nil {
			t.Fatal("failed to unmarsh create service response.")
		}
		t.Log("check service")
		checkCreateMetaData(serviceGot.Metadata, t)
		t.Log("check policy")
		checkCreateMetaData(serviceGot.Policies[0].Metadata, t)
		t.Log("check rolepolicy")
		checkCreateMetaData(serviceGot.RolePolicies[0].Metadata, t)
	} else {
		t.Fatal("failed to create service. status:", resp.StatusCode)
	}

}

func TestCreatePolicyPrincipalHeader(t *testing.T) {
	var policy = pmsapi.Policy{
		Name:   "p1",
		Effect: "deny",
	}
	policyData, err := json.Marshal(policy)
	if err != nil {
		t.Fatal("failed to marsh policy data")
	}
	req, err := http.NewRequest("POST", testserver.URL+svcs.PolicyMgmtPath+"service/fakeservice/policy", bytes.NewBuffer(policyData))
	if err != nil {
		t.Fatal("failed to make test request")
	}
	addPrincipalHeader(req)
	var client *http.Client
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	if resp.StatusCode == http.StatusCreated {
		policyGot := pmsapi.Policy{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response.")
		}
		err = json.Unmarshal(body, &policyGot)
		if err != nil {
			t.Fatal("failed to unmarsh response.")
		}
		checkCreateMetaData(policyGot.Metadata, t)
	} else {
		t.Fatal("failed to create policy. status:", resp.StatusCode)
	}
}

func TestCreateRolePolicyPrincipalHeader(t *testing.T) {

	var rpolicy = pmsapi.RolePolicy{
		Name:   "rp1",
		Effect: "deny",
	}
	policyData, err := json.Marshal(rpolicy)
	if err != nil {
		t.Fatal("failed to marsh rolepolicy data")
	}
	req, err := http.NewRequest("POST", testserver.URL+svcs.PolicyMgmtPath+"service/fakeservice/role-policy", bytes.NewBuffer(policyData))
	if err != nil {
		t.Fatal("failed to make test request")
	}
	addPrincipalHeader(req)
	var client *http.Client
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	if resp.StatusCode == http.StatusCreated {
		policyGot := pmsapi.RolePolicy{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response.")
		}
		err = json.Unmarshal(body, &policyGot)
		if err != nil {
			t.Fatal("failed to unmarsh response.")
		}
		checkCreateMetaData(policyGot.Metadata, t)
	} else {
		t.Fatal("failed to create rolepolicy. status:", resp.StatusCode)
	}
}

func TestCreateFunctionPrincipalHeader(t *testing.T) {
	var customFunc = pmsapi.Function{Name: "f1", FuncURL: "http://fakeurl", ResultCachable: false, ResultTTL: 256}
	funcData, err := json.Marshal(customFunc)
	if err != nil {
		t.Fatal("failed to marsh function data")
	}
	req, err := http.NewRequest("POST", testserver.URL+svcs.PolicyMgmtPath+"function", bytes.NewBuffer(funcData))
	if err != nil {
		t.Fatal("failed to make test request")
	}
	addPrincipalHeader(req)
	var client *http.Client
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("failed get response")
	}
	if resp.StatusCode == http.StatusCreated {
		funcGot := pmsapi.Function{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response.")
		}
		err = json.Unmarshal(body, &funcGot)
		if err != nil {
			t.Fatal("failed to unmarsh response.")
		}
		checkCreateMetaData(funcGot.Metadata, t)
	} else {
		t.Fatal("failed to create rolepolicy. status:", resp.StatusCode)
	}
}

func addPrincipalHeader(req *http.Request) {
	/*user := &ads.Principal{"user", creator, "wercker"}
	group := &ads.Principal{"group", "group1", "wercker"}
	var principals = []*ads.Principal{user, group}
	data, _ := json.Marshal(principals)*/
	req.Header.Add(svcs.PrincipalsHeader, creator)
}
