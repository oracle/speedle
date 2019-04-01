//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package etcd

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/cfg"
	"github.com/oracle/speedle/pkg/store"
)

var storeConfig *cfg.StoreConfig

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	etcd, etcdDir, err := StartEmbeddedEtcd("")
	defer CleanEmbeddedEtcd(etcd, etcdDir)
	if err != nil {
		log.Fatal("Fail to start embed etcd")
		os.Exit(1)
	}

	storeConfig, err = cfg.ReadStoreConfig("./etcdStoreConfig.json")
	if err != nil {
		log.Fatal("fail to read config file", err)
	}
	return m.Run()
}

func TestWriteReadPolicyStore(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()

	if psOrigin, err := store.ReadPolicyStore(); err != nil {
		t.Fatal("fail to read etcd3 store:", err)
	} else {
		t.Log("existing number of apps:", len(psOrigin.Services))
	}

	var ps pms.PolicyStore
	for i := 0; i < 10; i++ {
		service := pms.Service{Name: fmt.Sprintf("app%d", i), Type: pms.TypeApplication}
		ps.Services = append(ps.Services, &service)
	}
	err = store.WritePolicyStore(&ps)
	if err != nil {
		t.Fatal("fail to write policy store:", err)
	}
	var psr *pms.PolicyStore
	psr, err = store.ReadPolicyStore()
	if err != nil {
		t.Fatal("fail to read policy store:", err)
	}
	if 10 != len(psr.Services) {
		t.Error("should have 10 applications in the store")
	}
	for _, app := range psr.Services {
		t.Log(app.Name, " ")
	}
}

func TestWriteReadDeleteService(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()
	//clean the service firstly
	err = store.DeleteService("service1")
	t.Log("deleteing service1, err:", err)

	app := pms.Service{Name: "service1", Type: pms.TypeApplication}
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
	i = 0
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
	err = store.CreateService(&app)
	if err != nil {
		t.Log("fail to create application:", err)
		t.FailNow()
	}
	appr, errr := store.GetService("service1")
	if errr != nil {
		t.Log("fail to get application:", err)
		t.FailNow()
	}
	if "service1" != appr.Name {
		t.Log("app name should be service1")
		t.FailNow()
	}
	if pms.TypeApplication != appr.Type {
		t.Log("app type should be ", pms.TypeApplication)
		t.FailNow()
	}
	if num != len(appr.RolePolicies) {
		t.Logf("role policy number should be %d, but %d.", num, len(appr.RolePolicies))
		t.FailNow()
	}
	if num != len(appr.Policies) {
		t.Log("policy number should be ", num)
		t.FailNow()
	}
	err = store.DeleteService("service1")
	if err != nil {
		t.Log("fail to delete application:", err)
		t.FailNow()
	}
	appr, err = store.GetService("service1")
	t.Log("get non exist service:", err)
	if err == nil {
		t.Log("should fail as app is already deleted")
		t.FailNow()
	}
	err = store.DeleteService("nonexist-service")
	t.Log("delete non exist service:", err)
	if err == nil {
		t.Log("should fail as the service does not exist")
		t.FailNow()
	}
}

