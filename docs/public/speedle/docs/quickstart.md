Speedle, as an open source project in authorization area, allows you to create policies in user-friendly SPDL language, manage these policies. It accepts authorization requests (i.e. who can do what in what context), evaluates applicable policies and returns GRANT/DENY decision. This document guides you a quick tour with Speedle.

# Prerequisites

## 1. Install Golang with version 1.10.1

Please refer to https://golang.org/doc/install

## 2. Export GOPATH , GOROOT

```bash
  export GOROOT=/scratch/tools/go
  export GOPATH=/scratch/xuwwang/go
```

Notes:
GOROOT is the folder where your golang installed.
GOPATH is the folder where your projects stored

## 3. Install dep:

```bash
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

More info about dep, please refer to Installation part in https://github.com/golang/dep

# Sync or Download Speedle

## 1. Create bin, src folder under \$GOPATH

```bash
   cd $GOPATH
   mkdir bin
   mkdir src
```

## 2. Pull speedle source code into src folder

```bash
   cd $GOPATH/src
   git clone git@github.com:oracle/speedle.git
   mkdir -p github.com/oracle
   mv speedle github.com/oracle/
```

Notes:
you need to get account of gitlab and permission of this project.
the last two steps are used to make Oracle SSO works.

# Build Speedle

Run below commands to build source code

```bash
cd $GOPATH/src/github.com/oracle/speedle
make clean build
```

Then in \$GOPATH/bin, there would be 3 binary file generated
spctl, speedle-pms, speedle-ads

# Run Speedle

## 1. Start PMS (Policy Management Service)

```bash
cd $GOPATH/bin
./speedle-pms --store-type file
```

Notes: default policy store file is located in /tmp/speedle-test-file-store.json

## 2. Create policies via spctl command

### 2.1. Create services:

```bash
./spctl create service srv1
```

### 2.2. Create a policy:

```bash
./spctl create policy policy1 --pdl-command "grant user user1 get,del res1" --service-name=srv1
./spctl create policy policy2 --pdl-command "grant role role2 get,del res2" --service-name=srv1
```

### 2.3. Create a role policy:

```bash
./spctl create rolepolicy rolepolicy1 --pdl-command "grant user user2 role2 on res2" --service-name=srv1
```

### 2.4. Get all services:

```bash
cd $GOPATH/bin
./spctl get service --all
```

## 3. Start ADS (Authorization Decision Service)

```bash
./speedle-ads --store-type file
```

## 4. Check authorization result by curl command (result is allowd=true)

```bash
curl -X POST --data '{"subject":{"principals":[{"type":"user","name":"user1"}]},"serviceName":"srv1","resource":"res1","action":"get"}' http://127.0.0.1:6734/authz-check/v1/is-allowed

curl -X POST --data '{"subject":{"principals":[{"type":"user","name":"user2"}]},"serviceName":"srv1","resource":"res2","action":"get"}' http://127.0.0.1:6734/authz-check/v1/is-allowed
```

Now you see the policies you defined above take effect.

# TLS and HTTPS

HTTPS is configured by specifying TLS credentials via command line flags.

## 1. Generate Self-signed Certificates

### 1.1. Download cfssl

```bash
mkdir ~/bin
curl -s -L -o ~/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
curl -s -L -o ~/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
chmod +x ~/bin/{cfssl,cfssljson}
export PATH=$PATH:~/bin
```

### 1.2. Create Directory to Store Certificates:

```bash
mkdir ~/cfssl
cd ~/cfssl
```

### 1.3. Generate CA and Certificates

```bash
echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
export ADDRESS=192.168.122.68,ext1.example.com,coreos1.local,coreos1
export NAME=server
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
export ADDRESS=
export NAME=client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
```

### 1.4. Verify data

```bash
openssl x509 -in ca.pem -text -noout
openssl x509 -in server.pem -text -noout
openssl x509 -in client.pem -text -noout
```

## 2. PMS and ADS

### 2.1. Command line flags:

```
--insecure=true|false specifies when TLS is enabled or not, true: disabled, false: enabled.
--cert=<path> specifies the path of the file containing the TLS certificate.
--key=<path> specifies the path of the file containing the TLS private key.
--client-cert=<path> specifies the path of the file containing the trusted CA File for client certificate authentication.
--force-client-cert=true|false specifies if the client certificate authentication is forced or not.
```

### 2.2. Start Speedle with TLS enabled

```bash
./speedle-pms --store-type file --insecure=false --key=/home/diazhao/cfssl/server-key.pem --cert=/home/diazhao/cfssl/server.pem --client-cert=/home/diazhao/cfssl/ca.pem --endpoint="127.0.0.1:6733"
./speedle-ads --store-type file --insecure=false --key=/home/diazhao/cfssl/server-key.pem --cert=/home/diazhao/cfssl/server.pem --client-cert=/home/diazhao/cfssl/ca.pem --endpoint="127.0.0.1:6735"
use absolute path for the cert files.
```

### 2.3. Access the API with HTTPS

curl -k https://localhost:6733/policy-mgmt/v1/

    We have to use cURL's -k/--insecure flag because we are using a self-signed certificate.

## 3. Spctl

### 3.1. Command line flags:

```
--skipverify=true|false specifies if skip verifing PMS server's TLS certificate.
--cert=<path> specifies the path of the file containing the TLS certificate.
--key=<path> specifies the path of the file containing the TLS private key.
--cacert=<path> specifies the path of the file containing the trusted CA File for server certificate authentication.
--endpoint= pms endpoint like 127.0.0.1:6733
```

### 3.2. Start Spctl to Access TLS Enabled PMS

```bash
./spctl --cert=/home/diazhao/cfssl/client.pem --key=/home/diazhao/cfssl/client-key.pem --cacert=/home/diazhao/cfssl/ca.pem --skipverify=true --pms-endpoint="https://127.0.0.1:6733/policy-mgmt/v1/" 
```

---

**Looking for more deployment types? please check [deployment page](../deployment)**
