+++
title = "Deployment"
description = "Recommended deployment topology for Speedle"
weight = 60
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
categories = ["docs"]
bref = "Deploy Speedle in embedded mode or as a service"
+++

## Deploy Speedle in embedded mode

Speedle works as a Go library. To deploy Speedle in embedded mode, simply pull Speedle from the repository and import it in golang code.
The Policy Management Service (PMS) golang API is called to do policy management, and the Authorization Decision Service (ADS) golang API is called to do runtime authorization checks.

## Deploy Speedle as a service

You can deploy Speedle as a service on Kubernetes in these modes:

- **Dev mode** - where Speedle components run using a file as a policy store.

- **Production mode** - where Speedle components run on an etcd cluster as a policy store.

In both modes, users can access the Speedle service using the REST API or the `spctl` CLI tool.

### Deploy Speedle on Kubernetes

Deploy Speedle on Kubernetes in both dev and production modes using Helm, or manually. Both methods are described here.

#### Prerequisites

- Golang 1.10.0 environment installed on your development box
- Docker 1.12 or higher installed on your development box
- A Docker registry to publish Docker images
- The latest code from the Speedle repository synced to folder `$GOPATH/src/github.com/oracle/speedle/` on your dev box
- A Kubernetes cluster or an OKE cluster
- A Docker-registry secret. To create one:

```bash
kubectl create secret docker-registry reg-speedle --docker-server=$DOCKER_LOGIN_SERVER --docker-username=$DOCKER_LOGIN_USER --docker-password=$DOCKER_LOGIN_PASSWORD --docker-email="youremail@yourcorp.com"
```

#### Deploy Speedle using Helm

You can use Helm to deploy Speedle in both dev and production modes.

##### Dev mode

```bash
$ helm install -n speedle-dev deployment/helm/speedle-dev
NAME:   speedle-dev
LAST DEPLOYED: Tue Nov 27 23:21:18 2018
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/Service
NAME         TYPE       CLUSTER-IP      EXTERNAL-IP  PORT(S)            AGE
speedle-dev  ClusterIP  10.101.133.188  <none>       6733/TCP,6734/TCP  0s

==> v1beta2/Deployment
NAME         DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
speedle-dev  1        1        1           0          0s

==> v1/Pod(related)
NAME                          READY  STATUS             RESTARTS  AGE
speedle-dev-7568f8dc44-mwh5m  0/2    ContainerCreating  0         0s
```

NOTES:

```bash
$ kubectl port-forward svc/speedle-dev 6733:6733 6734:6734
```

- To manage policies, access http://127.0.0.1:6733/policy-mgmt/v1/
- To check permissions: access http://127.0.0.1:6734/authz-check/v1/is-allowed

##### Production mode

Before you deploy Speedle in production mode, you need to install and configure etcd. If you are using TLS, you also have to create a Kubernetes secret.

###### Install etcd

