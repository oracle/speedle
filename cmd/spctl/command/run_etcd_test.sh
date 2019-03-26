#!/bin/bash
set -ex
#source pkg/svcs/pmsrest/start_etcd.sh
rm -rf ./speedle.etcd
source ${GOPATH}/src/github.com/oracle/speedle/setTestEnv.sh
go clean -testcache

#Reconfig spctl
${GOPATH}/bin/spctl config pms-endpoint http://localhost:6733/policy-mgmt/v1/

startPMS etcd --config-file pkg/svcs/pmsrest/config_etcd.json
go test ${TEST_OPTS} github.com/oracle/speedle/cmd/spctl/command -run=TestMats

