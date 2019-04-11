<p align="center">
    <img src="/docs/images/Speedle_logo_b.svg" height="50%" width="50%" class="center"/>
</p>
<p align="center">
    <a href="https://join.slack.com/t/speedleproject/shared_invite/enQtNTUzODM3NDY0ODE2LTg0ODc0NzQ1MjVmM2NiODVmMThkMmVjNmMyODA0ZWJjZjQ3NDc2MjdlMzliN2U4MDRkZjhlYzYzMDEyZTgxMGQ">
        <img src="https://img.shields.io/badge/slack-speedle-red.svg">
    </a>
    <a href="https://github.com/oracle/speedle/tags">
        <img src="https://img.shields.io/github/tag/oracle/speedle.svg">
    </a>
    <a href="https://github.com/oracle/speedle/issues">
        <img src="https://img.shields.io/github/issues/oracle/speedle.svg">
    </a>
    <a href="https://goreportcard.com/report/github.com/oracle/speedle">
        <img src="https://goreportcard.com/badge/github.com/oracle/speedle">
    </a>
    <a href="https://app.wercker.com/project/byKey/07abf3ef318b376c1171c95346333083">
        <img alt="Wercker status" src="https://app.wercker.com/status/07abf3ef318b376c1171c95346333083/s/master">
    </a>
</p>

<p align="right">
<a href="README.zh-cn.md">中文版</a>
</p>

# Speedle

Speedle is a general purpose authorization engine. It allows users to construct their policy model with user-friendly policy definition language and get authorization decision in milliseconds based on the policies. Speedle is very user-friendly, efficient, and extremely scalable. 

Speedle open source project consits of a policy definition language, policy management module, authorization runtime module, commandline tool, and integration samples with popular systems.

## Documentation

Latest documentations are available at <https://speedle.io/docs>.

## Get Started

See Getting Started at <https://speedle.io/quick-start/>.

## Build

### Prerequisites

-   GO 1.10.1 or greater <https://golang.org/doc/install>
-   Set `GOROOT` and `GOPATH` properly

### Step

```
$ go get github.com/oracle/speedle/cmd/...
$ ls $GOPATH/bin
spctl  speedle-ads  speedle-pms
```

## Test

```
$ cd $GOPATH/src/github.com/oracle/speedle
$ make test
```

## Get Help

-   Join us on Slack: [#speedle-users](https://join.slack.com/t/speedleproject/shared_invite/enQtNTUzODM3NDY0ODE2LTg0ODc0NzQ1MjVmM2NiODVmMThkMmVjNmMyODA0ZWJjZjQ3NDc2MjdlMzliN2U4MDRkZjhlYzYzMDEyZTgxMGQ)
-   Mailing List: speedle-users@googlegroups.com

## Get Involved

-   Learn how to [contribute](CONTRIBUTING.md)
-   See [issues](https://github.com/oracle/speedle/issues) for issues you can help with
