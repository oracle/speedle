#/bin/bash

homeDir=$(cd $(dirname ${BASH_SOURCE[0]})/.. > /dev/null; pwd -P)

if [ ! -e ${homeDir}/bin ]; then
    mkdir ${homeDir}/bin
fi
rmdir /go/bin
ln -s ${homeDir}/bin /go/bin

if [ ! -e $WERCKER_CACHE_DIR/pkg ]; then
    mkdir $WERCKER_CACHE_DIR/pkg
fi
if [ -e /go/pkg ]; then
    rmdir /go/pkg
fi
ln -s $WERCKER_CACHE_DIR/pkg /go/pkg

export GOPATH=/go
