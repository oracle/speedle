#!/bin/bash

shell_dir=$(dirname $0)
temp_policy_file=/tmp/speedle-test-file-store.json

set -ex
source ${GOPATH}/src/github.com/oracle/speedle/setTestEnv.sh

startPMS file --config-file ${shell_dir}/../pmsrest/config_file.json
startADS --config-file ${shell_dir}/../pmsrest/config_file.json

go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsgrpc -run=TestMats
