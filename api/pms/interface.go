//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pms

type StoreManager interface {
	ReadPolicyStore() (*PolicyStore, error)
	WritePolicyStore(*PolicyStore) error
	Type() string
}

type FunctionManager interface {
	CreateFunction(function *Function) (*Function, error)
	DeleteFunction(funcName string) error
	DeleteFunctions() error
	GetFunction(funcName string) (*Function, error)
	ListAllFunctions(filter string) ([]*Function, error)
	GetFunctionCount() (int64, error)
}

type ServiceManager interface {
	CreateService(service *Service) error
	DeleteService(serviceName string) error
	DeleteServices() error
	GetService(serviceName string) (*Service, error)
	ListAllServices() ([]*Service, error)
	GetServiceCount() (int64, error)
	GetServiceNames() ([]string, error)
	GetPolicyAndRolePolicyCounts() (map[string]*PolicyAndRolePolicyCount, error)
}

type PolicyManager interface {
	CreatePolicy(serviceName string, policy *Policy) (*Policy, error)
	DeletePolicy(serviceName string, id string) error
	DeletePolicies(serviceName string) error
	GetPolicy(serviceName string, id string) (*Policy, error)
	ListAllPolicies(serviceName string, filter string) ([]*Policy, error)
	GetPolicyCount(serviceName string) (int64, error)
}

type RolePolicyManager interface {
	CreateRolePolicy(serviceName string, policy *RolePolicy) (*RolePolicy, error)
	DeleteRolePolicy(serviceName string, id string) error
	DeleteRolePolicies(serviceName string) error
	GetRolePolicy(serviceName string, id string) (*RolePolicy, error)
	ListAllRolePolicies(serviceName string, filter string) ([]*RolePolicy, error)
	GetRolePolicyCount(serviceName string) (int64, error)
}

type PolicyStoreWatcher interface {
	Watch() (StorageChangeChannel, error)
	StopWatch()
}

type PolicyStoreManager interface {
	ServiceManager
	StoreManager
	PolicyManager
	RolePolicyManager
	FunctionManager
	PolicyStoreWatcher
}

type PolicyStoreManagerADS interface {
	Type() string
	ReadPolicyStore() (*PolicyStore, error)
	GetService(serviceName string) (*Service, error)
	GetPolicy(serviceName string, id string) (*Policy, error)
	GetRolePolicy(serviceName string, id string) (*RolePolicy, error)
	GetFunction(funcName string) (*Function, error)
	PolicyStoreWatcher
}
