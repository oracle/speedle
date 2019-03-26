#!/bin/bash
set -ex

#Reconfig spctl
${GOPATH}/bin/spctl config pms-endpoint http://localhost:6733/policy-mgmt/v1/
source ${GOPATH}/src/gitlab-odx.oracledx.com/wcai/speedle/setTestEnv.sh

startPMS file --config-file pkg/svcs/pmsrest/config_file.json
go test ${TEST_OPTS} gitlab-odx.oracledx.com/wcai/speedle/cmd/spctl/command -run=TestMats

