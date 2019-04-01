//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package adsrest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/cfg"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/eval"
	"github.com/oracle/speedle/pkg/httputils"
	"github.com/oracle/speedle/pkg/logging"

	"github.com/oracle/speedle/pkg/svcs"
	log "github.com/sirupsen/logrus"
)

type JsonAttribute struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type JsonPrincipal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

type JsonSubject struct {
	Principals []*JsonPrincipal `json:"principals,omitempty"`
	TokenType  string           `json:"tokenType"`
	Token      string           `json:"token"`
}

type JsonContext struct {
	Subject     *JsonSubject     `json:"subject"`
	ServiceName string           `json:"serviceName"`
	Resource    string           `json:"resource"`
	Action      string           `json:"action"`
	Attributes  []*JsonAttribute `json:"attributes"`
}

type RESTService struct {
	Evaluator eval.InternalEvaluator
}

type IsAllowedResponse struct {
	Allowed      bool   `json:"allowed"`
	Reason       int32  `json:"reason"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type AuditEvaluationResult struct {
	Allowed string `json:"allowed"`
	Reason  string `json:"reason"`
}

type PermissionResponse struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

type PolicyResponse struct {
	Status      string             `json:"status,omitempty"`
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Effect      string             `json:"effect,omitempty"`
	Permissions []Permission       `json:"permissions,omitempty"`
	Principals  [][]string         `json:"principals,omitempty"`
	Condition   EvaluatedCondition `json:"condition,omitempty"`
}

type RolePolicyResponse struct {
	Status              string             `json:"status,omitempty"`
	ID                  string             `json:"id,omitempty"`
	Name                string             `json:"name,omitempty"`
	Effect              string             `json:"effect,omitempty"`
	Roles               []string           `json:"roles,omitempty"`
	Principals          []string           `json:"principals,omitempty"`
	Resources           []string           `json:"resources,omitempty"`
	ResourceExpressions []string           `json:"resourceExpressions,omitempty"`
	Condition           EvaluatedCondition `json:"condition,omitempty"`
}

type Permission struct {
	Resource           string   `json:"resource,omitempty"`
	ResourceExpression string   `json:"resourceExpression,omitempty"`
	Actions            []string `json:"actions,omitempty"`
}

type EvaluatedCondition struct {
	ConditionExpression string `json:"conditionExpression,omitempty"`
	EvaluationResult    string `json:"evaluationResult,omitempty"`
}

// Should we add Both of ReasonCode and ReasonMessage
type EvaluationDebugResponse struct {
	Allowed        bool                   `json:"allowed"`
	Reason         string                 `json:"reason"`
	RequestContext JsonContext            `json:"requestContext,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
	GrantedRoles   []string               `json:"grantedRoles,omitempty"`
	RolePolicies   []RolePolicyResponse   `json:"rolePolicies,omitempty"`
	Policies       []PolicyResponse       `json:"policies,omitempty"`
}

func NewRESTService(conf *cfg.Config) (*RESTService, error) {
	Evaluator, err := eval.NewFromConfig(conf)
	if err != nil {
		return nil, err
	}
	return &RESTService{
		Evaluator: Evaluator,
	}, nil
}

func NewRESTServiceWithEvaluator(evaluator eval.InternalEvaluator) (*RESTService, error) {
	return &RESTService{
		Evaluator: evaluator,
	}, nil
}

func DecodeJSONContext(r *http.Request) (*JsonContext, error) {
	decoder := json.NewDecoder(r.Body)
	var request JsonContext
	if err := decoder.Decode(&request); err != nil {
		return nil, errors.Wrap(err, errors.InvalidRequest, "unable to decode request")
	}
	return &request, nil
}

func DuplicateAttributeMap(attrs map[string]interface{}) map[string]interface{} {
	if attrs == nil {
		return nil
	}
	ret := make(map[string]interface{})
	for key, value := range attrs {
		ret[key] = value
	}
	return ret
}

