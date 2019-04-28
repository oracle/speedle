+++
title = "用户身份域"
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

## 什么是用户身份域?

身份域是用户和组的逻辑名称空间，通常表示物理数据存储中的一组离散用户和组。 每个身份域独立管理用户和组， 用户名和组名在标识域中必须是唯一的。

## 问题

在集成环境中，身份可能来自多个身份域。 例如，考虑以下场景，具有身份为`user1`的用户被授权对资源`book`执行`rent`操作。 此授权策略表示来自任何身份域的身份为`user1`的用户都具有权限以租用图书资源。

```bash

./spctl create policy -c "grant user user1 rent book" --service-name=booksvc

```

但是，该服务只希望来自身份域`github`的`user1`具有`rent`权限。 因此，Speedle 必须能够区分来自不同身份域的用户。

## 解决方案

要确保仅授予来自预期身份域的用户权限，您需要根据用户/组标识符和用户/组身份域构造新的用户/组标识符。 然后，您可以在 Speedle 策略中使用新的用户/组标识符。 新标识符结构定义如下，其中 IDD 表示身份域属性。

```go

type Principal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}

```

**备注:** 身份提供商（如[google account]（https://account.google.com/））不支持多个身份域，而[IDCS]（https://www.oracle.com/cloud/paas/identity-cloud -service.html）确实支持多个身份域。 身份域属性的值取决于实现。

## 使用身份域定义授权策略

要指定从指定的标识域为用户定义策略，请使用[SPDL]（../ spdl）关键字: _**from &lt;identity domain&gt;**_

POLICY = EFFECT SUBJECT _**from &lt;identity domain&lt;**_ ACTION RESOURCE if CONDITION

### 示例

这个策略表明只有来自`github`身份提供者的`user1`才能对`book`资源执行`read`动作：

```bash

# 授权来自 github 的名为 user1 的用户可以对资源： book具有 read的权限
./spctl create policy -c "grant user user1 from github read book" --service-name=booksvc

```

这个策略表明只有来自身份提供商`IDCS`的身份域`tenant01`的`user1`才能对`book`资源执行`read`动作：

```bash

# 授权来自身份提供商IDCS的身份域tenant01的名为 user1 的用户可以对资源： book具有 read的权限
./spctl create policy -c "grant user user1 from IDCS.tenant01 read book" --service-name=booksvc

```

### 授权策略管理的 REST APIs

需要为每一个用户身份加上前缀: _**"idd=&lt;identity domain&gt;:"**_

#### 示例

```bash

# 创建一个授权策略: 允许 user1 来自身份域: github 可以对资源: book 执行： read 操作
curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["read"]}],"principals": [["idd=github:user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

```

### 具有身份域的策略在策略决策中的作用

- **具有身份域的授权策略**

      	决策引擎首先测试传入请求中的用户身份域是否严格匹配策略中定义的身份域。 如果存在完全匹配，则会继续评估该策略。 如果没有完全匹配，则定义的策略的评估结果直接为否。

- **没有身份域的授权策略**

      	决策引擎在不考虑传入请求中的身份域的情况下评估策略，意味着所定义的授权策略将匹配来自任何身份域的用户。

## 示例

这些示例显示如何使用`spctl` CLI 和 Policy Management API 以及预期结果创建具有标识域的策略。 示例策略演示了具有相同名称的用户`user1`如何被授予不同的权限，具体取决于用户来自哪个身份提供者。

### 使用 spctl CLI 创建具有身份域的策略

```bash

./spctl create service booksvc
# 允许 user1 来自身份域: github 可以对资源: book 执行： read 操作
./spctl create policy -c "grant user user1 from github read book" --service-name=booksvc
# 允许 user1 来自身份域: google 可以对资源: book 执行： write 操作
./spctl create policy -c "grant user user1 from google write book" --service-name=booksvc
# 允许 user1 来自任何身份域可以对资源: book 执行： rent 操作
./spctl create policy -c "grant user user1 rent book" --service-name=booksvc

```

JSON 格式的策略数据

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

### 使用策略管理 API 创建具有身份域的策略

```bash

curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["read"]}],"principals": [["idd=github:user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["write"]}],"principals": [["idd=google:user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

curl -v -X POST -d '{"name": "policy1","effect": "grant","permissions": [{"resource": "book","actions": ["rent"]}],"principals": [["user:user1"]]}' http://127.0.0.1:6733/policy-mgmt/v1/service/booksvc/policy

```

### 策略请求的评估结果

以下策略评估结果基于上述示例中定义的策略。

```bash
# 评估结果为： 是.
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"github"}] },"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# 评估结果为： 否, 因为评估请求中的用户来自不同的身份域: gitlab
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"gitlab"}] },"serviceName":"booksvc","resource":"book","action":"read"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# 评估结果为： 是, 因为来自任何身份域的用户 user1 都可以对资源： book执行 rent 操作
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1"}] },"serviceName":"booksvc","resource":"book","action":"rent"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# 评估结果为： 是, 因为来自任何身份域的用户 user1 都可以对资源： book执行 rent 操作
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"google"}] },"serviceName":"booksvc","resource":"book","action":"rent"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

# 评估结果为： 否, 因为评估请求中的用户来自不同的身份域: notgoogle
curl -v -X POST -d '{ "subject": {"principals":[{"type":"user","name":"user1","idd":"notgoogle"}] },"serviceName":"booksvc","resource":"book","action":"write"}'  http://127.0.0.1:6734/authz-check/v1/is-allowed

```
