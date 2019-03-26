#!/bin/bash

shell_dir=$(dirname $0)
temp_policy_file=/tmp/speedle-test-file-store.json

set -ex
source ${GOPATH}/src/gitlab-odx.oracledx.com/wcai/speedle/setTestEnv.sh

startPMS file --config-file ${shell_dir}/../pmsrest/config_file.json

go test ${TEST_OPTS} gitlab-odx.oracledx.com/wcai/speedle/pkg/svcs/pmsgrpc -run=TestMats
