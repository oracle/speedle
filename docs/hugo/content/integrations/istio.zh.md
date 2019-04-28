+++
title = "Istio"
description = "Integrate with Istio Speedle adapter for policy decisions in the service mesh"
weight = 300
draft = false
iconurl = "23. Istio.svg"
bref= "This sample demonstrates how to connect Istio with Speedle via Istio Mixer Adapter"
toc = true
tocheading = "h3"
tags = ["istio"]
categories = ["integration", "cloudnative"]
+++

## Speedle Istio Mixer Adapter for authorization

### Prerequisite

-   To setup speedle istio adapter, a running Kubernetes cluster with istio is required.
-   Speedle services need to be deployed. Please see [Deploy Speedle](../../docs/deployment)
-   [Istio bookinfo sample](https://istio.io/docs/examples/bookinfo/) installed.
-   The kubectl configuration is set so that kubectl can be used to manage the cluster.
-   A docker repository is required to push build result.
-   [Istio Adapter Before you start](https://github.com/istio/istio/wiki/Mixer-Out-Of-Process-Adapter-Walkthrough#before-you-start)

### Build

```bash
git clone git@github.com:oracle/speedle.git
cd samples/integration/istio-integration
```

copy set-env.sh.template to set-env.sh and edit it according to your environment.

```bash
. set-env.sh
make init init_istio
make build-grpc-adapter
```

### Install Speedle adapter

The adapter/speedlegrpcadapter/operator_cfg.yaml.template file is configured to protect [Istio bookinfo](https://github.com/istio/istio/tree/master/samples/bookinfo) services. It can be edited to "match" attribute according to services to be protected.

```yaml
match: destination.labels["app"] == "details" || destination.labels["app"] == "productpage" || destination.labels["app"] == "reviews" || destination.labels["app"] == "ratings"
```

To install Speedle Istio adapter:

```bash
make install-speedle-grpc-adapter
```

### Install Speedle adapter in discover mode

By default, Speedle Istio adapter runs in normal authorization check mode. The Speedle Istio adapter can run in "discover" mode, in which all authorization requests will be allowed. The authorization requests will be collected at Speedle ADS. These requests can be retrieved. The application developers can use the collected requests to define policies.

```bash
export SPEEDLE_ADS_ENDPOINT="http://<speedle host>:6734/authz-check/v1/discover"
# Or you can edit set-env.sh to use 'discover' instead of 'is-allowed' in SPEEDLE_ADS_ENDPOINT

make install-speedle-grpc-adapter
```

After Speedle Istio adapter is installed, you can run some tests against your application. Then you can use spctl command line tool to get collected requests and policies (for reference only):

```bash
# get discovered requests
spctl discover request --service-name=istio

# get discovered policies (for refernece only)
spctl discover policy --service-name=istio
```

### Uninstall Speedle Adapter

To uninstall Speedle Istio adapter:

```bash
make uninstall-speedle-grpc-adapter
```

### References

-   [Istio Bookinfo sample](https://github.com/istio/istio/tree/master/samples/bookinfo)
-   [Istio policy attribute vocabulary](https://istio.io/docs/reference/config/policy-and-telemetry/attribute-vocabulary)
-   [Mixer Out Of Process Adapter Dev Guide](https://github.com/istio/istio/wiki/Mixer-Out-Of-Process-Adapter-Dev-Guide)
