+++
title = "Quick Start"
description = "Up and running in under a minute"
weight = 10
draft = false
toc = false
tocheading = "h2"
bref = "Speedle is an open source authorization project that allows you to create policies using the user-friendly SPDL language, and manage the policies using the policy management framework. Speedle accepts authorization requests (that is, who can do what in what context), evaluates the applicable policies, and then returns a GRANT/DENY decision using the authorization decision service. Let's take a quick tour of how to use Speedle"
icon = "9. Getting Started.svg"
+++

## Prerequisites

### 1. Install the Go programming language

The minimum version required is 1.10.1.

See https://golang.org/doc/install

If you use Oracle Linux, simply run one of below commmands.
```bash
$ sudo yum install -y oracle-golang-release-el7
$ sudo yum install -y golang
```

### 2. Set the Go environment variables

`GOROOT` is the folder where you installed Go.  
`GOPATH` is the folder where your Go projects are stored.

```bash
$ export GOROOT=/scratch/tools/go
$ export GOPATH=/scratch/xuwwang/go
```


## Download Source Code and Build Speedle

```bash
$ go get github.com/oracle/speedle/cmd/...
$ ls $GOPATH/bin
spctl  speedle-ads  speedle-pms
```

Three binary files are generated in the `$GOPATH/bin` directory:

- spctl - Speedle command line interface
- speedle-pms - Speedle policy management service
- speedle-ads - Speedle authorization decision service

## Run Speedle

### 1. Start the Policy Management Service

```bash
$ cd $GOPATH/bin
$ ./speedle-pms --store-type file
```

Note: A default policy store file is created at `/tmp/speedle-test-file-store.json`.

### 2. Create policies using spctl

2.1.  Create a service container for the authorization and role policies.

```bash
$ ./spctl create service mysvc
```

2.2. Create authorization policies in the mysvc service.

```bash
$ ./spctl create policy -c "grant user user1 get,del res1" --service-name=mysvc
$ ./spctl create policy -c "grant role role2 get,del res2" --service-name=mysvc
```

2.3. Create a role policy in the mysvc service.

```bash
$ ./spctl create rolepolicy -c "grant user user2 role2 on res2" --service-name=mysvc
```

2.4. List all services

```bash
$ cd $GOPATH/bin
$ ./spctl get service --all
```

The content of service mysvc displays.

### 3. Start the Authorization Decision Service

```bash
$ ./speedle-ads --store-type file
```

### 4. Verify the authorization result

To see the policies you defined take effect, run these commands in a separate command window:

```bash
$ curl -X POST --data '{"subject":{"principals":[{"type":"user","name":"user1"}]},"serviceName":"mysvc","resource":"res1","action":"get"}' http://127.0.0.1:6734/authz-check/v1/is-allowed

$ curl -X POST --data '{"subject":{"principals":[{"type":"user","name":"user2"}]},"serviceName":"mysvc","resource":"res2","action":"get"}' http://127.0.0.1:6734/authz-check/v1/is-allowed
```

The result for both commands is allowed:true.

### TLS/HTTPS

For TLS configuration, see [Message security](../docs/security#message-security-tls).

Looking for more deployment types? See [Deployment page](../docs/deployment).
