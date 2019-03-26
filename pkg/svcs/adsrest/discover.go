//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package adsrest

import (
	"net/http"

	"github.com/oracle/speedle/pkg/httputils"
	"github.com/oracle/speedle/pkg/logging"
)

func (e *RESTService) Discover(w http.ResponseWriter, r *http.Request) {
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

	// assert token
	e.Evaluator.AssertToken(context)

	result, reason, err := e.Evaluator.Discover(*context)
	response := IsAllowedResponse{
		Allowed: result,
		Reason:  int32(reason),
	}

	// Audit log
	if err != nil {
		response.ErrorMessage = err.Error()
		logging.WriteSimpleFailedAuditLog("Discovery", context, err.Error())
	} else {
		logging.WriteSimpleSucceededAuditLog("Discovery", context, nil)
	}

	httputils.SendOKResponse(w, &response)
}
