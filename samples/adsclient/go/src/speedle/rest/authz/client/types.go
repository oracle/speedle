package client

type jsonAttribute struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type jsonPrincipal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

type jsonSubject struct {
	Principals []*jsonPrincipal `json:"principals"`
	Attributes []*jsonAttribute `json:"attributes"`
	TokenType  string           `json:"tokenType"`
	Token      string           `json:"token"`
	Asserted   bool             `json:"asserted"`
}

type jsonContext struct {
	Subject     *jsonSubject     `json:"subject"`
	ServiceName string           `json:"serviceName"`
	Resource    string           `json:"resource"`
	Action      string           `json:"action"`
	Attributes  []*jsonAttribute `json:"attributes"`
}

type jsonIsAllowedResponse struct {
	Allowed      bool   `json:"allowed"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
