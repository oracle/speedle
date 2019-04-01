//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/store"
	log "github.com/sirupsen/logrus"
)

type discoverRequestStore struct {
	FileLocation string
	rwLock       sync.RWMutex
}

type RequestItem struct {
	Index   int64               `json:"index"`
	Request *ads.RequestContext `json:"request"`
}

type StoreContent struct {
	Requests []*RequestItem `json:"requests"`
}

const (
	discoverStoreFileName = "speedle_discover_requests.json"
)

//read policy store from file
func (s *discoverRequestStore) ReadDiscoverRequestStore() (*StoreContent, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.readDiscoverRequestStoreWithoutLock()
}

func (s *discoverRequestStore) readDiscoverRequestStoreWithoutLock() (*StoreContent, error) {
	var drs StoreContent
	raw, err := ioutil.ReadFile(s.FileLocation)
	if err != nil {
		return &drs, errors.Wrapf(err, errors.StoreError, "unable to read file %q", s.FileLocation)
	}
	if err := json.Unmarshal(raw, &drs); err != nil {
		log.Warnf("error reading discover request store because of error %v", err)
		return &StoreContent{}, err
	}
	return &drs, nil
}

func (s *discoverRequestStore) WriteDiscoverRequestStore(drs *StoreContent) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	return s.writeDiscoverRequestStoreWithoutLock(drs)
}

func (s *discoverRequestStore) writeDiscoverRequestStoreWithoutLock(drs *StoreContent) error {
	jsonFile, err := os.Create(s.FileLocation)
	defer jsonFile.Close()
	if err != nil {
		return errors.Wrapf(err, errors.StoreError, "unable to create file %q", s.FileLocation)
	}
	drsB, err := json.MarshalIndent(*drs, "", "    ")
	if err != nil {
		log.Errorf("marshal indent filed becuase of %v", err)
		return errors.Wrap(err, errors.SerializationError, "marshal indent failed")
	}
	if _, err := jsonFile.Write(drsB); err != nil {
		log.Errorf("write to file failed becuase of %v", err)
		return errors.Wrapf(err, errors.StoreError, "unable to write to file %q", s.FileLocation)
	}
	return nil
}

func getDiscoverRequestStore(s *Store) (*discoverRequestStore, error) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	if s.discoverStore == nil {
		dir, _ := filepath.Split(s.FileLocation)
		discoverFileLocation := filepath.Join(dir, discoverStoreFileName)
		if dir == "./" {
			discoverFileLocation = dir + discoverStoreFileName
		}
		log.Infof("discover store file location:%s\n", discoverFileLocation)
		if _, err := os.Stat(discoverFileLocation); os.IsNotExist(err) {
			log.Infof("discover store file does not exist, create one...")
			if err1 := ioutil.WriteFile(discoverFileLocation, []byte("{}"), 0644); err1 != nil {
				log.Errorf("error creating discover store file: %v\n", err1)
				return nil, err1
			}
		}
		s.discoverStore = &discoverRequestStore{FileLocation: discoverFileLocation}
	}
	return s.discoverStore, nil

}

// SaveDiscoverRequest saves discover request
func (s *Store) SaveDiscoverRequest(discoverRequest *ads.RequestContext) error {
	if discoverStore, err := getDiscoverRequestStore(s); err == nil {
		return discoverStore.saveDiscoverRequest(discoverRequest)
	} else {
		return err
	}
}
func (s *discoverRequestStore) saveDiscoverRequest(discoverRequest *ads.RequestContext) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	sContent, err := s.readDiscoverRequestStoreWithoutLock()
	if err != nil {
		return err
	}
	var idx = int64(0)
	if sContent.Requests != nil && len(sContent.Requests) > 0 {
		if math.MaxInt64 == sContent.Requests[len(sContent.Requests)-1].Index {
			return errors.New(errors.ExceedLimit, "reach largest revision number")
		}
		idx = sContent.Requests[len(sContent.Requests)-1].Index + 1

	}
	sContent.Requests = append(sContent.Requests, &RequestItem{Index: idx, Request: discoverRequest})
	return s.writeDiscoverRequestStoreWithoutLock(sContent)
}

// GetLastDiscoverRequest gets last request log
func (s *Store) GetLastDiscoverRequest(serviceName string) (*ads.RequestContext, int64, error) {
	if discoverStore, err := getDiscoverRequestStore(s); err == nil {
		return discoverStore.getLastDiscoverRequest(serviceName)
	} else {
		return nil, -1, err
	}
}
func (s *discoverRequestStore) getLastDiscoverRequest(serviceName string) (*ads.RequestContext, int64, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	sContent, err := s.readDiscoverRequestStoreWithoutLock()
	if err != nil {
		return nil, -1, err
	}

	if sContent.Requests != nil && len(sContent.Requests) > 0 {
		if serviceName == "" {
			reqItem := sContent.Requests[len(sContent.Requests)-1]
			return reqItem.Request, reqItem.Index, nil
		}
		for i := len(sContent.Requests) - 1; i >= 0; i-- {
			if sContent.Requests[i].Request.ServiceName == serviceName {
				reqItem := sContent.Requests[i]
				return reqItem.Request, reqItem.Index, nil
			}
		}
	}
	return nil, -1, errors.Errorf(errors.EntityNotFound, "no discover request found for service %q", serviceName)
}

