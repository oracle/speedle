+++
title = "Policy Discovery"
description = "Discover policies on the fly!"
weight = 4
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["policy", "core", "discovery"]
categories = ["docs"]
bref = "Speedle can be used for authorization decisions in systems with a large number of protected resources. Before those protected resources can be accessed, the Speedle authorization decision API (is-allowed API) is called. Defining the authorization policies for those protected resources is difficult if you don't know what actions are occuring on the protected resources, or how a protected resource is expressed in an authorization decision request. Policy discovery mode can make this task much easier"
+++

## What can discover mode do?

### \* Help you discover existing authorization requests

The _is-allowed_ API endpoint and the _discover_ API endpoint use the same request and response format. When you replace _is-allowed_ with _discover_ in the URL of the API, each authorization request that is returned shows _WHO_ carries out _WHAT_ action on _WHICH_ resource. Discover mode does not evaluate the authorization decision; it simply records the authorization decision requests, along with the actions and the resources protected, and returns "is-allowed=true".

Note that when using discover mode, the discover endpoint must be called by every service. Therefore, Speedle must be integrated into your system so that every authorization request flows through it.

### \* Generate policies using the discovered requests

You can generate policies based on the authorization requests that you record. In one scenario, a developer runs a test suite that mimics all the interactions in your system, and Speedle records all the necessary authorization calls. For example, Service A (a UI) calls Service B (a database). Speedle can effectively record - "Service A called Service B with these attributes..." and generates a policy that will allow such a request (is-allowed=true). You can then import the policy directly into Speedle.

All of the policies generated are _is-allowed_ policies because they are recorded directly from allowed authorization requests on protected resources.

### \* Disable authorization without changing the code

When using Speedle for authorization checks, calls made to the _is-allowed_ endpoint return true or false, depending on whether the call is allowed or not.
However, when you use discover mode and replace _is-allowed_ with _discover_ in the url, all authorization checks return an "is-allowed=true" result. Because discover mode does not perform authorization checks, it provides a convenient way for developers to disable authorization checks when diagnosing failures.

## How do I use discover mode?

### Step 1. Change all **is-allowed** calls in your system to **discover** calls

If your system calls the authorization decision API through REST, the REST endpoint of the authorization decision API is usually configured somewhere in your system. In this case, just change the REST endpoint of the authorization decision API from an _is-allowed_ endpoint to a _discover_ endpoint. Be sure to apply the configuration changes.

```
http://localhost:6734/authz-check/v1/is-allowed ---> http://localhost:6734/authz-check/v1/discover
```

If your system calls the authorization decision API through Grpc or a golang API, then change all the authorization calls from _is-allowed_ to _discover_, rebuild the changed code, and restart your system to make the change take effect.

### Step 2. [optional] Discover the authorization requests

Use the `spctl discover request` command to list all request details for the specified service. Using the `--force` option continuously discovers the last request.

```
spctl discover request --last --force --service-name=YOUR_SERVICE_NAME
```

Keep the window open so that you can review the requests in the next step.

### Step 3. Access the resources that are protected in your system

How you access the resources depends on your system, such as through a UI or through REST endpoints. Accessing protected resources triggers authorization decision requests to Speedle, which Speedle can then record. Depending on your system, this may trigger a lot of authorization requests.

### Step 4. Generate policies for your service based on the authorization requests

Use the `spctl discover policy` command to generate a JSON-based policy definition for the specified service. In this example, the policy definition is named `service.json`.

```
spctl discover policy --service-name=YOUR_SERVICE_NAME > service.json
```

### Step 5. [optional] Import the policies into Speedle

Use the `spctl create` command to create a service with policies using the json service definition created in step 4.

```
spctl create service --json-file service.json
```

### Step 6. [optional] Change all the **discover** calls in your system back to **is-allowed** calls

After you create your policies and import them into your system, you should update the endpoints to remove `discover` and use `is-allowed` instead. Be sure to apply the configuration changes.

## Discover mode command line reference

```
$ ./spctl discover --help
discover request or policy for services

Usage:
  spctl discover (request/policy/reset  | --service-name=NAME | --last | --force | --principal-name=USERNAME) [flags]

Examples:

        # List all request details for all services
        spctl discover request

        # List all request details for the given service
        spctl discover request --service-name="foo"

        # List the last request details for service "foo"
        spctl discover request --last --service-name="foo"

        # List the latest request details for service "foo", doesn't exit until you kill it using "Ctrl-C"
        spctl discover request --last --service-name="foo" -f

        # cleanup all requests
        spctl discover reset

        # clean up the requests for service "foo"
        spctl discover reset --service-name="foo"

        # Generate JSON based policy definition, all users are converted to a role. For example, user Jon visited resourceA. Then the following policy is generated "grant role role_Jon visit resourceA"
        spctl discover policy  --service-name="foo"

        # Generate JSON based policy definition, only for discover requests triggered by principal which has name 'Jon'
        spctl discover policy --principal-name="Jon" --service-name="foo"

Flags:
  -f, --force                   continuously discover last request
  -h, --help                    help for discover
  -l, --last                    list last request
      --principal-IDD string    principal Identity Domain
      --principal-name string   principal name
      --principal-type string   principal type, could be 'user', 'group','entity'
  -s, --service-name string     service name


```
