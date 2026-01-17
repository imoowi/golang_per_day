#!bin/bash
kind delete cluster --name golang-per-day
kind create cluster --config kind-config.yaml --name golang-per-day

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
kubectl create namespace monitoring
helm install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring
kubectl get pods -n monitoring

kubectl create namespace codee-jun

helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm install jaeger jaegertracing/jaeger

