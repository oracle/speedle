#!/bin/bash

shell_dir=$(dirname $0)

set -ex
source ${shell_dir}/../../../setTestEnv.sh


startPMS etcd --config-file ${shell_dir}/config_etcd.json

go clean -testcache github.com/oracle/speedle/pkg/svcs/pmsrest
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/pmsrest $*
