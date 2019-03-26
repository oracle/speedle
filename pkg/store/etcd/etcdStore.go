//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package etcd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/suid"

	"github.com/oracle/speedle/api/pms"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/embed"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

const (
	requestTimeout  = 10 * time.Second
	KeySeparator    = "/"
	PoliciesKey     = "policies"
	RolePoliciesKey = "role_policies"
	ServicesKey     = "services"
	FunctionsKey    = "functions"
	ServiceTypeKey  = "type"
	pageSize        = 1000
)

type Store struct {
	client       *clientv3.Client
	Config       *clientv3.Config
	KeyPrefix    string
	stop         chan struct{}
	embeddedInst *embed.Etcd
	embeddedDir  string
}

func (s *Store) destroy() error {
	err := s.client.Close()
	if s.embeddedInst != nil {
		CleanEmbeddedEtcd(s.embeddedInst, s.embeddedDir)
	}
	if err != nil {
		return errors.New(errors.StoreError, "unable to close connection to etcd server")
	}
	return nil
}

//read policy store from etcd3
func (s *Store) ReadPolicyStore() (*pms.PolicyStore, error) {
	serviceNames, err := s.GetServiceNames()
	if err != nil {
		return nil, err
	}
	var ps pms.PolicyStore
	for _, serviceName := range serviceNames {
		service, err := s.GetService(serviceName)
		if err != nil {
			return nil, err
		}
		ps.Services = append(ps.Services, service)
	}
	functions, err := s.ListAllFunctions("")
	if err != nil {
		return nil, err
	}
	ps.Functions = functions
	return &ps, nil
}

//write policy store to etcd3
func (s *Store) WritePolicyStore(ps *pms.PolicyStore) error {
	err := s.DeleteServices()
	if err != nil {
		return err
	}
	for _, service := range ps.Services {
		err := s.CreateService(service)
		if err != nil {
			return err
		}
	}
	for _, function := range ps.Functions {
		_, err := s.CreateFunction(function)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetServiceNames() (serviceNames []string, err error) {
	serviceKeyPrefix := s.KeyPrefix + ServicesKey + KeySeparator
	responses, err := s.prefixGet(serviceKeyPrefix, clientv3.WithKeysOnly())
	if err != nil {
		return serviceNames, err
	}
	for _, resp := range responses {
		if resp == nil {
			continue
		}
		for _, kv := range resp.Kvs {
			key := strings.TrimPrefix(string(kv.Key), serviceKeyPrefix)
			key = strings.TrimSuffix(key, KeySeparator)
			if len(key) > 0 && strings.Index(key, KeySeparator) == -1 {
				serviceNames = append(serviceNames, key)
			}
		}
	}
	return serviceNames, nil
}

func (s *Store) GetPolicyAndRolePolicyCounts() (map[string]*pms.PolicyAndRolePolicyCount, error) {
	serviceNames, err := s.GetServiceNames()
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]*pms.PolicyAndRolePolicyCount)

	for _, srvName := range serviceNames {
		var counts pms.PolicyAndRolePolicyCount

		policyCount, err := s.GetPolicyCount(srvName)
		if err != nil {
			return nil, err
		}
		counts.PolicyCount = policyCount

		rolePolicyCount, err := s.GetRolePolicyCount(srvName)
		if err != nil {
			return nil, err
		}
		counts.RolePolicyCount = rolePolicyCount

		countMap[srvName] = &counts
	}

	return countMap, nil
}

func (s *Store) GetServiceCount() (int64, error) {
	serviceNames, err := s.GetServiceNames()
	if err != nil {
		return 0, err
	}

	if nil == serviceNames {
		return 0, nil
	} else {
		return int64(len(serviceNames)), nil
	}
}

