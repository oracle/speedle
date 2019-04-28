+++
title = "日志"
description = "Speedle logging"
weight = 320
draft = false
toc = true
tocheading = "h2"
tocsidebar = false
tags = ["v0.1.0", "logging"]
categories = ["docs"]
bref = ""
+++

## 概述

Speedle 使用了下面两个开源项目来处理日志,

```
1. https://github.com/sirupsen/logrus
2. https://github.com/natefinch/lumberjack
```

Logrus 是一个针对 go 语言的结构化的日志处理器，它与 go 语言标准的 logger 库完全 API 兼容。因为 logrus 不支持对日志循环滚动处理（rotation），所以我们引入了另外一个开源项目 lumberjack。Logrus 与 lumberjack 之间通过 io.Writer 接口交互，如下图所示,

![Speedle Logging](../public/speedle/docs/img/logger.jpg)

## 日志配置

可以通过两种方式对 Speedle 的日志模块进行配置，即命令行参数或者配置文件。

### 通过命令行参数配置日志

所有的命令行参数如下所示。前两个参数供 logrus 使用，其它的参数则是为 lumberjack 准备的。

```
--log-level string
--log-formatter string
--log-reportcaller bool
--log-filename string
--log-maxsize int
--log-maxbackups int
--log-maxage int
--log-localtime bool
--log-compress bool
```

上面每一个参数的详细信息包含在下面的表格中。所有的参数都收可选的。 如果所有参数都没有配置，那么日志消息默认会输出到标准错误输出（stderr），而且默认的日志级别是 info。

| 配置项             | 描述                                                                                                 | 默认值                                                                                  |
| ------------------ | ---------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| --log-level        | 日志级别, 有效的日志级别包括：panic、fatal、error、warn、info、debug                                 | info                                                                                    |
| --log-formatter    | 日志格式, 有效的日志格式包括：text、json                                                             | text                                                                                    |
| --log-reportcaller | 用于控制是否将调用信息 (文件名, 行数和函数名)包含在日志消息中                                        | false                                                                                   |
| --log-filename     | 日志文件名                                                                                           | 没有默认值。 如果没有配置，那么日志消息默认会输出到标准错误输出（stderr）               |
| --log-maxsize      | 每个日志文件的最大尺寸。 注意该参数只有在"--log-filename"已经配置的情况下才有效                      | 100M                                                                                    |
| --log-maxbackups   | 历史日志文件保留的最大数量。 注意该参数只有在"--log-filename"已经配置的情况下才有效                  | 0, 保留所有历史日志文件（注意历史日志文件可能由于达到了 maxAge 设定的时间限制而被删除） |
| --log-maxage       | 历史日志文件保留的最多天数。注意该参数只有在"--log-filename"已经配置的情况下才有效                   | 0, 不删除任何历史日志文件（注意历史日志文件可能由于达到了 maxbackups 的限制而被删除）   |
| --log-localtime    | 用于控制在生成滚动日志文件时是否使用本地时间。注意该参数只有在"--log-filename"已经配置的情况下才有效 | false, 使用 UTC 时间                                                                    |
| --log-compress     | 用于控制在生成滚动日志文件时是否需要压缩。注意该参数只有在"--log-filename"已经配置的情况下才有效     | false, 不压缩                                                                           |

下面是两个日志消息的示例，第一个使用 text 格式输出，第二个则是使用 json 格式输出，

```
# text
time="2015-03-26T01:27:38-04:00" level=debug msg="Started observing beach" animal=walrus number=8
```

```
# json
{"animal":"walrus","level":"info","msg":"A group of walrus emerges from the
ocean","size":10,"time":"2014-03-10 19:57:38.562264131 -0400 EDT"}
```

### 通过配置文件配置日志

下面就是一个配置文件的例子，每一项的含义与对应的命令行参数完全一样。同样，所有配置项都是可选的。

```go
{
    "logConfig": {
        "level": "info",
        "formatter": "text",
        "setReportCaller": "false",
        "rotationConfig": {
            "filename": "/mnt/logs/speedle.log",
            "maxSize": 20,
            "maxBackups": 5,
            "maxAge": 0,
            "LocalTime": false,
            "compress": false
        }
    }
}
```

## 配置日志的最佳实践

如果 Speedle 以 docker 容器的方式运行，那么推荐使用 docker 的日志驱动 json-file 来配置日志文件的大小以及滚动等，而不是使用 natefinch/lumberjack。这种情况下，仅仅需要用 sirupsen/logrus 来配置日志的级别以及格式。

Docker 的所有配置信息，包括日志配置，都包含在配置文件/etc/docker/daemon.json 中。关于如何配置 docker 的日志驱动 json-file，具体请参考[docker json-file](https://docs.docker.com/config/containers/logging/json-file/)。下面是 json-file 的一个配置示例,

```
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "20m",
    "max-file": "10"
  }
}
```

其实在大多数情况下，用户也不需要配置 sirupsen/logrus，因为默认的日志级别（info）和默认格式的日志格式（text）基本满足大多数的需求。通常在用户希望改变日志格式的时候才需要配置 logrus，例如改成 json 格式。

## 如何增加新的日志消息

首先确保在\*.go 源文件中导入了 logrus 的 package，

```go
import (
  log "github.com/sirupsen/logrus"
)
```

然后就可以增加任何你希望的日志消息。例如：

```go
  log.Debug("This is a debug log entry")
  log.Info("This is a info log entry")
  log.Warn("This is a warning log entry")
  log.Warning("This is a warning log entry again")
  log.Error("This is a error log entry")
```

具体请参考[speedle/logging/demo/demo.go](https://github.com/oracle/speedle/blob/master/pkg/logging/demo/demo.go).

## 与第三方日志系统集成

Speedle 的日志可以很容易地与第三方日志系统集成，例如 fluentd、filebeat。用户所需要做的就是配置 fluentd 或 filebeat 去监控 Speedle 生成的日志文件。下面的例子就是配置 fluentd 的 in_tail 插件来收集 docker 容器的日志，如果 Speedle 作为 docker 容器运行的话，那么产生的日志自然也会被收集。

```
    <source>
      @type tail
      path /var/lib/docker/containers/*/*.log
      pos_file /data/fluentd/fluentd.log.pos
      tag speedle.log
      <parse>
        @type json
        time_key time
        keep_time_key true
      </parse>
      refresh_interval 5
    </source>
```

请参考[fluentd](https://www.fluentd.org/)和[filebeat](https://www.elastic.co/products/beats/filebeat)来获得更多信息。
