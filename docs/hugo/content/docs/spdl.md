+++
title = "SPDL - Security Policy Definition Language"
description = "Understand the basics of SPDL"
weight = 5
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["v0.1.0", "core", "policy", "spdl"]
categories = ["docs"]
bref = "Basics of SPDL"
+++

## Syntax

<img src="/img/speedle/spdl-syntax.png"/>

### Keywords

The reserved keywords in SPDL are as follows. You cannot use any of these keywords as user name, group name, action, resource, attribute name, and so on.

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

The keywords are all case-insensitive, which means that you cannot use one of "role", "ROLE", "Role", "rOLe", etc as user name, group name, action, resource, attribute name, and so on.

### Naming Convention

<pre>
User Name: a user name consists of letters, decimal digits, punctuation marks except for comma  i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Group Name: a group name consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Role Name: a role name consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Action: an action consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+

Resource: a resource consists of letters, decimal digits, punctuation marks i.e. [\p{L}\p{Nd}\p{Punct}]+
</pre>

Please see [Unicode Standard](http://www.unicode.org/reports/tr18/) and [Javadoc](https://docs.oracle.com/javase/7/docs/api/java/util/regex/Pattern.html) for the definition of letter, decimal digit, and punctuation mark.

### Syntax

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

## Condition

### 1. Overview

Condition is supported in both policy and role policy definitions. The policy or role policy can take effect _only_ when the condition is met.

This document focuses on what a condition is, and how to write a condition.

### 2. Condition

A condition is a bool expression that is constructed using attributes, functions, constants, operators, comparators or parenthesis and produces a bool value.

#### 2.1 Data Types

The data type of an attribute value and constant can be a string, numeric, bool, datetime or an array of string, numeric, bool, datetime.  
The data types and the corresponding supported operators and comparators are listed in following table.

 <table class="bordered striped">
    <thead>
      <tr>
        <th>Data Type</th>
        <th>Operators</th>
        <th>Comparators</th>
        <th>Comment</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>string</td>
        <td>+</td>
        <td>== <br> != <br> =~ <br> > <br> >= <br> < <br> <= <br>
        </td>
        <td>'+' is for string concatenation <br> '=~' is for regular expr matching, see the annotation below the table.</td>
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

Most comparators are self explanatory, here is the annotation for `"=~"`:  
`"=~"` is the regex comparator. It uses Go's standard regexp flavor of regex. The left side is expected to be the candidate string; the right side is the pattern.
`"=~"` returns whether or not the candidate string matches the regex pattern given on the right.

- Left side: string
- Right side: string
- Returns: bool

#### 2.2 Attributes

An attribute in a condition represents a variable. The attribute value is determined at runtime, obtained from either the Authorization Decision Service (ADS) runtime or the customer.
Built-in attributes are those values that are populated by the ADS runtime.
Customer attributes are those values that are provided by the customer when asking the ADS for authorization decisions.

##### 2.2.1 Built-in Attributes

Built-in attributes are as follows:

<table class="bordered striped">
    <thead>
      <tr>
        <th>Built-in Attribute Name</th>
        <th>Data Type </th>
        <th>Sample Value</th>
        <th>Definition</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>request_user</td>
        <td>string</td>
        <td>"Alice"</td>
        <td>The user information in the subject who is requesting to act on a resource</td>
      </tr>
      <tr>
        <td>request_groups</td>
        <td>[]string</td>
        <td>[]string{"managers"}</td>
        <td>The groups information in the subject who is requesting to act on a resource</td>
      </tr>
      <tr>
        <td>request_entity</td>
        <td>string</td>
        <td>"/org1/service1"</td>
        <td>The service name or other entity info in the subject which the service or program is requesting to act on a resource</td>
      </tr>
      <tr>
        <td>request_resource</td>
        <td>string</td>
        <td>"commercialLoans"</td>
        <td>The resource on which the subject is to act</td>
      </tr>
      <tr>
        <td>request_action</td>
        <td>string</td>
        <td>"issue" </td>
        <td>The action that the subject is to carry out on a resource </td>
      </tr>
      <tr>
        <td>request_time</td>
        <td>datetime</td>
        <td>'2019-01-02T15:04:05-07:00'</td>
        <td>The date and time when the request happens</td>
      </tr>
      <tr>
        <td>request_year</td>
        <td>int</td>
        <td>2019</td>
        <td>The year when the request happens</td>
      </tr>
      <tr>
        <td>request_month</td>
        <td>int</td>
        <td>1, 2, ... 12</td>
        <td>The month when the request happens</td>
      </tr>
      <tr>
        <td>request_day</td>
        <td>int</td>
        <td>1, 2, ... 31</td>
        <td>The day in a month when the request happens</td>
      </tr>
      <tr>
        <td>request_hour</td>
        <td>int</td>
        <td>0, 1, ... 23</td>
        <td>The hour in a day when the request happens</td>
      </tr>
      <tr>
        <td>request_weekday</td>
        <td>string</td>
        <td>"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"</td>
        <td>The day of the week when the request happens</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

##### 2.2.2 Customer Attributes

###### 2.2.2.1 name

A valid customer attribute name is a sequence of letters, digits or the underscore character `" _ "`. It must begin with a letter, and is limited to 255 characters. Customer attribute names must adhere to this rule, and should not use any reserved words (built-in attribute, built-in function name, operators, comparators) as attribute names.

###### 2.2.2.2 value

When customer attributes are used in a condition, the customer needs to pass the customer attribute values when requesting an isAllowed result.

- Passing attribute values in Golang API  
   Adhere to these rules when passing customer attribute values to a Golang API:

  - Use Golang bool type for bool attribute value.
  - Use Golang string type for string attribute value.
  - Use Golang float64 type for any numeric attribute value.
  - For datetime attribute value, use a Golang float64 representation of that datetime's Unix time (using `time.Time.Unix()`).
  - Use Golang []interface{} for array attribute value.

- Passing attribute values in REST API  
   Adhere to these rules when passing customer attribute values to a REST API:
  - Attribute is a struct, which contains name, type, value of the attribute. See the REST API for details.
  - Attribute type can be only "string", "numeric", "bool" or "datetime".
  - Attribute value can be a single value or a slice.

#### 2.3 Constants

Supported data types:

- string: single quotes, 'foobar'
- numeric: 10, 3.1415926
- bool: true or false
- datetime: single quotes, conform to RFC3339. Datetime of RFC3339 format is YYYY-MM-DDTHH:mm:SS[.sssssssss]Z, Z is [+|-]HH:mm.
- array: array of type string, numeric, bool, datetime

The following table shows sample constants.

<table class="bordered striped">
    <thead>
      <tr>
        <th>Data Type</th>
        <th>Constant Samples </th>
        <th>Array Constant Samples</th>
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

#### 2.4 Functions

Use functions when existing operators or comparators do not meet customer requirements.

##### 2.4.1 Built-in Functions

Speedle provides the following built-in functions.

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
        <td>Square root of a numeric</td>
        <td>One numeric parameter</td>
        <td>numeric</td>
        <td>Sqrt(x)<br>Sqrt(64)</td>
      </tr>
       <tr>
        <td>Max</td>
        <td>Get the max numeric in a set</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Max(1, 4, x)</td>
      </tr>
       <tr>
        <td>Min</td>
        <td>Get the min numeric in a set</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Min(x, 5, z)</td>
      </tr>
       <tr>
        <td>Sum</td>
        <td>Get the sum of a set of numeric</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Sum(1, 3, 5, 7, x)</td>
      </tr>
      <tr>
        <td>Avg</td>
        <td>Get the average value for a set of numeric</td>
        <td>1+ numeric parameters</td>
        <td>numeric</td>
        <td>Avg(x, 8, 10)</td>
      </tr>
      <tr>
        <td>IsSubSet</td>
        <td>Check if the first set/array is a subset of the second set/array</td>
        <td>2 sets/arrays, elements of the 2 sets/arrays have same data type</td>
        <td>bool</td>
        <td>IsSubset(s1, s2))</td>
      </tr>
    </tbody>
    <tfoot>
    </tfoot>
  </table>

