#!/bin/bash

eval $(minikube docker-env)
docker build -t asiman61/re-platform-front-dev ./front -f ./front/Dockerfile.dev
