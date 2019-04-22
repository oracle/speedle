#/bin/bash

homeDir=$(cd $(dirname ${BASH_SOURCE[0]})/.. > /dev/null; pwd -P)

if [ ! -e ${homeDir}/bin ]; then
    mkdir ${homeDir}/bin
fi
rmdir /go/bin
ln -s ${homeDir}/bin /go/bin

export GOPATH=/go
#mkdir -p $GOPATH/src/github.com/oracle
#ln -s $homeDir $GOPATH/src/github.com/oracle/speedle
