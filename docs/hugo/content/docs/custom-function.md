+++
title = "Custom Function"
description = "Write policy condition with your own functions"
date = 2019-01-21T10:43:02+08:00
weight = 600
draft = false
bref = "When the condition in a policy or role policy is too complicated to draft using an expression and built-in functions, you can write your own functions to include in the condition instead"
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
+++

## Why use custom functions?

In Speedle, policies and role policies may include conditions. You can create conditions using expressions. or use the built-in functions that are provided with Speedle. If the condition is very complicated, it may be too difficult to create using expressions and built-in functions. In that case, you can use a custom function that you define, which extends the ability of the authorization engine. Custom functions supplement the built-in functions that Speedle provides, increasing the flexibility for defining conditions.

## How do I use custom functions?

These steps demonstrate how to define a custom function and use it in a condition.

### 1) Implement your function logic and expose the function through the public REST API endpoint

Speedle calls custom functions with a predefined REST API. The REST endpoint of a custom function must be as follows:

#### HTTP Verb: POST

#### HTTP Request Body:

```
type CustomerFunctionRequest struct {
	Params []interface{} `json:"params"`
}
```

#### HTTP Response Body:

```
type CustomerFunctionResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}
```

#### Sample customer function implementation

```
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	funcServerCert = `-----BEGIN CERTIFICATE-----
MIID7TCCAtWgAwIBAgIJALM3l/OZ9uJKMA0GCSqGSIb3DQEBCwUAMIGMMQswCQYD
VQQGEwJjbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzEPMA0G
A1UECgwGb3JhY2xlMQwwCgYDVQQLDANpZG0xEjAQBgNVBAMMCWxvY2FsaG9zdDEm
MCQGCSqGSIb3DQEJARYXY3ludGhpYS5kaW5nQG9yYWNsZS5jb20wHhcNMTgwNDI1
MDc1MDMwWhcNMTkwNDI1MDc1MDMwWjCBjDELMAkGA1UEBhMCY24xEDAOBgNVBAgM
B2JlaWppbmcxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBm9yYWNsZTEMMAoG
A1UECwwDaWRtMRIwEAYDVQQDDAlsb2NhbGhvc3QxJjAkBgkqhkiG9w0BCQEWF2N5
bnRoaWEuZGluZ0BvcmFjbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAn/AFElluGOZYfvlzBGfHfkd/Q9SuQFsSnQt7Qp63Yuf5Ie/q4NACzWPC
B/L6nQrut4OMxJHvhVAswJozRZrQxXvX/vUxkg+TmALj3U9ejF/5arGtjy5v+yGi
wci7zM4r7VNFJGRkfluNRC1kJi4AY6jk6Gl4d/bX4tBXE8mEFY1rUswYtat3OMja
jVAoocClk6WcaQuK9R1uB+BPyxHLJ04RyKRuepPYRBQjgwHK5kMF3s5p07Os+2JH
5jyJYW2NPs6pQe0k8GWpaar/yZ2eut9gsgHnu5JCWnyedo4nEx6I/G4GSaX+0SeU
/Wb2aqq1QGfVOESml7CVcEa/buTeUwIDAQABo1AwTjAdBgNVHQ4EFgQU5i7CO32N
spQ5AaG/aRU0LX2koYwwHwYDVR0jBBgwFoAU5i7CO32NspQ5AaG/aRU0LX2koYww
DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAbQuCMPK8f8QuEmTpZBFv
aka9qruT/0/TrxrbxEh68N4moXSTVv4tSrDTmdkwUiiwayuGS7fvKjSV6hwGkQbV
zGbFDdwOw1tPE2OwnA7/+RPl4KmE4iTHnnIanyg9CKmBW/tMp/vUyv5nIt7Xw5n4
tx3C9/hme+Rlx+SVPIAwAjl4nVFNLfzyG+JDBnQWygySm88SzzK0WRgh5V+gyXCK
ucDW5rA6X9/CM3QrSY50mSM6dbyYDMtmTI4dX7E9STTBCNsNNcmgYkX0N9lm5RoF
uBsAcPmp1SVIbXelDHJiIXxMKzwZy8riZQ8+Dw6LMs6wZX7COVvMWN4Dfcuo89av
IQ==
-----END CERTIFICATE-----`
	funcServerKey = `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQCf8AUSWW4Y5lh+