The best way to install etcd is using the etcd-operator. Here we show the simplest steps to install a sample etcd cluster. For more information, see [etcd-operator](https://github.com/coreos/etcd-operator) on GitHub.

1. Install the etcd-operator using Helm

```bash
$ helm install stable/etcd-operator --name my-release
```

2. Create an etcd cluster

```bash
$ kubectl create -f https://raw.githubusercontent.com/coreos/etcd-operator/master/example/example-etcd-cluster.yaml
```

A 3 member etcd cluster is created.

```bash
$ kubectl get pods
NAME                            READY     STATUS    RESTARTS   AGE
example-etcd-cluster-gxkmr9ql7z   1/1       Running   0          1m
example-etcd-cluster-m6g62x6mwc   1/1       Running   0          1m
example-etcd-cluster-rqk62l46kw   1/1       Running   0          1m
```

Create a TLS-enabled etcd cluster using a static cluster TLS policy as described in
[Cluster TLS policy](https://github.com/coreos/etcd-operator/blob/master/doc/user/cluster_tls.md)

###### Create a Kubernetes secret for TLS

If you want to enable TLS for the Speedle server (see [Security](../security)), you need to create the TLS Kubernetes secret as described in these steps.

1. Install the CloudFlare PKI/TLS Toolkit (CFSSL)

```bash
curl -s -L -o /usr/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
curl -s -L -o /usr/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
```

2. Generate the CA and certificates

```bash
echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
# replace this address with your IP address or DNS
export ADDRESS=localhost,127.0.0.1,speedle-ads.default.svc
export NAME=server
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
export ADDRESS=
export NAME=client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
```

```bash
mv server.pem server.crt
mv server-key.pem server.key
cp ca.pem server-ca.crt
```

```bash
mv client.pem client.crt
mv client-key.pem client.key
cp ca.pem client-ca.crt
```

3. Create the TLS secret

```bash
kubectl create secret generic speedle-server-tls --from-file=server-ca.crt --from-file=server.crt --from-file=server.key
```

###### Deploy Speedle

- To deploy Speedle with TLS disabled:

```bash
$ helm install -n speedle deployment/helm/speedle-prod --set store.etcd.endpoint=http://example-etcd-cluster-client:2379,tls=

```

- To deploy Speedle with TLS enabled:

```bash
$ helm install -n speedle deployment/helm/speedle-prod --set store.etcd.endpoint=http://example-etcd-cluster-client:2379,tls.certSecret=speedle-server-tls,tls.certFile=server.crt,tls.keyFile=server.key,tls.trustedCAFile=server-ca.crt,tls.forceClientCert=true

```

If the etcd cluster is **TLS enabled**, set the following variables (see [operatorSecret](https://github.com/coreos/etcd-operator/blob/master/doc/user/cluster_tls.md#operatorsecret) in the etcd-operator documentation on GitHub):

```yaml
store:
  etcd:
    endpoint: https://example-etcd-cluster-client.default.svc:2379
    etcdClientCertSecret: etcd-client-tls
    certFile: etcd-client.crt
    keyFile: etcd-client.key
    trustedCAFile: etcd-client-ca.crt
```

#### Deploy Speedle manually

Before deploying Speedle manually, you need to build Speedle and build and push the Speedle Docker images.

##### Build Speedle

You can build Speedle by executing the following command. This command assumes that you have already set the GOPATH environment variable, and that the Speedle code is located in the directory `$GOPATH/src/github.com/oracle/speedle/`

```bash
cd $GOPATH/src/github.com/oracle/speedle
make
```

When the build completes, three files are located in the `$GOPATH/bin` directory:

- spctl
- speedle-pms
- speedle-ads

##### Build and push Speedle Docker images

```bash
export pmsImageRepo=Docker repository of the Speedle policy management service
export pmsImageImageVersion=Docker image version of the Speedle PMS
export adsImageRepo=Docker repository of the Speedle authorization decision service
export adsImageImageVersion=Docker image version of the Speedle ADS

cd $GOPATH/src/github.com/oracle/speedle
make image
```

##### Deploy Speedle in dev mode

###### Create the dev mode Kubernetes deployment

Create a YAML file, named `speedle-dev.yaml`, using the values appropriate for your environment, that can deploy both speedle-pms and speedle-ads in one Kubernetes deployment. A sample `speedle-dev.yaml`file is available in the [Speedle GIT repository](https://github.com/oracle/speedle/blob/master/deployment/k8s/speedle-dev.yaml).

```yaml
kind: Service
apiVersion: v1
metadata:
  name: speedle-dev
spec:
  selector:
    app: speedle
  ports:
    - protocol: TCP
      port: 6733
      targetPort: 6733
      name: pms
    - protocol: TCP
      port: 6734
      targetPort: 6734
      name: ads

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: speedle
  labels:
    app: speedle
spec:
  replicas: 1
  selector:
    matchLabels:
      app: speedle
      name: ads
  template:
    metadata:
      labels:
        app: speedle
        name: ads
    spec:
      containers:
        - name: ads
          image: r.authz.fun/speedle-ads:v0.1  // please update image location
          ports:
            - containerPort: 6734
          volumeMounts:
            - mountPath: /var/lib/speedle
              name: policy-store
        - name: pms
          image: r.authz.fun/speedle-pms:v0.1  // please update image location
          ports:
            - containerPort: 6733
          volumeMounts:
            - mountPath: /var/lib/speedle
              name: policy-store
      volumes:
        - name: policy-store
          hostPath:
            path: /home/opc/policystore   // Depends on your file system, there should be a valid file "policies.json" under this folder. A valid policy store should contain at least one line `{}` for an empty policy store.
      nodeSelector:
        kubernetes.io/hostname: host-ad1-vm1-2-1  // please replace host-ad1-vm1-2-1 with actual node
```

###### Install the dev mode Kubernetes deployment

Use the `kubectl create` command to create the Speedle deployment service using the `speedle-dev.yaml` file that you created in the previous step.

```bash
kubectl create -f speedle-dev.yaml
```

You can then use the `kubectl get` command to find one Speedle deployment, one pod and one service.

```bash
$ kubectl get deployment speedle
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
speedle   1         1         1            1           37m
$ kubectl get pods -l app=speedle
NAME                       READY     STATUS    RESTARTS   AGE
speedle-65f68c67fd-qqzhj   2/2       Running   0          38m
$ kubectl get services
NAME      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
speedle-dev   ClusterIP   10.108.146.148   <none>        6733/TCP,6734/TCP   39m
```

##### Deploy Speedle in production mode

###### Create the production mode Kubernetes deployment

Create a YAML file, named `speedle-prod.yaml`, using the values appropriate for your environment, that can deploy both speedle-pms and speedle-ads in one Kubernetes deployment. A sample `speedle-prod.yaml` file is available
in the [Speedle GIT repository](https://github.com/oracle/speedle/blob/master/deployment/k8s/speedle-prod.yaml)

In this sample deployment file, etcd runs in one Kubernetes Pod.

A high availability etcd cluster depends on the user's environment. If you want to run Speedle on an etcd cluster, see the [etcd Clustering Guide](https://coreos.com/etcd/docs/latest/v2/clustering.html). You must first create the etcd cluster, then modify the `speedle-prod.yaml` file with the etcd
cluster settings, and lastly deploy the modified `speedle-prod.yaml` file.

```yaml
kind: Service
apiVersion: v1
metadata:
  name: speedle-pms
spec:
  selector:
    app: speedle
    name: pms
  ports:
    - protocol: TCP
      port: 6733
      targetPort: 6733
      name: pms

---
kind: Service
apiVersion: v1
metadata:
  name: speedle-ads
spec:
  selector:
    app: speedle
    name: ads
  ports:
    - protocol: TCP
      port: 6734
      targetPort: 6734
      name: ads

---
kind: Service
apiVersion: v1
metadata:
  name: speedle-etcd
spec:
  selector:
    app: speedle
    name: etcd
  ports:
    - protocol: TCP
      port: 2379
      targetPort: 2379
      name: etcd

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: speedle-etcd
  labels:
    app: speedle
    name: etcd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: speedle
      name: etcd
  template:
    metadata:
      labels:
        app: speedle
        name: etcd
    spec:
      containers:
        - name: pms
          image: quay.io/coreos/etcd:v3.2
          command: ['/usr/local/bin/etcd']
          args:
            [
              '--listen-client-urls',
              'http://0.0.0.0:2379',
              '--advertise-client-urls',
              'http://0.0.0.0:2379',
            ]
          ports:
            - containerPort: 2379

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: speedle-pms
  labels:
    app: speedle
    name: pms
spec:
  replicas: 3
  selector:
    matchLabels:
      app: speedle
      name: pms
  template:
    metadata:
      labels:
        app: speedle
        name: pms
    spec:
      containers:
        - name: pms
          image: r.authz.fun/speedle-pms:v0.1  // please update image location
          command: ['speedle-pms']
          args:
            [
              '--endpoint',
              '0.0.0.0:6733',
              '--store-type',
              'etcd',
              '--etcdstore-endpoint',
              'speedle-etcd:2379',
            ]
          ports:
            - containerPort: 6733

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: speedle-ads
  labels:
    app: speedle
    name: ads
spec:
  replicas: 3
  selector:
    matchLabels:
      app: speedle
      name: ads
  template:
    metadata:
      labels:
        app: speedle
        name: ads
    spec:
      containers:
        - name: ads
          image: r.authz.fun/speedle-ads:v0.1  // please update image location
          command: ['speedle-ads']
          args:
            [
              '--endpoint',
              '0.0.0.0:6734',
              '--store-type',
              'etcd',
              '--etcdstore-endpoint',
              'speedle-etcd:2379',
            ]
          ports:
            - containerPort: 6734
```

###### Install the production mode Kubernetes deployment

Use the `kubectl create` command to create a Speedle deployment service using the `speedle-prod.yaml` file that you created in the previous step.

```bash
kubectl create -f speedle-prod.yaml
```

You can then use the `kubectl get` command to find one Speedle deployment, one pod and one service.

```bash
$ kubectl get deployment speedle
NAME          DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
speedle-pms   1         1         1            1           37m
speedle-ads   3         3         3            3           37m
speedle-etcd  1         1         1            1           37m

$ kubectl get services
NAME          TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
speedle-pms   ClusterIP   10.108.146.148   <none>        6733/TCP   38m
speedle-ads   ClusterIP   10.108.146.136   <none>        6734/TCP   38m
speedle-etcd  ClusterIP   10.108.146.118   <none>        2379/TCP   38m
```

### Try it out

#### Prerequisites

- Ensure that you have built Speedle as described in <a href="#build-speedle">Build Speedle</a>.
- Verify that the files `spctl`, `speedle-pms` and `speedle-ads` were created in the `$GOPATH/bin` after executing the build.

#### PMS and ADS endpoints

If you access the PMS/ADS service on any of the Kubernetes nodes, you can use the CLUSTER-IPs of the Kubernetes services of speedle-dev (in dev mode) or speedle-pms/speedle-ads (in production mode).

If you want to run the Speedle CLI `spctl` out of the Kubernetes cluster, you can port-forward the services to your local machine.

Dev mode:

```bash
kubectl port-forward svc/speedle-dev 6733:6733 6734:6734 &
```

Production mode:

```bash
kubectl port-forward svc/speedle-pms 6733:6733 &
kubectl port-forward svc/speedle-ads 6734:6734 &
```

Or, you can change the service type from ClusterIP to NodePort or LoadBalancer as described in the Kubernetes [Services](https://kubernetes.io/docs/concepts/services-networking/service/) documentation.

The following steps assume that you have port-forwarded the PMS/ADS service locally.

#### Configure the pms-endpoint using `spctl`

```bash
$ spctl config pms-endpoint http://localhost:6733/policy-mgmt/v1/
$ spctl config --list
cacert =
cert =
key =
pms-endpoint = http://localhost:6733/policy-mgmt/v1/
timeout = 5s

# If TLS is enabled
$ spctl config skipverify false cacert $tls_config_path/server-ca.crt cert $tls_config_path/client.crt key $tls_config_path/client.key pms-endpoint "https://localhost:6733/policy-mgmt/v1/"
$ spctl config --list
cacert = <tls_config_path>/server-ca.crt
cert = <tls_config_path>/client.crt
key = <tls_config_path>/client.key
pms-endpoint = https://localhost:6733/policy-mgmt/v1/
skipverify = false
timeout = 5s
```

#### Create a service named "test" using `spctl`

```bash
$ spctl create service test
service created
{"name":"test","type":"application"}
```

#### Create a policy using `spctl`

```bash
$ spctl create policy -c "grant user jiefu read book" --service-name test
policy created
{"id":"ta55v3kyzux5ssiy3wwr","name":"","effect":"grant","permissions":[{"resource":"book","actions":["read"]}],"principals":[["user:jiefu"]]}
```

#### Trigger a policy evaluation

- If TLS is disabled:

```bash
$ curl -X POST -d '{"subject":{"principals":[{"type":"user","name":"jiefu"}]},"serviceName":"test","resource":"book","action":"read"}' https://localhost:6734/authz-check/v1/is-allowed
{"allowed":true,"reason":0}
```

- If TLS is enabled:

```bash
$ curl -i --cacert $tls_config_path/server-ca.crt --cert $tls_config_path/client.crt --key $tls_config_path/client.key -X POST -d '{"subject":{"principals":[{"type":"user","name":"jiefu"}]},"serviceName":"test","resource":"book","action":"read"}' https://localhost:6734/authz-check/v1/is-allowed

```
