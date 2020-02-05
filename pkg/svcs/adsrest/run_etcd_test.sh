#!/bin/bash

shell_dir=$(dirname $0)

set -ex
#source pkg/svcs/pmsrest/start_etcd.sh
rm -rf ./speedle.etcd
source ${shell_dir}/../../../setTestEnv.sh

#Reconfig spctl
${GOPATH}/bin/spctl config ads-endpoint http://localhost:6734/authz-check/v1/
${GOPATH}/bin/spctl config pms-endpoint http://localhost:6733/policy-mgmt/v1/

startPMS etcd --config-file ${shell_dir}/../pmsrest/config_etcd.json
startADS --config-file ${shell_dir}/../pmsrest/config_etcd.json

go clean -testcache github.com/oracle/speedle/pkg/svcs/adsrest
${GOPATH}/bin/spctl delete service --all
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsrest -tags=runtime_test_prepare
${GOPATH}/bin/spctl get service --all
sleep 2
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsrest -tags="runtime_test runtime_cache_test" -run=TestMats

rm -rf ./speedle.etcd