+XMEZ8d+R39D1K5AWxKdC3tCnrdi5/kh7+rg0ALNY8IH8vqdCu63g4zEke+FUCzA
mjNFmtDFe9f+9TGSD5OYAuPdT16MX/lqsa2PLm/7IaLByLvMzivtU0UkZGR+W41E
LWQmLgBjqOToaXh39tfi0FcTyYQVjWtSzBi1q3c4yNqNUCihwKWTpZxpC4r1HW4H
4E/LEcsnThHIpG56k9hEFCODAcrmQwXezmnTs6z7YkfmPIlhbY0+zqlB7STwZalp
qv/JnZ6632CyAee7kkJafJ52jicTHoj8bgZJpf7RJ5T9ZvZqqrVAZ9U4RKaXsJVw
Rr9u5N5TAgMBAAECggEBAJCavJsojGiq61xyQVG8WxyLnD9B7gJ11VB0bw9+3SPp
xNCwUNbOe5okFeyF/Z07ozX9FKstnzgTk0LYqH7ISPYk0NfN7PG4b6PDCS6xcjTN
GX8kAl4wiEKw2K0IxvOXfRPoc91Bf7LXJ9R6jdAPS37P15ditO8SGYMTB4f2bRvm
AMcaVGFZVNrvxbzC7bw+HWX1k1dezjkGf4jl9M6u252K/gMWZz8KaeDlqj48Kwg/
k6hAnk5h28qCFe0wENOOZySgkrGaIV0hO/PeGydALUAbEL1Fn+20+eblT0IhD1v0
6icQMhaxSTtTgnk9GbHF4L+T2quRAwQ4PAR+9t9Vk6ECgYEA1KZWRzBv4q4u075W
Xu7NORmYJ9HuqL5QdQSb5qsFU7MfXQpK5BP9VDc38qTL5zyD5D40Tpjio7WX6BCp
yWwH+brBmijneZFJ3r/pI8Oyut0m7uX/bbkOmTFrrIwE8lpeOzcUObtUw7uwByWV
lmHP/fYiGMVNitxuTauxVejN6QsCgYEAwIrD3TRC6B9VHRLNiU/F1iuZTPWrCN+j
XXA4ihzlcZFXqBNqlCrN0BDn7cOloob0zLB5T+8hoL5hsIiVg8tn7qDVqnvgql0e
XKcasw0nGv4zI8cpHx7X7V8FxZUyyIAVd3rf33MO8qMTzBrshTh+JSsg5A1eVkjd
9cONGAgLfNkCgYEAidvYPUiqkGNp2j4gEmVwSF9OZCpWNbFDyckWJQGkb3HFmHTO
vnQzHIC71aN+yUdTHgoxsO6up4FXnMwItptBxGWNk5qHDinhoPX7eAMsALbUwbX7
1S9OxoPikTcpEdECHBOGGjNXLZmk8c0s4BRDWhpSWoq2zZpALDxtuAs4SqcCgYEA
m9y5CQwRTU5v3AUolQsan3DTvFTyi1BeMnlxi3ww0GpThx+Qmzi7Or80wGgsYRDW
ggwpZ+ewVStIcVtfjTzPeYCA9m0pRT/0IBS1rFPtYBB+3WuPgj25ldHiHjvUzDHD
LuEs8Pl3FDum/waciItesj/jdDjOMRLzess+IEIC6qECgYEAiA9S52BhIs8R3bYl
5urE2VBVchLLNejgCEyhCv6rE4MibqX4oWwsmLM960rFhw+8UGwPKb0fwwq456f1
vFKFOLz3XCMP6kg1g4dDcB7oRhO+9B4dk+QrUt349VlSaBA/sZ96mEL9ajudec5v
5Amt4mqzY19IT2D8tE+saR2/U4I=
-----END PRIVATE KEY-----`
)

type CustomerFunctionRequest struct {
	Params []interface{} `json:"params"`
}

type CustomerFunctionResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func CustomFunctionIsValid(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request CustomerFunctionRequest
	var response CustomerFunctionResponse
	httpSatus := http.StatusOK

	if err := decoder.Decode(&request); err != nil {
		fmt.Println(err)
		response = CustomerFunctionResponse{
			Error: "error decoding request",
		}
		httpSatus = http.StatusBadRequest
	} else {
		fmt.Printf("request = %v\n", request)
		maxValue := float64(50000)
		minValue := float64(100)
		isValid := true
		for index, param := range request.Params {
			fmt.Printf("param %d: value=%v\n", index, param)
			paramValue := param.(float64)
			if paramValue < minValue || paramValue > maxValue {
				fmt.Println("invalid")
				isValid = false
				break
			}
		}
		response = CustomerFunctionResponse{
			Result: isValid,
		}
	}
	payload, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("repsonse=", string(payload))
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpSatus)
	w.Write(payload)
}

func main() {
	http.HandleFunc("/func/isValid", CustomFunctionIsValid)
	// cp the content of funcServerCert to ./funcServer.crt, cp the content of funcServerKey to ./funcServer.key
	if _, err := os.Stat("./funcServer.crt"); os.IsNotExist(err) {
		if err1 := ioutil.WriteFile("./funcServer.crt", []byte(funcServerCert), 0644); err1 != nil {
			fmt.Println("error creating crt file")
		}
	}
	if _, err := os.Stat("./funcServer.key"); os.IsNotExist(err) {
		if err1 := ioutil.WriteFile("./funcServer.key", []byte(funcServerKey), 0644); err1 != nil {
			fmt.Println("error creating key file")
		}
	}
	err := http.ListenAndServeTLS(":23456", "./funcServer.crt", "./funcServer.key", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

```

