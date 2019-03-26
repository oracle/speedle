# Speedle - data store

## Overview
Speedle now support two kinds of data store: file store and etcd store.
You can implement your own data store (e.g. with mongodb, etc)

* Please notice the data store need support the watch function.

This document walks through step-by-step instructions to implement a data store.

## step 1: Write store code to implement the PolicyStoreManager interface

Create a "mystore" directory under store directory and navigate to it.
Create a store code file like mystore.go, and implment the PolicyStoreManager interface in this file.

Example in store/etcd/etcdStore.go:
```
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

Please pay attention about the "Watch" function, this function will monitor the changes of your data store. This function need return a StorageChangeChannel and every store change event(please check api/pms/types/StoreChangeEvent for details) will be send to this channel. ADS will receive these change events and update its cache immediately.

## step 2: Write storeBuilder code

### Understand the store configuration in speedle
Speedle read store config info from three different sources: flags, environment variables and config file.

The flags have highest priority, then the environment variable, then the config file. That mean the config item in higher priority source will override the same config item in lower priority source.

The store owner need provide the flags defination in init function, we use [pflag](https://github.com/spf13/pflag) to define the flags. 

Flags defination in store/etcd/storeBuilder.go:
```
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

The environment variable name is the transformation of flag name, the rule is: add a "SPDL_" prefix, replace every "-" to "_", and convert all the letters to upper case. For example:
"etcdstore-endpoint" -> "SPDL_ETCDSTORE_ENDPOINT"

The store owner also need define the store property name used in config file, through provide a {flagName:storePropName} map. This map will show the correspondence between flag and store property. 

Config file example:
```
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
StoreBuild will provide the function about create store and funcation get this store realted parameters.

Example in store/etcd/storeBuilder.go:
```
type Etcd3StoreBuilder struct{}

func (esb Etcd3StoreBuilder) NewStore(config map[string]interface{}) (pms.PolicyStoreManager, error) {
   ...
}
func (fs FileStoreBuilder) GetStoreParams() map[string]string {
   ...
}
```

NowStore method need return a store instance implemented PolicyStoreManager interface. 
The input parameter is a config map correspond to the "storeProps" in config file, which merged with the configs from flags and environment variables. 
You can read the config item's value from this config map to build the store. If you want to read the "EtcdEndpoint" value, you can do like this: 
```
etcdEndpoint, ok := config[EtcdEndpointKey].(string)
```

Please notice because the config value may come from different sources(flag, env variable, config file), the config value type maybe is your expect type(like int or bool) or string type. So you need check the value type and may need convert the value from string type to your expect type. 

The GetStoreParams() interface need return a {flagName:storePropName} map. This map will show the correspondence between flag and store property.

etcd storeBuilder GetStoreParams() function example:
```
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
And in init function need declare all the flags this store needed. 

Example in store/etcd/storeBuilder.go:
```
func init() {
    pflag.String(EtcdEndpointFlagName, DefaultEtcdStoreEndpoint, "Store config: endpoint of etcd store.")
    pflag.Bool(IsEmbeddedEtcdFlagName, DefaultEtcdStoreIsEmbedded, "Store config: is embedded etcd store or not.")
    ...

    store.Register(StoreType, Etcd3StoreBuilder{})
}
```

## step 3: Link the new store to speedle
In cmd/speedle-ads folder and cmd/speedle-pms folder, you can find a stores.go file with below content:

```
package main

import (
    _ "github.com/oracle/speedle/pkg/store/etcd"
    _ "github.com/oracle/speedle/pkg/store/file"
)
```

In this file, we link every store implmention with a side-effect import (using a blank import name). You can add your own store here too.
If you want to use speedle as in-process mode, you can copy this stores.go to your own package and modify the package name to your own package name.
