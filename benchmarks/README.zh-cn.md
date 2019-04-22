# 基准测试

## 运行基准测试
```bash
go test -bench .
```

## 测试报告

[eval_b_test.go](eval_b_test.go)用来得到策略决策的时间开销。测试平台为：

2 * `Intel(R) Xeon(R) Platinum 8167M CPU @ 2.00GHz 2核, 4线程`

该表格展示了API `IsAllowed`的时间开销。`op`的意思是对`IsAllowed`的一次调用。

用例|模型|时间开销 (μs/op)
----|----|----
One|1 policy, 1 role policy, 1 user|6.6
Tiny|10 policies, 10 role policies, 10 users/role|7.3
Small|100 policies, 100 role policies, 10 users/role|8.5
Medium|1K policies, 1K role policies, 10 users/role|8.6
Large|10K policies, 10K role policies, 10 users/role|8.7
Huge|100K policies, 100K role policies, 10 users/role|8.7
LargeExp|10K policies with resource expression, 10K role policies, 10 users/role|43.0
LargeCond|10K policies with condition, 10K role policies, 10 users/role|8.0

以下表格展示了API `GetAllGrantedPermissions`的时间开销。

用例|模型|时间开销 (μs/op)
----|----|----
LargePerm|10K policies, 10K role policies, 10 users/role|8.9

