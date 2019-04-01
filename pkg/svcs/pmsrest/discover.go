//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pmsrest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/httputils"
	"github.com/oracle/speedle/pkg/logging"
	"github.com/oracle/speedle/pkg/store"
	log "github.com/sirupsen/logrus"
)

type GetDiscoverRequestsResponse struct {
	Requests []*ads.RequestContext `json:"requests"`
	Revision int64                 `json:"revision"`
}

type GetDiscoverPoliciesResponse struct {
	Services []*pms.Service `json:"services"`
	Revision int64          `json:"revision"`
}

func (e *RESTService) GetAllDiscoverRequests(w http.ResponseWriter, r *http.Request) {
	if !e.checkPolicyForDiscover(w) {
		return
	}

	last := r.URL.Query().Get("last")
	revisionStr := r.URL.Query().Get("revision")

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"last":     last,
		"revision": revisionStr,
	}

	if strings.EqualFold("true", last) { //get last discover request
		request, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GetLastDiscoverRequest("")
		if err != nil {
			log.Errorf("%v, Cause: %v", err, errors.Cause(err))
			httputils.HandleError(w, err)

			// Audit log
			logging.WriteFailedAuditLog("GetAllDiscoverRequests", ctxFields, err.Error())
			return
		}
		response := GetDiscoverRequestsResponse{Requests: []*ads.RequestContext{}, Revision: revision}
		if request != nil {
			response.Requests = append(response.Requests, request)
		}

		httputils.SendOKResponse(w, &response)

		// Audit log
		logging.WriteSucceededAuditLog("GetAllDiscoverRequests", ctxFields, map[string]interface{}{"lastRequest": request})
	} else if len(revisionStr) != 0 { //get discover requests since revision
		revision, err := strconv.ParseInt(revisionStr, 10, 64)
		if err != nil {
			err = errors.Errorf(errors.InvalidRequest, "invlid revision number %q", revision)
			log.Error(err)
			httputils.HandleError(w, err)

			// Audit log
			logging.WriteFailedAuditLog("GetAllDiscoverRequests", ctxFields, err.Error())
			return
		}
		requests, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GetDiscoverRequestsSinceRevision("", revision)
		if err != nil {
			log.Errorf("%v, Cause: %v", err, errors.Cause(err))
			httputils.HandleError(w, err)
			// Audit log
			logging.WriteFailedAuditLog("GetAllDiscoverRequests", ctxFields, err.Error())
			return
		}
		response := GetDiscoverRequestsResponse{Requests: requests, Revision: revision}
		httputils.SendOKResponse(w, &response)

		// Audit log
		logging.WriteSucceededAuditLog("GetAllDiscoverRequests", ctxFields, map[string]interface{}{"requestCount": len(requests)})
	} else {
		//get all service requests
		requests, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GetDiscoverRequests("")
		if err != nil {
			log.Errorf("%v, Cause: %v", err, errors.Cause(err))
			httputils.HandleError(w, err)
			// Audit log
			logging.WriteFailedAuditLog("GetAllDiscoverRequests", ctxFields, err.Error())
			return
		}
		response := GetDiscoverRequestsResponse{Requests: requests, Revision: revision}
		httputils.SendOKResponse(w, &response)
		// Audit log
		logging.WriteSucceededAuditLog("GetAllDiscoverRequests", ctxFields, map[string]interface{}{"requestCount": len(requests)})
	}
}

func (e *RESTService) GetDiscoverRequests(w http.ResponseWriter, r *http.Request) {
	if !e.checkPolicyForDiscover(w) {
		return
	}

	serviceName, err := getServiceNameFromRequest(w, r)
	if err != nil {
		return
	}

	last := r.URL.Query().Get("last")
	revisionStr := r.URL.Query().Get("revision")

	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serverName": serviceName,
		"last":       last,
		"revision":   revisionStr,
	}

	if strings.EqualFold("true", last) { //get last discover request
		request, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GetLastDiscoverRequest(serviceName)
		if err != nil {
			log.Errorf("%v, Cause: %v", err, errors.Cause(err))
			httputils.HandleError(w, err)

			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return
		}
		response := GetDiscoverRequestsResponse{Requests: []*ads.RequestContext{}, Revision: revision}
		if request != nil {
			response.Requests = append(response.Requests, request)
		}
		httputils.SendOKResponse(w, &response)

		// Audit log
		logging.WriteSucceededAuditLog("GetDiscoverRequests", ctxFields, map[string]interface{}{"lastRequest": request})
	} else if len(revisionStr) != 0 { //get discover requests since revision
		revision, err := strconv.ParseInt(revisionStr, 10, 64)
		if err != nil {
			err = errors.Errorf(errors.InvalidRequest, "invlid revision number %q", revision)
			log.Error(err)
			httputils.HandleError(w, err)

			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return
		}
		requests, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GetDiscoverRequestsSinceRevision(serviceName, revision)
		if err != nil {
			log.Errorf("%v, Cause: %v", err, errors.Cause(err))
			httputils.HandleError(w, err)

			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return
		}
		response := GetDiscoverRequestsResponse{Requests: requests, Revision: revision}
		httputils.SendOKResponse(w, &response)

		// Audit log
		logging.WriteSucceededAuditLog("GetDiscoverRequests", ctxFields, map[string]interface{}{"requestCount": len(requests)})
	} else {
		//get all service requests
		requests, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GetDiscoverRequests(serviceName)
		if err != nil {
			log.Errorf("%v, Cause: %v", err, errors.Cause(err))
			httputils.HandleError(w, err)
			// Audit log
			logging.WriteFailedAuditLog("GetDiscoverRequests", ctxFields, err.Error())
			return
		}
		response := GetDiscoverRequestsResponse{Requests: requests, Revision: revision}
		httputils.SendOKResponse(w, &response)
		// Audit log
		logging.WriteSucceededAuditLog("GetDiscoverRequests", ctxFields, map[string]interface{}{"requestCount": len(requests)})
	}
}