func (s *Store) ListAllServices() (services []*pms.Service, err error) {
	serviceNames, err := s.GetServiceNames()
	if err != nil {
		return nil, err
	}
	for _, serviceName := range serviceNames {
		service, err := s.GetService(serviceName)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

//TODO: to be implemented
func (s *Store) GetServices(startName string, amount int, retrivePolcies bool) ([]*pms.Service, string, error) {
	var services []*pms.Service
	serviceNames, err := s.GetServiceNames()
	if err != nil {
		return nil, "", err
	}
	if amount < 0 { //set to retrieve all services when amount is less than 0.
		amount = len(serviceNames) + 1
	}
	nextService := ""
	foundStart := false
	for _, serviceName := range serviceNames {
		if serviceName >= startName {
			foundStart = true
		}
		if foundStart {
			if amount == 0 {
				nextService = serviceName
				break
			} else {
				var service *pms.Service
				if retrivePolcies {
					service, err = s.GetService(serviceName)
				} else {
					service, err = s.GetServiceItself(serviceName)
				}
				if err != nil {
					return nil, "", err
				}
				services = append(services, service)
				amount--
			}
		}
	}
	return services, nextService, nil

}
func (s *Store) GetServiceItself(serviceName string) (*pms.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	serviceKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName
	getOpts := []clientv3.OpOption{clientv3.WithKeysOnly(), clientv3.WithPrefix(), clientv3.WithLimit(int64(1))}
	resp, err := s.client.Get(ctx, serviceKey, getOpts...)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.Errorf(errors.EntityNotFound, "service %q is not found", serviceName)
	}
	service := pms.Service{Name: serviceName}

	resp, err = s.client.Get(ctx, serviceKey+KeySeparator+ServiceTypeKey)
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		service.Type = string(kv.Value)
	}

	return &service, nil
}

func (s *Store) GetService(serviceName string) (*pms.Service, error) {
	var service pms.Service
	serviceKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator
	responses, err := s.prefixGet(serviceKey)
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 || len(responses[0].Kvs) == 0 {
		return nil, errors.Errorf(errors.EntityNotFound, "service %q is not found", serviceName)
	}
	service.Name = serviceName
	for _, resp := range responses {
		for _, kv := range resp.Kvs {
			if strings.Compare(string(kv.Key), serviceKey+ServiceTypeKey) == 0 {
				//service type
				service.Type = string(kv.Value)
			}
			if strings.HasPrefix(string(kv.Key), serviceKey+PoliciesKey) {
				//policies
				var policy pms.Policy
				err := json.Unmarshal(kv.Value, &policy)
				if err != nil {
					return nil, errors.Errorf(errors.SerializationError, "failed to unmarshal policy %q", kv.Value)
				}
				service.Policies = append(service.Policies, &policy)
			}
			if strings.HasPrefix(string(kv.Key), serviceKey+RolePoliciesKey) {
				//role policies
				var rolePolicy pms.RolePolicy
				err := json.Unmarshal(kv.Value, &rolePolicy)
				if err != nil {
					return nil, errors.Errorf(errors.SerializationError, "failed to unmarshal role policy %q", kv.Value)
				}
				service.RolePolicies = append(service.RolePolicies, &rolePolicy)
			}
		}
	}
	return &service, nil
}

func (s *Store) timeOutGet(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.client.Get(ctx, key, opts...)
}

func (s *Store) prefixGet(prefix string, opts ...clientv3.OpOption) ([]*clientv3.GetResponse, error) {
	end := clientv3.GetPrefixRangeEnd(prefix)
	getOpts := []clientv3.OpOption{clientv3.WithPrefix(), clientv3.WithLimit(pageSize), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)}
	getOpts = append(getOpts, opts...)
	ret := []*clientv3.GetResponse{}
	for {
		getResp, err := s.timeOutGet(prefix, getOpts...)
		if err != nil {
			return nil, err
		}
		ret = append(ret, getResp)
		if getResp.More {
			lastKey := string(getResp.Kvs[pageSize-1].Key)
			prefix = clientv3.GetPrefixRangeEnd(lastKey)
			getOpts = []clientv3.OpOption{clientv3.WithRange(end), clientv3.WithLimit(pageSize), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)}
			getOpts = append(getOpts, opts...)
		} else {
			break
		}
	}

	return ret, nil
}

func (s *Store) getPutOps(service *pms.Service) ([]clientv3.Op, error) {
	var ops []clientv3.Op
	for _, policy := range service.Policies {
		if policy.ID == "" {
			policy.ID = suid.New().String()
		}
		key := s.KeyPrefix + ServicesKey + KeySeparator + service.Name + KeySeparator + PoliciesKey + KeySeparator + policy.ID
		value, err := json.Marshal(policy)
		if err != nil {
			return nil, errors.Errorf(errors.SerializationError, "failed to marshal policy")
		}
		ops = append(ops, clientv3.OpPut(key, string(value)))
	}
	for _, rolePolicy := range service.RolePolicies {
		if rolePolicy.ID == "" {
			rolePolicy.ID = suid.New().String()
		}
		key := s.KeyPrefix + ServicesKey + KeySeparator + service.Name + KeySeparator + RolePoliciesKey + KeySeparator + rolePolicy.ID
		value, err := json.Marshal(rolePolicy)
		if err != nil {
			return nil, errors.Errorf(errors.SerializationError, "failed to marshal role policy")
		}
		ops = append(ops, clientv3.OpPut(key, string(value)))
	}
	ops = append(ops, clientv3.OpPut(s.KeyPrefix+ServicesKey+KeySeparator+service.Name+KeySeparator+ServiceTypeKey, service.Type))
	//make sure updating service key is the last operation, so watch could work correctly
	ops = append(ops, clientv3.OpPut(s.KeyPrefix+ServicesKey+KeySeparator+service.Name+KeySeparator, ""))
	return ops, nil

}

