# SPDL Policy Definition language

# Syntax

<img src="../img/spdl-syntax.png"/> 

## Keywords

Here is a list of keywords in PDL. You cannot use any of the following as user name, group name, action, resource, attribute name, etc.  

*role* 
*user*
*group*
*entity*
*grant*
*deny*
*if*
*in*
*on*
*from*

The keywords are all case-insensitive. That means you cannot use one of "role", "ROLE", "Role", "rOLe", etc as user name, group name, action, resource, attribute name, etc.

## Naming Convention

`User Name: a user name consists of letters, decimal digits, punctuation marks except for comma  i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+ 
Group Name: a group name consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
Role Name: a role name consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
Action: an action consists of letters, decimal digits, punctuation marks except for comma i.e. [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
Resource: a resource consists of letters, decimal digits, punctuation marks i.e. [\p{L}\p{Nd}\p{Punct}]+
`

Please refer to [Unicode Standard](http://www.unicode.org/reports/tr18/) and [Javadoc](https://docs.oracle.com/javase/7/docs/api/java/util/regex/Pattern.html) for the definition of letter, decimal digit, and punctuation mark

## Syntax  

`POLICY = EFFECT SUBJECT ACTION RESOURCE if CONDITION
EFFECT = grant | deny
SUBJECT = AND_PRINCIPALS (, AND_PRINCIPALS)*
AND_PRINCIPALS = PRINCIPAL | \( PRINCIPAL_LIST \)
PRINCIPAL_LIST = PRINCIPAL (, PRINCIPAL)*
PRINCIPAL = PRINCIPAL_TYPE PRINCIPAL_NAME [PRINCIPAL_IDD]
PRINCIPAL_TYPE = user|group|entity|role
PRINCIPAL_IDD = from IDD_IDENTIFIER
IDD_IDENTIFIER = [\p{L}\p{Nd}\p{Punct}]+
ACTION = (ACTION_IDENTIFIER)(,ACTION_IDENTIFIER)*
RESOURCE = RESOURCE_IDENTIFIER
PRINCIPAL_NAME = [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
ACTION_IDENTIFIER = [\p{L}\p{Nd}[\p{Punct}&&[^,]]]+
RESOURCE_IDENTIFIER = [\p{L}\p{Nd}\p{Punct}]+
`

`ROLE_POLICY = EFFECT SUBJECT ROLE (on RESOURCE)? if CONDITION
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
`

# Condition

## 1. Overview    
Condition is supported in policy or role policy definition. Only when condition is met, the policy or role policy could take effect.     

This doc focus on what is supported in condition and how to write conditions.

## 2. Condition    
Condition is a bool expression which is constructed by attributes, functions, constants, operators, comparators or parenthesis and produces a bool value.   

### 2.1 Data Types:
Data type of attribute value and constant could be string, numeric, bool, datetime or array of string, numeric, bool,datetime.    
The data types and the corresponding supported operators and comparators are listed in following table. 

|Data Type| Operators | Comparators | Comment |
|--|--|--|--|
|string|+|==<br>  !=<br>  =~<br>  ><br>  >=<br>  <<br>  <=<br>| '+' here is for string concatenation<br> '=~' is for regular expr matching, more details see below explaination<br> |
|numeric|+<br> -<br> *<br> /<br> %<br>|==<br>  !=<br>  ><br>  >=<br>  <<br>  <=<br>| |
|bool|&&<br> &#124;&#124;<br> !<br>|==<br>  !=<br>| |
|datetime||==<br>  !=<br>  ><br>  >=<br>  <<br>  <=<br>| |
|array||in|membership in operator: left side shoule be a single type(string, numeric, bool, datetime), right side should be an array |


Most comparators are self explained, here is the explain for "=~", the regex comparator:    
"=~" uses go's standard regexp flavor of regex. The left side is expected to be the candidate string, the right side is the pattern. 
"=~" returns whether or not the candidate string matches the regex pattern given on the right.
- Left side: string
- Right side: string
- Returns: bool


### 2.2 Attributes
Attribute in condition stands for variable, value of which is known at runtime, either given by ADS runtime or given by customer.
Built-in attributes are those values of which are given by ADS runtime.
Customer attributes are those values of which are given by customer when asking ADS for authorization decisions.

#### 2.2.1 built-in attributes    
Built-in attributes are listed as follows:   
Think about a scenario in a bank, somebody is issuing commercialLoans, before he/she could issue commercial loans, authorization check is carrried out...

| built-in attribute name | data type | sample value | meaning |
|--|--|--|--|
| request_user| string| "Alice" | the user info of the subject who is requesting to act on a resource|
| request_groups| []string | []string{"managers"}| the groups info of the subject who is requesting to act on a resource|
| request_resource| string | "commercialLoans" | the resource on which the subject is to act |
| request_action| string | "issue" | the action that the subject is to carry out on a resource |
| request_time| datetime |'2017-01-02T15:04:05-07:00' | the datetime when the request happens |
| request_year|int |2017 | year when the request happens |
| request_month|int | 1,2,...12| month when the request happens |
| request_day|int |1,2,...31 | day in the month when the request happens |
| request_hour|int |0,1,...23|hour in a day when the request happens|
| request_weekday|string |"Sunday", "Monday", "Tuesday","Wednesday","Thursday",	"Friday","Saturday", |weekday when the request happens|


#### 2.2.2 customer attributes   
##### 2.2.2.1 name    
A leagal customer attribute name is a limited-length(255) sequence of letters, digits or the underscore character " _ ", beginning with a letter. Customer attribute name should strictly follow the rule, and should avoid using any reserved words(built-in attribute, built-in function name, operators, comparators) as attribute names.

##### 2.2.2.2 value    
When customer attribute is used in condition, customer need to pass customer attributes value when requesting isAllowed result.
- passing attributes values in Golang API    
    When passing customer attributes values to Golang API, pls follow the following rules :
    - pls use Golang bool type for bool attribute value.
    - pls use Golang string type for string attribute value.
    - pls use Golang float64 type for any numeric attribute value.
    - For datetime attribute value, pls use a Golang float64 representation of that datetime's unix time (using time.Time.Unix()).
    - Pls use Golang []interface{} for array attribute value.

- passing attributes values in REST API        
    When passing customer attributes values to REST API, pls follow the following rule:   
    - TODO



### 2.3 Constants
Supported data types:    
- string: single quotes, 'foobar'
- numeric: 10, 3.1415926
- bool: true or false
- datetime: single quotes, conform to RFC3339. Datetime of RFC3339 format is YYYY-MM-DDTHH:mm:SS[.sssssssss]Z, Z is [+|-]HH:mm.
- array: array of type string,numeric,bool,datetime

Following are samples constants:    

|Data Type| Constant Samples| Array Constant Samples |
|--|--|--|
|string|'a string'|('string1','string2')|
|bool|true<br>false|(true,false,true)|
|numeric/float64|1<br>2.1|(1,2,3.1)|
|datetime|'2006-01-02T15:04:05-07:00'|('2016-01-02T15:04:05-07:00','2017-01-02T15:04:05-07:00')|



### 2.4 Functions
Existing operators or comparators may not meet customer requirements for some use cases. We also provide some built-in functions.
Only built-in functions could be used in condition.   

| built-in function name | functionality | input & data type | output & data type| sample usage |
|--|--|--|--|--|
|Sqrt| square root of a numeric| one numeric parameter| numeric |Sqrt(x)<br>Sqrt(64) |
|Max | get the max numeric in a set | 1+ numeric parameters| numeric | Max(1,4,x)|
|Min | get the min numeric in a set |1+ numeric parameters | numeric| Min(x,5,z)|
|Sum | get the sum of a set of numeric |1+ numeric parameters | numeric| Sum(1,3,5,7,x) |
|Avg | get the average value for a set of numeric |1+ numeric parameters | numeric |Avg(x,8,10) |
|IsSubSet | check if the first set/array is a subset of the second set/array | 2 sets/arrays, elements of the 2 sets/arrays have same data type| bool| IsSubset(s1, s2))|

