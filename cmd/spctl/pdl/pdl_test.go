//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package pdl

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/oracle/speedle/api/pms"
)

func TestParsingPolicy(t *testing.T) {
	cmd := `grant user yufyu,group Developers role1, role2 on r1,r2, r3 in k8s if 'i > 30 && j < 4','t >= "2012-05-06"'`
	v, _, _ := ParseRolePolicy(cmd, "p1")
	fmt.Println(v)
}

func TestEffect(t *testing.T) {
	testCases := []struct {
		cmd  string
		want string
	}{
		{
			cmd:  "grant user ...",
			want: grant,
		},
		{
			cmd:  "deny user ...",
			want: deny,
		},
		{
			cmd:  "GRANT user ...",
			want: grant,
		},
		{
			cmd:  "DENY user ...",
			want: deny,
		},
		{
			cmd:  "Grant user ...",
			want: grant,
		},
		{
			cmd:  "Deny user ...",
			want: deny,
		},
		{
			cmd:  "  Grant user ...",
			want: grant,
		},
		{
			cmd:  "  Deny user ...",
			want: deny,
		},
		{
			cmd: `		Grant user ...`,
			want: grant,
		},
		{
			cmd: `		Deny user ...`,
			want: deny,
		},
		{
			cmd:  "\tGrant user ...",
			want: grant,
		},
		{
			cmd:  "\tDeny user ...",
			want: deny,
		},
	}

	for _, tc := range testCases {
		got, i, err := getEffect(tc.cmd)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if got != tc.want {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}
		if !strings.HasPrefix(tc.cmd[i:], "user") {
			t.Errorf("cmd: %s, output cmd: %s, should be %s", tc.cmd, tc.cmd[i:], "user")
		}
	}
}

func TestEffectNeg(t *testing.T) {
	testCases := []struct {
		cmd string
	}{
		{
			cmd: "",
		},
		{
			cmd: " ",
		},
		{
			cmd: "grent user ...",
		},
		{
			cmd: "Dany user ...",
		},
		{
			cmd: "grantuser ...",
		},
		{
			cmd: "denyuser ...",
		},
		{
			cmd: "grant\tuser ...",
		},
		{
			cmd: "deny\tuser ...",
		},
		{
			cmd: "grant, user ...",
		},
		{
			cmd: "deny, user ...",
		},
	}

	for _, tc := range testCases {
		_, _, err := getEffect(tc.cmd)
		if err == nil {
			t.Errorf("cmd: %s, no error", tc.cmd)
		}
	}
}

