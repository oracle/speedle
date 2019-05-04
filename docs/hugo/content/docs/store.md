+++
title = "Pluggable Storage"
description = "Bring your own persistence store!"
weight = 130
draft = false
toc = true
tocheading = "h2"
tags = ["storage", "pluggable"]
categories = ["docs"]
bref = "By default, Speedle supports etcd and file as persistence stores. However, It is easy to plug in a different persistence store, say, mongodb. This document describes how one could go about implementing/using a different store."
+++


## Overview
Speedle now supports two kinds of data store, OOTB: file store and etcd store.
However, you can implement your own data store (e.g. with mongodb, etc)

* Please note the data store needs to support the `watch` function.

This document walks through step-by-step instructions to implement a data store.

## Write store code to implement the PolicyStoreManager interface

Create a "mystore" directory under store directory and navigate to it.
Create a store code file like `mystore.go`, and implment the `PolicyStoreManager` interface in this file.

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

Please pay attention to the `Watch` function. This function will monitors the changes of your data store. This function needs to return a `StorageChangeChannel` and every store change event (please check `api/pms/types/StoreChangeEvent` for details) will be send to this channel. Authorization Decision Service (ADS) will receives these change events and updates its cache immediately.

## Write storeBuilder code

### Understand the store configuration in speedle
Speedle reads store config info from three different sources: flags, environment variables and config file.

The flags have highest priority, then the environment variable, then the config file. That means the config items in higher priority source will override the same config items in lower priority source.

The store owner needs to provide the flags definition in init function, we use [pflag](https://github.com/spf13/pflag) to define the flags. 

Flags definition in `store/etcd/storeBuilder.go`:
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

The environment variable name is the transformation of flag name, the rule is: add a `SPDL_` prefix, replace every `-` to `_`, and convert all the letters to upper case. For example:
`etcdstore-endpoint` -> `SPDL_ETCDSTORE_ENDPOINT`

The store owner also needs to define the store property name used in config file, through provides a {`flagName:storePropName`} map. This map will shows the correspondence between flag and store property. 

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
    },
}
```

### Implement the StoreBuilder interface in storeBuild code. 
StoreBuild will provides the function about create store and the function about get this store realted parameters.

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

`NewStore` method needs to return a store instance implementing `PolicyStoreManager` interface. 
The input parameter is a config map correspond to the `storeProps` in config file, which merged with the configs from flags and environment variables. 
You can read the config item's value from this config map to build the store. If you want to read the `EtcdEndpoint` value, you can do like this: 
```
etcdEndpoint, ok := config[EtcdEndpointKey].(string)
```

Please note that the config value may comes from different sources(flag, env variable, config file), the config value type maybe is your expect type(like int or bool) or string type. So you need to check the value type and may need to convert the value from string type to your expect type. 

The `GetStoreParams()` interface needs to return a {`flagName:storePropName`} map. This map will shows the correspondence between flag and store property.

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


### Register the storeBuilder
Register the storeBuilder will makes a type of store available by the provided name. 
And the init function needs to declare all the flags this store needed. 

Example in `store/etcd/storeBuilder.go`:
```golang
func init() {
    pflag.String(EtcdEndpointFlagName, DefaultEtcdStoreEndpoint, "Store config: endpoint of etcd store.")
    pflag.Bool(IsEmbeddedEtcdFlagName, DefaultEtcdStoreIsEmbedded, "Store config: is embedded etcd store or not.")
    ...

    store.Register(StoreType, Etcd3StoreBuilder{})
}
```

## Link the new store to Speedle
In cmd/speedle-ads folder and cmd/speedle-pms folder, you can find a stores.go file with below content:

```golang
package main

import (
    _ "github.com/oracle/speedle/store/etcd"
    _ "github.com/oracle/speedle/store/file"
)
```

In this file, we link every store implmention with a side-effect import (using a blank import name). You can add your own store here too.
If you want to use Speedle in `in-process` mode, you can copy this `stores.go` to your own package and modify the package name to your own package name.