func VerifyAttributeName(attrName string) error {
	// Currently don't verify attribute name
	return nil
}

// Key is data type in json
// Value is the data type in go
var dataTypeMap = map[string]string{
	"string":   "string",
	"numeric":  "float64",
	"bool":     "bool",
	"datetime": "string",
}

var supportDateTimeLayout = []string{
	time.RFC3339Nano,
	time.RubyDate,
	time.UnixDate,
}

func ParseDateTime(value string) (*time.Time, error) {
	for _, layout := range supportDateTimeLayout {
		ret, err := time.Parse(layout, value)
		if err == nil {
			return &ret, nil
		}
	}

	return nil, errors.Errorf(errors.InvalidRequest, "value %q is not a supported date time", value)
}

func ConvSingleValue(dataType string, value interface{}) (interface{}, error) {
	if len(dataType) == 0 {
		// If data type is empty, return value directly
		return value, nil
	}

	valueType, ok := dataTypeMap[dataType]
	if !ok {
		// Data type is not match
		return nil, errors.Errorf(errors.InvalidRequest, "inputted data type %s is not supported", dataType)
	}

	if valueType != reflect.TypeOf(value).String() {
		return nil, errors.Errorf(errors.InvalidRequest, "value data type %T is not equals to inputted data type %s", value, dataType)
	}

	switch dataType {
	case "datetime":
		strValue, _ := value.(string)
		retTime, err := ParseDateTime(strValue)
		if err != nil {
			return nil, err
		}
		return float64(retTime.Unix()), nil
	}
	return value, nil
}

func ConvMultipleValues(dataType string, values interface{}) (interface{}, error) {
	v := reflect.ValueOf(values)
	ret := []interface{}{}
	var prevType string
	for i := 0; i < v.Len(); i = i + 1 {
		vi := v.Index(i)
		item, err := ConvSingleValue(dataType, vi.Interface())
		if err != nil {
			return nil, err
		}
		if i == 0 {
			prevType = reflect.TypeOf(item).String()
		} else if prevType != reflect.TypeOf(item).String() {
			return nil, errors.New(errors.InvalidRequest, "types of all items in a array should be same")
		}
		ret = append(ret, item)
	}
	return ret, nil
}

func ConvValue(dataType string, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, errors.New(errors.InvalidRequest, "null value is not allowed")
	}
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		return ConvMultipleValues(dataType, value)
	}
	return ConvSingleValue(dataType, value)
}

func DumpRequestAttributes(attrs []*JsonAttribute) (map[string]interface{}, error) {
	// No attributes found
	if attrs == nil || len(attrs) == 0 {
		return nil, nil
	}

	attrMap := make(map[string]interface{})
	for _, attr := range attrs {
		if err := VerifyAttributeName(attr.Name); err != nil {
			return nil, err
		}
		value, err := ConvValue(attr.Type, attr.Value)
		if err != nil {
			return nil, err
		}
		attrMap[attr.Name] = value
	}
	return attrMap, nil
}

func DumpPrincipals(principals []*JsonPrincipal) []*adsapi.Principal {
	if principals == nil {
		return nil
	}
	ret := []*adsapi.Principal{}
	for _, princ := range principals {
		ret = append(ret, &adsapi.Principal{
			Type: princ.Type,
			Name: princ.Name,
			IDD:  princ.IDD,
		})
	}
	return ret
}

func ConvertJSONRequestToContext(ctxContext *JsonContext) (*adsapi.RequestContext, error) {
	subject := adsapi.Subject{}
	if ctxContext.Subject != nil {
		apiPrincipals := DumpPrincipals(ctxContext.Subject.Principals)
		subject = adsapi.Subject{
			Principals: apiPrincipals,
			TokenType:  ctxContext.Subject.TokenType,
			Token:      ctxContext.Subject.Token,
		}
	}

	contextAttr, err := DumpRequestAttributes(ctxContext.Attributes)
	if err != nil {
		return nil, err
	}

	context := adsapi.RequestContext{
		Subject:     &subject,
		ServiceName: ctxContext.ServiceName,
		Resource:    ctxContext.Resource,
		Action:      ctxContext.Action,
		Attributes:  contextAttr,
	}
	return &context, nil
}

