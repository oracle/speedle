+++
title = "存储扩展"
description = "Bring your own persistence store!"
weight = 130
draft = false
toc = true
tocheading = "h2"
tags = ["storage", "pluggable"]
categories = ["docs"]
bref = ""
+++

## 概述

Speedle 现在支持两种数据存储，OOTB:文件存储和 etcd 存储。
但是，您可以实现自己的数据存储(例如使用 mongodb 等)

- 请注意数据存储需要支持“watch”功能

本文档将说明如何逐步实现一个数据存储。

## 实现 PolicyStoreManager 接口

在 store 目录下创建一个“mystore”目录并导航到它目录下面。
创建一个类似 mystore 之类的存储代码文件，并在此文件中实现“PolicyStoreManager”接口。

Example in `store/etcd/etcdStore.go`:

```golang
type Store struct {
    ...
}
func (s *Store) ReadPolicyStore() (*pms.PolicyStore, error) {
  ...
}
func (s *Store) CreateService(service *pms.Service) error {
  ...
}
...
```

请注意这个 watch 的功能。该功能将监视数据存储的更改。这个函数需要返回一个“StorageChangeChannel”对象，每个 store change event(请查看“api/pms/types/StoreChangeEvent”以获取详细信息)都会被发送到这个 Channel。ADS 将接收这些更改事件并立即更新其缓存。

## 编写 storeBuilder 代码

### 理解 Speedle 存储配置

Speedle 从三个不同的源读取存储配置信息:flags、环境变量和配置文件。

flags 具有最高优先级，然后是环境变量，然后是配置文件。这意味着高优先级源中的配置项将覆盖低优先级源中的相同配置项。

