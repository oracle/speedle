#!/bin/bash

shell_dir=$(dirname $0)

set -ex
source ${GOPATH}/src/gitlab-odx.oracledx.com/wcai/speedle/setTestEnv.sh

startPMS file --config-file ${shell_dir}/config_file.json

go test ${TEST_OPTS} gitlab-odx.oracledx.com/wcai/speedle/pkg/svcs/pmsrest $*
