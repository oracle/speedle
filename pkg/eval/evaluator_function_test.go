//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/ext"
)

var (
	funcServerCert = `-----BEGIN CERTIFICATE-----\nMIID7TCCAtWgAwIBAgIUUwffVEFk+HL/aJARgvq3RgOWFFcwDQYJKoZIhvcNAQEL\nBQAwgYUxCzAJBgNVBAYTAmNuMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdC\nZWlqaW5nMRwwGgYDVQQKDBNEZWZhdWx0IENvbXBhbnkgTHRkMRIwEAYDVQQDDAls\nb2NhbGhvc3QxIDAeBgkqhkiG9w0BCQEWEWJpbGw4MjhAZ21haWwuY29tMB4XDTE5\nMDgxNDA4NTU1OVoXDTIyMDgxMzA4NTU1OVowgYUxCzAJBgNVBAYTAmNuMRAwDgYD\nVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdCZWlqaW5nMRwwGgYDVQQKDBNEZWZhdWx0\nIENvbXBhbnkgTHRkMRIwEAYDVQQDDAlsb2NhbGhvc3QxIDAeBgkqhkiG9w0BCQEW\nEWJpbGw4MjhAZ21haWwuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC\nAQEAn/AFElluGOZYfvlzBGfHfkd/Q9SuQFsSnQt7Qp63Yuf5Ie/q4NACzWPCB/L6\nnQrut4OMxJHvhVAswJozRZrQxXvX/vUxkg+TmALj3U9ejF/5arGtjy5v+yGiwci7\nzM4r7VNFJGRkfluNRC1kJi4AY6jk6Gl4d/bX4tBXE8mEFY1rUswYtat3OMjajVAo\nocClk6WcaQuK9R1uB+BPyxHLJ04RyKRuepPYRBQjgwHK5kMF3s5p07Os+2JH5jyJ\nYW2NPs6pQe0k8GWpaar/yZ2eut9gsgHnu5JCWnyedo4nEx6I/G4GSaX+0SeU/Wb2\naqq1QGfVOESml7CVcEa/buTeUwIDAQABo1MwUTAdBgNVHQ4EFgQU5i7CO32NspQ5\nAaG/aRU0LX2koYwwHwYDVR0jBBgwFoAU5i7CO32NspQ5AaG/aRU0LX2koYwwDwYD\nVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAhydQLmuoQhEwkrVQ0+KY\naiZv6a6mDTA/WPB47mdp8fr/PRggQdEwn39pPmzXCqG0N3k/EwJQm2xyQ4cpwWE/\nWrHbsE3k8cwfpaBPTuiTWwL++TUBS+/7FzT4wI0a0fpFJZx39UuWStaFR849VDdL\n75o4lIqfTcvCOzeA57VwTpHLTy2HMtx+aOP0FUP0Civfo1HhF9RPjo8zGz/33cYx\nUOUIkfLqJdOaJo6DeAjSLg2LsXn+yzscJ9OQV761F/3100kxWdxB3mG5T5ZPIj5n\n0xWSJUwwKmMWhu29T1/k9gzPs1WnBY2+YvhRGu+93N5CYAlgpXmLBSSFRi938w3z\nyQ==\n-----END CERTIFICATE-----\n`
	funcServerKey  = `-----BEGIN PRIVATE KEY-----\nMIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQCf8AUSWW4Y5lh+\n+XMEZ8d+R39D1K5AWxKdC3tCnrdi5/kh7+rg0ALNY8IH8vqdCu63g4zEke+FUCzA\nmjNFmtDFe9f+9TGSD5OYAuPdT16MX/lqsa2PLm/7IaLByLvMzivtU0UkZGR+W41E\nLWQmLgBjqOToaXh39tfi0FcTyYQVjWtSzBi1q3c4yNqNUCihwKWTpZxpC4r1HW4H\n4E/LEcsnThHIpG56k9hEFCODAcrmQwXezmnTs6z7YkfmPIlhbY0+zqlB7STwZalp\nqv/JnZ6632CyAee7kkJafJ52jicTHoj8bgZJpf7RJ5T9ZvZqqrVAZ9U4RKaXsJVw\nRr9u5N5TAgMBAAECggEBAJCavJsojGiq61xyQVG8WxyLnD9B7gJ11VB0bw9+3SPp\nxNCwUNbOe5okFeyF/Z07ozX9FKstnzgTk0LYqH7ISPYk0NfN7PG4b6PDCS6xcjTN\nGX8kAl4wiEKw2K0IxvOXfRPoc91Bf7LXJ9R6jdAPS37P15ditO8SGYMTB4f2bRvm\nAMcaVGFZVNrvxbzC7bw+HWX1k1dezjkGf4jl9M6u252K/gMWZz8KaeDlqj48Kwg/\nk6hAnk5h28qCFe0wENOOZySgkrGaIV0hO/PeGydALUAbEL1Fn+20+eblT0IhD1v0\n6icQMhaxSTtTgnk9GbHF4L+T2quRAwQ4PAR+9t9Vk6ECgYEA1KZWRzBv4q4u075W\nXu7NORmYJ9HuqL5QdQSb5qsFU7MfXQpK5BP9VDc38qTL5zyD5D40Tpjio7WX6BCp\nyWwH+brBmijneZFJ3r/pI8Oyut0m7uX/bbkOmTFrrIwE8lpeOzcUObtUw7uwByWV\nlmHP/fYiGMVNitxuTauxVejN6QsCgYEAwIrD3TRC6B9VHRLNiU/F1iuZTPWrCN+j\nXXA4ihzlcZFXqBNqlCrN0BDn7cOloob0zLB5T+8hoL5hsIiVg8tn7qDVqnvgql0e\nXKcasw0nGv4zI8cpHx7X7V8FxZUyyIAVd3rf33MO8qMTzBrshTh+JSsg5A1eVkjd\n9cONGAgLfNkCgYEAidvYPUiqkGNp2j4gEmVwSF9OZCpWNbFDyckWJQGkb3HFmHTO\nvnQzHIC71aN+yUdTHgoxsO6up4FXnMwItptBxGWNk5qHDinhoPX7eAMsALbUwbX7\n1S9OxoPikTcpEdECHBOGGjNXLZmk8c0s4BRDWhpSWoq2zZpALDxtuAs4SqcCgYEA\nm9y5CQwRTU5v3AUolQsan3DTvFTyi1BeMnlxi3ww0GpThx+Qmzi7Or80wGgsYRDW\nggwpZ+ewVStIcVtfjTzPeYCA9m0pRT/0IBS1rFPtYBB+3WuPgj25ldHiHjvUzDHD\nLuEs8Pl3FDum/waciItesj/jdDjOMRLzess+IEIC6qECgYEAiA9S52BhIs8R3bYl\n5urE2VBVchLLNejgCEyhCv6rE4MibqX4oWwsmLM960rFhw+8UGwPKb0fwwq456f1\nvFKFOLz3XCMP6kg1g4dDcB7oRhO+9B4dk+QrUt349VlSaBA/sZ96mEL9ajudec5v\n5Amt4mqzY19IT2D8tE+saR2/U4I=\n-----END PRIVATE KEY-----`
)

