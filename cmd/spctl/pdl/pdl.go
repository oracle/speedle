//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pdl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	adsapi "github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/subjectutils"
)

const (
	grant = "grant"
	deny  = "deny"

	res_expr_prefix = "expr:"
)

func ParsePolicy(cmd, name string) (*pms.Policy, io.Reader, error) {
	effect, i, err := getEffect(cmd)
	if err != nil {
		return nil, nil, err
	}
	principals, i, err := getOrPrincipals(cmd, i)
	if err != nil {
		return nil, nil, err
	}
	perms, i, err := getPermissions(cmd, i)
	if err != nil {
		return nil, nil, err
	}
	if len(perms) == 0 {
		return nil, nil, errors.New("No permission found")
	}
	condition, _, err := getCondition(cmd, i)
	if err != nil {
		return nil, nil, err
	}

	policy := pms.Policy{
		Name:        name,
		Effect:      effect,
		Principals:  principals,
		Permissions: perms,
		Condition:   condition,
	}

	return &policy, toJSON(policy), nil
}

func ParseRolePolicy(cmd, name string) (*pms.RolePolicy, io.Reader, error) {
	effect, i, err := getEffect(cmd)
	if err != nil {
		return nil, nil, err
	}
	principals, i, err := getRolePolicyPrincipals(cmd, i)
	if err != nil {
		return nil, nil, err
	}
	roles, i, err := getRoles(cmd, i)
	if err != nil {
		return nil, nil, err
	}
	if len(roles) == 0 {
		return nil, nil, errors.New("No role found")
	}
	resources, resExps, i, err := getResources(cmd, i)
	if err != nil {
		return nil, nil, err
	}
	condition, _, err := getCondition(cmd, i)
	if err != nil {
		return nil, nil, err
	}
	rolePolicy := pms.RolePolicy{
		Name:                name,
		Effect:              effect,
		Principals:          principals,
		Resources:           resources,
		ResourceExpressions: resExps,
		Roles:               roles,
		Condition:           condition,
	}
	return &rolePolicy, toJSON(rolePolicy), nil
}

func toJSON(i interface{}) io.Reader {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(i)
	return buf
}

func getEffect(cmd string) (string, int, error) {
	i := skipSpaces(cmd, 0)
	if (i+6 <= len(cmd)) && strings.EqualFold("grant ", cmd[i:i+6]) {
		return grant, i + 6, nil
	} else if i+5 <= len(cmd) && strings.EqualFold("deny ", cmd[i:i+5]) {
		return deny, i + 5, nil
	} else {
		return "", -1, getError("Not found valid effect", cmd, i)
	}
}

func getRolePolicyPrincipals(cmd string, i int) ([]string, int, error) {
	ret := []string{}
	onePrincipal, i, err := getPrincipal(cmd, i)
	if err != nil {
		return nil, -1, err
	}
	ret = append(ret, onePrincipal)
	i = skipSpaces(cmd, i)
	if i >= len(cmd) {
		return nil, -1, getError("Unexpected EOF found", cmd, i)
	}
	for cmd[i] == ',' {
		// Skip the comma
		i++
		onePrincipal, i, err = getPrincipal(cmd, i)
		if err != nil {
			return nil, -1, err
		}
		ret = append(ret, onePrincipal)
		i = skipSpaces(cmd, i)
		if i >= len(cmd) {
			return nil, -1, getError("Unexpected EOF found", cmd, i)
		}
	}
	return ret, i, nil
}

func getOrPrincipals(cmd string, i int) ([][]string, int, error) {
	ret := [][]string{}
	andPrincipals, i, err := getAndPrincipals(cmd, i)
	if err != nil {
		return nil, -1, err
	}
	ret = append(ret, andPrincipals)
	i = skipSpaces(cmd, i)
	if i >= len(cmd) {
		return nil, -1, getError("Unexpected EOF found", cmd, i)
	}
	for cmd[i] == ',' {
		// Skip the comma
		i++
		andPrincipals, i, err = getAndPrincipals(cmd, i)
		if err != nil {
			return nil, -1, err
		}
		ret = append(ret, andPrincipals)
		i = skipSpaces(cmd, i)
		if i >= len(cmd) {
			return nil, -1, getError("Unexpected EOF found", cmd, i)
		}
	}
	return ret, i, nil
}

