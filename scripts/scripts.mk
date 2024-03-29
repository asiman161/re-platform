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

.PHONY: build-peerserver-dev
build-peerserver-dev:
	scripts/peerserver.sh

.PHONY: reload-peerserver-dev
reload-peerserver-dev:
	kubectl apply -f kube-peerserver.dev.yaml
	kubectl delete deploy re-platform-peerserver-dev
	kubectl apply -f kube-peerserver.dev.yaml

.PHONY: infra-up-dev
infra-up-dev:
	kubectl apply -f kube-pvc.dev.yaml
	kubectl apply -f kube-infra.dev.yaml
	reload-peerserver-dev

.PHONY: peerserver-up-dev
peerserver-up-dev:
	kubectl apply -f kube-peerserver.dev.yaml


.PHONY: kube-front-dev
kube-front-dev: build-front-dev reload-front-dev

.PHONY: kube-back-dev
kube-back-dev: build-back-dev reload-back-dev

.PHONY: kube-peerserver-dev
kube-back-dev: build-peerserver-dev reload-peerserver-dev

.PHONY: start-dev
up-dev: reload-front-dev reload-back-dev reload-peerserver-dev infra-up-dev

.PHONY: start-simple-dev
start-simple-dev: infra-up-dev peerserver-up-dev


.PHONY: shutdown-dev
shutdown-dev:
	kubectl apply -f kube-front.dev.yaml
	kubectl apply -f kube-back.dev.yaml
	kubectl apply -f kube-infra.dev.yaml
	kubectl apply -f kube-peerserver.dev.yaml
	kubectl delete deploy re-platform-front-dev
	kubectl delete deploy re-platform-back-dev
	kubectl delete deploy re-platform-pg-dev
	kubectl delete deploy re-platform-redis-dev
	kubectl delete deploy re-platform-peerserver-dev