存储所有者需要在 init 函数中提供 flags 定义，我们使用[pflag](https://github.com/spf13/pflag)来定义 flags。

`store/etcd/storeBuilder.go` 中的 flags 定义:

```golang
const (
    IsEmbeddedEtcdFlagName             = "etcdstore-isembedded"
    EmbeddedEtcdDataDirFlagName        = "etcdstore-embeddedDataDir"
    EtcdEndpointFlagName               = "etcdstore-endpoint"
    EtcdKeyPrefixFlagName              = "etcdstore-keyprefix"
    EtcdTLSClientCertFileFlagName      = "etcdstore-tls-cert"
    EtcdTLSClientKeyFileFlagName       = "etcdstore-tls-key"
    EtcdTLSClientTrustedCAFileFlagName = "etcdstore-tls-ca"
    EtcdTLSAllowedCNFlagName           = "etcdstore-tls-allowedCN"
    EtcdTLSServerNameFlagName          = "etcdstore-tls-serverName"
    EtcdTLSCRLFileFlagName             = "etcdstore-tls-crlFile"
    EtcdTLSInsecureSkipVerifyFlagName  = "etcdstore-tls-insecureSkipVerify"

    //default property values
    DefaultKeyPrefix           = "/speedle_ps/"
    DefaultEtcdStoreEndpoint   = "localhost:2379"
    DefaultEtcdStoreKeyPrefix  = "/speedle_ps/"
    DefaultEtcdStoreIsEmbedded = false
)

func init() {
    pflag.String(EtcdEndpointFlagName, DefaultEtcdStoreEndpoint, "Store config: endpoint of etcd store.")
    pflag.String(EtcdKeyPrefixFlagName, DefaultEtcdStoreKeyPrefix, "Store config: key prefix to store speedle policy data in etcd store.")
    pflag.Bool(IsEmbeddedEtcdFlagName, DefaultEtcdStoreIsEmbedded, "Store config: is embedded etcd store or not.")
    pflag.String(EmbeddedEtcdDataDirFlagName, "", "Store config: data dir for embedded etcd store.")
    pflag.String(EtcdTLSClientCertFileFlagName, "", "Store config: etcd x509 client cert.")
    pflag.String(EtcdTLSClientKeyFileFlagName, "", "Store config: etcd x509 client key.")
    pflag.String(EtcdTLSClientTrustedCAFileFlagName, "", "Store config: etcd x509 client CA cert.")
    pflag.String(EtcdTLSAllowedCNFlagName, "", "Store config: etcd x509 allowed CN.")
    pflag.String(EtcdTLSServerNameFlagName, "", "Store config: etcd x509 server name.")
    pflag.String(EtcdTLSCRLFileFlagName, "", "Store config: etcd x509 CRL file.")
    pflag.Bool(EtcdTLSInsecureSkipVerifyFlagName, false, "Store config: etcd x509 insecure skip verify.")
}
```

环境变量名称是 flag 名称的转换，规则是:添加一个“SPDL*”前缀，将每个“-”替换为“*”，并将所有字母转换为大写。例如:
`etcdstore-endpoint` -> `SPDL_ETCDSTORE_ENDPOINT`

存储所有者还需要定义配置文件中使用的存储属性名，通过提供{' flagName:storePropName '} map。这个映射将显示 flag 和 store 属性之间的对应关系。

Config file example:

```json
{
	"storeConfig": {
		"storeType": "etcd",
		"storeProps": {
			"EtcdEndpoint": "localhost:2379",
			"EtcdKeyPrefix": "/opss_ps/",
			"IsEmbeddedEtcd": true,
			"EmbeddedEtcdDataDir": "./speedle.etcd"
		}
	}
}
```

### 实现 StoreBuilder 接口

StoreBuild 将提供关于创建存储的函数，以及获取该存储相关参数的函数。

Example in `store/etcd/storeBuilder.go`:

```golang
type Etcd3StoreBuilder struct{}

func (esb Etcd3StoreBuilder) NewStore(config map[string]interface{}) (pms.PolicyStoreManager, error) {
   ...
}
func (fs FileStoreBuilder) GetStoreParams() map[string]string {
   ...
}
```

“NewStore”方法需要返回实现“PolicyStoreManager”接口的存储实例。
输入参数是一个配置 map，对应于配置文件中的“storeProps”，它将于与来自 flags 和环境变量的配置信息合并在一起。
您可以从这个配置 map 中读取配置项的值来构建存储。如果您想读取' EtcdEndpoint '值，可以这样做:

```
etcdEndpoint, ok := config[EtcdEndpointKey].(string)
```

请注意，配置值可能来自不同的源(flag、环境变量、配置文件)，配置值类型可能是您的期望类型(如 int 或 bool)或字符串类型。因此，您需要检查值类型，并可能需要将值从 string 类型转换为您期望的类型。

接口“GetStoreParams()”需要返回一个{“flagName:storePropName”}map。这个 map 将显示 flag 和 store 属性之间的对应关系。

etcd storeBuilder `GetStoreParams`() function example:

```golang
func (esb Etcd3StoreBuilder) GetStoreParams() map[string]string {
    return map[string]string{

        IsEmbeddedEtcdFlagName:             IsEmbeddedEtcdKey,
        EmbeddedEtcdDataDirFlagName:        EmbeddedEtcdDataDirKey,
        EtcdEndpointFlagName:               EtcdEndpointKey,
        EtcdKeyPrefixFlagName:              EtcdKeyPrefixKey,
        EtcdTLSClientCertFileFlagName:      EtcdTLSClientCertFileKey,
        EtcdTLSClientKeyFileFlagName:       EtcdTLSClientKeyFileKey,
        EtcdTLSClientTrustedCAFileFlagName: EtcdTLSClientTrustedCAFileKey,
        EtcdTLSAllowedCNFlagName:           EtcdTLSAllowedCNKey,
        EtcdTLSServerNameFlagName:          EtcdTLSServerNameKey,
        EtcdTLSCRLFileFlagName:             EtcdTLSCRLFileKey,
        EtcdTLSInsecureSkipVerifyFlagName:  EtcdTLSInsecureSkipVerifyKey,
    }

}
```

### 注册 storeBuilder

注册 storeBuilder 将根据提供的名称使这种类型的存储可用。
init 函数需要声明这个存储所需的所有 flags。

Example in `store/etcd/storeBuilder.go`:

```golang
func init() {
    pflag.String(EtcdEndpointFlagName, DefaultEtcdStoreEndpoint, "Store config: endpoint of etcd store.")
    pflag.Bool(IsEmbeddedEtcdFlagName, DefaultEtcdStoreIsEmbedded, "Store config: is embedded etcd store or not.")
    ...

    store.Register(StoreType, Etcd3StoreBuilder{})
}
```

## 链接新存储到 Speedle

在目录 cmd/speedle-ads 和目录 cmd/speedle-pms 中，你可以找到一个 store.go 文件，其具有下面的内容：

```golang
package main

import (
    _ "github.com/oracle/speedle/store/etcd"
    _ "github.com/oracle/speedle/store/file"
)
```

在这个文件中，我们用一个 side-effect 导入(使用一个空白导入名称)来链接每个存储的实现。您也可以在这里添加自己的存储。
如果你想在“in-process”模式下使用 Speedle，你可以复制这个 store.go 文件。转到您自己的 package 中，将 package 名字修改为您自己的 package 名字。
