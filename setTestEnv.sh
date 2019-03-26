#/bin/sh

export PMS_ENDPOINT=http://127.0.0.1:6733
export ADS_ENDPOINT=http://127.0.0.1:6734
export PMS_ADMIN_TOKEN=
export ADS_ADMIN_TOKEN=
export SP_APP_NAME=spctl

echo "----------------------------------"
echo "PMS_ENDPOINT=${PMS_ENDPOINT}"
echo "SP_APP_NAME=${SP_APP_NAME}"
echo "----------------------------------"


echo "===================Set Local Test Environment for Speedle===================="
exit_trap_command=""
function cleanup {
    eval "$exit_trap_command"
}
trap cleanup EXIT

function add_exit_trap {
    local to_add=$1
    if [[ -z "$exit_trap_command" ]]
    then
        exit_trap_command="$to_add"
    else
        exit_trap_command="$to_add; $exit_trap_command"
    fi
}

# helper methods
function log::showGoodMessage() {
    echo -e "\033[32m$@\033[0m"
}

function log::showBadMessage() {
    echo -e "\033[31m$@\033[0m"
}

function log::showYellowMessage() {
    echo -e "\033[33m$@\033[0m"
}

function log::showFailedMessage() {
    log::showBadMessage "Failed to $1. Execution time:$[ `date +%s` - $2 ]s"
}

function log::showFailedMessage() {
    log::showGoodMessage "$1 OK. Execution time:$[ `date +%s` - $2 ]s"
}

function log::showTitle() {
    log::showYellowMessage "\n=======================$1======================="
}

function log::log() {
    echo -e $(date): "$@"
}

function tryCurl() {
    log::log request: `echo "$@" | sed "s/-H .*Bearer [^ ]* //" | sed "s/-u [^ ]* //"`
    http_response=$(curl --connect-timeout 600 -m 600 --retry-max-time 2 --silent --write-out "status:%{http_code}" "$@") || (echo $http_response && return 0)
    http_body=$(echo $http_response | sed -e 's/status\:.*//g')
    http_status=$(echo $http_response | tr -d '\n' | sed -e 's/.*status://')
    log::log response:
    log::log " code: $http_status"
    log::log " body: $http_body"
}

function util::retry {
    echo "---------------util::retry"
    RETRY_COUNT=60
    log::showYellowMessage "retry $@..."
    # set -x
    sleep $1
    shift
    for (( i = 0; i < $RETRY_COUNT; i++))
        do
            eval "$@" && return 0 ||  echo "wait..." && sleep 1
    done
    log::showBadMessage "timeout! retry last time..."
    # set -x
    eval "$@" || return 1
}

function waitService  {
    echo "Wait for service $1 to be ready..."
    # util::retry echo "8"
    util::retry 3 "cat < /dev/null > /dev/tcp/localhost/$2"
    if [  $? -ne 0 ]; then
        echo "service $1 is not ready"
    else
        echo "service $1 is ready"
    fi
    # sleep 1000000
}

shopt -s expand_aliases
if [ `uname -s` == "Darwin" ] ; then
  alias pms=${GOPATH}/bin/speedle-pms-mac
  alias ads=${GOPATH}/bin/speedle-ads-mac
else
  alias pms=${GOPATH}/bin/speedle-pms
  alias ads=${GOPATH}/bin/speedle-ads
fi

function ensureTestDir() {
    if [ "$1" == "file" ];then
        temp_policy_file=/tmp/speedle-test-file-store.json
        rm -f ${temp_policy_file}
        echo "{}" > ${temp_policy_file}
        add_exit_trap "rm -f ${temp_policy_file}"
    else
        rm -rf ./speedle.etcd
        add_exit_trap "rm -rf ./speedle.etcd"
    fi
}

function startPMS() {
    ensureTestDir $1
    shift

    pms $@ &
    serverPID=$!
    add_exit_trap "kill ${serverPID}"
    waitService PMS 6733 || exit 1
}

function startADS() {
    ads $@ &
    serverPID=$!
    add_exit_trap "kill ${serverPID}"
    waitService ADS 6734 || exit 1
}

function failTest {
    log::showBadMessage "Failed"
    exit 1
}