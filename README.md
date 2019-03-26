
<img src="/docs/images/Speedle_logo_b.svg" height="50%" width="50%" class="center"/> 

# Speedle

Speedle is an open source project for authorization management. It consits of a policy definition language, policy management module, authorization runtime module, commandline tool, and integration samples with popular systems.   


## Documentation

Latest documentation and javadocs are available at <https://speedle.io/docs/usecases/>.

## Get Started

See Getting Started at <https://speedle.io/quick-start/>.


## Build

### Prerequisites

* GO 1.10.1 or greater <https://golang.org/doc/install>
* Set `GOROOT` and `GOPATH` properly

### Steps

1. Fetch Speedle source code
```
[opc@wcai-speedle-host ~]$ go get github.com/oracle/speedle
```

2. Move to Speedle source code directory

```
[opc@wcai-speedle-host gopath]$ cd $GOPATH/src/github.com/oracle/speedle
```

3. Compile

```
[opc@wcai-speedle-host speedle]$ make
go build -ldflags "-X main.gitCommit=a61b32e -X main.productVersion=18.4.1 -X main.goVersion=go1.12.1" -o /home/opc/gopath/bin/speedle-pms github.com/oracle/speedle/cmd/speedle-pms
go build -ldflags "-X main.gitCommit=a61b32e -X main.productVersion=18.4.1 -X main.goVersion=go1.12.1" -o /home/opc/gopath/bin/speedle-ads github.com/oracle/speedle/cmd/speedle-ads
go build -ldflags "-X main.gitCommit=a61b32e -X main.productVersion=18.4.1 -X main.goVersion=go1.12.1" -o /home/opc/gopath/bin/spctl  github.com/oracle/speedle/cmd/spctl
```

## Get Help

* Join us on Slack: [#speedle-users](https://join.slack.com/t/speedleproject/shared_invite/enQtNTUzODM3NDY0ODE2LTg0ODc0NzQ1MjVmM2NiODVmMThkMmVjNmMyODA0ZWJjZjQ3NDc2MjdlMzliN2U4MDRkZjhlYzYzMDEyZTgxMGQ)
* Mailing List: speedle-users@googlegroups.com

## Get Involved

* Learn how to [contribute](CONTRIBUTING.md)
* See [issues](https://github.com/oracle/speedle/issues) for issues you can help with

