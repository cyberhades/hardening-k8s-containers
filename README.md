## Intentionally Vulnerable Note Taking Application

This repository contains a very simple and vulnerable application used to demonstrate the differences between deploying this application with the default settings to a Kubernetes cluster, versus adding some extra security settings that **WILL NOT** fix the application's vulnerabilities, but will restrain an attacker from doing bad things.

To try this yourself, you need access to a Kubernetes cluster, because this application is highly vulnerable **you do not want** to deploy this in a real cluster unless such cluster is properly isolated.

## Minikube
If you want to use minikube, you just can run the provided `init.sh` script or go step by step:

#### Start minikube with at least cni plugin

    minikube start --extra-config=apiserver.authorization-mode=RBAC --network-plugin=cni --memory=4096 --vm-driver=virtualbox

#### Install Cilium cni plugin

    minikube ssh -- sudo mount bpffs -t bpf /sys/fs/bpf
    kubectl create -f https://raw.githubusercontent.com/cilium/cilium/1.6.5/install/kubernetes/quick-install.yaml

#### Change your docker context to the one in minikube, otherwise you'll need to push the docker images to a registry

    eval $(minikube docker-env)
    
## Build your "default" docker image

    docker build -t notes:v1 .
    
## Build the second docker image with some hardening configuration settings

    docker build -t notes:v2 -f Dockerfile.v2 .
    
## Deploy to minikube

    kubectl apply -f kube/v1-deploy.yaml -f kube/v1-service.yaml -f kube/v2-deploy.yaml -f kube/v2-service.yaml

Have fun!


