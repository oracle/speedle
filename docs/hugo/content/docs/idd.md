+++
title = "Identity Domain"
description = "Speedle supports user identities from multiple identity domains"
date = 2019-01-17T13:20:07+08:00
weight = 40
draft = false
bref = ""
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
tags = ["identity domain", "identity"]
+++

## What is an identity domain?

An identity domain is a logical namespace for users and groups, typically representing a discrete set of users and groups in the physical data store. Each identity domain manages users and groups independently. User and group names must be unique in an identity domain.

## How Speedle handles multiple identity domains

In an integrated environment, identities may come from multiple identity domains. Consider, for example, the following policy where a user with the identifier `user1` is authorized to perform a `rent` action on the resource `book`. This policy grants permission for `user1` from any identity domain to rent the book resource.

```bash

./spctl create policy -c "grant user user1 rent book" --service-name=booksvc

```

However, the service only expects `user1` from the identity domain `github` to have the `rent` permission. Therefore, Speedle must be able to distinguish between users coming from different identity domains.

With Speedle, to ensure that only the user from the expected identity domain is granted permission, you construct a new user/group identifier based on the incoming user/group identifier and the user/group identity domain. You can then use the new user/group identifier in the Speedle policies. The new identifier structure is defined as follows, where IDD represents the identity domain property.

```go

type Principal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

```

**Note:** Identity providers such as [Google account](https://account.google.com/) don't support multiple identity domains, while [Oracle Identity Cloud Service (IDCS)](https://www.oracle.com/cloud/paas/identity-cloud-service.html) does support multiple identity domains. The value of the identity domain property is implementation dependent.

## Defining policies with identity domains

To specify that the policy is defined for a user from the specified identity domain, use the [SPDL](../spdl) keyword: _**from &lt;identity domain&gt;**_

POLICY = EFFECT SUBJECT _**from &lt;identity domain&lt;**_ ACTION RESOURCE if CONDITION

### Examples

This policy demonstrates that only `user1` from the `github` identity provider can perform the `read` action on the `book` resource:

```bash

# grant user1 from github permission to perform the action: read on the resource: book
./spctl create policy -c "grant user user1 from github read book" --service-name=booksvc

```

This policy demonstrates that only `user1` from identity domain `tenant01` of identity provider `IDCS` can perform the `read` action on the `book` resource:

```bash

# grant user1 from tenant01 of IDCS permission to perform action: read on resource: book
./spctl create policy -c "grant user user1 from IDCS.tenant01 read book" --service-name=booksvc

```

### Policy Management REST APIs

Prefix each principal with _**"idd=&lt;identity domain&gt;:"**_

#### Example

```bash

# Create a policy: allow user1 coming from identity domain: github to perform action: read on resource: book
curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["read"]}],"principals": [["idd=github:user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

```

### How policies with identity domains are evaluated

- **Policy defined with an identity domain**

      	The evaluation engine first tests whether the user's identity domain in the incoming request strictly matches the identity domain defined in the policy. If there is an exact match, it then evaluates the policy. If there is not an exact match, the evaluation result for the defined policy is false.

- **Policy defined without an identity domain**

      	The evaluation engine evaluates the policy without considering the identity domain in the incoming request and will match principals from any identity domains.

## Samples

These samples show how to create a policy with identity domains using the `spctl` CLI and the Policy Management API, and the expected results. The sample policy demonstrates how a user with the same name, `user1`, is granted different permissions depending on which identity provider the user is from.

### Creating a policy with identity domains using the spctl CLI

```bash

./spctl create service booksvc
# grant user1 coming from github permission to perform action: read on resource: book
./spctl create policy -c "grant user user1 from github read book" --service-name=booksvc
# grant user1 coming from google permission to perform action: write on resource: book
./spctl create policy -c "grant user user1 from google write book" --service-name=booksvc
# grant user1 coming from any identity provider permission to perform action: rent on resource: book
./spctl create policy -c "grant user user1 rent book" --service-name=booksvc

```

JSON format of policy data

```json
{
	"services": [
		{
			"name": "booksvc",
			"policies": [
				{
					"id": "policy1",
					"effect": "grant",
					"permissions": [
						{
							"resource": "book",
							"actions": ["read"]
						}
					],
					"principals": [["idd=github:user:user1"]]
				},
				{
					"id": "policy2",
					"effect": "grant",
					"permissions": [
						{
							"resource": "book",
							"actions": ["write"]
						}
					],
					"principals": [["idd=google:user:user1"]]
				},
				{
					"id": "policy3",
					"effect": "grant",
					"permissions": [
						{
							"resource": "book",
							"actions": ["rent"]
						}
					],
					"principals": [["user:user1"]]
				}
			]
		}
	]
}
```

### Creating a policy with identity domains using the Policy Management API

```bash

curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["read"]}],"principals": [["idd=github:user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["write"]}],"principals": [["idd=google:user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["rent"]}],"principals": [["user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

```

### Policy evaluation results

The following policy evaluation results are based on the policies defined in the samples above.

```bash
# The evaluation result is true.
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"github"}] },"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# The evaluation result is false because of different identity domain: gitlab
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"gitlab"}] },"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# The evaluation result is true because of user1 coming from any identity providers can perform rent action
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1"}] },"serviceName":"booksvc","resource":"book","action":"rent"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"google"}] },"serviceName":"booksvc","resource":"book","action":"rent"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# The evaluation result is false because of different identity domain "idd":"notgoogle"
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"notgoogle"}] },"serviceName":"booksvc","resource":"book","action":"write"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

```
