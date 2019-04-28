+++
title = "Security"
description = "Secure your Speedle instances"
weight = 310
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["security", "guide", "tls", "https"]
categories = ["docs"]
bref = ""
+++

## Overview

### API Endpoint Security / Authentication and Authorization

Authentication and authorization for client requests (other than TLS mutual auth) are not supported by Speedle.

If you want to protect Speedle API endpoints, you can use any existing/stock solutions to secure these APIs (e.g. an API Gateway like Ambassador with tokens etc).

### Message security / TLS

TLS is a cryptographic protocol that provides communications security, it offers many different ways of exchanging keys for authentication, encrypting data, and guaranteeing message integrity.

Speedle supports using TLS to secure the message transport inside untrusted environments.
TLS mutual authentication is also supported for the authentication of the client.

In the bellow sections, we'll describe how to enable TLS in Speedle server as a standalone application and how to use Speedle CLI/curl to access the TLS-enabled Speedle server.

Looking for securing your kubernetes deployment with TLS? - click [here](../deployment)

## Prerequisites

In order to implement TLS, Speedle must have an associated `Certificate` configured for its external interface (IP address or DNS name) that accepts secure connections.

This certificate is cryptographically signed by a trusted third party. These are called `Certificate Authorities` (`CA`s). To obtain a signed certificate, you need to choose a CA and follow the instructions your chosen CA provides to obtain your certificate.

In your test enviroment, you can create a `self-signed` certificate. Self-signed certificates are simply user generated certificates which have not been signed by a well-known CA and are, therefore, not really guaranteed to be authentic at all. They are not suitable for production use.

For convenience, the cfssl tool provides an easy interface to certificate generation.
The steps to use [cfssl](https://github.com/cloudflare/cfssl) to generate self-signed certificates are as following:

### Install cfssl

```bash
curl -s -L -o /usr/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
curl -s -L -o /usr/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
```

Or install by Go:

```bash
go get -u github.com/cloudflare/cfssl/cmd/cfssl
go get -u github.com/cloudflare/cfssl/cmd/cfssljson
```

### Generate CA and Certificates

```bash
echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
# replace this with your ip address or dns
export ADDRESS=localhost,127.0.0.1
export NAME=server
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
export ADDRESS=
export NAME=client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

mv server.pem server.crt
mv server-key.pem server.key
cp ca.pem server-ca.crt

mv client.pem client.crt
mv client-key.pem client.key
cp ca.pem client-ca.crt
```

## TLS-enabled Speedle Server

### Command Line Flags

Use the following flags to configure TLS for `PMS` (Policy Management Service) and `ADS` (Authorization Decision Service).

| Name              | Value                | Default | Comments                                                                                             |
| ----------------- | -------------------- | ------- | ---------------------------------------------------------------------------------------------------- |
| insecure          | true, false          | false   | specifies when TLS is enabled or not, true: disabled, false: enabled.                                |
| cert              | TLS certificate path |         | specifies the path of the file containing the TLS certificate.                                       |
| key               | TLS private key path |         | specifies the path of the file containing the TLS private key.                                       |
| client-cert       | client CA path       |         | specifies the path of the file containing the trusted CA File for client certificate authentication. |
| force-client-cert | true, false          | false   | specifies if the client certificate authentication is forced or not.                                 |

### Example

```bash
$ speedle-pms --store-type file --endpoint="127.0.0.1:6733" --insecure=false --key=$tls_config_path/server.key --cert=$tls_config_path/server.crt --client-cert=$tls_config_path/client-ca.crt

$ speedle-ads --store-type file --endpoint="127.0.0.1:6735" --insecure=false --force-client-cert=true --key=$tls_config_path/server.key --cert=$tls_config_path/server.crt --client-cert=$tls_config_path/client-ca.crt
```

**Note: Please use absolute pathes for the cert files.**

## Use `spctl` CLI to Access TLS-enabled Speedle

### Command Line Flags

Use the following flags to configure TLS for Spctl CLI.

| Name              | Value                | Default | Comments                                                                                             |
| ----------------- | -------------------- | ------- | ---------------------------------------------------------------------------------------------------- |
| skipverify        | true, false          | false   | specifies if skip verifing PMS server's TLS certificate.                                             |
| cert              | TLS certificate path |         | specifies the path of the file containing the client TLS certificate.                                |
| key               | TLS private key path |         | specifies the path of the file containing the client TLS private key.                                |
| client-cert       | client CA path       |         | specifies the path of the file containing the trusted CA File for server certificate authentication. |
| force-client-cert | true, false          | false   | specifies if the client certificate authentication is forced or not.                                 |

### Example

```bash
$ spctl config skipverify false cacert $tls_config_path/server-ca.crt cert $tls_config_path/client.crt key $tls_config_path/client.key pms-endpoint "https://localhost:6733/policy-mgmt/v1/"
$ spctl create service test
$ spctl get service --all
```

## Use `curl` to Access TLS-enabled Speedle

### Example

```bash
$ curl --cacert $tls_config_path/server-ca.crt --cert $tls_config_path/client.crt --key $tls_config_path/client.key https://localhost:6733/policy-mgmt/v1/service
```
