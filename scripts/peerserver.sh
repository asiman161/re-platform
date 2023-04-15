#!/bin/bash

eval $(minikube docker-env)
docker build -t asiman61/re-platform-peerserver-dev ./peerserver -f ./peerserver/Dockerfile.dev