##### 2.4.2 Custom Functions

Customers can also expose their own functions through a REST API, and use custom functions in a condition expression.  
For details, see [custom function](../custom-function/).

#### 2.5 Operator/Comparator Precedence

##### 2.5.1 Precedence order

When two operators share an operand, the operator with the higher precedence goes first. For example, `1 + 2 * 3` is treated as `1 + (2 * 3)`, whereas `1 * 2 + 3` is treated as `(1 * 2) + 3` because multiplication has a higher precedence than addition.

##### 2.5.2 Associativity

When an expression has two operators with the same precedence, the expression is evaluated according to its associativity. `72 / 2 / 3`is treated as `(72 / 2) / 3` because the `/` operator has left-to-right associativity. Some operators are not associative: for example, the expressions `(x <= y <= z)` and `x++--` are invalid.

##### 2.5.3 Precedence and Associativity of Supported Operators and Comparators

The following table lists all supported operators/comparators from highest to lowest precedence, along with their associativity.

<table class="bordered striped">
    <thead>
      <tr>
        <th>Precedence</th>
        <th>Operator/Comparator</th>
        <th>Description</th>
        <th>Associativity</th>
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

#### 2.6 Sample Conditions

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

### Condition Definition

Condition should be a valid bool expression. The bool expression Speedle supports is strictly defined as follows:

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

## Appendix

### Full Syntax of SPDL

<img src="/img/speedle/spdl-syntax-full.png"/>

<br/>
<br/>
