# Docker authorization plugin

Docker users can implement a docker authorization plugin for authorization check, the docker official document "Access authorization plugin" shows the details.


# Build

Please make sure the speedle golang adsclient code(../../adsclient/go/src/speedle/) is in your $GOPATH/src package, if you are in the same directory of this readme file:

$cp -r ../../adsclient/go/src/speedle/ $GOPATH/src

And the docker authorization package is also needed:

$go get github.com/docker/go-plugins-helpers/authorization

Build:

$make

Executable speedle-docker-plugin could be found in $HOME/go/bin

# Run speedle
Please refer the [speedle quickstart](../../../docs/public/speedle/docs/quickstart.md) 

# Run docker and the plugin

1. First stop docker

$sudo systemctl stop docker

2. Plugin must be run before docker

$sudo $HOME/go/bin/speedle-docker-plugin <speedle host> <service name>

If the plugin is run behind a HTTPS proxy, please firstly set proxy.

3. Run docker engine with plugin

$sudo /usr/bin/dockerd --selinux-enabled --authorization-plugin=speedle-docker-plugin

eg: if speedle is running on localhost with service name=docker. $sudo $HOME/go/bin/speedle-docker-plugin localhost docker

# Testing

1. Test if containers could be listed. Expected: denied, because service "docker" is not created.

$docker ps
Error response from daemon: authorization denied by plugin speedle-docker-plugin:

2. Create a service and grant root to all resources.

$spctl create service docker
service created
{"name":"docker","type":"application"}
$spctl create policy root-policy -c "grant user root GET,POST,PUT,DELETE expr:.*" --service-name=docker
policy created
{"id":"e6e0ec73-1b7e-4816-8c7e-b4bc56561cf2","name":"root-policy","effect":"grant","permissions":[{"resourceExpression":".*","actions":["GET","POST","PUT","DELETE"]}],"principals":["user:root"]}

3. Test if containers could be listed.

$docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES

4. Deny user root to list containers.

$spctl create policy root-deny-policy -c "deny user root GET expr:.*" --service-name docker
policy created
{"id":"26edcfb5-22ee-493a-8fde-f757c615dd9c","name":"root-deny-policy","effect":"deny","permissions":[{"resource":"containers","actions":["GET"]}],"principals":["user:root"]
Test if container could be listed

$docker ps
Error response from daemon: authorization denied by plugin speedle-docker-plugin:
