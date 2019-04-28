+++
title = "Token asserter"
description = "Token asserter"
date = 2019-01-21T09:28:30+08:00
weight = 60
draft = false
bref = "Evaluation requests can contain an identity token that represents a user. In this case, the evaluation engine can invoke the asserter service to obtain the explicit identities of the user"
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
tags = ["Identity", "Token asserter"]
+++

## Benefits of the token asserter

An evaluation request can contain an identity token issued by any identity provider as an incoming user identity, instead of specifying the user identities (user identifier and those groups user belongs to) explicitly. When you integrate Speedle into your service, your service does not need to validate and parse identity tokens. Speedle can do it for you.

The Speedle evaluation engine checks whether the incoming request contains an identity token. If yes, then the evaluation engine invokes the token asserter webhook to assert the identity token and obtains the explicit user identities (user identifier and groups). The evaluation engine can then execute the policy evaluation based on the user identifier and groups.

## How to evaluate authorization requests containing identity tokens

To evaluate authorization requests that contain an identity token, follow these steps.

### 1. Implement the webhook interface of the asserter

The asserter service which implements [Token Assertion Plugin API](../api/asserter_api) takes the identity token, the identity provider, and the allowedIDD as inputs and performs token validation and parsing. The service then retrieves explicit identities (user identifier and groups) represented by the identity token.

Note that if the principal's identity domain is set, then the asserted identity may contain the identity domain of the user/group.

#### Sample asserter service

**Note:** This sample asserter service is for testing purposes only.

Sample asserter service source code:

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	// user type of Principal
	PRINCIPAL_TYPE_USER   = "user"
	// group type of Principal
	PRINCIPAL_TYPE_GROUP  = "group"
	// entity type of Principal
	PRINCIPAL_TYPE_ENTITY = "entity"
)

// Principal of Speedle
type Principal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

