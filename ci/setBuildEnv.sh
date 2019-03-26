#/bin/bash

homeDir=$(cd $(dirname ${BASH_SOURCE[0]})/.. > /dev/null; pwd -P)

export GOPATH=$homeDir
mkdir -p $GOPATH/src/github.com/oracle
ln -s $homeDir $GOPATH/src/github.com/oracle/speedle
