+++
title = "Authorization Decisions"
description = "Get authorization decisions for your service interactions"
weight = 30
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["pdp", "policy", "core"]
categories = ["docs"]
bref = "Get authorization decisions"
+++

## What is an authorization decision?

- An authorization decision determines whether a subject performing an action on a resource is allowed.

- An authorization decision is the result of real-time evaluation based on policies and attributes.

## Ways to get authorization decisions

Authorization decisions can be performed by the Authorization Decision Service or an by an embedded evaluator:

- Authorization Decision Service (ADS)
  - REST API
  - Grpc API
- Embedded Evaluator
  - Golang API

## APIs and Samples

The ADS decision APIs make authorization decisions based on policies that describe the actions, permissions, and roles granted to a subject.

### Get decision

Get a decision on whether a subject performing an action on a resource is allowed.

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
  - Get a decision on whether user Alan is allowed to download a book from an online bookstore
  - Decision is based on policies defined in a service named "onlineBookStore"

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

Here, reason '0' means that the ADS found the grant policy. The list of reasons and definitions are as follows:

 <table class="bordered striped">
    <thead>
      <tr>
        <th>Reason</th>
        <th>Definition</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td> 0 </td>
        <td> GRANT_POLICY_FOUND </td>
      </tr>
      <tr>
        <td> 1 </td>
        <td> DENY_POLICY_FOUND </td>
      </tr>
      <tr>
        <td> 2 </td>
        <td> SERVICE_NOT_FOUND </td>
      </tr>
      <tr>
        <td> 3 </td>
        <td> NO_APPLICABLE_POLICIES </td>
      </tr>
      <tr>
        <td> 4 </td>
        <td> ERROR_IN_EVALUATION </td>
      </tr>
      <tr>
        <td> 5 </td>
        <td> DISCOVER_MODE </td>
      </tr>
   </tbody>
 </table>

### Get Roles

Get all the roles granted to the subject in a request.

- API overview

  - IN
    - Given the subject
    - Given the runtime attributes \*\*optional\*\*
    - Given the service scope
  - OUT
    - Returns a slice of roles granted to current subject
    - Returns errors if an error occurs

- Sample
  - Get the roles granted to the user Alan
  - Decision is based on policies defined in service named "onlineBookStore"

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

### Get Permissions

Get all permissions granted to the subject in a request.

- API overview

  - IN
    - Given the subject
    - Given the runtime attributes \*\*optional\*\*
    - Given the service scope
  - OUT
    - Returns a slice of (actions, resource) pairs, current subject is allowed to perform.
    - Returns errors if an error occurs

- Sample
  - Get all permissions granted to user Alan
  - Decision is based on policies defined in service named "onlineBookStore"

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
