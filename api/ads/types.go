//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package ads

import "github.com/oracle/speedle/api/pms"

type Principal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

func (p *Principal) String() string {
	return "{" + "\"type\": \"" + p.Type + "\", \"name\": \"" + p.Name + "\", \"idd\":\"" + p.IDD + "\"}"
}

type Subject struct {
	Principals []*Principal `json:"principals,omitempty"`
	TokenType  string       `json:"tokenType,omitempty"`
	Token      string       `json:"token,omitempty"`
	Asserted   bool         `json:"asserted,omitempty"`
}

type RequestContext struct {
	Subject     *Subject               `json:"subject,omitempty"`
	ServiceName string                 `json:"serviceName,omitempty"`
	Resource    string                 `json:"resource,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

type EvaluationResult struct {
	Allowed      bool                   `json:"allowed"`
	Reason       Reason                 `json:"reason"`
	RequestCtx   *RequestContext        `json:"requestContext,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	GrantedRoles []string               `json:"grantedRoles,omitempty"`
	RolePolicies []*EvaluatedRolePolicy `json:"rolePolicies,omitempty"`
	Policies     []*EvaluatedPolicy     `json:"policies,omitempty"`
}

type EvaluatedPolicy struct {
	Status      string              `json:"status,omitempty"`
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Effect      string              `json:"effect,omitempty"`
	Permissions []pms.Permission    `json:"permissions,omitempty"`
	Principals  [][]string          `json:"principals,omitempty"`
	Condition   *EvaluatedCondition `json:"condition,omitempty"`
}

type EvaluatedRolePolicy struct {
	Status              string              `json:"status,omitempty"`
	ID                  string              `json:"id,omitempty"`
	Name                string              `json:"name,omitempty"`
	Effect              string              `json:"effect,omitempty"`
	Roles               []string            `json:"roles,omitempty"`
	Principals          []string            `json:"principals,omitempty"`
	Resources           []string            `json:"resources,omitempty"`
	ResourceExpressions []string            `json:"resourceExpression,omitempty"`
	Condition           *EvaluatedCondition `json:"condition,omitempty"`
}

type EvaluatedCondition struct {
	ConditionExpression string `json:"conditionExpression,omitempty"`
	EvaluationResult    string `json:"evaluationResult,omitempty"`
}

const (
	Evaluation_TakeEffect      string = "takeEffect"
	Evaluation_ConditionFailed string = "conditionFailed"
	Evaluation_Ignored         string = "ignored"
)

//reason for evaluation result
type Reason int32

const (
	GRANT_POLICY_FOUND Reason = iota
	DENY_POLICY_FOUND
	SERVICE_NOT_FOUND
	NO_APPLICABLE_POLICIES
	ERROR_IN_EVALUATION
	DISCOVER_MODE
	REASON_NOT_AVAILABLE
)

const (
	BuiltIn_Role_Anonymous     = "anonymous_role"
	BuiltIn_Role_Authenticated = "authenticated_role"
	BuiltIn_Role_Everyone      = "everyone_role"

	BuiltIn_Attr_RequestUser     = "request_user"
	BuiltIn_Attr_RequestGroups   = "request_groups"
	BuiltIn_Attr_RequestResource = "request_resource"
	BuiltIn_Attr_RequestAction   = "request_action"
	BuiltIn_Attr_RequestEntity   = "request_entity"

	BuiltIn_Attr_RequestTime    = "request_time"
	BuiltIn_Attr_RequestYear    = "request_year"
	BuiltIn_Attr_RequestMonth   = "request_month"
	BuiltIn_Attr_RequestDay     = "request_day"
	BuiltIn_Attr_RequestHour    = "request_hour"
	BuiltIn_Attr_RequestWeekday = "request_weekday"
)

var reason = []string{
	"GRANT_POLICY_FOUND",
	"DENY_POLICY_FOUND",
	"SERVICE_NOT_FOUND",
	"NO_APPLICABLE_POLICIES",
	"ERROR_IN_EVALUATION",
	"DISCOVER_MODE",
	"REASON_NOT_AVAILABLE",
}

const (
	PRINCIPAL_TYPE_USER   = "user"
	PRINCIPAL_TYPE_GROUP  = "group"
	PRINCIPAL_TYPE_ROLE   = "role"
	PRINCIPAL_TYPE_ENTITY = "entity"
)

// String returns the English name of the Reason
func (m Reason) String() string { return reason[m] }
