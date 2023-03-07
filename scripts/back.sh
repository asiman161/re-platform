#!/bin/bash

eval $(minikube docker-env)
docker build -t re-platform-back-dev ./back -f ./back/Dockerfile
