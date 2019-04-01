//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsrest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/httputils"
	"github.com/oracle/speedle/pkg/logging"
	"github.com/oracle/speedle/pkg/svcs/pmsimpl"

	"github.com/gorilla/mux"

	"github.com/oracle/speedle/api/pms"

	"time"

	"github.com/oracle/speedle/pkg/svcs"
	log "github.com/sirupsen/logrus"
)

type RESTService struct {
	PolicyStore pms.PolicyStoreManager
}

type serviceRequestBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewRestService(s pms.PolicyStoreManager) (*RESTService, error) {
	return &RESTService{PolicyStore: s}, nil
}

// returns:
//     1. ServiceName
//     2. policy/role-policy ID
func ParseRequestURI(r *http.Request) (string, string) {
	segs := strings.Split(r.URL.Path, "/")
	segLength := len(segs)
	if segLength > 4 {
		if segLength > 6 {
			return segs[4], segs[6]
		}
		return segs[4], ""
	}
	return "", ""
}

// ParseForFilters parse query filter from request
// returns
// value of query parameter - filter
func ParseForFilters(r *http.Request) string {
	filterStr := r.URL.Query().Get("filter")
	if len(filterStr) <= 0 {
		return ""
	}
	return filterStr
}

func decodeServiceRequest(r *http.Request) (*serviceRequestBody, error) {
	decoder := json.NewDecoder(r.Body)
	var request serviceRequestBody
	err := decoder.Decode(&request)
	if err != nil {
		return nil, errors.Wrap(err, errors.InvalidRequest, "failed to decode request body")
	}

	// TODO: Verify if context is good
	return &request, nil
}

func decodeRequestBody(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields()  //temporarily remove this go 1.10 feature since 1.10 cannot debug evaluator test due to golang issue #23733
	err := decoder.Decode(obj)
	if err != nil {
		return errors.Wrap(err, errors.InvalidRequest, "failed to decode request body")
	}

	// TODO: Verify if context is good
	return nil
}

func getCreateMetaData(r *http.Request) map[string]string {
	var createMetaData = make(map[string]string)
	creator := r.Header.Get(svcs.PrincipalsHeader)
	if creator != "" { //set creteby meta data only when asserter returned creator info
		createMetaData["createby"] = creator
	}
	createMetaData["createtime"] = time.Unix(time.Now().Unix(), 0).Format(time.RFC3339)
	return createMetaData
}

// Service management
func (mgr *RESTService) CreateService(w http.ResponseWriter, r *http.Request) {
	var service pms.Service
	err := decodeRequestBody(r, &service)
	if err != nil {
		httputils.HandleError(w, err)
		return
	}

	err = pmsimpl.CheckService(&service, mgr.PolicyStore)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreateService", &service, err.Error())
		return
	}

	if _, err := mgr.PolicyStore.GetService(service.Name); err == nil {
		// servcie already exists.
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Service already exists.",
		})
		logging.WriteSimpleFailedAuditLog("CreateService", &service, "Service already exists")
		return
	}

	//set createby and createtime
	var metaData = getCreateMetaData(r)
	service.Metadata = metaData
	for _, policy := range service.Policies {
		policy.Metadata = metaData
	}
	for _, rolepolicy := range service.RolePolicies {
		rolepolicy.Metadata = metaData
	}

	if err := mgr.PolicyStore.CreateService(&service); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreateService", &service, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("CreateService", &service, nil)
	httputils.SendCreatedResponse(w, &service)
}