func (e *RESTService) ResetDiscoverRequests(w http.ResponseWriter, r *http.Request) {
	if !e.checkPolicyForDiscover(w) {
		return
	}

	serviceName, err := getServiceNameFromRequest(w, r)
	if err != nil {
		return
	}

	if err := e.PolicyStore.(store.DiscoverRequestManager).ResetDiscoverRequests(serviceName); err != nil {
		log.Errorf("%v, Cause: %v", err, errors.Cause(err))
		httputils.HandleError(w, err)

		// Audit log
		logging.WriteSimpleFailedAuditLog("ResetDiscoverRequests", serviceName, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// Audit log
	logging.WriteSimpleSucceededAuditLog("ResetDiscoverRequests", serviceName, nil)
}

func (e *RESTService) ResetAllDiscoverRequests(w http.ResponseWriter, r *http.Request) {
	if !e.checkPolicyForDiscover(w) {
		return
	}

	err := e.PolicyStore.(store.DiscoverRequestManager).ResetDiscoverRequests("")
	if err != nil {
		log.Errorf("%v, Cause: %v", err, errors.Cause(err))
		httputils.HandleError(w, err)
		// Audit log
		logging.WriteSimpleFailedAuditLog("ResetAllDiscoverRequests", nil, err.Error())
	}
	w.WriteHeader(http.StatusNoContent)

	// Audit log
	logging.WriteSimpleSucceededAuditLog("ResetAllDiscoverRequests", nil, nil)
}

func (e *RESTService) GetDiscoverPolicies(w http.ResponseWriter, r *http.Request) {
	if !e.checkPolicyForDiscover(w) {
		return
	}

	vads := mux.Vars(r)
	serviceName := vads["serviceName"]
	principalType := r.URL.Query().Get("principalType")
	principalName := r.URL.Query().Get("principalName")
	principalIDD := r.URL.Query().Get("principalIDD")
	// Audit contextual fields for request
	ctxFields := map[string]interface{}{
		"serverName":    serviceName,
		"principalType": principalType,
		"principalName": principalName,
		"principalIDD":  principalIDD,
	}

	serviceMap, revision, err := e.PolicyStore.(store.DiscoverRequestManager).GeneratePolicies(serviceName, principalType, principalName, principalIDD)
	if err != nil {
		log.Error(err)
		httputils.HandleError(w, err)

		// Audit log
		logging.WriteFailedAuditLog("GetDiscoverPolicies", ctxFields, err.Error())
		return
	}
	services := []*pms.Service{}
	for _, value := range serviceMap {
		services = append(services, value)
	}
	response := GetDiscoverPoliciesResponse{
		Services: services,
		Revision: revision,
	}
	httputils.SendOKResponse(w, &response)

	// Audit log
	serviceCount := len(services)
	logging.WriteSucceededAuditLog("GetDiscoverPolicies", ctxFields, map[string]interface{}{
		"revision":     revision,
		"serviceCount": serviceCount,
	})

}

func (e *RESTService) checkPolicyForDiscover(w http.ResponseWriter) bool {
	if _, ok := e.PolicyStore.(store.DiscoverRequestManager); !ok {
		err := errors.Errorf(errors.InvalidRequest, "%q policy store doesn't support discover request management", e.PolicyStore.Type())
		log.Error(err)
		httputils.HandleError(w, err)
		return false
	}
	return true
}

func getServiceNameFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	vads := mux.Vars(r)
	serviceName := vads["serviceName"]
	if len(serviceName) == 0 {
		err := errors.New(errors.InvalidRequest, "no service name found in request")
		log.Error(err)
		httputils.HandleError(w, err)

		// Audit log
		logging.WriteSimpleFailedAuditLog("GetDiscoverRequests", serviceName, err.Error())
		return "", err
	}

	return serviceName, nil
}