func (s *Store) CreateService(service *pms.Service) error {
	ops, err := s.getPutOps(service)
	if err != nil {
		return err
	}
	//currently etcd transaction only support up to 128 operations in one transaction.
	//https://github.com/coreos/etcd/issues/7826, it seems the MaxOpsPerTxn is configurable in later release.
	maxOps := int(embed.DefaultMaxTxnOps)
	startIndex := 0
	var endIndex int
	fail := false
	for startIndex < len(ops) {
		if startIndex+maxOps < len(ops) {
			endIndex = startIndex + maxOps
		} else {
			endIndex = len(ops)
		}
		txnResp, err := s.client.KV.Txn(context.TODO()).If(
			clientv3.Compare(clientv3.Version(s.KeyPrefix+ServicesKey+KeySeparator+service.Name+KeySeparator), "=", 0), //service key does not exist
		).Then(
			ops[startIndex:endIndex]...,
		).Commit()
		if err != nil {
			fail = true
			break
		}
		if !txnResp.Succeeded {
			return errors.Errorf(errors.EntityAlreadyExists, "service %q already exists", service.Name)
		}
		startIndex = endIndex
	}
	if fail { //clean all data inserted
		_, err := s.client.KV.Txn(context.TODO()).Then(
			clientv3.OpDelete(s.KeyPrefix+ServicesKey+KeySeparator+service.Name+KeySeparator, clientv3.WithPrefix()),
		).Commit()
		if err != nil {
			return err
		}
		return errors.Errorf(errors.StoreError, "failed to create service %q", service.Name)
	}
	return nil

}

//delete application from etcd3
func (s *Store) DeleteService(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator), ">", 0), //key exist
	).Then(
		clientv3.OpDelete(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator, clientv3.WithPrefix()),
	).Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.Errorf(errors.EntityNotFound, "service %q is not found", serviceName)
	}
	return nil
}

func (s *Store) DeleteServices() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.KV.Txn(ctx).Then(
		clientv3.OpDelete(s.KeyPrefix+ServicesKey+KeySeparator, clientv3.WithPrefix()),
	).Commit()
	if err != nil {
		return err
	}
	return nil
}

//get the storage type of the store
func (s *Store) Type() string {
	return StoreType
}

func validateFunc(function *pms.Function) error {
	if function.Name == "" || function.FuncURL == "" {
		return errors.New(errors.InvalidRequest, "\"name\" and \"funcURL\" in function definition can not be empty")
	}
	return nil
}

func (s *Store) CreateFunction(function *pms.Function) (*pms.Function, error) {
	if err := validateFunc(function); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	functionKey := s.KeyPrefix + FunctionsKey + KeySeparator + function.Name
	value, err := json.Marshal(*function)
	if err != nil {
		return nil, errors.Wrap(err, errors.StoreError, "failed to marshal function")
	}
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(functionKey), "=", 0), //policy key does not exist
	).Then(
		clientv3.OpPut(functionKey, string(value)),
	).Commit()
	if err != nil {
		return nil, errors.Wrap(err, errors.StoreError, "failed to insert function into etcd server")
	}
	if !txnResp.Succeeded {
		return nil, errors.Errorf(errors.EntityAlreadyExists, "function %q already exists", function.Name)
	}
	return function, nil
}

func (s *Store) DeleteFunction(funcName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	functionKey := s.KeyPrefix + FunctionsKey + KeySeparator + funcName
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(functionKey), ">", 0), //key exist
	).Then(
		clientv3.OpDelete(functionKey),
	).Commit()
	if err != nil {
		return errors.Wrap(err, errors.StoreError, "failed to delete function from etcd server")
	}
	if !txnResp.Succeeded {
		return errors.Errorf(errors.EntityNotFound, "function %q is not found", funcName)
	}
	return nil
}

