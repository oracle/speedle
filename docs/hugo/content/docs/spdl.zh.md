+++
title = "SPDL - 策略定义语言"
description = "Understand the basics of SPDL"
weight = 5
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["v0.1.0", "core", "policy", "spdl"]
categories = ["docs"]
bref = ""
+++

## 基本的策略定义语言

![基本策略定义语法](/img/speedle/spdl-syntax.png)

### 保留关键字

SPDL 的保留关键字如下. 这些关键字不能用作 user name, group name, action, resource, attribute name 等等。

- _role_
- _user_
- _group_
- _entity_
- _grant_
- _deny_
- _if_
- _in_
- _on_
- _from_

这些关键字均大小写不敏感。这意味着，"role", "ROLE", "Role", "rOLe" 都不能用作 user name, group name, action, resource, attribute name 等等。

### 命名规范

<pre>
User Name: a user name consists of letters, decimal digits, punctuation marks except for comma  i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Group Name: a group name consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Role Name: a role name consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Action: an action consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Resource: a resource consists of letters, decimal digits, punctuation marks i.e. [\p{L}\p{Nd}\p{Punct}]+
</pre>

Please see [Unicode Standard](http://www.unicode.org/reports/tr18/) and [Javadoc](https://docs.oracle.com/javase/7/docs/api/java/util/regex/Pattern.html) for the definition of letter, decimal digit, and punctuation mark.

### 语法

<pre>
POLICY = EFFECT SUBJECT ACTION RESOURCE if CONDITION
EFFECT = grant | deny
SUBJECT = AND_PRINCIPALS (, AND_PRINCIPALS)*
AND_PRINCIPALS = PRINCIPAL | \( PRINCIPAL_LIST \)
PRINCIPAL_LIST = PRINCIPAL (, PRINCIPAL)*
PRINCIPAL = PRINCIPAL_TYPE PRINCIPAL_NAME [PRINCIPAL_IDD]
PRINCIPAL_TYPE = user|group|entity|role
PRINCIPAL_IDD = from IDD_IDENTIFIER
IDD_IDENTIFIER = [\p{L}\p{Nd}\p{Punct}]+
ACTION = (ACTION_IDENTIFIER)(, ACTION_IDENTIFIER)*
RESOURCE = RESOURCE_IDENTIFIER
PRINCIPAL_NAME = [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
ACTION_IDENTIFIER = [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
RESOURCE_IDENTIFIER = [\p{L}\p{Nd}\p{Punct}]+
</pre>
<pre>
ROLE_POLICY = EFFECT SUBJECT ROLE (on RESOURCE)? if CONDITION
EFFECT = grant | deny
SUBJECT = PRINCIPAL (, PRINCIPAL)*
PRINCIPAL = PRINCIPAL_TYPE PRINCIPAL_NAME [PRINCIPAL_IDD]
PRINCIPAL_TYPE = user|group|entity|role
PRINCIPAL_IDD = from IDD_IDENTIFIER
IDD_IDENTIFIER = [\p{L}\p{Nd}\p{Punct}]+
ROLE = (role)? SUBJECT_IDENTIFIER
RESOURCE = RESOURCE_IDENTIFIER
SUBJECT_IDENTIFIER = [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
RESOURCE_IDENTIFIER = [\p{L}\p{Nd}\p{Punct}]+
</pre>

## 条件（Condition）

### 1. 概述

Policy 和 role policy 都支持 Condition。只有 condition 满足了， policy 或 role policy 才会生效。

这一小节主要介绍什么是 condition， 以及如何构建 condition。

### 2. Condition

Condition 就是一个布尔表达式。由属性(attributes), 函数(functions), 常量(constants), 操作符(operators), 比较运算符(comparators) or 括号(parenthesis)构建的布尔表达式。

#### 2.1 数据类型 （Data Types）

属性和常量的数据类型可以是 string, numeric, bool, datetime，或者由 string, numeric, bool, datetime 构成的数组。  
数据类型及其支持的操作符，比较运算符如下表所示：

 <table class="bordered striped">
    <thead>
      <tr>
        <th>数据类型<br>Data Type</th>
        <th>操作符<br>Operators</th>
        <th>比较运算符<br>Comparators</th>
        <th>备注<br>Comment</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>string</td>
        <td>+</td>
        <td>== <br> != <br> =~ <br> > <br> >= <br> < <br> <= <br>
        </td>
        <td>'+' 用于字符串连接操作 <br><br> '=~' 用于正则表达式匹配。  <br>左边是匹配的字符串， 右边是正则表达式.`"=~"`返回ture如果匹配成功，或者false如果匹配失败。</td>
      </tr>
      <tr>
        <td>numeric</td>
        <td> + <br> - <br> * <br> / <br> %</td>
        <td>== <br> != <br> =~ <br> > <br> >= <br> < <br> <= </td>
        <td></td>
      </tr>
      <tr>
        <td>bool</td>
        <td>&& <br> || <br> !</td>
        <td>== <br>!=</td>
        <td></td>
      </tr>
      <tr>
        <td>datetime</td>
        <td></td>
        <td>== <br> != <br> > <br> >= <br> < <br> <=</td>
        <td></td>
      </tr>
      <tr>
        <td>array</td>
        <td></td>
        <td>in</td>
        <td>membership 'in' operator: left side should be a single type(string, numeric, bool, datetime), right side should be an array</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

#### 2.2 属性(Attributes)

属性(attribute)代表一个变量。属性分为内置属性和用户属性两大类。
内置属性(Built-in attributes)是 Speedle 预定义的，它们的值是在决策运算中由 Authorization Decision Service (ADS)运行时指定的。
用户属性（customer attributes)的值是用户在授权请求(authorization decision request)中传入的。

##### 2.2.1 内置属性(Built-in Attributes)

Speedle 预定义的内置属性如下:

<table class="bordered striped">
    <thead>
      <tr>
        <th>内置属性名<br>Built-in Attribute Name</th>
        <th>数据类型<br>Data Type </th>
        <th>例子<br>Sample Value</th>
        <th>定义<br>Definition</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>request_user</td>
        <td>string</td>
        <td>"Alice"</td>
        <td>Authorization请求中(Subject)的user信息</td>
      </tr>
      <tr>
        <td>request_groups</td>
        <td>[]string</td>
        <td>[]string{"managers"}</td>
        <td>Authorization请求中(Subject)的groups信息</td>
      </tr>
      <tr>
        <td>request_entity</td>
        <td>string</td>
        <td>"/org1/service1"</td>
        <td>Authorization请求中(Subject)的entity信息</td>
      </tr>
      <tr>
        <td>request_resource</td>
        <td>string</td>
        <td>"commercialLoans"</td>
        <td>Authorization请求中的资源</td>
      </tr>
      <tr>
        <td>request_action</td>
        <td>string</td>
        <td>"issue" </td>
        <td>Authorization请求中对资源的操作</td>
      </tr>
      <tr>
        <td>request_time</td>
        <td>datetime</td>
        <td>'2019-01-02T15:04:05-07:00'</td>
        <td>Authorization请求时的日期和时间</td>
      </tr>
      <tr>
        <td>request_year</td>
        <td>int</td>
        <td>2019</td>
        <td>Authorization请求时的年份</td>
      </tr>
      <tr>
        <td>request_month</td>
        <td>int</td>
        <td>1, 2, ... 12</td>
        <td>Authorization请求时的月份</td>
      </tr>
      <tr>
        <td>request_day</td>
        <td>int</td>
        <td>1, 2, ... 31</td>
        <td>Authorization请求时是一个月中的哪一天</td>
      </tr>
      <tr>
        <td>request_hour</td>
        <td>int</td>
        <td>0, 1, ... 23</td>
        <td>Authorization请求时是一天中的哪个时辰</td>
      </tr>
      <tr>
        <td>request_weekday</td>
        <td>string</td>
        <td>"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"</td>
        <td>Authorization请求时是星期几</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

##### 2.2.2 用户属性(Customer Attributes)

###### 2.2.2.1 属性名(name)

一个合法的属性名由字母，数字和下划线`" _ "`组成，必须以字母开头，小于 255 个字符，且不能用保留关键字。

###### 2.2.2.2 属性值(value)

当用户属性应用于 condition 时, 当用户向 ADS 发送 Authorization Decision 请求时，用户需要将属性值随请求一并传入。

- 通过 Golang API 传入属性值  
   需遵循如下规则：

  - bool 型属性值使用 Golang bool 类型.
  - string 型属性值使用 Golang string 类型.
  - numeric 型属性值使用 Golang float64 类型.
  - datetime 型属性值使用 Golang float64 类型, 也就是 Unix time (using `time.Time.Unix()`).
  - 数组属性使用 Golang []interface{}.

- 通过 REST API 传入属性值  
   须遵循如下规则：
  - 属性是一个结构体(struct),包含属性名(name), 属性的数据类型(type), 属性值(value). 详细信息参见 REST API.
  - 属性的数据类型(type)只能是 "string", "numeric", "bool" or "datetime".
  - 属性值可以是单个值， 也可以是数组.

#### 2.3 常量(Constants)

支持的数据类型:

- string: single quotes, 'foobar'
- numeric: 10, 3.1415926
- bool: true or false
- datetime: single quotes, conform to RFC3339. Datetime of RFC3339 format is YYYY-MM-DDTHH:mm:SS[.sssssssss]Z, Z is [+|-]HH:mm.
- array: array of type string, numeric, bool, datetime

各种常量的例子如下表所示：

<table class="bordered striped">
    <thead>
      <tr>
        <th>数据类型<br>Data Type</th>
        <th>常量例子<br>Constant Samples </th>
        <th>常量数组例子<br>Array Constant Samples</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>string</td>
        <td>'a string'</td>
        <td>('string1', 'string2')</td>
      </tr>
      <tr>
        <td>bool</td>
        <td>true <br> false</td>
        <td>(true, false, true)</td>
      </tr>
      <tr>
        <td>numeric / float</td>
        <td>2.1</td>
        <td>(1, 2, 3.1)</td>
      </tr>
      <tr>
        <td>datetime</td>
        <td>'2006-01-02T15:04:05-07:00'</td>
        <td>('2016-01-02T15:04:05-07:00', '2019-01-02T15:04:05-07:00')</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

#### 2.4 函数(Functions)

当简单的运算或比较操作符不满足需求时，可以使用函数。函数也分为内置函数和用户自定义函数。

##### 2.4.1 内置函数(Built-in Functions)

Speedle 提供以下内置函数：

<table class="bordered striped">
    <thead>
      <tr>
        <th>Built-in Function Name</th>
        <th>Functionality </th>
        <th>Input and Data Type </th>
        <th>Output and Data Type</th>
        <th>Sample Usage</th>
      </tr>
    </thead>
    <tbody>
       <tr>
        <td>Sqrt</td>
        <td>求平方根</td>
        <td>One numeric parameter</td>
        <td>numeric</td>
        <td>Sqrt(x)<br>Sqrt(64)</td>
      </tr>
       <tr>
        <td>Max</td>
        <td>取集合中的最大值</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Max(1, 4, x)</td>
      </tr>
       <tr>
        <td>Min</td>
        <td>取集合中的最小值</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Min(x, 5, z)</td>
      </tr>
       <tr>
        <td>Sum</td>
        <td>求和</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Sum(1, 3, 5, 7, x)</td>
      </tr>
      <tr>
        <td>Avg</td>
        <td>求平均值</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Avg(x, 8, 10)</td>
      </tr>
      <tr>
        <td>IsSubSet</td>
        <td>第一个参数是否是第二个参数的子集</td>
        <td>2 sets/arrays, elements of the 2 sets/arrays have same data type</td>
        <td>bool</td>
        <td>IsSubset(s1, s2))</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

##### 2.4.2 用户自定义函数(Custom Functions）

用户可以向 Speedle 暴露自己定义的函数, 并将自定义函数用于 condition.  
更多细节, 参见 [custom function](../custom-function/).

#### 2.5 运算比较操作符的优先级(operator/Comparator Precedence)

##### 2.5.1 优先顺序(Precedence order)

当两个运算符共享一个操作数时，优先级较高的运算符优先。例如, `1 + 2 * 3` 被处理成 `1 + (2 * 3)`, 但是 `1 * 2 + 3` 被处理成 `(1 * 2) + 3`。 因为乘法比加法的优先级高。

##### 2.5.2 关联(Associativity)

当表达式具有两个具有相同优先级的运算符时，将根据其关联性来计算表达式。 `72/2/3`被视为`（72/2）/ 3`，因为`/`运算符具有从左到右的关联性。 有些运算符不是关联的：例如，表达式`（x <= y <= z）`和`x ++ -`无效。

##### 2.5.3 Precedence and Associativity of Supported Operators and Comparators

下表按优先级列出了所有运算和比较操作符及其关联性.

<table class="bordered striped">
    <thead>
      <tr>
        <th>优先级<br>Precedence</th>
        <th>运算/比较操作符<br>Operator/Comparator</th>
        <th>描述<br>Description</th>
        <th>关联<br>Associativity</th>
      </tr>
    </thead>
    <tbody>
    <tr>
        <td>7</td>
        <td>( )</td>
        <td>parentheses</td>
        <td>N/A</td>
      </tr>
    <tr>
        <td>6</td>
        <td>function call</td>
        <td></td>
        <td>N/A</td>
      </tr>
       <tr>
        <td>5</td>
        <td>* <br>/<br>%<br></td>
        <td></td>
        <td>Left to right</td>
      </tr>
       <tr>
        <td>4</td>
        <td>+<br>-<br></td>
        <td>1+ numeric parameters</td>
        <td>Left to right</td>
      </tr>
       <tr>
        <td>3</td>
        <td>!</td>
        <td></td>
        <td>N/A</td>
      </tr>
       <tr>
        <td>2</td>
        <td>&&</td>
        <td></td>
        <td>Left to right</td>
      </tr>
      <tr>
        <td>1</td>
        <td>||</td>
        <td></td>
        <td>Left to right</td>
      </tr>
      <tr>
        <td>0</td>
        <td>==<br>!=<br>=~<br>><br>>=<br><<br><=<br>in<br></td>
        <td></td>
        <td>N/A</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

#### 2.6 Condition 示例

<table class="bordered striped">
    <thead>
      <tr>
        <th>Condition</th>
        <th>Comment</th>
      </tr>
    </thead>
    <tbody>
    <tr><td>a=='abc'</td><td> Value of attribute 'a' equals 'abc'</td></tr>
    <tr><td>a!='abc'</td><td></td></tr>
    <tr><td>a>='abc'</td><td></td></tr>
    <tr><td>a+b=='ab'</td><td></td></tr>
    <tr><td>a=~'\^get.*'</td><td>Value of attribute 'a' matches regular expression '\^get.*'</td></tr>
<tr><td>a=123</td><td></td></tr>
<tr><td>a-b>123</td><td></td></tr>
<tr><td>a in (1, 2, 3) </td><td>Value of attribute 'a' is one of 1, 2, 3 </td></tr>
<tr><td>'manager' in a </td><td>'a' is an attribute of string array, 'manager' is one of the array element </td></tr>
<tr><td>IsSubSet(e, ('s1', 's2', 's3'))</td><td>'e' is an attribute of string array/set, and e is subset of array/set ('s1', 's2', 's3')</td></tr>
<tr><td>a in (1, 2, 3) && (b==c &#124;&#124; d==3) && IsSubSet(e, ('s1', 's2', 's3'))</td><td></td></tr>
<tr><td>request_year==2019 && request_month==12 </td><td> request_year and request_month are built-in attributes. The year when a resource is accessed equals 2019, and the month when the resource is accessed is 12</td></tr>
 </tbody>
    <tfoot>
    </tfoot>
  </table>

### Condition 定义

Condition 必须是一个合法的布尔表达式。 Speedle 支持的布尔表达式严格定义如下:

```
BoolExpr: ('!')BoolExpr
          | BoolExpr ('&&'|'||') BoolExpr
		  | BoolConstant
		  | Attribute
		  | Function
		  | RelationalExpr
		  | '(' BoolExpr ')'


RelationalExpr: NumericExpr ('=='|'!='|'>'|'>='|'<'|'<=') NumericExpr
              | StringExpr ('=='|'!='|'=~'|'>'|'>='|'<'|'<=') StringExpr
			  | BoolExpr ('=='|'!=') BoolExpr
			  | DateTimeExpr ('=='|'!='|'>'|'>='|'<'|'<=') DateTimeExpr
			  | (NumericExpr|StringExpr|BoolExpr|DateTimeExpr) ('in') ArrayExpr


NumericExpr: NumericExpr('+'|'-'|'*'|'/'|'%')NumericExpr
           | NumericConstant
		   | Attribute
		   | Function
		   | '(' NumericExpr ')'

StringExpr: StringExpr('+') StringExpr
           | StringConstant
		   | Attribute
		   | Function
           | '(' StringExpr ')'

DateTimeExpr: DateTimeConstant
           | Attribute
           | Function

ArrayExpr: ArrayConstant
         | Attribute
         | Function

ArrayConstant: '('Constant [, Constant]* ')'

Constant: NumericConstant
         |StringConstant
         |BoolConstant
         |DateTimeConstant

NumericConstant: any numeric float64 data
StringConstant: single quoted string, for example, 'string1'
BoolConstant: true|false
DateTimeConstant: single quoted datetime, datetime should conform to the format defined by RFC3339: YYYY-MM-DDTHH:mm:SS[.sssssssss]Z, Z is [+|-]HH:mm. For example, '2016-01-02T15:04:05-07:00'

Attribute: attribute name should conform to  [a-zA-Z]+[a-zA-Z0-9_]*, length should be <=255
Function: FuntionName '(' Argument [,Argument]* ')'

```

## 附录

### 完整策略定义语法

![完整策略定义语法](/img/speedle/spdl-syntax-full.png)

<br/>
<br/>
