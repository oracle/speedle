//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/store"
)

func TestSaveGetLastRequest(t *testing.T) {
	s, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new file store:", err)
	}
	store.MaxDiscoverRequestNum = int64(100)
	store.DeleteNumWhenReachMaxDiscoverRequest = int64(10)
	discover := s.(store.DiscoverRequestManager)

	i := 0
	for i < 100 {
		user := ads.Principal{Type: "user", Name: "user" + strconv.Itoa(i%10)}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		serviceName := "erp" + strconv.Itoa(i%10)
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request in store", err)
		}

		req, _, err := discover.GetLastDiscoverRequest(serviceName)
		if err != nil {
			t.Error("fail to get last request for service")
		}
		if req.Resource != resName {
			t.Error("the last request for service is incorrect")
		}

		i++
	}

}

func TestGetLastRequestContinously(t *testing.T) {
	s, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd store")
	}
	discover := s.(store.DiscoverRequestManager)
	i := 0
	for i < 5 {
		user := ads.Principal{Type: "user", Name: "user" + strconv.Itoa(i%10)}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		serviceName := "erp"
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Fatal("fail to put request in store")
		}
		i++
	}
	request, revision, err := discover.GetLastDiscoverRequest("erp")
	if err != nil {
		t.Errorf("fail to GetLastDiscoverRequest:%v", err)
	}
	if request.Resource != "/res4" {
		t.Error("last request is incorrect.")
	}
	for i < 10 {
		user := ads.Principal{Type: "user", Name: "user" + strconv.Itoa(i%10)}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		serviceName := "erp"
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Fatal("fail to put request in store")
		}
		i++
	}
	requests, _, err := discover.GetDiscoverRequestsSinceRevision("erp", revision)
	if err != nil {
		t.Errorf("fail to GetDiscoverRequestsSinceRevision:%v", err)
	}
	if len(requests) != 5 {
		t.Error("requests number should be 5")
	}
	if requests[0].Resource != "/res5" ||
		requests[1].Resource != "/res6" ||
		requests[2].Resource != "/res7" ||
		requests[3].Resource != "/res8" ||
		requests[4].Resource != "/res9" {
		t.Error("requests sequence or content is incorrect")
	}
}

func TestResetDiscoverRequests(t *testing.T) {
	s, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd store")
	}
	discover := s.(store.DiscoverRequestManager)
	i := 0
	for i < 100 {
		user := ads.Principal{Type: "user", Name: "user" + strconv.Itoa(i%10)}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		serviceName := "erp" + strconv.Itoa(i%10)
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request.")
		}
		i++
	}
	err = discover.ResetDiscoverRequests("erp0")
	if err != nil {
		t.Fatal("Fail to reset service requests")
	}
	zeroRequests, _, err := discover.GetDiscoverRequests("erp0")
	if err != nil {
		t.Fatal("Fail to GetDiscoverRequests")
	}
	if len(zeroRequests) != 0 {
		t.Fatal("Should have no requests now, as requests are reset")
	}

	err = discover.ResetDiscoverRequests("")
	if err != nil {
		t.Fatal("Fail to reset all requests")
	}
	zeroRequests, _, _ = discover.GetDiscoverRequests("")
	if len(zeroRequests) != 0 {
		t.Fatal("Should have no requests now, as requests are reset")
	}

}

func TestGetRequests(t *testing.T) {
	s, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd store")
	}
	store.MaxDiscoverRequestNum = int64(1000)
	store.DeleteNumWhenReachMaxDiscoverRequest = int64(100)
	discover := s.(store.DiscoverRequestManager)
	err = discover.ResetDiscoverRequests("")
	if err != nil {
		t.Fatal("Fail to reset all requests")
	}

	requestNum := 220
	i := 0
	for i < requestNum {
		user := ads.Principal{Type: "user", Name: "user" + strconv.Itoa(i%10)}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		serviceName := "erp" + strconv.Itoa(i%10)
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request.")
		}
		i++
	}
	fmt.Println("starting get requests, time is:", time.Now())
	requests, _, err := discover.GetDiscoverRequests("erp0")
	fmt.Println("finish get requests, time is:", time.Now())
	if err != nil {
		t.Error("fail to GetDiscoverRequests.")
	}

	if len(requests) != 22 {
		t.Error("number incorrect, expected is:", requestNum, ",but is:", len(requests))
	}
	for index, req := range requests {
		if req.Resource != "/res"+strconv.Itoa(index*10) {
			t.Error("sequence is incorrect,", "expected is:", "/res"+strconv.Itoa(index), "but is:", req.Resource)
		}
	}
}

