//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/oracle/speedle/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func errorCodeToHTTPStatus(err error) int {
	switch errors.Code(err) {
	case errors.EntityNotFound:
		return http.StatusNotFound
	case errors.EntityAlreadyExists:
		return http.StatusConflict
	case errors.SerializationError:
		return http.StatusInternalServerError
	case errors.StoreError:
		return http.StatusInternalServerError
	case errors.InvalidRequest:
		return http.StatusBadRequest
	case errors.ExceedLimit:
		return http.StatusForbidden
	default:
		// Unknown status
		return http.StatusInternalServerError
	}
}

func HandleError(w http.ResponseWriter, err error) {
	log.Warningf("Handle error, err: %+v", err)
	SendResponse(w, errorCodeToHTTPStatus(err), &ErrorResponse{
		Error: err.Error(),
	})
}

func SendBadRequestResponse(w http.ResponseWriter, object interface{}) {
	SendResponse(w, http.StatusBadRequest, object)
}

func SendInternalErrorResponse(w http.ResponseWriter, object interface{}) {
	SendResponse(w, http.StatusInternalServerError, object)
}

func SendBasicAtzErrorResponse(w http.ResponseWriter, realm string) {
	headValue := "Basic"
	if len(realm) != 0 {
		headValue = fmt.Sprintf("%s realm=\"%s\"", headValue, realm)
	}
	w.Header().Set("WWW-Authenticate", headValue)
	w.WriteHeader(http.StatusUnauthorized)
}

// VerifyContentType reads content type from HTTP header,
// if the content type can not be found in expectedContentTypes,
// write response with status 415 Unsupported Media Type and returns false.
func VerifyContentType(w http.ResponseWriter, r *http.Request, expectedContentTypes []string) bool {
	requestContentType := r.Header.Get("Content-Type")
	if len(requestContentType) > 0 {
		// Content-Type: "application/json; charset=utf-8"
		cts := strings.Split(requestContentType, ";")
		ct := strings.TrimSpace(cts[0])
		for _, expectedContentType := range expectedContentTypes {
			if ct == expectedContentType {
				// ignore charset validates
				return true
			}
		}
	}
	w.WriteHeader(http.StatusUnsupportedMediaType)
	return false
}

func SendBearerAtzErrorResponse(w http.ResponseWriter, realm string, err string, errDesc string) {
	headValue := "Bearer"
	value := ""
	if len(realm) != 0 {
		value = "realm=\"" + realm + "\""
	}
	if len(err) != 0 {
		if len(value) != 0 {
			value = value + ", "
		}
		value = value + "error=\"" + err + "\""
	}
	if len(errDesc) != 0 {
		if len(value) != 0 {
			value = value + ", "
		}
		value = value + "error_description=\"" + errDesc + "\""
	}
	if len(value) != 0 {
		headValue = headValue + " " + value
	}

	w.Header().Set("WWW-Authenticate", headValue)
	w.WriteHeader(http.StatusUnauthorized)
}

func SendPageNotFoundResponse(w http.ResponseWriter) {
	http.Error(w, "", http.StatusNotFound)
}

func SendEmptyListResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[]"))
}

func SendResponse(w http.ResponseWriter, statusCode int, object interface{}) {
	var payload []byte
	if object != nil {
		payload, _ = json.Marshal(object)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if object != nil {
		w.Write(payload)
	}
}

func SendCreatedResponse(w http.ResponseWriter, object interface{}) {
	SendResponse(w, http.StatusCreated, object)
}

func SendOKResponse(w http.ResponseWriter, object interface{}) {
	SendResponse(w, http.StatusOK, object)
}
