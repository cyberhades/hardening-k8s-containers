#!/bin/bash

# Start minikube
minikube start --extra-config=apiserver.authorization-mode=RBAC --network-plugin=cni --memory=4096 --driver=virtualbox --kubernetes-version v1.15.0 

# Setup Cilium
minikube ssh -- sudo mount bpffs -t bpf /sys/fs/bpf
kubectl create -f https://raw.githubusercontent.com/cilium/cilium/1.8.1/install/kubernetes/quick-install.yaml

# Switch docker context
eval $(minikube docker-env)

# Build images
docker build -t notes:v1 .
docker build -t notes:v2 -f Dockerfile.v2 .

# Deploy the applications
kubectl apply -f kube/v1-deploy.yaml -f kube/v1-service.yaml -f kube/v2-deploy.yaml -f kube/v2-service.yaml