func (s *Store) DeleteFunctions() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.KV.Txn(ctx).Then(
		clientv3.OpDelete(s.KeyPrefix+FunctionsKey, clientv3.WithPrefix()),
	).Commit()
	if err != nil {
		return errors.Wrap(err, errors.StoreError, "failed to delete all functions from etcd server")
	}
	return nil
}

func (s *Store) GetFunction(funcName string) (*pms.Function, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	functionKey := s.KeyPrefix + FunctionsKey + KeySeparator + funcName
	getResp, err := s.client.Get(ctx, functionKey)
	if err != nil {
		return nil, errors.Wrapf(err, errors.StoreError, "failed to get function %q from etcd server", funcName)
	}
	if len(getResp.Kvs) == 0 {
		return nil, errors.Errorf(errors.EntityNotFound, "function %q is not found", funcName)
	}
	var function pms.Function
	err = json.Unmarshal(getResp.Kvs[0].Value, &function)
	if err != nil {
		return nil, errors.Errorf(errors.SerializationError, "failed to unmarshal function %q", getResp.Kvs[0].Value)
	}
	return &function, nil
}

func (s *Store) ListAllFunctions(filter string) ([]*pms.Function, error) {
	f := parseFilter(filter)

	functionKeyPrefix := s.KeyPrefix + FunctionsKey + KeySeparator
	responses, err := s.prefixGet(functionKeyPrefix)
	if err != nil {
		return nil, err
	}
	var functions []*pms.Function
	for _, resp := range responses {
		for _, kv := range resp.Kvs {
			var function pms.Function
			err := json.Unmarshal(kv.Value, &function)
			if err != nil {
				return nil, errors.Errorf(errors.SerializationError, "failed to unmarshal function %q", kv.Value)
			}
			isExpected := true
			if f != nil {
				isExpected = nameFilter(function.Name, f)
			}

			if isExpected {
				functions = append(functions, &function)
			}
		}
	}
	return functions, nil
}

func (s *Store) GetFunctionCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	functionKeyPrefix := s.KeyPrefix + FunctionsKey + KeySeparator
	getResp, err := s.client.Get(ctx, functionKeyPrefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, errors.Wrap(err, errors.StoreError, "failed to get function count from etcd server")
	}

	if nil == getResp {
		return 0, nil
	}

	return getResp.Count, nil
}

func (s *Store) Watch() (pms.StorageChangeChannel, error) {
	log.Info("Entering Watch...")
	evalChan := make(chan pms.StoreChangeEvent)
	s.stop = make(chan struct{})
	errChan := make(chan error)
	stopChan := make(chan struct{})
	go func() {
		defer func() {
			close(evalChan)
			close(s.stop)
			close(errChan)
			close(stopChan)
			log.Info("Exiting Watch...")
		}()
	loop:
		for {
			//TODO: reload policy store in case missing any changes during watch failure
			go watch(evalChan, s, errChan, stopChan)
			select {
			case err := <-errChan:
				log.Warningf("Error %v happens, restart watching...\n", err)
				continue
			case <-stopChan:
				log.Warning("Receiving stop signal, stop Watching...")
				break loop
			}
		}

	}()
	return evalChan, nil
}