func TestPrincipals(t *testing.T) {
	testCases := []struct {
		cmd  string
		want [][]string
	}{
		{
			cmd:  "user yufei.yu list ...",
			want: [][]string{{"user:yufei.yu"}},
		},
		{
			cmd:  "   user      yufei.yu list ...",
			want: [][]string{{"user:yufei.yu"}},
		},
		{
			cmd:  "group Developers list ...",
			want: [][]string{{"group:Developers"}},
		},
		{
			cmd:  "   group      Developers list ...",
			want: [][]string{{"group:Developers"}},
		},
		{
			cmd:  "role Dev list ...",
			want: [][]string{{"role:Dev"}},
		},
		{
			cmd:  "   role     Dev list ...",
			want: [][]string{{"role:Dev"}},
		},
		{
			cmd:  "USER yufei.yu list ...",
			want: [][]string{{"user:yufei.yu"}},
		},
		{
			cmd:  "   USER      yufei.yu list ...",
			want: [][]string{{"user:yufei.yu"}},
		},
		{
			cmd:  "Group Developers list ...",
			want: [][]string{{"group:Developers"}},
		},
		{
			cmd:  "   Group      Developers list ...",
			want: [][]string{{"group:Developers"}},
		},
		{
			cmd:  "RoLE Dev list ...",
			want: [][]string{{"role:Dev"}},
		},
		{
			cmd:  "   RoLE     Dev list ...",
			want: [][]string{{"role:Dev"}},
		},
		{
			cmd:  `user "Yufei Yu" list ...`,
			want: [][]string{{"user:Yufei Yu"}},
		},
		{
			cmd:  `user 'Yufei Yu' list ...`,
			want: [][]string{{"user:Yufei Yu"}},
		},
		{
			cmd:  `user yufei.yu, user william.cai, user jiexiang.fu list ...`,
			want: [][]string{{"user:yufei.yu"}, {"user:william.cai"}, {"user:jiexiang.fu"}},
		},
		{
			cmd:  `user yufei.yu from cisco, user william.cai, user jiexiang.fu from oracle list ...`,
			want: [][]string{{"idd=cisco:user:yufei.yu"}, {"user:william.cai"}, {"idd=oracle:user:jiexiang.fu"}},
		},
		{
			cmd:  `user 'Yufei Yu', user 'William Cai', user 'Bill Fu' list ...`,
			want: [][]string{{"user:Yufei Yu"}, {"user:William Cai"}, {"user:Bill Fu"}},
		},
		{
			cmd:  `user 'Yufei Yu',group Developers  ,   role 'Dev' list ...`,
			want: [][]string{{"user:Yufei Yu"}, {"group:Developers"}, {"role:Dev"}},
		},
		{
			cmd:  `(user william) list ...`,
			want: [][]string{{"user:william"}},
		},
		{
			cmd:  `(user william,group employee) list ...`,
			want: [][]string{{"user:william", "group:employee"}},
		},
		{
			cmd:  `     (     user      william     ,     group     employee     )     list ...`,
			want: [][]string{{"user:william", "group:employee"}},
		},
		{
			cmd:  `(user william,group employee),role finance list ...`,
			want: [][]string{{"user:william", "group:employee"}, {"role:finance"}},
		},
		{
			cmd:  `   (    user    william    ,   group   employee   )   ,   role    finance    list ...`,
			want: [][]string{{"user:william", "group:employee"}, {"role:finance"}},
		},
		{
			cmd:  `role finance,(user william,group employee) list ...`,
			want: [][]string{{"role:finance"}, {"user:william", "group:employee"}},
		},
		{
			cmd:  `    role    finance   ,   (   user    william   ,   group    employee   )    list ...`,
			want: [][]string{{"role:finance"}, {"user:william", "group:employee"}},
		},
		{
			cmd:  `(role finance),(user william,group employee) list ...`,
			want: [][]string{{"role:finance"}, {"user:william", "group:employee"}},
		},
		{
			cmd:  `  (   role    finance  )  ,   (   user    william   ,   group    employee   )    list ...`,
			want: [][]string{{"role:finance"}, {"user:william", "group:employee"}},
		},
		{
			cmd:  `(role finance,role erp),(user william,group employee) list ...`,
			want: [][]string{{"role:finance", "role:erp"}, {"user:william", "group:employee"}},
		},
		{
			cmd:  `  (   role    finance    ,    role   erp    )  ,   (   user    william   ,   group    employee   )    list ...`,
			want: [][]string{{"role:finance", "role:erp"}, {"user:william", "group:employee"}},
		},
	}

	for _, tc := range testCases {
		got, i, err := getOrPrincipals(tc.cmd, 0)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if len(got) != len(tc.want) {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
			break
		}

		if !strings.HasPrefix(tc.cmd[i:], "list") {
			t.Errorf("cmd: %s, output cmd: %s, should be %s", tc.cmd, tc.cmd[i:], "list")
		}
	}
}

func TestRoles(t *testing.T) {
	testCases := []struct {
		cmd  string
		want []string
	}{
		{
			cmd:  "Role1 on ...",
			want: []string{"Role1"},
		},
		{
			cmd:  "       Role1 on ...",
			want: []string{"Role1"},
		},
		{
			cmd:  "Role1,Role2,Role3 on ...",
			want: []string{"Role1", "Role2", "Role3"},
		},
		{
			cmd:  "   Role1,   Role2    ,Role3   ,  Role4 on ...",
			want: []string{"Role1", "Role2", "Role3", "Role4"},
		},
		{
			cmd:  "role Role1 on ...",
			want: []string{"Role1"},
		},
		{
			cmd:  "   role Role1,   role Role2    , role Role3   ,  role Role4 on ...",
			want: []string{"Role1", "Role2", "Role3", "Role4"},
		},
	}

	for _, tc := range testCases {
		got, i, err := getRoles(tc.cmd, 0)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if len(got) != len(tc.want) {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}
		for j := range got {
			if got[j] != tc.want[j] {
				t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
				break
			}
		}

		if !strings.HasPrefix(tc.cmd[i:], "on") {
			t.Errorf("cmd: %s, output cmd: %s, should be %s", tc.cmd, tc.cmd[i:], "on")
		}
	}
}