### 2.5 operator/comparator precedence　　　　  

#### 2.5.1 Precedence order.    
When two operators share an operand the operator with the higher precedence goes first. For example, 1 + 2 * 3 is treated as 1 + (2 * 3), whereas 1 * 2 + 3 is treated as (1 * 2) + 3 since multiplication has a higher precedence than addition.

#### 2.5.2 Associativity.     
When an expression has two operators with the same precedence, the expression is evaluated according to its associativity.  72 / 2 / 3 is treated as (72 / 2) / 3 since the / operator has left-to-right associativity. Some operators are not associative: for example, the expressions (x <= y <= z) and x++-- are invalid.

#### 2.5.3 Precedence and associativity of supported operarors and comparartors.     
The table below shows all supported operators/comparators from highest to lowest precedence, along with their associativity.

|precedence|operator/comparator|description|associativity|
|--|--|--|--|
| 7|() |parentheses |N/A|
| 6| function call| |N/A|
| 5|*<br>/<br>%<br> | |left to right|
| 4|+<br>-<br>| |left to right|
| 3|!| |N/A|
| 2|&& | |left to right|
| 1| &#124;&#124;| |left to right|
| 0|==<br>!=<br>=~<br>><br>>=<br><<br><=<br>in<br>|  |N/A|

### 2.6 Sample conditions
|condition|comment|
|--|--|
|a=='abc'| value of attribute 'a' equals 'abc'|
|a!='abc'| |
|a>='abc'| |
|a+b=='ab'| |
|a=~'\^get.*'| value of attribute 'a' matches regular expression '\^get.*'|
|a=123| |
|a-b>123| |
|a in (1, 2, 3)| value of attribute 'a' is one of 1,2,3 |
|'manager' in a| 'a' is an attribute of string array, 'manager' is one of the array element |
|IsSubSet(e, ('s1','s2','s3'))| 'e' is an attribute of string array/set, and e is subset of array/set ('s1','s2','s3')|
|a in (1,2,3) && (b==c &#124;&#124; d==3) && IsSubSet(e, ('s1','s2','s3'))| |
|request_year==2017 && request_month==12| request_year and request_month are buitl-int attributes. The year when a resource is accessed equals 2017, and the month when the resource is accesed is 12|


## Appendix: Condition Definition
Condition should be a valid bool expression. And bool expression we support is strictly defined as follows:
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
