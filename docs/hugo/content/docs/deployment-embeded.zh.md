+++
title = "Speedle嵌入模式"
description = ""
weight = 12
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
bref = ""
+++

## 什么是 Speedle 嵌入模式

区别于 Speedle 服务模式，Speedle 嵌入模式把 Speedle 的策略决策引擎嵌入到调用者进程内部，作为进程的一部分运行。调用者可以直接调用 Speedle 的 API 获得策略决策的结果。

## 嵌入模式的局限

- 只支持 Golang，推荐使用 Golang 1.10。
- 只接受文件存储策略。

## 如何使用嵌入模式

### 把 Speedle 代码添加到你的工程中

在嵌入模式中，需要先将 Speedle 源代码从代码仓库中下载下来。主要有一下几种方式：

1. 使用 go get

```bash
go get github.com/oracle/speedle
```

2. 使用 dep 工具

在你的 Gopkg.toml 文件中添加以下几行

```toml
[[constraint]]
  name = "github.com/oracle/speedle"
  branch = "master"
```

使用 dep 工具将 Speedle 源代码添加到`vendor`目录下：

```bash
dep ensure -update github.com/oracle/speedle
```

### 在代码中初始化一个 Evalautor 实例

方法`func eval.NewFromFile(loc string, isWatch bool) (ads.PolicyEvaluator, error)`使用一个策略定义文件初始化一个 Evaluator 实例。该方法接受两个参数：

- loc：决策定义文件的路径。支持两种决策文件格式：JSON 格式和 SPDL 格式。
- isWatch: 如果值为`true`，当策略定义文件发生改动，Speedle 决策引擎能自动以最新的文件内容进行决策。如果值为`false`，Speedle 决策引擎不监视策略定义文件的任何改动。

该方法返回一个 Evaluator 示例。

以下代码片段展示如何初始化一个 Evaluator 实例。

```go
import (
  "github.com/oracle/speedle/pkg/eval"
)

func foo() {
    eval, err := evaluator.NewFromFile(spdlLoc, true)
}
```

### 调用 Evaluator API 进行决策

方法`func ads.IsAllowed(c ads.RequestContext) (allowed bool, reason Reason, err error)`依据策略定义文件中定义的策略进行决策。

该方法接受一个参数`ads.RequestContext`，该结构的定义如下：

```go
typedef RequestContext {
  Subject     *ads.Subject                 // Subject包含用户信息，用户名，组名等。
  ServiceName string                       // Speedle Service
  Resource    string                       // Resource Name
  Action      string                       // Action Name
  Attributes  map[string]interface{}       // Attributes
}
```

该方法有 3 个返回值：

1. allowed bool: `true`，`Subject`指定的用户有权限操作资源(Action, Resource)。`false`，没有权限操作资源(Action, Resource)。
2. reason Reason: 做出决策的原因。
3. err error: 决策过程中出现的异常。如果没有异常，返回 nil。

## 嵌入模式的例子

例子可以在
https://github.com/oracle/speedle/tree/master/samples/embedded/expenses
中找到。
