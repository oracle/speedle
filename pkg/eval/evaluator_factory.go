//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"github.com/oracle/speedle/pkg/cfg"
	"github.com/oracle/speedle/pkg/store"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"

	log "github.com/sirupsen/logrus"
)

//New creates a policy evaluator based on the given configuration file
func New(configFile string) (InternalEvaluator, error) {
	conf, err := cfg.ReadConfig(configFile)
	if err != nil {
		return nil, err
	}
	return NewFromConfig(conf)
}

// NewFromFile loads policies from a policy file, and returns an evaluator instance
func NewFromFile(fileLoc string, isWatch bool) (adsapi.PolicyEvaluator, error) {
	storeConfig := cfg.StoreConfig{
		StoreType: "file",
		StoreProps: map[string]interface{}{
			"FileLocation": fileLoc,
		},
	}

	// For file store, watchuing a policy store is disabled
	return NewFromConfig(&cfg.Config{
		StoreConfig: &storeConfig,
		EnableWatch: isWatch,
	})
}

//NewFromConfig creates a policy evaluator based on the given configuration file
func NewFromConfig(conf *cfg.Config) (InternalEvaluator, error) {
	s, err := store.NewStore(conf.StoreConfig.StoreType, conf.StoreConfig.StoreProps)
	if err != nil {
		return nil, err
	}

	return NewWithStore(conf, s)
}

// NewWithStore creates a policy evaluator with policy store
func NewWithStore(conf *cfg.Config, s pms.PolicyStoreManagerADS) (InternalEvaluator, error) {
	ps, err := s.ReadPolicyStore()
	if err != nil {
		return nil, err
	}

	var updateChan pms.StorageChangeChannel
	if conf.EnableWatch {
		log.Info("Watch policy store.")
		updateChan, err = s.Watch()
		if err != nil {
			return nil, err
		}
	}

	runtimePolicyStore := NewRuntimePolicyStore()
	runtimePolicyStore.init(ps, conf.FuncsvcEndpoint)

	p := &PolicyEvalImpl{
		RuntimePolicyStore: runtimePolicyStore,
		Store:              s,
	}

	// start a goroutine watching to the channel for update events and
	// refresh runtime cache accordingly once receiving any events
	if updateChan != nil {
		go p.updateRuntimeCacheWithStoreChange(updateChan)
	}

	p.cleanExpiredFunctionResultPeriodically()

	return p, nil
}
