//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/suid"

	"github.com/fsnotify/fsnotify"
	"github.com/oracle/speedle/api/pms"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	FileLocation  string
	stop          chan struct{}
	rwLock        sync.RWMutex
	discoverStore *discoverRequestStore
}

// ReadPolicyStore reads policy store from a file
func (s *Store) ReadPolicyStore() (*pms.PolicyStore, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.readPolicyStoreWithoutLock()
}

func (s *Store) readPolicyStoreWithoutLock() (*pms.PolicyStore, error) {
	if strings.HasSuffix(s.FileLocation, ".spdl") {
		return s.readSPDLWithoutLock()
	}

	var ps pms.PolicyStore

	f, err := os.Open(s.FileLocation)
	if err != nil {
		return &ps, errors.Wrapf(err, errors.StoreError, "unable to open file %q", s.FileLocation)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Warnf("Error when closing file %s", s.FileLocation)
		}
	}()

	decoder := json.NewDecoder(bufio.NewReader(f))
	if err := decoder.Decode(&ps); err != nil {
		log.Warnf("Unable to parse %s in JSON format because of error %v", s.FileLocation, err)
		return &pms.PolicyStore{}, err
	}

	return &ps, nil
}

// WritePolicyStore writes policies to a file
func (s *Store) WritePolicyStore(ps *pms.PolicyStore) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	return s.writePolicyStoreWithoutLock(ps)
}

func (s *Store) writePolicyStoreWithoutLock(ps *pms.PolicyStore) error {
	jsonFile, err := os.Create(s.FileLocation)
	defer jsonFile.Close()
	if err != nil {
		return errors.Wrapf(err, errors.StoreError, "unable to create file %q", s.FileLocation)
	}
	psB, err := json.MarshalIndent(ps, "", "    ")
	if err != nil {
		return errors.Wrap(err, errors.StoreError, "marshal indent failed")
	}
	if _, err := jsonFile.Write(psB); err != nil {
		return errors.Wrapf(err, errors.StoreError, "unable to write to file %q", s.FileLocation)
	}
	return nil
}

// ListAllServices lists all the services
func (s *Store) ListAllServices() ([]*pms.Service, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.getServicesWithoutLock()
}

func (s *Store) getServicesWithoutLock() ([]*pms.Service, error) {
	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return nil, err
	}
	return ps.Services, nil
}

// GetServiceNames reads all the service names
func (s *Store) GetServiceNames() ([]string, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.getServiceNamesWithoutLock()
}

func (s *Store) getServiceNamesWithoutLock() ([]string, error) {

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return nil, err
	}

	var serviceNames []string

	for _, srv := range ps.Services {
		serviceNames = append(serviceNames, srv.Name)
	}

	return serviceNames, nil
}

// GetPolicyAndRolePolicyCounts returns a map, in which the key is the service name, and the value is the count of both policies and role policies in the service.
func (s *Store) GetPolicyAndRolePolicyCounts() (map[string]*pms.PolicyAndRolePolicyCount, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	serviceNames, err := s.getServiceNamesWithoutLock()
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]*pms.PolicyAndRolePolicyCount)

	for _, srvName := range serviceNames {
		var counts pms.PolicyAndRolePolicyCount

		policyCount, err := s.getPolicyCountWithoutLock(srvName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get policy count for service: %s", srvName)
		}
		counts.PolicyCount = policyCount

		rolePolicyCount, err := s.getRolePolicyCountWithoutLock(srvName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get rolePolicy count for service: %s", srvName)
		}
		counts.RolePolicyCount = rolePolicyCount

		countMap[srvName] = &counts
	}

	return countMap, nil
}

// GetServiceCount gets the service count
func (s *Store) GetServiceCount() (int64, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return 0, err
	}

	if nil == ps.Services {
		return 0, nil
	} else {
		return int64(len(ps.Services)), nil
	}
}

// GetService gets the detailed info of a service
func (s *Store) GetService(serviceName string) (*pms.Service, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.getServiceWithoutLock(serviceName)
}

