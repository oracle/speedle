#!/bin/bash

shell_dir=$(dirname $0)

set -ex
source ${shell_dir}/../../../setTestEnv.sh

startPMS file --config-file ${shell_dir}/config_file.json

go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/pmsrest $*
