#/bin/bash

homeDir=$(cd $(dirname ${BASH_SOURCE[0]})/.. > /dev/null; pwd -P)

echo "**** ${PWD}"
echo "**** ${homeDir}"

rmdir /go/bin
mkdir ${WERCKER_SOURCE_DIR}/bin
ln -s ${WERCKER_SOURCE_DIR}/bin /go/bin

export GOPATH=/go
#mkdir -p $GOPATH/src/github.com/oracle
#ln -s $homeDir $GOPATH/src/github.com/oracle/speedle
