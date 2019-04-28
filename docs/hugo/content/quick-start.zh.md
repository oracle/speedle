+++
title = "Speedle 五分钟入门"
description = ""
weight = 10
draft = false
toc = false
tocheading = "h2"
bref = "Speedle 是一个非常强大的企业级权限管理方案。不同于传统企业级应用，Speedle 简单易学。使用者可以在 5 分钟之内编译，部署，和简单使用 Speedle。"
icon = "9. Getting Started.svg"
+++

## 第一步 编译

1. 先安装 GO 编译器，安装文件在https://golang.org/dl/
2. 设置 GOPATH 环境变量
3. 运行如下命令

```bash
$ go get github.com/oracle/speedle/cmd/…
```

该命令会从 github 下载代码并自动编译。命令执行完毕后，我们可以检查\$GOPATH/bin 目录，应该有三个文件在该目录下：

```bash
$ ls $GOPATH/bin
spctl speedle-ads speedle-pms
```

speedle-pms 是授权策略（Policy）管理服务  
speedle-ads 是授权决定服务（运行时）  
spctl 是命令行工具

## 第二步 运行

启动 PMS 和 ADS

```bash
$ cd $GOPATH/bin
$ ./speedle-pms --store-type file &
$ ./speedle-ads --store-type file &
```

## 第三步 使用

作为权限控制软件，最基本的功能有两个：

1. 管理授权策略（Policy)
2. 处理授权请求，根据定义的 Policy 得出授权决定

例如，对于图书馆管理系统，我们要定义一个 Policy 说“张三可以借书”。我们可以进行如下操作：

```bash
$ ./spctl create service library
$ ./spctl create policy -c "grant user ZhangSan borrow book" --service-name=library
```

这样 Policy 就存储在 Speedle 中，接下来我们测试一下 Speedle 可不可以正确处理授权请求。

## 测试

问问它张三可不可以借书：

```bash
$ curl -X POST --data '{"subject":{"principals":[{"type":"user","name":"ZhangSan"}]},"serviceName":"library","resource":"book","action":"borrow"}' http://127.0.0.1:6734/authz-check/v1/is-allowed
```

它回答可以

再问问它李四可不可以借书

```bash
$ curl -X POST --data '{"subject":{"principals":[{"type":"user","name":"LiSi"}]},"serviceName":"library","resource":"book","action":"borrow"}' http://127.0.0.1:6734/authz-check/v1/is-allowed
```

它回答不可以

就这么简单！下一步您可以深入了解 SPDL 语言的用法。它可以支持普通的 ACL,也支持 RBAC 和 ABAC，也可以 RBAC 和 ABAC 混合使用。能满足各种应用场景。期待您去尝试。