func TestGeneratePolicies(t *testing.T) {
	s, err := store.NewStore("file", storeConfig)
	if err != nil {
		t.Fatal("fail to new etcd store")
	}
	discover := s.(store.DiscoverRequestManager)
	err = discover.ResetDiscoverRequests("")
	if err != nil {
		t.Fatal("Fail to reset all requests")
	}
	serviceName := "erp"
	i := 0
	for i < 5 {
		user := ads.Principal{Type: "user", Name: "user1"}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request.")
		}
		i++
	}
	for i < 10 {
		user := ads.Principal{Type: "user", Name: "user2"}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "write", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request.")
		}
		i++
	}
	serviceName = "erp1"
	i = 0
	for i < 5 {
		user := ads.Principal{Type: "user", Name: "user1"}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "write", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request.")
		}
		i++
	}
	for i < 10 {
		user := ads.Principal{Type: "user", Name: "user2"}
		subj := ads.Subject{Principals: []*ads.Principal{&user}}
		resName := "/res" + strconv.Itoa(i)
		request := ads.RequestContext{Subject: &subj, ServiceName: serviceName, Resource: resName, Action: "read", Attributes: map[string]interface{}{}}
		err := discover.SaveDiscoverRequest(&request)
		if err != nil {
			t.Error("fail to put request.")
		}
		i++
	}
	serviceMap, _, err := discover.GeneratePolicies("erp", "", "user1", "")
	if err != nil {
		t.Error("fail to generate policy for user:", err)
	}
	if serviceMap["erp"] == nil {
		t.Error("service should exist.")
	} else {
		if len(serviceMap["erp"].Policies) != 5 {
			t.Error("should have generated 5 policies")
		} else {
			for index, policy := range serviceMap["erp"].Policies {
				if (policy.Permissions[0].Resource != "/res"+strconv.Itoa(index)) && (policy.Permissions[0].Actions[0] != "read") {
					t.Error("policy error")
				}
			}
		}
	}
	serviceMap, _, err = discover.GeneratePolicies("erp1", "", "user1", "")
	if err != nil {
		t.Error("fail to generate policy for user:", err)
	}
	if serviceMap["erp1"] == nil {
		t.Error("service should exist")
	} else {
		if len(serviceMap["erp1"].Policies) != 5 {
			t.Error("should have generated 5 policies")
		} else {
			for index, policy := range serviceMap["erp1"].Policies {
				if (policy.Permissions[0].Resource != "/res"+strconv.Itoa(index)) && (policy.Permissions[0].Actions[0] != "write") {
					t.Error("policy error")
				}
			}
		}
	}
	serviceMap, _, err = discover.GeneratePolicies("erp", "", "user2", "")
	if err != nil {
		t.Error("fail to generate policy for user:", err)
	}
	if serviceMap["erp"] == nil {
		t.Error("policy should not be nil")
	} else {
		if len(serviceMap["erp"].Policies) != 5 {
			t.Error("should have generated 5 policies")
		} else {
			for index, policy := range serviceMap["erp"].Policies {
				if (policy.Permissions[0].Resource != "/res"+strconv.Itoa(5+index)) && (policy.Permissions[0].Actions[0] != "write") {
					t.Error("policy error")
				}
			}
		}
	}
	serviceMap, _, err = discover.GeneratePolicies("erp1", "", "user2", "")
	if err != nil {
		t.Error("fail to generate policy for user:", err)
	}
	if serviceMap["erp1"] == nil {
		t.Error("policy should not be nil")
	} else {
		if len(serviceMap["erp1"].Policies) != 5 {
			t.Error("should have generated 5 policies")
		} else {
			for index, policy := range serviceMap["erp1"].Policies {
				if (policy.Permissions[0].Resource != "/res"+strconv.Itoa(5+index)) && (policy.Permissions[0].Actions[0] != "read") {
					t.Error("policy error")
				}
			}
		}
	}

	serviceMap, _, err = discover.GeneratePolicies("erp", "", "", "")
	if err != nil {
		t.Error("fail to generate policies:", err)
	} else {
		if len(serviceMap["erp"].RolePolicies) != 2 {
			t.Error("should have generated 2 role policies")
		}

		if len(serviceMap["erp"].Policies) != 10 {
			t.Error("should have generated 10 policies")
		}
	}

	serviceMap, _, err = discover.GeneratePolicies("erp1", "", "", "")
	if err != nil {
		t.Error("fail to generate policies:", err)
	} else {
		if len(serviceMap["erp1"].RolePolicies) != 2 {
			t.Error("should have generated 2 role policies")
		}

		if len(serviceMap["erp1"].Policies) != 10 {
			t.Error("should have generated 10 policies")
		}
	}

}
