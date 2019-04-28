+++
title = "令牌用户身份断言"
description = "Token asserter"
date = 2019-01-21T09:28:30+08:00
weight = 60
draft = false
bref = ""
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
tags = ["Identity", "Token asserter"]
+++

## 令牌用户身份断言的好处

评估请求可以包含由任何身份提供者作为传入用户身份发布的身份令牌，而不是明确指定用户身份（用户标识符和用户所属的那些组）。 将 Speedle 集成到服务中时，您的服务无需验证和解析身份令牌。 Speedle 可以为您做到。

Speedle 评估引擎检查传入请求是否包含身份令牌。 如果是，则评估引擎调用令牌断言服务的 webhook 以获得显式用户身份（用户标识符和组）。 然后，评估引擎可以基于用户身份域和组来执行策略评估。

## 如何评估包含身份令牌的授权请求

要评估包含身份令牌的授权请求，请按照下列步骤操作。

### 1. 实现断言服务的 webhook 接口

实现[Token Assertion Plugin API]（../ api / asserter_api）的断言服务将身份令牌，身份提供者和允许的身份域作为输入，并执行令牌验证和解析。 然后，该服务返回由身份令牌表示的显式身份（用户标识符和组）。

请注意，如果设置了用户的身份域，则经过断言服务的用户标识可能包含用户/组的身份域。

#### 示例断言服务

**备注:** 此示例断言器服务仅用于测试目的。

示例断言器服务源代码：

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

a. 编译

```bash

go build src/asserter/asserter.go

```

b. 启动示例断言服务

```bash
./asserter
2019/01/25 14:42:54 ServeHTTP, token: test-token, idp: github
2019/01/25 14:42:54 ServeHTTP, asserted subject: {"principals":[{"type":"user","name":"user1","idd":"github"}],"errCode":0}
2019/01/25 14:43:18 ServeHTTP, token: test-token, idp: google
2019/01/25 14:43:18 ServeHTTP, asserted subject: {"principals":[{"type":"user","name":"user1","idd":"google"}],"errCode":0}
2019/01/25 14:43:29 ServeHTTP, token: , idp: google
2019/01/25 14:43:29 ServeHTTP, asserted subject: {"errCode":400,"errMessage":"token or idp is empty"}

```

c. 测试示例断言服务

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

### 2. 启动断言服务

```bash
./asserter
```

### 3. 配置决策服务使用断言服务

在 config.json 文件中为授权决策服务（ADS）配置断言服务的 webhook，并启动 ADS 服务。

示例配置 config.json:

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

更新 `asserterWebhookConfig` 部分为如下所示.

```json

"asserterWebhookConfig": {
        "endpoint": "http://localhost:8080/v1/assert",
        "clientCert": "",
        "clientKey": "",
        "caCert": ""
    }

```

本例中:

- `endpoint` - 断言服务的 url
- `clientCert` - 如果断言器服务需要双向 SSL 验证，则为客户端证书文件的路径
- `clientKey` - 如果断言器服务需要双向 SSL 验证，则为客户端私钥文件的路径
- `caCert` - 如果断言器服务作为 HTTPS 服务公开，则为断言服务的 CA 证书文件的路径

### 4. 创建测试策略

以下策略在[identity domain]（../ idd）页面上定义。

```bash

./spctl create service booksvc
# grant user1 coming from github to perform action: read on resource: book
./spctl create policy -c "grant user user1 from github read book" --service-name=booksvc
# grant user1 coming from google to perform action: write on resource: book
./spctl create policy -c "grant user user1 from google write book" --service-name=booksvc
# grant user1 coming from any identity providers to perform action: rent on resource: book
./spctl create policy -c "grant user user1 rent book" --service-name=booksvc

```

### 5. 从身份提供者中获取身份令牌

此步骤取决于您的服务如何与身份提供商集成。 如果您的服务与支持[OpenID Connect]（https://openid.net/connect/）或[OAuth]（https://oauth.net/2/）协议的身份提供商集成，那么您的服务始终可以获取身份提供者颁发的有效身份令牌或访问令牌。

有关获取身份令牌的详细步骤，请参阅身份提供者提供的文档。

### 6. 使用由不同身份提供者发布的身份令牌获取策略评估结果

以下策略评估结果基于上一节中定义的测试策略。

```bash
# 评估结果为： 是.
curl -v -k -X POST -d '{ "subject": {"token": "githubtoken", "tokenType":"github"},"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# 评估结果为： 否, 因为评估请求中的用户来自不同的身份域: gitlab
curl -v -k -X POST -d '{ "subject": {"token": "gitlabtoken", "tokenType":"github"},"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# 评估结果为： 是, 因为来自任何身份域的用户 user1 都可以对资源： book执行 rent 操作
curl -v -k -X POST -d '{ "subject": {"token": "githubtoken", "tokenType":"github"},"serviceName":"booksvc","resource":"book","action":"rent"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed


# 评估结果为： 否, 因为评估请求中的用户来自不同的身份域: notgoogle
curl -v -k -X POST -d '{ "subject": {"token": "id token not issued by google", "tokenType":"google"},"serviceName":"booksvc","resource":"book","action":"write"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed
```