func constructEvaluationResultForAudit(allowed bool, reason adsapi.Reason) *AuditEvaluationResult {
	evaResult := "denied"
	if allowed {
		evaResult = "allowed"
	}

	auditResult := AuditEvaluationResult{
		Allowed: evaResult,
		Reason:  reason.String(),
	}

	return &auditResult
}

func (e *RESTService) IsAllowed(w http.ResponseWriter, r *http.Request) {
	jsonRequest, err := DecodeJSONContext(r)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	context, err := ConvertJSONRequestToContext(jsonRequest)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	result, reason, err := e.Evaluator.IsAllowed(*context)
	response := IsAllowedResponse{
		Allowed: result,
		Reason:  int32(reason),
	}
	// Audit log
	responseForAudit := constructEvaluationResultForAudit(result, reason)

	//Token assertion is done in e.Evaluator.IsAllowed(). Now context has been populated with subject info
	if len(context.Subject.Principals) > 0 {
		for _, principal := range context.Subject.Principals {
			if principal.Type == adsapi.PRINCIPAL_TYPE_USER {
				w.Header().Add(svcs.PrincipalsHeader, principal.Name)
				break
			}
		}
	}

	if err != nil {
		response.ErrorMessage = err.Error()
		logging.WriteFailedAuditLog("IsAllowed", log.Fields{"requestContext": context, "evaluationResult": responseForAudit}, response.ErrorMessage)
	} else {
		logging.WriteSucceededAuditLog("IsAllowed", log.Fields{"requestContext": context}, log.Fields{"evaluationResult": responseForAudit})
	}

	httputils.SendOKResponse(w, &response)
}

func (e *RESTService) GetAllGrantedRoles(w http.ResponseWriter, r *http.Request) {
	jsonRequest, err := DecodeJSONContext(r)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	context, err := ConvertJSONRequestToContext(jsonRequest)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	roles, err := e.Evaluator.GetAllGrantedRoles(*context)
	if err != nil {
		httputils.HandleError(w, err)
		// Audit log
		logging.WriteFailedAuditLog("GetAllGrantedRoles", log.Fields{"requestContext": context}, err.Error())
		return
	}

	// Audit log
	logging.WriteSucceededAuditLog("GetAllGrantedRoles", log.Fields{"requestContext": context}, log.Fields{"roles": roles})

	if len(roles) == 0 {
		httputils.SendEmptyListResponse(w)
		return
	}
	httputils.SendOKResponse(w, roles)
}

func (e *RESTService) GetAllGrantedPermissions(w http.ResponseWriter, r *http.Request) {
	jsonRequest, err := DecodeJSONContext(r)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	context, err := ConvertJSONRequestToContext(jsonRequest)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	permissions, err := e.Evaluator.GetAllGrantedPermissions(*context)
	if err != nil {
		httputils.HandleError(w, err)
		// Audit log
		logging.WriteFailedAuditLog("GetAllGrantedPermissions", log.Fields{"requestContext": context}, err.Error())
		return
	}

	var retPermissions []PermissionResponse
	for _, permission := range permissions {
		retPermissions = append(retPermissions, PermissionResponse{
			Resource: permission.Resource,
			Actions:  permission.Actions,
		})
	}

	// Audit log
	logging.WriteSucceededAuditLog("GetAllGrantedPermissions", log.Fields{"requestContext": context}, log.Fields{"permissions": retPermissions})

	if len(retPermissions) == 0 {
		httputils.SendEmptyListResponse(w)
		return
	}
	httputils.SendOKResponse(w, retPermissions)
}

