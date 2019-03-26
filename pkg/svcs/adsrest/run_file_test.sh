#!/bin/bash

shell_dir=$(dirname $0)

set -ex
source ${GOPATH}/src/github.com/oracle/speedle/setTestEnv.sh

#Reconfig spctl
${GOPATH}/bin/spctl config ads-endpoint http://localhost:6734/authz-check/v1/
${GOPATH}/bin/spctl config pms-endpoint http://localhost:6733/policy-mgmt/v1/


startPMS file --config-file ${shell_dir}/../pmsrest/config_file.json
${GOPATH}/bin/spctl delete service --all
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsrest -tags=runtime_test_prepare
$GOPATH/bin/spctl get service --all

startADS --config-file ${shell_dir}/../pmsrest/config_file.json
sleep 2
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsrest -tags="runtime_test runtime_cache_test" -run=TestMats