func (mgr *RESTService) DeleteService(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}

	if err := mgr.PolicyStore.DeleteService(serviceName); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("DeleteService", serviceName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("DeleteService", serviceName, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) DeleteServices(w http.ResponseWriter, r *http.Request) {
	if err := mgr.PolicyStore.DeleteServices(); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("DeleteServices", nil, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("DeleteServices", nil, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) GetService(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}

	service, err := mgr.PolicyStore.GetService(serviceName)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("GetService", serviceName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("GetService", serviceName, nil)
	httputils.SendOKResponse(w, &service)
}

func (mgr *RESTService) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := mgr.PolicyStore.ListAllServices()
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("ListServices", nil, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("ListServices", nil, len(services))

	if len(services) == 0 {
		httputils.SendEmptyListResponse(w)
		return
	}
	httputils.SendOKResponse(w, &services)
}

func (mgr *RESTService) ListPolicyAndRolePolicyCounts(w http.ResponseWriter, r *http.Request) {
	countMap, err := mgr.PolicyStore.GetPolicyAndRolePolicyCounts()
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("ListPolicyCounts", nil, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("ListPolicyCounts", nil, countMap)
	httputils.SendOKResponse(w, countMap)
}

// Policy management
func (mgr *RESTService) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}
	var policy pms.Policy
	if err := decodeRequestBody(r, &policy); err != nil {
		httputils.HandleError(w, err)

		return
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName": serviceName,
		"policy":      &policy,
	}

	err := pmsimpl.CheckPolicy(serviceName, &policy, mgr.PolicyStore)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreatePolicy", ctxFields, err.Error())
		return
	}

	policy.Metadata = getCreateMetaData(r)
	ret, err := mgr.PolicyStore.CreatePolicy(serviceName, &policy)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteFailedAuditLog("CreatePolicy", ctxFields, err.Error())
		return
	}

	logging.WriteSucceededAuditLog("CreatePolicy", ctxFields, nil)
	httputils.SendCreatedResponse(w, &ret)
}

func (mgr *RESTService) DeletePolicies(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}

	if err := mgr.PolicyStore.DeletePolicies(serviceName); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("DeletePolicies", serviceName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("DeletePolicies", serviceName, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	serviceName, policyIDStr := ParseRequestURI(r)
	if len(serviceName) == 0 || len(policyIDStr) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name or policy ID.",
		})
		return
	}

	// Audit contextual fields for request
	ctxFields := log.Fields{
		"serviceName": serviceName,
		"policyId":    policyIDStr,
	}

	if err := mgr.PolicyStore.DeletePolicy(serviceName, policyIDStr); err != nil {
		httputils.HandleError(w, err)
		logging.WriteFailedAuditLog("DeletePolicy", ctxFields, err.Error())
		return
	}

	logging.WriteSucceededAuditLog("DeletePolicy", ctxFields, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) GetPolicy(w http.ResponseWriter, r *http.Request) {
	serviceName, policyIDStr := ParseRequestURI(r)
	if len(serviceName) == 0 || len(policyIDStr) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name or policy ID.",
		})
		return
	}

	// Audit log for request
	ctxFields := log.Fields{
		"serviceName": serviceName,
		"policyId":    policyIDStr,
	}

	policy, err := mgr.PolicyStore.GetPolicy(serviceName, policyIDStr)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteFailedAuditLog("GetPolicy", ctxFields, err.Error())
		return
	}

	logging.WriteSucceededAuditLog("GetPolicy", ctxFields, map[string]interface{}{"policy": policy})
	httputils.SendOKResponse(w, &policy)
}

func (mgr *RESTService) ListPolicies(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}
	filters := ParseForFilters(r)
	policies, err := mgr.PolicyStore.ListAllPolicies(serviceName, filters)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("ListPolicies", serviceName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("ListPolicies", serviceName, len(policies))
	if len(policies) == 0 {
		httputils.SendEmptyListResponse(w)
		return
	}
	httputils.SendOKResponse(w, policies)
}

// Role policy management
func (mgr *RESTService) CreateRolePolicy(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}
	var rolePolicy pms.RolePolicy
	if err := decodeRequestBody(r, &rolePolicy); err != nil {
		httputils.HandleError(w, err)
		return
	}

	// Audit log for request
	ctxFields := log.Fields{
		"serviceName": serviceName,
		"rolePolicy":  &rolePolicy,
	}

	err := pmsimpl.CheckRolePolicy(serviceName, &rolePolicy, mgr.PolicyStore)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreateRolePolicy", ctxFields, err.Error())
		return
	}

	rolePolicy.Metadata = getCreateMetaData(r)
	ret, err := mgr.PolicyStore.CreateRolePolicy(serviceName, &rolePolicy)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteFailedAuditLog("CreateRolePolicy", ctxFields, err.Error())
		return
	}

	logging.WriteSucceededAuditLog("CreateRolePolicy", ctxFields, nil)
	httputils.SendCreatedResponse(w, &ret)
}

