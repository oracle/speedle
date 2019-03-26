SHELL := /bin/bash
MAKE = make


init:
	mkdir -p workdir
	mkdir -p workdir/speedlegrpcadapter

init_istio:
	rm -rf ${GOPATH}/src/istio.io/istio
	git clone https://github.com/istio/istio  ${GOPATH}/src/istio.io/istio --single-branch --branch release-1.0

build-adapter:

#Copy the speedle adapter to istio and build it:
	cd ${MIXER_REPO}/adapter
	cp -r adapter/speedle ${MIXER_REPO}/adapter/
	cd ${MIXER_REPO}/adapter/speedle ; \
	go generate ./... ; \
	go build ./...
	echo "************* go build done ************ "

#Plug speedle adapter to Mixer
	echo 'speedle: "istio.io/istio/mixer/adapter/speedle"' >> ${MIXER_REPO}/adapter/inventory.yaml

	cd ${MIXER_REPO} ; \
	go generate ${MIXER_REPO}/adapter/doc.go

	cd ${ISTIO}/istio && make mixs
	cd ${ISTIO}/istio && make HUB=${BUILD_DOCKER_HUB} TAG=${BUILD_DOCKER_TAG} push.docker.mixer
#	${GOPATH}/out/linux_amd64/release/mixs crd adapter

build-grpc-adapter:

#Copy the speedle adapter from speedle-pub and build it:
	cd ${MIXER_REPO}/adapter
	cp -r adapter/speedle ${MIXER_REPO}/adapter/
	cp -r adapter/speedlegrpcadapter ${MIXER_REPO}/adapter/
	cd ${MIXER_REPO}/adapter/speedlegrpcadapter ; \
	go generate ./...

	echo "  go generate done"
	cd ${MIXER_REPO}/adapter/speedlegrpcadapter/cmd && $(MAKE)
	echo "************* speedle grpc adapter build done ************ "


install-bookinfo:
	cd ${GOPATH}/src/istio.io/istio ; kubectl apply -f <(istioctl kube-inject -f samples/bookinfo/platform/kube/bookinfo.yaml)
	cd ${GOPATH}/src/istio.io/istio ; kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml

uninstall-bookinfo:
	cd ${GOPATH}/src/istio.io/istio ; samples/bookinfo/platform/kube/cleanup.sh


_run-install-grpc-adapter:

	kubectl apply -f ${MIXER_REPO}/template/authorization/template.yaml
	kubectl apply -f ${MIXER_REPO}/adapter/speedlegrpcadapter/speedlegrpcadapter.yaml
	kubectl apply -f ./workdir/speedlegrpcadapter/deployment.yaml
	istioctl create -f ./workdir/speedlegrpcadapter/operator_cfg.yaml


_init-grpc-adapter-files: init
	cp ./adapter/speedlegrpcadapter/deployment.yaml.template ./workdir/speedlegrpcadapter/deployment.yaml
	cp ./adapter/speedlegrpcadapter/operator_cfg.yaml.template ./workdir/speedlegrpcadapter/operator_cfg.yaml

install-speedle-grpc-adapter: _prepare-speedle-grpc-adater-files _run-install-grpc-adapter


_prepare-speedle-grpc-adater-files: _init-grpc-adapter-files

	sed -i 's#<BUILD_OOP_DOCKER_HUB>#${BUILD_OOP_DOCKER_HUB}#g' ./workdir/speedlegrpcadapter/deployment.yaml
	sed -i 's#<BUILD_OOP_DOCKER_TAG>#${BUILD_OOP_DOCKER_TAG}#g' ./workdir/speedlegrpcadapter/deployment.yaml
	sed -i 's#<OOP_ADAPTER_SERVICE_PORT>#${OOP_ADAPTER_SERVICE_PORT}#g' ./workdir/speedlegrpcadapter/deployment.yaml

	sed -i 's#<ADS_ENDPOINT>#${SPEEDLE_ADS_ENDPOINT}#g' ./workdir/speedlegrpcadapter/operator_cfg.yaml
	sed -i 's#<SERVICE_NAME>#${SPEEDLE_SERVICE_NAME}#g' ./workdir/speedlegrpcadapter/operator_cfg.yaml
	sed -i 's#<OOP_ADAPTER_SERVICE_PORT>#${OOP_ADAPTER_SERVICE_PORT}#g' ./workdir/speedlegrpcadapter/operator_cfg.yaml


uninstall-speedle-grpc-adapter: _prepare-speedle-grpc-adater-files _run-uninstall-grpc-adapter


_run-uninstall-grpc-adapter:
	istioctl delete -f ./workdir/speedlegrpcadapter/operator_cfg.yaml
	kubectl delete -f ./workdir/speedlegrpcadapter/deployment.yaml
	kubectl delete -f ${MIXER_REPO}/adapter/speedlegrpcadapter/speedlegrpcadapter.yaml
	kubectl delete -f ${MIXER_REPO}/template/authorization/template.yaml

clear:
	rm -rf ./workdir