func startFunctionService() {
	http.HandleFunc("/funcs/testsum", CustomFunctionTestSum)

	go http.ListenAndServe("0.0.0.0:12345", nil)

	//We have an assumption that on speedle/sphinx side, certificate is issued by well known CA.
	/*caCert, err := ioutil.ReadFile("client.crt")
	if err != nil {
		log.Fatal(err)
	}*/
	caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		//ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      "0.0.0.0:23456",
		TLSConfig: tlsConfig,
	}
	server.ListenAndServeTLS("./funcServer.crt", "./funcServer.key")

}

func CustomFunctionTestSum(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request ext.CustomerFunctionRequest
	var response ext.CustomerFunctionResponse
	httpSatus := http.StatusOK

	if err := decoder.Decode(&request); err != nil {
		fmt.Println(err)
		response = ext.CustomerFunctionResponse{
			Error: "error decoding request",
		}
		httpSatus = http.StatusBadRequest
	} else {
		fmt.Printf("request = %v\n", request)
		sum := float64(0)
		for index, param := range request.Params {
			fmt.Printf("param %d: value=%v, type=%t\n", index, param, param)
			sum = sum + param.(float64)
		}
		response = ext.CustomerFunctionResponse{
			Result: sum,
		}
	}
	payload, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
		fmt.Println("repsonse=", string(payload))
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpSatus)
	w.Write(payload)
}

