package main

import (
	"fmt"

	"github.com/docker/go-plugins-helpers/authorization"
	"net/url"
	"speedle/api/authz"
	"speedle/rest/authz/client"
	"strings"
)

type speedleDockerPlugin struct {
	serviceName   string
	speedleClient client.ADSClient
}

// newPlugin creates a new Speedle authorization plugin
func newPlugin(connProperties map[string]string, serviceName string) (*speedleDockerPlugin, error) {
	speedleClient, err := client.New(connProperties)
	if err != nil {
		return nil, err
	}
	return &speedleDockerPlugin{
		speedleClient: speedleClient,
		serviceName:   serviceName,
	}, nil
}

// AuthZReq authorizes the docker client command.
// The command is allowed only if it matches a Speedle policy rule.
// Otherwise, the request is denied!
func (plugin *speedleDockerPlugin) AuthZReq(req authorization.Request) authorization.Response {
	fmt.Printf("AuthZReq - Request: %v.\n", req)
	fmt.Printf("User: %s from %s.\n", req.User, req.UserAuthNMethod)
	if req.RequestURI == "/_ping" {
		return authorization.Response{
			Allow: true,
		}
	}
	result, err := plugin.speedleClient.IsAllowed(*plugin.toRequestContext(&req))
	message := ""
	if err != nil {
		message = err.Error()
	}
	return authorization.Response{
		Allow: result,
		Err:   message,
	}
}

// AuthZRes authorizes the docker client response.
// All responses are allowed by default.
func (plugin *speedleDockerPlugin) AuthZRes(req authorization.Request) authorization.Response {
	// Allowed by default.
	fmt.Printf("AuthZRes - Request: %v.\n", req)
	fmt.Printf("User: %s from %s.\n", req.User, req.UserAuthNMethod)
	if req.RequestURI == "/_ping" {
		return authorization.Response{
			Allow: true,
		}
	}
	result, err := plugin.speedleClient.IsAllowed(*plugin.toRequestContext(&req))
	message := ""
	if err != nil {
		message = err.Error()
	}
	return authorization.Response{
		Allow: result,
		Err:   message,
	}
}

func (plugin *speedleDockerPlugin) toRequestContext(req *authorization.Request) *authz.RequestContext {
	user := "root"
	if len(req.User) != 0 {
		user = req.User
	}
	subject := authz.Subject{
		Principals: []*authz.Principal{
			{
				Type: "user",
				Name: user,
			},
		},
	}
	segs := strings.Split(req.RequestURI, "/")
	if len(segs) < 3 {
		fmt.Printf("Unknown request URI %s", req.RequestURI)
	}

	urlObj, _ := url.ParseRequestURI(req.RequestURI)

	attributes := make(map[string]interface{})
	values := urlObj.Query()
	for key, value := range values {
		if len(value) == 1 {
			attributes[key] = value[0]
		} else {
			attributes[key] = value
		}
	}

	fmt.Printf("Subject: %v\n", subject)
	fmt.Printf("Service Name: %s\n", plugin.serviceName)
	fmt.Printf("Resource: %s\n", urlObj.EscapedPath())
	fmt.Printf("Action: %s\n", req.RequestMethod)
	fmt.Printf("Attributes: %v\n", attributes)

	return &authz.RequestContext{
		Subject:     &subject,
		ServiceName: plugin.serviceName,
		Resource:    urlObj.EscapedPath(),
		Action:      req.RequestMethod,
		Attributes:  attributes,
	}
}
