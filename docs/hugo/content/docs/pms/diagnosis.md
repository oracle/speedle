+++
title = "Policy Diagnosis"
description = "Authorization policy diagnosis"
date = 2019-01-18T15:46:59+08:00
weight = 3
draft = false
bref = "This feature is used to diagnose the evaluation process of authorization. When a user isn't allowed to operate on a resource, then this feature is useful to find out the reason, i.e., what's the policy which denies the operation on the resource."
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
+++

## What's policy diagnosis

Policy diagnosis is one of the Speedle's advanced features, which is used to diagnose the evaluation process of authorization. When the evaluation result isn't expected or a user isn't allowed to operate on a resource, then this feature is useful to find out the reason, i.e., what's the policy which denies the operation on the resource.

## How to use policy diagnosis

The usage of policy diagnosis request is almost identical to the usage of authorization decision request. They have exactly the same parameters in function calls or request payload formats. The only difference is that they have different API names, please check the following table for detailed info,

| API Type                         | Authorization Decision     | Policy Diagnosis         |
| -------------------------------- | -------------------------- | ------------------------ |
| Golang API name in embedded mode | IsAllowed                  | Diagnose                 |
| gRPC API name                    | IsAllowed                  | Diagnose                 |
| REST API name                    | /authz-check/v1/is-allowed | /authz-check/v1/diagnose |

## Example

Let's take REST API as example. Assuming you submitted the following authorization decision request,

```bash
curl -X POST http://localhost:6734/authz-check/v1/is-allowed -d @- << EOF
{
	"subject":{"principals":[{"type":"user", "name":"user1"}]},
	"serviceName": "srv1",
	"action": "get",
	"resource": "/api/v1/example/res1"
}
EOF
```

And you got the following response,

```bash
{"allowed":false,"reason":1}
```

Afterwards, you want to figure out why the request was denied using the policy diagnosis feature. You just need to replace "is-allowed" with "diagnose" in the REST request path, and keep using exactly the same request payload,

```bash
curl -X POST http://localhost:6734/authz-check/v1/diagnose -d @- << EOF
{
	"subject":{"principals":[{"type":"user", "name":"user1"}]},
	"serviceName": "srv1",
	"action": "get",
	"resource": "/api/v1/example/res1"
}
EOF
```

Then you will get a response something as below.

```bash
{
  "allowed": false,
  "reason": "DENY_POLICY_FOUND",
  "requestContext": {
    "subject": {
      "principals": [
        {
          "type": "user",
          "name": "user1"
        }
      ],
      "tokenType": "",
      "token": ""
    },
    "serviceName": "srv1",
    "resource": "/api/v1/example/res1",
    "action": "get",
    "attributes": null
  },
  "attributes": {
    "request_action": "get",
    "request_day": 28,
    "request_groups": [

    ],
    "request_hour": 17,
    "request_month": 1,
    "request_resource": "/api/v1/example/res1",
    "request_time": 1548666167,
    "request_user": "user1",
    "request_weekday": "Monday",
    "request_year": 2019
  },
  "policies": [
    {
      "status": "takeEffect",
      "id": "lre2z6nbklw7yxv2uxbb",
      "name": "policy2",
      "effect": "deny",
      "permissions": [
        {
          "resource": "/api/v1/example/res1",
          "actions": [
            "get"
          ]
        }
      ],
      "principals": [
        [
          "user:user1"
        ]
      ],
      "condition": {

      }
    },
    {
      "status": "ignored",
      "id": "6ww73cvfypkml46oibk2",
      "name": "policy1",
      "effect": "grant",
      "permissions": [
        {
          "resourceExpression": "/api/v1/example/.*",
          "actions": [
            "get"
          ]
        }
      ],
      "principals": [
        [
          "user:user1"
        ]
      ],
      "condition": {

      }
    }
  ]
}
```

From the above policy diagnosis response, we can easily tell that there were two policies which matched the request, one denies the user("user1") to operate("get") on the resource("/api/v1/example/res1"), while the other one allows the user("user1") to operate("get") on all resources which match the pattern("/api/v1/example/.\*"). If we translate the two policies into SPDL, then they are as below,

```
deny user user1 get /api/v1/example/res1
```

```
grant user user1 get expr:/api/v1/example/.*
```

It's obvious that user "user1" is allowed to get all resources that match the pattern "/api/v1/example/.\*", but except the resource "/api/v1/example/res1". So the previous authorization decision was denied.
