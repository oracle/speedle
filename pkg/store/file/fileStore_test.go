//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/store"
)

var storeConfig map[string]interface{} = make(map[string]interface{})

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	defer os.Remove("ps.json")
	storeConfig["FileLocation"] = "./ps.json"

	return m.Run()
}

func TestWriteReadPolicyStore(t *testing.T) {
	store, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new file store:", err)
	}
	psr, err := store.ReadPolicyStore()
	if err != nil {
		t.Fatal("fail to read policy store:", err)
	} else {
		t.Log("app num in the store is:", len(psr.Services))
	}

	var ps pms.PolicyStore
	for i := 0; i < 10; i++ {
		service := pms.Service{Name: fmt.Sprintf("app%d", i)}
		ps.Services = append(ps.Services, &service)
	}
	err = store.WritePolicyStore(&ps)
	if err != nil {
		t.Fatal("fail to write policy store:", err)
	}

	psr, err = store.ReadPolicyStore()
	if err != nil {
		t.Fatal("fail to read policy store:", err)
	}
	if 10 != len(psr.Services) {
		t.Error("should have 10 services in the store")
	}
	for _, service := range psr.Services {
		log.Printf(service.Name)
	}

}

func TestWriteReadService(t *testing.T) {
	store, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new file store:", err)
	}

	service := pms.Service{Name: "app1", Type: pms.TypeApplication}
	err = store.(*Store).WriteService(&service)
	if err != nil {
		t.Fatal("fail to write service:", err)
	}
	servicer, errr := store.GetService("app1")
	if errr != nil {
		t.Fatal("fail to read service:", err)
	}
	if "app1" != servicer.Name {
		t.Error("app name should be app1")
	}
	err = store.DeleteService("app1")
	if err != nil {
		t.Fatal("fail to delete application:", err)
	}
	servicer, err = store.GetService("app1")
	t.Log(err)
	if err == nil {
		t.Fatal("should fail as app is already deleted")
	}

}

