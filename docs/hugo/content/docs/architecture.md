+++
title = "Architecture"
description = "Speedle building blocks"
weight = 10
draft = false
toc = false
tocheading = "h2"
tocsidebar = false
bref = "Speedle architecture"
+++

## Modules

Speedle is an authorization engine comprised of these components:

- Policy Management Service (PMS) API - Manages the authorization and role policies, and the objects from which they are created.

- Policy repository - Stores all the policy artifacts. The policy repository can be a json file, or a persistent store such as a database or etcd.

- Authorization Decision Service (ADS) API - Evaluates the authorization requests against the applicable policies and returns GRANT/DENY decisions.

This diagram shows the interaction between these components.

<img src="/img/speedle/spdlarch.jpg" />

(1) Users create/manage policies through the Policy Management Service API.  
(2) The Policy Management Service persists the policies in the policy repository.  
(3) The policies are provisioned from the policy repository to the Authorization Decision Service API.  
(4) Users systems invoke the Authorization Decision Service API for authorization checks.

If you are familiar with the XACML model, the Policy Management Service (PMS) serves as the Policy Administration Point (PAP), and the Authorization Decision Service (ADS) serves as the Policy Decision Point (PDP).