func TestPermissions(t *testing.T) {
	testCases := []struct {
		cmd  string
		want []pms.Permission
	}{
		{
			cmd:  "get core/pods ...",
			want: []pms.Permission{{Actions: []string{"get"}, Resource: "core/pods"}},
		},
		{
			cmd:  "   get core/pods ...",
			want: []pms.Permission{{Actions: []string{"get"}, Resource: "core/pods"}},
		},
		{
			cmd:  "get,list,watch core/pods ...",
			want: []pms.Permission{{Actions: []string{"get", "list", "watch"}, Resource: "core/pods"}},
		},
		{
			cmd:  "get, list , watch core/pods ...",
			want: []pms.Permission{{Actions: []string{"get", "list", "watch"}, Resource: "core/pods"}},
		},
		{
			cmd:  "get, list , watch 'res with whitespace' ...",
			want: []pms.Permission{{Actions: []string{"get", "list", "watch"}, Resource: "res with whitespace"}},
		},
	}

	for _, tc := range testCases {
		got, i, err := getPermissions(tc.cmd, 0)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if len(got) != len(tc.want) {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}
	outer:
		for j := range got {
			if len(got[j].Actions) != len(tc.want[j].Actions) {
				t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got[j].Actions, tc.want[j].Actions)
				break
			}
			for k := range got[j].Actions {
				if got[j].Actions[k] != tc.want[j].Actions[k] {
					t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got[j].Actions, tc.want[j].Actions)
					break outer
				}
			}
			if got[j].Resource != tc.want[j].Resource {
				t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got[j].Resource, tc.want[j].Resource)
				break
			}
		}

		if !strings.HasPrefix(tc.cmd[i:], "...") {
			t.Errorf("cmd: %s, output cmd: %s, should be %s", tc.cmd, tc.cmd[i:], "...")
		}
	}
}

func TestResources(t *testing.T) {
	testCases := []struct {
		cmd  string
		want []string
	}{
		{
			cmd:  "ON c1/default/core/pods/* ...",
			want: []string{"c1/default/core/pods/*"},
		},
		{
			cmd:  "on c1/default/core/pods/* ...",
			want: []string{"c1/default/core/pods/*"},
		},
		{
			cmd:  "ON       c1/default/core/pods/* ...",
			want: []string{"c1/default/core/pods/*"},
		},
		{
			cmd:  "       ON       c1/default/core/pods/* ...",
			want: []string{"c1/default/core/pods/*"},
		},
		{
			cmd:  "ON c1/default/core/pods/*,c2/test/core/pods/* ...",
			want: []string{"c1/default/core/pods/*", "c2/test/core/pods/*"},
		},
		{
			cmd:  "ON r1, r2,r3 ...",
			want: []string{"r1", "r2", "r3"},
		},
	}

	for _, tc := range testCases {
		got, _, i, err := getResources(tc.cmd, 0)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if len(got) != len(tc.want) {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
				break
			}
		}
		if !strings.HasPrefix(tc.cmd[i:], "...") {
			t.Errorf("cmd: %s, output cmd: %s, should be %s", tc.cmd, tc.cmd[i:], " ...")
		}
	}
}

func TestService(t *testing.T) {
	testCases := []struct {
		cmd  string
		want string
	}{
		{
			cmd:  "IN k8s ...",
			want: "k8s",
		},
		{
			cmd:  "in k8s ...",
			want: "k8s",
		},
		{
			cmd:  "IN       k8s ...",
			want: "k8s",
		},
		{
			cmd:  "       IN       k8s ...",
			want: "k8s",
		},
	}

	for _, tc := range testCases {
		got, i, err := getService(tc.cmd, 0)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if got != tc.want {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}

		if !strings.HasPrefix(tc.cmd[i:], " ...") {
			t.Errorf("cmd: %s, output cmd: %s, should be %s", tc.cmd, tc.cmd[i:], " ...")
		}
	}
}

func TestConditions(t *testing.T) {
	testCases := []struct {
		cmd  string
		want string
	}{
		{
			cmd:  "if a=3",
			want: "a=3",
		},
		{
			cmd:  "IF a=3",
			want: "a=3",
		},
		{
			cmd:  "if        a=3      ",
			want: "a=3",
		},
		{
			cmd:  "if a = 3",
			want: "a = 3",
		},

		{
			cmd:  "if a = 3 &&   b == 4     ",
			want: "a = 3 &&   b == 4",
		},
	}

	for _, tc := range testCases {
		got, _, err := getCondition(tc.cmd, 0)
		if err != nil {
			t.Errorf("cmd: %s, error: %v", tc.cmd, err)
		}
		if len(got) != len(tc.want) {
			t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
		}
		for j := range got {
			if got[j] != tc.want[j] {
				t.Errorf("cmd: %s, got %v, want %v", tc.cmd, got, tc.want)
				break
			}
		}
	}
}

func TestFullCmd(t *testing.T) {
	cmd := "grant user  user_bool_equal1   get,del res_equal1 if x == false"
	_, _, err := ParsePolicy(cmd, "test")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
