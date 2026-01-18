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
helm install jaeger jaegertracing/jaeger --namespace monitoring

kubectl apply -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/main/deploy/crds/jaegertracing.io_jaegers_crd.yaml
# 这里直接用一个简单的 deployment 运行 jaeger
kubectl run jaeger --image=jaegertracing/all-in-one:latest -n monitoring --expose --port=16686
kubectl expose pod jaeger -n monitoring --type=ClusterIP --port=4317 --target-port=4317 --name=jaeger-collector

