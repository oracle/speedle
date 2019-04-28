+++
title = "授权诊断"
description = "Authorization policy diagnosis"
date = 2019-01-18T15:46:59+08:00
weight = 3
draft = false
bref = ""
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
+++

授权诊断是 Speedle 的一个高级特性，用于诊断授权评估的过程。当一个用户被拒绝操作一个资源时，这个特性就可以用来找出拒绝的原因，例如，是哪条策略（policy）导致用户对资源的操作被拒绝的。

## 授权诊断响应包的格式

授权诊断响应包的格式如下：

```go
type EvaluationDebugResponse struct {
	Allowed       bool                    `json:"allowed"`
	Reason        string                  `json:"reason"`
	RequestCtx    *RequestContext         `json:"requestContext,omitempty"`
	Attributes    map[string]interface{}  `json:"attributes,omitempty"`
	GrantedRoles  []string                `json:"grantedRoles,omitempty"`
	RolePolicies  []*EvaluatedRolePolicy  `json:"rolePolicies,omitempty"`
	Policies      []*EvaluatedPolicy      `json:"policies,omitempty"`
}

type EvaluatedPolicy struct {
	Status      string              `json:"status,omitempty"`
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Effect      string              `json:"effect,omitempty"`
	Permissions []Permission        `json:"permissions,omitempty"`
	Principals  []string            `json:"principals,omitempty"`
	Condition   *EvaluatedCondition `json:"condition,omitempty"`
}

type EvaluatedRolePolicy struct {
	Status              string              `json:"status,omitempty"`
	ID                  string              `json:"id,omitempty"`
	Name                string              `json:"name,omitempty"`
	Effect              string              `json:"effect,omitempty"`
	Roles               []string            `json:"roles,omitempty"`
	Principals          []string            `json:"principals,omitempty"`
	Resources           []string            `json:"resources,omitempty"`
	ResourceExpressions []string            `json:"resourceExpression,omitempty"`
	Condition           *EvaluatedCondition `json:"condition,omitempty"`
}

type EvaluatedCondition struct {
	ConditionExpression  string  `json:"conditionExpression,omitempty"`
	EvaluationResult     string  `json:"evaluationResult,omitempty"`
}
```

## 示例

注意：原请求中的属性以及 Speedle 内置的属性都会包含在诊断响应包的属性列表中。

策略（policy）和角色策略（rolePolicy）中的字段"status"有三个有效的值，即：“takeEffect”、“conditionFailed”和"ignored"。

- takeEffect

"takeEffect"意思是策略（policy）或角色策略（rolePolicy）匹配并且已经被评估过了。

- conditioFailed

"conditionFailed"意思是策略（policy）或角色策略（rolePolicy）与请求中的服务名称（service name）、主体（subject）、资源（resource）以及操作（action）匹配，但是条件的评估结果是 false。

- ignored

"ignored"意思是授权评估过程已经结束，因此该策略（policy）已经没有评估的必要了。

下面就是授权诊断响应包的一个例子，

```go
{
  "Allowed": "true",
  "requestContext": {
    "subject": {
      "user": "user1",
      "groups": null,
      "attributes": null
    },
    "serviceName": "srv1",
    "resource": "res1",
    "action": "read",
    "attributes": null,
    "token": null
  },
  "attributes": {
    "request_action": "read",
    "request_day": 23,
    "request_groups": null,
    "request_month": "November",
    "request_resource": "res1",
    "request_time": 1511406017,
    "request_user": "user1",
    "request_weekday": "Thursday",
    "request_year": 2017
  },
  "grantedRoles": [
    "role1"
  ],
  "rolePolicies": [
    {
      "status": "takeEffect",
      "id": "c8087db3-60cf-4dad-aa9d-033eb6da0b15",
      "name": "rp01",
      "effect": "grant",
      "roles": [
        "role1"
      ],
      "principals": [
        "user:user1",
        "user:user2"
      ],
      "resources": [
        "res1",
        "res2"
      ],
      "condition": {

      }
    }
  ],
  "policies": [
    {
      "status": "takeEffect",
      "id": "f56b494f-dd6b-42af-962e-a109c890b7a0",
      "name": "p01",
      "effect": "grant",
      "permissions": [
        {
          "resource": "res1",
          "actions": [
            "list",
            "read",
            "write"
          ]
        },
        {
          "resource": "res2",
          "actions": [
            "list"
          ]
        }
      ],
      "principals": [
        "user:user1",
        "user:user2"
      ],
      "condition": {
        "conditionExpression": "request_year ==2017",
        "evaluationResult": "true"
      }
    }
  ]
}
```