### 2) Create the custom function definition in Speedle and associate the function name with the REST endpoint

The definition of a custom function is as follows:

```
type Function struct {
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	FuncURL        string            `json:"funcURL"`                  //endpoint of the function
	CA             string            `json:"ca,omitempty"`             //tls related configurations
	ResultCachable bool              `json:"resultCachable,omitempty"` //false by default
	ResultTTL      int64             `json:"resultTTL,omitempty"`      //TTL of function result in second
}
```

You can create a custom function using either the Speedle Policy Management Service (PMS) REST API, or using the Speedle CLI (`spctl`).

The following example shows how to create a custom function `isValid` using the `spctl` CLI and the JSON file (`function.json`) that contains the custom function definition of `isValid`.

`function.json` file:

```
{
    "name" : "isValid",
    "funcURL" : "https://localhost:23456/func/isValid",
    "resultCachable": true,
    "resultTTL": 300,
    "ca":"-----BEGIN CERTIFICATE-----\nMIID7TCCAtWgAwIBAgIJALM3l/OZ9uJKMA0GCSqGSIb3DQEBCwUAMIGMMQswCQYD\nVQQGEwJjbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzEPMA0G\nA1UECgwGb3JhY2xlMQwwCgYDVQQLDANpZG0xEjAQBgNVBAMMCWxvY2FsaG9zdDEm\nMCQGCSqGSIb3DQEJARYXY3ludGhpYS5kaW5nQG9yYWNsZS5jb20wHhcNMTgwNDI1\nMDc1MDMwWhcNMTkwNDI1MDc1MDMwWjCBjDELMAkGA1UEBhMCY24xEDAOBgNVBAgM\nB2JlaWppbmcxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBm9yYWNsZTEMMAoG\nA1UECwwDaWRtMRIwEAYDVQQDDAlsb2NhbGhvc3QxJjAkBgkqhkiG9w0BCQEWF2N5\nbnRoaWEuZGluZ0BvcmFjbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAn/AFElluGOZYfvlzBGfHfkd/Q9SuQFsSnQt7Qp63Yuf5Ie/q4NACzWPC\nB/L6nQrut4OMxJHvhVAswJozRZrQxXvX/vUxkg+TmALj3U9ejF/5arGtjy5v+yGi\nwci7zM4r7VNFJGRkfluNRC1kJi4AY6jk6Gl4d/bX4tBXE8mEFY1rUswYtat3OMja\njVAoocClk6WcaQuK9R1uB+BPyxHLJ04RyKRuepPYRBQjgwHK5kMF3s5p07Os+2JH\n5jyJYW2NPs6pQe0k8GWpaar/yZ2eut9gsgHnu5JCWnyedo4nEx6I/G4GSaX+0SeU\n/Wb2aqq1QGfVOESml7CVcEa/buTeUwIDAQABo1AwTjAdBgNVHQ4EFgQU5i7CO32N\nspQ5AaG/aRU0LX2koYwwHwYDVR0jBBgwFoAU5i7CO32NspQ5AaG/aRU0LX2koYww\nDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAbQuCMPK8f8QuEmTpZBFv\naka9qruT/0/TrxrbxEh68N4moXSTVv4tSrDTmdkwUiiwayuGS7fvKjSV6hwGkQbV\nzGbFDdwOw1tPE2OwnA7/+RPl4KmE4iTHnnIanyg9CKmBW/tMp/vUyv5nIt7Xw5n4\ntx3C9/hme+Rlx+SVPIAwAjl4nVFNLfzyG+JDBnQWygySm88SzzK0WRgh5V+gyXCK\nucDW5rA6X9/CM3QrSY50mSM6dbyYDMtmTI4dX7E9STTBCNsNNcmgYkX0N9lm5RoF\nuBsAcPmp1SVIbXelDHJiIXxMKzwZy8riZQ8+Dw6LMs6wZX7COVvMWN4Dfcuo89av\nIQ==\n-----END CERTIFICATE-----"
}
```

Use the `spctl create function` command to create the custom function:

```
./spctl create function --json-file=function.json
```

#### 3) Use the custom function in the condition of a policy or role-policy

The following example shows how to create a policy using the `isValid` custom function in the condition.

```
./spctl create policy -c "grant user Ally access library if isValid(attr1)" --service-name=service1
```

**Note:**
You must ensure that the parameters of the function in the condition match the parameters accepted by the function's REST endpoint.