func watch(evalChan chan pms.StoreChangeEvent, s *Store, errChan chan error, stopChan chan struct{}) {
	watchID := time.Now().Unix()
	log.Infof("Entering watch %v...", watchID)
	cli, err := clientv3.New(*s.Config)
	if err != nil {
		log.Warningf("Error happens when new etcd client, %v, exiting watch...\n", err)
		err := errors.Wrapf(err, errors.StoreError, "failed to connect to etcd server")
		errChan <- err
		return
	}
	defer func() {
		cli.Close()
		log.Infof("Exiting watch %v...", watchID)
	}()

	// Session represents a lease kept alive for the lifetime of a client. Fault-tolerant applications may use sessions to reason about liveness.
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(60))
	if err != nil {
		log.Warningf("Error happens when new session, %v\n", err)
		err = errors.Wrap(err, errors.StoreError, "failed to create session with etcd server")
		errChan <- err
		return
	}

	etcdChan := cli.Watch(context.Background(), s.KeyPrefix, clientv3.WithPrefix())

	for {
		select {
		// receive watch response from etcd
		case resp := <-etcdChan:
			if err := resp.Err(); err != nil {
				log.Warningf("Error happens in watch response, %v\n", err)
				err = errors.Wrap(err, errors.StoreError, "error found in watch response")
				errChan <- err
				return
			}
			for _, e := range resp.Events {
				id := time.Now().Unix()
				//Note: In each policy/rolePolicy creation/deletion, service node (s.KeyPrefix+serviceName+keySeparator) will be updated.
				//so we could only check the event on service node.
				if clientv3.EventTypeDelete == e.Type {
					if strings.HasPrefix(string(e.Kv.Key), s.KeyPrefix+ServicesKey+KeySeparator) {
						serviceName := strings.TrimPrefix(string(e.Kv.Key), s.KeyPrefix+ServicesKey+KeySeparator)
						serviceName = strings.TrimSuffix(serviceName, KeySeparator)
						if strings.Index(serviceName, KeySeparator) == -1 {
							evalChan <- pms.StoreChangeEvent{Type: pms.SERVICE_DELETE, ID: id, Content: []string{serviceName}}
						}
					} else if strings.HasPrefix(string(e.Kv.Key), s.KeyPrefix+FunctionsKey+KeySeparator) {
						functionName := strings.TrimPrefix(string(e.Kv.Key), s.KeyPrefix+FunctionsKey+KeySeparator)
						evalChan <- pms.StoreChangeEvent{Type: pms.FUNCTION_DELETE, ID: id, Content: []string{functionName}}
					}

				} else if clientv3.EventTypePut == e.Type {
					if strings.HasPrefix(string(e.Kv.Key), s.KeyPrefix+ServicesKey+KeySeparator) {
						serviceName := strings.TrimPrefix(string(e.Kv.Key), s.KeyPrefix+ServicesKey+KeySeparator)
						serviceName = strings.TrimSuffix(serviceName, KeySeparator)
						if strings.Index(serviceName, KeySeparator) == -1 {
							service, err := s.GetService(serviceName)
							if err != nil {
								log.Warningf("Unable get service due to error %v.\n", err)
								continue
							}
							evalChan <- pms.StoreChangeEvent{Type: pms.SERVICE_ADD, ID: id, Content: service}
						}
					} else if strings.HasPrefix(string(e.Kv.Key), s.KeyPrefix+FunctionsKey+KeySeparator) {
						functionName := strings.TrimPrefix(string(e.Kv.Key), s.KeyPrefix+FunctionsKey+KeySeparator)
						function, err := s.GetFunction(functionName)
						if err != nil {
							log.Warningf("Unable to get function due to error %v.\n", err)
						}
						evalChan <- pms.StoreChangeEvent{Type: pms.FUNCTION_ADD, ID: id, Content: function}

					}
				}
			}
			// receive the stop signal
		case <-s.stop:
			log.Warning("Receiving stop signal")
			stopChan <- struct{}{}
			return

		case <-session.Done(): // closed by etcd
			log.Warning("Session is closed by etcd")
			errChan <- errors.New(errors.StoreError, "watch session is closed by remote etcd server")
			return
		}
	}

}

func (s *Store) StopWatch() {
	if s.stop != nil {
		s.stop <- struct{}{}
	}
}

// For policy manager
func (s *Store) ListAllPolicies(serviceName string, filter string) ([]*pms.Policy, error) {
	f := parseFilter(filter)

	policyKeyPrefix := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + PoliciesKey
	responses, err := s.prefixGet(policyKeyPrefix)
	if err != nil {
		return nil, err
	}
	var policies []*pms.Policy
	for _, resp := range responses {
		for _, kv := range resp.Kvs {
			var policy pms.Policy
			err := json.Unmarshal(kv.Value, &policy)
			if err != nil {
				return nil, errors.Wrap(err, errors.SerializationError, "failed to unmarshal policies")
			}
			isExpected := true
			if f != nil {
				isExpected = nameFilter(policy.Name, f)
			}

			if isExpected {
				policies = append(policies, &policy)
			}
		}
	}
	return policies, nil
}

func (s *Store) GetPolicyCount(serviceName string) (int64, error) {
	var policyCount int64 = 0
	if len(serviceName) > 0 {
		// Get the policy count in the specified service
		return s.getPolicyCountImpl(serviceName)
	} else {
		// Get the policy count in all services
		serviceNames, err := s.GetServiceNames()
		if err != nil {
			return 0, err
		}
		for _, curServiceName := range serviceNames {
			curCount, err := s.getPolicyCountImpl(curServiceName)
			if err != nil {
				return 0, err
			}
			policyCount += curCount
		}
	}

	return policyCount, nil
}

