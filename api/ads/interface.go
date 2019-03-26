//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package ads

import (
	"github.com/oracle/speedle/api/pms"
)

type PolicyEvaluator interface {
	// IsAllowed returns if the subject has been granted to a resource specified by a request context
	IsAllowed(c RequestContext) (allowed bool, reason Reason, err error)

	// GetAllGrantedRoles returns the granted app roles in an application.
	GetAllGrantedRoles(c RequestContext) ([]string, error)

	// GetAllGrantedPermissions returns the granted resources in an application.
	GetAllGrantedPermissions(cl RequestContext) ([]pms.Permission, error)

	Refresh() error

	Discover

	Diagnose
}

type Discover interface {
	// always returns true in Resource Discovery Mode
	Discover(c RequestContext) (allowed bool, reason Reason, err error)
}

type Diagnose interface {
	// returns all the policies related to a subject
	Diagnose(c RequestContext) (*EvaluationResult, error)
}
