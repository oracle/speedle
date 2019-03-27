<p align="center">
    <img src="/docs/images/Speedle_logo_b.svg" height="50%" width="50%" class="center"/>
</p>
<p align="center">
    <a href="https://app.wercker.com/project/byKey/07abf3ef318b376c1171c95346333083">
    <img alt="Wercker status" src="https://app.wercker.com/status/07abf3ef318b376c1171c95346333083/s/master">
    </a>
</p>

# Speedle

Speedle is an open source project for authorization management. It consits of a policy definition language, policy management module, authorization runtime module, commandline tool, and integration samples with popular systems.

## Documentation

Latest documentation and javadocs are available at <https://speedle.io/docs/usecases/>.

## Get Started

See Getting Started at <https://speedle.io/quick-start/>.

## Build

### Prerequisites

-   GO 1.10.1 or greater <https://golang.org/doc/install>
-   Set `GOROOT` and `GOPATH` properly

### Step

```
[opc@wcai-speedle-host gopath]$ go get github.com/oracle/speedle/cmd/...
[opc@wcai-speedle-host gopath]$ ls $GOPATH/bin
spctl  speedle-ads  speedle-pms
```

## Test

```
[opc@wcai-speedle-host ~]$ cd $GOPATH/src/github.com/oracle/speedle
[opc@wcai-speedle-host speedle]$ make test
```

## Get Help

-   Join us on Slack: [#speedle-users](https://join.slack.com/t/speedleproject/shared_invite/enQtNTUzODM3NDY0ODE2LTg0ODc0NzQ1MjVmM2NiODVmMThkMmVjNmMyODA0ZWJjZjQ3NDc2MjdlMzliN2U4MDRkZjhlYzYzMDEyZTgxMGQ)
-   Mailing List: speedle-users@googlegroups.com

## Get Involved

-   Learn how to [contribute](CONTRIBUTING.md)
-   See [issues](https://github.com/oracle/speedle/issues) for issues you can help with
