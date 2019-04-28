+++
date = "2017-04-10T16:41:54+01:00"
weight = 100
description = "Integrate with Docker Auth"
title = "Docker"
draft = false
bref= "This sample demonstrates how to protect your docker registry by Speedle authorization engine."
toc = true
tocheading = "h3"
tags = ["docker"]
categories = ["integration", "cloudnative"]
iconurl = "24. Docker.svg"
+++

## Docker authorization plugin

Docker users can implement a docker authorization plugin for authorization check, the docker official document "Access authorization plugin" shows the details.

### Build

```bash
git clone git@github.com:oracle/speedle.git

# make sure the speedle golang adsclient code is in your $GOPATH/src package.

cp -r samples/adsclient/go/src/speedle/ $GOPATH/src

# Get docker authorization plugin

go get github.com/docker/go-plugins-helpers/authorization

cd samples/integration/docker-integration
make
```

Executable `speedle-docker-plugin` could be found in \$HOME/go/bin

### Run speedle

refer to speedle [quickstart](../quick-start)

### Run docker and the plugin

- First stop docker

```bash
sudo systemctl stop docker
```

- Plugin must be run before docker. If the plugin is run behind a HTTPS proxy, set proxy first.

```bash
# sudo $HOME/go/bin/speedle-docker-plugin  <speedle host> <speedle name>
# e.g. if speedle is running on localhost with service name=docker.

sudo $HOME/go/bin/speedle-docker-plugin localhost docker
```

- Run docker engine with plugin

```bash
sudo /usr/bin/dockerd --selinux-enabled --authorization-plugin=speedle-docker-plugin
```

### Testing

- Test if containers could be listed. Expected: denied, because service "docker" is not created.

```bash
docker ps
Error response from daemon: authorization denied by plugin speedle-docker-plugin:
```

- Create a service and grant root to all resources.

```bash
spctl create service docker
service created
{"name":"docker","type":"application"}

spctl create policy root-policy -c "grant user root GET,POST,PUT,DELETE expr:.*" --service-name=docker
policy created
{"id":"e6e0ec73-1b7e-4816-8c7e-b4bc56561cf2","name":"root-policy","effect":"grant","permissions":[{"resourceExpression":".*","actions":["GET","POST","PUT","DELETE"]}],"principals":["user:root"]}
```

- Test if containers could be listed.

```bash
docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

- Deny user root to list containers.

```bash
spctl create policy root-deny-policy -c "deny user root GET expr:.*" --service-name docker
policy created
{"id":"26edcfb5-22ee-493a-8fde-f757c615dd9c","name":"root-deny-policy","effect":"deny","permissions":[{"resource":"containers","actions":["GET"]}],"principals":["user:root"]
```

- Test if container could be listed

```bash
docker ps
Error response from daemon: authorization denied by plugin speedle-docker-plugin:
```
