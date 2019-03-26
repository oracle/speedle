#/bin/sh
set -ex
source ${GOPATH}/src/github.com/oracle/speedle/setTestEnv.sh

#Reconfig spctl
${GOPATH}/bin/spctl config ads-endpoint http://localhost:6734/authz-check/v1/
${GOPATH}/bin/spctl config pms-endpoint http://localhost:6733/policy-mgmt/v1/

ads --config-file ../pmsrest/config_file.json &
serverPID=$!
add_exit_trap "kill ${serverPID}"
waitService ADS 6734 || exit 1

#go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsrest -v -tags="runtime_test runtime_cache_test" -run=TestMats
go test ${TEST_OPTS} github.com/oracle/speedle/pkg/svcs/adsrest -v -tags="runtime_cache_test" -run=TestMats

