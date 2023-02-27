.PHONY: build-front-dev
build-front-dev:
	scripts/front.sh

.PHONY: reload-front-dev
reload-front-dev:
	kubectl apply -f kube-front.dev.yaml
	kubectl delete deploy re-platform-front-dev
	kubectl apply -f kube-front.dev.yaml

.PHONY: kube-front-dev
kube-front-dev: build-front-dev reload-front-dev
