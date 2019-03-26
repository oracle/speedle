#!/bin/bash

shell_dir=$(dirname $0)
# set -x
source ${GOPATH}/src/github.com/oracle/speedle/setTestEnv.sh
${GOPATH}/bin/spctl config pms-endpoint https://localhost:6733/policy-mgmt/v1/

log::showTitle "Start pms service" && \
startPMS ${shell_dir}/config_file.json --cert ${shell_dir}/tls/server.crt --key=${shell_dir}/tls/server.key --client-cert ${shell_dir}/tls/server-ca.crt --insecure false && \ 
log::showTitle "Test skipverify=true without any cert settings" && \
$GOPATH/bin/spctl --skipverify=true  get service --all  && \
log::showGoodMessage "OK" && \
log::showTitle "Test skipverify=true + cert && key" && \
$GOPATH/bin/spctl --cert=${shell_dir}/tls/client.crt --key=${shell_dir}/tls/client.key --skipverify=true  get service --all && \
log::showGoodMessage "OK" && \
log::showTitle "Test skipverify=false + cert && key && cacert" && \
$GOPATH/bin/spctl --cert=${shell_dir}/tls/client.crt --key=${shell_dir}/tls/client.key --cacert=${shell_dir}/tls/client-ca.crt  --skipverify=false  get service --all  && \
log::showGoodMessage "OK"  && \
log::showTitle "Test skipverify=false + cacert" && \
$GOPATH/bin/spctl --cacert=${shell_dir}/tls/client-ca.crt --skipverify=false  get service --all  && \
log::showGoodMessage "OK" && \
log::showTitle "Test bad certificate" && \
${GOPATH}/bin/spctl config pms-endpoint https://127.0.0.1:6733/policy-mgmt/v1/
$GOPATH/bin/spctl --cacert=${shell_dir}/tls/client-ca.crt --skipverify=false  get service --all && log::showBadMessage "Failed" || \
log::showGoodMessage "OK" && \
exit 0

failTest
