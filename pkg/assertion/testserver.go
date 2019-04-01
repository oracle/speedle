//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package assertion

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
)

type asserterHandler struct {
	tb testing.TB
}

// Just prepend "vault:v1:" prefix as encrypted text.
func (h *asserterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		token := r.Header.Get("x-token")
		subj := &AssertResponse{}

		if len(token) == 0 {
			subj.ErrMessage = "token is empty"
			sendResp(w, http.StatusBadRequest, subj)
			return
		} else if token == "test-token" {
			subj.ErrCode = http.StatusBadRequest
			subj.ErrMessage = "invalid token"

			sendResp(w, http.StatusBadRequest, subj)
		} else {
			subj.ErrCode = 0
			subj.ErrMessage = ""
			subj.Principals = []*adsapi.Principal{
				{
					Type: adsapi.PRINCIPAL_TYPE_USER,
					Name: "testUser",
				},
			}

			sendResp(w, http.StatusOK, subj)
		}
	}

}

func sendResp(w http.ResponseWriter, status int, data *AssertResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	raw, _ := json.Marshal(data)
	w.Write(raw)
}

func defaultTestHandlers(tb testing.TB) map[string]http.Handler {
	return map[string]http.Handler{
		"/assert": &asserterHandler{tb},
	}
}

func NewTestServer(tb testing.TB, handlers map[string]http.Handler) *httptest.Server {
	mux := http.NewServeMux()
	if handlers == nil {
		handlers = defaultTestHandlers(tb)
	}
	for path, handler := range handlers {
		mux.Handle(path, handler)
	}
	server := httptest.NewUnstartedServer(mux)

	server.Start()

	return server
}
