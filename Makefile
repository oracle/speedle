.PHONY: all test

gopath := $(shell go env GOPATH)
gitCommit := $(shell git rev-parse --short HEAD)
# go version output is "go version go1.11.2 linux/amd64"
goVersion := $(word 3,$(shell go version))
goLDFlags := -ldflags "-X main.gitCommit=${gitCommit} -X main.productVersion=0.1 -X main.goVersion=${goVersion}"

pmsImageRepo := speedle-pms
pmsImageTag := v0.1
adsImageRepo := speedle-ads
adsImageTag := v0.1

all: build

build: buildPms buildAds buildSpctl

buildPms:
	go build ${goLDFlags} -o ${gopath}/bin/speedle-pms github.com/oracle/speedle/cmd/speedle-pms

buildAds:
	go build ${goLDFlags} -o ${gopath}/bin/speedle-ads github.com/oracle/speedle/cmd/speedle-ads

buildSpctl:
	go build ${goLDFlags} -o ${gopath}/bin/spctl  github.com/oracle/speedle/cmd/spctl

image: imagePms imageAds

imagePms:
	cp ${gopath}/bin/speedle-pms deployment/docker/speedle-pms/.
	docker build -t ${pmsImageRepo}:${pmsImageTag} --rm --no-cache deployment/docker/speedle-pms
	rm deployment/docker/speedle-pms/speedle-pms

imageAds:
	cp ${gopath}/bin/speedle-ads deployment/docker/speedle-ads/.
	docker build -t ${adsImageRepo}:${adsImageTag} --rm --no-cache deployment/docker/speedle-ads
	rm deployment/docker/speedle-ads/speedle-ads

test: testAll

testAll: speedleUnitTests testSpeedleRest testSpeedleGRpc testSpctl testSpeedleRestADSCheck testSpeedleGRpcADSCheck testSpeedleTls

speedleUnitTests:
	go test ${TEST_OPTS} github.com/oracle/speedle/pkg/cfg 
	go test ${TEST_OPTS} github.com/oracle/speedle/pkg/eval 
	go test ${TEST_OPTS} github.com/oracle/speedle/pkg/store/file
	go test ${TEST_OPTS} github.com/oracle/speedle/pkg/store/etcd
	go test ${TEST_OPTS} github.com/oracle/speedle/cmd/spctl/pdl
	go test ${TEST_OPTS} github.com/oracle/speedle/pkg/suid
	go test ${TEST_OPTS} github.com/oracle/speedle/pkg/assertion
	go clean -testcache
	STORE_TYPE=etcd go test ${TEST_OPTS} github.com/oracle/speedle/pkg/eval

testSpeedleRest:
	pkg/svcs/pmsrest/run_file_test.sh
	pkg/svcs/pmsrest/run_etcd_test.sh

testSpeedleGRpc:
	pkg/svcs/pmsgrpc/run_file_test.sh
	pkg/svcs/pmsgrpc/run_etcd_test.sh

testSpeedleRestADSCheck:
	pkg/svcs/adsrest/run_file_test.sh
	pkg/svcs/adsrest/run_etcd_test.sh

testSpeedleGRpcADSCheck:
	pkg/svcs/adsgrpc/run_file_test.sh
	pkg/svcs/adsgrpc/run_etcd_test.sh

testSpctl:
	cmd/spctl/command/run_file_test.sh
	cmd/spctl/command/run_etcd_test.sh

testSpeedleTls:
	pkg/svcs/pmsrest/tls_test.sh
	pkg/svcs/pmsrest/tls_test-force-client-cert.sh
clean:
	rm -rf ${gopath}/pkg/linux_amd64/github.com/oracle/speedle
	rm -f ${gopath}/bin/speedle-pms
	rm -f ${gopath}/bin/speedle-ads
	rm -f ${gopath}/bin/spctl
