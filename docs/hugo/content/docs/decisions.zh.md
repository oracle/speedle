+++
title = "授权查询"
description = "Get authorization decisions for your service interactions"
weight = 30
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["pdp", "policy", "core"]
categories = ["docs"]
bref = ""
+++

## 1. 什么是授权查询?

- 授权查询是 Speedle ADS(Authorization Decision Service)提供的服务接口， 一般用于查询某个主体(subject)对某个资源(resource)实施某项操作(action)是否被允许。

- 授权查询的结果是基于角色策略(role-policies)和策略(policies)的实时运算。

## 2. 授权查询的方式

Speedle 支持以下 3 种方式进行授权查询：

- REST API provided by Authorization Decision Service(ADS)
- Grpc API provided by Authorization Decision Service(ADS)
- Golang API

## 3. 授权查询 API 及其示例

The ADS decision APIs make authorization decisions based on policies that describe the actions, permissions, and roles granted to a subject.

### 3.1 查询授权决定

查询某个主体(subject)对某个资源(resource)实施某项操作(action)是否被允许

- API overview
  - IN
    - Given the request: subject, action, resource
    - Given the runtime attributes \*\*optional\*\*
    - Given the service scope
  - OUT
    - Returns _true_ if allowed, _false_ if _NOT_ allowed
    - Returns reason for the decision
    - Returns errors if an error occurs
- Sample
  - 查询 user Alan 从 onlineBookStore 应用 下载 HarryPotter 这本书是否被允许。
  - 授权结果基于定义在 "onlineBookStore" 这个 service 中的所有角色策略(role-policies)和策略(policies)的。

**REST API example:**

_Request:_

```
curl -X POST  http://localhost:6734/authz-check/v1/is-allowed \
-d @- << EOF
{
 "subject": {"principals":[{"type":"user", "name":"Alan"}]},
 "action": "download",
 "resource":"/books/HarryPotter",
 "serviceName": "onlineBookStore"
}
EOF
```

_Response:_

```
{"allowed":true,"reason":0}
```

这里 reason '0'表示 ADS 找到了授权策略. 下表列出了所有原因的定义:

 <table class="bordered striped">
    <thead>
      <tr>
        <th>原因<br>Reason</th>
        <th>定义<br>Definition</th>
        <th>含义<br>Comment</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td> 0 </td>
        <td> GRANT_POLICY_FOUND </td>
        <td> 找到了授权策略 </td>
      </tr>
      <tr>
        <td> 1 </td>
        <td> DENY_POLICY_FOUND </td>
        <td> 找到了拒绝授权策略 </td>
      </tr>
      <tr>
        <td> 2 </td>
        <td> SERVICE_NOT_FOUND </td>
        <td> 没找到服务 </td>
      </tr>
      <tr>
        <td> 3 </td>
        <td> NO_APPLICABLE_POLICIES </td>
        <td> 没找到匹配的策略 </td>
      </tr>
      <tr>
        <td> 4 </td>
        <td> ERROR_IN_EVALUATION </td>
        <td> 策略运算中出现错误 </td>
      </tr>
      <tr>
        <td> 5 </td>
        <td> DISCOVER_MODE </td>
        <td> 处于Discovery Mode </td>
      </tr>
   </tbody>
 </table>

### 3.2 查询某一主体(subject)的所有角色(Roles)

取得某一主体(subject)的所有角色(roles)

- API overview

  - IN
    - Given the subject
    - Given the runtime attributes \*\*optional\*\*
    - Given the service scope
  - OUT
    - Returns a slice of roles granted to current subject
    - Returns errors if an error occurs

- Sample
  - 取得 user Alan 被授予的所有角色(roles)
  - 结果基于定义在 "onlineBookStore" 这个 service 中的所有角色策略(role-policies)。

**REST API example:**  
_Request:_

```
curl -X POST  http://localhost:6734/authz-check/v1/all-granted-roles \
-d @- << EOF
{
 "subject": {"principals":[{"type":"user", "name":"Alan"}]},
 "serviceName": "onlineBookStore"
}
EOF
```

_Response:_

```
["role1", "role2"]
```

### 3.3 查询某一主体(subject)被授予的所有权限(Permissions)

取得授予某一主体(subject)的所有的权限(permissions).

- API overview

  - IN
    - Given the subject
    - Given the runtime attributes \*\*optional\*\*
    - Given the service scope
  - OUT
    - Returns a slice of (actions, resource) pairs, current subject is allowed to perform.
    - Returns errors if an error occurs

- Sample
  - 取得授予 user Alan 的所有的权限(permissions).
  - 结果基于定义在 "onlineBookStore" 这个 service 中的所有角色策略(role-policies)和策略(policies)。

**REST API example:**  
_Request:_

```
curl -X POST  http://localhost:6734/authz-check/v1/all-granted-permissions \
-d @- << EOF
{
 "subject": {"principals":[{"type":"user", "name":"Alan"}]},
 "serviceName": "onlineBookStore"
}
EOF
```

_Response:_

```
[{
    "resource":"/books/HarryPotter",
    "actions":["download","read"]
 },
 {
    "resource":"/books/ThreeBodyProblem",
    "actions":["borrow"]
 }]
```

For details, see [Authorization Runtime/Decision API](../api/decision_api).
