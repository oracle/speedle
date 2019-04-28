+++
title = "Global Policy"
description = "Policies take effect globally"
date = 2019-01-18T21:19:44+08:00
weight = 2
draft = false
bref = "When there are policies that are common to all services, it's tedious and difficult to manage if the same policies are created in each service. In this scenario, you can use global policies to simplify your policy management tasks"
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
+++

## How global policies simplify policy management

In Speedle, role policies and policies are defined and evaluated within a scope, which is called a 'service'. For example, when creating a policy, you need to specify the service in which to create the policy.
When calling the Authorization Decision Service API, you need to specify a service name so that the policies and role policies are evaluated within the scope of that service.

However, in a large system with many sub-systems, for example, each sub-system has its own specific authorization requirement, and each sub-system defines role policies and policies within its own service scope. Because all of the sub-systems are under the same large system, it is likely that there are role policies or policies that are common across all sub-systems. In this case, the same policies have to be defined in each service, and it becomes increasingly hard to manage as the number of policies grows.

To avoid creating the same role policies or policies repeatedly in each service, you can create a 'global policy' in a 'global service', which is a policy or role policy that can take effect globally across all services.

**Note:** Global authorization policies are not currently supported.

## How to use global role policies

You define global role policies in a special 'global' service. You can create and delete global services, and call the ADS API in the scope of a global service in the same way you create or call any other service. Global role policies are evaluated in global services just as if they were created in an ordinary service.

To use global policies in your system, follow these steps.

### 1. Create a global service

Create a service named 'global'.

```
./spctl create service global

```

### 2. Create a global role policy

Create a role policy in the global service granting user Emma the Admin role.

For example:

```
./spctl create rolepolicy -c "grant user Emma AdminRole" --service-name=global
```

### 3. Create an ordinary service

Create an ordinary service named 'library'.

For example:

```
./spctl create service library
```

### 4. Create a policy in the ordinary service

Create a policy that allows users in the AdminRole permission to borrow books.

For example:

```
./spctl create policy -c "grant role AdminRole borrow books" --service-name=library
```

### 5. Call the ADS is-allowed API within the scope of the ordinary service

Call the `is-allowed` API in the ordinary service scope, where the global role policy will also take effect.

For example:

```
curl -X POST  http://localhost:6734/authz-check/v1/is-allowed \
-d @- << EOF
{
 "subject": {"principals":[{"type":"user", "name":"Emma"}]},
 "serviceName": "library",
 "resource":"books",
 "action": "borrow"
}
EOF
```

In this example, the API returns `allowed = true`, because the role policy defined in the global service (user Emma is assigned the Admin role) takes effect.
