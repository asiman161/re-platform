include scripts/scripts.mk

.PHONY: build-dev
build-dev: kube-front-dev

.PHONY: tunnel
tunnel:
	minikube tunnel
