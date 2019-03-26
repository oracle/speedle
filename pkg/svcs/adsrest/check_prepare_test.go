//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

//+build runtime_test_prepare

package adsrest

import (
	"testing"

	"github.com/oracle/speedle/testutil"
)

//Prepare service/policy/rolepolicy in store for runtime_test
func TestPrepareCheckEnvironment(t *testing.T) {

	cc := testutil.NewCmdClient()

	serviceSimplePDL := []string{
		"role-policies:",
		"grant user user_condtion1 role_condition1 if age > 10",
		"grant user  userWithRole1  role1 on res_allow",
		"grant user  userWithRole1  role1 on res_deny",
		"deny  group groupWithRole2 role2 on res_deny",
		"policies:",
		"grant user  user1  get,del res_allow",
		"grant group group1 get,del res_allow",
		"grant role  role1  get,del res_allow",

		"deny  user  user1  get,del res_deny",
		"deny  group group1 get,del res_deny",
		"deny  role  role1  get,del res_deny",
		"grant role  role2  get,del res_deny",
	}

	if !cc.CreateServiceWithPDL("service-simple", serviceSimplePDL) {
		t.Fatalf("Create service(service-simple) failed via pdl ")
	}

	serviceBothGrantDenyPDL := []string{
		"role-policies:",
		"grant user  userWithRole1     role1 on res1",
		"deny  group groupWithoutRole1 role1 on res1",
		"policies:",
		"deny  user  userWithRole1     get,del res1",
		"grant role  role1             get,del res1",
		"grant group groupWithoutRole1 get,del res1",

		"grant user  user_allowed,group group_allowed  get,del res1",
		"deny user  user_denied, group group_denied    get,del res1",
	}

	if !cc.CreateServiceWithPDL("service-both-grant-deny", serviceBothGrantDenyPDL) {
		t.Fatalf("Create service(service-both-grant-deny) failed via pdl ")
	}

	serviceForCondFuncPDL := []string{
		"policies:",
		//x=4 deny; x=6 deny, x=8 allow
		"grant user  user_cond_sqrt1    get,del res1 if Sqrt(16)<x",
		"deny  user  user_cond_sqrt1    get,del res1 if Sqrt(64)>x",

		//x=4 deny; x=6 allow, x=8 deny
		"grant user  user_cond_sqrt2   get,del res1 if Sqrt(16) < x && Sqrt(64) > x",

		//s1=16, s2=4 allow
		//s1=16.0, s2=4 allow
		//s1=16, s2=4.00 allow
		//s1=-16, s2=-4.00 allow
		//s1=-16, s2=4.00 deny
		"grant user  user_cond_sqrt3   get,del res1 if Sqrt(s1) < s2",

		//s1=1, s2=4 allow
		//s1=6, s2=9 deny
		//s1=10.00, s2=12.2 allow
		"grant user  user_cond_sum1   get,del res1 if Sum(s1,s2) < 10 || Sum(s1,s2) > 20",

		//s1=1, s2=4 allow
		//s1=6, s2=9 deny
		//s1=12.55, s2=8.46 allow
		"grant user  user_cond_avg1   get,del res1 if Avg(s1,s2) < 5 || Avg(s1,s2) > 10",

		//s1=1, s2=4 allow
		//s1=6, s2=9 deny
		//s1=11.55, s2=12 allow
		"grant user  user_cond_min1   get,del res1 if Min(s1,s2) < 5.5 || Min(s1,s2) > 10.5",

		//s1=5.4, s2=4.4 deny
		//s1=6.4, s2=3.4 allow
		"grant user  user_cond_max1   get,del res1 if !(Max(s1,s2) < 5.5)",

		//s1=4, s2=3 deny
		//s1=4, s2=5.0 allow
		"grant user  user_cond_max2   get,del res1 if Max(s1,s2) == 5",

		//s1=('1','2'), s2=('2','1') allow
		//s1=('1','2'), s2=('2','1', '3') allow
		//s1=('1','2','3'), s2=('2','1') allow
		//s1=('1'), s2=('2','1') allow
		//s1=('1'), s2=('1') allow
		"grant user  user_subset_str   get,del res1 if IsSubSet(s1,s2)",

		//s1=(1,2), s2=(2,1) allow
		//s1=(1,2), s2=(2.0,1.000,3) allow
		//s1=(1.00,2.0), s2=(2,1,3) allow
		//s1=(1,2), s2=(2,1) allow
		//s1=(1), s2=(1) allow
		//s1=(1,2), s2=(1) deny
		//s1=(1,2), s2=(1,3) deny
		"grant user  user_subset_int   get,del res1 if IsSubSet(s1,s2)",

		//s1=('2017-01-02T15:04:05-07:00'), s2=('2017-01-02T15:04:05-07:00') allow
		//s1=('2017-01-02T15:04:05-07:00'), s2=('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00') allow
		//s1=('2017-01-02T15:04:05-19:00'), s2=('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00') deny
		//s1=('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00'), s2=('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00') allow
		//s1=('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00'), s2=('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-10:00') deny
		"grant user  user_subset_datetime   get,del res1 if IsSubSet(s1,s2)",
	}
	if !cc.CreateServiceWithPDL("service-condition-func", serviceForCondFuncPDL) {
		t.Fatalf("Create service(service-condition-func) failed via pdl ")
	}

	serviceForCondCustomAttriPDL := []string{
		"policies:",
		//------------------------Boolean-------------------------
		//x=false allow; x=true deny; x=abc error
		"grant user  user_bool_equal1   get,del res_equal1 if x == false",

		//x=false deny; x=true allow; x=abc error
		"grant user  user_bool_equal2   get,del res_equal2 if !x == false",

		//x=false allow; x=true allow; x=abc error
		"grant user  user_bool_notequal   get,del res_equal2 if x != true",

		//s1=true,s2=false,s3=true allow
		//s1=true,s2=true,s3=true allow
		//s1=false,s2=false,s3=true allow
		//s1=true,s2=true,s3=false deny
		//s1=false,s2=true,s3=false deny
		//s1=abc,s2=true,s3=abc error
		//s1=False,s2=True,s3=False error
		"grant user  user_bool_complex   get,del res_complex if s1 && !s2 || s3 == true",

		//s1=true,s2=false,s3=true allow.
		//Could create success if use &&!, but fail to parse it when reload the file
		//"grant user  user_bool1   get,del res1 if s1&&!s2 || s3 == true",
		"grant user  user_bool1   get,del res1 if s1 && !s2 || s3 == true",

		//s1=true,s2=false,s3=true allow
		"grant user  user_bool2   get,del res1 if s1 && (!s2)||s3 == true",

		//s1=true,s2=false,s3=false allow
		"grant user  user_bool3   get,del res1 if s1 && (!s2||s3) == true",

		//------------------------Numberic-------------------------
		//s1=20, s2=2, s3=3  allow
		//s1=20, s2=2, s3=2  deny
		//s1=20, s2=2, s3=-3  deny
		//s1=-22, s2=2, s3=-3  allow
		//s1=-22, s2=2, s3=3  deny
		//s1=12, s2=-2, s3=3  allow
		//s1=abc, s2=-2, s3=3  error
		"grant user  user_num_equal   get,del res_num1 if (s1+5-s2*2)/3%4 == s3",

		//s1=-10, s2=3, s3=-1 allow
		//s1=-10, s2=3, s3=1 deny
		//s1=0, s2=-3, s3=1 allow
		//s1=0, s2=-3, s3=-1 deny
		//s1=1.2, s2=-3.0, s3=1 allow
		"grant user  user_num_equal2   get,del res_num2 if 5*(s1+5)%s2 == s3",

		//s1=0, s2=10, deny
		//s1=1, s2=10, allow
		//s1=abc, s2=10, error
		"grant user  user_num_notequal  get,del res_num3 if (s1+5)*2 != s2",

		//s1=-5, s2=0, deny
		//s1=-2, s2=0, allow
		//s1=0, s2=0, allow
		//s1=-10, s2=-10, allow
		"grant user  user_num_greater   get,del res_num4 if (s1+5)%2 > s2",

		//s1=-5, s2=0, allow
		//s1=-20, s2=-20, allow
		//s1=30, s2=10.02, allow
		"grant user  user_num_greater_equal   get,del res_num5 if (s1+5)/2.0 >= s2",

		//s1=1, s2=0.4 deny
		//s1=1, s2=0.5 allow
		//s1="1", s2=".5" error
		//s1="a", s2=".5" error
		"grant user  user_num_greater_equal2  get,del res_num6 if s1+s2 >= 1.41421356237309504880168872420969807856967187537694807317667974",

		//------------------------String-------------------------------
		//s1='1', s2='23' allow
		//s1='1', s2='22' deny
		//s1='', s2='123' allow
		//s1=null, s2='123' allow
		//s1=1, s2=23 error
		"grant user  user_str_equal1  get,del res_str1 if s1+s2 == '123'",

		//s1='1', s2='1' deny
		//s1='1', s2='-1' allow
		//s1="1", s2=1 error
		//s1=nil, s2='nil' allow
		//s1="", s2='nil' allow
		"grant user  user_str_notequal  get,del res_str2 if s1 != s2",

		//s1=abc, s2=abc allow
		//s1=abc, s2="123" deny
		"grant user  user_str_regx  get,del res_str3 if s1+s2=~'.*abc$'",

		//s1='abc', s2='Abc' allow
		//s1='abc', s2='abc' deny
		//s1='abc', s2='' allow
		"grant user  user_str_greater  get,del res_str4 if s1>s2",

		//s1='abc', s2='Abc' allow
		//s1='abc', s2='abc' allow
		//s1='null', s2='' allow
		"grant user  user_str_greater_equal  get,del res_str5 if s1>=s2",

		//s1='abc', s2='Abc' deny
		//s1='abc', s2='abc' allow
		//s1='null', s2='' deny
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-02T15:04:05-07:00' allow
		"grant user  user_str_less_equal  get,del res_str6 if s1<=s2",

		//s1='abc', s2='Abc' deny
		//s1='abc', s2='abc' deny
		//s1='null', s2='' deny
		//s1='a bc', s2='a1bc' allow
		"grant user  user_str_less  get,del res_str7 if s1<s2",

		//------------------------Datetime-------------------------------
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-02T15:04:05-09:00' deny
		//s1='2017-01-02T15:04:05-07:00', s1='2017-01-22T15:04:05-07:00' deny
		//s1='2017-02-29T15:04:05-07:00', s1='2017-01-22T15:04:05-07:00' error
		"grant user  user_date_equal  get,del res1 if s1==s2",

		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' deny
		"grant user  user_date_notequal  get,del res1 if s1!=s2",

		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' deny
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' deny
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00' allow
		"grant user  user_date_greater  get,del res1 if s1>s2",

		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' deny
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00' allow
		"grant user  user_date_greater_equal  get,del res1 if s1>=s2",

		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' allow
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00' deny
		"grant user  user_date_less_equal  get,del res1 if s1<=s2",

		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00' deny
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-09:00' allow
		//s1='2017-01-02T15:04:05-09:00', s2='2017-01-02T15:04:05-07:00' deny
		"grant user  user_date_less  get,del res1 if s1 < s2",

		//------------------------Array-------------------------------
		//x='a' allow
		//x='b' deny
		//x='1' allow
		//x='' deny
		//x=1 error
		"grant user  user_array_str  get,del res1 if x in ('a','c','1')",

		//x=1 allow
		//x=0 deny
		"grant user  user_array_num_single  get,del res1 if x in (1)",

		//x=1 allow
		//x=2 allow
		//x=2.0 allow
		//x=-2.1 allow
		//x=3.5 deny
		//x=3.567 allow
		//x='1' error
		"grant user  user_array_num  get,del res_in2 if x in (1,2.0,-2.1,3.567)",

		//x='a' allow
		//x='b' deny
		"grant user  user_array_str_single  get,del res1 if x in ('a')",
		"grant user  user_array_str_single1  get,del res1 if s1 in s2",

		//s1='1', s2='1', s3='2' allow
		//s1='1', s2='2', s3='2' deny
		//s1=1, s2=1, s3=2 allow
		//s1=1, s2=2, s3=2 deny
		//s1=1.1, s2=1.1, s3=2 allow
		//s1=-2.0, s2=-2, s3=3 allow
		//s1=-2, s2=-2.0, s3=3 allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-07:00', s3='2017-01-02T15:04:05-09:00' allow
		//s1='2017-01-02T15:04:05-07:00', s2='2017-01-02T15:04:05-10:00', s3='2017-01-02T15:04:05-09:00' deny
		//s1=('1','2'),s2=('1','5'),s3=('2','6'),allow
		//s1=(1,2),s2=(1,5),s3=(2,6),allow
		//s1=('2017-01-02T15:04:05-07:00','2017-01-02T15:04:05-17:00'),s2=('2017-01-02T15:04:05-07:00','2017-01-02T15:04:05-09:00'),s3=('2017-01-02T15:04:05-09:00', '2017-01-02T15:04:05-17:00'),allow
		"grant user  user_array_str_multi  get,del res1 if s1 in (s2, s3)",

		//x='2017-01-02T15:04:05-07:00' allow
		//x='2017-01-02T15:04:05-17:00' deny
		"grant user  user_array_date  get,del res_in2 if x in ('2017-01-02T15:04:05-07:00', '2017-01-02T15:04:05-09:00')",

		//x='2017-01-02T15:04:05-07:00' allow
		//x='2017-01-02T15:04:05-17:00' deny
		"grant user  user_array_date_single  get,del res_in2 if x in ('2017-01-02T15:04:05-07:00')",

		"grant user  user_array_issubset1  get,del res_subset if IsSubSet(s1,s2)",

		"grant user  user_array_issubset2  get,del res_subset if IsSubSet(s1,('a1','b1','c1'))",
	}
	if !cc.CreateServiceWithPDL("service-condition-attribute", serviceForCondCustomAttriPDL) {
		t.Fatalf("Create service(service-condition-attribute) failed via pdl ")
	}

	serviceForBuiltinAttributePDL := []string{
		"policies:",
		"grant user  user_attri1    get,del res_request_time1 if request_time > '2017-09-04 12:00:00'",

		//user=admin, groups=group_attri1 allowed; user=admin1, groups=group1 deny.   user=admin, groups=group2 deny
		"grant group  group_attri1    get,del res_request_user1 if request_user == 'admin'",

		//user=admin, groups="group1","group2" allowed; user=admin1, groups="group1","group5" deny.   user=admin, groups=group2 deny
		"grant group  group1    get,del res_request_groups1 if  IsSubSet(request_groups, ('group1','group2','group3'))",

		//user=user1, actions=wget grant; user=user1, actions=wget grant; user=user1, actions=getx deny;
		"grant user  user_attri1    wget,uget,del res_request_action1 if request_action =~ 'get'",

		//"request_resource =~ '^/node1/*?'",
		//user=user1, actions=wget grant; user=user1, actions=wget grant; user=user1, actions=getx deny;
		"grant user  user_attri1   get,del res_request_resource1  if request_resource =~ 'resource1'",
		"grant user  user_attri1   get,del res_request_resource2  if request_resource =~ 'resource1'",

		//"request_weekday == 'Monday'",
		"grant user  user_attri1   get,del res_request_weekday1 if request_weekday == 'Monday'",

		//"request_year == 2017",
		"grant user  user_attri1   get,del res_request_year_equal_2017 if request_year == 2017",
		"grant user  user_attri1   get,del res_request_year_greater_2017 if request_year >= 2017",

		//"request_month == 'September'",
		"grant user  user_attri1   get,del res_request_month_nov if request_month == 'November'",

		//"request_day == 14",
		"grant user  user_attri1   get,del res_request_day_14 if request_month == 14",
	}
	if !cc.CreateServiceWithPDL("service-builtin-attribute", serviceForBuiltinAttributePDL) {
		t.Fatalf("Create service(service-builtin-attribute) failed via pdl ")
	}

	serviceComplexPDL := []string{
		"role-policies:",
		"grant user user_complex1, user user_complex1A,user user_complex1B role_complex1",
		"policies:",
		"grant role role_complex1 get,del res_complex1 if request_user != 'user_complex1'",
		"deny user user_complex1A del res_complex1",
	}
	if !cc.CreateServiceWithPDL("service-complex", serviceComplexPDL) {
		t.Fatalf("Create service(service-complex) failed via pdl ")
	}

	//role3---role2---role1---user1
	//			|		|-----user11
	//			|		|-----group1
	//			|
	//			|-----user2
	//			|-----user22
	//			|-----group2
	//
	//user1 get res1: allowed
	//user1 del res2: allowed
	//user1 del res3: allowed
	//user11 del res3: denied
	//user11,group1 del res3: denied
	//user11,group1 get res3: allowed
	//userAny,group1 get res3: allowed
	//user2 get res2: allowed
	//user2 get res1: allowed
	//user2 del res1: denied
	//user2 get res3: allowed
	//user22 get res1: denied
	//user22 del res3: denied
	//user22 get res3: allowed
	//user1 get res9: allowed

	serviceComplexRoleEmbededPDL := []string{
		"role-policies:",
		"grant user user1, user user11, group group1 role1",
		"grant user user2, user user22, group group2, role role1 role2",
		"grant role role2 role3",
		"grant role role3 role4",
		"grant role role4 role5",
		"grant role role5 role6",
		"grant role role6 role7",
		"grant role role7 role8",
		"grant role role8 role9",
		"grant user user11, user user22 role-denined",
		"grant role role-denined role-denined1",
		"deny user user22 role5",

		"grant user userRes1 role10  on res1",
		"grant user userRes1 role11 if request_action == 'get'",
		"grant role role10 role12 on res1",
		"grant role role10 role12 on res2",

		"policies:",
		"grant role role1 get,del res1",
		"grant role role2 get,del res2",
		"grant role role3 get,del res3",
		"grant role role9 get,del res9",
		"deny role role-denined1 del res3",
		"grant user user2 get res1",
	}
	if !cc.CreateServiceWithPDL("service-complex-role", serviceComplexRoleEmbededPDL) {
		t.Fatalf("Create service(service-complex-role) failed via pdl ")
	}

	//user1 get res: allowed
	//user1 get res 2: allowed
	//user1 get res*: allowed
	//user2 get res*: allowed
	//user2 get res2*: allowed
	//user2 get *res2*: allowed
	//user2 get 22res222: allowed
	//user2 get ?res?: allowed
	//user1 get res-denied*: denied
	//user1 get res-denied: denied
	//user1 get 11res-denied11: denied
	//user2 get 11res-denied11: denied
	serviceWithResExprPDL := []string{
		//https://golang.org/pkg/regexp/syntax
		"role-policies:",
		"grant user user2 role2 on expr:.*res.*",
		"policies:",
		"grant user user1 get,del expr:res*",
		"grant role role2 get,del expr:.*res.*",
		"deny user user1, user user2 get,del expr:.*res-denied*",
	}
	if !cc.CreateServiceWithPDL("service-with-resexpr", serviceWithResExprPDL) {
		t.Fatalf("Create service(service-with-resexpr) failed via pdl ")
	}

	//no test case yet
	serviceWithNegResExprPDL := []string{
		//https://golang.org/pkg/regexp/syntax
		"policies:",
		"grant user user1 get,del expr:*res*",
	}
	if !cc.CreateServiceWithPDL("service-with-neg-resexpr", serviceWithNegResExprPDL) {
		t.Fatalf("Create service(service-with-neg-resexpr) failed via pdl ")
	}

	serviceWithComplexPrincipleInPDL := []string{
		//https://golang.org/pkg/regexp/syntax
		"policies:",
		"grant (user user1, user user11, group group1), (user user2, group group2, group group22), (group group3, group group33) get,del res1",
		"grant (user user1, group group1), (user user2, group group2, group group22) get,del res2",
	}
	if !cc.CreateServiceWithPDL("service-with-complex-principle-policy", serviceWithComplexPrincipleInPDL) {
		t.Fatalf("Create service(service-with-complex-principle-policy) failed via pdl ")
	}

	serviceWithComplexPrincipleInRolePolicyPDL := []string{
		//https://golang.org/pkg/regexp/syntax
		"role-policies:",
		"grant (user user1, user user11, group group1), (user user2, group group2, group group22), (group group3, group group33) role1",
		"grant role role1 role2",
		"policies:",
		"grant role role1 get,del res1",
		"grant (group group4, role role2),(user user4, role role2) get,del res4",
		"grant group group5, role role2 get,del res5",
	}
	if !cc.CreateServiceWithPDL("service-with-complex-principle_rolePolicy", serviceWithComplexPrincipleInRolePolicyPDL) {
		t.Fatalf("Create service(service-with-complex-principle_rolePolicy) failed via pdl ")
	}

	serviceWithEntityPrinciplePDL := []string{
		//https://golang.org/pkg/regexp/syntax
		"role-policies:",
		"grant entity schema://domain.name/path1, entity schema://domain.name/path2 role1",
		//"grant (entity spiffe://domain.name/path1, entity spiffe://domain.name/ns/user1) role2", //multi entities is not supported when do arsrest

		"policies:",
		"grant role role1  get,del res1",
		"grant (group group1, entity spiffe://acme.com/9eebccd2-12bf-40a6-b262-65fe0487d453), role role1 get,del res2",
		"deny entity schema://domain.name/path2 get,del res2",
		"grant (entity special-schema.1+2://user1:pwd@domain1/path1/path-2/a), role role1 get,del res3",
	}
	if !cc.CreateServiceWithPDL("service-with-entity-principle", serviceWithEntityPrinciplePDL) {
		t.Fatalf("Create service(service-with-entity-principle) failed via pdl ")
	}

}
