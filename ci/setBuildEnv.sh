#/bin/bash

homeDir=$(cd $(dirname ${BASH_SOURCE[0]})/.. > /dev/null; pwd -P)

export GOPATH=$homeDir
mkdir -p $GOPATH/src/gitlab-odx.oracledx.com/wcai
ln -s $homeDir $GOPATH/src/gitlab-odx.oracledx.com/wcai/speedle
