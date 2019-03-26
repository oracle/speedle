//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package msg

import (
	"fmt"
)

func OUTPUT_SERVICE_CREATED() string {
	return "service created"
}

func OUTPUT_SERVICE_DELETED(sName string) string {
	return fmt.Sprintf("service %s deleted", sName)
}

func OUTPUT_POLICY_CREATED() string {
	return "policy created"
}

func OUTPUT_POLICY_DELETED(policyName string) string {
	return fmt.Sprintf("policy %s deleted", policyName)
}

func OUTPUT_ROLEPOLICY_CREATED() string {
	return "rolepolicy created"
}

func OUTPUT_ROLEPOLICY_DELETED(rpolicyName string) string {
	return fmt.Sprintf("rolepolicy %s deleted", rpolicyName)
}

func OUTPUT_POLICY_NOTFOUND(serverName string, policyID string) string {
	return fmt.Sprintf("service %s policy %s not found", serverName, policyID)
}

func OUTPUT_ROLEPOLICY_NOTFOUND(serverName string, rolePolicyID string) string {
	return fmt.Sprintf("service %s role-policy %s not found", serverName, rolePolicyID)
}
