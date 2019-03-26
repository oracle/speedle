//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package expenses

import (
	"net/http"
	"strings"
)

// ExpenseHTTPHandler is the handler for all handler functions
var ExpenseHTTPHandler http.Handler

func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/reports", reportHandler)

	ExpenseHTTPHandler = mux
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.ToLower(r.Method) + "ing an expense report is done\n"))
}
