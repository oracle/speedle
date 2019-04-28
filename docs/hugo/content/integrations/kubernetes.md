+++
title = "Kubernetes"
description = "Up and running with Kubernetes authorization webhook"
weight = 200
draft = false
toc = true
iconurl = "22. Kubernetes.svg"
tocheading = "h3"
tags = ["k8s", "kubernetes", "webhook"]
categories = ["integration", "cloudnative"]
bref = "This sample demonstrates how to connect Kubernetes with Speedle via Kubernetes Authorization Webhook such that Kubernetes operations can be protected by Speedle authorization engine."
+++

## Kubernetes Authorization Webhook Sample

### Build

```bash
git clone git@github.com:oracle/speedle.git

# make sure the speedle golang adsclient code is in your $GOPATH/src package.

cp -r samples/adsclient/go/src/speedle/ $GOPATH/src

# get additional required packages:

$go get k8s.io/api/authorization/v1beta1
$go get k8s.io/client-go/kubernetes
$go get k8s.io/client-go/tools/clientcmd

# build
cd samples/integration/kubernetes-integration
go build webhook.go
#Executable webhook could be found in same folder
```

### Create Certificates

Generate self-signed certificates used by webhook by refering to [this document](https://kubernetes.io/docs/concepts/cluster-administration/certificates/)
BTW, let's assume we store the generated certificates and other needed files to path `/path/to/webhook/`, we will use this path in below sections.

### Start Webhook

```bash
sudo ./webhook -key /path/to/webhook/server.key \
-cert /path/to/server.crt \
-client-cert /etc/kubernetes/pki/ca.crt \
-cluster-name kubernates \
-speedle-host <speedle host ip> \
-kubeconfig ~/.kube/config
```

Default listen port is `8843`, can use `./webhook -h` to see the usage.

### Enable webhook

#### 1. Mount the localhost folder /path/to/webhook/ to kube-apiserver pod

Since this is just a sample, so we use hostpath to mount the `/path/to/webhook`(which contains all the certificates and other needed files) into `kube-apiserver` pod (mountPath is `/etc/webhook`). Modify `/etc/kubernetes/manifests/kube-apiserver.yaml`, add below `-hostPath` and `-mountPath` section:

```yaml
  volumeMounts:
    ...
    - mountPath: /etc/webhook
      name: webhook
      readOnly: true

  volumes:
  ...
  - hostPath:
      path: /path/to/webhook
      type: DirectoryOrCreate
    name: webhook
```

#### 2. Modify the cert file path and server URL in webhook.yaml

All the path used in `webhook.yaml` should be the `mountPath` in `kube-apiserver` pod (`/etc/webhook`).

```yaml
# clusters refers to the remote service.
clusters:
    - name: name-of-remote-authz-service
      cluster:
          # CA for verifying the remote service.
          certificate-authority: /etc/webhook/ca.crt
          # URL of remote service to query. Must use 'https'. May not include parameters.
          server: https://<webhook host ip>:8443/

# users refers to the API Server's webhook configuration.
users:
    - name: name-of-api-server
      user:
          client-certificate: /etc/webhook/server.crt # cert for the webhook plugin to use
          client-key: /etc/webhook/server.key # key matching the cert
```

#### 3. Enable webhook in `kube-apiserver.yaml`

Enable three authorization modes, `Node`, `RBAC`, and `Webhook`
(invoking Speedle underneath), and update webhook.yaml, user.csv path
like below in `/etc/kubernetes/manifests/kube-apiserver.yaml`

```bash
- --authorization-mode=Node,RBAC,Webhook
- --authorization-webhook-config-file=/etc/webhook/webhook.yaml
- --authorization-webhook-cache-authorized-ttl=110ms
- --authorization-webhook-cache-unauthorized-ttl=110ms
- --basic-auth-file=/etc/webhook/user.csv
```

### Testing

-   Check if user `joe` has permission to read pod `etcd-slc00bqb` in namespace `kube-system`

    1.Check if user `joe` has permission to read pod `etcd-slc00bqb` in namespace `kube-system`, this pod has a label defined `component=etcd`.

```bash
curl -k -v -u joe:joe https://10.96.0.1:443/api/v1/namespaces/kube-system/pods/etcd-slc00bqb

{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {

  },
  "status": "Failure",
  "message": "pods \"etcd-slc00bqb\" is forbidden: User \"joe\" cannot get pods in the namespace \"kube-system\"",
  "reason": "Forbidden",
  "details": {
    "name": "etcd-slc00bqb",
    "kind": "pods"
  },
  "code": 403
* Connection #0 to host 10.96.0.1 left intact
}
```

Currently no policy defined for user joe, so code 403 is expected.

2.Now grant joe the permission in speedle service

```bash
$./spctl create service kubernetes
service created
{"name":"kubernetes","type":"application","metadata":{"createby":"","createtime":"2018-11-28T21:13:17-08:00"}}

$./spctl create policy joe-pod --service-name=kubernetes -c "grant user joe get expr:/pods/.* if labels_component == \"etcd\""
policy created
{"id":"3ttrj5xtj2k4teyz4ja4","name":"joe-pod","effect":"grant","permissions":[{"resourceExpression":"/pods/.*","actions":["get"]}],"principals":[["user:joe"]],"condition":"labels_component == \"etcd\"","metadata":{"createby":"","createtime":"2018-11-29T18:40:49-08:00"}}
```

Above commands create a service named `kubernetes`, created a policy, user `joe` has permission to get all pods if pod has a label `component=etcd`.

3.Check again, expected result is `allowed`

```
$curl -k -v -u joe:joe https://10.96.0.1:443/api/v1/namespaces/kube-system/pods/etcd-slc00bqb
{
  "kind": "Pod",
  "apiVersion": "v1",
  "metadata": {
    "name": "etcd-slc00bqb",
    "namespace": "kube-system",
    "selfLink": "/api/v1/namespaces/kube-system/pods/etcd-slc00bqb",
    "uid": "1aba7673-e879-11e8-9c67-00163e7a4530",
    "resourceVersion": "533818",
    "creationTimestamp": "2018-11-15T01:52:27Z",
    "labels": {
      "component": "etcd",
      "tier": "control-plane"
    },
    "annotations": ....
}
```

you can also try other pod which doesn't have label `component=etcd`. you will get code 403

-   Check if user `joe` has permission to read pods belonging to a deployment `kube-dns`

There only show how to define a policy, allows user `joe` has permission to read all pods owned by a deployment "kube-dns".

```bash
$./spctl create policy joe_policy -c "grant user joe get,watch expr:/pods/.* if owner_Deployment == \"kube-dns\"" --service-name=kubernetes
```
