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
rmdir /go/pkg
ln -s $WERCKER_CACHE_DIR/pkg /go/pkg

export GOPATH=/go
#mkdir -p $GOPATH/src/github.com/oracle
#ln -s $homeDir $GOPATH/src/github.com/oracle/speedle
