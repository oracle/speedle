# 什么是discover mode?
当一个系统使用Speedle作为权限控制引擎时，在所有保护资源被访问之前， 都会向Speedle的ARS(Authorization Runtime Service)发authorization请求。所有的authorization请求都被被发送到Speedle ARS的is-allowed RESTful endpoint。ARS根据系统中的所有policy计算出当前请求的资源访问是否允许。

当系统中需要保护的资源越来越多，为这些资源创建policy就是一件比较痛苦的事情。因为policy的制定者需要知道如何在policy中正确表述资源。discover mode就是为了解决这一痛点而设计的。简单来说， 我们提供了一个和is-allowed有着相同输入和输出的接口discover, 这个接口永远返回allowed, 同时记录下authorization请求。并提供命令行工具查询被discover接口记录下的authorization请求,甚至为这些请求生成Policy.

当我们把系统中的is-allowed接口统统换成discover接口,我们称系统工作在discover mode.
# discover mode 能帮我们做什么?

* 记录authorization请求的内容

* 根据记录的authorization请求生成Policy

* 关闭权限检查    
因为discover API总是返回is-allowed=true, 所以discover mode相当于关闭了权限检查。

# 如何使用discover mode?

Step 1. 将系统中所有is-allowed调用改成discover调用
如果你的系统是通过RESTFul API来调is-allowed, 这种情况下，只需将is-allowed endpoint改成discover endpoint, 如下所示：
```
http://localhost:6734/authz-check/v1/is-allowed ---> http://localhost:6734/authz-check/v1/discover
```
如果你的系统是通过Grpc或者golang API来调is-allowed, 那么需要将所有的is-allowed调用改成discover调用。重新编译，并重启系统确保修改生效。

Step 2. [optional] 使用命令行工具不间断发现authorization请求
使用 spctl discover request 命令来发现某一个服务下的所有authorization请求。 使用 --force 来不间断发现authorization请求。
```
spctl discover request --last --force --service-name=YOUR_SERVICE_NAME
```
保持窗口打开，这样你可以看到下一步中的authotization请求。

Step 3. 访问系统中的被保护资源
不同的系统访问资源的方式不同，有通过UI访问的，又通过接口访问的。访问保护资源将触发authotization请求送往Speedle,Speedle会记录下收到的请求。

Step 4. 基于访问生成对应的Policy
使用 spctl discover policy 命令来为某个service生成json格式的policy定义。在这个例子中, 生成的policy存入service.json文件.
```
spctl discover policy --service-name=YOUR_SERVICE_NAME > service.json
```

Step 5. [optional] 将生成的policy导入到Speedle中
使用 spctl create service 将上一步生成的policy导入Speedle.
```
spctl create service --json-file service.json
```

Step 6. [optional]将系统中所有discover调用改成is-allowed调用
最后别忘了将discover mode切回正常模式。也就是step 1 的逆操作。

# Discover 命令参考
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
