+++
title = "Speedle服务的部署方式"
description = "Recommended deployment topology for Speedle"
weight = 14
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
+++

Speedle 支持两种部署方式

1. 开发模式

在开发模式下，Speedle 使用一个文件作为策略仓库，用户使用 RESTful API 和 CLI 工具访问 Speedle 的各项服务。

2. 生产模式

在生产模式下，Speedle 使用`etcd`作为策略仓库，用户使用 RESTful API 和 CLI 工具访问 Speedle 的各项服务。

## 准备工作

1. Golang 1.10.0 开发环境已经安装好。
2. Docker 1.12 或更高的版本已经安装好。
3. 一个 Docker registry。该 Docker registry 被用作推送和拉取 Speedle 的 Docker 映像。
4. 在目录`\$GOPATH/src/github.com/oracle/speedle/`下同步最新的 Speedle 代码。
5. 一个 Kubernetes 集群。
6. 为你的 Docker registry 设置一个 Kubernetes Secret

```bash
# 为Docker registry设置Kubernetes Secret
$ kubectl create secret docker-registry reg-speedle --docker-server=$DOCKER_LOGIN_SERVER --docker-username=$DOCKER_LOGIN_USER --docker-password=$DOCKER_LOGIN_PASSWORD
```

## 开发环境下部署 Speedle 服务

开发环境下可以使用`helm`部署 Speedle 服务。

```bash
$ helm install -n speedle deployment/helm/speedle-dev
NAME:   speedle
LAST DEPLOYED: Tue Nov 27 23:19:37 2018
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/Service
NAME         TYPE       CLUSTER-IP      EXTERNAL-IP  PORT(S)   AGE
speedle-pms  ClusterIP  10.101.117.233  <none>       6733/TCP  1s
speedle-ads  ClusterIP  10.104.5.111    <none>       6734/TCP  1s

==> v1beta2/Deployment
NAME         DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
speedle-pms  1        1        1           0          1s
speedle-ads  1        1        1           0          1s

==> v1/Pod(related)
NAME                          READY  STATUS             RESTARTS  AGE
speedle-pms-867787bc95-kftbr  0/1    ContainerCreating  0         1s
speedle-ads-7877ffbbf7-bwwtd  0/1    ContainerCreating  0         1s


NOTES:
1. $ kubectl port-forward svc/speedle-pms 6733:6733
  Then access http://127.0.0.1:6733/policy-mgmt/v1/ to manage policies.

2. $ kubectl port-forward svc/speedle-ads 6734:6734
  Then access http://127.0.0.1:6734/authz-check/v1/is-allowed to check permissions.
```

## 生产环境下部署 Speedle 服务

### 安装 etcd 集群

强烈建议使用`etcd-operator`部署`etcd`。https://github.com/coreos/etcd-operator

使用`helm`部署`etcd-operator`：

```bash
$ helm install stable/etcd-operator --name my-release
```

安装一个新的`etcd`集群：

```bash
$ kubectl create -f https://raw.githubusercontent.com/coreos/etcd-operator/master/example/example-etcd-cluster.yaml
```

一个 3 节点的`etcd`集群将被安装：

```bash
$ kubectl get pods
NAME                            READY     STATUS    RESTARTS   AGE
example-etcd-cluster-gxkmr9ql7z   1/1       Running   0          1m
example-etcd-cluster-m6g62x6mwc   1/1       Running   0          1m
example-etcd-cluster-rqk62l46kw   1/1       Running   0          1m
```

如果你希望安装一个支持 TLS 协议的`etcd`集群，请参考一下文档：
https://github.com/coreos/etcd-operator/blob/master/doc/user/cluster_tls.md

### 部署 Speedle

编辑文件`values.yaml`，修改跟`etcd`的监听地址。

```
store:
  type: etcd
  etcd:
    endpoint: http://<etcdServiceName>:<etcdServicePort>
```

你也可以在执行命令`helm install`的时候，使用参数`--set store.etcd.endpoint=http://<etcdServiceName>:<etcdServicePort>`指定`etcd`的监听地址。

```bash
$ helm install -n speedle deployment/helm/speedle-prod
NAME:   speedle
LAST DEPLOYED: Tue Nov 27 23:19:37 2018
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/Service
NAME         TYPE       CLUSTER-IP      EXTERNAL-IP  PORT(S)   AGE
speedle-pms  ClusterIP  10.101.117.233  <none>       6733/TCP  1s
speedle-ads  ClusterIP  10.104.5.111    <none>       6734/TCP  1s

==> v1beta2/Deployment
NAME         DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
speedle-pms  1        1        1           0          1s
speedle-ads  1        1        1           0          1s

==> v1/Pod(related)
NAME                          READY  STATUS             RESTARTS  AGE
speedle-pms-867787bc95-kftbr  0/1    ContainerCreating  0         1s
speedle-ads-7877ffbbf7-bwwtd  0/1    ContainerCreating  0         1s


NOTES:
1. $ kubectl port-forward svc/speedle-pms 6733:6733
  Then access http://127.0.0.1:6733/policy-mgmt/v1/ to manage policies.

2. $ kubectl port-forward svc/speedle-ads 6734:6734
  Then access http://127.0.0.1:6734/authz-check/v1/is-allowed to check permissions.
```

如果你的`etcd`启用了 TLS，请在`values.yaml`中设置以下变量：

```yaml
store:
  etcd:
    endpoint: https://example-client.default.svc:2379
    etcdClientCertSecret: etcd-client-tls
    certFile: etcd-client.crt
    keyFile: etcd-client.key
    trustedCAFile: etcd-client-ca.crt
```

## 使用 Speedle

### 配置 `spctl`

_如果你在任何一个 Kubernetes 节点上运行`spctl`，PMS 的 IP 地址是 Kubernetes Service `Speedle`的 Cluster-IP. 在这个例子中，IP 地址是 10.101.117.233。如果你想在 Kubernetes 集群外运行`spctl`，请将 Kubernetes Service `Speedle`定义为一个`nodePort` Service，请参考 Kubernetes 文档https://kubernetes.io/docs/concepts/services-networking/service/。_

```bash
$ spctl config pms-endpoint http://10.108.146.148:6733/policy-mgmt/v1/
$ spctl config --list
cacert =
cert =
key =
pms-endpoint = http://10.108.146.148:6733/policy-mgmt/v1/
timeout = 5s
```

### 用`spctl`新建一个 Speedle Service

```bash
$ spctl create service test
service created
{"name":"test","type":"application"}
```

### 用`spctl`新建一条策略

```bash
$ spctl create policy -c "grant user jiefu read book" --service-name test
policy created
{"id":"ta55v3kyzux5ssiy3wwr","name":"","effect":"grant","permissions":[{"resource":"book","actions":["read"]}],"principals":[["user:jiefu"]]}
```

### 用`curl`测试策略

```bash
$ curl -X POST -d '{"subject":{"principals":[{"type":"user","name":"jiefu"}]},"serviceName":"test","resource":"book","action":"read"}' http://10.108.146.148:6734/authz-check/v1/is-allowed
{"allowed":true,"reason":0
```

## 疑难解答

如果 Speedle 的映像放在你私人的 Speedle Registry 上，在部署 Speedle 服务之前，请先添加一个 Kubernetes Secret。

```bash
kubectl create secret docker-registry <secret name> --docker-username=<user name> --docker-password="<Password>" --docker-email="your@email.com"
--docker-server=<server-add>
```

更多的关于 Kubernetes Secret 详细，请参考：
https://kubernetes.io/docs/concepts/configuration/secret/
