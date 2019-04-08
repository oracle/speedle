# 示例：Speedle嵌入模式

# 构建
请确认已经正确的设置了两个Go语言环境变量`GOROOT`和`GOPATH`。并且你已经在该目录下。
```bash
$ go build -o expenses_sample cmd/expenses_sample
```

当构建完成后，会在当前目录下生成一个可执行文件`expenses_sample`。

# 为该示例定义的策略
策略被定义在文件expenses.spdl中。
```
[service.expenses]
[policy]
GRANT ROLE employee get, post, delete /reports
GRANT ROLE auditor get, modify /reports
[rolepolicy]
GRANT USER alice employee
```

在此文件中，定义了一个service `expenses`。在该service中，定义了两个策略，角色`employee`有权限`get`、
`post`、`delete`资源`/reports`; 角色`auditor`有权限`get`、`modify`资源`/reports`.

用户`alice`是一个`employee`。

# 运行并测试
```bash
# 该程序会在localhost:8080上监听。
$ ./expenses_sample expenses.spdl
```

```bash
# 验证用户alice是否能删除/reports
$ curl -X DELETE -u alice:afdsa http://localhost:8080/reports
deleteing an expense report is done

# 验证用户bob是否能删除/reports
$ curl -X DELETE -u bob:afdsa http://localhost:8080/reports
forbidden.
```

# 把用户bob设置成角色`employee`，并测试
```bash
# 把用户bob设置成角色employee
$ echo "GRANT USER bob employee" >> expenses.spdl
```

```bash
# 验证用户bob是否能删除/reports
$ curl -X DELETE -u bob:afdsa http://localhost:8080/reports
deleteing an expense report is done
```