// AssertResponse assertion response
type AssertResponse struct {
	Principals []*Principal           `json:"principals,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	// non zero indicates errors happen
	ErrCode    int                    `json:"errCode"`
	ErrMessage string                 `json:"errMessage,omitempty"`
}

// SampleAsserter for testing only
type SampleAsserter struct {
}

func (a SampleAsserter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token := r.Header.Get("x-token")
		idp := r.Header.Get("x-idp")

		log.Printf("ServeHTTP, token: %s, idp: %s \n", token, idp)

		subj := AssertResponse{}

		if len(token) == 0 || len(idp) == 0 {
			subj.ErrCode = http.StatusBadRequest
			subj.ErrMessage = "token or idp is empty"
			sendResp(w, http.StatusBadRequest, &subj)
			return
		} else {
			// Parse token and validate token
			// Retrieve groups etc. from token issuer
			// Here we just return a sample result
			subj.ErrCode = 0
			subj.ErrMessage = ""
			subj.Principals = []*Principal{
				&Principal{
					Type: PRINCIPAL_TYPE_USER,
					Name: "user1",
					IDD:  idp,
				},
			}

			sendResp(w, http.StatusOK, &subj)
		}
	}
}

func sendResp(w http.ResponseWriter, status int, data *AssertResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	raw, _ := json.Marshal(data)

	log.Printf("ServeHTTP, asserted subject: %s \n", string(raw))

	w.Write(raw)
}

func main() {

	mux := http.NewServeMux()

	asserter := SampleAsserter{}

	mux.Handle("/v1/assert", asserter)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("start server error: %v", err)
	}

}


```

a. Compile the sample

```bash
go build src/asserter/asserter.go
```

b. Start the sample asserter service

```bash
./asserter
2019/01/25 14:42:54 ServeHTTP, token: test-token, idp: github
2019/01/25 14:42:54 ServeHTTP, asserted subject: {"principals":[{"type":"user","name":"user1","idd":"github"}],"errCode":0}
2019/01/25 14:43:18 ServeHTTP, token: test-token, idp: google
2019/01/25 14:43:18 ServeHTTP, asserted subject: {"principals":[{"type":"user","name":"user1","idd":"google"}],"errCode":0}
2019/01/25 14:43:29 ServeHTTP, token: , idp: google
2019/01/25 14:43:29 ServeHTTP, asserted subject: {"errCode":400,"errMessage":"token or idp is empty"}

```

c. Test the sample

```bash
curl -v -H "x-token:test-token" -H "x-idp:github" http://localhost:8080/v1/assert
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /v1/assert HTTP/1.1
> User-Agent: curl/7.29.0
> Host: localhost:8080
> Accept: */*
> x-token:test-token
> x-idp:github
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 25 Jan 2019 06:42:54 GMT
< Content-Length: 74
<
* Connection #0 to host localhost left intact
{"principals":[{"type":"user","name":"user1","idd":"github"}],"errCode":0}


curl -v -H "x-token:test-token" -H "x-idp:google" http://localhost:8080/v1/assert
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /v1/assert HTTP/1.1
> User-Agent: curl/7.29.0
> Host: localhost:8080
> Accept: */*
> x-token:test-token
> x-idp:google
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 25 Jan 2019 06:43:18 GMT
< Content-Length: 74
<
* Connection #0 to host localhost left intact
{"principals":[{"type":"user","name":"user1","idd":"google"}],"errCode":0}

curl -v -H "x-token:" -H "x-idp:google" http://localhost:8080/v1/assert
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /v1/assert HTTP/1.1
> User-Agent: curl/7.29.0
> Host: localhost:8080
> Accept: */*
> x-idp:google
>
< HTTP/1.1 400 Bad Request
< Content-Type: application/json
< Date: Fri, 25 Jan 2019 06:43:29 GMT
< Content-Length: 50
<
* Connection #0 to host localhost left intact
{"errCode":400,"errMessage":"token or idp is empty"}

```

### 2. Start the asserter service

Start the sample asserter.

```bash
./asserter
```

### 3. Configure the identity asserter webhook

Configure the asserter webhook in the config.json file for the authorization decision service (ADS), and start the ADS service.

Sample config.json:

```json
{
  "storeConfig": {
    "storeType": "file",
    "storeProps": {
      "FileLocation": "./ps.json"
    }
  },
  "enableWatch": true,
  "asserterWebhookConfig": {
    "endpoint": "http://host:port/v1/assert",
    "clientCert": "",
    "clientKey": "",
    "caCert": ""
  },
  "serverConfig": {
    "endpoint": "",
    "insecure": "",
    "certPath": "",
    "clientCertPath": "",
    "keyPath": ""
  },
  "logConfig": {
    "level": "info",
    "formatter": "text",
    "rotationConfig": {
      "filename": ".speedle.log",
      "maxSize": 10,
      "maxBackups": 5,
      "maxAge": 0,
      "LocalTime": false,
      "compress": false
    }
  }
}
```

Update the `asserterWebhookConfig` section of the config.json file to correspond to URI of the asserter service.

```json

"asserterWebhookConfig": {
        "endpoint": "http://localhost:8080/v1/assert",
        "clientCert": "",
        "clientKey": "",
        "caCert": ""
    }

```

In this example:

- `endpoint` - Endpoint of the asserter service
- `clientCert` - Path to the client certificate file, if the asserter service requires two-way SSL verification
- `clientKey` - Path to the client private key file, if the asserter service requires two-way SSL verification
- `caCert` - Path to the asserter service's CA certificate file, if the asserter service is exposed as an HTTPS service

### 4. Create test policies

The following policies are defined on the [identity domain](../idd) page.

```bash

./spctl create service booksvc
# grant user1 coming from github to perform action: read on resource: book
./spctl create policy -c "grant user user1 from github read book" --service-name=booksvc
# grant user1 coming from google to perform action: write on resource: book
./spctl create policy -c "grant user user1 from google write book" --service-name=booksvc
# grant user1 coming from any identity providers to perform action: rent on resource: book
./spctl create policy -c "grant user user1 rent book" --service-name=booksvc

```

### 5. Retrieve identity token from identity provider

This step depends on how your service integrates with the identity provider. If your service integrates with an identity provider that supports the [OpenID Connect](https://openid.net/connect/) or [OAuth ](https://oauth.net/2/) protocols, then your service can always get a valid identity token or access token issued by the identity provider.

For detailed steps for retrieving the identity token, see the documents provided by the identity provider.

### 6. Evaluate the policy with identity tokens issued by different identity providers

The following policy evaluation results are based on the test policies defined in the previous section.

```bash
# The evaluation result is true.
curl -v -k -X POST -d '{ "subject": {"token": "githubtoken", "tokenType":"github"},"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# The evaluation result is false because the identity token was issued by a different identity provider: gitlab
curl -v -k -X POST -d '{ "subject": {"token": "gitlabtoken", "tokenType":"github"},"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# The evaluation result is true.
curl -v -k -X POST -d '{ "subject": {"token": "githubtoken", "tokenType":"github"},"serviceName":"booksvc","resource":"book","action":"rent"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# The evaluation result is false because of different identity domain "idd":"notgoogle"
curl -v -k -X POST -d '{ "subject": {"token": "id token not issued by google", "tokenType":"google"},"serviceName":"booksvc","resource":"book","action":"write"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

```
