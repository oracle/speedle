package authz

type Principal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

type Subject struct {
	Principals []*Principal           `json:"principals,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	TokenType  string                 `json:"tokenType,omitempty"`
	Token      string                 `json:"token,omitempty"`
	Asserted   bool                   `json:"asserted,omitempty"`
}

type RequestContext struct {
	Subject     *Subject               `json:"subject,omitempty"`
	ServiceName string                 `json:"serviceName,omitempty"`
	Resource    string                 `json:"resource,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

const (
	HOST_PROP      = "host"
	ADS_PORT_PROP  = "ads.port"
	IS_SECURE_PROP = "is.secure"
)