func (s *Store) getPolicyCountImpl(serviceName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	policyKeyPrefix := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + PoliciesKey
	getResp, err := s.client.Get(ctx, policyKeyPrefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, errors.Wrap(err, errors.StoreError, "failed to get policy count from etcd server")
	}

	if nil == getResp {
		return 0, nil
	}

	return getResp.Count, nil
}

func (s *Store) GetPolicy(serviceName string, id string) (*pms.Policy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	policyKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + PoliciesKey + KeySeparator + id
	getResp, err := s.client.Get(ctx, policyKey)
	if err != nil {
		return nil, errors.Wrap(err, errors.StoreError, "failed to get policy from etcd server")
	}
	if len(getResp.Kvs) == 0 {
		return nil, errors.Errorf(errors.EntityNotFound, "policy %q is not found in service %q", id, serviceName)
	}
	var policy pms.Policy
	err = json.Unmarshal(getResp.Kvs[0].Value, &policy)
	if err != nil {
		return nil, errors.Wrapf(err, errors.SerializationError, "failed to unmarshal a policy")
	}
	return &policy, nil
}

//TODO: to be implemented
func (s *Store) GetRolePolicies(serviceName string, startID string, amount int) (policies []*pms.RolePolicy, nextID string, err error) {
	if amount <= 0 {
		return nil, "", errors.Errorf(errors.InvalidRequest, "invalid amount %d", amount)
	}
	policyPrefix := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + RolePoliciesKey + KeySeparator
	end := clientv3.GetPrefixRangeEnd(policyPrefix)
	getOpts := []clientv3.OpOption{clientv3.WithRange(end), clientv3.WithLimit(int64(amount)), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)}
	resp, err := s.timeOutGet(policyPrefix+startID, getOpts...)
	if err != nil {
		return nil, "", err
	}
	for _, kv := range resp.Kvs {
		var policy pms.RolePolicy
		err = json.Unmarshal(kv.Value, &policy)
		if err != nil {
			return nil, "", errors.Wrap(err, errors.SerializationError, "failed to unmarshal role policy")
		}
		policies = append(policies, &policy)
	}

	if len(policies) == amount {
		lastKey := string(resp.Kvs[amount-1].Key)
		startOfNextRange := clientv3.GetPrefixRangeEnd(lastKey)
		getOpts = []clientv3.OpOption{clientv3.WithRange(end), clientv3.WithKeysOnly(), clientv3.WithLimit(1), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)}
		resp, err = s.timeOutGet(startOfNextRange, getOpts...)
		if err != nil {
			return nil, "", err
		}
		for _, kv := range resp.Kvs {
			nextID = strings.TrimPrefix(string(kv.Key), policyPrefix)
		}
	}
	return policies, nextID, err
}

func (s *Store) DeletePolicy(serviceName string, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	policyKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + PoliciesKey + KeySeparator + id
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(policyKey), ">", 0), //key exist
	).Then(
		clientv3.OpDelete(policyKey),
		//make sure updating service key is the last operation, so watch could work correctly
		clientv3.OpPut(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator, ""),
	).Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.Errorf(errors.EntityNotFound, "policy %q is not found in service %q", id, serviceName)
	}
	return nil
}

func (s *Store) DeletePolicies(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.KV.Txn(ctx).Then(
		clientv3.OpDelete(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator+PoliciesKey, clientv3.WithPrefix()),
		//make sure updating service key is the last operation, so watch could work correctly
		clientv3.OpPut(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator, ""),
	).Commit()
	if err != nil {
		return errors.Wrap(err, errors.StoreError, "failed to delete all policies from etcd server")
	}
	return nil
}

func (s *Store) CreatePolicy(serviceName string, policy *pms.Policy) (*pms.Policy, error) {
	//TODO:validate policy
	dupPolicy := *policy
	if policy.ID == "" {
		dupPolicy.ID = suid.New().String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	serviceKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator
	policyKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + PoliciesKey + KeySeparator + dupPolicy.ID
	value, err := json.Marshal(dupPolicy)
	if err != nil {
		return nil, errors.Wrap(err, errors.SerializationError, "falied to marshal policy")
	}
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(serviceKey), ">", 0), //service key exist
		clientv3.Compare(clientv3.Version(policyKey), "=", 0),  //policy key does not exist
	).Then(
		clientv3.OpPut(policyKey, string(value)),
		//make sure updating service key is the last operation, so watch could work correctly
		clientv3.OpPut(serviceKey, ""),
	).Commit()
	if err != nil {
		return nil, errors.Wrapf(err, errors.StoreError, "falied to create a policy in service %q", serviceName)
	}
	if !txnResp.Succeeded {
		return nil, errors.Errorf(errors.EntityAlreadyExists, "policy %q already exists in service %q", policy.ID, serviceName)
	}
	return &dupPolicy, nil
}

