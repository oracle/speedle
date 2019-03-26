# Deploy Speedle

## Deploy Speedle in Embedded Mode

Speedle works as a golang library.  
Just pull speedle from repository and import it in golang code.  
PMS golang API is called to do policy management, and ADS golang API is called to do runtime authorization check.

T.B.D.

## Deploy Speedle as a Service

This chapter introduces how to deploy Speedle as a service.

Two modes are supported:

1. Dev mode

In `Dev mode`, Speedle components run with a `file` as a policy store, and users can
visit Speedle service via REST API or CLI tool.

2. Production mode

In `Production mode`, Speedle components run on an `etcd` cluster as a policy store, and users can visit Speedle service via REST API or CLI tool.

### Prerequisite

1. Golang 1.10.0 environment has already been installed on dev box.
2. Docker 1.12 or higher has already been installed on dev box.
3. A docker registry to publish docker images.
4. Sync up the latest code from Speedle repo to folder \$GOPATH/src/github.com/oracle/speedle/ on dev box.
5. A Kubernetes cluster or an OKE cluster
6. Create the docker-registry secret:

```bash
kubectl create secret docker-registry reg-speedle --docker-server=$DOCKER_LOGIN_SERVER --docker-username=$DOCKER_LOGIN_USER --docker-password=$DOCKER_LOGIN_PASSWORD
```

### Deploy Speedle by Helm

#### Dev Mode

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


NOTES:
$ kubectl port-forward svc/speedle-dev 6733:6733 6734:6734
  Then access http://127.0.0.1:6733/policy-mgmt/v1/ to manage policies, access http://127.0.0.1:6734/authz-check/v1/is-allowed to check permissions.

```

#### Production Mode

##### Install etcd

It is recommendted to install etcd by etcd-operator: https://github.com/coreos/etcd-operator

Install etcd-operator by helm:

```bash
$ helm install stable/etcd-operator --name my-release
```

Create an etcd cluster:

```bash
$ kubectl create -f https://raw.githubusercontent.com/coreos/etcd-operator/master/example/example-etcd-cluster.yaml
```

A 3 member etcd cluster will be created.

```bash
$ kubectl get pods
NAME                            READY     STATUS    RESTARTS   AGE
example-etcd-cluster-gxkmr9ql7z   1/1       Running   0          1m
example-etcd-cluster-m6g62x6mwc   1/1       Running   0          1m
example-etcd-cluster-rqk62l46kw   1/1       Running   0          1m
```

Create a TLS enabled etcd cluster with Static cluster TLS policy:
https://github.com/coreos/etcd-operator/blob/master/doc/user/cluster_tls.md

##### Install Speedle

Update the etcd endpoint in `values.yaml`:

```
store:
  type: etcd
  etcd:
    endpoint: http://<etcdServiceName>:<etcdServicePort>
```

Or you can override it by using `--set store.etcd.endpoint=http://<etcdServiceName>:<etcdServicePort>` when running `helm install`.

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

If the etcd cluster is TLS enabled, please set the following variables:

```yaml
store:
    etcd:
        endpoint: https://example-client.default.svc:2379
        etcdClientCertSecret: etcd-client-tls
        certFile: etcd-client.crt
        keyFile: etcd-client.key
        trustedCAFile: etcd-client-ca.crt
```

### Deploy Speedle Manually

#### Build Speedle

Assume system environment variable GOPATH has already been set, and speedle code could be found under directory \$GOPATH/src/github.com/oracle/speedle/

```bash
cd $GOPATH/src/github.com/oracle/speedle
make
```

Files "spctl", "speedle-pms" and "speedle-ads" should be found under \$GOPATH/bin after building.

#### Build and Push Speedle Docker Images

```bash
export pmsImageRepo=docker repository of pms
export pmsImageImageVesion=docker image version of pms
export adsImageRepo=docker repository of ads
export adsImageImageVesion=docker image version of ads

cd $GOPATH/src/github.com/oracle/speedle
make image
```

#### Deploy Speedle in Dev Mode

##### Create the Dev Mode Kubernetes Deployment

