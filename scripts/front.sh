#!/bin/bash

eval $(minikube docker-env)
docker build -t re-platform-front-dev ./front -f ./front/Dockerfile.dev