func getAndPrincipals(cmd string, i int) ([]string, int, error) {
	// Skip spaces front of and principals first
	i = skipSpaces(cmd, i)
	if i >= len(cmd) {
		return nil, -1, getError("Unexpected EOF found", cmd, i)
	}
	if cmd[i] != '(' {
		// No ( found, only one principal found
		onePrincipal, i, err := getPrincipal(cmd, i)
		return []string{onePrincipal}, i, err
	}

	i++
	if i >= len(cmd) {
		return nil, -1, getError("Unexpected EOF found", cmd, i)
	}

	// principals should be begin with ( and end with )
	i = skipSpaces(cmd, i)
	if i >= len(cmd) {
		return nil, -1, getError("Unexpected EOF found", cmd, i)
	}

	principals := []string{}
	// End of and principals
	var principal string
	var err error
	principal, i, err = getPrincipal(cmd, i)
	if err != nil {
		return nil, -1, err
	}
	if principal == "" {
		// This is an error, there isn't a principal between ()
		return nil, -1, getError("No principal found between ()", cmd, i)
	}
	principals = append(principals, principal)
	i = skipSpaces(cmd, i)
	if i >= len(cmd) {
		return nil, -1, getError("Unexpected EOF found", cmd, i)
	}

	for cmd[i] != ')' {
		i = skipSpaces(cmd, i)
		if i >= len(cmd) {
			return nil, -1, getError("Unexpected EOF found", cmd, i)
		}
		if cmd[i] != ',' {
			return nil, -1, getError("Principals should be sepereated by commas", cmd, i)
		}
		// read ,
		i++
		if i >= len(cmd) {
			return nil, -1, getError("Unexpected EOF found", cmd, i)
		}

		principal, i, err = getPrincipal(cmd, i)
		if err != nil {
			return nil, -1, err
		}
		principals = append(principals, principal)
		i = skipSpaces(cmd, i)
		if i >= len(cmd) {
			return nil, -1, getError("Unexpected EOF found", cmd, i)
		}
	}

	// read ')'
	i++
	return principals, i, nil
}

func getPrincipal(cmd string, i int) (string, int, error) {
	i = skipSpaces(cmd, i)

	var principal adsapi.Principal
	if (i+5 <= len(cmd)) && strings.EqualFold("user ", cmd[i:i+5]) {
		i += 5
		principal.Type = adsapi.PRINCIPAL_TYPE_USER
	} else if i+6 <= len(cmd) && strings.EqualFold("group ", cmd[i:i+6]) {
		i += 6
		principal.Type = adsapi.PRINCIPAL_TYPE_GROUP
	} else if i+5 <= len(cmd) && strings.EqualFold("role ", cmd[i:i+5]) {
		i += 5
		principal.Type = adsapi.PRINCIPAL_TYPE_ROLE
	} else if i+7 <= len(cmd) && strings.EqualFold("entity ", cmd[i:i+7]) {
		i += 7
		principal.Type = adsapi.PRINCIPAL_TYPE_ENTITY
	} else {
		return "", -1, getError("Not found principal type (user|group|role)", cmd, i)
	}

	principal.Name, i = getToken(cmd, i)
	if principal.Name == "" {
		return "", -1, getError("Not found principal name", cmd, i)
	}

	i = skipSpaces(cmd, i)
	if i+5 <= len(cmd) && strings.EqualFold("from ", cmd[i:i+5]) {
		// principal has idd with key word "from"
		i += 5
		i = skipSpaces(cmd, i)
		principal.IDD, i = getToken(cmd, i)
		if len(principal.IDD) == 0 {
			// no IDD found
			return "", -1, getError("No idd found after key word \"from\"", cmd, i)
		}
	}

	return subjectutils.EncodePrincipal(&principal), i, nil
}

func getRoles(cmd string, i int) ([]string, int, error) {
	tokens := []string{}
	t, i := getToken(cmd, i)
	if strings.EqualFold("role", t) {
		t, i = getToken(cmd, i)
	}
	if t == "" {
		return tokens, i, nil
	} 

	tokens = append(tokens, t)
	i = skipSpaces(cmd, i)
	for i < len(cmd) && cmd[i] == ',' {
		t, i = getToken(cmd, i+1)
		if strings.EqualFold("role", t) {
			t, i = getToken(cmd, i)
		}
		if t == "" {
			return nil, -1, getError("Not found role", cmd, i)
		}
		tokens = append(tokens, t)
		i = skipSpaces(cmd, i)
	}
	return tokens, i, nil
}

func getPermissions(cmd string, i int) ([]*pms.Permission, int, error) {
	perms := []*pms.Permission{}
	p, i, err := getPermission(cmd, i)
	if err != nil {
		return nil, -1, err
	}
	if p == nil {
		return perms, i, nil
	}
	perms = append(perms, p)
	i = skipSpaces(cmd, i)
	for i < len(cmd) && cmd[i] == ',' {
		p, i, err = getPermission(cmd, i+1)
		if err != nil {
			return nil, -1, err
		}
		if p == nil {
			return nil, -1, getError("Not found permission", cmd, i)
		}
		perms = append(perms, p)
		i = skipSpaces(cmd, i)
	}
	return perms, i, nil
}

