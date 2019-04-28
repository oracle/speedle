.PHONY: ci
retryCount=5
waitInteval=2
ci:
	docker build $$docker_build_options -t phx.ocir.io/kubernetessecurity/speedle/ci ci/docker/ci
	docker push phx.ocir.io/kubernetessecurity/speedle/ci

public:
	hugo

docker-login:
	docker login -u "$$DOCKER_LOGIN_USER" -p "$$DOCKER_LOGIN_PASSWORD" $$DOCKER_LOGIN_SERVER

build/bin/spctl:
	docker run --rm -v "${PWD}/build":/go golang:1.10 go get github.com/oracle/speedle/cmd/spctl

build/bin/speedle-pms:
	docker run --rm -v "${PWD}/build":/go golang:1.10 go get github.com/oracle/speedle/cmd/speedle-pms

build/bin/speedle-ads:
	docker run --rm -v "${PWD}/build":/go golang:1.10 go get github.com/oracle/speedle/cmd/speedle-ads

pms: build/bin/speedle-pms
	cp build/bin/speedle-pms ci/docker/pms
	docker build -t phx.ocir.io/kubernetessecurity/speedle/pms ci/docker/pms
	docker push phx.ocir.io/kubernetessecurity/speedle/pms
	rm ci/docker/pms/speedle-pms

ads: build/bin/speedle-ads
	cp build/bin/speedle-ads ci/docker/ads
	docker build -t phx.ocir.io/kubernetessecurity/speedle/ads ci/docker/ads
	docker push phx.ocir.io/kubernetessecurity/speedle/ads
	rm ci/docker/ads/speedle-ads

build-playground:
	@+make -C playground

build-backend: public build-playground
	rm -rf backend/public
	cp -R public/ backend/public
	cp -R playground/build/ backend/public/playground

site: build-backend build/bin/spctl
	rm -rf ci/docker/playground/app; mkdir ci/docker/playground/app
	cp -R backend/public/ ci/docker/playground/app/public

	cp backend/package.json ci/docker/playground
	cp backend/app.js ci/docker/playground/app
	mkdir ci/docker/playground/app/views
	cp -R backend/routes ci/docker/playground/app
	cp -R backend/bin ci/docker/playground/app

	cp build/bin/spctl ci/docker/playground

	docker build -t phx.ocir.io/kubernetessecurity/speedle/site-playground ci/docker/playground
	docker push phx.ocir.io/kubernetessecurity/speedle/site-playground

	rm -rf ci/docker/playground/app
	rm ci/docker/playground/package.json
	rm ci/docker/playground/spctl

docker-build-apache:
	docker build -t phx.ocir.io/kubernetessecurity/speedle/apache ci/docker/apache
	docker push phx.ocir.io/kubernetessecurity/speedle/apache

docker-build-site:
	docker build -t phx.ocir.io/kubernetessecurity/speedle/site ci/docker/site
	docker push phx.ocir.io/kubernetessecurity/speedle/site
	rm -rf ci/docker/site/public

publish: site

deploy: publish
	pods=(`kubectl get pod -l app=speedlesite --no-headers=true -n default -o jsonpath="{..metadata.name}"`); \
	echo "Deleting pod: $${pods[0]}"; \
	kubectl delete po $${pods[0]} -n default && \
	for ((i = 1; i <= $(retryCount); i++)) \
	do \
		pods1=(`kubectl get pod -l app=speedlesite --field-selector status.phase=Running --no-headers=true -n default -o jsonpath="{..metadata.name}"`); \
		[ "$${#pods1[@]}" -eq "2" ] && break; \
		echo "Waiting for pod ready: $${i}"; \
		sleep $(waitInteval); \
	done; \
	echo "Deleting pod: pod $${pods[1]}"; \
	kubectl delete po $${pods[1]} -n default

clean:
	rm -rf build
	rm -rf public