#/bin/bash

homeDir=$(cd $(dirname ${BASH_SOURCE[0]})/.. > /dev/null; pwd -P)

rm /go/bin
ln -s ${WERCKER_OUTPUT_DIR} /go/bin

#export GOPATH=$homeDir
#mkdir -p $GOPATH/src/github.com/oracle
#ln -s $homeDir $GOPATH/src/github.com/oracle/speedle
