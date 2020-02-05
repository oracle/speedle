#!/bin/bash

shell_dir=$(dirname $0)
rm -rf ./speedle.etcd

set -ex
source ${shell_dir}/../../../setTestEnv.sh


startPMS etcd --config-file ${shell_dir}/../pmsrest/config_etcd.json
startADS --config-file ${shell_dir}/../pmsrest/config_etcd.json

go clean -testcache github.com/oracle/speedle/pkg/svcs/adsgrpc
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsgrpc -run=TestMats
rm -rf ./speedle.etcd
