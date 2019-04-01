//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package eval

import (
	"fmt"
	"testing"
	"time"

	adsapi "github.com/oracle/speedle/api/ads"
)

func TestConditions(t *testing.T) {
	subject := adsapi.Subject{
		Principals: []*adsapi.Principal{
			{
				Type: adsapi.PRINCIPAL_TYPE_USER,
				Name: "admin",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "manager",
			},
			{
				Type: adsapi.PRINCIPAL_TYPE_GROUP,
				Name: "tester",
			},
		},
	}
	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true}},
			want:      true,
		},
		{
			condition: "a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false}},
			want:      false,
		},
		{
			condition: "a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 1}},
			want:      false,
		},
		{
			condition: "a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{}},
			want:      false,
		},
		{
			condition: "!a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "!a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false}},
			want:      true,
		},
		{
			condition: "a==false",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a == false"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false}},
			want:      true,
		},
		{
			condition: "attr == true",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a == true"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false}},
			want:      false,
		},
		{
			condition: "a && b",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a && b"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true}},
			want:      false,
		},
		{
			condition: "a || b",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a || b"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true}},
			want:      true,
		},
		{
			condition: "a==1.5",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a==1.5"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 1.5}},
			want:      true,
		},
		{
			condition: "a>b",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a>b"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 35, "b": 34.99}},
			want:      true,
		},
		{
			condition: "a*(b+c)**2==d",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a*(b+c)**2==d"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 2, "b": 1, "c": 2, "d": 18}},
			want:      true,
		},
		{
			condition: "a IN (1,3,5)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a IN (1,3,5)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 3}},
			want:      true,
		},
		{
			condition: "a IN (1,3,5)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a IN (1,3,5)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 2}},
			want:      false,
		},
		{
			condition: "a IN ('BJ', 'SH', 'GZ', 'SZ')",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a IN ('BJ', 'SH', 'GZ', 'SZ')"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "SZ"}},
			want:      true,
		},
		{
			condition: "a IN ('BJ', 'SH', 'GZ', 'SZ')",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a IN ('BJ', 'SH', 'GZ', 'SZ')"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "TJ"}},
			want:      false,
		},
		{
			condition: "a==b",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a==b"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "lol", "b": "lol"}},
			want:      true,
		},
		{
			condition: "a==b",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a==b"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "lol", "b": "LoL"}},
			want:      false,
		},
		{
			condition: "a =~ '.+@oracle.com'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a =~ '.+@oracle.com'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "yufei.yu@oracle.com"}},
			want:      true,
		},
		{
			condition: "a =~ '.+@oracle.com'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a =~ '.+@oracle.com'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "yufei.yu@yahoo.com"}},
			want:      false,
		},
		{
			condition: "a !~ '.+@oracle.com'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a !~ '.+@oracle.com'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "yufei.yu@yahoo.com"}},
			want:      true,
		},
		{
			condition: "a < '2026-11-02'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a < '2026-11-02'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": time.Now().Unix()}},
			want:      true,
		},
		{
			condition: "request_time > '2017-09-04 12:00:00 '",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "request_time > '2017-09-04 12:00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{}},
			want:      true,
		},
		{
			condition: "request_user == 'admin'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "request_user == 'admin'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{}},
			want:      true,
		},
		{
			condition: "request_action =~'^get*?'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_action =~'^get*?'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      true,
		},
		{
			condition: "request_resource =~ '^/node1/*?'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_resource =~ '^/node1/*?'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      true,
		},
		{
			condition: "request_weekday == 'Monday'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_weekday == 'Monday'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      time.Now().Weekday().String() == "Monday",
		},
		{
			condition: "request_year == 2017",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_year==2017"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      time.Now().Year() == 2017,
		},
		{
			condition: "request_month == 11",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_month == 11"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      time.Now().Month() == 11,
		},
		{
			condition: "request_day == 13",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_day == 13"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      time.Now().Day() == 13,
		},
		{
			condition: "request_hour == 14",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "condition": "request_hour == 14"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      time.Now().Hour() == 14,
		},
		{
			condition: "'manager' IN request_groups",
			stream:    `{"services": [{"name": "crm","policies": [{"name": "p1", "effect": "grant", "condition": "'manager' IN request_groups"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: &subject, ServiceName: "crm", Resource: "/node1/a/b/c", Action: "getABC", Attributes: map[string]interface{}{}},
			want:      true,
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
			got, reason, err := eval.IsAllowed(tc.ctx)
			if err != nil && reason != adsapi.ERROR_IN_EVALUATION {
				t.Errorf("condition: %s, context: %v, error: %v", tc.condition, tc.ctx.Attributes, err)
			}
			if got != tc.want {
				t.Errorf("condition: %s, context: %v, got %v, want %v", tc.condition, tc.ctx.Attributes, got, tc.want)
			}
		}
	}
}

func TestStringConditionsPostive(t *testing.T) {
	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "a=='abc'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a=='abc'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "abc"}},
			want:      true,
		},
		{
			condition: "a!='abc'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a!='abc'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "abcd"}},
			want:      true,
		},
		{
			condition: "a>'abc'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a>'abc'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "abcd"}},
			want:      true,
		},
		{
			condition: "a>='abc'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a>='abc'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "abc"}},
			want:      true,
		},
		{
			condition: "a<'abc'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a<'abc'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "ab"}},
			want:      true,
		},
		{
			condition: "a<='abc'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a<='abc'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "abc"}},
			want:      true,
		},
		{
			condition: "a+b=='ab'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a+b=='ab'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "a", "b": "b"}},
			want:      true,
		},
		{
			condition: "a=~'^get.*'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a=~'get.*'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "get_1"}},
			want:      true,
		},
		{
			condition: "!(a=~'^get.*')",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "!(a=~'get.*')"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "got_1"}},
			want:      true,
		},
		{
			condition: "a=~'get.*'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a=~'get.*'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "aget_1"}},
			want:      true,
		},
		{
			condition: "a!~'^delete.*'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a!~'^delete.*'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "adelete"}},
			want:      true,
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

func TestNumericConditionsPostive(t *testing.T) { //float64
	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "a==123",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a==123"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 123}},
			want:      true,
		},
		{
			condition: "a!=123",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a!=123"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 223}},
			want:      true,
		},
		{
			condition: "a>123",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a>123"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 127}},
			want:      true,
		},
		{
			condition: "a>=123",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a>=123"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 123}},
			want:      true,
		},
		{
			condition: "a<123",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a<123"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 122}},
			want:      true,
		},
		{
			condition: "a<=123",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a<=123"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 123}},
			want:      true,
		},
		{
			condition: "a+b==246",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a+b==246"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 123, "b": 123}},
			want:      true,
		},
		{
			condition: "a-b==0",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a-b==0"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 123, "b": 123}},
			want:      true,
		},
		{
			condition: "a*b==8",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a*b==8"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 2, "b": 4}},
			want:      true,
		},
		{
			condition: "a/b==2",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a/b==2"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 8, "b": 4}},
			want:      true,
		},
		{
			condition: "a%b==2",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a%b==2"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 7, "b": 5}},
			want:      true,
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

func TestDatetimeConditionsPostive(t *testing.T) {
	date20170102 := time.Date(2017, 1, 2, 0, 0, 0, 0, time.UTC)
	location, err := time.LoadLocation("Local")
	if err != nil {
		t.Fatal("get localtion failed:", err)
	}
	date20170102Local := time.Date(2017, 1, 2, 0, 0, 0, 0, location)
	fmt.Println(date20170102Local)

	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "a=='2017-01-02T00:00:00-00:00'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a=='2017-01-02T00:00:00-00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": date20170102.Unix()}},
			want:      true,
		},
		{
			condition: "a=='2017-01-02'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a=='2017-01-02'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": date20170102Local.Unix()}},
			want:      true,
		},
		{
			condition: "a!='2017-01-02T00:00:00-00:00'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a!='2017-01-02T00:00:00-00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": time.Now().Unix()}},
			want:      true,
		},
		{
			condition: "a >'2017-01-02T00:00:00-00:00'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a >'2017-01-02T00:00:00-00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": time.Now().Unix()}},
			want:      true,
		},
		{
			condition: "a <'2017-01-03T00:00:00-00:00'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a <'2017-01-03T00:00:00-00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": date20170102.Unix()}},
			want:      true,
		},
		{
			condition: "a >='2017-01-02T00:00:00-00:00'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a >='2017-01-02T00:00:00-00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": time.Now().Unix()}},
			want:      true,
		},
		{
			condition: "a <='2017-01-03T00:00:00-00:00'",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a <='2017-01-03T00:00:00-00:00'"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": date20170102.Unix()}},
			want:      true,
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

func TestBoolConditionsPostive(t *testing.T) {
	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "a==true",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a==true"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true}},
			want:      true,
		},
		{
			condition: "a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true}},
			want:      true,
		},
		{
			condition: "!a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "!a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false}},
			want:      true,
		},
		{
			condition: "a||b||c",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a||b||c"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false, "b": false, "c": true}},
			want:      true,
		},
		{
			condition: "a&&b&&c",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a&&b&&c"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true, "b": true, "c": true}},
			want:      true,
		},
		{
			condition: "a||b&&c",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a||b&&c"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": false, "b": true, "c": true}},
			want:      true,
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

func TestArrayConditionsPostive(t *testing.T) { //float64, string, datetime
	date20170102 := time.Date(2017, 1, 2, 0, 0, 0, 0, time.UTC)
	date20180102 := time.Date(2018, 1, 2, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{
		{
			condition: "a in (1, 2, 3)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in (1, 2, 3)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 1}},
			want:      true,
		},

		{
			condition: "1 in a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "1 in a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": []float64{1, 2, 3}}},
			want:      true,
		},

		{
			condition: "a in ('a','b','c')",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in ('a','b','c')"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": "a"}},
			want:      true,
		},
		{
			condition: "'a' in a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "'a' in a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": []string{"a", "b", "c"}}},
			want:      true,
		},
		{
			condition: "a in (true, false)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in (true,false)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": true}},
			want:      true,
		},
		{
			condition: "true in a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "true in a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": []bool{true, false}}},
			want:      true,
		},
		//TODO:
		{
			condition: "'2017-01-02T00:00:00-00:00' in a",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "'2017-01-02T00:00:00-00:00' in a"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": []float64{float64(date20170102.Unix()), float64(date20180102.Unix())}}},
			want:      true,
		},
		{
			condition: "a in ('2017-01-02T00:00:00-00:00', '2018-01-02T00:00:00-00:00')",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in ('2017-01-02T00:00:00-00:00', '2018-01-02T00:00:00-00:00')"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": date20170102.Unix()}},
			want:      true,
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

func TestComplexConditionsPostive(t *testing.T) {
	testCases := []struct {
		condition string
		stream    string
		ctx       adsapi.RequestContext
		want      bool
	}{

		{
			condition: "a in (1, 2, 3) && (b==c || d==3)",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in (1, 2, 3)"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 1, "b": "string1", "c": "string2", "d": 3}},
			want:      true,
		},
		{
			condition: "a in (1,2,3) && (b==c || d==3) && 's1' in e",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in (1,2,3) && (b==c || d==3) || 's1' in e"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 2, "b": "string1", "c": "string2", "d": 3, "e": []string{"s1", "s2"}}},
			want:      true,
		},
		{
			condition: "a in (1,2,3) && (b==c || d==3) && IsSubSet(e, ('s1','s2','s3'))",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "a in (1,2,3) && (b==c || d==3) || IsSubSet(e, ('s1','s2','s3'))"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": 2, "b": "string1", "c": "string2", "d": 3, "e": []string{"s1", "s2"}}},
			want:      true,
		},
		{
			condition: "IsSubSet((1,3),a) && (b==c || d==3) && IsSubSet(e, ('s1','s2','s3'))",
			stream:    `{"services": [{"name": "crm","policies": [{"id": "p1", "effect": "grant", "permissions": [{"resource": "/node1","actions": ["get"]}],"condition": "IsSubSet((1,3),a) && (b==c || d==3) && IsSubSet(e, ('s1','s2','s3'))"}]}]}`,
			ctx:       adsapi.RequestContext{Subject: nil, ServiceName: "crm", Resource: "/node1", Action: "get", Attributes: map[string]interface{}{"a": []float64{1, 2, 3, 4}, "b": "string1", "c": "string2", "d": 3, "e": []string{"s1", "s2"}}},
			want:      true,
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