// For role policy manager
func (s *Store) ListAllRolePolicies(serviceName string, filter string) ([]*pms.RolePolicy, error) {
	f := parseFilter(filter)
	rolePolicyKeyPrefix := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + RolePoliciesKey
	responses, err := s.prefixGet(rolePolicyKeyPrefix)
	if err != nil {
		return nil, err
	}
	var rolePolicies []*pms.RolePolicy
	for _, resp := range responses {
		for _, kv := range resp.Kvs {
			var rolePolicy pms.RolePolicy
			err := json.Unmarshal(kv.Value, &rolePolicy)
			if err != nil {
				return nil, errors.New(errors.SerializationError, "failed to unmarshal role policy")
			}
			isExpected := true
			if f != nil {
				isExpected = nameFilter(rolePolicy.Name, f)
			}
			if isExpected {
				rolePolicies = append(rolePolicies, &rolePolicy)
			}
		}
	}
	return rolePolicies, nil
}

func (s *Store) GetRolePolicyCount(serviceName string) (int64, error) {
	var rolePolicyCount int64 = 0
	if len(serviceName) > 0 {
		// Get the rolePolicy count in the specified service
		return s.getRolePolicyCountImpl(serviceName)
	} else {
		// Get the rolePolicy count in all services
		serviceNames, err := s.GetServiceNames()
		if err != nil {
			return 0, err
		}
		for _, curServiceName := range serviceNames {
			curCount, err := s.getRolePolicyCountImpl(curServiceName)
			if err != nil {
				return 0, err
			}
			rolePolicyCount += curCount
		}
	}

	return rolePolicyCount, nil
}

func (s *Store) getRolePolicyCountImpl(serviceName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	rolePolicyKeyPrefix := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + RolePoliciesKey
	getResp, err := s.client.Get(ctx, rolePolicyKeyPrefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, errors.New(errors.StoreError, "failed to get role policy count from etcd server")
	}

	if nil == getResp {
		return 0, nil
	}

	return getResp.Count, nil
}

//TODO: to be implemented
func (s *Store) GetPolicies(serviceName string, startID string, amount int) (policies []*pms.Policy, nextID string, err error) {
	if amount <= 0 {
		return nil, "", errors.Errorf(errors.InvalidRequest, "invalid input amount %d", amount)
	}
	policyPrefix := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + PoliciesKey + KeySeparator
	end := clientv3.GetPrefixRangeEnd(policyPrefix)
	getOpts := []clientv3.OpOption{clientv3.WithRange(end), clientv3.WithLimit(int64(amount)), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)}
	resp, err := s.timeOutGet(policyPrefix+startID, getOpts...)
	if err != nil {
		return nil, "", err
	}
	for _, kv := range resp.Kvs {
		var policy pms.Policy
		err = json.Unmarshal(kv.Value, &policy)
		if err != nil {
			return nil, "", errors.Wrap(err, errors.SerializationError, "failed to unmarshal policies")
		}
		policies = append(policies, &policy)
	}
	if len(policies) == amount {
		lastKey := string(resp.Kvs[amount-1].Key)
		startOfNextRange := clientv3.GetPrefixRangeEnd(lastKey)
		getOpts = []clientv3.OpOption{clientv3.WithRange(end), clientv3.WithKeysOnly(), clientv3.WithLimit(1), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)}
		resp, err = s.timeOutGet(startOfNextRange, getOpts...)
		if err != nil {
			return nil, "", err
		}
		for _, kv := range resp.Kvs {
			nextID = strings.TrimPrefix(string(kv.Key), policyPrefix)
		}
	}
	return policies, nextID, err
}

