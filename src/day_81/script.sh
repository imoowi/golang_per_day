#!/bin/bash

# 创建 kind 集群
# 使用 kind-config.yaml 配置文件创建名为 codee-jun 的集群
kind create cluster --config kind-config.yaml --name coding

# 检查节点是否 Ready
# 验证集群是否成功创建
kubectl get nodes
# 安装基础服务
kubectl apply -f mysql.yaml 
kubectl apply -f redis.yaml

# 设置 LINKERD2_VERSION 环境变量，指定要安装的 Linkerd 版本
# 如果未设置，则会安装最新的 edge 版本
export LINKERD2_VERSION=edge-25.10.7

# 下载并运行 Linkerd 安装脚本
# 使用 TLSv1.2 协议确保安全下载
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install-edge | sh

# 将 Linkerd CLI 添加到 PATH 环境变量中
export PATH=$HOME/.linkerd2/bin:$PATH

# 检查 Linkerd 版本
# 验证 CLI 是否正确安装
linkerd version

# 安装 Gateway API 标准资源
# 使用 server-side 应用模式安装 v1.4.0 版本
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.0/standard-install.yaml

# 运行 Linkerd 预检查
# 验证集群是否满足安装 Linkerd 的条件
linkerd check --pre

# 安装 Linkerd CRDs
# 先安装自定义资源定义
linkerd install --crds | kubectl apply -f -

# 安装 Linkerd 控制平面
# 安装核心组件
linkerd install | kubectl apply -f -

# 检查 Linkerd 部署状态
# 验证控制平面组件是否成功部署
kubectl -n linkerd get deploy

# 运行 Linkerd 健康检查
# 验证整个 Linkerd 安装是否健康
linkerd check


# 安装 Linkerd 可视化组件
# 安装 dashboard 等监控工具
linkerd viz install | kubectl apply -f -

# 启动 Linkerd 仪表盘
# 打开 Web 界面查看服务网格状态
# linkerd viz dashboard

# 安装 APISIX API 网关
# 参考 https://apisix.apache.org/docs/apisix/getting-started/README/

# 添加 APISIX Helm 仓库
helm repo add apisix https://charts.apiseven.com
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# We use Apisix 3.0 in this example. If you're using Apisix v2.x, please set to v2
ADMIN_API_VERSION=v3

helm install apisix apisix/apisix \
  --set service.type=NodePort \
  --set gateway.http.nodePort=30080 \
  --set ingress-controller.enabled=true \
  --set ingress-controller.config.apisix.serviceNamespace=apisix \
  --set ingress-controller.config.apisix.adminAPIVersion=$ADMIN_API_VERSION
# 检查 APISIX 部署状态
# 验证 APISIX 组件是否成功部署
kubectl get pods 
# 安装 APISIX 仪表盘
helm install apisix-dashboard apisix/apisix-dashboard 

kubectl apply -f app.configmap.yaml
kubectl apply -f app.yaml
# 验证 user-service 是否成功部署
kubectl get svc  user-service

kubectl port-forward svc/user-service 8000:8000 

kubectl port-forward svc/apisix-dashboard 9080:80 
kubectl apply -f apisix-route.yaml
# 取apisix-admin的 apikey
kubectl get cm apisix  -o yaml | grep key
# 使用 APISIX Admin API 来直接检查当前的路由配置
kubectl run curl  --image=curlimages/curl --restart=Never -it --rm -- sh
$ curl -s -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" http://apisix-admin:9180/apisix/admin/routes
$ curl -s -X POST -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" -H "Content-Type: application/json" http://apisix-admin:9180/apisix/admin/routes -d '{"uri":"/api/user/*","upstream":{"type":"roundrobin","nodes":{"user-service.apisix.svc.cluster.local:8000":1}}}'

kubectl edit svc ingress--gateway 
# 修改 nodePort: 30080 ，保存退出，k8s自动重启服务