func (mgr *RESTService) DeleteRolePolicies(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}

	if err := mgr.PolicyStore.DeleteRolePolicies(serviceName); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("DeleteRolePolicies", serviceName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("DeleteRolePolicies", serviceName, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) DeleteRolePolicy(w http.ResponseWriter, r *http.Request) {
	serviceName, rolePolicyIDStr := ParseRequestURI(r)
	if len(serviceName) == 0 || len(rolePolicyIDStr) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name or role policy ID.",
		})
		return
	}

	// Audit contextual fields for request
	ctxFields := log.Fields{
		"serviceName":  serviceName,
		"rolePolicyId": rolePolicyIDStr,
	}

	if err := mgr.PolicyStore.DeleteRolePolicy(serviceName, rolePolicyIDStr); err != nil {
		httputils.HandleError(w, err)
		logging.WriteFailedAuditLog("DeleteRolePolicy", ctxFields, err.Error())
		return
	}

	logging.WriteSucceededAuditLog("DeleteRolePolicy", ctxFields, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) GetRolePolicy(w http.ResponseWriter, r *http.Request) {
	serviceName, rolePolicyIDStr := ParseRequestURI(r)
	if len(serviceName) == 0 || len(rolePolicyIDStr) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name or role policy ID.",
		})
		return
	}

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serviceName":  serviceName,
		"rolePolicyId": rolePolicyIDStr,
	}

	rolePolicy, err := mgr.PolicyStore.GetRolePolicy(serviceName, rolePolicyIDStr)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteFailedAuditLog("GetRolePolicy", ctxFields, err.Error())
		return
	}

	logging.WriteSucceededAuditLog("GetRolePolicy", ctxFields, map[string]interface{}{"rolePolicy": rolePolicy})
	httputils.SendOKResponse(w, &rolePolicy)
}

func (mgr *RESTService) ListRolePolicies(w http.ResponseWriter, r *http.Request) {
	serviceName, _ := ParseRequestURI(r)
	if len(serviceName) == 0 {
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: "Invalid service name.",
		})
		return
	}
	filters := ParseForFilters(r)
	rolePolicies, err := mgr.PolicyStore.ListAllRolePolicies(serviceName, filters)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("ListRolePolicies", serviceName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("ListRolePolicies", serviceName, len(rolePolicies))

	if len(rolePolicies) == 0 {
		httputils.SendEmptyListResponse(w)
		return
	}

	httputils.SendOKResponse(w, &rolePolicies)
}

func (mgr *RESTService) CreateFunction(w http.ResponseWriter, r *http.Request) {
	var cf pms.Function
	err := decodeRequestBody(r, &cf)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreateFunction", nil, err.Error())
		return
	}

	err = pmsimpl.CheckFunction(&cf, mgr.PolicyStore)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreateFunction", &cf, err.Error())
		return
	}
	cf.Metadata = getCreateMetaData(r)
	ret, err := mgr.PolicyStore.CreateFunction(&cf)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("CreateFunction", &cf, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("CreateFunction", &cf, nil)
	httputils.SendCreatedResponse(w, ret)
}

func (mgr *RESTService) DeleteFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	funcName, ok := vars["functionName"]
	if !ok || funcName == "" {
		msg := "functionName is not specified"
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: msg,
		})
		logging.WriteSimpleFailedAuditLog("DeleteFunction", nil, msg)
		return
	}

	if err := mgr.PolicyStore.DeleteFunction(funcName); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("DeleteFunction", funcName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("DeleteFunction", funcName, nil)
	w.WriteHeader(http.StatusNoContent)

}

func (mgr *RESTService) DeleteFunctions(w http.ResponseWriter, r *http.Request) {
	if err := mgr.PolicyStore.DeleteFunctions(); err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("DeleteFunctions", nil, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("DeleteFunctions", nil, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (mgr *RESTService) GetFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	funcName, ok := vars["functionName"]
	if !ok || funcName == "" {
		msg := "functionName is not specified"
		httputils.SendBadRequestResponse(w, &httputils.ErrorResponse{
			Error: msg,
		})
		logging.WriteSimpleFailedAuditLog("GetFunction", nil, msg)
		return
	}

	cf, err := mgr.PolicyStore.GetFunction(funcName)
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("GetFunction", funcName, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("GetFunction", funcName, nil)
	httputils.SendOKResponse(w, cf)
}

func (mgr *RESTService) ListFunctions(w http.ResponseWriter, r *http.Request) {
	functions, err := mgr.PolicyStore.ListAllFunctions("")
	if err != nil {
		httputils.HandleError(w, err)
		logging.WriteSimpleFailedAuditLog("ListFunctions", nil, err.Error())
		return
	}

	logging.WriteSimpleSucceededAuditLog("ListFunctions", nil, len(functions))
	if len(functions) == 0 {
		httputils.SendEmptyListResponse(w)
		return
	}
	httputils.SendOKResponse(w, functions)
}
