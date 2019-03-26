//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.
package adsrest

import (
	"net/http/httptest"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/assertion"
	"github.com/oracle/speedle/pkg/cfg"
	"github.com/oracle/speedle/pkg/eval"
	_ "github.com/oracle/speedle/pkg/store/file"
)

var storeLocation = "./fakestore.json"

func NewTestServerWithConfig(conf *cfg.Config) (*httptest.Server, error) {
	e, err := newEvaluator(conf)
	if err != nil {
		return nil, err
	}
	routers, err := NewRouter(e)
	if err != nil {
		return nil, err
	}
	server := httptest.NewUnstartedServer(routers)
	server.Start()

	return server, nil
}

func NewTestServer() (*httptest.Server, error) {
	conf := GenerateServerConfig()
	evaluator, err := newEvaluator(conf)
	if err != nil {
		return nil, err
	}
	routers, err := NewRouter(evaluator)
	if err != nil {
		return nil, err
	}
	server := httptest.NewUnstartedServer(routers)
	server.Start()

	return server, nil
}

func GenerateServerConfig() *cfg.Config {
	var conf cfg.Config
	var storeConf cfg.StoreConfig
	storeConf.StoreType = cfg.StorageTypeFile
	storeConf.StoreProps = make(map[string]interface{})
	storeConf.StoreProps["FileLocation"] = storeLocation
	conf.StoreConfig = &storeConf
	conf.EnableWatch = false
	return &conf
}

func newEvaluator(conf *cfg.Config) (eval.InternalEvaluator, error) {
	evaluator, err := eval.NewFromConfig(conf)
	if err != nil {
		return nil, err
	}

	as, errLoadAsserter := assertion.NewAsserter(conf.AsserterWebhookConfig, nil)
	if errLoadAsserter != nil {
		return nil, errLoadAsserter
	} else {
		f := func(ctx *adsapi.RequestContext) error {
			if ctx.Subject != nil &&
				len(ctx.Subject.TokenType) != 0 &&
				len(ctx.Subject.Token) != 0 {
				tokenType := ctx.Subject.TokenType
				token := ctx.Subject.Token
				s, err := as.AssertToken(token, tokenType, "", nil)
				if err == nil {
					for _, p := range s.Principals {
						ctx.Subject.Principals = append(ctx.Subject.Principals, p)
					}
				} else {
					return err
				}
			}
			return nil
		}
		evaluator.SetAsserterFunc(f)
	}

	return evaluator, nil
}