// GetDiscoverRequestsSinceRevision gets request logs since a revision.
func (s *Store) GetDiscoverRequestsSinceRevision(serviceName string, revision int64) ([]*ads.RequestContext, int64, error) {
	if discoverStore, err := getDiscoverRequestStore(s); err == nil {
		return discoverStore.getDiscoverRequestsSinceRevision(serviceName, revision)
	} else {
		return nil, -1, err
	}
}
func (s *discoverRequestStore) getDiscoverRequestsSinceRevision(serviceName string, revision int64) ([]*ads.RequestContext, int64, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	sContent, err := s.readDiscoverRequestStoreWithoutLock()
	if err != nil {
		return nil, -1, err
	}
	requests := []*ads.RequestContext{}
	if sContent.Requests != nil && len(sContent.Requests) > 0 {
		for _, reqItem := range sContent.Requests {
			if reqItem.Index > revision {
				if serviceName == "" || serviceName == reqItem.Request.ServiceName {
					requests = append(requests, reqItem.Request)
				}
			}
		}
		return requests, sContent.Requests[len(sContent.Requests)-1].Index, nil
	}
	return nil, -1, errors.Errorf(errors.EntityNotFound, "no discover request found for service %q.", serviceName)
}

// GetDiscoverRequests gets request logs for a service.
// Get all requests when serviceName is empty.
func (s *Store) GetDiscoverRequests(serviceName string) ([]*ads.RequestContext, int64, error) {
	if discoverStore, err := getDiscoverRequestStore(s); err == nil {
		return discoverStore.getDiscoverRequests(serviceName)
	} else {
		return nil, -1, err
	}
}
func (s *discoverRequestStore) getDiscoverRequests(serviceName string) ([]*ads.RequestContext, int64, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	sContent, err := s.readDiscoverRequestStoreWithoutLock()
	if err != nil {
		return nil, -1, err
	}
	requests := []*ads.RequestContext{}

	if sContent.Requests != nil && len(sContent.Requests) > 0 {
		for _, reqItem := range sContent.Requests {
			if serviceName == "" || serviceName == reqItem.Request.ServiceName {
				requests = append(requests, reqItem.Request)
			}
		}
		return requests, sContent.Requests[len(sContent.Requests)-1].Index, nil
	}

	return nil, -1, errors.Errorf(errors.EntityNotFound, "no discover request found for service %q.", serviceName)
}

// ResetDiscoverRequests cleans request logs for a service.
// Clean all request logs when serviceName is empty.
func (s *Store) ResetDiscoverRequests(serviceName string) error {
	discoverStore, err := getDiscoverRequestStore(s)
	if err != nil {
		return err
	}

	return discoverStore.resetDiscoverRequests(serviceName)
}

func (s *discoverRequestStore) resetDiscoverRequests(serviceName string) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	sContent, err := s.readDiscoverRequestStoreWithoutLock()
	if err != nil {
		return err
	}
	if sContent.Requests != nil && len(sContent.Requests) != 0 {
		if serviceName == "" { //reset all
			sContent = &StoreContent{Requests: []*RequestItem{}}
		} else {
			for i := len(sContent.Requests) - 1; i >= 0; i-- {
				if sContent.Requests[i].Request.ServiceName == serviceName {
					sContent.Requests = append(sContent.Requests[:i], sContent.Requests[i+1:]...)
				}
			}
		}
		return s.writeDiscoverRequestStoreWithoutLock(sContent)
	}
	return nil
}

//Generate policies for principal based on existing request logs. Generate policies for all principals when principalXXX are empty.
func (s *Store) GeneratePolicies(serviceName, principalType, principalName, principalIDD string) (map[string]*pms.Service, int64, error) {
	if discoverStore, err := getDiscoverRequestStore(s); err == nil {
		return discoverStore.generatePolicies(serviceName, principalType, principalName, principalIDD)
	} else {
		return nil, -1, err
	}
}
func (s *discoverRequestStore) generatePolicies(serviceName, principalType, principalName, principalIDD string) (map[string]*pms.Service, int64, error) {
	requests, revision, err := s.getDiscoverRequests(serviceName)
	if err != nil {
		return nil, -1, err
	}
	serviceMap, err := store.GeneratePoliciesFromDiscoverRequests(requests, principalType, principalName, principalIDD)
	if err != nil {
		return nil, -1, err
	}
	return serviceMap, revision, nil

}