func TestEtcdStore_GetPolicyByName(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()
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

func TestEtcdStore_GetRolePolicyByName(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()
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

func TestManagePolicies(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()
	//clean the service firstly
	store.DeleteService("service1")
	app := pms.Service{Name: "service1", Type: pms.TypeApplication}
	err = store.CreateService(&app)
	if err != nil {
		t.Fatal("fail to create application:", err)
	}
	var policy pms.Policy
	policy.Name = fmt.Sprintf("policy1")
	policy.Effect = "grant"
	policy.Permissions = []*pms.Permission{
		{
			Resource: "/node1",
			Actions:  []string{"get", "create", "delete"},
		},
	}
	policy.Principals = [][]string{{"user:Alice"}}
	policyR, err := store.CreatePolicy("service1", &policy)
	if err != nil {
		t.Fatal("fail to create policy:", err)
	}
	policyR1, err := store.GetPolicy("service1", policyR.ID)
	t.Log(policyR1)
	if err != nil {
		t.Fatal("fail to get policy:", err)
	}

	policies, err := store.ListAllPolicies("service1", "")
	if err != nil {
		t.Fatal("fail to list policies:", err)
	}
	if len(policies) != 1 {
		t.Fatal("should have 1 policy")
	}

	_, err = store.GetPolicy("service1", "nonexistID")
	t.Log(err)
	if err == nil {
		t.Fatal("should fail to get policy")
	}

	err = store.DeletePolicy("service1", "nonexistID")
	t.Log(err)
	if err == nil {
		t.Fatal("should fail to delete policy")
	}

	err = store.DeletePolicy("service1", policyR.ID)
	if err != nil {
		t.Fatal("fail to delete policy:", err)
	}
}

func TestManageRolePolicies(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()

	//clean the service firstly
	store.DeleteService("service1")
	app := pms.Service{Name: "service1", Type: pms.TypeApplication}
	err = store.CreateService(&app)
	if err != nil {
		t.Fatal("fail to create application:", err)
	}
	var rolePolicy pms.RolePolicy
	rolePolicy.Name = "rp1"
	rolePolicy.Effect = "grant"
	rolePolicy.Roles = []string{"role1"}
	rolePolicy.Principals = []string{"user:Alice"}

	policyR, err := store.CreateRolePolicy("service1", &rolePolicy)
	if err != nil {
		t.Fatal("fail to create role policy:", err)
	}
	policyR1, err := store.GetRolePolicy("service1", policyR.ID)
	t.Log(policyR1)
	if err != nil {
		t.Fatal("fail to get role policy:", err)
	}

	rolePolicies, err := store.ListAllRolePolicies("service1", "")
	if err != nil {
		t.Fatal("fail to list role policies:", err)
	}
	if len(rolePolicies) != 1 {
		t.Fatal("should have 1 role policy")
	}

	_, err = store.GetRolePolicy("service1", "nonexistID")
	t.Log(err)
	if err == nil {
		t.Fatal("should fail to get role policy")
	}

	err = store.DeleteRolePolicy("service1", "nonexistID")
	t.Log(err)
	if err == nil {
		t.Fatal("should fail to delete role policy")
	}

	err = store.DeleteRolePolicy("service1", policyR.ID)
	if err != nil {
		t.Fatal("fail to delete role policy:", err)
	}
}

func TestCheckItemsCount(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
	}
	defer store.(*Store).destroy()

	// clean the services
	store.DeleteServices()

	// Create service1
	app1 := pms.Service{Name: "service1", Type: pms.TypeApplication}
	err = store.CreateService(&app1)
	if err != nil {
		t.Fatal("fail to create service:", err)
	}
	// Check service count
	serviceCount, err := store.GetServiceCount()
	if err != nil {
		t.Fatal("Failed to get service count:", err)
	}
	if serviceCount != 1 {
		t.Fatalf("Service count doesn't match, expected: 1, actual: %d", serviceCount)
	}

	// Create policies
	policies := []pms.Policy{
		{Name: "p01", Effect: "grant", Principals: [][]string{{"user:user1"}}},
		{Name: "p02", Effect: "grant", Principals: [][]string{{"user:user2"}}},
		{Name: "p03", Effect: "grant", Principals: [][]string{{"user:user3"}}},
	}
	for _, policy := range policies {
		_, err := store.CreatePolicy("service1", &policy)
		if err != nil {
			t.Fatal("fail to create policy:", err)
		}
	}
	// Check policy count
	policyCount, err := store.GetPolicyCount("service1")
	if err != nil {
		t.Fatal("Failed to get the policy count: ", err)
	}
	if policyCount != int64(len(policies)) {
		t.Fatalf("Policy count doesn't match, expected:%d, actual:%d", len(policies), policyCount)
	}

	// Create Role Policies
	rolePolicies := []pms.RolePolicy{
		{Name: "p01", Effect: "grant", Principals: []string{"user:user1"}, Roles: []string{"role1"}},
		{Name: "p02", Effect: "grant", Principals: []string{"user:user2"}, Roles: []string{"role2"}},
	}
	for _, rolePolicy := range rolePolicies {
		_, err := store.CreateRolePolicy("service1", &rolePolicy)
		if err != nil {
			t.Fatal("Failed to get role policy count:", err)
		}
	}
	// Check role Policy count
	rolePolicyCount, err := store.GetRolePolicyCount("service1")
	if err != nil {
		t.Fatal("Failed to get the role policy count")
	}
	if rolePolicyCount != int64(len(rolePolicies)) {
		t.Fatalf("RolePolicy count doesn't match, expected:%d, actual:%d", len(rolePolicies), rolePolicyCount)
	}

	// Create service2
	app2 := pms.Service{Name: "service2", Type: pms.TypeApplication}
	err = store.CreateService(&app2)
	if err != nil {
		t.Fatal("fail to create service:", err)
	}
	// Check service count
	serviceCount, err = store.GetServiceCount()
	if err != nil {
		t.Fatal("Failed to get service count:", err)
	}
	if serviceCount != 2 {
		t.Fatalf("Service count doesn't match, expected: 2, actual: %d", serviceCount)
	}

	// Create policies in service2
	for _, policy := range policies {
		_, err := store.CreatePolicy("service2", &policy)
		if err != nil {
			t.Fatal("fail to create policy:", err)
		}
	}
	// Check policy count in service2
	policyCount, err = store.GetPolicyCount("service2")
	if err != nil {
		t.Fatal("Failed to get the policy count: ", err)
	}
	if policyCount != int64(len(policies)) {
		t.Fatalf("Policy count doesn't match, expected:%d, actual:%d", len(policies), policyCount)
	}
	// Check policy count in both service1 and service2
	policyCount, err = store.GetPolicyCount("")
	if err != nil {
		t.Fatal("Failed to get the policy count: ", err)
	}
	if policyCount != int64(len(policies)*2) {
		t.Fatalf("Policy count doesn't match, expected:%d, actual:%d", len(policies)*2, policyCount)
	}

	// Create rolePolicy in service2
	for _, rolePolicy := range rolePolicies {
		_, err := store.CreateRolePolicy("service2", &rolePolicy)
		if err != nil {
			t.Fatal("Failed to get role policy count:", err)
		}
	}
	// Check role Policy count in service2
	rolePolicyCount, err = store.GetRolePolicyCount("service2")
	if err != nil {
		t.Fatal("Failed to get the role policy count")
	}
	if rolePolicyCount != int64(len(rolePolicies)) {
		t.Fatalf("RolePolicy count doesn't match, expected:%d, actual:%d", len(rolePolicies), rolePolicyCount)
	}
	// Check role Policy count in both service1 and service2
	rolePolicyCount, err = store.GetRolePolicyCount("")
	if err != nil {
		t.Fatal("Failed to get the role policy count")
	}
	if rolePolicyCount != int64(len(rolePolicies)*2) {
		t.Fatalf("RolePolicy count doesn't match, expected:%d, actual:%d", len(rolePolicies)*2, rolePolicyCount)
	}
}

func TestWatch(t *testing.T) {
	store, err := store.NewStore(storeConfig.StoreType, storeConfig.StoreProps)
	defer store.StopWatch()
	defer store.(*Store).destroy()
	if err != nil {
		t.Fatal("fail to new etcd3 store:", err)
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

	select {
	case <-time.After(5 * time.Second):
		t.Errorf("fail to receive policy update event")
	case e := <-ch:
		if e.Type != pms.SERVICE_ADD {
			t.Errorf("expected event type: %d, received event type :%d\n", pms.SERVICE_ADD, e.Type)
		}
	}

	//delete app
	store.DeleteService("app1_new")
	select {
	case <-time.After(5 * time.Second):
		t.Errorf("fail to receive policy update event")
	case e := <-ch:
		if e.Type != pms.SERVICE_DELETE {
			t.Errorf("expected event type: %d, received event type :%d\n", pms.SERVICE_DELETE, e.Type)
		}
	}

}