func (s *Store) getServiceWithoutLock(serviceName string) (*pms.Service, error) {
	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return nil, err
	}
	for _, service := range ps.Services {
		if serviceName == service.Name {
			return service, nil
		}
	}
	return nil, errors.Errorf(errors.EntityNotFound, "service %q is not found", serviceName)
}

// CreateService creates a new service
func (s *Store) CreateService(service *pms.Service) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return err
	}
	for _, value := range ps.Services {
		if service.Name == value.Name {
			return errors.Errorf(errors.EntityAlreadyExists, "service %q already exists", service.Name)
		}
	}
	serviceWithIDs, err := generateID(service)
	if err == nil {
		ps.Services = append(ps.Services, serviceWithIDs)
		err = s.writePolicyStoreWithoutLock(ps)
	}
	return err
}

func generateID(service *pms.Service) (*pms.Service, error) {
	var result pms.Service
	result = *service
	for _, policy := range result.Policies {
		policy.ID = suid.New().String()
	}
	for _, rolePolicy := range result.RolePolicies {
		rolePolicy.ID = suid.New().String()
	}
	return &result, nil
}

// WriteService writes a service into a file
func (s *Store) WriteService(service *pms.Service) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	return s.writeServiceWithoutLock(service)
}

func (s *Store) writeServiceWithoutLock(service *pms.Service) error {

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return err
	}
	for index, value := range ps.Services {
		if service.Name == value.Name {
			ps.Services = append(ps.Services[:index], ps.Services[index+1:]...)
			break
		}
	}
	ps.Services = append(ps.Services, service)
	if err := s.writePolicyStoreWithoutLock(ps); err != nil {
		return err
	}
	return nil
}

// DeleteService deletes a service named ${serviceName} from a file
func (s *Store) DeleteService(serviceName string) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return err
	}
	found := false
	for index, value := range ps.Services {
		if serviceName == value.Name {
			ps.Services = append(ps.Services[:index], ps.Services[index+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf(errors.EntityNotFound, "service %q is not found", serviceName)
	}
	s.writePolicyStoreWithoutLock(ps)
	return nil

}

// DeleteServices deletes all services from a file
func (s *Store) DeleteServices() error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return err
	}
	ps.Services = []*pms.Service{}

	return s.writePolicyStoreWithoutLock(ps)
}

func (s *Store) Watch() (pms.StorageChangeChannel, error) {
	log.Info("Enter Watch...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("Failed to create a new watcher, error: %v", err)
		return nil, errors.Wrap(err, errors.StoreError, "fsnotify new watcher failed")
	}

	err = watcher.Add(s.FileLocation)
	if err != nil {
		log.Errorf("Failed to add the file %q into the watch list, error: %v", s.FileLocation, err)
		return nil, errors.Wrapf(err, errors.StoreError, "Failed to add the file %q into the watch list", s.FileLocation)
	}

	var storeChangeChan pms.StorageChangeChannel
	storeChangeChan = make(chan pms.StoreChangeEvent)

	s.stop = make(chan struct{})

	go func() {
		defer func() {
			watcher.Close()
			close(storeChangeChan)
			close(s.stop)
		}()
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Info("Reloading the file store...")
					reloadEvent := pms.StoreChangeEvent{Type: pms.FULL_RELOAD}
					storeChangeChan <- reloadEvent
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					if _, err := os.Lstat(s.FileLocation); os.IsNotExist(err) {
						log.Fatalf("The policy file %q has already been renamed, please double check", s.FileLocation)
					} else {
						// This is just a workaround for the issue https://github.com/fsnotify/fsnotify/issues/282
						log.Info("Reloading the file store....")
						reloadEvent := pms.StoreChangeEvent{Type: pms.FULL_RELOAD}
						storeChangeChan <- reloadEvent

						err = watcher.Add(s.FileLocation)
						if err != nil {
							log.Fatalf("Failed to add the file %q into the watch list again, error: %v", s.FileLocation, err)
						}
					}
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Fatalf("The policy file %q has already been removed, please double check", s.FileLocation)
				default:
					log.Infof("Operation %q was detected on the policy file %q", event.Op, s.FileLocation)
				}
			case err := <-watcher.Errors:
				log.Warningf("Error happened when watching the policy file, error: %v", err)
			case <-s.stop:
				log.Warning("Received stop signal")
				return
			}
		}
	}()

	return storeChangeChan, nil
}

