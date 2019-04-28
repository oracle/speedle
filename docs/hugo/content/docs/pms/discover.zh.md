+++
title = "策略发现"
description = "Discover policies on the fly!"
weight = 4
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["policy", "core", "discovery"]
categories = ["docs"]
bref = ""
+++

## 什么是 discovery mode?

当一个系统使用 Speedle 作为权限控制引擎时，在所有保护资源被访问之前， 都会向 Speedle 的 ARS(Authorization Runtime Service)发 authorization 请求。所有的 authorization 请求都被被发送到 Speedle ARS 的 is-allowed RESTful endpoint。ARS 根据系统中的所有 policy 计算出当前请求的资源访问是否允许。

当系统中需要保护的资源越来越多，为这些资源创建 policy 就是一件比较痛苦的事情。因为 policy 的制定者需要知道如何在 policy 中正确表述资源。discover mode 就是为了解决这一痛点而设计的。简单来说， 我们提供了一个和 is-allowed 有着相同输入和输出的接口 discover, 这个接口永远返回 allowed, 同时记录下 authorization 请求。并提供命令行工具查询被 discover 接口记录下的 authorization 请求,甚至为这些请求生成 Policy.

当我们把系统中的 is-allowed 接口统统换成 discover 接口,我们称系统工作在 discovery mode.

## discovery mode 能帮我们做什么?

- 记录 authorization 请求的内容

- 根据记录的 authorization 请求生成 Policy

- 关闭权限检查  
  因为 discover API 总是返回 is-allowed=true, 所以 discovery mode 相当于关闭了权限检查。

## 如何使用 discovery mode?

### Step 1. 将系统中所有 is-allowed 调用改成 discover 调用

如果你的系统是通过 RESTFul API 来调 is-allowed, 这种情况下，只需将 is-allowed endpoint 改成 discover endpoint, 如下所示：

```
http://localhost:6734/authz-check/v1/is-allowed ---> http://localhost:6734/authz-check/v1/discover
```

如果你的系统是通过 Grpc 或者 golang API 来调 is-allowed, 那么需要将所有的 is-allowed 调用改成 discover 调用。重新编译，并重启系统确保修改生效。

### Step 2. [optional] 使用命令行工具不间断发现 authorization 请求

使用 spctl discover request 命令来发现某一个服务下的所有 authorization 请求。 使用 --force 来不间断发现 authorization 请求。

```
spctl discover request --last --force --service-name=YOUR_SERVICE_NAME
```

保持窗口打开，这样你可以看到下一步中的 authotization 请求。

### Step 3. 访问系统中的被保护资源

不同的系统访问资源的方式不同，有通过 UI 访问的，有通过接口访问的。访问保护资源将触发 authotization 请求送往 Speedle,Speedle 会记录下收到的请求。

### Step 4. 基于访问生成对应的 Policy

使用 spctl discover policy 命令来为某个 service 生成 json 格式的 policy 定义。在这个例子中, 生成的 policy 存入 service.json 文件.

```
spctl discover policy --service-name=YOUR_SERVICE_NAME > service.json
```

### Step 5. [optional] 将生成的 policy 导入到 Speedle 中

使用 spctl create service 将上一步生成的 policy 导入 Speedle.

```
spctl create service --json-file service.json
```

### Step 6. [optional]将系统中所有 discover 调用改成 is-allowed 调用

最后别忘了将 discovery mode 切回正常模式。也就是 step 1 的逆操作。

## Discover 命令参考

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
