+++
title = "策略管理"
description = "Manage policy lifecycle "
weight = 1
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["pms", "policy", "core"]
categories = ["docs"]
bref = "Basics of policy management"
+++

## What is a Speedle policy?

A Speedle policy is a set of criteria that specify whether a user is granted access to a particular protected resource or assignment to a particular role. You manage Speedle policies using the Speedle Policy Management Service(PMS).

## Understanding the Speedle Policy Module

**Note:** The Speedle syntax used in this document is defined in [SPDL - Security Policy Definition Language](../../spdl).

#### Policy store

The policy store maintains all policy artifacts and can be persisted to an etcd store or a JSON file.

<img src="/img/speedle/policystore.png"/>

#### Service

A service is a container that contains a set of authorization and role policies that exist only in the scope of that service. Policies and role policies are evaluated within the scope of the service in which they were defined, not in the entire policy store. You can manage multiple services with Speedle.

You can also define global policies in a global service. Global policies take effect globally across all services. For details, see [Global Policy](../global-policy).

#### Authorization policy

An authorization policy defines the criteria that controls access to protected resources.

<img src="/img/speedle/authzpolicy.png"/>

You create authorization policies to grant or deny principals (user/role/group/entity) permission to perform specific actions on specific resources if the condition is true.

Sample:

```
grant group Administrators list,watch,get expr:c1/default/core/pods/*
```

This sample grants the group "Administrators" permission to perform "list", "watch", and "get" operations on the resource that matches the name expression `c1/default/core/pods/*`.

#### Role policy

A role policy defines the criteria that controls how principals (user/role/group/entity) are granted or denied membership to roles created using Speedle.

<img src="/img/speedle/rolepolicy.png"/>

You create role policies to grant or deny roles, which you created using Speedle, to principals (user/role/group/entity) on specific resources if the condition is true.

Sample:

```
grant user alan manager on res1
```

This sample grants user "alan" the "manager" role on the resource "res1". In other words, user "alan" can perform operations on the resource "res1" because "alan" has the permissions assigned to the role "manager".

#### Policy elements

##### Effect

Effect has two values: "grant" or "deny".  
When Speedle evaluates policies, the final authorization decision is based on the "DENY overrides" combining algorithm. For example, if there is a policy that grants permission to a subject at the same time as a policy that denies the same permission to the subject, then the "deny" policy takes effect and overrides the "grant" policy.

##### Principal

In authorization and role policies, the principal is the identity object to which the access rights or roles can be granted or denied. A principal can be a user, a group, an entity or a role. Most frequently, it is a role.

<img src="/img/speedle/principal.png"/>

User, group and entity are principals from the identity store and are usually obtained after authentication or token assertion. Users and groups represent a human identity; an entity represents a non-human identity such as a service, a Kubernetes pod, and so on.

#### AND principal

AND principal is a combination of a small set of principals, separated by commas. If a policy uses AND principal, the policy can take effect only when all of these principles are matched.

<img src="/img/speedle/andprincipal.png"/>

Sample:

```
grant role (designer, dba) update db_design_doc
```

In this sample, only a user with both roles "designer" and "dba" can update the resource "db_design_doc".

##### Resource

A resource is a protected object to which access is granted or denied. A resource represents the application component or business object that is secured by an authorization policy.

<img src="/img/speedle/resource.png"/>

resourceNameExpression supports regular expressions.

##### Action

An action is an operation that can be performed on the protected resource. Action is just a string in a policy. You can define any actions when you create the policy.

##### Condition

A condition is a bool expression that is constructed using attributes, functions, constants, operators, comparators or parenthesis and produces a bool value. Conditions are supported in both role and authorization policies. The policy or role policy can take effect only when the condition is met.

For details, see [SPDL - Security Policy Definition Language](../../spdl).

## Managing Speedle policies

Use the Speedle Policy Management Service (PMS) to manage authorization and role policies, and the security objects from which they are created.