func TestFileStore_GetPolicyByName(t *testing.T) {
	store, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	//clean the service firstly
	serviceName := "service1"
	err = store.DeleteService(serviceName)
	t.Log("deleteing service1, err:", err)

	app := pms.Service{Name: serviceName, Type: pms.TypeApplication}
	num := 10
	i := 0
	for i < num {
		var policy pms.Policy
		policy.Name = fmt.Sprintf("policy%d", i)
		policy.Effect = "grant"
		policy.Permissions = []*pms.Permission{
			{
				Resource: "/node1",
				Actions:  []string{"get", "create", "delete"},
			},
		}
		policy.Principals = [][]string{{"user:Alice"}}
		app.Policies = append(app.Policies, &policy)
		i++
	}
	blankNamePolicy := pms.Policy{
		Effect: "grant",
		Permissions: []*pms.Permission{
			{
				Resource: "/node1",
				Actions:  []string{"get", "create", "delete"},
			},
		},
		Principals: [][]string{{"user:Alice"}},
	}
	app.Policies = append(app.Policies, &blankNamePolicy)
	duplicateNamePolicy := pms.Policy{
		Name:   "policy0",
		Effect: "grant",
		Permissions: []*pms.Permission{
			{
				Resource: "/node1",
				Actions:  []string{"get", "create", "delete"},
			},
		},
		Principals: [][]string{{"user:Alice"}},
	}
	app.Policies = append(app.Policies, &duplicateNamePolicy)

	err = store.CreateService(&app)
	if err != nil {
		t.Log("fail to create application:", err)
		t.FailNow()
	}
	service, errr := store.GetService(serviceName)
	if errr != nil {
		t.Log("fail to get application:", err)
		t.FailNow()
	}
	poilcyName := "policy0"

	policyArrListed, err := store.ListAllPolicies(service.Name, "name eq "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}

	if len(policyArrListed) != 2 { //2 policy0 policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name co "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 2 { //2 policy0 policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name sw "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 2 { //2 policy0 policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name gt "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != num-1 { //all policy name great than policy0
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name ge "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != num+1 { //all policy name great than or equals to policy0
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name lt "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 1 { //1 blank name policy
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name le "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 3 { //1 blank name policy and 2 duplicate policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name le ''")
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 1 { //1 blank name policy
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllPolicies(service.Name, "name pr")
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != num+1 {
		t.Fatal("Get none blank name poclies failed! ")
	}

}

func TestFileStore_GetRolePolicyByName(t *testing.T) {
	store, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	//clean the service firstly
	serviceName := "service1"
	err = store.DeleteService(serviceName)
	t.Log("deleteing service1, err:", err)

	app := pms.Service{Name: serviceName, Type: pms.TypeApplication}
	num := 1000
	i := 0
	for i < num {
		var rolePolicy pms.RolePolicy
		rolePolicy.Name = fmt.Sprintf("rp%d", i)
		rolePolicy.Effect = "grant"
		rolePolicy.Roles = []string{fmt.Sprintf("role%d", i)}
		rolePolicy.Principals = []string{"user:Alice"}
		app.RolePolicies = append(app.RolePolicies, &rolePolicy)
		i++
	}
	blankNameRolePolicy := pms.RolePolicy{
		Effect:     "grant",
		Roles:      []string{fmt.Sprintf("role%d", i)},
		Principals: []string{"user:Alice"},
	}
	app.RolePolicies = append(app.RolePolicies, &blankNameRolePolicy)

	duplicateNameRolePolicy := pms.RolePolicy{
		Name:       "rp0",
		Effect:     "grant",
		Roles:      []string{fmt.Sprintf("role%d", i)},
		Principals: []string{"user:Alice"},
	}
	app.RolePolicies = append(app.RolePolicies, &duplicateNameRolePolicy)

	err = store.CreateService(&app)
	if err != nil {
		t.Log("fail to create application:", err)
		t.FailNow()
	}
	service, errr := store.GetService(serviceName)
	if errr != nil {
		t.Log("fail to get application:", err)
		t.FailNow()
	}
	poilcyName := "rp0"

	policyArrListed, err := store.ListAllRolePolicies(service.Name, "name eq "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}

	if len(policyArrListed) != 2 { //2 policy0 policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name co "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 2 { //2 policy0 policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name sw "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 2 { //2 policy0 policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name gt "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != num-1 { //all policy name great than policy0
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name ge "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != num+1 { //all policy name great than or equals to policy0
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name lt "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 1 { //1 blank name policy
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name le "+poilcyName)
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 3 { //1 blank name policy and 2 duplicate policies
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name le ''")
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != 1 { //1 blank name policy
		t.Fatal("get poilcy by name didn't get expected policies! ")
	}

	policyArrListed, err = store.ListAllRolePolicies(service.Name, "name pr")
	if err != nil {
		t.Fatal("Failed to list polices for service:", service.Name, err)
	}
	if len(policyArrListed) != num+1 {
		t.Fatal("Get none blank name poclies failed! ")
	}

}

func TestFunctionManagement(t *testing.T) {
	store, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	//clean the store
	store.DeleteFunctions()

	testFunc := &pms.Function{
		Name:           "testFunc",
		Description:    "test function",
		FuncURL:        "https://localhost:23456/testFunc",
		ResultCachable: true,
		ResultTTL:      300,
		CA:             "-----BEGIN CERTIFICATE-----\nMIID7TCCAtWgAwIBAgIJALM3l/OZ9uJKMA0GCSqGSIb3DQEBCwUAMIGMMQswCQYD\nVQQGEwJjbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzEPMA0G\nA1UECgwGb3JhY2xlMQwwCgYDVQQLDANpZG0xEjAQBgNVBAMMCWxvY2FsaG9zdDEm\nMCQGCSqGSIb3DQEJARYXY3ludGhpYS5kaW5nQG9yYWNsZS5jb20wHhcNMTgwNDI1\nMDc1MDMwWhcNMTkwNDI1MDc1MDMwWjCBjDELMAkGA1UEBhMCY24xEDAOBgNVBAgM\nB2JlaWppbmcxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBm9yYWNsZTEMMAoG\nA1UECwwDaWRtMRIwEAYDVQQDDAlsb2NhbGhvc3QxJjAkBgkqhkiG9w0BCQEWF2N5\nbnRoaWEuZGluZ0BvcmFjbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAn/AFElluGOZYfvlzBGfHfkd/Q9SuQFsSnQt7Qp63Yuf5Ie/q4NACzWPC\nB/L6nQrut4OMxJHvhVAswJozRZrQxXvX/vUxkg+TmALj3U9ejF/5arGtjy5v+yGi\nwci7zM4r7VNFJGRkfluNRC1kJi4AY6jk6Gl4d/bX4tBXE8mEFY1rUswYtat3OMja\njVAoocClk6WcaQuK9R1uB+BPyxHLJ04RyKRuepPYRBQjgwHK5kMF3s5p07Os+2JH\n5jyJYW2NPs6pQe0k8GWpaar/yZ2eut9gsgHnu5JCWnyedo4nEx6I/G4GSaX+0SeU\n/Wb2aqq1QGfVOESml7CVcEa/buTeUwIDAQABo1AwTjAdBgNVHQ4EFgQU5i7CO32N\nspQ5AaG/aRU0LX2koYwwHwYDVR0jBBgwFoAU5i7CO32NspQ5AaG/aRU0LX2koYww\nDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAbQuCMPK8f8QuEmTpZBFv\naka9qruT/0/TrxrbxEh68N4moXSTVv4tSrDTmdkwUiiwayuGS7fvKjSV6hwGkQbV\nzGbFDdwOw1tPE2OwnA7/+RPl4KmE4iTHnnIanyg9CKmBW/tMp/vUyv5nIt7Xw5n4\ntx3C9/hme+Rlx+SVPIAwAjl4nVFNLfzyG+JDBnQWygySm88SzzK0WRgh5V+gyXCK\nucDW5rA6X9/CM3QrSY50mSM6dbyYDMtmTI4dX7E9STTBCNsNNcmgYkX0N9lm5RoF\nuBsAcPmp1SVIbXelDHJiIXxMKzwZy8riZQ8+Dw6LMs6wZX7COVvMWN4Dfcuo89av\nIQ==\n-----END CERTIFICATE-----",
	}
	//test create function
	_, err = store.CreateFunction(testFunc)
	if err != nil {
		t.Fatal("Failed to create function:", err)
	}

	//test get function
	_, err = store.GetFunction(testFunc.Name)
	if err != nil {
		t.Fatal("Failed to get function:", err)
	}

	//test delete function
	err = store.DeleteFunction(testFunc.Name)
	if err != nil {
		t.Fatal("Failed to delete function:", err)
	}
	_, err = store.GetFunction(testFunc.Name)
	if err == nil {
		t.Fatal("Should failed to get function as it is delete")
	}

	//test listAllFunctions
	i := 0
	for i < 10 {
		testFunc := &pms.Function{
			Name:           "testFunc" + strconv.Itoa(i),
			Description:    "test function" + strconv.Itoa(i),
			FuncURL:        "https://localhost:23456/testFunc" + strconv.Itoa(i),
			ResultCachable: true,
			ResultTTL:      300,
			CA:             "-----BEGIN CERTIFICATE-----\nMIID7TCCAtWgAwIBAgIJALM3l/OZ9uJKMA0GCSqGSIb3DQEBCwUAMIGMMQswCQYD\nVQQGEwJjbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzEPMA0G\nA1UECgwGb3JhY2xlMQwwCgYDVQQLDANpZG0xEjAQBgNVBAMMCWxvY2FsaG9zdDEm\nMCQGCSqGSIb3DQEJARYXY3ludGhpYS5kaW5nQG9yYWNsZS5jb20wHhcNMTgwNDI1\nMDc1MDMwWhcNMTkwNDI1MDc1MDMwWjCBjDELMAkGA1UEBhMCY24xEDAOBgNVBAgM\nB2JlaWppbmcxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBm9yYWNsZTEMMAoG\nA1UECwwDaWRtMRIwEAYDVQQDDAlsb2NhbGhvc3QxJjAkBgkqhkiG9w0BCQEWF2N5\nbnRoaWEuZGluZ0BvcmFjbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAn/AFElluGOZYfvlzBGfHfkd/Q9SuQFsSnQt7Qp63Yuf5Ie/q4NACzWPC\nB/L6nQrut4OMxJHvhVAswJozRZrQxXvX/vUxkg+TmALj3U9ejF/5arGtjy5v+yGi\nwci7zM4r7VNFJGRkfluNRC1kJi4AY6jk6Gl4d/bX4tBXE8mEFY1rUswYtat3OMja\njVAoocClk6WcaQuK9R1uB+BPyxHLJ04RyKRuepPYRBQjgwHK5kMF3s5p07Os+2JH\n5jyJYW2NPs6pQe0k8GWpaar/yZ2eut9gsgHnu5JCWnyedo4nEx6I/G4GSaX+0SeU\n/Wb2aqq1QGfVOESml7CVcEa/buTeUwIDAQABo1AwTjAdBgNVHQ4EFgQU5i7CO32N\nspQ5AaG/aRU0LX2koYwwHwYDVR0jBBgwFoAU5i7CO32NspQ5AaG/aRU0LX2koYww\nDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAbQuCMPK8f8QuEmTpZBFv\naka9qruT/0/TrxrbxEh68N4moXSTVv4tSrDTmdkwUiiwayuGS7fvKjSV6hwGkQbV\nzGbFDdwOw1tPE2OwnA7/+RPl4KmE4iTHnnIanyg9CKmBW/tMp/vUyv5nIt7Xw5n4\ntx3C9/hme+Rlx+SVPIAwAjl4nVFNLfzyG+JDBnQWygySm88SzzK0WRgh5V+gyXCK\nucDW5rA6X9/CM3QrSY50mSM6dbyYDMtmTI4dX7E9STTBCNsNNcmgYkX0N9lm5RoF\nuBsAcPmp1SVIbXelDHJiIXxMKzwZy8riZQ8+Dw6LMs6wZX7COVvMWN4Dfcuo89av\nIQ==\n-----END CERTIFICATE-----",
		}
		_, err := store.CreateFunction(testFunc)
		if err != nil {
			t.Fatal("Failed to create function:", err)
		}
		i++
	}
	testFunctions, err := store.ListAllFunctions("")
	if err != nil {
		t.Fatal("Failed to list all functions:", err)
	}
	if len(testFunctions) != 10 {
		t.Fatal("Failed to list all function:", err)
	}

	//test deleteAllFunctions
	err = store.DeleteFunctions()
	if err != nil {
		t.Fatal("Failed to delete all functions:", err)
	}
	testFunctions, err = store.ListAllFunctions("")
	if err != nil {
		t.Fatal("Failed to list all functions:", err)
	}
	if len(testFunctions) != 0 {
		t.Fatal("Failed to delete all function:", err)
	}
}

func TestWatch(t *testing.T) {
	store, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new file store:", err)
	}
	//defer store.StopWatch()
	if err != nil {
		t.Fatal("fail to new file store:", err)
	}

	ch, err := store.Watch()
	if err != nil {
		t.Fatal("fail to watch:", err)
	}
	time.Sleep(2 * time.Second)

	//add new app
	rolePolicy1 := pms.RolePolicy{Name: "rp1", Effect: "grant", Roles: []string{"role1"}, Principals: []string{"user:Alice"}}
	rolePolicy2 := pms.RolePolicy{Name: "rp2", Effect: "grant", Roles: []string{"role2"}, Principals: []string{"user:Bill"}}
	service := pms.Service{
		Name:         "app1_new",
		Type:         pms.TypeApplication,
		RolePolicies: []*pms.RolePolicy{&rolePolicy1, &rolePolicy2},
	}
	err = store.CreateService(&service)
	if err != nil {
		t.Fatal("fail to write application:", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for e := range ch {
			i = i + 1
			t.Logf("Receive one event, type is %d\n", e.Type)
			if e.Type != pms.FULL_RELOAD {
				t.Errorf("expected event type: %d, received event type :%d\n", pms.FULL_RELOAD, e.Type)
			}
		}
		if i < 2 {
			t.Errorf("Not receive enough event")
		}
	}()

	//delete app
	store.DeleteService("app1_new")

	time.Sleep(2 * time.Second)
	store.StopWatch()
	wg.Wait()
}
