#!/bin/bash

eval $(minikube docker-env)
docker build -t re-platform-peerserver-dev ./peerserver -f ./peerserver/Dockerfile.dev