Speedle allows administrators to perform create, read, and delete operations on all policy objects. You can do this in any of the following ways:

-   Using the Speedle command line interface `spctl` (as described here. This is the recommended method.)

-   Using the PMS Golang Management API in Embedded Mode (as described in the [Speedle API doc](https://github.com/oracle/speedle/tree/master/api/pms).

-   Using the PMS REST Service (as described in the [Speedle Policy Management API](../docs/api/management_api)).

-   Using the PMS gRPC Service (as described in the [Speedle GRPC document](/protobuf/pms.proto)).

#### Managing services

You create a service as the overall container for authorization and role policies.
You can perform the following management operations on service instances.

-   Create a "test" service:

```bash
$ ./spctl create service test
service created
{"name":"test","type":"application","metadata":{"createby":"","createtime":"2019-02-12T22:51:19-08:00"}}
```

-   Get the "test" service:

```bash
$ ./spctl get service test
{
    "name": "test",
    "type": "application",
    "metadata": {
        "createby": "",
        "createtime": "2019-02-12T22:51:19-08:00"
    }
}
```

-   Get all services:

```bash
$ ./spctl get service --all
[
    {
        "name": "test",
        "type": "application",
        "metadata": {
            "createby": "",
            "createtime": "2019-02-12T22:51:19-08:00"
        }
    }
]
```

-   Delete the "test" service:

```bash
$ ./spctl delete service test
service test deleted.
```

#### Managing authorization policies

You can perform the following management operations on authorization policies.

-   Create a policy named "policy1" in the "test" service:

```bash
$ ./spctl create policy policy1 -c "grant user alan read book" --service-name test
policy created
{"id":"ao3olis24hrzchwjduea","name":"policy1","effect":"grant","permissions":[{"resource":"book","actions":["read"]}],"principals":[["user:alan"]],"metadata":{"createby":"","createtime":"2019-02-12T22:57:46-08:00"}}
```

-   Get "policy1" in the "test" service using the policy id:

```bash
$ ./spctl get policy ao3olis24hrzchwjduea --service-name=test
{
    "effect": "grant",
    "id": "ao3olis24hrzchwjduea",
    "metadata": {
        "createby": "",
        "createtime": "2019-02-12T22:57:46-08:00"
    },
    "name": "policy1",
    "permissions": [
        {
            "actions": [
                "read"
            ],
            "resource": "book"
        }
    ],
    "principals": [
        [
            "user:alan"
        ]
    ]
}
```

-   Delete "policy1" in the "test" service using the policy id:

```bash
$ ./spctl delete policy ao3olis24hrzchwjduea --service-name=test
policy ao3olis24hrzchwjduea deleted.
```

#### Managing role policies

You can perform the following management operations on role policies.

-   Create a new role policy named "rolepolicy01" in the "test" service:

```bash
$ ./spctl create rolepolicy rolepolicy01 -c "grant user alan manager" --service-name test
rolepolicy created
{"id":"4gskmqamoiebmidyw2fi","name":"rolepolicy01","effect":"grant","roles":["manager"],"principals":["user:alan"],"metadata":{"createby":"","createtime":"2019-02-12T23:00:44-08:00"}}
```

-   Get the role policy using the policy id:

```bash
$ ./spctl get rolepolicy 4gskmqamoiebmidyw2fi --service-name test
{
    "effect": "grant",
    "id": "4gskmqamoiebmidyw2fi",
    "metadata": {
        "createby": "",
        "createtime": "2019-02-12T23:00:44-08:00"
    },
    "name": "rolepolicy01",
    "principals": [
        "user:alan"
    ],
    "roles": [
        "manager"
    ]
}

```

-   Delete the role policy using the policy id:

```bash
$ ./spctl delete rolepolicy 4gskmqamoiebmidyw2fi --service-name test
rolepolicy 4gskmqamoiebmidyw2fi deleted.
```
