+++
title = "安全"
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

## 概要

### API Endpoint 安全/认证和授权

Speedle API Endpoint 自身仅支持 TLS mutual auth（TLS 相互身份验证）这种认证和授权.

您可以使用像 Ambassador API Gateway 等既存的方案对 Speedle API endpoint 进行保护。

### 消息安全 / TLS

TLS 是一种提供通信安全性的加密协议，它提供了许多不同的交换密钥进行身份验证，加密数据和保证消息完整性的方法。

Speedle 支持使用 TLS 来保护在非信任环境中的消息传输，也支持使用 TLS mutual auth 来验证客户端身份。

在下文中，我们将介绍如何在 Speedle（作为一个独立应用程序）中启用 TLS，以及如何使用 Speedle CLI / curl 访问启用 TLS 的 Speedle 服务。

如果您希望通过 TLS 保护在 Kubernetes 中的 Speedle 部署，请单击 [此处](../deployment)

## 先决条件

为了启用 TLS，Speedle 必须为其外部接口（IP 地址或 DNS 名称）配置相关的“证书”，以接受安全连接。

此证书由受信任的第三方（证书颁发机构）加密签名。要获取签名证书，您需要选择 CA 并按照所选 CA 提供的说明获取证书。

在测试环境中，可以创建一个“自签名”证书。自签名证书是未经知名 CA 签名的用户自己生成的证书，无法保证证书的真实性，因此它们不适合在生产环境中使用。

为方便起见，您可以使用 [cfssl](https://github.com/cloudflare/cfssl) 工具来生成自签名证书，步骤如下：

### 安装 cfssl

```bash
curl -s -L -o /usr/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
curl -s -L -o /usr/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
```

或者使用 Go 来安装：

```bash
go get -u github.com/cloudflare/cfssl/cmd/cfssl
go get -u github.com/cloudflare/cfssl/cmd/cfssljson
```

### 生成 CA 和证书

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

## 在 Speedle 中启用 TLS

### 命令行参数

Speedle `PMS`（Policy Management Service）和`ADS`（Authorization Decision Service）的 TLS 相关配置参数：

| 名称              | 值                   | 默认  | 注释                                             |
| ----------------- | -------------------- | ----- | ------------------------------------------------ |
| insecure          | true, false          | false | 指定是否为非安全方式， true: 禁用, false: 启用。 |
| cert              | TLS certificate path |       | 指定 TLS 证书文件的目录。                        |
| key               | TLS private key path |       | 指定 TLS 私钥文件的目录。                        |
| client-cert       | client CA path       |       | 指定用于客户端认证的 CA 证书的目录。             |
| force-client-cert | true, false          | false | 指定是否强制客户端证书认证。                     |

### 示例

```bash
$ speedle-pms --store-type file --endpoint="127.0.0.1:6733" --insecure=false --key=$tls_config_path/server.key --cert=$tls_config_path/server.crt --client-cert=$tls_config_path/client-ca.crt

$ speedle-ads --store-type file --endpoint="127.0.0.1:6735" --insecure=false --force-client-cert=true --key=$tls_config_path/server.key --cert=$tls_config_path/server.crt --client-cert=$tls_config_path/client-ca.crt
```

**注意: 请在指定 TLS 相关目录时使用绝对目录。**

## 使用`spctl` CLI 访问启用 TLS 的 Speedle

### 命令行参数

Spctl CLI TLS 相关配置参数：

| 名称              | 值                   | 默认  | 注释                                                                 |
| ----------------- | -------------------- | ----- | -------------------------------------------------------------------- |
| skipverify        | true, false          | false | 指定是否跳过 Speedle 的 TLS 证书校验 certificate.                    |
| cert              | TLS certificate path |       | 指定客户端 TLS 证书的目录。                                          |
| key               | TLS private key path |       | 指定客户端 TLS 私钥的目录。                                          |
| client-cert       | client CA path       |       | 指定 CA 证书的目录。                                                 |
| force-client-cert | true, false          | false | specifies if the client certificate authentication is forced or not. |

### 示例

```bash
$ spctl config skipverify false cacert $tls_config_path/server-ca.crt cert $tls_config_path/client.crt key $tls_config_path/client.key pms-endpoint "https://localhost:6733/policy-mgmt/v1/"
$ spctl create service test
$ spctl get service --all
```

## 使用`curl`访问启用 TLS 的 Speedle

### 示例

```bash
$ curl --cacert $tls_config_path/server-ca.crt --cert $tls_config_path/client.crt --key $tls_config_path/client.key https://localhost:6733/policy-mgmt/v1/service
```
