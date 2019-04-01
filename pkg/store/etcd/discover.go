//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package etcd

import (
	"encoding/json"
	"fmt"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/store"
	"golang.org/x/net/context"

	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	DiscoverPrefix  = "/isAllowedRequests/"
	DefaultPageSize = int64(1000)
)

func (s *Store) SaveDiscoverRequest(request *ads.RequestContext) error {
	_, err := s.PutRequest(request)
	return err
}

func (s *Store) PutRequest(request *ads.RequestContext) (int64, error) {
	value, err := json.Marshal(request)
	if err != nil {
		return -1, errors.Wrap(err, errors.SerializationError, "failed to marshal request")
	}
	succeed := false
	for !succeed {
		key := DiscoverPrefix + request.ServiceName + "/" + time.Now().String()
		txnResp, err := s.client.KV.Txn(context.TODO()).If(
			clientv3.Compare(clientv3.CreateRevision(key), "=", 0), //key does not exist
		).Then(
			clientv3.OpPut(key, string(value)),
			clientv3.OpGet(DiscoverPrefix, clientv3.WithPrefix(), clientv3.WithCountOnly()), //get number of requests
			clientv3.OpGet(DiscoverPrefix, clientv3.WithLimit(store.DeleteNumWhenReachMaxDiscoverRequest), clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend)), //get oldest keys
		).Commit()
		if err != nil {
			return -1, err
		}
		if txnResp.Succeeded { //if not succeed, the key already exist, try again
			succeed = true
			count := txnResp.Responses[1].GetResponseRange().Count
			if count >= store.MaxDiscoverRequestNum { //reach Max number of requests, remove the oldest ones.
				keys := []string{}
				for _, kv := range txnResp.Responses[2].GetResponseRange().Kvs {
					keys = append(keys, string(kv.Key))
				}
				go s.DeleteRequests(keys)
			}
			return count, nil
		}
		fmt.Println("key already exist, try with a new key...")
	}
	return -1, nil //should not go here
}

func (s *Store) DeleteRequests(keys []string) error {
	deleteOps := []clientv3.Op{}
	for _, key := range keys {
		deleteOps = append(deleteOps, clientv3.OpDelete(key))
	}
	_, err := s.client.KV.Txn(context.TODO()).Then(deleteOps...).Commit()
	if err != nil {
		return errors.Wrapf(err, errors.StoreError, "unable to delete all discover requests %v", keys)
	}

	return nil
}

func (s *Store) GetLastDiscoverRequest(serviceName string) (*ads.RequestContext, int64, error) {
	getOpts := append(clientv3.WithLastCreate(), clientv3.WithPrefix())
	keyPrefix4Search := DiscoverPrefix
	if len(serviceName) > 0 {
		keyPrefix4Search = keyPrefix4Search + serviceName + KeySeparator
	}
	getResp, err := s.client.Get(context.TODO(), keyPrefix4Search, getOpts...)
	if err != nil {
		return nil, -1, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, -1, errors.Wrapf(err, errors.EntityNotFound, "no request found for service %q", serviceName)
	}
	var request ads.RequestContext
	err = json.Unmarshal(getResp.Kvs[0].Value, &request)
	if err != nil {
		return nil, -1, errors.Wrapf(err, errors.SerializationError, "failed to unmarshal request context %q for service %q", getResp.Kvs[0].Value, serviceName)
	}
	return &request, getResp.Header.Revision, nil
}

func (s *Store) GetDiscoverRequestsSinceRevision(serviceName string, revision int64) ([]*ads.RequestContext, int64, error) {
	getOpts := []clientv3.OpOption{clientv3.WithMinCreateRev(revision + 1), clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend)}
	keyPrefix4Search := DiscoverPrefix
	if len(serviceName) > 0 {
		keyPrefix4Search = keyPrefix4Search + serviceName + KeySeparator
	}
	getResp, err := s.client.Get(context.TODO(), keyPrefix4Search, getOpts...)
	if err != nil {
		return nil, revision, errors.Wrapf(err, errors.StoreError, "unable to get discover request for service %q with revision %d", serviceName, revision)
	}
	requests := []*ads.RequestContext{}
	for _, kv := range getResp.Kvs {
		var req ads.RequestContext
		err = json.Unmarshal(kv.Value, &req)
		if err != nil {
			return nil, -1, errors.Wrapf(err, errors.SerializationError, "failed to unmarshal request context %q for service %q", kv.Value, serviceName)
		}
		requests = append(requests, &req)
	}
	return requests, getResp.Header.Revision, nil

}

func (s *Store) GetRequests(keyPrefix string, pageSize int64) ([]*ads.RequestContext, int64, error) {
	requests := []*ads.RequestContext{}
	getOpts := []clientv3.OpOption{clientv3.WithPrefix(), clientv3.WithLimit(pageSize), clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend)}
	var revision int64
	for {
		getResp, err := s.client.Get(context.TODO(), keyPrefix, getOpts...)
		if err != nil {
			return nil, -1, errors.Wrapf(err, errors.StoreError, "unable to get discover requests from etcd server for prefix %q", keyPrefix)
		}
		for _, kv := range getResp.Kvs {
			var req ads.RequestContext
			err = json.Unmarshal(kv.Value, &req)
			if err != nil {
				return nil, -1, errors.Wrapf(err, errors.SerializationError, "failed to unmarshal request context %q for prefix %q", kv.Value, keyPrefix)
			}
			requests = append(requests, &req)
		}
		fmt.Println("len=", len(getResp.Kvs), "more:", getResp.More, "revision:", getResp.Header.Revision)
		if getResp.More {
			revision := getResp.Kvs[pageSize-1].CreateRevision
			getOpts = []clientv3.OpOption{clientv3.WithPrefix(), clientv3.WithMinCreateRev(revision + 1), clientv3.WithLimit(pageSize), clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend)}
		} else {
			revision = getResp.Header.Revision
			break
		}
	}
	return requests, revision, nil

}

func (s *Store) GetDiscoverRequests(serviceName string) ([]*ads.RequestContext, int64, error) {
	if len(serviceName) == 0 {
		return s.GetRequests(DiscoverPrefix, DefaultPageSize)
	} else {
		return s.GetRequests(DiscoverPrefix+serviceName+KeySeparator, DefaultPageSize)
	}

}

func (s *Store) ResetDiscoverRequests(serviceName string) error {
	var err error
	if len(serviceName) == 0 {
		_, err = s.client.Delete(context.TODO(), DiscoverPrefix, clientv3.WithPrefix())
	} else {
		_, err = s.client.Delete(context.TODO(), DiscoverPrefix+serviceName+KeySeparator, clientv3.WithPrefix())
	}
	if err != nil {
		return errors.Errorf(errors.StoreError, "unable to reset discover requests from service %q", serviceName)
	}
	return nil
}

//This method is implemented as common method at evaluator part
func (d *Store) GeneratePolicies(serviceName, principalType, principalName, principalIDD string) (map[string]*pms.Service, int64, error) {
	requests, revision, err := d.GetDiscoverRequests(serviceName)
	if err != nil {
		return nil, -1, err
	}
	serviceMap, err := store.GeneratePoliciesFromDiscoverRequests(requests, principalType, principalName, principalIDD)
	if err != nil {
		return nil, -1, err
	}
	return serviceMap, revision, nil
}
