//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package expenses

import (
	"net/http"
	"strings"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/eval"
)

// Wrap wraps an HTTP handler with a new HTTP handler for authorization
func Wrap(handler http.Handler, service, spdlLoc string) (http.HandlerFunc, error) {
	ev, err := eval.NewFromFile(spdlLoc, true)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		act := strings.ToLower(r.Method)
		res := r.RequestURI
		user, _, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized.\n"))
			return
		}

		// Construct context for evaluator
		reqCtx := ads.RequestContext{
			Subject: &ads.Subject{
				Principals: []*ads.Principal{
					{
						Type: "user",
						Name: user,
					},
				},
			},
			ServiceName: service,
			Action:      act,
			Resource:    res,
		}

		// Call evaluator
		allowed, _, _ := ev.IsAllowed(reqCtx)

		if !allowed {
			// Not allowed, send 403
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden.\n"))
			return
		}

		// Allowed, call delecated handler
		handler.ServeHTTP(w, r)
	}, nil
}
