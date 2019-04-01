//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/cmd/spctl/pdl"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/suid"
	log "github.com/sirupsen/logrus"
)

type phase int

const (
	phaseRoot phase = iota
	phaseService
	phasePolicy
	phaseRolepolicy
	phaseUnknown
)

func (t phase) String() string {
	name := []string{"phaseRoot", "phaseService", "phasePolicy", "phaseRolepolicy"}
	i := int(t)
	switch {
	case i < int(phaseUnknown):
		return name[i]
	default:
		return name[int(phaseUnknown)]
	}
}

type lineType int

const (
	lineEmpty lineType = iota
	lineSection
	linePolicyDef
	lineUnknown
)

func (t lineType) String() string {
	name := []string{"lineEmpty", "lineSection", "linePolicyDef", "lineUnknown"}
	i := int(t)
	switch {
	case i < int(lineUnknown):
		return name[i]
	default:
		return name[int(lineUnknown)]
	}
}

type lineCtx struct {
	no      int
	origin  string
	trimed  string
	ltype   lineType
	errLoc  int
	section string
	phs     phase
	service *pms.Service
}

var emptyPS pms.PolicyStore

func (s *Store) readSPDLWithoutLock() (*pms.PolicyStore, error) {
	var ps pms.PolicyStore

	f, err := os.Open(s.FileLocation)
	if err != nil {
		return &emptyPS, errors.Wrapf(err, errors.StoreError, "unable to open file %q", s.FileLocation)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Warnf("Unable to close file %s because of error %v", s.FileLocation, err)
		}
	}()

	lc := lineCtx{}
	r := bufio.NewReader(f)
	for {
		if err := readLine(r, &lc); err != nil {
			if err == io.EOF {
				break
			} else {
				return &emptyPS, err
			}
		}

		if err := determineType(&lc); err != nil {
			return &emptyPS, err
		}
		switch lc.ltype {
		case lineEmpty:
			continue
		case lineSection:
			if err := processSection(&ps, &lc); err != nil {
				return &emptyPS, err
			}
		case linePolicyDef:
			if err := processPolicyDef(&ps, &lc); err != nil {
				return &emptyPS, err
			}
		default:
			return &emptyPS, fmt.Errorf("Unknown line type %s", lc.ltype)
		}
	}

	return &ps, nil
}

func getServiceSection(ps *pms.PolicyStore, service string) *pms.Service {
	for _, svc := range ps.Services {
		if svc.Name == service {
			return svc
		}
	}
	return nil
}

func processPolicySection(ps *pms.PolicyStore, lc *lineCtx) error {
	// Check if in correct service
	if lc.phs == phaseRoot {
		return fmt.Errorf("Policy section %s is in wrong service section", lc.section)
	}
	lc.phs = phasePolicy

	return nil
}

func processRolePolicySection(ps *pms.PolicyStore, lc *lineCtx) error {
	// Check if in correct service
	if lc.phs == phaseRoot {
		return fmt.Errorf("Policy section %s is in wrong service section", lc.section)
	}
	lc.phs = phaseRolepolicy

	return nil
}

func processServiceSection(ps *pms.PolicyStore, lc *lineCtx) error {
	serviceName := lc.section[len("service."):]
	service := getServiceSection(ps, serviceName)
	if service != nil {
		lc.service = service
		return nil
	}
	lc.service = &pms.Service{
		Name: serviceName,
	}
	lc.phs = phaseService
	ps.Services = append(ps.Services, lc.service)
	return nil
}

func processSection(ps *pms.PolicyStore, lc *lineCtx) error {
	// There are three kinds of sections, service, policy and rolepolicy
	switch {
	case strings.HasPrefix(lc.section, "service."):
		return processServiceSection(ps, lc)
	case lc.section == "policy":
		return processPolicySection(ps, lc)
	case lc.section == "rolepolicy":
		return processRolePolicySection(ps, lc)
	default:
		return fmt.Errorf("Unknown section %s", lc.section)
	}
}

func processPolicyPDL(ps *pms.PolicyStore, lc *lineCtx) error {
	policy, _, err := pdl.ParsePolicy(lc.trimed, "")
	if err != nil {
		return err
	}
	policy.ID = suid.New().String()
	lc.service.Policies = append(lc.service.Policies, policy)
	return nil
}

func processRolePolicyPDL(ps *pms.PolicyStore, lc *lineCtx) error {
	rolePolicy, _, err := pdl.ParseRolePolicy(lc.trimed, "")
	if err != nil {
		return err
	}
	rolePolicy.ID = suid.New().String()
	lc.service.RolePolicies = append(lc.service.RolePolicies, rolePolicy)
	return nil
}

func processPolicyDef(ps *pms.PolicyStore, lc *lineCtx) error {
	switch lc.phs {
	case phasePolicy:
		return processPolicyPDL(ps, lc)
	case phaseRolepolicy:
		return processRolePolicyPDL(ps, lc)
	default:
		return fmt.Errorf("Wrong policy definition at line %d", lc.no)
	}
}

func determineType(lc *lineCtx) error {
	lc.trimed = lc.origin
	// Trim comments
	idx := strings.Index(lc.origin, "#")
	if idx != -1 {
		lc.trimed = lc.origin[0:idx]
	}

	// Trim spaces
	lc.trimed = strings.TrimSpace(lc.trimed)
	if len(lc.trimed) == 0 || lc.trimed[0] == '#' {
		// blank line
		lc.ltype = lineEmpty
		return nil
	}

	if lc.trimed[0] == ']' {
		return fmt.Errorf("Syntax error near %s at line %d", lc.trimed, lc.no)
	}

	if lc.trimed[0] == '[' {
		if len(lc.trimed) == 2 || lc.trimed[len(lc.trimed)-1] != ']' {
			// This is an error, begin with [, but don't end with ]
			return fmt.Errorf("Syntax error near %s at line %d", lc.trimed, lc.no)
		}
		// length > 2 and wrap with []
		lc.section = lc.trimed[1:(len(lc.trimed) - 1)]
		lc.ltype = lineSection
		return nil
	}

	// Don't start with [, assume this is a line for SPDL
	lc.ltype = linePolicyDef

	return nil
}

func readLine(r *bufio.Reader, lc *lineCtx) error {
	lc.origin = ""
	for {
		bs, isp, err := r.ReadLine()
		if err != nil {
			return err
		}
		lc.origin += string(bs)
		if !isp {
			lc.no++
			break
		}
	}

	return nil
}
