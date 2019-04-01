package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"speedle/api/authz"
	"strconv"
	"time"
)

type restADSClient struct {
	isAllowedEndpoint string
	httpClient        *http.Client
}

var propertyDefalutValues = map[string]string{
	authz.HOST_PROP:      "127.0.0.1",
	authz.ADS_PORT_PROP:  "6734",
	authz.IS_SECURE_PROP: "false",
}

func getPropertyValue(properties map[string]string, propName string) (string, error) {
	value, ok := properties[propName]
	if !ok {
		value, ok = propertyDefalutValues[propName]
		if !ok {
			return "", fmt.Errorf("Property %s is not defined", propName)
		}
	}
	return value, nil
}

func New(properties map[string]string) (ADSClient, error) {
	host, _ := getPropertyValue(properties, authz.HOST_PROP)
	adsPort, _ := getPropertyValue(properties, authz.ADS_PORT_PROP)
	isSecureStr, _ := getPropertyValue(properties, authz.IS_SECURE_PROP)

	protocol := "https"
	isSecure, err := strconv.ParseBool(isSecureStr)
	if err != nil {
		return nil, err
	}
	if !isSecure {
		protocol = "http"
	}

	isAllowedEndpoint := fmt.Sprintf("%s://%s:%s/authz-check/v1/is-allowed", protocol, host, adsPort)

	log.Printf("Endpoint of is-allowed API: %s\n", isAllowedEndpoint)

	httpClient, err := newHttpClient(isAllowedEndpoint)
	if err != nil {
		return nil, err
	}

	return &restADSClient{
		isAllowedEndpoint: isAllowedEndpoint,
		httpClient:        httpClient,
	}, nil
}

// New construct a new ADS client instance
func newHttpClient(isAllowedEndpoint string) (*http.Client, error) {
	req, err := http.NewRequest(http.MethodGet, isAllowedEndpoint, nil)
	if err != nil {
		return nil, err
	}
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		return nil, err
	}

	httpTrans := &http.Transport{}
	if proxyURL != nil {
		log.Printf("Using proxy %s", proxyURL)
		httpTrans.Proxy = http.ProxyURL(proxyURL)
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: httpTrans,
	}, nil
}

func convertToJSONAttribute(attrs map[string]interface{}) []*jsonAttribute {
	ret := []*jsonAttribute{}
	if attrs == nil {
		return ret
	}
	for key, value := range attrs {
		ret = append(ret, &jsonAttribute{
			Name:  key,
			Value: value,
		})
	}
	return ret
}

func convertToJSONRequest(context *authz.RequestContext) *jsonContext {
	var subject *jsonSubject
	if context.Subject != nil {

		var principals []*jsonPrincipal
		for _, principal := range context.Subject.Principals {
			principals = append(principals, &jsonPrincipal{principal.Type, principal.Name, principal.IDD})
		}

		subject = &jsonSubject{
			Principals: principals,
			Attributes: convertToJSONAttribute(context.Subject.Attributes),
		}
	}
	return &jsonContext{
		Subject:     subject,
		ServiceName: context.ServiceName,
		Resource:    context.Resource,
		Action:      context.Action,
		Attributes:  convertToJSONAttribute(context.Attributes),
	}
}

func (c *restADSClient) IsAllowed(context authz.RequestContext) (bool, error) {

	payload, err := json.Marshal(convertToJSONRequest(&context))
	if err != nil {
		return false, err
	}

	log.Printf("IsAllowed request payload: %s\n", payload)

	req, err := http.NewRequest(http.MethodPost, c.isAllowedEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return false, err
	}
	req.Header.Add("Context-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Status %d returned", resp.StatusCode)
	}

	var result jsonIsAllowedResponse
	reader := json.NewDecoder(resp.Body)
	if err := reader.Decode(&result); err != nil {
		return false, err
	}
	return result.Allowed, nil
}
