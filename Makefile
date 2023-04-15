include scripts/scripts.mk

.PHONY: build
build: kube-front-dev kube-back-dev

.PHONY: tunnel
tunnel:
	minikube tunnel

.PHONY: up-dev
up-dev: start-dev

# starts only infrastructure and peerserver
.PHONY: up-simple-dev
up-simple-dev: start-simple-dev tunnel

.PHONY: down-dev
down-dev: shutdown-dev
