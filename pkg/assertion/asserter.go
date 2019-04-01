//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package assertion

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	adsapi "github.com/oracle/speedle/api/ads"
	log "github.com/sirupsen/logrus"
)

const (
	// TokenKey http header key name for asserting token
	TokenKey = "x-token"
	// IdpTypeKey http header key name for idp
	IdpTypeKey = "x-idp"
	// AllowedIDDKey http header key name for allowed IDD
	AllowedIDDKey = "x-allowedIDD"
	// RequestHeaderKey http header key name for extra headers which will be passed to IDP
	RequestHeaderKey = "x-ecid"
)

// AssertResponse assertion response
type AssertResponse struct {
	Principals []*adsapi.Principal    `json:"principals,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	ErrCode    int                    `json:"errCode"`
	ErrMessage string                 `json:"errMessage,omitempty"`
}

// TokenAsserter asserter interface
type TokenAsserter interface {
	// AssertToken assert token and generate subject to represent the identity
	AssertToken(token string, idpType string, allowedIDD string, requestHeaders map[string]string) (*AssertResponse, error)
}

// AsserterConfig asserter webhook client configuration
type AsserterConfig struct {
	Endpoint    string `json:"endpoint"`
	CACert      string `json:"caCert"`
	ClientCert  string `json:"clientCert"`
	ClientKey   string `json:"clientKey"`
	HTTPTimeout int    `json:"httpTimeout"`
}

// WebHookAsserter implements asserter client interface
type WebHookAsserter struct {
	ServerEndpoint string
	caCert         string
	clientCert     string
	clientKey      string
	arstenant      *string
	httpClient     *http.Client
}

// NewAsserter create asserter webhook client
func NewAsserter(conf *AsserterConfig, tenant *string) (TokenAsserter, error) {
	log.Debugf("conf: %v, tenant: %v", conf, tenant)
	if conf == nil || len(conf.Endpoint) == 0 {
		log.Errorf("asserter configuration is nil or endpoint is emtpy")
		return nil, fmt.Errorf("asserter configuration is nil or endpoint is emtpy")
	}
	if conf.HTTPTimeout <= 0 {
		conf.HTTPTimeout = 10
	}
	a := WebHookAsserter{
		ServerEndpoint: strings.ToLower(conf.Endpoint),
		caCert:         conf.CACert,
		clientCert:     conf.ClientCert,
		clientKey:      conf.ClientKey,
		arstenant:      tenant,
	}

	tr := http.Transport{
		MaxIdleConns:    1000,
		IdleConnTimeout: 60 * time.Second,
		Proxy:           http.ProxyFromEnvironment,
	}

	if strings.HasPrefix(a.ServerEndpoint, "https") {
		tlsConf := &tls.Config{}
		if len(a.caCert) > 0 {
			caCert, err := ioutil.ReadFile(a.caCert)
			if err != nil {
				return nil, err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConf.RootCAs = caCertPool
		}
		if len(a.clientCert) > 0 && len(a.clientKey) > 0 {
			cert, err := tls.LoadX509KeyPair(a.clientCert, a.clientKey)
			if err != nil {
				return nil, err
			}
			tlsConf.Certificates = []tls.Certificate{cert}
		}
		tr.TLSClientConfig = tlsConf
	}

	a.httpClient = &http.Client{
		Transport: &tr,
		Timeout:   time.Duration(conf.HTTPTimeout) * time.Second,
	}

	log.Debugf("asserter created: %v", a)

	return &a, nil
}

// AssertToken assert token via webhook
func (a *WebHookAsserter) AssertToken(token string, idpType string, allowedIDD string, requestHeaders map[string]string) (*AssertResponse, error) {
	log.Debugf("token: %s, idpType: %s, allowedIDD: %s, requestHeaders: %v", token, idpType, allowedIDD, requestHeaders)

	if len(token) == 0 {
		log.Errorf("token is empty")
		return nil, fmt.Errorf("token is empty")
	}

	requestURL := a.ServerEndpoint
	if a.arstenant != nil {
		requestURL += *a.arstenant
	}

	req, errReq := http.NewRequest(http.MethodGet, requestURL, nil)
	if errReq != nil {
		log.Errorf("NewRequest error: %v", errReq)
		return nil, errReq
	}
	req.Header.Add(TokenKey, token)
	req.Header.Add(IdpTypeKey, idpType)

	if len(allowedIDD) > 0 {
		req.Header.Add(AllowedIDDKey, allowedIDD)
	}
	if requestHeaders != nil {
		keys := ""
		for k, v := range requestHeaders {
			req.Header.Add(k, v)
			keys = keys + k + ","
		}
		if len(keys) > 1 {
			keys = keys[0 : len(keys)-1]
		}
		req.Header.Add(RequestHeaderKey, keys)
	}

	resp, errResp := a.httpClient.Do(req)
	if errResp != nil {
		log.Errorf("Do error: %v", errResp)
		return nil, errResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("assertion error, status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("asserter error, status code: %d", resp.StatusCode)
	}

	raw, errRaw := ioutil.ReadAll(resp.Body)
	if errRaw != nil {
		log.Errorf("ReadAll error: %v", errRaw)
		return nil, errRaw
	}

	var ar AssertResponse

	errJSON := json.Unmarshal(raw, &ar)
	if errJSON != nil {
		log.Errorf("Unmarshal error: %v", errJSON)
		return nil, errJSON
	}

	// flag error if asserter indicates failure
	if ar.ErrCode != 0 {
		log.Errorf("assertion failure: %v", ar)
		return nil, fmt.Errorf("ErrCode: %d, ErrMsg: %s", ar.ErrCode, ar.ErrMessage)
	}

	log.Debugf("asserted: %v", ar)

	return &ar, nil
}