func (s *Store) GetRolePolicy(serviceName string, id string) (*pms.RolePolicy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	rolePolicyKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + RolePoliciesKey + KeySeparator + id
	getResp, err := s.client.Get(ctx, rolePolicyKey)
	if err != nil {
		return nil, errors.Wrap(err, errors.StoreError, "failed to get a role policy from etcd server")
	}
	if len(getResp.Kvs) == 0 {
		return nil, errors.Errorf(errors.EntityNotFound, "role policy %q is not found in service %q", id, serviceName)
	}
	var rolePolicy pms.RolePolicy
	err = json.Unmarshal(getResp.Kvs[0].Value, &rolePolicy)
	if err != nil {
		return nil, errors.Wrap(err, errors.SerializationError, "failed to unmarshal role policy")
	}
	return &rolePolicy, nil
}

func (s *Store) DeleteRolePolicy(serviceName string, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	rolePolicyKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + RolePoliciesKey + KeySeparator + id
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(rolePolicyKey), ">", 0), //key exist
	).Then(
		clientv3.OpDelete(rolePolicyKey),
		//make sure updating service key is the last operation, so watch could work correctly
		clientv3.OpPut(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator, ""),
	).Commit()
	if err != nil {
		return errors.Wrap(err, errors.StoreError, "failed to delete a role policy from etcd server")
	}
	if !txnResp.Succeeded {
		return errors.Errorf(errors.EntityNotFound, "role policy %q is not found in service %q", id, serviceName)
	}
	return nil
}

func (s *Store) DeleteRolePolicies(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.KV.Txn(ctx).Then(
		clientv3.OpDelete(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator+RolePoliciesKey, clientv3.WithPrefix()),
		//make sure updating service key is the last operation, so watch could work correctly
		clientv3.OpPut(s.KeyPrefix+ServicesKey+KeySeparator+serviceName+KeySeparator, ""),
	).Commit()
	if err != nil {
		return errors.Wrap(err, errors.StoreError, "failed to delete all policies from etcd server")
	}
	return nil
}

func (s *Store) CreateRolePolicy(serviceName string, rolePolicy *pms.RolePolicy) (*pms.RolePolicy, error) {
	//TODO: validate rolePolicy
	dupRolePolicy := *rolePolicy
	if rolePolicy.ID == "" {
		dupRolePolicy.ID = suid.New().String()
	}
	serviceKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator
	rolePolicyKey := s.KeyPrefix + ServicesKey + KeySeparator + serviceName + KeySeparator + RolePoliciesKey + KeySeparator + dupRolePolicy.ID
	value, err := json.Marshal(dupRolePolicy)
	if err != nil {
		return nil, errors.Wrap(err, errors.SerializationError, "failed to marshal role policy")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	txnResp, err := s.client.KV.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(serviceKey), ">", 0),    //service key exist
		clientv3.Compare(clientv3.Version(rolePolicyKey), "=", 0), //role policy key does not exist
	).Then(
		clientv3.OpPut(rolePolicyKey, string(value)),
		//make sure updating service key is the last operation, so watch could work correctly
		clientv3.OpPut(serviceKey, ""),
	).Commit()
	if err != nil {
		return nil, errors.Wrap(err, errors.StoreError, "failed to create role policy in etcd server")
	}
	if !txnResp.Succeeded {
		return nil, errors.Errorf(errors.EntityAlreadyExists, "role policy %q already exists in service %q", dupRolePolicy.ID, serviceName)
	}
	return &dupRolePolicy, nil
}

type filter struct {
	field    string
	operator string
	target   string
}

func (f filter) String() string {
	return fmt.Sprint(f.field, f.operator, f.target)
}

func parseFilter(filterStr string) *filter {
	if len(filterStr) == 0 {
		return nil
	}
	values := strings.Split(filterStr, " ")
	if len(values) >= 2 {
		f := &filter{
			field:    values[0],
			operator: values[1],
		}
		if len(values) > 2 {
			f.target = values[2]
		}
		return f
	} else {
		log.Error("invalid filter string:", filterStr)
		return nil
	}

}

func nameFilter(name string, f *filter) bool { //this filter function return true when the input filter is invalid.
	if f.field != "name" {
		log.Error("invalid name filter. filter is:", f)
		return true
	}
	switch f.operator {
	case "eq":
		return name == f.target
	case "co":
		return strings.Contains(name, f.target)
	case "sw":
		return strings.HasPrefix(name, f.target)
	case "pr":
		return len(name) > 0
	case "gt":
		return name > f.target
	case "ge":
		return name >= f.target
	case "lt":
		return name < f.target
	case "le":
		return name <= f.target
	default:
		log.Error("invalid name filter:", f)
		return true
	}
}
