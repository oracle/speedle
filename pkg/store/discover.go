//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package store

import (
	"reflect"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/subjectutils"
)

const (
	DefaultMaxDiscoverRequestNum                = int64(50000)
	DefaultDeleteNumWhenReachMaxDiscoverRequest = int64(100)
)

var MaxDiscoverRequestNum, DeleteNumWhenReachMaxDiscoverRequest int64

func init() {
	MaxDiscoverRequestNum = DefaultMaxDiscoverRequestNum
	DeleteNumWhenReachMaxDiscoverRequest = DefaultDeleteNumWhenReachMaxDiscoverRequest
}

type DiscoverRequestManager interface {
	//Save discover request
	SaveDiscoverRequest(discoverRequest *ads.RequestContext) error
	//Get last request log
	GetLastDiscoverRequest(serviceName string) (*ads.RequestContext, int64, error)
	//Get request logs since a revision.
	GetDiscoverRequestsSinceRevision(serviceName string, revision int64) ([]*ads.RequestContext, int64, error)
	//Get request logs for a service. Get all requests when serviceName is empty.
	GetDiscoverRequests(serviceName string) ([]*ads.RequestContext, int64, error)
	//Clean request logs for a service. Clean all request logs when serviceName is empty.
	ResetDiscoverRequests(serviceName string) error
	//Generate policies for principal based on existing request logs. Generate policies for all principals when principalXXX are empty.
	GeneratePolicies(serviceName, principalType, principalName, principalIDD string) (map[string]*pms.Service, int64, error)
}

func GeneratePoliciesFromDiscoverRequests(requests []*ads.RequestContext, principalType, principalName, principalIDD string) (map[string]*pms.Service, error) {
	serviceMap := map[string]*pms.Service{}
	principalMap := map[string]bool{}
	policyMap := map[string]*pms.Policy{}
	for _, req := range requests {
		if _, ok := serviceMap[req.ServiceName]; !ok {
			serviceMap[req.ServiceName] = &pms.Service{Name: req.ServiceName,
				RolePolicies: []*pms.RolePolicy{},
				Policies:     []*pms.Policy{}}
		}
		if len(req.Subject.Principals) == 0 { //anonymous_role
			roleName := ads.BuiltIn_Role_Anonymous
			policyKey := "svc=" + req.ServiceName + ";res=" + req.Resource + ";role=" + roleName
			if p, ok := policyMap[policyKey]; ok {
				alreadyExist := false
				for _, act := range p.Permissions[0].Actions {
					if act == req.Action {
						alreadyExist = true
						break
					}
				}
				if !alreadyExist {
					p.Permissions[0].Actions = append(p.Permissions[0].Actions, req.Action)
				}
			} else {
				policy := pms.Policy{
					Effect:     "grant",
					Principals: [][]string{{"role:" + roleName}},
					Permissions: []*pms.Permission{
						{
							Resource: req.Resource,
							Actions:  []string{req.Action},
						},
					},
				}
				serviceMap[req.ServiceName].Policies = appendPolicyIfUnique(serviceMap[req.ServiceName].Policies, &policy)
				policyMap[policyKey] = &policy
			}
		} else {
			for _, princ := range req.Subject.Principals {
				if principalType != "" && principalType != princ.Type ||
					principalName != "" && principalName != princ.Name ||
					principalIDD != "" && principalIDD != princ.IDD {
					break
				}
				encodedPrincipal := subjectutils.EncodePrincipal(princ)
				roleName := "role_" + encodedPrincipal
				policyKey := "svc=" + req.ServiceName + ";res=" + req.Resource + ";role=" + roleName

				if _, ok := principalMap[encodedPrincipal]; !ok {
					rolePolicy := pms.RolePolicy{
						Effect:     "grant",
						Principals: []string{encodedPrincipal},
						Roles:      []string{roleName},
					}
					serviceMap[req.ServiceName].RolePolicies = append(serviceMap[req.ServiceName].RolePolicies, &rolePolicy)
					principalMap[encodedPrincipal] = true
				}
				if p, ok := policyMap[policyKey]; ok {
					alreadyExist := false
					for _, act := range p.Permissions[0].Actions {
						if act == req.Action {
							alreadyExist = true
							break
						}
					}
					if !alreadyExist {
						p.Permissions[0].Actions = append(p.Permissions[0].Actions, req.Action)
					}

				} else {
					policy := pms.Policy{
						Effect:     "grant",
						Principals: [][]string{{"role:" + roleName}},
						Permissions: []*pms.Permission{
							{
								Resource: req.Resource,
								Actions:  []string{req.Action},
							},
						},
					}
					serviceMap[req.ServiceName].Policies = appendPolicyIfUnique(serviceMap[req.ServiceName].Policies, &policy)
					policyMap[policyKey] = &policy
				}
			}

		}
	}
	return serviceMap, nil
}

func appendPolicyIfUnique(policies []*pms.Policy, newPolicy *pms.Policy) []*pms.Policy {
	for _, policy := range policies {
		if reflect.DeepEqual(*policy, *newPolicy) {
			return policies
		}
	}
	return append(policies, newPolicy)
}
