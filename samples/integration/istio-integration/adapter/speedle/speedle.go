//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -f mixer/adapter/speedle/config/config.proto

package speedle

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"istio.io/istio/mixer/adapter/speedle/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/status"
	"istio.io/istio/mixer/template/authorization"
)

type (
	builder struct {
		//types         map[string]*authz.Type
		adapterConfig *config.Params
		configError   bool
	}

	handler struct {
		adsEndpoint    string
		speedleService string
		env            adapter.Env
		configError    bool
		failClose      bool
		expireTime     time.Time
	}

	/*
		request struct {
			subject     map[string]string
			servicename string
			resource    string
			action      string
		}*/

	JsonAttribute struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}

	JsonPrincipal struct {
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	}

	JsonSubject struct {
		Principals []*JsonPrincipal `json:"principals,omitempty"`
	}

	JsonContext struct {
		Subject     *JsonSubject     `json:"subject"`
		ServiceName string           `json:"serviceName"`
		Resource    string           `json:"resource"`
		Action      string           `json:"action"`
		Attributes  []*JsonAttribute `json:"attributes"`
	}
)

///////////////// Configuration Methods ///////////////

func (b *builder) SetAuthorizationTypes(types map[string]*authorization.Type) {}

func (b *builder) SetAdapterConfig(cfg adapter.Config) {
	b.adapterConfig = cfg.(*config.Params)
}

// To support fail close, Validate will not append errors to ce
// It set configError to true, then HandleAuthz return permission denied respose
func (b *builder) Validate() (ce *adapter.ConfigErrors) {
	b.configError = false
	if len(b.adapterConfig.AdsEndpoint) == 0 {
		b.configError = true
		return
	}

	return
}

func (b *builder) Build(context context.Context, env adapter.Env) (adapter.Handler, error) {

	return &handler{
		adsEndpoint:    b.adapterConfig.AdsEndpoint,
		speedleService: b.adapterConfig.SpeedleService,
		env:            env,
		configError:    b.configError,
		failClose:      b.adapterConfig.FailClose,
	}, nil
}

////////////////// Runtime Methods //////////////////////////

func convertActionObjectToMap(action *authorization.Action) map[string]interface{} {
	result := map[string]interface{}{}
	resource := ""

	if len(action.Namespace) > 0 {
		resource = resource + action.Namespace
	}

	if len(action.Service) > 0 {
		if resource == "" {
			resource = action.Service
		} else {
			resource = resource + "/" + action.Service
		}
	}

	if len(action.Path) > 0 {
		if resource == "" {
			resource = action.Path
		} else {
			resource = resource + action.Path
		}
	}

	result["resource"] = resource

	if len(action.Method) > 0 {
		result["action"] = action.Method
	} else {
		result["action"] = "access"
	}

	if action.Properties != nil {
		properties := map[string]interface{}{}
		count := 0
		for key, val := range action.Properties {
			properties[key] = val
			count++
		}
		if count > 0 {
			result["properties"] = properties
		}
	}

	return result
}

func convertSubjectObjectToMap(subject *authorization.Subject) map[string]interface{} {
	result := map[string]interface{}{}
	if len(subject.User) > 0 {
		result["user"] = subject.User
	}

	if len(subject.Groups) > 0 {
		result["groups"] = []interface{}{subject.Groups}
	}

	if subject.Properties != nil {
		properties := map[string]interface{}{}
		count := 0
		for key, val := range subject.Properties {
			properties[key] = val
			count++
		}
		if count > 0 {
			result["properties"] = properties
		}
	}

	return result
}

