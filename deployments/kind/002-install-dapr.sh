#!/bin/sh

helm repo add dapr https://daprio.azurecr.io/helm/v1/repo
helm repo update

kubectl create namespace dapr-system
helm install dapr dapr/dapr --namespace dapr-system

echo "Run 'kubectl get pods -n dapr-system -w' to verify instalation"