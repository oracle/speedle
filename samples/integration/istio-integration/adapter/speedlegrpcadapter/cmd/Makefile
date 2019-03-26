.PHONY: deploy docker

all: speedlegrpcadapter docker

speedlegrpcadapter: main.go
	../../../../bin/gobuild.sh speedlegrpcadapter istio.io/istio/mixer/adapter/speedlegrpcadapter/cmd .
docker: speedlegrpcadapter
	docker build . -t ${BUILD_OOP_DOCKER_HUB}:${BUILD_OOP_DOCKER_TAG}