func ConvertAPIPolicy2PolicyResponse(apiPolicy *adsapi.EvaluatedPolicy, policyResp *PolicyResponse) {
	if apiPolicy == nil || policyResp == nil {
		// It shouldn't happen
		return
	}

	var retPermission []Permission

	for _, permission := range apiPolicy.Permissions {
		retPermission = append(retPermission, Permission{
			Resource:           permission.Resource,
			Actions:            permission.Actions,
			ResourceExpression: permission.ResourceExpression,
		})
	}

	policyResp.Status = apiPolicy.Status
	policyResp.ID = apiPolicy.ID
	policyResp.Name = apiPolicy.Name
	policyResp.Effect = apiPolicy.Effect
	policyResp.Permissions = retPermission
	policyResp.Principals = apiPolicy.Principals

	if apiPolicy.Condition != nil {
		policyResp.Condition = EvaluatedCondition{
			ConditionExpression: apiPolicy.Condition.ConditionExpression,
			EvaluationResult:    apiPolicy.Condition.EvaluationResult,
		}
	}

}

func ConvertAPIRolePolicy2RolePolicyResponse(apiRolePolicy *adsapi.EvaluatedRolePolicy, rolePolicyResp *RolePolicyResponse) {
	if apiRolePolicy == nil || rolePolicyResp == nil {
		// It shouldn't happen
		return
	}

	rolePolicyResp.Status = apiRolePolicy.Status
	rolePolicyResp.ID = apiRolePolicy.ID
	rolePolicyResp.Name = apiRolePolicy.Name
	rolePolicyResp.Effect = apiRolePolicy.Effect
	rolePolicyResp.Roles = apiRolePolicy.Roles
	rolePolicyResp.Principals = apiRolePolicy.Principals
	rolePolicyResp.Resources = apiRolePolicy.Resources
	rolePolicyResp.ResourceExpressions = apiRolePolicy.ResourceExpressions

	if apiRolePolicy.Condition != nil {
		rolePolicyResp.Condition = EvaluatedCondition{
			ConditionExpression: apiRolePolicy.Condition.ConditionExpression,
			EvaluationResult:    apiRolePolicy.Condition.EvaluationResult,
		}
	}
}

func (e *RESTService) Diagnose(w http.ResponseWriter, r *http.Request) {
	jsonRequest, err := DecodeJSONContext(r)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	context, err := ConvertJSONRequestToContext(jsonRequest)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	evaResult, err := e.Evaluator.Diagnose(*context)
	if err != nil {
		httputils.HandleError(w, err)
		// Audit log
		logging.WriteSimpleFailedAuditLog("Diagnose", context, err.Error())
		return
	}

	// Convert all the returned policies
	var retPolicies []PolicyResponse
	for _, policy := range evaResult.Policies {
		var policyResp PolicyResponse
		ConvertAPIPolicy2PolicyResponse(policy, &policyResp)
		retPolicies = append(retPolicies, policyResp)
	}

	// Convert all the returned role policies
	var retRolePolicies []RolePolicyResponse
	for _, rolePolicy := range evaResult.RolePolicies {
		var rolePolicyResp RolePolicyResponse
		ConvertAPIRolePolicy2RolePolicyResponse(rolePolicy, &rolePolicyResp)
		retRolePolicies = append(retRolePolicies, rolePolicyResp)
	}

	// Construct & return the response
	response := EvaluationDebugResponse{
		Allowed:        evaResult.Allowed,
		Reason:         evaResult.Reason.String(),
		RequestContext: *jsonRequest,
		Attributes:     evaResult.Attributes,
		GrantedRoles:   evaResult.GrantedRoles,
		RolePolicies:   retRolePolicies,
		Policies:       retPolicies,
	}

	// Audit log
	logging.WriteSimpleSucceededAuditLog("Diagnose", context, &response)

	httputils.SendOKResponse(w, &response)
}