func getPermission(cmd string, i int) (*pms.Permission, int, error) {
	acts, i, err := getTokens(cmd, i, "action")
	if err != nil {
		return nil, i, err
	}
	if len(acts) == 0 {
		return nil, i, getError("Not found permission", cmd, i)
	}
	res, i := getToken(cmd, i)
	if res == "" {
		return nil, i, getError("Not found permission", cmd, i)
	}
	isResExpr, resExpr := isResExpr(res)
	if isResExpr {
		return &pms.Permission{ResourceExpression: resExpr, Actions: acts}, i, nil
	} else {
		return &pms.Permission{Resource: res, Actions: acts}, i, nil
	}

}

func isResExpr(res string) (bool, string) {
	if strings.HasPrefix(res, res_expr_prefix) {
		return true, strings.TrimPrefix(res, res_expr_prefix)
	} else {
		return false, res
	}
}

func getResources(cmd string, i int) ([]string, []string, int, error) {
	i = skipSpaces(cmd, i)
	if i+3 <= len(cmd) && strings.EqualFold("on ", cmd[i:i+3]) {
		i += 3
		tokens, i, err := getTokens(cmd, i, "resource")
		if err != nil {
			return nil, nil, i, err
		}
		if len(tokens) == 0 {
			return nil, nil, -1, getError("Not found resource", cmd, i)
		}
		var resources, resExps []string
		for _, token := range tokens {
			isResExpr, resExp := isResExpr(token)
			if !isResExpr {
				resources = append(resources, token)
			} else {
				resExps = append(resExps, resExp)
			}
		}
		return resources, resExps, i, nil
	}
	return nil, []string{}, i, nil
}

func getService(cmd string, i int) (string, int, error) {
	i = skipSpaces(cmd, i)
	if i+3 <= len(cmd) && strings.EqualFold("in ", cmd[i:i+3]) {
		i += 3
		var serv string
		serv, i = getToken(cmd, i)
		if serv == "" {
			return "", -1, getError("Not found service", cmd, i)
		}
		return serv, i, nil
	}
	return "", i, nil
}

func getCondition(cmd string, i int) (string, int, error) {
	i = skipSpaces(cmd, i)
	if i+3 <= len(cmd) && strings.EqualFold("if ", cmd[i:i+3]) {
		i += 3
		if i >= len(cmd) {
			return "", -1, getError("Unexpected EOF found", cmd, i)
		}
		i = skipSpaces(cmd, i)
		if i >= len(cmd) {
			return "", -1, getError("Unexpected EOF found", cmd, i)
		}
		ret := cmd[i:]
		ret = strings.TrimRight(ret, " ")
		i = len(cmd)
		return ret, i, nil
	}
	return "", i, nil
}

func getTokens(cmd string, i int, e string) ([]string, int, error) {
	tokens := []string{}
	t, i := getToken(cmd, i)
	if t == "" {
		return tokens, i, nil
	}
	tokens = append(tokens, t)
	i = skipSpaces(cmd, i)
	for i < len(cmd) && cmd[i] == ',' {
		t, i = getToken(cmd, i+1)
		if t == "" {
			return nil, -1, getError(fmt.Sprintf("Not found %s", e), cmd, i)
		}
		tokens = append(tokens, t)
		i = skipSpaces(cmd, i)
	}
	return tokens, i, nil
}

func getToken(cmd string, i int) (string, int) {
	i = skipSpaces(cmd, i)
	if i >= len(cmd) {
		return "", i
	}
	end := cmd[i]
	if end == '"' || end == '\'' {
		i++
	}

	var buffer bytes.Buffer
	for ; i < len(cmd); i++ {
		if end == '"' || end == '\'' {
			if cmd[i] == end {
				i++
				break
			}
		} else if cmd[i] == ' ' || cmd[i] == ',' || cmd[i] == '(' || cmd[i] == ')' {
			break
		}
		buffer.WriteByte(cmd[i])
	}
	return buffer.String(), i
}

func skipSpaces(cmd string, i int) int {
	var j int
	var v rune
	for j, v = range cmd[i:] {
		if !unicode.IsSpace(v) {
			return i + j
		}
	}
	return i + j
}

func getErrorIndicator(cmd string, pos int) string {
	var buffer bytes.Buffer
	for i := 0; i < pos; i++ {
		buffer.WriteRune(' ')
	}
	buffer.WriteRune('^')
	return fmt.Sprintf("%s\n%s", cmd, buffer.String())
}

func getError(msg, cmd string, pos int) error {
	return fmt.Errorf("%s\n%s", msg, getErrorIndicator(cmd, pos))
}
