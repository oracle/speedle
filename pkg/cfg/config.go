//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package cfg

import (
	"encoding/json"
	"io/ioutil"

	"github.com/oracle/speedle/pkg/assertion"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/logging"
)

const (
	StorageTypeFile = "file"
)

type StoreConfig struct {
	StoreType  string                 `json:"storeType"`
	StoreProps map[string]interface{} `json:"storeProps"`
}

type ServerConfig struct {
	Endpoint        string `json:"endpoint,omitempty"`
	Insecure        string `json:"insecure,omitempty"`
	EnableAuthz     string `json:"enableAuthz,omitempty"`
	KeyPath         string `json:"keyPath,omitempty"`
	CertPath        string `json:"certPath,omitempty"`
	ClientCertPath  string `json:"clientCertPath,omitempty"`
	ForceClientCert bool   `json:"forceClientCert,omitempty"`
}

type Config struct {
	StoreConfig           *StoreConfig              `json:"storeConfig"`
	EnableWatch           bool                      `json:"enableWatch,omitempty"`
	AsserterWebhookConfig *assertion.AsserterConfig `json:"asserterWebhookConfig,omitempty"`
	FuncsvcEndpoint       string                    `json:"funcsvcEndpoint,omitempty"`
	ServerConfig          *ServerConfig             `json:"serverConfig,omitempty"`
	LogConfig             *logging.LogConfig        `json:"logConfig,omitempty"`
	AuditLogConfig        *logging.LogConfig        `json:"auditLogConfig,omitempty"`
}

func ReadConfig(configFileLocation string) (*Config, error) {
	var config Config
	raw, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return &config, errors.Wrapf(err, errors.ConfigError, "failed to read configure file %s", configFileLocation)
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		err = errors.Wrapf(err, errors.ConfigError, "fauiled to unmarshal configure file %s", configFileLocation)
	}

	return &config, err
}

func ReadStoreConfig(configFileLocation string) (*StoreConfig, error) {
	var storeConfig StoreConfig
	raw, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return nil, errors.Wrapf(err, errors.ConfigError, "failed to read store configure file %s", configFileLocation)
	}
	err = json.Unmarshal(raw, &storeConfig)
	if err != nil {
		err = errors.Wrapf(err, errors.ConfigError, "failed to unmarshal store configure file %s", configFileLocation)
	}
	return &storeConfig, err
}
