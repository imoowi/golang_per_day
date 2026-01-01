#!/bin/bash
#删除
kubectl delete -f deploy.yaml
kubectl delete -f server.yaml
kubectl delete -f storageclass.yaml
kubectl delete -f pv-0.yaml
kubectl delete -f pv-1.yaml
kubectl delete -f pv-2.yaml
kubectl delete -f configmap.yaml
#部署
kubectl apply -f configmap.yaml
kubectl apply -f pv-0.yaml
kubectl apply -f pv-1.yaml
kubectl apply -f pv-2.yaml
kubectl apply -f storageclass.yaml
kubectl apply -f server.yaml
kubectl apply -f deploy.yaml
kubectl get pods -l app=nats