func (h *handler) HandleAuthorization(context context.Context, instance *authorization.Instance) (adapter.CheckResult, error) {
	fmt.Println("Incoming a new handle authorization request")
	fmt.Printf("%+v\n", instance)
	// Handle configuration error
	if h.configError {
		retStatus := status.OK

		if h.failClose {
			retStatus = status.WithPermissionDenied("speedle: request was rejected")
		}

		return adapter.CheckResult{
			Status: retStatus,
		}, nil
	}

	subjectMap := convertSubjectObjectToMap(instance.Subject)
	actionMap := convertActionObjectToMap(instance.Action)

	subject := JsonSubject{}
	var principals []*JsonPrincipal
	user, ok := subjectMap["user"]
	if ok {
		if strings.HasPrefix(user.(string), "user=") {
			splits := strings.Split(user.(string), "user=")
			if len(splits) > 1 {
				user = splits[1]
			}
		}
		principals = append(principals, &JsonPrincipal{"user", user.(string)})
	}
	groups, ok := subjectMap["groups"]
	if ok {
		groups := groups.([]interface{})
		for _, v := range groups {
			principals = append(principals, &JsonPrincipal{"group", v.(string)})
		}
	}
	sproperties, ok := subjectMap["properties"]
	if ok {
		var sourcenamespace, sourceservice string
		pmap := sproperties.(map[string]interface{})
		for k, v := range pmap {
			if k == "source.namespace" {
				sourcenamespace = v.(string)
			} else if k == "source.service" {
				sourceservice = v.(string)
			}
		}
		if sourcenamespace != "" || sourceservice != "" {
			entity := "service://" + sourcenamespace + "/" + sourceservice
			principals = append(principals, &JsonPrincipal{"entity", entity})
		}
	}
	subject.Principals = principals

	requestContext := JsonContext{}

	requestContext.Subject = &subject
	requestContext.ServiceName = h.speedleService

	resource, ok := actionMap["resource"]
	if ok {
		requestContext.Resource = resource.(string)
	}

	action, ok := actionMap["action"]
	if ok {
		requestContext.Action = action.(string)
	}

	var attributes []*JsonAttribute
	subjectproperties, ok := subjectMap["properties"]
	if ok {
		for k, v := range subjectproperties.(map[string]interface{}) {
			attributes = append(attributes, &JsonAttribute{"subject:" + k, v})
		}
	}

	actionproperties, ok := actionMap["properties"]
	if ok {
		for k, v := range actionproperties.(map[string]interface{}) {
			attributes = append(attributes, &JsonAttribute{"action:" + k, v})
		}
	}
	requestContext.Attributes = attributes

	jsonStr, err := json.Marshal(requestContext)
	if err != nil {
		return adapter.CheckResult{
			Status: status.WithPermissionDenied(fmt.Sprintf("Failed to generate request context: %v", err)),
		}, nil
	}

	fmt.Println(string(jsonStr))

	req, err := http.NewRequest(http.MethodPost, h.adsEndpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return adapter.CheckResult{
			Status: status.WithPermissionDenied(fmt.Sprintf("speedle: request was rejected: %v", err)),
		}, nil
	}

	req.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return adapter.CheckResult{
			Status: status.WithPermissionDenied(fmt.Sprintf("speedle: request was rejected: %v", err)),
		}, nil

	}
	defer resp.Body.Close()

	code := resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if code != 200 {
		return adapter.CheckResult{
			Status: status.WithPermissionDenied(fmt.Sprintf("speedle: request was rejected: %v, %v, %v", h.adsEndpoint, code, string(body))),
		}, nil
	}

	var f map[string]interface{}
	json.Unmarshal(body, &f)
	result := f["allowed"].(bool)
	if !result {
		return adapter.CheckResult{
			Status: status.WithPermissionDenied(fmt.Sprintf("speedle: allowed == false: %v, %v, %v", h.adsEndpoint, code, string(body))),
		}, nil
	}

	return adapter.CheckResult{
		Status: status.OK,
	}, nil
}

func (h *handler) Close() error {
	return nil
}

////////////////// Bootstrap //////////////////////////

// GetInfo returns the Info associated with this adapter implementation.
func GetInfo() adapter.Info {
	return adapter.Info{
		Name:        "speedle",
		Impl:        "istio.io/istio/mixer/adapter/speedle",
		Description: "Istio Authorization with Speedle engine",
		SupportedTemplates: []string{
			authorization.TemplateName,
		},
		DefaultConfig: &config.Params{},
		NewBuilder:    func() adapter.HandlerBuilder { return &builder{} },
	}
}
