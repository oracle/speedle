//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsimpl

import (
	"encoding/json"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/errors"
)

var (
	//The value <= 0 mean don't check the max number/size
	MaxServiceNum  = int64(-1) // Maximum number of service for a tenant
	MaxPolicyNum   = int64(-1) // Maximum number of Policy + RolePolicy per tenant
	MaxFunctionNum = int64(-1) //Maximum number of function defined by customer
	MaxPolicySize  = int64(-1) // Maximum size in bytes for a Policy or RolePolicy
)

/*
Check the following items:
	1. The maximum number of service;
	2. The maximum number of Policy + RolePolicy;
	3. The size of each Policy and RolePolicy;
*/
func CheckService(service *pms.Service, policyStore pms.PolicyStoreManager) error {
	// Check the number of the service
	srvCount, err := policyStore.GetServiceCount()
	if nil != err {
		return err
	}
	if MaxServiceNum > 0 && srvCount+1 > MaxServiceNum {
		return errors.Errorf(errors.ExceedLimit, "reached the maximum number of service, existingCount: %d", srvCount)
	}

	// Check the number of policy and rolePolicy
	existingCount, err := getPolicyAndRolePolicyCount("", policyStore)
	if nil != err {
		return err
	}
	creatingCount := 0
	if nil != service.Policies {
		creatingCount += len(service.Policies)
	}
	if nil != service.RolePolicies {
		creatingCount += len(service.RolePolicies)
	}
	if MaxPolicyNum > 0 && int64(existingCount)+int64(creatingCount) > MaxPolicyNum {
		return errors.Errorf(errors.ExceedLimit, "reached the maximum number of policy and rolePolicy, existingCount: %d, creatingCount: %d", existingCount, creatingCount)
	}

	// Check the size of each policy and RolePolicy
	for _, policy := range service.Policies {
		sizeValid, err := checkMaxSize(*policy, MaxPolicySize)
		if !sizeValid {
			return err
		}
	}
	for _, rolePolicy := range service.RolePolicies {
		sizeValid, err := checkMaxSize(*rolePolicy, MaxPolicySize)
		if !sizeValid {
			return err
		}
	}

	return nil
}

/*
Check the following items:
	1. The maximum number of Policy + RolePolicy;
	2. The size of the Policy;
    3. If the effect field of policy is empty;
*/
func CheckPolicy(serviceName string, policy *pms.Policy, policyStore pms.PolicyStoreManager) error {
	// Check global service
	if serviceName == pms.GlobalService {
		return errors.New(errors.InvalidRequest, "global policy doesn't support authorization policies")
	}

	if len(policy.Effect) <= 0 {
		return errors.New(errors.InvalidRequest, "no effect provided in policy.")
	}

	// Check the number of Policy + RolePolicy
	existingCount, err := getPolicyAndRolePolicyCount("", policyStore)
	if nil != err {
		return err
	}
	if MaxPolicyNum > 0 && existingCount+1 > MaxPolicyNum {
		return errors.Errorf(errors.ExceedLimit, "reached the maximum number of policy and rolePolicy: %d", existingCount)
	}

	// Check the size of the Policy
	sizeValid, err := checkMaxSize(*policy, MaxPolicySize)
	if !sizeValid {
		return err
	}

	return nil
}

/*
Check the following items:
	1. The maximum number of Policy + RolePolicy;
	2. The size of the RolePolicy;
    3. If the effect field of RolePolicy is empty;
*/
func CheckRolePolicy(serviceName string, rolePolicy *pms.RolePolicy, policyStore pms.PolicyStoreManager) error {
	if len(rolePolicy.Effect) <= 0 {
		return errors.New(errors.InvalidRequest, "no effect provided in role policy.")
	}

	// Check the number of Policy + RolePolicy
	existingCount, err := getPolicyAndRolePolicyCount("", policyStore)
	if nil != err {
		return err
	}
	if MaxPolicyNum > 0 && existingCount+1 > MaxPolicyNum {
		return errors.Errorf(errors.ExceedLimit, "reached the maximum number of policy and rolePolicy: %d", existingCount)
	}

	// Check the size of the RolePolicy
	sizeValid, err := checkMaxSize(*rolePolicy, MaxPolicySize)
	if !sizeValid {
		return err
	}

	return nil
}

// get the existing number of policy + rolePolicy
func getPolicyAndRolePolicyCount(serviceName string, policyStore pms.PolicyStoreManager) (int64, error) {
	policyCount, err := policyStore.GetPolicyCount(serviceName)
	if nil != err {
		return 0, err
	}

	rolePolicyCount, err := policyStore.GetRolePolicyCount(serviceName)
	if nil != err {
		return 0, err
	}

	return (policyCount + rolePolicyCount), nil
}

// check the size of policy or rolePolicy
func checkMaxSize(val interface{}, maxSize int64) (bool, error) {
	value, err := json.Marshal(val)
	if err != nil {
		return false, errors.Wrap(err, errors.SerializationError, "failed to marshal object")
	}
	if maxSize > 0 && int64(len(value)) > maxSize {
		return false, errors.Errorf(errors.ExceedLimit, "the actual size: %d exceeded the maximum size: %d", len(value), maxSize)
	}

	return true, nil
}

/*
Check the following items:
	1. The maximum number of function;
*/
func CheckFunction(function *pms.Function, policyStore pms.PolicyStoreManager) error {
	// Check the number of Policy + RolePolicy
	existingCount, err := policyStore.GetFunctionCount()
	if nil != err {
		return err
	}
	if MaxFunctionNum > 0 && existingCount >= MaxFunctionNum {
		return errors.Errorf(errors.ExceedLimit, "reached the maximum number of function: %d", existingCount)
	}
	return nil
}
