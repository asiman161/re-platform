.PHONY: build-front-dev
build-front-dev:
	scripts/front.sh

.PHONY: reload-front-dev
reload-front-dev:
	kubectl apply -f kube-front.dev.yaml
	kubectl delete deploy re-platform-front-dev
	kubectl apply -f kube-front.dev.yaml

.PHONY: build-back-dev
build-back-dev:
	scripts/back.sh

.PHONY: reload-back-dev
reload-back-dev:
	kubectl apply -f kube-back.dev.yaml
	kubectl delete deploy re-platform-back-dev
	kubectl apply -f kube-back.dev.yaml

.PHONY: kube-front-dev
kube-front-dev: build-front-dev reload-front-dev

.PHONY: kube-back-dev
kube-back-dev: build-back-dev reload-back-dev

.PHONY: start-dev
start-dev: reload-front-dev reload-back-dev
