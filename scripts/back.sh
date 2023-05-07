#!/bin/bash

eval $(minikube docker-env)
docker build -t asiman61/re-platform-back-dev ./goback -f ./goback/Dockerfile