This is a sample YAML file to deploy both speedle-pms and speedle-ads in one kubernetes deployment. This file can be found from GIT repo: https://github.com/oracle/speedle/blob/master/deployment/k8s/speedle-dev.yaml

```yaml
kind: Service
apiVersion: v1
metadata:
    name: speedle
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
                      path: /home/opc/policystore   // Depends on your file system, there should be a valid file "policies.json" under this folder. A valid policy store should have at least one line "{}" (don't include the quotas) for empty policy store.
            nodeSelector:
                kubernetes.io/hostname: sphinx-ad1-vm1-2-1  // please replace sphinx-ad1-vm1-2-1 with actual node
```

##### Install the Dev Mode Kubernetes Deployment

Assume the YAML file is stored as name "speedle-dev.yaml".

```bash
kubectl create -f speedle-dev.yaml
```

Then you can find one deployment, one pod and one service for speedle.

```bash
$ kubectl get deployment speedle
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
speedle   1         1         1            1           37m
$ kubectl get pods -l app=speedle
NAME                       READY     STATUS    RESTARTS   AGE
speedle-65f68c67fd-qqzhj   2/2       Running   0          38m
$ kubectl get services
NAME      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
speedle   ClusterIP   10.108.146.148   <none>        6733/TCP,6734/TCP   39m
```

#### Deploy Speedle in Production Mode

##### Create the Production Mode Kubernetes Deployment

This is a sample YAML file to deploy both speedle-pms and speedle-ads in one kubernetes
deployment. This file can be found from GIT repo:
https://github.com/oracle/speedle/blob/master/deployment/k8s/speedle-prod.yaml.

In this deployment file, `etcd` runs in one Kubernetes POD. An HA `etcd` cluster
depends on users environment, if want to run Speedle on an `etcd` cluster,
please follow https://coreos.com/etcd/docs/latest/v2/clustering.html to create
an `etcd` cluster first, and modify speedle-prod.yaml occording to the `etcd`
cluster settings, then deploy the modified speedle-prod.yaml.

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
                      ['--listen-client-urls', 'http://0.0.0.0:2379', '--advertise-client-urls', 'http://0.0.0.0:2379']
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

##### Install the Production Mode Kubernetes Deployment

Assume the YAML file is stored as name "speedle-prod.yaml".

```bash
kubectl create -f speedle-prod.yaml
```

Then you can find one deployment, one pod and one service for speedle.

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

## Try it out

### Config `spctl`

_If you run spctl on any of the kubernetes nodes, the IP address in PMS endpoint can be the CLUSTER-IP of the k8s service "speedle", in this example, the IP Address should be 10.108.146.148; or if you want to run spctl out of the cluster, please define the service as a nodePort service, please refer to k8s doc https://kubernetes.io/docs/concepts/services-networking/service/ for details._

```bash
$ spctl config pms-endpoint http://10.108.146.148:6733/policy-mgmt/v1/
$ spctl config --list
cacert =
cert =
key =
pms-endpoint = http://10.108.146.148:6733/policy-mgmt/v1/
timeout = 5s
```

### Create a Service "test" with `spctl`

```bash
$ spctl create service test
service created
{"name":"test","type":"application"}
```

### Create a Policy

```bash
$ spctl create policy -c "grant user jiefu read book" --service-name test
policy created
{"id":"ta55v3kyzux5ssiy3wwr","name":"","effect":"grant","permissions":[{"resource":"book","actions":["read"]}],"principals":[["user:jiefu"]]}
```

### Trigger a Policy Evalution

```bash
$ curl -X POST -d '{"subject":{"principals":[{"type":"user","name":"jiefu"}]},"serviceName":"test","resource":"book","action":"read"}' http://10.108.146.148:6734/authz-check/v1/is-allowed
{"allowed":true,"reason":0
```

## TroubleShooting

If you need pull docker image from a private server, you creat a secret and add it to your yaml file. Try,

```bash
kubectl create secret docker-registry <secret name> --docker-username=<user name> --docker-password="<Password>" --docker-email="your@email.com"
--docker-server=<server-add>
```

More information about k8s secret, please refer k8s doc
https://kubernetes.io/docs/concepts/configuration/secret/

```

```