func (s *Store) StopWatch() {
	if s.stop != nil {
		s.stop <- struct{}{}
	}
}

func (s *Store) Type() string {
	return StoreType
}

// For policy manager
func (s *Store) ListAllPolicies(serviceName string, filter string) ([]*pms.Policy, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	f := parseFilter(filter)
	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return nil, err
	}
	ret := []*pms.Policy{}
	for _, policy := range service.Policies {
		isExpected := true
		if f != nil {
			isExpected = nameFilter(policy.Name, f)
		}
		if isExpected {
			ret = append(ret, policy)
		}
	}
	return ret, nil
}

func (s *Store) GetPolicyCount(serviceName string) (int64, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.getPolicyCountWithoutLock(serviceName)
}

func (s *Store) getPolicyCountWithoutLock(serviceName string) (int64, error) {
	var policyCount int64 = 0
	if len(serviceName) > 0 {
		// Get the policy count in the specified service
		return s.getPolicyCountImpl(serviceName)
	} else {
		// Get the policy count in all services
		services, err := s.getServicesWithoutLock()
		if err != nil {
			return 0, err
		}
		for _, curService := range services {
			curCount, err := s.getPolicyCountImpl(curService.Name)
			if err != nil {
				return 0, err
			}
			policyCount += curCount
		}
	}

	return policyCount, nil
}

func (s *Store) getPolicyCountImpl(serviceName string) (int64, error) {
	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return 0, err
	}

	if nil == service.Policies {
		return 0, nil
	} else {
		return int64(len(service.Policies)), nil
	}
}

func (s *Store) GetPolicy(serviceName string, id string) (*pms.Policy, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return nil, err
	}
	for _, policy := range service.Policies {
		if policy.ID == id {
			// Found
			return policy, nil
		}
	}

	return nil, errors.Errorf(errors.EntityNotFound, "unable to find policy %q in service %q", id, serviceName)
}

func (s *Store) DeletePolicy(serviceName string, id string) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return err
	}
	for index, policy := range service.Policies {
		if policy.ID == id {
			// Found
			service.Policies = append(service.Policies[:index], service.Policies[index+1:]...)
			return s.writeServiceWithoutLock(service)
		}
	}

	return errors.Errorf(errors.EntityNotFound, "unable to find policy %q in service %q", id, serviceName)
}

func (s *Store) DeletePolicies(serviceName string) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return err
	}
	service.Policies = []*pms.Policy{}
	if err := s.writeServiceWithoutLock(service); err != nil {
		return err
	}
	return nil
}

func (s *Store) CreatePolicy(serviceName string, policy *pms.Policy) (*pms.Policy, error) {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return nil, err
	}
	dupPolicy := *policy
	dupPolicy.ID = suid.New().String()

	service.Policies = append(service.Policies, &dupPolicy)
	if err := s.writeServiceWithoutLock(service); err != nil {
		return nil, err
	}
	return &dupPolicy, nil
}

// For role policy manager
func (s *Store) ListAllRolePolicies(serviceName string, filter string) ([]*pms.RolePolicy, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	f := parseFilter(filter)
	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return nil, err
	}
	ret := []*pms.RolePolicy{}
	for _, rolePolicy := range service.RolePolicies {
		isExpected := true
		if f != nil {
			isExpected = nameFilter(rolePolicy.Name, f)
		}
		if isExpected {
			ret = append(ret, rolePolicy)
		}
	}
	return ret, nil
}

func (s *Store) GetRolePolicyCount(serviceName string) (int64, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.getRolePolicyCountWithoutLock(serviceName)
}

func (s *Store) getRolePolicyCountWithoutLock(serviceName string) (int64, error) {
	var rolePolicyCount int64 = 0
	if len(serviceName) > 0 {
		// Get the policy count in the specified service
		return s.getRolePolicyCountImpl(serviceName)
	} else {
		// Get the policy count in all services
		services, err := s.getServicesWithoutLock()
		if err != nil {
			return 0, err
		}
		for _, curService := range services {
			curCount, err := s.getRolePolicyCountImpl(curService.Name)
			if err != nil {
				return 0, err
			}
			rolePolicyCount += curCount
		}
	}

	return rolePolicyCount, nil
}

