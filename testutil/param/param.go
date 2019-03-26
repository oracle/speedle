//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package param

import (
	"fmt"
)

func CREATE_SERVICE(serverName string, serverType string) string {
	return fmt.Sprintf("create service %s --service-type=%s", serverName, serverType)
}

func CREATE_SERVICE_WITH_FLAG(serverName string, serverType string, flag string) string {
	return fmt.Sprintf("create service %s --service-type=%s %s", serverName, serverType, flag)
}

func CREATE_SERVICE_WITH_JSONFILE(jsonFile string) string {
	return fmt.Sprintf("create service --json-file %s", jsonFile)
}

func CREATE_SERVICE_WITH_PDLFILE(serverName string, serverType string, pdlFile string) string {
	return fmt.Sprintf("create service %s --service-type=%s --pdl-file %s", serverName, serverType, pdlFile)
}

func GET_SERVICE(serverName string) string {
	return fmt.Sprintf("get service %s", serverName)
}

func GET_SERVICE_ALL() string {
	return fmt.Sprintf("get service --all")
}

func DELETE_SERVICE(serverName string) string {
	return fmt.Sprintf("delete service %s", serverName)
}

func DELETE_POLICY(serverName string, policyID string) string {
	return fmt.Sprintf("delete policy %s --service-name=%s", policyID, serverName)
}

func DELETE_ROLEPOLICY(serverName string, rolePolicyID string) string {
	return fmt.Sprintf("delete rolepolicy %s --service-name=%s", rolePolicyID, serverName)
}

func CREATE_POLICY(serverName string, policyName string, policyPDL string) string {
	return fmt.Sprintf("create policy %s -c \"%s\" --service-name=%s", policyName, policyPDL, serverName)
}

func CREATE_ROLEPOLICY(serverName string, rolePolicyName string, rolePolicyPDL string) string {
	return fmt.Sprintf("create rolepolicy %s -c \"%s\" --service-name=%s", rolePolicyName, rolePolicyPDL, serverName)
}

func GET_POLICY(serverName string, policyID string) string {
	return fmt.Sprintf("get policy %s --service-name=%s", policyID, serverName)
}

func GET_ROLEPOLICY(serverName string, rolePolicyID string) string {
	return fmt.Sprintf("get rolepolicy %s --service-name=%s", rolePolicyID, serverName)
}

func GET_POLICY_ALL(serverName string) string {
	return fmt.Sprintf("get policy --all --service-name=%s", serverName)
}

func GET_ROLEPOLICY_ALL(serverName string) string {
	return fmt.Sprintf("get rolepolicy --all --service-name=%s", serverName)
}

func CREATE_POLICY_WITH_FLAG(serverName string, policyName string, policyItem string, flag string) string {
	return fmt.Sprintf("create policy %s -c \"%s\"  --service-name=%s %s", policyName, policyItem, serverName, flag)
}

func CREATE_POLICY_VIA_FILE(serverName string, policyName string, fileName string) string {
	return fmt.Sprintf("create policy %s -f %s --service-name=%s", policyName, fileName, serverName)
}

func CREATE_POLICY_VIA_FILE_WITH_FLAG(serverName string, policyName string, fileName string, flag string) string {
	return fmt.Sprintf("create policy %s -f %s --service-name=%s %s", policyName, fileName, serverName, flag)
}