func TestFunctions(t *testing.T) {
	go startFunctionService()

	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "testsum(1,2) <4",
			stream:    `{"functions":[{"name":"testsum","funcURL":"http://localhost:12345/funcs/testsum"}],"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "testsum(1,2) <4"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 7.99}},
			want:      true,
		},
		{
			condition: "testsum1(1,2) <4",
			stream:    `{"functions":[{"name":"testsum1","funcURL":"https://localhost:23456/funcs/testsum", "CA" : "` + funcServerCert + `"}],"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "testsum1(1,2) <4"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 7.99}},
			want:      true,
		},
		{
			condition: "Sqrt(64) > x",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "Sqrt(64) > x"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 7.99}},
			want:      true,
		},
		{
			condition: "Sqrt(64) > x",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "Sqrt(64) > x"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 8.01}},
			want:      false,
		},
		{
			condition: "Sqrt(x) > 7.99",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "Sqrt(x) > 7.99"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 64}},
			want:      true,
		},
		{
			condition: "Sqrt(x) > 8.01",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "Sqrt(x) > 8.01"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 64}},
			want:      false,
		},
		{
			condition: "Max(-3, x, 5) > y",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "Max(-3, x, 5) > y"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 7, "y": 6}},
			want:      true,
		},
		{
			condition: "Max(-3, x, 5) > y",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "Max(-3, x, 5) > y"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"x": 4, "y": 6}},
			want:      false,
		},

		{
			condition: "IsSubSet(s1,s2)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(s1,s2)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s1": []int{1, 3}, "s2": []int{1, 2, 3, 4}}},
			want:      true,
		},
		{
			condition: "IsSubSet(s1,s2)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(s1,s2)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s1": []int{1}, "s2": []int{1, 2, 3, 4}}},
			want:      true,
		},
		{
			condition: "IsSubSet(s1,s2)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(s1,s2)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s1": []int{1, 5}, "s2": []int{1, 2, 3, 4}}},
			want:      false,
		},
		{
			condition: "IsSubSet(s,('BJ','SH','GZ','SZ'))",
			stream:    `{"services": [{"name": "crm","policies": [{"name": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(s,('BJ','SH','GZ','SZ'))"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s": []string{"GZ", "SH"}}},
			want:      true,
		},
		{
			condition: "IsSubSet(s,('BJ','SH','GZ','SZ'))",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(s,('BJ','SH','GZ','SZ'))"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s": []string{"GZ", "TJ"}}},
			want:      false,
		},
		{
			condition: "IsSubSet(('BJ', 'SZ'), s)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(('BJ', 'SZ'), s)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s": []string{"BJ", "GZ", "SH", "SZ"}}},
			want:      true,
		},
		{
			condition: "IsSubSet(('BJ', 'TJ'), s)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet(('BJ', 'TJ'), s)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"s": []string{"BJ", "GZ", "SH", "SZ"}}},
			want:      false,
		},
	}

	for _, tc := range testCases {
		preparePolicyDataInStore([]byte(tc.stream), t)
		eval, err := NewWithStore(conf, testPS)
		if err != nil {
			t.Errorf("error creating evaluator : %v", err)
			continue
		}
		// Run 3 times
		for i := 0; i < 3; i++ {
			got, _, err := eval.IsAllowed(tc.ctx)
			if err != nil {
				t.Errorf("condition: %s, context: %v, error: %v", tc.condition, tc.ctx.Attributes, err)
			}
			if got != tc.want {
				t.Errorf("condition: %s, context: %v, got %v, want %v", tc.condition, tc.ctx.Attributes, got, tc.want)
			}
		}
	}
}
