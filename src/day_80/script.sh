#!/bin/bash

# 创建 kind 集群
# 使用 kind-config.yaml 配置文件创建名为 codee-jun 的集群
kind create cluster --config kind-config.yaml --name codee-jun

# 检查节点是否 Ready
# 验证集群是否成功创建
kubectl get nodes

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

# 安装示例应用 emojivoto
# 部署一个示例应用来测试服务网格
kubectl apply -f https://run.linkerd.io/emojivoto.yml

# 为 emojivoto 应用注入 Linkerd 边车
# 将 Linkerd 代理注入到应用的 Pod 中
kubectl -n emojivoto get deploy -o yaml \
  | linkerd inject - \
  | kubectl apply -f -

# 检查 emojivoto 应用的 Pod 状态
# 验证应用是否成功部署并且边车已注入
kubectl -n emojivoto get pods

# 安装 Linkerd 可视化组件
# 安装 dashboard 等监控工具
linkerd viz install | kubectl apply -f -

# 启动 Linkerd 仪表盘
# 打开 Web 界面查看服务网格状态
linkerd viz dashboard

# 安装 APISIX API 网关
# 参考 https://apisix.apache.org/docs/apisix/getting-started/README/

# 添加 APISIX Helm 仓库
helm repo add apisix https://charts.apiseven.com

# 更新 Helm 仓库
helm repo update

# 创建 apisix 命名空间
kubectl create ns apisix

# 检查 apisix 命名空间是否已正确创建
kubectl get ns apisix -o yaml | grep linkerd

# 使用 Helm 安装 APISIX
# 配置：
# - 设置网关类型为 NodePort
# - 设置 HTTP 节点端口为 30080
# - 禁用 TLS
# - 禁用 etcd 持久化
# - 启用 ingress 控制器
# - 禁用 Gateway API
# helm install apisix apisix/apisix -n apisix --create-namespace \
#   --set gateway.type=NodePort \
#   --set gateway.http.nodePort=30080 \
#   --set gateway.tls.enabled=false \
#   --set etcd.persistence.enabled=false \
#   --set ingress-controller.enabled=true \
#   --set ingress-controller.config.kubernetes.enableGatewayAPI=false
helm install apisix apisix/apisix --create-namespace --namespace apisix
# 检查 APISIX 部署状态
# 验证 APISIX 组件是否成功部署
kubectl get pods -n apisix

kubectl apply -f app.yaml
kubectl apply -f apisix-route.yaml

kubectl patch svc apisix-gateway -n apisix \
  -p '{"spec":{"ports":[{"port":80,"targetPort":80,"nodePort":30080}]}}'