func (s *Store) getRolePolicyCountImpl(serviceName string) (int64, error) {
	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return 0, err
	}

	if nil == service.RolePolicies {
		return 0, nil
	} else {
		return int64(len(service.RolePolicies)), nil
	}
}

func (s *Store) GetRolePolicy(serviceName string, id string) (*pms.RolePolicy, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return nil, err
	}
	for _, rolePolicy := range service.RolePolicies {
		if rolePolicy.ID == id {
			// Found
			return rolePolicy, nil
		}
	}

	return nil, errors.Errorf(errors.EntityNotFound, "unable to find role policy %q in service %q", id, serviceName)
}

func (s *Store) DeleteRolePolicy(serviceName string, id string) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return err
	}
	for index, rolePolicy := range service.RolePolicies {
		if rolePolicy.ID == id {
			// Found
			service.RolePolicies = append(service.RolePolicies[:index], service.RolePolicies[index+1:]...)
			return s.writeServiceWithoutLock(service)
		}
	}
	return errors.Errorf(errors.EntityNotFound, "unable to find role policy %q in service %q", id, serviceName)
}

func (s *Store) DeleteRolePolicies(serviceName string) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return err
	}
	service.RolePolicies = []*pms.RolePolicy{}

	return s.writeServiceWithoutLock(service)
}

func (s *Store) CreateRolePolicy(serviceName string, rolePolicy *pms.RolePolicy) (*pms.RolePolicy, error) {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	service, err := s.getServiceWithoutLock(serviceName)
	if err != nil {
		return nil, err
	}
	dupRolePolicy := *rolePolicy
	dupRolePolicy.ID = suid.New().String()

	service.RolePolicies = append(service.RolePolicies, &dupRolePolicy)
	if err := s.writeServiceWithoutLock(service); err != nil {
		return nil, err
	}
	return &dupRolePolicy, nil
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
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return nil, err
	}
	for _, value := range ps.Functions {
		if function.Name == value.Name {
			return nil, errors.Errorf(errors.EntityAlreadyExists, "function %q already exists", function.Name)
		}
	}
	ps.Functions = append(ps.Functions, function)

	err = s.writePolicyStoreWithoutLock(ps)
	if err != nil {
		return nil, err
	}

	return function, nil
}

func (s *Store) DeleteFunction(funcName string) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return err
	}

	for index, value := range ps.Functions {
		if funcName == value.Name {
			ps.Functions = append(ps.Functions[:index], ps.Functions[index+1:]...)
			return s.writePolicyStoreWithoutLock(ps)
		}
	}
	return errors.Errorf(errors.EntityNotFound, "function %q is not found", funcName)
}

func (s *Store) DeleteFunctions() error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return err
	}
	ps.Functions = []*pms.Function{}
	return s.writePolicyStoreWithoutLock(ps)
}

func (s *Store) GetFunction(funcName string) (*pms.Function, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return nil, err
	}
	for _, value := range ps.Functions {
		if funcName == value.Name {
			return value, nil
		}
	}
	return nil, errors.Errorf(errors.EntityNotFound, "function %q is not found", funcName)
}

func (s *Store) ListAllFunctions(filter string) ([]*pms.Function, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	f := parseFilter(filter)
	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return nil, err
	}
	ret := []*pms.Function{}
	for _, value := range ps.Functions {
		isExpected := true
		if f != nil {
			isExpected = nameFilter(value.Name, f)
		}
		if isExpected {
			ret = append(ret, value)
		}
	}
	return ret, nil
}

func (s *Store) GetFunctionCount() (int64, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	ps, err := s.readPolicyStoreWithoutLock()
	if err != nil {
		return -1, err
	}
	if nil == ps.Functions {
		return 0, nil
	} else {
		return int64(len(ps.Functions)), nil
	}
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
	if !strings.HasPrefix(filterStr, "name") {
		log.Error("unsupported filter string:", filterStr)
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
