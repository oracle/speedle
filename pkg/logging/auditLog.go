//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package logging

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	Response_Succeeded string = "succeeded"
	Response_Failed    string = "failed"

	Request_default_key  string = "requestValue"
	Response_default_key string = "responseValue"

	Audit_Prefix string = "[AUDIT]"

	TenantID_key     = "tenantid"
	Target_Index_Key = "target_index"
)

var (
	auditLogger = log.New()

	tenantID = ""
)

func SetTenantID(tid string) {
	tenantID = tid
}

func generateTargetIndex() string {
	t := time.Now()
	return fmt.Sprintf("audit-%s-%d%02d", tenantID, t.Year(), t.Month())
}

func AuditLog() *log.Logger {
	return auditLogger
}

// Initialize Audit Logger
func InitAuditLog(cfg *LogConfig) error {
	return InitLogInstance(auditLogger, cfg)
}

// The request or response have only one contextual field
func WriteSimpleSucceededAuditLog(apiName string, reqField interface{}, respField interface{}) {
	var reqFields, respFields map[string]interface{}

	if reqField != nil {
		reqFields = map[string]interface{}{
			Request_default_key: reqField,
		}
	}

	if respField != nil {
		respFields = map[string]interface{}{
			Response_default_key: respField,
		}
	}

	WriteSucceededAuditLog(apiName, reqFields, respFields)
}

// The request has only one contextual field
func WriteSimpleFailedAuditLog(apiName string, reqField interface{}, reason string) {
	var reqFields map[string]interface{}

	if reqField != nil {
		reqFields = map[string]interface{}{
			Request_default_key: reqField,
		}
	}

	WriteFailedAuditLog(apiName, reqFields, reason)
}

//The request or response may have multiple contextual fields
func WriteSucceededAuditLog(apiName string, reqFields map[string]interface{}, respFields map[string]interface{}) {
	if respFields != nil {
		respFields["response"] = Response_Succeeded
	} else {
		respFields = map[string]interface{}{
			"response": Response_Succeeded,
		}
	}
	writeAuditLog(apiName, reqFields, respFields, log.InfoLevel)
}

//The request may have multiple contextual fields
func WriteFailedAuditLog(apiName string, reqFields map[string]interface{}, reason string) {
	respFields := map[string]interface{}{
		"response": Response_Failed,
		"reason":   reason,
	}

	writeAuditLog(apiName, reqFields, respFields, log.ErrorLevel)
}

func writeAuditLog(apiName string, reqFields map[string]interface{}, respFields map[string]interface{}, level log.Level) {
	// Merge the reqFields and respFields, so that we only need to generate one audit log entry .
	// If we want to generate two audit log entries, then we shouldn't merge them.
	allFields := make(map[string]interface{})
	for k, v := range reqFields {
		allFields[k] = v
	}
	for k, v := range respFields {
		allFields[k] = v
	}

	// Set tenantID and target_index
	allFields[TenantID_key] = tenantID
	allFields[Target_Index_Key] = generateTargetIndex()

	ctxLogger := AuditLog().WithFields(allFields)
	msg := Audit_Prefix + apiName

	switch level {
	case log.InfoLevel:
		ctxLogger.Info(msg)
	case log.ErrorLevel:
		ctxLogger.Error(msg)
	case log.WarnLevel:
		ctxLogger.Warn(msg)
	case log.DebugLevel:
		ctxLogger.Debug(msg)
	case log.FatalLevel:
		ctxLogger.Fatal(msg)
	case log.PanicLevel:
		ctxLogger.Panic(msg)
	}
}
